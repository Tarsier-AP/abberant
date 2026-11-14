package api

import (
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

/** Activity pages show at most this many recent events. */
const activityLimit = 100

type activityResponse struct {
	/** `comment`, `vote` or `bookmark`. */
	Type      string             `json:"type"`
	CreatedAt string             `json:"createdAt"`
	Article   activityArticleRef `json:"article"`
	// Comment events.
	CommentID *int64  `json:"commentId,omitempty"`
	Body      *string `json:"body,omitempty"`
	// Vote events.
	Stars *int `json:"stars,omitempty"`
}

type activityArticleRef struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type activityEvent struct {
	response  activityResponse
	createdAt time.Time
}

func (s *Server) handleActivity(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	events := make([]activityEvent, 0, activityLimit)

	comments, err := s.comments.ListByUser(ctx.Context(), user.ID, activityLimit)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		articleID := comment.ArticleID
		commentID := comment.ID
		body := comment.Body

		events = append(events, activityEvent{
			response: activityResponse{
				Type:      "comment",
				Article:   activityArticleRef{ID: articleID},
				CommentID: &commentID,
				Body:      &body,
			},
			createdAt: comment.CreatedAt,
		})
	}

	votes, err := s.ratings.VoteRecords(ctx.Context(), user.ID, activityLimit)
	if err != nil {
		return err
	}

	for _, vote := range votes {
		stars := vote.Stars

		events = append(events, activityEvent{
			response: activityResponse{
				Type:    "vote",
				Article: activityArticleRef{ID: vote.ArticleID},
				Stars:   &stars,
			},
			createdAt: vote.UpdatedAt,
		})
	}

	bookmarks, err := s.bookmarks.Records(ctx.Context(), user.ID, activityLimit)
	if err != nil {
		return err
	}

	for _, bookmark := range bookmarks {
		events = append(events, activityEvent{
			response: activityResponse{
				Type:    "bookmark",
				Article: activityArticleRef{ID: bookmark.ArticleID},
			},
			createdAt: bookmark.CreatedAt,
		})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].createdAt.After(events[j].createdAt)
	})

	if len(events) > activityLimit {
		events = events[:activityLimit]
	}

	// Titles for every article referenced by the kept events.
	ids := make([]string, 0, len(events))

	for _, event := range events {
		ids = append(ids, event.response.Article.ID)
	}

	titles, err := s.articles.TitlesByIDs(ctx.Context(), ids)
	if err != nil {
		return err
	}

	response := make([]activityResponse, 0, len(events))

	for _, event := range events {
		event.response.CreatedAt = event.createdAt.Format(time.RFC3339)
		event.response.Article = activityArticleRef{
			ID:    event.response.Article.ID,
			Title: titles[event.response.Article.ID],
		}

		response = append(response, event.response)
	}

	return ctx.JSON(fiber.Map{"activity": response})
}
