package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const maxBackoff = 15 * time.Minute

const freeAttempts = 5

type Attempts struct {
	pool *pgxpool.Pool
}

func NewAttempts(pool *pgxpool.Pool) *Attempts {
	return &Attempts{pool: pool}
}

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

func (a *Attempts) RetryAfter(ctx context.Context, key string) (time.Duration, error) {
	var (
		failCount  int
		lastFailAt time.Time
	)

	err := a.pool.QueryRow(ctx, `
		SELECT fail_count, last_fail_at FROM login_attempts WHERE key = $1
	`, key).Scan(&failCount, &lastFailAt)

	if err != nil {
		return 0, nil
	}

	wait := BackoffFor(failCount)

	remaining := time.Until(lastFailAt.Add(wait))
	if remaining > 0 {
		return remaining, nil
	}

	return 0, nil
}

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

func (a *Attempts) Clear(ctx context.Context, key string) error {
	_, err := a.pool.Exec(ctx, `DELETE FROM login_attempts WHERE key = $1`, key)

	return err
}
