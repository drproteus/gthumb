package main

import (
	"time"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

// Metadata is the subset of EXIF data used by the gallery.
// See https://exiftool.org/TagNames/EXIF.html for tag definitions.
type Metadata struct {
	Orientation uint16
	Camera      string
	Taken       string
}

// readMetadata extracts EXIF data from a JPEG file. Files without EXIF data
// yield the default metadata (orientation 1, no camera, no date) rather than an error.
func readMetadata(path string) Metadata {
	meta := Metadata{Orientation: Horizontal}

	rawExif, err := exif.SearchFileAndExtractExif(path)
	if err != nil {
		return meta
	}
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return meta
	}
	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return meta
	}

	if ite := findTag(index.RootIfd, "Orientation"); ite != nil {
		if value, err := ite.Value(); err == nil {
			if values, ok := value.([]uint16); ok && len(values) > 0 {
				meta.Orientation = values[0]
			}
		}
	}
	if ite := findTag(index.RootIfd, "Model"); ite != nil {
		meta.Camera, _ = ite.FormatFirst()
	}
	if exifIfd, ok := index.Lookup["IFD/Exif"]; ok {
		if ite := findTag(exifIfd, "DateTimeOriginal"); ite != nil {
			if raw, err := ite.FormatFirst(); err == nil {
				if t, err := time.Parse("2006:01:02 15:04:05", raw); err == nil {
					meta.Taken = t.Format("2 Jan 2006")
				}
			}
		}
	}
	return meta
}

func findTag(ifd *exif.Ifd, name string) *exif.IfdTagEntry {
	results, err := ifd.FindTagWithName(name)
	if err != nil || len(results) == 0 {
		return nil
	}
	return results[0]
}
