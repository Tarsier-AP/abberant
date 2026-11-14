-- Remove the placeholder article seed that shipped with the initial schema:
-- mock news, reviews and walkthroughs that used the SVG placeholder cover.
-- The real November 2026 news batch (000010) and editor-created articles stay.
-- Comments, bookmarks, notifications and votes cascade with the deleted rows.
DELETE FROM articles
WHERE id ~ '^news-(pc|xbox|playstation|nintendo|vr|hardware)-[1-6]$'
   OR id ~ '^review-[1-3]$'
   OR id ~ '^recent-[1-3]$'
   OR id ~ '^control-resonant-[1-9]$'
   OR id ~ '^silent-hill-townfall-[1-9]$'
   OR id ~ '^ghost-of-tsushima-[1-9]$'
   OR id ~ '^baldurs-gate-3-[1-9]$'
   OR id ~ '^walkthrough-[1-3]$';
