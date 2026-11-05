package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Options struct {
	Src       string
	Out       string
	Templates string
	Title     string
	ThumbSize int
	MaxSize   int
	Quality   int
}

// previewLimit is how many photos per album are shown before the expand button.
const previewLimit = 10

type Photo struct {
	Name   string
	Thumb  string
	Full   string
	Camera string
	Taken  string
	Hidden bool
}

type Album struct {
	Title  string
	Photos []Photo
	Hidden int
}

type Page struct {
	Title  string
	Albums []Album
	Total  int
}

// Build scans opts.Src for albums (directories containing JPEG files), writes
// thumbnails and display-size images under opts.Out, and renders index.html.
func Build(opts Options) error {
	absSrc, err := filepath.Abs(opts.Src)
	if err != nil {
		return err
	}
	absOut, err := filepath.Abs(opts.Out)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(filepath.Join(opts.Templates, "index.gohtml"))
	if err != nil {
		return fmt.Errorf("loading template: %w", err)
	}

	page := Page{Title: opts.Title}
	err = filepath.WalkDir(absSrc, func(dir string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if dir == absOut {
			return filepath.SkipDir
		}
		album, err := buildAlbum(opts, absSrc, dir)
		if err != nil {
			return err
		}
		if album != nil {
			page.Albums = append(page.Albums, *album)
			page.Total += len(album.Photos)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := os.MkdirAll(opts.Out, 0o755); err != nil {
		return err
	}
	out, err := os.Create(filepath.Join(opts.Out, "index.html"))
	if err != nil {
		return err
	}
	defer out.Close()
	return tmpl.ExecuteTemplate(out, "index.gohtml", page)
}

// buildAlbum processes the JPEG files directly inside dir. It returns nil when dir has none.
func buildAlbum(opts Options, absSrc, dir string) (*Album, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(absSrc, dir)
	if err != nil {
		return nil, err
	}
	urlDir := ""
	title := opts.Title
	if rel != "." {
		urlDir = filepath.ToSlash(rel)
		title = urlDir
	}

	album := Album{Title: title}
	for _, entry := range entries {
		if entry.IsDir() || !isJPEG(entry.Name()) {
			continue
		}
		photo, err := processPhoto(opts, filepath.Join(dir, entry.Name()), urlDir)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		album.Photos = append(album.Photos, photo)
	}
	if len(album.Photos) == 0 {
		return nil, nil
	}
	for i := previewLimit; i < len(album.Photos); i++ {
		album.Photos[i].Hidden = true
		album.Hidden++
	}
	return &album, nil
}

// processPhoto writes the display and thumbnail versions of one source image.
func processPhoto(opts Options, src, urlDir string) (Photo, error) {
	name := filepath.Base(src)
	meta := readMetadata(src)

	img, err := openImage(src)
	if err != nil {
		return Photo{}, err
	}
	img = applyOrientation(img, meta.Orientation)

	thumbURL := path.Join("thumbs", urlDir, name)
	fullURL := path.Join("photos", urlDir, name)
	thumbPath := filepath.Join(opts.Out, filepath.FromSlash(thumbURL))
	fullPath := filepath.Join(opts.Out, filepath.FromSlash(fullURL))

	if err := os.MkdirAll(filepath.Dir(thumbPath), 0o755); err != nil {
		return Photo{}, err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return Photo{}, err
	}
	if err := saveJPEG(thumbPath, squareThumbnail(img, opts.ThumbSize), opts.Quality); err != nil {
		return Photo{}, err
	}
	if err := saveJPEG(fullPath, fitWithin(img, opts.MaxSize), opts.Quality); err != nil {
		return Photo{}, err
	}

	return Photo{
		Name:   strings.TrimSuffix(name, filepath.Ext(name)),
		Thumb:  thumbURL,
		Full:   fullURL,
		Camera: meta.Camera,
		Taken:  meta.Taken,
	}, nil
}

func isJPEG(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg":
		return true
	}
	return false
}
