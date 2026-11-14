package webp

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

// encodePNG turns a tightly packed RGBA buffer into a PNG file. It holds
// the whole algorithm of the cgo-free fallback so it stays testable on any
// platform, including Linux where the fallback itself is never compiled.
func encodePNG(rgba []byte, width, height, stride int) ([]byte, error) {
	if err := validateArgs(rgba, width, height, stride); err != nil {
		return nil, err
	}

	// Copy row by row: the source buffer may have padding the PNG encoder
	// must not pick up.
	out := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		copy(out.Pix[y*out.Stride:y*out.Stride+width*4], rgba[y*stride:y*stride+width*4])
	}

	var buffer bytes.Buffer

	if err := png.Encode(&buffer, out); err != nil {
		return nil, fmt.Errorf("webp: png fallback: %w", err)
	}

	return buffer.Bytes(), nil
}

// validateArgs checks the invariants shared by every encoder.
func validateArgs(rgba []byte, width, height, stride int) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("webp: invalid dimensions %dx%d", width, height)
	}

	if stride < width*4 {
		return fmt.Errorf("webp: stride %d is smaller than width*4", stride)
	}

	if len(rgba) < stride*height {
		return fmt.Errorf("webp: buffer of %d bytes is too small for %dx%d", len(rgba), width, height)
	}

	return nil
}
