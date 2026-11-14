package auth

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID                    int64
	Email                 string
	Username              string
	PasswordHash          string
	Activated             bool
	Role                  string
	Avatar                string
	AdditionalEmail       string
	Phone                 string
	SmsNotifications      bool
	EmailNotifications    bool
	TwoFactorEnabled      bool
	TwoFactorMethod       string
	TotpSecret            string
	RegistrationTokenHash string
	CreatedAt             time.Time
}

const (
	MethodAuthenticator = "authenticator"
	MethodEmail         = "email"
)

const userColumns = `
	id, email, username, password_hash, activated, role, avatar,
	additional_email, phone, sms_notifications, email_notifications,
	two_factor_enabled, two_factor_method, totp_secret, registration_token_hash, created_at
`

type Users struct {
	pool *pgxpool.Pool
}

func NewUsers(pool *pgxpool.Pool) *Users { return &Users{pool: pool} }

func scanUser(row pgx.Row) (*User, error) {
	user := &User{}

	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Activated,
		&user.Role, &user.Avatar, &user.AdditionalEmail, &user.Phone,
		&user.SmsNotifications, &user.EmailNotifications, &user.TwoFactorEnabled,
		&user.TwoFactorMethod, &user.TotpSecret, &user.RegistrationTokenHash, &user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (u *Users) ByEmail(ctx context.Context, email string) (*User, error) {
	return scanUser(u.pool.QueryRow(ctx, `
		SELECT `+userColumns+` FROM users WHERE lower(email) = lower($1)
	`, email))
}

func (u *Users) ByID(ctx context.Context, id int64) (*User, error) {
	return scanUser(u.pool.QueryRow(ctx, `
		SELECT `+userColumns+` FROM users WHERE id = $1
	`, id))
}

func (u *Users) ByRegistrationToken(ctx context.Context, tokenHash string) (*User, error) {
	return scanUser(u.pool.QueryRow(ctx, `
		SELECT `+userColumns+` FROM users
		WHERE registration_token_hash = $1 AND registration_token_hash <> ''
	`, tokenHash))
}

func (u *Users) Create(ctx context.Context, email, username, registrationTokenHash string) (*User, error) {
	return scanUser(u.pool.QueryRow(ctx, `
		INSERT INTO users (email, username, registration_token_hash)
		VALUES ($1, $2, $3)
		RETURNING `+userColumns+`
	`, email, username, registrationTokenHash))
}

func (u *Users) SetRegistrationToken(ctx context.Context, userID int64, tokenHash string) error {
	_, err := u.pool.Exec(ctx, `
		UPDATE users SET registration_token_hash = $2 WHERE id = $1
	`, userID, tokenHash)

	return err
}

func (u *Users) SetPasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	_, err := u.pool.Exec(ctx, `
		UPDATE users SET password_hash = $2 WHERE id = $1
	`, userID, passwordHash)

	return err
}

func (u *Users) Activate(ctx context.Context, userID int64) error {
	_, err := u.pool.Exec(ctx, `
		UPDATE users SET activated = true WHERE id = $1
	`, userID)

	return err
}

func (u *Users) SaveProfile(ctx context.Context, user *User) error {
	_, err := u.pool.Exec(ctx, `
		UPDATE users
		SET username = $2,
		    avatar = $3,
		    additional_email = $4,
		    phone = $5,
		    sms_notifications = $6,
		    email_notifications = $7
		WHERE id = $1
	`, user.ID, user.Username, user.Avatar, user.AdditionalEmail,
		user.Phone, user.SmsNotifications, user.EmailNotifications)

	return err
}

func (u *Users) SetTwoFactor(ctx context.Context, userID int64, enabled bool, method, totpSecret string) error {
	_, err := u.pool.Exec(ctx, `
		UPDATE users
		SET two_factor_enabled = $2, two_factor_method = $3, totp_secret = $4
		WHERE id = $1
	`, userID, enabled, method, totpSecret)

	return err
}

func (u *Users) SetTotpSecret(ctx context.Context, userID int64, secret string) error {
	_, err := u.pool.Exec(ctx, `UPDATE users SET totp_secret = $2 WHERE id = $1`, userID, secret)

	return err
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at <= 0 || at == len(email)-1 {
		return false
	}

	if strings.Contains(email, " ") {
		return false
	}

	return true
}

func ValidUsername(username string) bool {
	if len([]rune(username)) < 3 || len([]rune(username)) > 20 {
		return false
	}

	for _, symbol := range username {
		switch {
		case unicode.IsLetter(symbol), unicode.IsDigit(symbol), symbol == '-', symbol == '_':
		default:
			return false
		}
	}

	return true
}
