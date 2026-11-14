package mailer

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"

	"abberant/backend/internal/config"
)

type Mailer interface {
	Send(to, subject, body string) error
}

func New(cfg config.Config) Mailer {
	if cfg.SMTPEnabled() {
		return &smtpMailer{cfg: cfg}
	}

	return &logMailer{}
}

type smtpMailer struct {
	cfg config.Config
}

func (m *smtpMailer) Send(to, subject, body string) error {
	addr := net.JoinHostPort(m.cfg.SMTPHost, m.cfg.SMTPPort)

	from := m.cfg.SMTPFrom
	message := buildMessage(from, to, subject, body)

	client, err := dialSMTP(addr, m.cfg.SMTPPort, m.cfg.SMTPHost)
	if err != nil {
		return err
	}

	defer client.Close()

	if m.cfg.SMTPPort != implicitTLSPort {
		if err := m.startTLS(client); err != nil {
			return err
		}
	}

	// Local SMTP catchers (MailPit, MailDev, postfix on localhost) usually
	// do not require authentication and do not advertise the AUTH
	// extension, so authenticate only when credentials are configured.
	if m.cfg.SMTPUser != "" || m.cfg.SMTPPass != "" {
		auth := smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPass, m.cfg.SMTPHost)

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}

	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return client.Quit()
}

// implicitTLSPort marks the port that speaks TLS from the first byte.
// It is a variable only so tests can point it at an ephemeral port.
var implicitTLSPort = "465"

// dialSMTP connects to the server. Port 465 speaks implicit TLS from the
// first byte, so it needs a TLS handshake instead of STARTTLS.
func dialSMTP(addr, port, host string) (*smtp.Client, error) {
	if port != implicitTLSPort {
		conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			return nil, fmt.Errorf("dial smtp: %w", err)
		}

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("smtp handshake: %w", err)
		}

		return client, nil
	}

	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return nil, fmt.Errorf("dial smtp over TLS: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("smtp handshake: %w", err)
	}

	return client, nil
}

func (m *smtpMailer) startTLS(client *smtp.Client) error {
	if !m.cfg.SMTPStartTL {
		return nil
	}

	if ok, _ := client.Extension("STARTTLS"); !ok {
		return nil
	}

	if err := client.StartTLS(&tls.Config{ServerName: m.cfg.SMTPHost}); err != nil {
		return fmt.Errorf("smtp STARTTLS: %w", err)
	}

	return nil
}

func buildMessage(from, to, subject, body string) []byte {
	var builder strings.Builder

	builder.WriteString("From: Abberant <" + from + ">\r\n")
	builder.WriteString("To: <" + to + ">\r\n")
	builder.WriteString("Subject: " + subject + "\r\n")
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	builder.WriteString("\r\n")
	builder.WriteString(body)

	return []byte(builder.String())
}

type logMailer struct{}

func (m *logMailer) Send(to, subject, body string) error {
	log.Printf(
		"[mailer] SMTP is not configured, message shown here instead of being sent:\n"+
			"  To: %s\n  Subject: %s\n%s", to, subject, body)

	return nil
}
