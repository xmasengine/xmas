package xres

import (
	"github.com/hajimehoshi/ebiten/v2"
)

import (
	"encoding/xml"
	"image"
	"image/color"
	_ "image/png"
	"io/fs"
	"slices"
)

// Icon is a sub image stored in an IconAtlas
type Icon struct {
	Name   string        `xml:"name,attr"`
	X      int           `xml:"x,attr"`
	Y      int           `xml:"y,attr"`
	Width  int           `xml:"width,attr"`
	Height int           `xml:"height,attr"`
	Image  *ebiten.Image `xml:"-"`
}

func (i Icon) Draw(screen *ebiten.Image, r image.Rectangle) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(float64(r.Dx())/float64(i.Width), float64(r.Dy())/float64(i.Height))
	opts.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	screen.DrawImage(i.Image, &opts)
}

func LoadImageAndPalette(from fs.FS, name string) (*ebiten.Image, color.Palette, error) {
	rd, err := from.Open(name)
	if err != nil {
		return nil, nil, err
	}
	if rd == nil {
		return nil, nil, fs.ErrNotExist
	}
	defer rd.Close()
	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, nil, err
	}

	var pal color.Palette
	pimg, ok := img.(image.PalettedImage)
	if ok {
		mod := pimg.ColorModel()
		pal, _ = mod.(color.Palette)
	}

	return ebiten.NewImageFromImage(img), pal, nil
}

// IconAtlas is a specific atlas for icons stored in both
// an XML and a image file.
type IconAtlas struct {
	Path  string        `xml:"path,attr"` // Path to the image file.
	Icons []*Icon       `xml:"Icon"`      // Icons defined in the XML file.
	Image *ebiten.Image `xml:"-"`
}

func (a IconAtlas) DrawIcon(screen *ebiten.Image, name string, r image.Rectangle) {
	icon := a.GetIcon(name)
	if icon == nil || icon.Image == nil {
		return
	}
	icon.Draw(screen, r)
}

func LoadIconAtlas(from fs.FS, name string) (*IconAtlas, error) {
	buf, err := fs.ReadFile(from, name)
	if err != nil {
		return nil, err
	}

	atlas := &IconAtlas{}
	err = xml.Unmarshal(buf, atlas)
	if err != nil {
		return nil, err
	}
	atlas.Image, _, err = LoadImageAndPalette(from, atlas.Path)
	for _, icon := range atlas.Icons {
		r := image.Rect(icon.X, icon.Y, icon.X+icon.Width, icon.Y+icon.Height)
		sub := atlas.Image.SubImage(r)
		icon.Image = sub.(*ebiten.Image)
	}

	return atlas, nil
}

func (a IconAtlas) GetIcon(name string) *Icon {
	idx := slices.IndexFunc(a.Icons, func(i *Icon) bool {
		return i != nil && i.Name == name
	})
	if idx < 0 {
		return nil
	}
	return a.Icons[idx]
}
