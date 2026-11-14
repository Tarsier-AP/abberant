// Package webp encodes RGBA pixel buffers for the media store.
//
// On Linux the encoder is libwebp (cgo): uploads become lossy WebP at the
// configured quality. See encode_cgo.go.
//
// On every other platform (Windows builds stay cgo-free and are
// cross-compilable without a C toolchain) EncodeRGBA produces PNG instead;
// files simply get a different extension. See encode_png.go.
package webp
