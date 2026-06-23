package main

import "github.com/hajimehoshi/ebiten/v2"
import "io/fs"
import "io"
import _ "image/png"
import _ "image/jpeg"
import _ "image/gif"
import "image"
import "os"

func FromName(name string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return os.Open(name)
	}
}

func FromRoot(root *os.Root, name string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return root.Open(name)
	}
}

func FromFS(sys fs.FS, name string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return sys.Open(name)
	}
}

func ToRoot(root *os.Root, name string) func() (io.WriteCloser, error) {
	return func() (io.WriteCloser, error) {
		return root.Create(name)
	}
}

func ToName(name string) func() (io.WriteCloser, error) {
	return func() (io.WriteCloser, error) {
		return os.Create(name)
	}
}

func FromFirst(cbs ...func() (io.Reader, error)) func() (io.Reader, error) {
	return func() (rd io.Reader, err error) {
		for _, cb := range cbs {
			rd, err = cb()
			if err != nil && rd != nil {
				return rd, nil
			}
		}
		return nil, err
	}
}

func ToFirst(cbs ...func() (io.Writer, error)) func() (io.Writer, error) {
	return func() (wr io.Writer, err error) {
		for _, cb := range cbs {
			wr, err = cb()
			if err != nil && wr != nil {
				return wr, nil
			}
		}
		return nil, err
	}
}

func FirstAvailable[T any](cbs ...func() (T, error)) func() (T, error) {
	return func() (rd T, err error) {
		for _, cb := range cbs {
			rd, err = cb()
			if err != nil {
				return rd, nil
			}
		}
		return rd, err
	}
}

func DecodeSurface(rd io.Reader) (*Surface, error) {
	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}
	eimg := ebiten.NewImageFromImage(img)
	return eimg, nil
}

func DecodeSurfaceAndImage(rd io.Reader) (*Surface, Image, error) {
	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, nil, err
	}
	eimg := ebiten.NewImageFromImage(img)
	return eimg, img, nil
}

func LoadSurface(cb func() (io.ReadCloser, error)) (*Surface, error) {
	rd, err := cb()
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return DecodeSurface(rd)
}

func LoadSurfaceAndImage(cb func() (io.ReadCloser, error)) (*Surface, Image, error) {
	rd, err := cb()
	if err != nil {
		return nil, nil, err
	}
	defer rd.Close()
	return DecodeSurfaceAndImage(rd)
}
