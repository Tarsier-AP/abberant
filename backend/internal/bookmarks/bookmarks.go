// Package bookmarks stores the articles a reader saved to return to later.
package bookmarks

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"abberant/backend/internal/articles"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// articleColumns is the article column list qualified for the bookmarks join;
// created_at exists in both tables, so every column needs the alias.
const articleColumns = `
	a.id, a.title, a.image, a.date, a.comments, a.kind, a.content, a.author, a.status,
	a.rating, a.platform, a.verdict, a.section, a.created_at
`

// IDs returns the bookmarked article ids of a user, newest bookmark first.
func (r *Repo) IDs(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT article_id FROM bookmarks
		WHERE user_id = $1
		ORDER BY created_at DESC, article_id
	`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]string, 0, 32)

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		list = append(list, id)
	}

	return list, rows.Err()
}

// Articles returns the bookmarked articles of a user, newest bookmark first.
func (r *Repo) Articles(ctx context.Context, userID int64) ([]*articles.Article, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+articleColumns+`
		FROM bookmarks b
		JOIN articles a ON a.id = b.article_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC, a.id
	`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]*articles.Article, 0, 32)

	for rows.Next() {
		article, err := articles.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, article)
	}

	return list, rows.Err()
}

// Add bookmarks the article for the user and reports whether a new row was created.
func (r *Repo) Add(ctx context.Context, userID int64, articleID string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `
		INSERT INTO bookmarks (user_id, article_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, userID, articleID)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

// Remove drops the bookmark and reports whether the row existed.
func (r *Repo) Remove(ctx context.Context, userID int64, articleID string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2
	`, userID, articleID)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}
