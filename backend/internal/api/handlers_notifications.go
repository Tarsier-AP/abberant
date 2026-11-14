package api

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/notifications"
)

/** Same list length as the bell panel shows. */
const notificationListLimit = 50

type notificationResponse struct {
	ID        int64                `json:"id"`
	Type      string               `json:"type"`
	CreatedAt string               `json:"createdAt"`
	Read      bool                 `json:"read"`
	Actor     *notificationActor   `json:"actor,omitempty"`
	Article   *notificationArticle `json:"article,omitempty"`
	CommentID *int64               `json:"commentId,omitempty"`
}

type notificationActor struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type notificationArticle struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func newNotificationResponse(item *notifications.Notification) notificationResponse {
	response := notificationResponse{
		ID:        item.ID,
		Type:      item.Type,
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		Read:      item.Read,
		CommentID: item.CommentID,
	}

	if item.ActorID != nil {
		response.Actor = &notificationActor{
			Username: item.ActorName,
			Avatar:   item.ActorAvatar,
		}
	}

	if item.ArticleID != nil {
		response.Article = &notificationArticle{
			ID:    *item.ArticleID,
			Title: item.ArticleTitle,
		}
	}

	return response
}

func (s *Server) handleListNotifications(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	list, err := s.notifications.List(ctx.Context(), user.ID, notificationListLimit)
	if err != nil {
		return err
	}

	response := make([]notificationResponse, 0, len(list))

	for _, item := range list {
		response = append(response, newNotificationResponse(item))
	}

	return ctx.JSON(fiber.Map{"notifications": response})
}

func (s *Server) handleUnreadCount(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	count, err := s.notifications.UnreadCount(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"count": count})
}

func (s *Server) handleMarkAllRead(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	marked, err := s.notifications.MarkAllRead(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok", "marked": marked})
}

func (s *Server) handleListSubscriptions(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	sections, err := s.notifications.Subscriptions(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	if sections == nil {
		sections = []string{}
	}

	return ctx.JSON(fiber.Map{"sections": sections})
}

func (s *Server) handleSubscribe(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	section := strings.TrimSpace(ctx.Params("section"))
	if section == "" || len(section) > 64 {
		return auth.Fail(400, "Некорректный идентификатор раздела.")
	}

	if _, err := s.notifications.Subscribe(ctx.Context(), user.ID, section); err != nil {
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{"status": "ok"})
}

func (s *Server) handleUnsubscribe(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	if _, err := s.notifications.Unsubscribe(ctx.Context(), user.ID, ctx.Params("section")); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}
