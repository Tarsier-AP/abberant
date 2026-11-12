package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// maxBackoff caps the progressive lockout window.
const maxBackoff = 15 * time.Minute

// freeAttempts is the number of failures allowed before backoff starts.
const freeAttempts = 5

// Attempts tracks failed attempts per key (account or IP) with backoff.
type Attempts struct {
	pool *pgxpool.Pool
}

// NewAttempts builds the attempt counter.
func NewAttempts(pool *pgxpool.Pool) *Attempts {
	return &Attempts{pool: pool}
}

// BackoffFor returns how long the key must stay idle after n failures.
func BackoffFor(failCount int) time.Duration {
	if failCount < freeAttempts {
		return 0
	}

	exponent := failCount - freeAttempts
	if exponent > 30 {
		exponent = 30
	}

	delay := time.Minute << uint(exponent)
	if delay > maxBackoff || delay <= 0 {
		return maxBackoff
	}

	return delay
}

// RetryAfter returns the remaining wait for the key or zero when a new
// attempt is allowed.
func (a *Attempts) RetryAfter(ctx context.Context, key string) (time.Duration, error) {
	var (
		failCount  int
		lastFailAt time.Time
	)

	err := a.pool.QueryRow(ctx, `
		SELECT fail_count, last_fail_at FROM login_attempts WHERE key = $1
	`, key).Scan(&failCount, &lastFailAt)

	if err != nil {
		// No row means no failures recorded yet.
		return 0, nil
	}

	wait := BackoffFor(failCount)

	remaining := time.Until(lastFailAt.Add(wait))
	if remaining > 0 {
		return remaining, nil
	}

	return 0, nil
}

// Record bumps the failure counter for the key.
func (a *Attempts) Record(ctx context.Context, key string) error {
	_, err := a.pool.Exec(ctx, `
		INSERT INTO login_attempts (key, fail_count, last_fail_at)
		VALUES ($1, 1, now())
		ON CONFLICT (key) DO UPDATE
		SET fail_count = CASE
				WHEN login_attempts.fail_count < 1000 THEN login_attempts.fail_count + 1
				ELSE login_attempts.fail_count
			END,
		    last_fail_at = now()
	`, key)

	return err
}

// Clear resets the failure counter after a successful attempt.
func (a *Attempts) Clear(ctx context.Context, key string) error {
	_, err := a.pool.Exec(ctx, `DELETE FROM login_attempts WHERE key = $1`, key)

	return err
}
