package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	KindActivation = "activation"
	KindTwoFactor  = "twofactor"
)

const CodeTTL = 5 * time.Minute

type Codes struct {
	pool *pgxpool.Pool
	ttl  time.Duration
}

func NewCodes(pool *pgxpool.Pool, ttl time.Duration) *Codes {
	return &Codes{pool: pool, ttl: ttl}
}

func NewCode() (string, error) {
	raw := make([]byte, 4)

	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	value := int(uint32(raw[0])<<24|uint32(raw[1])<<16|uint32(raw[2])<<8|uint32(raw[3])) % 1_000_000

	return fmt.Sprintf("%06d", value), nil
}

func HashCode(code string) string {
	digest := sha256.Sum256([]byte(code))

	return hex.EncodeToString(digest[:])
}

func (c *Codes) Issue(ctx context.Context, userID int64, email, kind string) (string, error) {
	code, err := NewCode()
	if err != nil {
		return "", err
	}

	_, err = c.pool.Exec(ctx, `
		UPDATE codes
		SET used_at = now()
		WHERE user_id = $1 AND kind = $2 AND used_at IS NULL
	`, userID, kind)
	if err != nil {
		return "", err
	}

	_, err = c.pool.Exec(ctx, `
		INSERT INTO codes (user_id, email, kind, code_hash, expires_at)
		VALUES ($1, $2, $3, $4, now() + $5)
	`, userID, email, kind, HashCode(code), c.ttl)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (c *Codes) Verify(ctx context.Context, userID int64, kind, code string) (bool, error) {
	var (
		codeHash  string
		expiresAt time.Time
	)

	err := c.pool.QueryRow(ctx, `
		SELECT code_hash, expires_at
		FROM codes
		WHERE user_id = $1 AND kind = $2 AND used_at IS NULL AND expires_at > now()
		ORDER BY id DESC
		LIMIT 1
	`, userID, kind).Scan(&codeHash, &expiresAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if subtle.ConstantTimeCompare([]byte(HashCode(code)), []byte(codeHash)) != 1 {
		return false, nil
	}

	tag, err := c.pool.Exec(ctx, `
		UPDATE codes
		SET used_at = now()
		WHERE user_id = $1 AND kind = $2 AND code_hash = $3 AND used_at IS NULL
	`, userID, kind, codeHash)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}
