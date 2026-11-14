//go:build linux

package webp

/*
#cgo LDFLAGS: -l:libwebp.so.7

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

// Extension is the file-name extension of the images produced by
// EncodeRGBA on this platform.
const Extension = ".webp"

// ContentType is the MIME type of the images produced by EncodeRGBA.
const ContentType = "image/webp"

// EncodeRGBA encodes a tightly packed RGBA buffer as a lossy WebP image.
// Quality follows the usual 0–100 libwebp scale. The returned slice holds a
// complete .webp file.
//
// Build requirement: the libwebp runtime library (soname libwebp.so.7,
// shipped by the libwebp7 / libwebp packages). No development headers are
// needed — the small part of the API used here is declared above.
func EncodeRGBA(rgba []byte, width, height, stride int, quality float64) ([]byte, error) {
	if err := validateArgs(rgba, width, height, stride); err != nil {
		return nil, err
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
