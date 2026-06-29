package xgal

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"image/color"
	"image/gif"
	"image/jpeg"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Texture loads an image file from fsys as a [Surface].
func Texture(fsys fs.FS, name string) (*Surface, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// Pixels loads an image file from fsys as an [Image].
func Pixels(fsys fs.FS, name string) (Image, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Crayon loads an image from fsys as a PalettedImage
func Crayon(fsys fs.FS, name string) (PalettedImage, error) {
	img, err := Pixels(fsys, name)
	if err != nil {
		return nil, err
	}
	if pimg, ok := img.(PalettedImage); ok {
		return pimg, nil
	}
	return nil, errors.New(name + " is not a paletted image")
}

// Bake converts an [image.Image] to a [Surface].
func Bake(img image.Image) *Surface {
	return ebiten.NewImageFromImage(img)
}

// Scoop converts a [Surface] to an [image.Image].
func Scoop(surf *Surface) image.Image {
	w, h := surf.Size()
	pixels := make([]byte, w*h*4)
	surf.ReadPixels(pixels)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, pixels)
	return img
}

// Snap saves screen as a PNG file.
func Snap(screen *Surface, name string) error {
	w, h := screen.Size()
	pixels := make([]byte, w*h*4)
	screen.ReadPixels(pixels)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, pixels)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

type Encoder interface {
	Encode(w io.Writer, m image.Image) error
}

type gifEncoder struct{}

func (g gifEncoder) Encode(w io.Writer, m image.Image) error {
	pimg, ok := m.(PalettedImage)
	if ok {
		model := pimg.ColorModel()
		pmod, ok := model.(color.Palette)
		if ok {
			opts := gif.Options{
				NumColors: len(pmod),
			}
			return gif.Encode(w, m, &opts)
		}
	}
	return gif.Encode(w, m, nil)
}

type jpegEncoder struct{}

func (j jpegEncoder) Encode(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, nil)
}

// returns an encoder for the file name by extension or nil if not found.
func encoderFor(name string) Encoder {
	ext := filepath.Ext(name)
	switch ext {
	case ".png":
		return &png.Encoder{}
	case ".gif":
		return &gifEncoder{}
	case ".jpg", ".jpeg":
		return &jpegEncoder{}
	default:
		return nil
	}

}

// Scribble writes a PalettedImage to a file
func Scribble(name string, pimg PalettedImage) error {
	encoder := encoderFor(name)
	if encoder == nil {
		return errors.New("file format not supported: " + name)
	}
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return encoder.Encode(f, pimg)
}

// Express (like paint on a palette) creates an empty paletted image.
func Express(rect Rectangle, pal color.Palette) *Paletted {
	return image.NewPaletted(rect, pal)
}

// Reduce reduces an Image to a Paletted image.
// However if it already is a paletted image, just returs itself.
func Reduce(src Image, pal Palette) *Paletted {
	if pimg, ok := src.(*Paletted); ok {
		return pimg
	}
	b := src.Bounds()
	img := image.NewPaletted(b, pal)
	draw.FloydSteinberg.Draw(img, b, src, b.Min)
	return img
}
