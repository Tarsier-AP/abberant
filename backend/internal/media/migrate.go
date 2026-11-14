package media

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// mimePattern matches the raster data URLs the migration converts. SVG data
// URLs (used by seeded example accounts) cannot be raster-decoded and are
// left untouched.
var mimePattern = regexp.MustCompile(`^data:image/(png|jpeg|jpg|webp);base64,(.+)$`)

// contentImagePattern matches markdown image syntax with a raster data URL,
// as inserted by the editor: ![alt](data:image/png;base64,... "caption").
var contentImagePattern = regexp.MustCompile(
	`!\[[^\]]*\]\(\s*data:image/(?:png|jpeg|jpg|webp);base64,[A-Za-z0-9+/=]+(?:\s+"[^"]*")?\s*\)`,
)

// mediaHashPattern extracts the hash directory of site-relative media URLs.
var mediaHashPattern = regexp.MustCompile(`/media/([0-9a-f]{64})/`)

// MigrateDataURLs converts legacy inline data URLs to media storage:
//
//   - article cover images stored as data URLs,
//   - raster images inside article content (markdown image syntax),
//   - user avatars (raster ones only; SVG data URLs are kept as-is).
//
// The migration is idempotent: rows already carrying /media/ URLs are never
// touched, so running it on every start is safe.
func MigrateDataURLs(ctx context.Context, pool *pgxpool.Pool, store *Store) error {
	if err := migrateArticleCovers(ctx, pool, store); err != nil {
		return fmt.Errorf("media: migrate covers: %w", err)
	}

	if err := migrateArticleContent(ctx, pool, store); err != nil {
		return fmt.Errorf("media: migrate content: %w", err)
	}

	if err := migrateAvatars(ctx, pool, store); err != nil {
		return fmt.Errorf("media: migrate avatars: %w", err)
	}

	return nil
}

