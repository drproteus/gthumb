package main

import (
	"image"
	"math"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// EXIF Orientations
const (
	Horizontal                  = 1
	MirrorHorizontal            = 2
	Rotate180                   = 3
	MirrorVertical              = 4
	MirrorHorizontalRotate270CW = 5
	Rotate90CW                  = 6
	MirrorHorizontalRotate90CW  = 7
	Rotate270CW                 = 8
)

func openImage(path string) (image.Image, error) {
	return imgio.Open(path)
}

func saveJPEG(path string, img image.Image, quality int) error {
	return imgio.Save(path, img, imgio.JPEGEncoder(quality))
}

// applyOrientation returns img transformed so that it is displayed upright,
// following the EXIF Orientation tag semantics.
func applyOrientation(img image.Image, orientation uint16) image.Image {
	switch orientation {
	case MirrorHorizontal:
		return transform.FlipH(img)
	case Rotate180:
		return rotate180(img)
	case MirrorVertical:
		return transform.FlipV(img)
	case MirrorHorizontalRotate270CW:
		return rotate270CW(transform.FlipH(img))
	case Rotate90CW:
		return rotate90CW(img)
	case MirrorHorizontalRotate90CW:
		return rotate90CW(transform.FlipH(img))
	case Rotate270CW:
		return rotate270CW(img)
	default:
		return img
	}
}

func rotate90CW(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(h-1-y, x, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func rotate270CW(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(y, w-1-x, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func rotate180(src image.Image) image.Image {
	return transform.FlipV(transform.FlipH(src))
}

// fitWithin scales img down so neither edge exceeds maxEdge, preserving aspect ratio.
// Images already small enough are returned unchanged.
func fitWithin(img image.Image, maxEdge int) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= maxEdge && h <= maxEdge {
		return img
	}
	scale := float64(maxEdge) / math.Max(float64(w), float64(h))
	nw := max(1, int(math.Round(float64(w)*scale)))
	nh := max(1, int(math.Round(float64(h)*scale)))
	return transform.Resize(img, nw, nh, transform.Linear)
}

// squareThumbnail scales img so its shorter edge equals size, then center-crops to size x size.
func squareThumbnail(img image.Image, size int) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	scale := float64(size) / math.Min(float64(w), float64(h))
	nw := max(size, int(math.Round(float64(w)*scale)))
	nh := max(size, int(math.Round(float64(h)*scale)))
	scaled := transform.Resize(img, nw, nh, transform.Linear)

	x0 := (nw - size) / 2
	y0 := (nh - size) / 2
	return transform.Crop(scaled, image.Rect(x0, y0, x0+size, y0+size))
}
