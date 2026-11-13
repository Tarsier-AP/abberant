package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
	"abberant/backend/internal/comments"
)

// Same limit as the comments table check constraint.
const maxCommentLength = 500

type commentAuthorResponse struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

type commentResponse struct {
	ID        int64                 `json:"id"`
	ArticleID string                `json:"articleId"`
	ParentID  *int64                `json:"parentId"`
	Body      string                `json:"body"`
	CreatedAt string                `json:"createdAt"`
	Author    commentAuthorResponse `json:"author"`
}

func newCommentResponse(comment *comments.Comment) commentResponse {
	return commentResponse{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		ParentID:  comment.ParentID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		Author: commentAuthorResponse{
			Username: comment.Author.Username,
			Avatar:   comment.Author.Avatar,
			Role:     comment.Author.Role,
		},
	}
}

// canSeeArticle reports whether the account may view a non-published article.
func canSeeArticle(article *articles.Article, user *auth.User) bool {
	if article.Status == articles.StatusPublished {
		return true
	}

	return user != nil && (canModerateArticles(user) || article.Author == user.Username)
}

func (s *Server) handleListComments(ctx *fiber.Ctx) error {
	_, user, err := s.resolve(ctx)
	if err != nil {
		return err
	}

	article, err := s.articles.ByID(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	if article == nil || !canSeeArticle(article, user) {
		return auth.Fail(404, "Статья не найдена.")
	}

	list, err := s.comments.ListByArticle(ctx.Context(), article.ID)
	if err != nil {
		return err
	}

	response := make([]commentResponse, 0, len(list))

	for _, comment := range list {
		response = append(response, newCommentResponse(comment))
	}

	return ctx.JSON(fiber.Map{"comments": response})
}

type createCommentRequest struct {
	Body     string `json:"body"`
	ParentID *int64 `json:"parentId"`
}

func (s *Server) handleCreateComment(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	article, err := s.articles.ByID(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	if article == nil || !canSeeArticle(article, user) {
		return auth.Fail(404, "Статья не найдена.")
	}

	if article.Status != articles.StatusPublished {
		return auth.Fail(403, "Комментировать можно только опубликованные статьи.")
	}

	var req createCommentRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	body := strings.TrimSpace(req.Body)
	if body == "" {
		return auth.Fail(400, "Введите текст комментария.")
	}

	if len([]rune(body)) > maxCommentLength {
		return auth.Fail(400, "Комментарий не может быть длиннее 500 символов.")
	}

	var parentID *int64

	if req.ParentID != nil {
		parent, err := s.comments.ByID(ctx.Context(), *req.ParentID)
		if err != nil {
			return err
		}

		switch {
		case parent == nil || parent.ArticleID != article.ID:
			return auth.Fail(400, "Комментарий, на который вы отвечаете, не найден.")

		case parent.ParentID != nil:
			return auth.Fail(400, "Отвечать можно только на комментарии верхнего уровня.")

		default:
			parentID = &parent.ID
		}
	}

	created, err := s.comments.Create(ctx.Context(), article.ID, user.ID, parentID, body)
	if err != nil {
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{"comment": newCommentResponse(created)})
}

func (s *Server) handleDeleteComment(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return auth.Fail(404, "Комментарий не найден.")
	}

	comment, err := s.comments.ByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if comment == nil {
		return auth.Fail(404, "Комментарий не найден.")
	}

	if user.ID != comment.UserID && !canModerateArticles(user) {
		return auth.Fail(403, "Вы можете удалять только свои комментарии.")
	}

	removed, err := s.comments.Delete(ctx.Context(), id)
	if err != nil {
		return err
	}

	if !removed {
		return auth.Fail(404, "Комментарий не найден.")
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}
