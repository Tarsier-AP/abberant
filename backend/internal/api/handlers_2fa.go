package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/session"
)

type codeOnlyRequest struct {
	Code string `json:"code"`
}

func (s *Server) handleTotpSetup(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	if !user.Activated {
		return auth.Fail(403, "Активируйте аккаунт, чтобы включить двухфакторную аутентификацию.")
	}

	enrollment, err := s.auth.SetupTOTP(ctx.Context(), user)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"status":     "ok",
		"secret":     enrollment.Secret,
		"otpauthUrl": enrollment.OtpauthURL,
		"qrDataURL":  enrollment.QRDataURL,
	})
}

func (s *Server) handleTotpConfirm(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	var req codeOnlyRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	codes, err := s.auth.ConfirmTOTP(ctx.Context(), user, strings.TrimSpace(req.Code))
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok", "backupCodes": codes})
}

func (s *Server) handleEmailSend(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)
	pending, _ := ctx.Locals("session").(*session.Session)

	if !user.Activated {
		return auth.Fail(403, "Активируйте аккаунт, чтобы включить двухфакторную аутентификацию.")
	}

	if pending != nil && pending.Pending && user.TwoFactorMethod != auth.MethodEmail {
		return auth.Fail(400, "Способ подтверждения — не электронная почта.")
	}

	devCode, err := s.auth.SendEmail2FACode(ctx.Context(), user, auth.ClientIP(ctx.IP()))
	if err != nil {
		return err
	}

	response := fiber.Map{"status": "ok", "emailCodeSent": true}
	if devCode != "" {
		response["devCode"] = devCode
	}

	return ctx.JSON(response)
}

func (s *Server) handleBackupRegenerate(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	codes, err := s.auth.RegenerateBackupCodes(ctx.Context(), user)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok", "backupCodes": codes})
}
