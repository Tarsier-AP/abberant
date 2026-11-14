package media

import (
	"os"
	"path/filepath"
	"time"
)

// gcGrace keeps directories younger than this even when unreferenced, so a
// hash written moments ago cannot be removed while its upload request is
// still finishing.
const gcGrace = time.Hour

// GarbageCollect removes hash directories that no article cover, article
// content or user avatar references anymore. Returns the number of removed
// directories.
func (s *Store) GarbageCollect(referenced map[string]struct{}) (int, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}

		return 0, err
	}

	cutoff := time.Now().Add(-gcGrace)

	removed := 0

	for _, entry := range entries {
		if !entry.IsDir() || !hashPattern.MatchString(entry.Name()) {
			continue
		}

		if _, ok := referenced[entry.Name()]; ok {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Directory mod time moves when files are added, which keeps fresh
		// uploads safe even with a coarse ReadDir stat.
		if info.ModTime().After(cutoff) {
			continue
		}

		if err := os.RemoveAll(filepath.Join(s.dir, entry.Name())); err != nil {
			return removed, err
		}

		removed++
	}

	return removed, nil
}
