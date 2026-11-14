// Package articles stores and validates editorial content for the site.
package articles

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	KindNews   = "news"
	KindReview = "review"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

// FallbackCover is used when an article has no cover image.
const FallbackCover = "/img/news-placeholder.svg"

const maxImageLength = 6 * 1024 * 1024

type Article struct {
	ID        string
	Title     string
	Image     string
	Date      time.Time
	Comments  int
	Kind      string
	Content   string
	Author    string
	Status    string
	Rating    int // 0 means the article has no rating (news)
	Platform  string
	Verdict   string
	Section   string
	CreatedAt time.Time
}

const articleColumns = `
	id, title, image, date, comments, kind, content, author, status,
	rating, platform, verdict, section, created_at
`

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// ScanRow builds an article from a database row. It is exported so other
// packages can join the articles table with their own column alias.
func ScanRow(row pgx.Row) (*Article, error) {
	var (
		article  Article
		rating   *int
		platform *string
		verdict  *string
		section  *string
	)

	err := row.Scan(
		&article.ID, &article.Title, &article.Image, &article.Date, &article.Comments,
		&article.Kind, &article.Content, &article.Author, &article.Status,
		&rating, &platform, &verdict, &section, &article.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	if rating != nil {
		article.Rating = *rating
	}

	if platform != nil {
		article.Platform = *platform
	}

	if verdict != nil {
		article.Verdict = *verdict
	}

	if section != nil {
		article.Section = *section
	}

	return &article, nil
}

// All returns every article, newest first.
func (r *Repo) All(ctx context.Context) ([]*Article, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+articleColumns+` FROM articles
		ORDER BY date DESC, created_at DESC, id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]*Article, 0, 64)

	for rows.Next() {
		article, err := ScanRow(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, article)
	}

	return list, rows.Err()
}

func (r *Repo) ByID(ctx context.Context, id string) (*Article, error) {
	return ScanRow(r.pool.QueryRow(ctx, `
		SELECT `+articleColumns+` FROM articles WHERE id = $1
	`, id))
}

// TitlesByIDs returns id → title for every article from the list. Unknown ids
// are simply missing from the map.
func (r *Repo) TitlesByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, title FROM articles WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	titles := make(map[string]string, len(ids))

	for rows.Next() {
		var (
			id    string
			title string
		)

		if err := rows.Scan(&id, &title); err != nil {
			return nil, err
		}

		titles[id] = title
	}

	return titles, rows.Err()
}

func (r *Repo) Create(ctx context.Context, article *Article) (*Article, error) {
	return ScanRow(r.pool.QueryRow(ctx, `
		INSERT INTO articles (id, title, image, date, comments, kind, content, author, status, rating, platform, verdict, section)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING `+articleColumns,
		article.ID, article.Title, article.Image, article.Date, article.Comments,
		article.Kind, article.Content, article.Author, article.Status,
		nullableRating(article.Rating), article.Platform, article.Verdict, article.Section,
	))
}

// Update rewrites every editable column of the article.
func (r *Repo) Update(ctx context.Context, article *Article) (*Article, error) {
	return ScanRow(r.pool.QueryRow(ctx, `
		UPDATE articles
		SET title = $2, image = $3, date = $4, comments = $5, kind = $6, content = $7,
		    author = $8, status = $9, rating = $10, platform = $11, verdict = $12, section = $13
		WHERE id = $1
		RETURNING `+articleColumns,
		article.ID, article.Title, article.Image, article.Date, article.Comments,
		article.Kind, article.Content, article.Author, article.Status,
		nullableRating(article.Rating), article.Platform, article.Verdict, article.Section,
	))
}

// Delete removes the article and reports whether the row existed.
func (r *Repo) Delete(ctx context.Context, id string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

// NextUserID returns the next free `user-N` identifier.
func (r *Repo) NextUserID(ctx context.Context) (string, error) {
	var maxUsed int

	err := r.pool.QueryRow(ctx, `
		SELECT coalesce(max(substring(id from '^user-([0-9]+)$')::int), 0)
		FROM articles WHERE id ~ '^user-[0-9]+$'
	`).Scan(&maxUsed)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("user-%d", maxUsed+1), nil
}

func nullableRating(rating int) any {
	if rating <= 0 {
		return nil
	}

	return rating
}

// IsUniqueViolation reports whether the error is a primary key conflict.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// ParseDate parses an article date in the YYYY-MM-DD form.
func ParseDate(value string) (time.Time, bool) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, false
	}

	return parsed, true
}

func ValidKind(kind string) bool {
	return kind == KindNews || kind == KindReview
}

func ValidStatus(status string) bool {
	return status == StatusDraft || status == StatusPublished
}

func ValidRating(rating int) bool {
	return rating >= 1 && rating <= 5
}

// ValidImage accepts an empty value, an absolute or site-relative URL,
// or an image data URL of at most 6 MB (same limit as avatars).
func ValidImage(image string) bool {
	if image == "" {
		return true
	}

	if len(image) > maxImageLength {
		return false
	}

	if strings.HasPrefix(image, "data:") {
		return strings.HasPrefix(image, "data:image/")
	}

	return true
}
