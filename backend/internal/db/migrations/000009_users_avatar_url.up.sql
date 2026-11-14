-- User avatars move from inline data URLs to media storage; the column now
-- holds either a legacy data URL or a /media/{hash}/... address.
ALTER TABLE users RENAME COLUMN avatar_data_url TO avatar;
