-- Test accounts for each role so the API can be exercised right after a fresh
-- database. All four are activated and share the password «Abberant2026!»:
--
--   owner@abberant.test     / owner      — full access (edit + moderate)
--   moderator@abberant.test / moderator  — moderate drafts and comments
--   creator@abberant.test   / creator    — create and edit own articles
--   user@abberant.test      / user       — regular account, no extra rights
INSERT INTO users (email, username, activated, role, password_hash) VALUES
    ('owner@abberant.test',     'owner',     true, 'owner',
     '$argon2id$v=19$m=65536,t=3,p=1$KFdnp8ptj2jMxYRBh8DiIQ$Lr1LH/XGG19p8z4OjFEfKB0sO5lWY6O1IeuoSq3+aas'),
    ('moderator@abberant.test', 'moderator', true, 'moderator',
     '$argon2id$v=19$m=65536,t=3,p=1$KFdnp8ptj2jMxYRBh8DiIQ$Lr1LH/XGG19p8z4OjFEfKB0sO5lWY6O1IeuoSq3+aas'),
    ('creator@abberant.test',   'creator',   true, 'creator',
     '$argon2id$v=19$m=65536,t=3,p=1$KFdnp8ptj2jMxYRBh8DiIQ$Lr1LH/XGG19p8z4OjFEfKB0sO5lWY6O1IeuoSq3+aas'),
    ('user@abberant.test',      'user',      true, 'user',
     '$argon2id$v=19$m=65536,t=3,p=1$KFdnp8ptj2jMxYRBh8DiIQ$Lr1LH/XGG19p8z4OjFEfKB0sO5lWY6O1IeuoSq3+aas')
ON CONFLICT (lower(email)) DO NOTHING;
