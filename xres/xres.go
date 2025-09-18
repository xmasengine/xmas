// xres is the xmas engine resource handling package.
package xres

import "github.com/hajimehoshi/ebiten/v2"
import "io/fs"
import "io"
import _ "image/png"
import _ "image/jpeg"
import _ "image/gif"
import "image"
import "os"

func DecodeImage(rd io.Reader) (*ebiten.Image, error) {
	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}
	eimg := ebiten.NewImageFromImage(img)
	return eimg, nil
}

func LoadImageFromFile(name string) (*ebiten.Image, error) {
	rd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return DecodeImage(rd)
}

func LoadImageFromFS(from fs.FS, name string) (*ebiten.Image, error) {
	rd, err := from.Open(name)
	if err != nil {
		println("open failed")
		return nil, err
	}
	if rd == nil {
		println("read nil")
		return nil, fs.ErrNotExist
	}
	defer rd.Close()
	return DecodeImage(rd)
}
