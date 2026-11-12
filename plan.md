# Plan: Authorization with 2FA support

Status: agreed design (grilled 12.11.2026). Backend does not exist yet; `backend/` is empty. Frontend has a full mocked auth flow that will be rewired to the real API.

## Scope

- First real backend implementation: Go + Fiber, PostgreSQL, behind nginx.
- Rewire the frontend (`src/stores/auth.ts` and auth views) off `localStorage` mocks to real API calls.
- Russian-language UI; all new user-facing copy in Russian.

## Decisions (from the interview)

| Topic | Decision |
| --- | --- |
| Scope | Backend + frontend rewiring |
| 2FA methods at launch | TOTP (authenticator app) + email codes; SMS shown as «скоро» (disabled), rejected at API level |
| Auth identifier | Email only; phone becomes a profile-only field |
| Sessions | Server-side sessions in Postgres, opaque token in HttpOnly cookie, revoke-all supported |
| Session lifetime | Expires after 7 days of inactivity (sliding renewal) |
| Registration flow | Replicate the mocked flow exactly: register → set password → activation code by email → activate → login |
| Activation delivery | Email via SMTP (env-configured); if SMTP unset, codes are logged server-side (dev mode) |
| Passwords | Argon2id (pure Go), minimum 8 characters |
| Anti-enumeration | Identical login errors («Неверный логин или пароль»), ambiguous registration responses |
| Brute force | Full hardening: single-use codes, expiry, constant-time compares, per-account/per-IP attempt counters with backoff |
| TOTP semantics | Standard timing: 30-second windows, ±1 step drift (~90 s validity) |
| DB on Windows dev box | `embedded-postgres` — real PG binaries run as a child process; no Docker, no admin rights needed; same code path on Linux |
| TOTP enrollment UX | Secret → QR (data-URL, no external CDN) → confirm with a live code → 2FA activates; ~8 one-time backup codes |
| Backend layout | Standard: `cmd/server`, `internal/{config,db,mailer,auth,session,api}`, embedded SQL migrations, `/api/v1` prefix |
| Plumbing | Vite dev proxy for `/api` + nginx reverse-proxy config in the repo |

## Backend

### Layout

```
backend/
  cmd/server/main.go
  internal/
    config/    — env-driven config (DSN, SMTP, cookie flags, dev mode)
    db/        — embedded-postgres bootstrap + migrations (golang-migrate, embedded)
    mailer/    — SMTP sender with log fallback
    auth/      — registration, activation, login, TOTP, email codes, backup codes
    session/   — server-side session store, sliding expiry, revoke-all
    api/       — Fiber routes, middleware, /api/v1 handlers
  migrations/  — SQL files
  go.mod
```

### Database schema (initial migration)

- `users`: id, email (unique, citext or lower-indexed), username, password_hash, activated, avatar_data_url, additional_email, sms_notifications, email_notifications, two_factor_enabled, two_factor_method, totp_secret (encrypted or plain-at-rest for now), created_at
- `sessions`: id, user_id, token_hash, created_at, last_seen_at (sliding 7-day expiry), revoked_at
- `codes`: id, user_id (nullable for pre-registration), email, kind (activation | twofactor), code_hash, expires_at, used_at — single-use, 5-minute expiry
- `login_attempts`: account key + IP, counters for backoff
- `totp_backup_codes`: user_id, code_hash, used_at

### API (`/api/v1`)

Auth:
- `POST /register` — start registration (email + username); ambiguous success
- `POST /register/password` — set password for pending registration
- `POST /register/activate` — submit activation code
- `POST /activate` — activation code check (mirrors mock flow)
- `POST /login` — password step; returns `2fa_required` + method when enabled
- `POST /login/2fa` — TOTP code / email code / backup code
- `POST /logout`
- `GET /me`, `PATCH /me` — account profile, settings (2FA switch, method, notifications)

TOTP enrollment:
- `POST /2fa/totp/setup` — generate secret, return QR as data-URL + secret
- `POST /2fa/totp/confirm` — verify a live code, activate 2FA, return backup codes
- `POST /2fa/email/send` — send an email 2FA code on demand
- `POST /2fa/backup/codes/regenerate`

Security invariants:
- Codes: single-use, 5-minute expiry, constant-time comparison, stored hashed.
- Attempt counters per account and per IP with backoff (see note below).
- Login never reveals account existence or 2FA state until password is correct.
- Passwords: Argon2id; minimum 8 chars.

## Frontend changes

- `src/stores/auth.ts`: replace localStorage mock with API client calls; keep the same public shape where possible so views need minimal edits.
- Login/2FA/register/activation views: submit to real endpoints, handle `2fa_required`, error states, dev-mode code hinting.
- Settings: SMS method shown but disabled («скоро»); new TOTP enrollment screen (QR + live-code confirm + backup codes display); «выйти со всех устройств» button.
- Registration form: email-only; phone stays editable as a profile field in settings.
- New copy in Russian.

## Plumbing

- `frontend/vite.config.ts`: dev proxy `/api` → `http://localhost:8080`.
- `nginx` conf in repo: serve frontend statics, proxy `/api` to the Fiber server.
- `docker-compose.yml` optional later for Linux/prod; dev works without it.

## Open point / correction watch

- Rate-limit wording: hardening package stands as chosen; «1 minute for TOTP as usual» is interpreted as standard TOTP timing (30 s windows, ±1 step drift). If a 1-minute lockout after failed attempts is intended instead, encode that in `login_attempts` and confirm.

## Build order

1. Backend skeleton: go.mod, config, embedded-postgres bootstrap, migrations.
2. Mailer (SMTP + log fallback), session store.
3. Auth endpoints: register → password → activation → login (+ anti-enumeration semantics).
4. TOTP + email 2FA, backup codes, enrollment endpoints.
5. Frontend rewiring (store, views, settings enrollment screen).
6. Vite proxy + nginx conf; verify `npm run build` and lint; commit.
