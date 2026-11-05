package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var opts Options
	flag.StringVar(&opts.Src, "src", "examples", "directory to scan for albums of JPEG images")
	flag.StringVar(&opts.Out, "out", "public", "output directory for the generated site")
	flag.StringVar(&opts.Templates, "templates", "templates", "directory containing index.gohtml")
	flag.StringVar(&opts.Title, "title", "gthumb", "site title")
	flag.IntVar(&opts.ThumbSize, "thumb", 400, "thumbnail edge length in pixels (square crop)")
	flag.IntVar(&opts.MaxSize, "max", 2048, "maximum edge length in pixels for display images")
	flag.IntVar(&opts.Quality, "quality", 80, "JPEG quality (1-100)")
	flag.Parse()

	if err := Build(opts); err != nil {
		fmt.Fprintln(os.Stderr, "gthumb:", err)
		os.Exit(1)
	}
	fmt.Println("Gallery written to", opts.Out)
}
