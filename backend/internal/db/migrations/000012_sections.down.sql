DROP TABLE IF EXISTS sections;

ALTER TABLE articles DROP CONSTRAINT articles_kind_check;
ALTER TABLE articles ADD CONSTRAINT articles_kind_check CHECK (kind IN ('news', 'review'));
