package main

import (
	"image"
	"image/color"
	"testing"
)

// testImage returns a 3x2 image whose pixel at (x,y) is encoded as color R=x, G=y.
func testImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), A: 255})
		}
	}
	return img
}

func TestApplyOrientation(t *testing.T) {
	tests := []struct {
		name         string
		orientation  uint16
		w, h         int
		x, y         int // source pixel to probe
		wantX, wantY int // where that source pixel should land
	}{
		{"horizontal", Horizontal, 3, 2, 0, 0, 0, 0},
		{"mirror horizontal", MirrorHorizontal, 3, 2, 0, 0, 2, 0},
		{"rotate 180", Rotate180, 3, 2, 0, 0, 2, 1},
		{"mirror vertical", MirrorVertical, 3, 2, 0, 0, 0, 1},
		{"rotate 90 cw", Rotate90CW, 2, 3, 0, 0, 1, 0},
		{"rotate 270 cw", Rotate270CW, 2, 3, 0, 0, 0, 2},
		{"mirror h + rotate 90 cw", MirrorHorizontalRotate90CW, 2, 3, 0, 0, 1, 2},
		{"mirror h + rotate 270 cw", MirrorHorizontalRotate270CW, 2, 3, 0, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := applyOrientation(testImage(), tt.orientation)
			b := out.Bounds()
			if b.Dx() != tt.w || b.Dy() != tt.h {
				t.Fatalf("size = %dx%d, want %dx%d", b.Dx(), b.Dy(), tt.w, tt.h)
			}
			src := testImage().At(tt.x, tt.y).(color.RGBA)
			got := out.At(tt.wantX, tt.wantY).(color.RGBA)
			if got.R != src.R || got.G != src.G {
				t.Fatalf("pixel at (%d,%d) = %v, want source (%d,%d) = %v",
					tt.wantX, tt.wantY, got, tt.x, tt.y, src)
			}
		})
	}
}

func TestFitWithin(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4000, 3000))
	b := fitWithin(img, 2000).Bounds()
	if b.Dx() != 2000 || b.Dy() != 1500 {
		t.Fatalf("size = %dx%d, want 2000x1500", b.Dx(), b.Dy())
	}
	small := image.NewRGBA(image.Rect(0, 0, 100, 50))
	if got := fitWithin(small, 2000).Bounds(); got.Dx() != 100 || got.Dy() != 50 {
		t.Fatalf("small image was resized to %v", got)
	}
}

func TestSquareThumbnail(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4000, 3000))
	b := squareThumbnail(img, 400).Bounds()
	if b.Dx() != 400 || b.Dy() != 400 {
		t.Fatalf("size = %dx%d, want 400x400", b.Dx(), b.Dy())
	}
}
