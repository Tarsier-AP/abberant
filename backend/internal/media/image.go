package media

import (
	"image"

	"golang.org/x/image/draw"
)

// toRGBA draws any image into a fresh RGBA buffer so pixel access and
// encoding work on a single predictable format.
func toRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	bounds := img.Bounds()

	out := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	draw.Draw(out, out.Bounds(), img, bounds.Min, draw.Src)

	return out
}

// squareCrop center-crops the image to a square the size of its shorter edge.
func squareCrop(img image.Image) *image.RGBA {
	bounds := img.Bounds()

	size := min(bounds.Dx(), bounds.Dy())

	cx := bounds.Min.X + (bounds.Dx()-size)/2
	cy := bounds.Min.Y + (bounds.Dy()-size)/2

	out := image.NewRGBA(image.Rect(0, 0, size, size))

	draw.Draw(out, out.Bounds(), img, image.Point{X: cx, Y: cy}, draw.Src)

	return out
}

// resize scales the image to exactly the requested size using a smooth
// interpolator; content is assumed to already be cropped/fitted by the caller.
func resize(source *image.RGBA, width, height int) *image.RGBA {
	if width == source.Rect.Dx() && height == source.Rect.Dy() {
		return source
	}

	out := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.CatmullRom.Scale(out, out.Bounds(), source, source.Bounds(), draw.Over, nil)

	return out
}
