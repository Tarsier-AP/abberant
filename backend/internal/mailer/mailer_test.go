package mailer

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"abberant/backend/internal/config"
)

// The STARTTLS test relies on SSL_CERT_FILE to make the standard x509
// verifier trust the self-signed test certificate. Go honours that
// variable on unix systems only, so the test is skipped elsewhere.
func TestSendDeliversOverSTARTTLS(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("SSL_CERT_FILE override works only on linux")
	}

	serverTLS, certFile := selfSignedTLSConfig(t)
	t.Setenv("SSL_CERT_FILE", certFile)

	srv := newFakeServer(t, serverTLS, true)

	cfg := testConfig(srv, true)
	mailer := &smtpMailer{cfg: cfg}

	assertFullDelivery(t, srv, mailer)
}

func TestSendDeliversWithoutSTARTTLS(t *testing.T) {
	srv := newFakeServer(t, nil, false)

	cfg := testConfig(srv, false)
	mailer := &smtpMailer{cfg: cfg}

	assertFullDelivery(t, srv, mailer)
}

func TestSendProceedsWhenSTARTTLSIsNotOffered(t *testing.T) {
	srv := newFakeServer(t, nil, false)

	// STARTTLS is enabled in the config, but the fake server does not
	// advertise the extension: delivery must continue without TLS.
	cfg := testConfig(srv, true)
	mailer := &smtpMailer{cfg: cfg}

	assertFullDelivery(t, srv, mailer)
}

func TestSendUsesImplicitTLSOnPort465(t *testing.T) {
	// Binding a real port 465 needs privileges in some environments, so
	// point the mailer's implicit-TLS branch at a free ephemeral port.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot allocate a port for the test: %v", err)
	}
	defer listener.Close()

	_, port, _ := net.SplitHostPort(listener.Addr().String())
	original := implicitTLSPort
	implicitTLSPort = port
	defer func() { implicitTLSPort = original }()

	firstByte := make(chan byte, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		conn.SetDeadline(time.Now().Add(5 * time.Second))

		buf := make([]byte, 1)
		if _, err := conn.Read(buf); err == nil {
			firstByte <- buf[0]
		}
	}()

	cfg := config.Config{
		SMTPHost:    "localhost",
		SMTPPort:    port,
		SMTPUser:    "user",
		SMTPPass:    "pass",
		SMTPFrom:    "abberant@localhost",
		SMTPStartTL: false,
	}

	err = (&smtpMailer{cfg: cfg}).Send("user@example.com", "Test subject", "body")
	if err == nil {
		t.Fatal("expected an error: the fake server does not speak TLS")
	}

	select {
	case b := <-firstByte:
		if b != 0x16 {
			t.Fatalf("first byte is %#x, want TLS ClientHello (0x16)", b)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no connection attempt was observed")
	}
}

func TestBuildMessage(t *testing.T) {
	got := string(buildMessage("from@example.com", "to@example.com", "Привет", "line one\r\nline two"))

	want := "From: Abberant <from@example.com>\r\n" +
		"To: <to@example.com>\r\n" +
		"Subject: Привет\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"line one\r\nline two"

	if got != want {
		t.Fatalf("unexpected message:\n got: %q\nwant: %q", got, want)
	}
}

func TestNewMailerSelection(t *testing.T) {
	if _, ok := New(config.Config{}).(*logMailer); !ok {
		t.Fatal("want logMailer when SMTP is not configured")
	}

	if _, ok := New(config.Config{SMTPHost: "smtp.example.com"}).(*smtpMailer); !ok {
		t.Fatal("want smtpMailer when SMTP host is set")
	}
}

// --- helpers ---

const (
	testFrom = "abberant@localhost"
	testTo   = "user@example.com"
	subject  = "Test subject"
	body     = "Тело письма с кодом 123456"
)

// fakeServer is a minimal SMTP server for tests. It records commands,
// captures the DATA payload and can upgrade to TLS via STARTTLS.
type fakeServer struct {
	listener net.Listener
	starttls bool
	tlsCfg   *tls.Config
	commands chan string
	message  chan string
}

func newFakeServer(t *testing.T, tlsCfg *tls.Config, starttls bool) *fakeServer {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	t.Cleanup(func() { listener.Close() })

	srv := &fakeServer{
		listener: listener,
		starttls: starttls,
		tlsCfg:   tlsCfg,
		commands: make(chan string, 64),
		message:  make(chan string, 1),
	}

	go srv.acceptLoop()

	return srv
}

func (s *fakeServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go s.handle(conn)
	}
}

