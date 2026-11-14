package api

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"abberant/backend/internal/auth"
	"abberant/backend/internal/session"
)

type Account struct {
	Login              string `json:"login"`
	Username           string `json:"username"`
	Avatar             string `json:"avatar"`
	AdditionalEmail    string `json:"additionalEmail"`
	Phone              string `json:"phone"`
	Activated          bool   `json:"activated"`
	TwoFactorEnabled   bool   `json:"twoFactorEnabled"`
	TwoFactorMethod    string `json:"twoFactorMethod"`
	SmsNotifications   bool   `json:"smsNotifications"`
	EmailNotifications bool   `json:"emailNotifications"`
	Role               string `json:"role"`
}

func newAccount(user *auth.User) *Account {
	return &Account{
		Login:              user.Email,
		Username:           user.Username,
		Avatar:             user.AvatarDataURL,
		AdditionalEmail:    user.AdditionalEmail,
		Phone:              user.Phone,
		Activated:          user.Activated,
		TwoFactorEnabled:   user.TwoFactorEnabled,
		TwoFactorMethod:    user.TwoFactorMethod,
		SmsNotifications:   user.SmsNotifications,
		EmailNotifications: user.EmailNotifications,
		Role:               user.Role,
	}
}

func (s *Server) setSessionCookie(ctx *fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     session.CookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.cfg.CookieSecure,
		SameSite: "Lax",
		MaxAge:   int(s.cfg.SessionTTL / time.Second),
	})
}

func (s *Server) clearSessionCookie(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     session.CookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.cfg.CookieSecure,
		SameSite: "Lax",
		MaxAge:   -1,
	})
}

func (s *Server) sessionToken(ctx *fiber.Ctx) string {
	return ctx.Cookies(session.CookieName)
}

func (s *Server) optionalSession(ctx *fiber.Ctx) error {
	session, user, err := s.resolve(ctx)
	if err != nil {
		return err
	}

	if session != nil && !session.Pending {
		ctx.Locals("session", session)
		ctx.Locals("user", user)
	}

	return ctx.Next()
}

func (s *Server) requireSession(ctx *fiber.Ctx) error {
	session, user, err := s.resolve(ctx)
	if err != nil {
		return err
	}

	if session == nil || session.Pending {
		return fail401()
	}

	ctx.Locals("session", session)
	ctx.Locals("user", user)

	return ctx.Next()
}

func (s *Server) requireEditor(ctx *fiber.Ctx) error {
	user := s.currentUser(ctx)

	if user == nil {
		return fail401()
	}

	if !auth.CanEditArticles(user.Role) {
		return auth.Fail(403, "Недостаточно прав для управления статьями.")
	}

	return ctx.Next()
}

func (s *Server) requireAnySession(ctx *fiber.Ctx) error {
	session, user, err := s.resolve(ctx)
	if err != nil {
		return err
	}

	if session == nil {
		return fail401()
	}

	ctx.Locals("session", session)
	ctx.Locals("user", user)

	return ctx.Next()
}

func (s *Server) resolve(ctx *fiber.Ctx) (*session.Session, *auth.User, error) {
	token := s.sessionToken(ctx)
	if token == "" {
		return nil, nil, nil
	}

	var (
		found *session.Session
		err   error
	)

	if found, err = s.store.GetActive(ctx.Context(), token); err != nil {
		return nil, nil, err
	}

	if found != nil {
		user, err := s.userByID(ctx, found.UserID)
		if err != nil {
			return nil, nil, err
		}

		return found, user, nil
	}

	if found, err = s.store.GetPending(ctx.Context(), token); err != nil {
		return nil, nil, err
	}

	if found == nil {
		return nil, nil, nil
	}

	user, err := s.userByID(ctx, found.UserID)
	if err != nil {
		return nil, nil, err
	}

	return found, user, nil
}

func (s *Server) userByID(ctx *fiber.Ctx, id int64) (*auth.User, error) {
	user, err := s.auth.UserByID(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fail401()
	}

	return user, nil
}
