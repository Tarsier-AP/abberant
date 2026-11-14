package auth

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"

	"abberant/backend/internal/config"
	"abberant/backend/internal/mailer"
	"abberant/backend/internal/session"
)

type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

func Fail(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}

func fail(status int, message string) *Error {
	return Fail(status, message)
}

var (
	ErrInvalidCredentials  = fail(401, "Неверный логин или пароль.")
	ErrNotActivated        = fail(403, "Аккаунт не активирован. Проверьте электронную почту.")
	ErrInvalidCode         = fail(401, "Неверный код.")
	ErrInvalidRegistration = fail(400, "Регистрация не найдена или устарела. Зарегистрируйтесь заново.")
	ErrTooManyAttempts     = fail(429, "Слишком много попыток. Повторите через несколько минут.")
	ErrSMSUnavailable      = fail(400, "Подтверждение по SMS пока недоступно.")
	ErrAuthenticatorSetup  = fail(400, "Подтвердите настройку приложения-аутентификатора кодом.")
	ErrMethodRequired      = fail(400, "Выберите способ подтверждения.")
	ErrBad2FAMethod        = fail(400, "Неизвестный способ подтверждения.")
)

type Service struct {
	pool     *pgxpool.Pool
	users    *Users
	codes    *Codes
	totp     *TOTP
	attempts *Attempts
	sessions *session.Store
	mailer   mailer.Mailer
	cfg      config.Config
}

func NewService(
	pool *pgxpool.Pool,
	sessions *session.Store,
	mailer mailer.Mailer,
	cfg config.Config,
) *Service {
	return &Service{
		pool:     pool,
		users:    NewUsers(pool),
		codes:    NewCodes(pool, cfg.CodeTTL),
		totp:     NewTOTP(pool),
		attempts: NewAttempts(pool),
		sessions: sessions,
		mailer:   mailer,
		cfg:      cfg,
	}
}

func (s *Service) SaveProfile(ctx context.Context, user *User) error {
	return s.users.SaveProfile(ctx, user)
}

func (s *Service) UserByID(ctx context.Context, id int64) (*User, error) {
	return s.users.ByID(ctx, id)
}

type LoginResult struct {
	Status              string
	SessionToken        string
	TwoFactorMethod     string
	Email2FADevCode     string
	Email2FASent        bool
	RegistrationDevCode string
}

func ClientIP(remote string) string {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		return remote
	}

	return host
}

func (s *Service) attemptKeys(email, ip string) []string {
	keys := []string{"ip:" + ip}

	if email != "" {
		keys = append(keys, "login:"+email)
	}

	return keys
}

func (s *Service) StartRegistration(ctx context.Context, email, username string) (string, error) {
	email = NormalizeEmail(email)
	username = NormalizeUsername(username)

	if !ValidEmail(email) {
		return "", fail(400, "Введите корректный адрес электронной почты.")
	}

	if !ValidUsername(username) {
		return "", fail(400, "Имя пользователя должно содержать от 3 до 20 символов: буквы, цифры, дефис или подчёркивание.")
	}

	existing, err := s.users.ByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	token, tokenHash, err := session.NewToken()
	if err != nil {
		return "", err
	}

	switch {
	case existing == nil:
		if _, err := s.users.Create(ctx, email, username, tokenHash); err != nil {
			return "", err
		}

	case existing.Activated:
		// A finished account cannot be registered again: report it instead of
		// issuing a registration token that leads nowhere.
		return "", fail(409, "Аккаунт с этой электронной почтой уже зарегистрирован.")

	default:
		// Not activated yet: restart the registration (a fresh password and
		// activation code will be issued) instead of leaving a dead end.
		if err := s.users.SetRegistrationToken(ctx, existing.ID, tokenHash); err != nil {
			return "", err
		}
	}

	return token, nil
}

