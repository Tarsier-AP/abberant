CREATE TABLE article_votes (
    user_id    bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    article_id text   NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    stars      smallint NOT NULL CHECK (stars BETWEEN 1 AND 5),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, article_id)
);

CREATE INDEX article_votes_article_idx ON article_votes (article_id);
