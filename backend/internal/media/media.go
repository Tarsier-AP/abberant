// Package media stores uploaded images on disk and serves their derivatives
// under content-addressed /media/{hash}/{file} URLs.
//
// Layout: every upload is hashed (SHA-256); all derivatives of that upload
// live in a single per-hash directory. Because identical bytes always map to
// the same directory, /media responses can be cached immutably forever.
// The original bytes are kept verbatim; generated variants use lossy WebP
// on Linux (libwebp via cgo) and PNG elsewhere (Windows builds stay cgo-free).
package media

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg" // register JPEG decoder
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "golang.org/x/image/webp" // register decoder

	"abberant/backend/internal/webp"
)

const (
	// Quality is the lossy quality (0–100) used for every generated WebP
	// variant; ignored by the PNG fallback encoder.
	Quality = 80.0

	// AvatarSize is the edge length of the square avatar derivative.
	AvatarSize = 160

	// MaxUpload is the largest accepted upload in bytes.
	MaxUpload = 10 << 20

	// maxDimension caps the pixel source used to generate variants so a
	// huge upload cannot monopolize decode memory. The original itself is
	// always stored verbatim.
	maxDimension = 2560
)

// CoverWidths are the generated variant widths for covers and body images.
// The set is a fixed contract: the frontend derives srcset attributes from
// it, so every raster upload gets all of these variants. The file format of
// a variant depends on the platform (webp.Extension).
var CoverWidths = []int{320, 480, 828, 1200}

// Kind selects how an upload is processed.
type Kind string

const (
	KindCover  Kind = "cover"
	KindBody   Kind = "body"
	KindAvatar Kind = "avatar"
)

// File describes one derivative written for an upload.
type File struct {
	// URL is the site-relative /media/... address of the file.
	URL string
	// Width and Height are the pixel dimensions of the derivative.
	Width, Height int
	// Original marks the untouched (only dimension-capped) upload.
	Original bool
}

// Result reports everything stored for one upload.
type Result struct {
	Hash  string
	Files []File
}

// MainURL returns the default image URL: the avatar derivative for avatars,
// otherwise the largest generated WebP variant.
func (r *Result) MainURL() string {
	for i := len(r.Files) - 1; i >= 0; i-- {
		if !r.Files[i].Original {
			return r.Files[i].URL
		}
	}

	if len(r.Files) > 0 {
		return r.Files[0].URL
	}

	return ""
}

// Store keeps uploaded media under a configurable directory.
type Store struct {
	dir string
}

func NewStore(dir string) *Store { return &Store{dir: dir} }

func (s *Store) Dir() string { return s.dir }

var hashPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// ValidMediaPath reports whether the path after /media/ addresses a file of
// the store: a hash directory plus a plain file name. Used by the static
// handler and by validators that accept site-relative media URLs.
func ValidMediaPath(path string) bool {
	hash, name, ok := strings.Cut(path, "/")
	if !ok || !hashPattern.MatchString(hash) {
		return false
	}

	if name == "" || strings.Contains(name, "/") || strings.Contains(name, `\`) {
		return false
	}

	return name == filepath.Base(name)
}

// Put decodes an upload, generates every derivative and writes them to disk.
func (s *Store) Put(data []byte) (*Result, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("media: пустой файл")
	}

	if len(data) > MaxUpload {
		return nil, fmt.Errorf("media: файл больше %d МБ", MaxUpload>>20)
	}

	decoded, format, err := decode(data)
	if err != nil {
		return nil, err
	}

	// Variants are generated from a dimension-capped copy; the original
	// itself is written verbatim.
	variants := capDimension(decoded)

	hash := hashBytes(data)
	dir := filepath.Join(s.dir, hash)

	original, err := s.writeOriginal(dir, hash, data, decoded, format)
	if err != nil {
		return nil, err
	}

	result := &Result{Hash: hash, Files: []File{original}}

	for _, width := range CoverWidths {
		file, err := s.writeVariant(dir, hash, variants, width)
		if err != nil {
			return nil, err
		}

		result.Files = append(result.Files, file)
	}

	return result, nil
}

// PutAvatar decodes an upload, center-crops it to a square, scales it to
// AvatarSize and writes the avatar derivative.
func (s *Store) PutAvatar(data []byte) (*Result, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("media: пустой файл")
	}

	if len(data) > MaxUpload {
		return nil, fmt.Errorf("media: файл больше %d МБ", MaxUpload>>20)
	}

	decoded, format, err := decode(data)
	if err != nil {
		return nil, err
	}

	hash := hashBytes(data)
	dir := filepath.Join(s.dir, hash)

	original, err := s.writeOriginal(dir, hash, data, decoded, format)
	if err != nil {
		return nil, err
	}

	rgba := toRGBA(squareCrop(capDimension(decoded)))
	scaled := resize(rgba, AvatarSize, AvatarSize)

	encoded, err := webp.EncodeRGBA(scaled.Pix, AvatarSize, AvatarSize, scaled.Stride, Quality)
	if err != nil {
		return nil, fmt.Errorf("media: webp: %w", err)
	}

	name := "avatar" + webp.Extension

	if err := writeFileSync(dir, name, encoded); err != nil {
		return nil, err
	}

	return &Result{
		Hash: hash,
		Files: []File{
			original,
			{URL: mediaURL(hash, name), Width: AvatarSize, Height: AvatarSize},
		},
	}, nil
}

func (s *Store) writeOriginal(dir, hash string, data []byte, decoded image.Image, format string) (File, error) {
	var extension string

	switch format {
	case "png":
		extension = ".png"
	case "jpeg":
		extension = ".jpg"
	case "webp":
		extension = ".webp"
	default:
		return File{}, fmt.Errorf("media: неизвестный формат %q", format)
	}

	name := "original" + extension

	if err := writeFileSync(dir, name, data); err != nil {
		return File{}, err
	}

	bounds := decoded.Bounds()

	return File{
		URL:      mediaURL(hash, name),
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
		Original: true,
	}, nil
}

func (s *Store) writeVariant(dir, hash string, img image.Image, width int) (File, error) {
	bounds := img.Bounds()
	sourceWidth, sourceHeight := bounds.Dx(), bounds.Dy()

	// Fit the source to the target width without distorting it.
	scale := float64(width) / float64(sourceWidth)
	targetHeight := int(float64(sourceHeight)*scale + 0.5)

	rgba := toRGBA(img)
	scaled := resize(rgba, width, targetHeight)

	encoded, err := webp.EncodeRGBA(scaled.Pix, scaled.Rect.Dx(), scaled.Rect.Dy(), scaled.Stride, Quality)
	if err != nil {
		return File{}, fmt.Errorf("media: variant %dw: %w", width, err)
	}

	name := fmt.Sprintf("%d%s", width, webp.Extension)

	if err := writeFileSync(dir, name, encoded); err != nil {
		return File{}, err
	}

	return File{
		URL:    mediaURL(hash, name),
		Width:  width,
		Height: targetHeight,
	}, nil
}

func writeFileSync(dir, name string, data []byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("media: mkdir: %w", err)
	}

	// Temp file plus rename so concurrent writes of the same hash never
	// expose a half-written derivative.
	temp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("media: temp: %w", err)
	}

	tempName := temp.Name()

	if _, err := temp.Write(data); err != nil {
		temp.Close()
		os.Remove(tempName)
		return fmt.Errorf("media: write: %w", err)
	}

	if err := temp.Close(); err != nil {
		os.Remove(tempName)
		return fmt.Errorf("media: close: %w", err)
	}

	// Temp files are private; derivatives must be readable by the static
	// file server (nginx) in production.
	if err := os.Chmod(tempName, 0o644); err != nil {
		os.Remove(tempName)
		return fmt.Errorf("media: chmod: %w", err)
	}

	if err := os.Rename(tempName, filepath.Join(dir, name)); err != nil {
		os.Remove(tempName)
		return fmt.Errorf("media: rename: %w", err)
	}

	return nil
}

func decode(data []byte) (image.Image, string, error) {
	decoded, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("media: не удалось прочитать изображение: %w", err)
	}

	switch format {
	case "png", "jpeg", "webp":
		return decoded, format, nil
	default:
		return nil, "", fmt.Errorf("media: формат %q не поддерживается (JPG, PNG или WebP)", format)
	}
}

func capDimension(img image.Image) image.Image {
	bounds := img.Bounds()

	longest := max(bounds.Dx(), bounds.Dy())
	if longest <= maxDimension {
		return img
	}

	scale := float64(maxDimension) / float64(longest)

	return resize(toRGBA(img), int(float64(bounds.Dx())*scale+0.5), int(float64(bounds.Dy())*scale+0.5))
}

func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}

func mediaURL(hash, name string) string {
	return "/media/" + hash + "/" + name
}
