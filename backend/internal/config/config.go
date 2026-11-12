package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr         string
	DSN          string
	DataDir      string
	EmbeddedPort int
	CookieSecure bool
	PublicOrigin string
	DevMode      bool

	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	SMTPStartTL bool

	SessionTTL time.Duration
	PendingTTL time.Duration
	CodeTTL    time.Duration
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func envBool(key string, fallback bool) bool {
	raw := env(key, "")

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return value
}

func envInt(key string, fallback int) int {
	raw := env(key, "")

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}

func Load() Config {
	return Config{
		Addr:         env("ABBERANT_ADDR", ":8080"),
		DSN:          env("ABBERANT_DATABASE_DSN", ""),
		DataDir:      env("ABBERANT_DB_DATA_DIR", ".data/postgres"),
		EmbeddedPort: envInt("ABBERANT_DB_PORT", 55432),
		CookieSecure: envBool("ABBERANT_COOKIE_SECURE", false),
		PublicOrigin: env("ABBERANT_PUBLIC_ORIGIN", "http://localhost:5173"),
		DevMode:      envBool("ABBERANT_DEV_MODE", true),

		SMTPHost:    env("SMTP_HOST", ""),
		SMTPPort:    env("SMTP_PORT", "587"),
		SMTPUser:    env("SMTP_USERNAME", ""),
		SMTPPass:    env("SMTP_PASSWORD", ""),
		SMTPFrom:    env("SMTP_FROM", "abberant@localhost"),
		SMTPStartTL: envBool("SMTP_STARTTLS", true),

		SessionTTL: 7 * 24 * time.Hour,
		PendingTTL: 5 * time.Minute,
		CodeTTL:    5 * time.Minute,
	}
}

func (c Config) SMTPEnabled() bool {
	return c.SMTPHost != ""
}

func (c Config) DeliveryDevMode() bool {
	return c.DevMode && !c.SMTPEnabled()
}
