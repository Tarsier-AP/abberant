// Package notifications collects in-app events for readers: replies to their
// comments and new material in the news sections they follow.
package notifications

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TypeReply   = "reply"
	TypeArticle = "article"
)

// Notification is one bell entry. Actor is set for replies and nil for
// editorial "new material" events.
type Notification struct {
	ID           int64
	Type         string
	ActorID      *int64
	ActorName    string
	ActorAvatar  string
	ArticleID    *string
	ArticleTitle string
	CommentID    *int64
	CreatedAt    time.Time
	Read         bool
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// List returns the newest notifications of a user with actor and article data.
func (r *Repo) List(ctx context.Context, userID int64, limit int) ([]*Notification, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT n.id, n.type, n.actor_id, u.username, u.avatar,
		       n.article_id, a.title, n.comment_id, n.created_at,
		       (n.read_at IS NOT NULL)
		FROM notifications n
		JOIN articles a ON a.id = n.article_id
		LEFT JOIN users u ON u.id = n.actor_id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC, n.id DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]*Notification, 0, limit)

	for rows.Next() {
		var (
			item        Notification
			actorName   *string
			actorAvatar *string
		)

		if err := rows.Scan(
			&item.ID, &item.Type, &item.ActorID, &actorName, &actorAvatar,
			&item.ArticleID, &item.ArticleTitle, &item.CommentID, &item.CreatedAt,
			&item.Read,
		); err != nil {
			return nil, err
		}

		if actorName != nil {
			item.ActorName = *actorName
		}

		if actorAvatar != nil {
			item.ActorAvatar = *actorAvatar
		}

		list = append(list, &item)
	}

	return list, rows.Err()
}

// UnreadCount returns how many notifications the user has not seen yet.
func (r *Repo) UnreadCount(ctx context.Context, userID int64) (int, error) {
	var count int

	err := r.pool.QueryRow(ctx, `
		SELECT count(*) FROM notifications
		WHERE user_id = $1 AND read_at IS NULL
	`, userID).Scan(&count)

	return count, err
}

// MarkAllRead marks every unread notification of the user and reports how
// many rows changed.
func (r *Repo) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE notifications SET read_at = now()
		WHERE user_id = $1 AND read_at IS NULL
	`, userID)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

// CreateReply notifies the parent comment author about a fresh reply.
func (r *Repo) CreateReply(ctx context.Context, userID, actorID int64, articleID string, commentID int64) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, actor_id, article_id, comment_id)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, TypeReply, actorID, articleID, commentID)

	return err
}

// CreateArticleEvents tells every subscriber of the section about a published
// article, except its author. Returns how many notifications were created.
func (r *Repo) CreateArticleEvents(ctx context.Context, section, articleID string, exceptUserID int64) (int64, error) {
	tag, err := r.pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, article_id)
		SELECT user_id, $2, $3 FROM section_subscriptions
		WHERE section = $1 AND user_id <> $4
	`, section, TypeArticle, articleID, exceptUserID)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

// Subscriptions returns the section ids the user follows, alphabetical.
func (r *Repo) Subscriptions(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT section FROM section_subscriptions
		WHERE user_id = $1
		ORDER BY section
	`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]string, 0, 8)

	for rows.Next() {
		var section string

		if err := rows.Scan(&section); err != nil {
			return nil, err
		}

		list = append(list, section)
	}

	return list, rows.Err()
}

// Subscribe follows the section and reports whether a new row was created.
func (r *Repo) Subscribe(ctx context.Context, userID int64, section string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `
		INSERT INTO section_subscriptions (user_id, section)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, userID, section)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}

// Unsubscribe drops the section follow and reports whether the row existed.
func (r *Repo) Unsubscribe(ctx context.Context, userID int64, section string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM section_subscriptions WHERE user_id = $1 AND section = $2
	`, userID, section)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}
