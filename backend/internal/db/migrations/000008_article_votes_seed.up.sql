-- Reader grades for the seeded content: every published article gets votes
-- from the example commenters (deactivated accounts that can never log in,
-- so all of their votes come from this seed). Stars are hashed per
-- (article, user) pair: ~45% fives, ~35% fours, ~15% threes, ~5% twos, and
-- each article is voted on by a varying subset of the commenters.

INSERT INTO article_votes (user_id, article_id, stars)
SELECT user_id,
       article_id,
       CASE
           WHEN h < 45 THEN 5
           WHEN h < 80 THEN 4
           WHEN h < 95 THEN 3
           ELSE 2
       END
FROM (
    SELECT u.id   AS user_id,
           a.id   AS article_id,
           abs(hashtext(a.id || u.username || 'stars')) % 100 AS h
    FROM articles a
    JOIN users u ON u.email LIKE '%@example.abberant'
    WHERE a.status = 'published'
      AND abs(hashtext(a.id || u.username)) % 10 < 6
) picks
ON CONFLICT (user_id, article_id) DO NOTHING;
