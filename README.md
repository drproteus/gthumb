# gthumb

go thumbnail and static website gallery generator

gthumb scans a directory for JPEG photos, applies their EXIF orientation, writes
square thumbnails and display-size copies, and produces a single-page static
gallery (`index.html`) that can be hosted anywhere.

## Usage

```sh
go run . -src examples -out public
```

Each directory under `-src` (including `-src` itself) that contains JPEG files becomes an album. The output looks like:

```
public/
  index.html
  thumbs/<album>/<photo>.jpg   square thumbnails (center-cropped)
  photos/<album>/<photo>.jpg   display-size copies (aspect ratio kept)
```

### Flags

| Flag         | Default     | Description                                       |
|--------------|-------------|---------------------------------------------------|
| `-src`       | `examples`  | directory to scan for albums                      |
| `-out`       | `public`    | output directory                                  |
| `-templates` | `templates` | directory containing `index.gohtml`               |
| `-title`     | `gthumb`    | site title                                        |
| `-thumb`     | `400`       | thumbnail edge length in pixels                   |
| `-max`       | `2048`      | maximum display image edge length in pixels       |
| `-quality`   | `80`        | JPEG quality (1-100)                              |

## How it works

* **Orientation:** the EXIF `Orientation` tag (values 1-8) is read and applied so
  images are stored upright, including phone photos shot in portrait.
* **Metadata:** camera model and capture date (`DateTimeOriginal`) are shown in each
  photo's caption when present. Files without EXIF data are still processed.
* **Thumbnails:** scaled so the shorter edge matches `-thumb`, then center-cropped to a square.
* **Display images:** scaled down to fit within `-max` on the longest edge. Images are never upscaled.
* **Template:** `templates/index.gohtml` is a plain `html/template`, so the layout can be changed without rebuilding.

## Development

```sh
go test ./...
```

## Dependencies

* [bild](https://github.com/anthonynsimon/bild): image decoding, encoding, resizing, and flipping
* [go-exif](https://github.com/dsoprea/go-exif/tree/v3.0.1): EXIF parsing

## References
* https://exiftool.org/TagNames/EXIF.html
* https://pkg.go.dev/github.com/anthonynsimon/bild
* https://github.com/dsoprea/go-exif/tree/v3.0.1
* https://thumbsup.github.io/
* https://www.shopify.com/blog/image-sizes
