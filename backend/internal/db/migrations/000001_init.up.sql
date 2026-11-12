CREATE TABLE users (
    id                      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email                   text        NOT NULL,
    username                text        NOT NULL,
    password_hash           text        NOT NULL DEFAULT '',
    activated               boolean     NOT NULL DEFAULT false,
    role                    text        NOT NULL DEFAULT 'reader',
    avatar_data_url         text        NOT NULL DEFAULT '',
    additional_email        text        NOT NULL DEFAULT '',
    phone                   text        NOT NULL DEFAULT '',
    sms_notifications       boolean     NOT NULL DEFAULT true,
    email_notifications     boolean     NOT NULL DEFAULT true,
    two_factor_enabled      boolean     NOT NULL DEFAULT false,
    two_factor_method       text        NOT NULL DEFAULT '',
    totp_secret             text        NOT NULL DEFAULT '',
    registration_token_hash text        NOT NULL DEFAULT '',
    created_at              timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_lower_key ON users (lower(email));

CREATE TABLE sessions (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      bigint      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash   text        NOT NULL UNIQUE,
    pending      boolean     NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL,
    revoked_at   timestamptz
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TABLE codes (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    bigint REFERENCES users (id) ON DELETE CASCADE,
    email      text        NOT NULL,
    kind       text        NOT NULL CHECK (kind IN ('activation', 'twofactor')),
    code_hash  text        NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at    timestamptz
);

CREATE INDEX codes_email_kind_idx ON codes (lower(email), kind);

CREATE TABLE totp_backup_codes (
    id       bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id  bigint      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    code_hash text       NOT NULL,
    used_at  timestamptz
);

CREATE INDEX totp_backup_codes_user_id_idx ON totp_backup_codes (user_id);

CREATE TABLE login_attempts (
    key          text PRIMARY KEY,
    fail_count   int         NOT NULL DEFAULT 0,
    last_fail_at timestamptz NOT NULL DEFAULT now()
);
