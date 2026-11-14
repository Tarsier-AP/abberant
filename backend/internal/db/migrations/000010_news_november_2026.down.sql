-- Revert the November 2026 news batch.
DELETE FROM articles
WHERE id IN (
    'news-pc-7',
    'news-pc-8',
    'news-xbox-7',
    'news-xbox-8',
    'news-playstation-7',
    'news-playstation-8',
    'news-nintendo-7',
    'news-nintendo-8',
    'news-vr-7',
    'news-hardware-7'
);
