CREATE TABLE notifications (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type       text   NOT NULL CHECK (type IN ('reply', 'article')),
    actor_id   bigint REFERENCES users (id) ON DELETE CASCADE,
    article_id text   REFERENCES articles (id) ON DELETE CASCADE,
    comment_id bigint REFERENCES comments (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    read_at    timestamptz
);

CREATE INDEX notifications_user_created_idx ON notifications (user_id, created_at DESC, id);
CREATE INDEX notifications_unread_idx ON notifications (user_id) WHERE read_at IS NULL;

CREATE TABLE section_subscriptions (
    user_id    bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    section    text   NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, section)
);

CREATE INDEX section_subscriptions_section_idx ON section_subscriptions (section);
