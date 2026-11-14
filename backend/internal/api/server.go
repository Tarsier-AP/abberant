package api

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
	"abberant/backend/internal/comments"
	"abberant/backend/internal/config"
	"abberant/backend/internal/session"
)

type Server struct {
	cfg      config.Config
	auth     *auth.Service
	store    *session.Store
	articles *articles.Repo
	comments *comments.Repo
	app      *fiber.App
}

func New(cfg config.Config, service *auth.Service, store *session.Store, articlesRepo *articles.Repo, commentsRepo *comments.Repo) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          errorHandler,
		DisableStartupMessage: false,
	})

	app.Use(logger.New(logger.Config{
		Format: "${status} ${method} ${path}\n",
	}))

	server := &Server{cfg: cfg, auth: service, store: store, articles: articlesRepo, comments: commentsRepo}

	v1 := app.Group("/api/v1")

	v1.Post("/register", server.handleRegister)
	v1.Post("/register/password", server.handleRegisterPassword)
	v1.Post("/register/activate", server.handleRegisterActivate)
	v1.Post("/register/resend", server.handleRegisterResend)
	v1.Post("/activate", server.handleActivate)
	v1.Post("/login", server.handleLogin)
	v1.Post("/login/2fa", server.handleLogin2FA)
	v1.Post("/logout", server.handleLogout)
	v1.Get("/me", server.handleMe)
	v1.Patch("/me", server.requireSession, server.handlePatchMe)
	v1.Post("/2fa/totp/setup", server.requireSession, server.handleTotpSetup)
	v1.Post("/2fa/totp/confirm", server.requireSession, server.handleTotpConfirm)
	v1.Post("/2fa/email/send", server.requireAnySession, server.handleEmailSend)
	v1.Post("/2fa/backup/codes/regenerate", server.requireSession, server.handleBackupRegenerate)

	v1.Get("/articles", server.handleListArticles)
	v1.Get("/articles/:id/comments", server.handleListComments)
	v1.Post("/articles/:id/comments", server.requireSession, server.handleCreateComment)
	v1.Delete("/comments/:id", server.requireSession, server.handleDeleteComment)
	v1.Post("/articles", server.requireSession, server.requireEditor, server.handleCreateArticle)
	v1.Patch("/articles/:id", server.requireSession, server.requireEditor, server.handleUpdateArticle)
	v1.Delete("/articles/:id", server.requireSession, server.handleDeleteArticle)

	return app
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	if apiErr, ok := err.(*auth.Error); ok {
		return ctx.Status(apiErr.Status).JSON(fiber.Map{"error": apiErr.Message})
	}

	var fiberErr *fiber.Error

	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(fiber.Map{"error": "Запрос не обработан."})
	}

	log.Printf("[api] internal error: %v", err)

	return ctx.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка сервера. Повторите попытку позже."})
}
