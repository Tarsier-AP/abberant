DELETE FROM users
WHERE lower(email) IN (
    'owner@abberant.test',
    'moderator@abberant.test',
    'creator@abberant.test',
    'user@abberant.test'
);
