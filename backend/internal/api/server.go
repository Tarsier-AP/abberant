package api

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/config"
	"abberant/backend/internal/session"
)

type Server struct {
	cfg   config.Config
	auth  *auth.Service
	store *session.Store
	app   *fiber.App
}

func New(cfg config.Config, service *auth.Service, store *session.Store) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          errorHandler,
		DisableStartupMessage: false,
	})

	app.Use(logger.New(logger.Config{
		Format: "${status} ${method} ${path}\n",
	}))

	server := &Server{cfg: cfg, auth: service, store: store}

	v1 := app.Group("/api/v1")

	v1.Post("/register", server.handleRegister)
	v1.Post("/register/password", server.handleRegisterPassword)
	v1.Post("/register/activate", server.handleRegisterActivate)
	v1.Post("/activate", server.handleActivate)
	v1.Post("/login", server.handleLogin)
	v1.Post("/login/2fa", server.handleLogin2FA)
	v1.Post("/logout", server.handleLogout)
	v1.Post("/logout/all", server.handleLogoutAll)
	v1.Get("/me", server.requireSession, server.handleMe)
	v1.Patch("/me", server.requireSession, server.handlePatchMe)
	v1.Post("/2fa/totp/setup", server.requireSession, server.handleTotpSetup)
	v1.Post("/2fa/totp/confirm", server.requireSession, server.handleTotpConfirm)
	v1.Post("/2fa/email/send", server.requireAnySession, server.handleEmailSend)
	v1.Post("/2fa/backup/codes/regenerate", server.requireSession, server.handleBackupRegenerate)

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
