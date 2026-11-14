-- Server-managed sections: news platform feeds and walkthrough game hubs.
-- Review groups stay static frontend config (they are feed groupings, not
-- article assignments). The `recent` pseudo-sections remain frontend-only.
-- Hubs are soft-deleted: `archived` hides them from pickers and public feeds
-- while assigned articles keep their section and their URLs.

CREATE TABLE sections (
    id         text PRIMARY KEY,
    type       text    NOT NULL CHECK (type IN ('news', 'walkthrough')),
    title      text    NOT NULL,
    archived   boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO sections (id, type, title) VALUES
    ('pc', 'news', 'PC'),
    ('xbox', 'news', 'Xbox Series X|S'),
    ('playstation', 'news', 'PlayStation 5'),
    ('nintendo', 'news', 'Nintendo'),
    ('vr', 'news', 'VR'),
    ('hardware', 'news', 'Железо'),
    ('control-resonant', 'walkthrough', 'CONTROL RESONANT'),
    ('silent-hill-townfall', 'walkthrough', 'SILENT HILL: TOWNFALL'),
    ('ghost-of-tsushima', 'walkthrough', 'GHOST OF TSUSHIMA'),
    ('baldurs-gate-3', 'walkthrough', 'BALDUR’S GATE 3');

-- Walkthroughs become a first-class article kind.
ALTER TABLE articles DROP CONSTRAINT articles_kind_check;
ALTER TABLE articles ADD CONSTRAINT articles_kind_check CHECK (kind IN ('news', 'review', 'walkthrough'));
