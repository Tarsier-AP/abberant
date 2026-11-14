// Package webp encodes RGBA pixel buffers to WebP via libwebp (cgo).
//
// Build requirements:
//   - Linux: the libwebp runtime library (soname libwebp.so.7, shipped by the
//     libwebp7 / libwebp packages). No development headers are needed: the
//     small part of the libwebp API used here is declared below.
//   - Windows: a C toolchain (mingw-w64) plus the libwebp import library
//     (webp.dll + libwebp.a) available to the linker.
package webp

/*
#cgo linux LDFLAGS: -l:libwebp.so.7
#cgo windows LDFLAGS: -lwebp

#include <stdlib.h>

size_t WebPEncodeRGBA(const unsigned char* rgba, int width, int height,
                      int stride, float quality, unsigned char** output);
void WebPFree(void* ptr);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// EncodeRGBA encodes a tightly packed RGBA buffer as a lossy WebP image.
// Quality follows the usual 0–100 libwebp scale. The returned slice holds a
// complete .webp file.
func EncodeRGBA(rgba []byte, width, height, stride int, quality float64) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("webp: invalid dimensions %dx%d", width, height)
	}

	if stride < width*4 {
		return nil, fmt.Errorf("webp: stride %d is smaller than width*4", stride)
	}

	if len(rgba) < stride*height {
		return nil, fmt.Errorf("webp: buffer of %d bytes is too small for %dx%d", len(rgba), width, height)
	}

	if quality < 0 {
		quality = 0
	}

	if quality > 100 {
		quality = 100
	}

	var output *C.uchar

	written := C.WebPEncodeRGBA(
		(*C.uchar)(unsafe.Pointer(&rgba[0])),
		C.int(width),
		C.int(height),
		C.int(stride),
		C.float(quality),
		&output,
	)

	if written == 0 {
		return nil, fmt.Errorf("webp: encoding %dx%d failed", width, height)
	}

	defer C.WebPFree(unsafe.Pointer(output))

	return C.GoBytes(unsafe.Pointer(output), C.int(written)), nil
}
