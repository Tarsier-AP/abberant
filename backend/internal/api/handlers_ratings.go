package api

import (
	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
	"abberant/backend/internal/ratings"
)

type voteRequest struct {
	Stars int `json:"stars"`
}

func (s *Server) handleVote(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	article, err := s.articles.ByID(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	if article == nil || !canSeeArticle(article, user) {
		return auth.Fail(404, "Статья не найдена.")
	}

	if !articles.Rateable(article) {
		return auth.Fail(400, "Оценивать можно только обзоры и прохождения.")
	}

	var req voteRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	if !articles.ValidRating(req.Stars) {
		return auth.Fail(400, "Оценка должна быть от 1 до 5.")
	}

	stats, err := s.ratings.Put(ctx.Context(), user.ID, article.ID, req.Stars)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"rating": newReaderRatingResponse(stats, req.Stars)})
}

func (s *Server) handleRemoveVote(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	article, err := s.articles.ByID(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	if article == nil || !canSeeArticle(article, user) {
		return auth.Fail(404, "Статья не найдена.")
	}

	removed, err := s.ratings.Delete(ctx.Context(), user.ID, article.ID)
	if err != nil {
		return err
	}

	if !removed {
		return auth.Fail(404, "Оценка не найдена.")
	}

	stats, err := s.ratings.StatsByArticle(ctx.Context(), article.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"rating": newReaderRatingResponse(stats, 0)})
}

func newReaderRatingResponse(stats ratings.Stats, myVote int) *readerRatingResponse {
	return &readerRatingResponse{
		Average: stats.Average,
		Votes:   stats.Votes,
		MyVote:  myVote,
	}
}
