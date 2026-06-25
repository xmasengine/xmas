package xgal

import (
	"image"
	"io/fs"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

// Pixels loads an image file from fsys as an [image.Image].
func Pixels(fsys fs.FS, name string) (image.Image, error) {
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

// Bake converts an [image.Image] to a [Surface] by baking its pixels onto the GPU.
func Bake(img image.Image) *Surface {
	return ebiten.NewImageFromImage(img)
}

// Scoop converts a [Surface] to an [image.Image] by scooping its pixels from the GPU.
func Scoop(surf *Surface) image.Image {
	w, h := surf.Size()
	pixels := make([]byte, w*h*4)
	surf.ReadPixels(pixels)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, pixels)
	return img
}
