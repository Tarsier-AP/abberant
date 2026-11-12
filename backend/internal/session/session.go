package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const CookieName = "abberant_session"

type Session struct {
	ID         int64
	UserID     int64
	TokenHash  string
	Pending    bool
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
}

type Store struct {
	pool       *pgxpool.Pool
	sessionTTL time.Duration
	pendingTTL time.Duration
}

func NewStore(pool *pgxpool.Pool, sessionTTL, pendingTTL time.Duration) *Store {
	return &Store{pool: pool, sessionTTL: sessionTTL, pendingTTL: pendingTTL}
}

func NewToken() (token string, tokenHash string, err error) {
	raw := make([]byte, 32)

	if _, err = rand.Read(raw); err != nil {
		return "", "", err
	}

	token = base64.RawURLEncoding.EncodeToString(raw)
	tokenHash = HashToken(token)

	return token, tokenHash, nil
}

func HashToken(token string) string {
	digest := sha256.Sum256([]byte(token))

	return hex.EncodeToString(digest[:])
}

func (s *Store) Create(ctx context.Context, userID int64, pending bool) (token string, session *Session, err error) {
	token, tokenHash, err := NewToken()
	if err != nil {
		return "", nil, err
	}

	now := time.Now()
	ttl := s.sessionTTL

	if pending {
		ttl = s.pendingTTL
	}

	row := s.pool.QueryRow(ctx, `
		INSERT INTO sessions (user_id, token_hash, pending, created_at, last_seen_at, expires_at)
		VALUES ($1, $2, $3, $4, $4, $5)
		RETURNING id, user_id, token_hash, pending, created_at, last_seen_at, expires_at, revoked_at
	`, userID, tokenHash, pending, now, now.Add(ttl))

	session = &Session{}

	if err = row.Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.Pending,
		&session.CreatedAt, &session.LastSeenAt, &session.ExpiresAt, &session.RevokedAt,
	); err != nil {
		return "", nil, err
	}

	return token, session, nil
}

func (s *Store) GetActive(ctx context.Context, token string) (*Session, error) {
	tokenHash := HashToken(token)

	row := s.pool.QueryRow(ctx, `
		UPDATE sessions
		SET last_seen_at = now(),
		    expires_at = now() + $2
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND pending = false
		  AND expires_at > now()
		RETURNING id, user_id, token_hash, pending, created_at, last_seen_at, expires_at, revoked_at
	`, tokenHash, s.sessionTTL)

	var session Session

	if err := row.Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.Pending,
		&session.CreatedAt, &session.LastSeenAt, &session.ExpiresAt, &session.RevokedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}

func (s *Store) GetPending(ctx context.Context, token string) (*Session, error) {
	tokenHash := HashToken(token)

	row := s.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, pending, created_at, last_seen_at, expires_at, revoked_at
		FROM sessions
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND pending = true
		  AND expires_at > now()
	`, tokenHash)

	var session Session

	if err := row.Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.Pending,
		&session.CreatedAt, &session.LastSeenAt, &session.ExpiresAt, &session.RevokedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}

func (s *Store) Upgrade(ctx context.Context, sessionID int64) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE sessions
		SET pending = false,
		    last_seen_at = now(),
		    expires_at = now() + $2
		WHERE id = $1
		  AND revoked_at IS NULL
		  AND pending = true
	`, sessionID, s.sessionTTL)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Store) Revoke(ctx context.Context, sessionID int64) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE sessions
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL
	`, sessionID)

	return err
}

func (s *Store) RevokeAllForUser(ctx context.Context, userID int64, keepSessionID int64) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE sessions
		SET revoked_at = now()
		WHERE user_id = $1
		  AND revoked_at IS NULL
		  AND id <> $2
	`, userID, keepSessionID)

	return err
}
