package media

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "golang.org/x/image/webp" // register decoder
)

func gradientPNG(t *testing.T, width, height int) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := img.PixOffset(x, y)

			img.Pix[offset] = uint8(x % 256)
			img.Pix[offset+1] = uint8(y % 256)
			img.Pix[offset+2] = uint8((x + y) % 256)
			img.Pix[offset+3] = 0xff
		}
	}

	var buffer strings.Builder

	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("encode test png: %v", err)
	}

	return []byte(buffer.String())
}

func TestPutGeneratesVariants(t *testing.T) {
	store := NewStore(t.TempDir())

	result, err := store.Put(gradientPNG(t, 640, 400))
	if err != nil {
		t.Fatalf("put: %v", err)
	}

	if len(result.Hash) != 64 || !hashPattern.MatchString(result.Hash) {
		t.Fatalf("unexpected hash %q", result.Hash)
	}

	// original + every fixed cover width, even though the source is smaller
	// than the two largest widths: the width set is a frontend contract.
	wantNames := []string{"original.png", "320.webp", "480.webp", "828.webp", "1200.webp"}

	if len(result.Files) != len(wantNames) {
		t.Fatalf("got %d files, want %d: %+v", len(result.Files), len(wantNames), result.Files)
	}

	dir := filepath.Join(store.Dir(), result.Hash)

	for index, name := range wantNames {
		file := result.Files[index]

		if !strings.HasSuffix(file.URL, "/"+name) {
			t.Fatalf("file %d url %q, want suffix %q", index, file.URL, name)
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		if name == "original.png" {
			decoded, format, err := image.Decode(strings.NewReader(string(data)))
			if err != nil || format != "png" {
				t.Fatalf("original did not survive as png: %v %q", err, format)
			}

			bounds := decoded.Bounds()

			if bounds.Dx() != 640 || bounds.Dy() != 400 {
				t.Fatalf("original resized: %dx%d", bounds.Dx(), bounds.Dy())
			}

			continue
		}

		decoded, format, err := image.Decode(strings.NewReader(string(data)))
		if err != nil {
			t.Fatalf("decode %s: %v", name, err)
		}

		if format != "webp" {
			t.Fatalf("%s is %q, want webp", name, format)
		}

		bounds := decoded.Bounds()

		if bounds.Dx() != file.Width || bounds.Dy() != file.Height {
			t.Fatalf("%s is %dx%d, want %dx%d", name, bounds.Dx(), bounds.Dy(), file.Width, file.Height)
		}
	}

	if result.MainURL() != result.Files[4].URL {
		t.Fatalf("MainURL %q, want the largest variant", result.MainURL())
	}
}

func TestPutAvatar(t *testing.T) {
	store := NewStore(t.TempDir())

	result, err := store.PutAvatar(gradientPNG(t, 400, 300))
	if err != nil {
		t.Fatalf("put avatar: %v", err)
	}

	if len(result.Files) != 2 {
		t.Fatalf("got %d files, want original + avatar.webp: %+v", len(result.Files), result.Files)
	}

	if result.MainURL() != "/media/"+result.Hash+"/avatar.webp" {
		t.Fatalf("MainURL %q", result.MainURL())
	}

	data, err := os.ReadFile(filepath.Join(store.Dir(), result.Hash, "avatar.webp"))
	if err != nil {
		t.Fatalf("read avatar: %v", err)
	}

	decoded, format, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("decode avatar: %v", err)
	}

	if format != "webp" {
		t.Fatalf("avatar is %q, want webp", format)
	}

	bounds := decoded.Bounds()

	if bounds.Dx() != AvatarSize || bounds.Dy() != AvatarSize {
		t.Fatalf("avatar is %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), AvatarSize, AvatarSize)
	}
}

func TestGarbageCollect(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	upload := gradientPNG(t, 400, 300)

	referenced, err := store.Put(upload)
	if err != nil {
		t.Fatalf("put: %v", err)
	}

	orphan, err := store.Put(gradientPNG(t, 401, 301))
	if err != nil {
		t.Fatalf("put: %v", err)
	}

	// Pretend the orphan is old enough to be collected.
	if err := os.Chtimes(filepath.Join(dir, orphan.Hash), time.Unix(0, 0), time.Unix(0, 0)); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	removed, err := store.GarbageCollect(map[string]struct{}{referenced.Hash: {}})
	if err != nil {
		t.Fatalf("gc: %v", err)
	}

	if removed != 1 {
		t.Fatalf("removed %d directories, want 1", removed)
	}

	if _, err := os.Stat(filepath.Join(dir, referenced.Hash)); err != nil {
		t.Fatalf("referenced upload was removed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, orphan.Hash)); !os.IsNotExist(err) {
		t.Fatalf("orphan survived: %v", err)
	}
}

func TestValidMediaPath(t *testing.T) {
	hash := strings.Repeat("a", 64)

	valid := []string{
		hash + "/1200.webp",
		hash + "/original.jpg",
		hash + "/avatar.webp",
	}

	for _, path := range valid {
		if !ValidMediaPath(path) {
			t.Errorf("%q should be valid", path)
		}
	}

	invalid := []string{
		"",
		hash + "",
		hash + "/",
		strings.Repeat("g", 64) + "/1200.webp",
		hash + "/../etc/passwd",
		"../../etc/passwd",
		hash + "/nested/1200.webp",
		hash + "/1200.webp/",
	}

	for _, path := range invalid {
		if ValidMediaPath(path) {
			t.Errorf("%q should be invalid", path)
		}
	}
}
