// Package comments stores reader discussions under each article.
package comments

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Comment is a single message under an article. A nil ParentID marks a
// top-level comment; replies are nested exactly one level deep.
type Comment struct {
	ID        int64
	ArticleID string
	UserID    int64
	ParentID  *int64
	Body      string
	CreatedAt time.Time
	Author    Author
}

// Author is the public view of the account that wrote the comment.
type Author struct {
	ID       int64
	Username string
	Avatar   string
	Role     string
}

const commentColumns = `
	c.id, c.article_id, c.user_id, c.parent_id, c.body, c.created_at,
	u.id, u.username, u.avatar_data_url, u.role
`

const commentJoin = `
	FROM comments c
	JOIN users u ON u.id = c.user_id
`

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

func scanComment(row pgx.Row) (*Comment, error) {
	var (
		comment  Comment
		parentID *int64
	)

	err := row.Scan(
		&comment.ID, &comment.ArticleID, &comment.UserID, &parentID,
		&comment.Body, &comment.CreatedAt,
		&comment.Author.ID, &comment.Author.Username, &comment.Author.Avatar, &comment.Author.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	comment.ParentID = parentID

	return &comment, nil
}

// ListByArticle returns the comments of an article, oldest first.
func (r *Repo) ListByArticle(ctx context.Context, articleID string) ([]*Comment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+commentColumns+commentJoin+`
		WHERE c.article_id = $1
		ORDER BY c.created_at, c.id
	`, articleID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]*Comment, 0, 32)

	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, comment)
	}

	return list, rows.Err()
}

func (r *Repo) ByID(ctx context.Context, id int64) (*Comment, error) {
	return scanComment(r.pool.QueryRow(ctx, `
		SELECT `+commentColumns+commentJoin+`
		WHERE c.id = $1
	`, id))
}

// Create adds a comment and returns it joined with the author data.
func (r *Repo) Create(ctx context.Context, articleID string, userID int64, parentID *int64, body string) (*Comment, error) {
	var id int64

	err := r.pool.QueryRow(ctx, `
		INSERT INTO comments (article_id, user_id, parent_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, articleID, userID, parentID, body).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.ByID(ctx, id)
}

// Delete removes the comment (its replies go with it) and reports whether the row existed.
func (r *Repo) Delete(ctx context.Context, id int64) (bool, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

// CountsByArticle returns the number of comments for every discussed article.
func (r *Repo) CountsByArticle(ctx context.Context) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT article_id, count(*) FROM comments GROUP BY article_id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	counts := make(map[string]int, 64)

	for rows.Next() {
		var articleID string

		var count int

		if err := rows.Scan(&articleID, &count); err != nil {
			return nil, err
		}

		counts[articleID] = count
	}

	return counts, rows.Err()
}
