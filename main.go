package main

import (
	"github.com/h2non/bimg"
)

func main() {
	buffer, err := bimg.Read("examples/image.jpg")
	if err != nil {
		panic(err)
	}
	newImage, err := bimg.NewImage(buffer).Thumbnail(300)
	bimg.Write("examples/thumb.jpg", newImage)
}