func migrateArticleCovers(ctx context.Context, pool *pgxpool.Pool, store *Store) error {
	rows, err := pool.Query(ctx, `
		SELECT id, image FROM articles
		WHERE image LIKE 'data:image/%'
	`)
	if err != nil {
		return err
	}

	defer rows.Close()

	type update struct {
		id  string
		url string
	}

	updates := make([]update, 0, 16)

	for rows.Next() {
		var (
			id    string
			image string
		)

		if err := rows.Scan(&id, &image); err != nil {
			return err
		}

		result, err := store.Put(DecodeDataURL(image))
		if err != nil {
			// Undecodable covers stay as they are instead of blocking
			// startup; they keep working through the data-URL validators.
			log.Printf("media: cover of %s stays a data URL: %v", id, err)
			continue
		}

		updates = append(updates, update{id: id, url: result.MainURL()})
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, item := range updates {
		if _, err := pool.Exec(ctx, `UPDATE articles SET image = $2 WHERE id = $1`, item.id, item.url); err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		log.Printf("media: migrated %d article cover(s)", len(updates))
	}

	return nil
}

func migrateArticleContent(ctx context.Context, pool *pgxpool.Pool, store *Store) error {
	rows, err := pool.Query(ctx, `
		SELECT id, content FROM articles WHERE content LIKE '%data:image/%'
	`)
	if err != nil {
		return err
	}

	defer rows.Close()

	type update struct {
		id      string
		content string
	}

	updates := make([]update, 0, 16)

	for rows.Next() {
		var (
			id      string
			content string
		)

		if err := rows.Scan(&id, &content); err != nil {
			return err
		}

		replacements := make(map[string]string)

		for _, snippet := range contentImagePattern.FindAllString(content, -1) {
			dataURL := contentImageDataURL(snippet)
			if dataURL == "" {
				continue
			}

			if _, ok := replacements[dataURL]; ok {
				continue
			}

			result, err := store.Put(DecodeDataURL(dataURL))
			if err != nil {
				log.Printf("media: content image of %s stays a data URL: %v", id, err)
				continue
			}

			// Keep the largest generated variant as the inline image; it is
			// also what the article zoom view shows.
			replacements[dataURL] = result.MainURL()
		}

		if len(replacements) == 0 {
			continue
		}

		for dataURL, url := range replacements {
			content = strings.ReplaceAll(content, dataURL, url)
		}

		updates = append(updates, update{id: id, content: content})
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, item := range updates {
		if _, err := pool.Exec(ctx, `UPDATE articles SET content = $2 WHERE id = $1`, item.id, item.content); err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		log.Printf("media: migrated content of %d article(s)", len(updates))
	}

	return nil
}

func migrateAvatars(ctx context.Context, pool *pgxpool.Pool, store *Store) error {
	rows, err := pool.Query(ctx, `
		SELECT id, avatar FROM users
		WHERE avatar LIKE 'data:image/%'
		  AND avatar NOT LIKE 'data:image/svg%'
	`)
	if err != nil {
		return err
	}

	defer rows.Close()

	type update struct {
		id  int64
		url string
	}

	updates := make([]update, 0, 16)

	for rows.Next() {
		var (
			id     int64
			avatar string
		)

		if err := rows.Scan(&id, &avatar); err != nil {
			return err
		}

		result, err := store.PutAvatar(DecodeDataURL(avatar))
		if err != nil {
			log.Printf("media: avatar of user %d stays a data URL: %v", id, err)
			continue
		}

		updates = append(updates, update{id: id, url: result.MainURL()})
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, item := range updates {
		if _, err := pool.Exec(ctx, `UPDATE users SET avatar = $2 WHERE id = $1`, item.id, item.url); err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		log.Printf("media: migrated %d avatar(s)", len(updates))
	}

	return nil
}

// DecodeDataURL extracts the raw bytes of a raster data URL. Malformed input
// yields nil, which the store rejects with a clear error.
func DecodeDataURL(dataURL string) []byte {
	match := mimePattern.FindStringSubmatch(strings.TrimSpace(dataURL))
	if match == nil {
		return nil
	}

	raw, err := base64.StdEncoding.DecodeString(match[2])
	if err != nil {
		return nil
	}

	return raw
}

// contentImageDataURL pulls the data URL out of a matched markdown snippet.
func contentImageDataURL(snippet string) string {
	match := regexp.MustCompile(`data:image/(?:png|jpeg|jpg|webp);base64,[A-Za-z0-9+/=]+`).
		FindString(snippet)

	return match
}

// ReferencedHashes walks every column that stores media URLs and returns the
// set of hash directories that must survive a garbage-collection sweep.
func ReferencedHashes(ctx context.Context, pool *pgxpool.Pool) (map[string]struct{}, error) {
	referenced := make(map[string]struct{})

	collect := func(query string, args ...any) error {
		rows, err := pool.Query(ctx, query, args...)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var value string

			if err := rows.Scan(&value); err != nil {
				return err
			}

			for _, hash := range mediaHashPattern.FindAllStringSubmatch(value, -1) {
				referenced[hash[1]] = struct{}{}
			}
		}

		return rows.Err()
	}

	if err := collect(`SELECT image FROM articles`); err != nil {
		return nil, err
	}

	if err := collect(`SELECT content FROM articles WHERE content LIKE '%/media/%'`); err != nil {
		return nil, err
	}

	if err := collect(`SELECT avatar FROM users WHERE avatar LIKE '/media/%'`); err != nil {
		return nil, err
	}

	return referenced, nil
}

// RunMaintenance performs the one-time data-URL migration and then garbage
// collects unreferenced uploads once a day until ctx is done.
func RunMaintenance(ctx context.Context, pool *pgxpool.Pool, store *Store) {
	if err := MigrateDataURLs(ctx, pool, store); err != nil {
		log.Printf("media: data-url migration failed: %v", err)
	}

	sweep := func() {
		referenced, err := ReferencedHashes(ctx, pool)
		if err != nil {
			log.Printf("media: gc: %v", err)
			return
		}

		removed, err := store.GarbageCollect(referenced)
		if err != nil {
			log.Printf("media: gc: %v", err)
			return
		}

		if removed > 0 {
			log.Printf("media: gc: removed %d unreferenced upload(s)", removed)
		}
	}

	sweep()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sweep()
		}
	}
}
