// Package ratings stores reader votes (1–5 stars) on articles.
package ratings

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Stats is the aggregated reader rating of one article.
type Stats struct {
	Average float64
	Votes   int
}

// VoteRecord is one reader vote with its last-change time.
type VoteRecord struct {
	ArticleID string
	Stars     int
	UpdatedAt time.Time
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// Put stores or updates the user's vote and returns the fresh article stats.
func (r *Repo) Put(ctx context.Context, userID int64, articleID string, stars int) (Stats, error) {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO article_votes (user_id, article_id, stars)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, article_id)
		DO UPDATE SET stars = EXCLUDED.stars, updated_at = now()
	`, userID, articleID, stars)
	if err != nil {
		return Stats{}, err
	}

	return r.StatsByArticle(ctx, articleID)
}

// Delete removes the user's vote and reports whether the row existed.
func (r *Repo) Delete(ctx context.Context, userID int64, articleID string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM article_votes WHERE user_id = $1 AND article_id = $2
	`, userID, articleID)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

// VoteOf returns the user's own stars on the article, 0 when there is none.
func (r *Repo) VoteOf(ctx context.Context, userID int64, articleID string) (int, error) {
	var stars int

	err := r.pool.QueryRow(ctx, `
		SELECT stars FROM article_votes
		WHERE user_id = $1 AND article_id = $2
	`, userID, articleID).Scan(&stars)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return stars, nil
}

// StatsByArticle aggregates the reader votes of one article.
func (r *Repo) StatsByArticle(ctx context.Context, articleID string) (Stats, error) {
	var (
		sum   int
		count int
	)

	err := r.pool.QueryRow(ctx, `
		SELECT coalesce(sum(stars), 0), count(*) FROM article_votes WHERE article_id = $1
	`, articleID).Scan(&sum, &count)

	if err != nil {
		return Stats{}, err
	}

	return statsFrom(sum, count), nil
}

// StatsByArticleMap aggregates the reader votes of every voted article.
func (r *Repo) StatsByArticleMap(ctx context.Context) (map[string]Stats, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT article_id, coalesce(sum(stars), 0), count(*)
		FROM article_votes GROUP BY article_id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[string]Stats, 32)

	for rows.Next() {
		var (
			articleID string
			sum       int
			count     int
		)

		if err := rows.Scan(&articleID, &sum, &count); err != nil {
			return nil, err
		}

		result[articleID] = statsFrom(sum, count)
	}

	return result, rows.Err()
}

// VotesOf returns the user's votes keyed by article id.
func (r *Repo) VotesOf(ctx context.Context, userID int64) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT article_id, stars FROM article_votes WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	votes := make(map[string]int, 16)

	for rows.Next() {
		var (
			articleID string
			stars     int
		)

		if err := rows.Scan(&articleID, &stars); err != nil {
			return nil, err
		}

		votes[articleID] = stars
	}

	return votes, rows.Err()
}

// VoteRecords returns the user's votes with their timestamps, newest first.
func (r *Repo) VoteRecords(ctx context.Context, userID int64, limit int) ([]VoteRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT article_id, stars, updated_at FROM article_votes
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]VoteRecord, 0, limit)

	for rows.Next() {
		var record VoteRecord

		if err := rows.Scan(&record.ArticleID, &record.Stars, &record.UpdatedAt); err != nil {
			return nil, err
		}

		list = append(list, record)
	}

	return list, rows.Err()
}

func statsFrom(sum, count int) Stats {
	if count == 0 {
		return Stats{Average: 0, Votes: 0}
	}

	return Stats{Average: float64(sum) / float64(count), Votes: count}
}
