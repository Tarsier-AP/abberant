package api

import (
	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
)

// bookmarkedResponses turns the bookmark rows into the same article response
// shape as the public list, with live comment counts.
func (s *Server) bookmarkedResponses(ctx *fiber.Ctx, user *auth.User) ([]articleResponse, error) {
	list, err := s.bookmarks.Articles(ctx.Context(), user.ID)
	if err != nil {
		return nil, err
	}

	counts, err := s.comments.CountsByArticle(ctx.Context())
	if err != nil {
		return nil, err
	}

	response := make([]articleResponse, 0, len(list))

	for _, article := range list {
		if !canSeeArticle(article, user) {
			continue
		}

		item := newArticleResponse(article)

		if count, ok := counts[article.ID]; ok {
			item.Comments = count
		}

		response = append(response, item)
	}

	return response, nil
}

func (s *Server) handleListBookmarks(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	response, err := s.bookmarkedResponses(ctx, user)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"articles": response})
}

func (s *Server) handleListBookmarkIDs(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	ids, err := s.bookmarks.IDs(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	if ids == nil {
		ids = []string{}
	}

	return ctx.JSON(fiber.Map{"ids": ids})
}

func (s *Server) handleAddBookmark(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	articleID := ctx.Params("articleId")

	article, err := s.articles.ByID(ctx.Context(), articleID)
	if err != nil {
		return err
	}

	if article == nil || !canSeeArticle(article, user) {
		return auth.Fail(404, "Статья не найдена.")
	}

	created, err := s.bookmarks.Add(ctx.Context(), user.ID, article.ID)
	if err != nil {
		return err
	}

	status := 200
	if created {
		status = 201
	}

	return ctx.Status(status).JSON(fiber.Map{"status": "ok"})
}

func (s *Server) handleRemoveBookmark(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	if _, err := s.bookmarks.Remove(ctx.Context(), user.ID, ctx.Params("articleId")); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}
