// Package sections stores the server-managed section metadata: news platform
// feeds and walkthrough game hubs. Review groups are feed groupings rendered
// from static frontend config, so they are not stored here.
package sections

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TypeNews        = "news"
	TypeWalkthrough = "walkthrough"
)

// Section is one manageable section row. Archived sections stay readable for
// their assigned articles but disappear from pickers and public feeds.
type Section struct {
	ID        string
	Type      string
	Title     string
	Archived  bool
	CreatedAt time.Time
}

// IsValidType reports whether sectionType is a stored section type.
func IsValidType(sectionType string) bool {
	return sectionType == TypeNews || sectionType == TypeWalkthrough
}

// IsUniqueViolation reports whether err is a PostgreSQL unique constraint
// failure, mirroring articles.IsUniqueViolation.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

const sectionColumns = `id, type, title, archived, created_at`

func scanSection(row pgx.Row) (*Section, error) {
	var s Section

	if err := row.Scan(&s.ID, &s.Type, &s.Title, &s.Archived, &s.CreatedAt); err != nil {
		return nil, err
	}

	return &s, nil
}

// List returns every section, oldest first, so frontend ordering stays stable.
func (r *Repo) List(ctx context.Context) ([]*Section, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+sectionColumns+` FROM sections ORDER BY created_at, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*Section, 0)

	for rows.Next() {
		s, err := scanSection(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, s)
	}

	return list, rows.Err()
}

// ByID returns one section, or nil when absent.
func (r *Repo) ByID(ctx context.Context, id string) (*Section, error) {
	s, err := scanSection(r.pool.QueryRow(ctx,
		`SELECT `+sectionColumns+` FROM sections WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return s, err
}

// Create inserts a section under an auto-generated slug derived from title.
// Slug collisions are resolved with a numeric suffix.
func (r *Repo) Create(ctx context.Context, sectionType, title string) (*Section, error) {
	base := Slugify(title)
	if base == "" {
		base = "section"
	}

	for attempt := 0; attempt < 50; attempt++ {
		id := base
		if attempt > 0 {
			id = base + "-" + strconv.Itoa(attempt+1)
		}

		s, err := r.insert(ctx, id, sectionType, title)
		if err != nil {
			if IsUniqueViolation(err) {
				continue
			}

			return nil, err
		}

		return s, nil
	}

	return nil, errors.New("sections: could not allocate a unique slug")
}

func (r *Repo) insert(ctx context.Context, id, sectionType, title string) (*Section, error) {
	return scanSection(r.pool.QueryRow(ctx,
		`INSERT INTO sections (id, type, title) VALUES ($1, $2, $3)
		 RETURNING `+sectionColumns, id, sectionType, title))
}

// Rename updates the display title; the id/slug is immutable.
func (r *Repo) Rename(ctx context.Context, id, title string) (*Section, error) {
	return scanSection(r.pool.QueryRow(ctx,
		`UPDATE sections SET title = $2 WHERE id = $1
		 RETURNING `+sectionColumns, id, title))
}

// Update persists title and archived state of an already loaded section.
func (r *Repo) Update(ctx context.Context, s *Section) (*Section, error) {
	return scanSection(r.pool.QueryRow(ctx,
		`UPDATE sections SET title = $2, archived = $3 WHERE id = $1
		 RETURNING `+sectionColumns, s.ID, s.Title, s.Archived))
}

// SetArchived soft-deletes (or restores) a section.
func (r *Repo) SetArchived(ctx context.Context, id string, archived bool) (*Section, error) {
	return scanSection(r.pool.QueryRow(ctx,
		`UPDATE sections SET archived = $2 WHERE id = $1
		 RETURNING `+sectionColumns, id, archived))
}

// Slugify converts a display title into a URL-safe latin slug, transliterating
// Cyrillic and collapsing runs of separators into single dashes.
func Slugify(title string) string {
	var b strings.Builder

	lastDash := true // suppress a leading dash

	for _, r := range strings.TrimSpace(strings.ToLower(title)) {
		if tr, ok := cyrillic[r]; ok {
			b.WriteString(tr)
			lastDash = false

			continue
		}

		switch {
		case unicode.IsLetter(r) && r < unicode.MaxLatin1 || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case unicode.IsLetter(r) && r >= unicode.MaxLatin1:
			// Unmapped letters degrade to a separator.
			fallthrough
		default:
			if !lastDash && r != '\'' && r != '’' {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(b.String(), "-")
}

var cyrillic = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "",
	'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya", 'і': "i", 'ї': "yi",
	'є': "ye", 'ґ': "g",
}
