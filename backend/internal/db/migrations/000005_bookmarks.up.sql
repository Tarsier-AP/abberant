CREATE TABLE bookmarks (
    user_id    bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    article_id text   NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, article_id)
);

CREATE INDEX bookmarks_user_created_idx ON bookmarks (user_id, created_at DESC);
