package webp

import (
	"bytes"
	"image"
	"testing"
)

func TestEncodePNGFallback(t *testing.T) {
	width, height, stride := 5, 3, 24 // stride has extra padding per row

	rgba := make([]byte, stride*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := y*stride + x*4

			rgba[offset] = uint8(x * 40)
			rgba[offset+1] = uint8(y * 80)
			rgba[offset+2] = uint8((x + y) * 20)
			rgba[offset+3] = 0xff

			// Padding bytes between rows would corrupt the image if the
			// encoder ignored the stride.
			if x == width-1 {
				copy(rgba[offset+4:(y+1)*stride], []byte{9, 9, 9, 9, 9})
			}
		}
	}

	data, err := encodePNG(rgba, width, height, stride)
	if err != nil {
		t.Fatalf("encodePNG: %v", err)
	}

	decoded, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if format != "png" {
		t.Fatalf("format %q, want png", format)
	}

	bounds := decoded.Bounds()

	if bounds.Dx() != width || bounds.Dy() != height {
		t.Fatalf("got %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), width, height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := decoded.At(x, y).RGBA()

			want := rgba[y*stride+x*4 : y*stride+x*4+4]

			if uint8(r>>8) != want[0] || uint8(g>>8) != want[1] || uint8(b>>8) != want[2] || uint8(a>>8) != want[3] {
				t.Fatalf("pixel (%d,%d) mismatch", x, y)
			}
		}
	}
}

func TestEncodePNGFallbackRejectsBadInput(t *testing.T) {
	cases := []struct {
		name          string
		rgba          []byte
		width, height int
		stride        int
	}{
		{"zero dims", make([]byte, 16), 0, 1, 4},
		{"negative dims", make([]byte, 16), -1, 1, 4},
		{"small stride", make([]byte, 16), 4, 1, 4},
		{"short buffer", make([]byte, 8), 1, 4, 4},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := encodePNG(testCase.rgba, testCase.width, testCase.height, testCase.stride); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
