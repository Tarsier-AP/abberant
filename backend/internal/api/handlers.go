package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/session"
)

func fail401() *auth.Error {
	return auth.Fail(401, "Требуется вход в аккаунт.")
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type registrationTokenRequest struct {
	RegistrationID string `json:"registrationId"`
}

type registerPasswordRequest struct {
	RegistrationID string `json:"registrationId"`
	Password       string `json:"password"`
}

type codeRequest struct {
	RegistrationID string `json:"registrationId"`
	Email          string `json:"email"`
	Code           string `json:"code"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type twoFactorRequest struct {
	Code string `json:"code"`
}

type patchMeRequest struct {
	Username           *string `json:"username"`
	Avatar             *string `json:"avatar"`
	AdditionalEmail    *string `json:"additionalEmail"`
	Phone              *string `json:"phone"`
	SmsNotifications   *bool   `json:"smsNotifications"`
	EmailNotifications *bool   `json:"emailNotifications"`
	TwoFactorEnabled   *bool   `json:"twoFactorEnabled"`
	TwoFactorMethod    *string `json:"twoFactorMethod"`
}

func decode(ctx *fiber.Ctx, out any) error {
	if err := ctx.BodyParser(out); err != nil {
		return auth.Fail(400, "Некорректный запрос.")
	}

	return nil
}

func devCodeField(code string) fiber.Map {
	if code == "" {
		return fiber.Map{}
	}

	return fiber.Map{"devCode": code}
}

func (s *Server) handleRegister(ctx *fiber.Ctx) error {
	var req registerRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	token, err := s.auth.StartRegistration(ctx.Context(), req.Email, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"status":         "ok",
		"registrationId": token,
	})
}

func (s *Server) handleRegisterPassword(ctx *fiber.Ctx) error {
	var req registerPasswordRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	result, err := s.auth.SetPassword(ctx.Context(), req.RegistrationID, req.Password, auth.ClientIP(ctx.IP()))
	if err != nil {
		return err
	}

	response := fiber.Map{"status": "ok"}
	for key, value := range devCodeField(result.RegistrationDevCode) {
		response[key] = value
	}

	return ctx.JSON(response)
}

func (s *Server) handleRegisterActivate(ctx *fiber.Ctx) error {
	var req codeRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	if err := s.auth.ActivateByToken(ctx.Context(), req.RegistrationID, req.Code, auth.ClientIP(ctx.IP())); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func (s *Server) handleActivate(ctx *fiber.Ctx) error {
	var req codeRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	if err := s.auth.ActivateByEmail(ctx.Context(), req.Email, req.Code, auth.ClientIP(ctx.IP())); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func (s *Server) handleLogin(ctx *fiber.Ctx) error {
	var req loginRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	result, user, err := s.auth.Login(ctx.Context(), req.Email, req.Password, auth.ClientIP(ctx.IP()))
	if err != nil {
		return err
	}

	s.setSessionCookie(ctx, result.SessionToken)

	response := fiber.Map{"status": result.Status}

	if result.Status == "2fa_required" {
		response["twoFactorMethod"] = result.TwoFactorMethod

		if result.Email2FASent {
			response["emailCodeSent"] = true
		}

		if result.Email2FADevCode != "" {
			response["devCode"] = result.Email2FADevCode
		}

		return ctx.JSON(response)
	}

	response["account"] = newAccount(user)

	return ctx.JSON(response)
}

func (s *Server) handleLogin2FA(ctx *fiber.Ctx) error {
	pending, err := s.pendingSession(ctx)
	if err != nil {
		return err
	}

	var req twoFactorRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	upgraded, user, err := s.auth.CompleteTwoFactor(ctx.Context(), pending, req.Code, auth.ClientIP(ctx.IP()))
	if err != nil {
		return err
	}

	ctx.Locals("session", upgraded)
	ctx.Locals("user", user)

	return ctx.JSON(fiber.Map{"status": "ok", "account": newAccount(user)})
}

func (s *Server) pendingSession(ctx *fiber.Ctx) (*session.Session, error) {
	token := s.sessionToken(ctx)
	if token == "" {
		return nil, fail401()
	}

	pending, err := s.store.GetPending(ctx.Context(), token)
	if err != nil {
		return nil, err
	}

	if pending == nil {
		return nil, fail401()
	}

	return pending, nil
}

func (s *Server) handleLogout(ctx *fiber.Ctx) error {
	token := s.sessionToken(ctx)

	if token != "" {
		session, err := s.store.GetActive(ctx.Context(), token)
		if err != nil {
			return err
		}

		if session != nil {
			if err := s.store.Revoke(ctx.Context(), session.ID); err != nil {
				return err
			}
		}
	}

	s.clearSessionCookie(ctx)

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func (s *Server) handleMe(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	return ctx.JSON(fiber.Map{"account": newAccount(user)})
}

func (s *Server) handlePatchMe(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	var req patchMeRequest

	if err := decode(ctx, &req); err != nil {
		return err
	}

	if req.Username != nil {
		username := auth.NormalizeUsername(*req.Username)
		if !auth.ValidUsername(username) {
			return auth.Fail(400, "Имя пользователя должно содержать от 3 до 20 символов: буквы, цифры, дефис или подчёркивание.")
		}

		user.Username = username
	}

	if req.Avatar != nil {
		if !validAvatar(*req.Avatar) {
			return auth.Fail(400, "Некорректное изображение аватара.")
		}

		user.AvatarDataURL = *req.Avatar
	}

	if req.AdditionalEmail != nil {
		email := auth.NormalizeEmail(*req.AdditionalEmail)
		if email != "" && !auth.ValidEmail(email) {
			return auth.Fail(400, "Введите корректный адрес электронной почты.")
		}

		user.AdditionalEmail = email
	}

	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		if len([]rune(phone)) > 32 {
			return auth.Fail(400, "Слишком длинный номер телефона.")
		}

		user.Phone = phone
	}

	if req.SmsNotifications != nil {
		user.SmsNotifications = *req.SmsNotifications
	}

	if req.EmailNotifications != nil {
		user.EmailNotifications = *req.EmailNotifications
	}

	if req.TwoFactorMethod != nil && *req.TwoFactorMethod != user.TwoFactorMethod {
		if err := s.auth.ChangeMethod(ctx.Context(), user, strings.ToLower(*req.TwoFactorMethod)); err != nil {
			return err
		}

		user.TwoFactorMethod = strings.ToLower(*req.TwoFactorMethod)
	}

	if req.TwoFactorEnabled != nil && *req.TwoFactorEnabled != user.TwoFactorEnabled {
		if err := s.changeTwoFactorState(ctx, user, *req.TwoFactorEnabled); err != nil {
			return err
		}
	}

	if err := s.auth.SaveProfile(ctx.Context(), user); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"account": newAccount(user)})
}

func (s *Server) changeTwoFactorState(ctx *fiber.Ctx, user *auth.User, enabled bool) error {
	if !enabled {
		if err := s.auth.Disable2FA(ctx.Context(), user.ID); err != nil {
			return err
		}

		user.TwoFactorEnabled = false
		user.TwoFactorMethod = ""
		user.TotpSecret = ""

		return nil
	}

	method := user.TwoFactorMethod

	if err := s.auth.Enable2FA(ctx.Context(), user, method); err != nil {
		return err
	}

	user.TwoFactorEnabled = true

	return nil
}

func validAvatar(dataURL string) bool {
	const maxAvatarLength = 6 * 1024 * 1024

	if dataURL == "" {
		return true
	}

	return strings.HasPrefix(dataURL, "data:image/") && len(dataURL) <= maxAvatarLength
}

func (s *Server) currentUser(ctx *fiber.Ctx) *auth.User {
	user, _ := ctx.Locals("user").(*auth.User)

	return user
}