func (s *Service) SetPassword(ctx context.Context, registrationToken, password, ip string) (*LoginResult, error) {
	user, err := s.userByRegistrationToken(ctx, registrationToken)
	if err != nil {
		return nil, err
	}

	if !ValidatePassword(password) {
		return nil, fail(400, "Пароль должен содержать не менее 8 символов, заглавную букву, цифру и специальный символ.")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	if err := s.users.SetPasswordHash(ctx, user.ID, hash); err != nil {
		return nil, err
	}

	return s.sendActivationCode(ctx, user)
}

func (s *Service) userByRegistrationToken(ctx context.Context, registrationToken string) (*User, error) {
	if registrationToken == "" {
		return nil, ErrInvalidRegistration
	}

	user, err := s.users.ByRegistrationToken(ctx, session.HashToken(registrationToken))
	if err != nil {
		return nil, err
	}

	if user == nil || user.Activated {
		return nil, ErrInvalidRegistration
	}

	return user, nil
}

func (s *Service) sendActivationCode(ctx context.Context, user *User) (*LoginResult, error) {
	code, err := s.codes.Issue(ctx, user.ID, user.Email, KindActivation)
	if err != nil {
		return nil, err
	}

	if err := s.mailer.Send(
		user.Email,
		"Код активации Abberant",
		fmt.Sprintf(
			"Здравствуйте!\n\nВаш код активации аккаунта Abberant: %s\nКод действует 5 минут и может быть использован один раз.\n\nЕсли вы не регистрировались — просто проигнорируйте это письмо.",
			code,
		),
	); err != nil {
		return nil, fmt.Errorf("send activation code: %w", err)
	}

	log.Printf("[auth] activation code issued for user %d (%s)", user.ID, user.Email)

	return &LoginResult{RegistrationDevCode: s.devCode(code)}, nil
}

func (s *Service) devCode(code string) string {
	if s.cfg.DeliveryDevMode() {
		return code
	}

	return ""
}

// ResendActivationCode issues and sends a fresh activation code for a
// pending registration. It mirrors SendEmail2FACode: the attempt keys
// throttle the endpoint after failures and are cleared on success.
func (s *Service) ResendActivationCode(ctx context.Context, registrationToken, ip string) (*LoginResult, error) {
	user, err := s.userByRegistrationToken(ctx, registrationToken)
	if err != nil {
		return nil, err
	}

	keys := s.attemptKeys(user.Email, ip)
	keys = append(keys, fmt.Sprintf("activation-send:%d", user.ID))

	if err := s.checkAttempts(ctx, keys); err != nil {
		return nil, err
	}

	result, err := s.sendActivationCode(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := s.clearAttempts(ctx, keys); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) ActivateByToken(ctx context.Context, registrationToken, code, ip string) error {
	user, err := s.userByRegistrationToken(ctx, registrationToken)
	if err != nil {
		return err
	}

	return s.activate(ctx, user, code, ip)
}

func (s *Service) ActivateByEmail(ctx context.Context, email, code, ip string) error {
	email = NormalizeEmail(email)

	user, err := s.users.ByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil || user.Activated || user.PasswordHash == "" {
		return ErrInvalidCode
	}

	return s.activate(ctx, user, code, ip)
}

func (s *Service) activate(ctx context.Context, user *User, code, ip string) error {
	if err := s.checkAttempts(ctx, s.attemptKeys(user.Email, ip)); err != nil {
		return err
	}

	ok, err := s.codes.Verify(ctx, user.ID, KindActivation, code)
	if err != nil {
		return err
	}

	if !ok {
		if err := s.recordAttempts(ctx, s.attemptKeys(user.Email, ip)); err != nil {
			return err
		}

		return ErrInvalidCode
	}

	if err := s.users.Activate(ctx, user.ID); err != nil {
		return err
	}

	if err := s.clearAttempts(ctx, s.attemptKeys(user.Email, ip)); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password, ip string) (*LoginResult, *User, error) {
	email = NormalizeEmail(email)

	if err := s.checkAttempts(ctx, s.attemptKeys(email, ip)); err != nil {
		return nil, nil, err
	}

	user, err := s.users.ByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if user == nil || user.PasswordHash == "" || !VerifyPassword(password, user.PasswordHash) {
		if err := s.recordAttempts(ctx, s.attemptKeys(email, ip)); err != nil {
			return nil, nil, err
		}

		return nil, nil, ErrInvalidCredentials
	}

	if err := s.clearAttempts(ctx, s.attemptKeys(email, ip)); err != nil {
		return nil, nil, err
	}

	if !user.Activated {
		return nil, nil, ErrNotActivated
	}

	if user.TwoFactorEnabled {
		token, _, err := s.sessions.Create(ctx, user.ID, true)
		if err != nil {
			return nil, nil, err
		}

		result := &LoginResult{
			Status:          "2fa_required",
			SessionToken:    token,
			TwoFactorMethod: user.TwoFactorMethod,
		}

		if user.TwoFactorMethod == MethodEmail {
			code, err := s.issueEmailCode(ctx, user)
			if err != nil {
				return nil, nil, err
			}

			result.Email2FASent = true
			result.Email2FADevCode = s.devCode(code)
		}

		return result, user, nil
	}

	token, _, err := s.sessions.Create(ctx, user.ID, false)
	if err != nil {
		return nil, nil, err
	}

	return &LoginResult{Status: "ok", SessionToken: token}, user, nil
}

func (s *Service) CompleteTwoFactor(ctx context.Context, pending *session.Session, code, ip string) (*session.Session, *User, error) {
	if err := s.checkAttempts(ctx, s.attemptKeys("", ip)); err != nil {
		return nil, nil, err
	}

	user, err := s.users.ByID(ctx, pending.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil || !user.Activated {
		return nil, nil, ErrInvalidCredentials
	}

	keys := s.attemptKeys(user.Email, ip)
	keys = append(keys, fmt.Sprintf("2fa:%d", user.ID))

	if err := s.checkAttempts(ctx, keys); err != nil {
		return nil, nil, err
	}

	verified := false

	switch user.TwoFactorMethod {
	case MethodAuthenticator:
		verified = s.totp.Verify(user.TotpSecret, code)
	case MethodEmail:
		verified, err = s.codes.Verify(ctx, user.ID, KindTwoFactor, code)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, ErrBad2FAMethod
	}

	if !verified {
		var backupErr error

		verified, backupErr = s.totp.ConsumeBackupCode(ctx, user.ID, code)
		if backupErr != nil {
			return nil, nil, backupErr
		}
	}

	if !verified {
		if err := s.recordAttempts(ctx, keys); err != nil {
			return nil, nil, err
		}

		return nil, nil, ErrInvalidCode
	}

	if err := s.sessions.Upgrade(ctx, pending.ID); err != nil {
		return nil, nil, err
	}

	if err := s.clearAttempts(ctx, keys); err != nil {
		return nil, nil, err
	}

	upgraded := &session.Session{
		ID: pending.ID, UserID: pending.UserID, Pending: false,
	}

	return upgraded, user, nil
}

func (s *Service) SendEmail2FACode(ctx context.Context, user *User, ip string) (string, error) {
	keys := s.attemptKeys(user.Email, ip)
	keys = append(keys, fmt.Sprintf("2fa-send:%d", user.ID))

	if err := s.checkAttempts(ctx, keys); err != nil {
		return "", err
	}

	code, err := s.issueEmailCode(ctx, user)
	if err != nil {
		return "", err
	}

	if err := s.clearAttempts(ctx, keys); err != nil {
		return "", err
	}

	return s.devCode(code), nil
}

func (s *Service) issueEmailCode(ctx context.Context, user *User) (string, error) {
	code, err := s.codes.Issue(ctx, user.ID, user.Email, KindTwoFactor)
	if err != nil {
		return "", err
	}

	if err := s.mailer.Send(
		user.Email,
		"Код подтверждения входа Abberant",
		fmt.Sprintf(
			"Ваш код для подтверждения входа: %s\nКод действует 5 минут и может быть использован один раз.",
			code,
		),
	); err != nil {
		return "", fmt.Errorf("send 2fa code: %w", err)
	}

	return code, nil
}

func (s *Service) SetupTOTP(ctx context.Context, user *User) (*Enrollment, error) {
	enrollment, err := s.totp.GenerateSecret(user.Email)
	if err != nil {
		return nil, err
	}

	if err := s.users.SetTotpSecret(ctx, user.ID, enrollment.Secret); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s *Service) ConfirmTOTP(ctx context.Context, user *User, code string) ([]string, error) {
	secret, err := s.users.ByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.TotpSecret == "" {
		return nil, fail(400, "Настройка не начата. Сгенерируйте QR-код заново.")
	}

	if !s.totp.Verify(secret.TotpSecret, code) {
		return nil, ErrInvalidCode
	}

	codes, err := GenerateBackupCodes(BackupCodeCount)
	if err != nil {
		return nil, err
	}

	if err := s.totp.StoreBackupCodes(ctx, user.ID, codes); err != nil {
		return nil, err
	}

	if err := s.users.SetTwoFactor(ctx, user.ID, true, MethodAuthenticator, secret.TotpSecret); err != nil {
		return nil, err
	}

	return codes, nil
}

func (s *Service) Disable2FA(ctx context.Context, userID int64) error {
	if err := s.users.SetTwoFactor(ctx, userID, false, "", ""); err != nil {
		return err
	}

	_, err := s.pool.Exec(ctx, `DELETE FROM totp_backup_codes WHERE user_id = $1 AND used_at IS NULL`, userID)

	return err
}

func (s *Service) ChangeMethod(ctx context.Context, user *User, method string) error {
	if method == "sms" {
		return ErrSMSUnavailable
	}

	if method != MethodAuthenticator && method != MethodEmail {
		return ErrBad2FAMethod
	}

	if method == MethodAuthenticator && user.TotpSecret == "" {
		return ErrAuthenticatorSetup
	}

	return s.users.SetTwoFactor(ctx, user.ID, user.TwoFactorEnabled, method, user.TotpSecret)
}

func (s *Service) Enable2FA(ctx context.Context, user *User, method string) error {
	switch method {
	case MethodEmail:
		return s.users.SetTwoFactor(ctx, user.ID, true, MethodEmail, "")
	case MethodAuthenticator:
		return ErrAuthenticatorSetup
	case "sms":
		return ErrSMSUnavailable
	default:
		return ErrMethodRequired
	}
}

func (s *Service) RegenerateBackupCodes(ctx context.Context, user *User) ([]string, error) {
	if !user.TwoFactorEnabled {
		return nil, fail(400, "Сначала включите двухфакторную аутентификацию.")
	}

	codes, err := GenerateBackupCodes(BackupCodeCount)
	if err != nil {
		return nil, err
	}

	if err := s.totp.StoreBackupCodes(ctx, user.ID, codes); err != nil {
		return nil, err
	}

	return codes, nil
}

func (s *Service) checkAttempts(ctx context.Context, keys []string) error {
	for _, key := range keys {
		wait, err := s.attempts.RetryAfter(ctx, key)
		if err != nil {
			return err
		}

		if wait > 0 {
			return ErrTooManyAttempts
		}
	}

	return nil
}

func (s *Service) recordAttempts(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := s.attempts.Record(ctx, key); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) clearAttempts(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := s.attempts.Clear(ctx, key); err != nil {
			return err
		}
	}

	return nil
}
