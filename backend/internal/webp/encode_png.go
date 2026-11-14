//go:build !linux

package webp

// Extension is the file-name extension of the images produced by
// EncodeRGBA on this platform.
const Extension = ".png"

// ContentType is the MIME type of the images produced by EncodeRGBA.
const ContentType = "image/png"

// EncodeRGBA encodes a tightly packed RGBA buffer as PNG.
//
// This is the cgo-free fallback for platforms without libwebp (Windows
// builds must stay pure Go): PNG is lossless and therefore heavier than
// WebP, but it renders everywhere and requires no external library.
// The quality parameter is ignored.
func EncodeRGBA(rgba []byte, width, height, stride int, quality float64) ([]byte, error) {
	return encodePNG(rgba, width, height, stride)
}
