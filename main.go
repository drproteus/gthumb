package main

import (
	"fmt"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	log "github.com/dsoprea/go-logging"
)

func main() {
	buffer, err := imgio.Open("examples/image.jpg")
	log.PanicIf(err)
	bounds := buffer.Bounds().Max
	thumbnail := transform.Resize(buffer, bounds.X/4, bounds.Y/4, transform.Linear)
	if err := imgio.Save("examples/thumb.jpg", thumbnail, imgio.JPEGEncoder(100)); err != nil {
		panic(err)
	}
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
	// EXIF Orientations
	// 1 = Horizontal (normal)
	// 2 = Mirror horizontal
	// 3 = Rotate 180
	// 4 = Mirror vertical
	// 5 = Mirror horizontal and rotate 270 CW
	// 6 = Rotate 90 CW
	// 7 = Mirror horizontal and rotate 90 CW
	// 8 = Rotate 270 CW
	value := valueRaw.([]uint16)[0]
	switch value {
	case 1:
		fmt.Println(value, "Normal")
	case 2:
		fmt.Println(value, "Mirror horizontal")
	case 3:
		fmt.Println(value, "Rotate 180")
	case 4:
		fmt.Println(value, "Mirror vertical")
	default:
		fmt.Println(value, "TODO")
	}
}