func (s *fakeServer) handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Fprint(conn, "220 test ESMTP\r\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		raw := strings.TrimSpace(line)
		cmd := strings.ToUpper(raw)
		s.commands <- raw

		switch {
		case strings.HasPrefix(cmd, "EHLO"), strings.HasPrefix(cmd, "HELO"):
			caps := "250-test\r\n250-8BITMIME\r\n250-AUTH PLAIN\r\n250 SMTPUTF8\r\n"
			if s.starttls {
				caps = "250-test\r\n250-STARTTLS\r\n250-AUTH PLAIN\r\n250 8BITMIME\r\n"
			}
			fmt.Fprint(conn, caps)
		case strings.HasPrefix(cmd, "STARTTLS"):
			if !s.starttls {
				fmt.Fprint(conn, "502 not implemented\r\n")
				continue
			}
			fmt.Fprint(conn, "220 ready\r\n")
			tlsConn := tls.Server(conn, s.tlsCfg)
			if err := tlsConn.Handshake(); err != nil {
				return
			}
			conn = tlsConn
			reader = bufio.NewReader(conn)
		case strings.HasPrefix(cmd, "AUTH"):
			fmt.Fprint(conn, "235 ok\r\n")
		case strings.HasPrefix(cmd, "MAIL"), strings.HasPrefix(cmd, "RCPT"):
			fmt.Fprint(conn, "250 ok\r\n")
		case strings.HasPrefix(cmd, "DATA"):
			fmt.Fprint(conn, "354 go\r\n")

			var payload strings.Builder
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				payload.WriteString(dataLine)
				if dataLine == ".\r\n" || dataLine == ".\n" {
					break
				}
			}
			s.message <- payload.String()
			fmt.Fprint(conn, "250 queued\r\n")
		case strings.HasPrefix(cmd, "QUIT"):
			fmt.Fprint(conn, "221 bye\r\n")
			return
		default:
			fmt.Fprint(conn, "502 not implemented\r\n")
		}
	}
}

func testConfig(srv *fakeServer, starttls bool) config.Config {
	_, port, _ := net.SplitHostPort(srv.listener.Addr().String())

	return config.Config{
		SMTPHost:    "localhost",
		SMTPPort:    port,
		SMTPUser:    "user",
		SMTPPass:    "pass",
		SMTPFrom:    testFrom,
		SMTPStartTL: starttls,
	}
}

func assertFullDelivery(t *testing.T, srv *fakeServer, mailer *smtpMailer) {
	t.Helper()

	err := mailer.Send(testTo, subject, body)
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	authCmd := expectCommand(t, srv, "AUTH PLAIN")
	assertPlainAuth(t, authCmd)

	mailCmd := expectCommand(t, srv, "MAIL FROM")
	if !strings.Contains(mailCmd, testFrom) {
		t.Fatalf("unexpected MAIL FROM: %q", mailCmd)
	}

	rcptCmd := expectCommand(t, srv, "RCPT TO")
	if !strings.Contains(rcptCmd, testTo) {
		t.Fatalf("unexpected RCPT TO: %q", rcptCmd)
	}

	expectCommand(t, srv, "DATA")

	select {
	case payload := <-srv.message:
		want := string(buildMessage(testFrom, testTo, subject, body)) + "\r\n.\r\n"
		if payload != want {
			t.Fatalf("unexpected message:\n got: %q\nwant: %q", payload, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("message was not received by the server")
	}
}

// assertPlainAuth verifies the SASL PLAIN initial response carries the
// configured credentials: base64("\0user\0pass").
func assertPlainAuth(t *testing.T, cmd string) {
	t.Helper()

	encoded, ok := strings.CutPrefix(cmd, "AUTH PLAIN ")
	if !ok {
		t.Fatalf("unexpected AUTH command: %q", cmd)
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode AUTH payload: %v", err)
	}

	if want := "\x00user\x00pass"; string(decoded) != want {
		t.Fatalf("AUTH credentials: got %q, want %q", decoded, want)
	}
}

func expectCommand(t *testing.T, srv *fakeServer, prefix string) string {
	t.Helper()

	deadline := time.After(5 * time.Second)
	for {
		select {
		case cmd := <-srv.commands:
			if strings.HasPrefix(strings.ToUpper(cmd), prefix) {
				return cmd
			}
		case <-deadline:
			t.Fatalf("command %q was not sent", prefix)
		}
	}
}

// selfSignedTLSConfig returns a server TLS config with a self-signed
// certificate for "localhost" and writes the certificate PEM so that
// SSL_CERT_FILE can point to it.
func selfSignedTLSConfig(t *testing.T) (*tls.Config, string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	certFile := filepath.Join(t.TempDir(), "cert.pem")
	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatalf("write certificate: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{der},
			PrivateKey:  key,
		}},
	}, certFile
}
