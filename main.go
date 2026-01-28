package main

import (
	"fmt"
	"image"
	"math"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	log "github.com/dsoprea/go-logging"
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
const THUMB_QUALITY = 80
const THUMB_SCALE = 8

func main() {
	buffer, err := imgio.Open("examples/image.jpg")
	log.PanicIf(err)
	bounds := buffer.Bounds().Max
	newWidth, newHeight := bounds.X/THUMB_SCALE, bounds.Y/THUMB_SCALE
	minWidth := math.Min(float64(newWidth), float64(newHeight))
	thumbnail := transform.Resize(buffer, newWidth, newHeight, transform.Linear)
	rawExif, err := exif.SearchFileAndExtractExif("examples/image.jpg")
	im, err := exifcommon.NewIfdMappingWithStandard()
	log.PanicIf(err)
	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	log.PanicIf(err)
	tagName := "Orientation"
	rootIfd := index.RootIfd
	results, err := rootIfd.FindTagWithName(tagName)
	if len(results) < 1 {
		panic("Failed to find 'Orientation' in EXIF!")
	}
	ite := results[0]
	valueRaw, err := ite.Value()
	log.PanicIf(err)
	value := valueRaw.([]uint16)[0]
	if value != Horizontal {
		// Must rotate
		if value == Rotate180 {
			thumbnail = transform.Rotate(thumbnail, 180, &transform.RotationOptions{})
		} else {
			// TODO: Handle other orientations
			fmt.Println("WARNING: Unhandled orientation", value)
		}
	} else {
		// Save as-is
	}
	// Crop to square, should try to center the cropped rectangle later
	thumbnail = transform.Crop(thumbnail, image.Rect(0, 0, int(minWidth), int(minWidth)))
	if err := imgio.Save("examples/thumb.jpg", thumbnail, imgio.JPEGEncoder(THUMB_QUALITY)); err != nil {
		panic(err)
	}
}
