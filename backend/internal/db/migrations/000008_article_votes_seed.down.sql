-- The example commenters have no password hash and cannot log in, so every
-- vote of theirs originates from the seed above.
DELETE FROM article_votes
WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.abberant'
);
