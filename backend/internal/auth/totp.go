package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// totpPeriod is the standard TOTP timing: 30-second windows.
const totpPeriod = 30

// BackupCodeCount is how many one-time backup codes are issued.
const BackupCodeCount = 8

// TOTP manages authenticator-app enrollment and backup codes.
type TOTP struct {
	pool *pgxpool.Pool
}

// NewTOTP builds the TOTP manager.
func NewTOTP(pool *pgxpool.Pool) *TOTP {
	return &TOTP{pool: pool}
}

// Enrollment describes a freshly generated TOTP secret.
type Enrollment struct {
	Secret     string
	OtpauthURL string
	QRDataURL  string
}

// GenerateSecret creates a new secret plus a QR code as a data URL.
func (t *TOTP) GenerateSecret(email string) (*Enrollment, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Abberant",
		AccountName: email,
		Period:      totpPeriod,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("generate totp secret: %w", err)
	}

	png, err := qrcode.New(key.URL(), qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("build qr code: %w", err)
	}

	var image bytes.Buffer

	if err := png.Write(256, &image); err != nil {
		return nil, fmt.Errorf("render qr code: %w", err)
	}

	return &Enrollment{
		Secret:     key.Secret(),
		OtpauthURL: key.URL(),
		QRDataURL:  "data:image/png;base64," + base64.StdEncoding.EncodeToString(image.Bytes()),
	}, nil
}

// Verify checks a live TOTP code, allowing one window of drift in either
// direction (~90 seconds of validity).
func (t *TOTP) Verify(secret, code string) bool {
	if secret == "" || code == "" {
		return false
	}

	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    totpPeriod,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return false
	}

	return valid
}

// GenerateBackupCode returns one code in the XXXX-XXXX form.
func GenerateBackupCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	raw := make([]byte, 8)

	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	out := make([]byte, 0, 9)

	for i, value := range raw {
		if i == 4 {
			out = append(out, '-')
		}

		out = append(out, alphabet[int(value)%len(alphabet)])
	}

	return string(out), nil
}

// HashBackupCode hashes a backup code for storage.
func HashBackupCode(code string) string {
	digest := sha256.Sum256([]byte(code))

	return hex.EncodeToString(digest[:])
}

// GenerateBackupCodes returns n fresh backup codes.
func GenerateBackupCodes(n int) ([]string, error) {
	codes := make([]string, 0, n)

	for range n {
		code, err := GenerateBackupCode()
		if err != nil {
			return nil, err
		}

		codes = append(codes, code)
	}

	return codes, nil
}

// StoreBackupCodes replaces unused backup codes with the given ones.
func (t *TOTP) StoreBackupCodes(ctx context.Context, userID int64, codes []string) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx) //nolint:errcheck // read-only rollback

	if _, err := tx.Exec(ctx, `
		DELETE FROM totp_backup_codes WHERE user_id = $1 AND used_at IS NULL
	`, userID); err != nil {
		return err
	}

	for _, code := range codes {
		if _, err := tx.Exec(ctx, `
			INSERT INTO totp_backup_codes (user_id, code_hash) VALUES ($1, $2)
		`, userID, HashBackupCode(code)); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// ConsumeBackupCode checks the code against every unused backup code of the
// user and marks the match as used.
func (t *TOTP) ConsumeBackupCode(ctx context.Context, userID int64, code string) (bool, error) {
	normalized := normalizeBackupCode(code)
	if normalized == "" {
		return false, nil
	}

	rows, err := t.pool.Query(ctx, `
		SELECT id, code_hash FROM totp_backup_codes WHERE user_id = $1 AND used_at IS NULL
	`, userID)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	target := HashBackupCode(normalized)

	var matchedID int64

	for rows.Next() {
		var (
			id       int64
			codeHash string
		)

		if err := rows.Scan(&id, &codeHash); err != nil {
			return false, err
		}

		if subtle.ConstantTimeCompare([]byte(target), []byte(codeHash)) == 1 {
			matchedID = id
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	if matchedID == 0 {
		return false, nil
	}

	tag, err := t.pool.Exec(ctx, `
		UPDATE totp_backup_codes SET used_at = now() WHERE id = $1 AND used_at IS NULL
	`, matchedID)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

// normalizeBackupCode strips spaces and a surrounding dash variant so users
// can paste codes with or without separators.
func normalizeBackupCode(code string) string {
	var builder []rune

	for _, symbol := range code {
		if symbol == ' ' {
			continue
		}

		if 'a' <= symbol && symbol <= 'z' {
			symbol -= 'a' - 'A'
		}

		builder = append(builder, symbol)
	}

	normalized := string(builder)
	if len(normalized) == 9 && normalized[4] == '-' {
		return normalized
	}

	if len(normalized) == 8 {
		return normalized[:4] + "-" + normalized[4:]
	}

	return ""
}
