package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"image"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Format is just the lowercase extension including the '.' prefix.
type Format string

func (f Format) Unmarshal(buf []byte, ptr any) error {
	switch f {
	case ".json", ".js", ".masite":
		return json.Unmarshal(buf, ptr)
	case ".xml", ".mas", ".maxite":
		return xml.Unmarshal(buf, ptr)
	default:
		return errors.New("format not supported: " + string(f))
	}
}

func (f Format) Marshal(ptr any) ([]byte, error) {
	switch f {
	case ".json", ".js", ".masite":
		return json.Marshal(ptr)
	case ".xml", ".mas", ".maxite":
		return xml.Marshal(ptr)
	default:
		return nil, errors.New("format not supported: " + string(f))
	}
}

type Flag byte

const (
	FlagExtended       Flag = 1
	FlagHorizontalFlip Flag = 2
	FlagVerticalFlip   Flag = 4
	FlagSpritePalette  Flag = 8
	FlagOnTop          Flag = 16
	FlagSolid          Flag = 32
	FlagBless          Flag = 64
	FlagHarm           Flag = 128
)

type Cell struct {
	Index byte `json:"index" xml:"index,attr"`
	Flag  Flag `json:"flag" xml:"flag,attr"`
}

type Row struct {
	Cells []Cell `json:"cells" xml:"cells"`
}

type Map struct {
	Width  int    `json:"width" xml:"width,attr"`
	Height int    `json:"height" xml:"height,,attr"`
	Tw     int    `json:"tw" xml:"tw,attr"`
	Th     int    `json:"th" xml:"th,attr"`
	Offset int    `json:"offset" xml:"offset,attr"`
	From   string `json:"from" xml:"from,attr"` // From where to load the images tiles.
	Rows   []Row  `json:"rows" xml:"rows"`      // Rows.

	Surface *Surface `json:"-" xml:"-"` // Ebiten Surface for display.
}

func FormatFor(name string) Format {
	return Format(strings.ToLower(filepath.Ext(name)))
}

const TW = 8
const TH = 8

func NewMap(w, h int, from string) (*Map, error) {
	res := &Map{Width: w, Height: h, Th: TH, Tw: TW}
	err := res.LoadSurface(from)
	if err != nil {
		return nil, err
	}

	for y := 0; y < h; y++ {
		cells := make([]Cell, w)
		row := Row{Cells: cells}
		res.Rows = append(res.Rows, row)
	}
	return res, nil
}

func LoadMap(from string) (*Map, error) {
	buf, err := os.ReadFile(from)
	if err != nil {
		println(from, err.Error())
		return nil, err
	}
	res := &Map{}
	err = FormatFor(from).Unmarshal(buf, res)
	if err != nil {
		return nil, err
	}
	err = res.LoadSurface(res.From)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *Map) LoadSurface(name string) error {
	img, err := LoadSurface(FromName(name))
	if err != nil {
		return errors.Join(errors.New("Cannot load image:"+name), err)
	}
	m.From = name
	m.Surface = img
	return nil
}

func (m *Map) ToTile(at Point, camera Rectangle) Point {
	off := at.Sub(camera.Min)
	return image.Pt(off.X/m.Tw, off.Y/m.Th)
}

func (m *Map) Put(atTile Point, cell Cell) {
	if atTile.X < 0 || atTile.X >= m.Width {
		return
	}
	if atTile.Y < 0 || atTile.Y >= m.Height {
		return
	}
	m.Rows[atTile.Y].Cells[atTile.X] = cell
}

func (m *Map) Save(to string) error {
	buf, err := FormatFor(to).Marshal(m)
	if err != nil {
		return err
	}
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(buf)
	return err
}

func (m *Map) Render(screen *Surface, camera Rectangle) {
	ab := m.Surface.Bounds()

	starty := camera.Min.Y / m.Th
	if starty < 0 {
		starty = 0
	}
	endy := min(camera.Max.Y/m.Th, len(m.Rows)-1)

	// This draws the whole layer. Only draw visible part using a camera.
	for ty := starty; ty < endy; ty++ {
		row := m.Rows[ty]

		startx := max(camera.Min.X/m.Tw, 0)
		endx := min(camera.Max.X/m.Tw, len(row.Cells)-1)
		for tx := startx; tx < endx; tx++ {
			cell := row.Cells[tx]
			id := int(cell.Index)
			if cell.Flag&FlagExtended != 0 {
				id += 255
			}
			idx := id % ab.Dx()
			idy := id / ab.Dx()
			fx := idx * m.Tw
			fy := idy * m.Th

			from := image.Rect(fx, fy, fx+m.Tw, fy+m.Th)
			sub := m.Surface.SubImage(from).(*Surface)
			opts := ebiten.DrawImageOptions{}
			if cell.Flag&FlagHorizontalFlip != 0 {
				opts.GeoM.Scale(-1, 1)
				opts.GeoM.Translate(float64(m.Tw), 0)
			}
			if cell.Flag&FlagVerticalFlip != 0 {
				opts.GeoM.Scale(1, -1)
				opts.GeoM.Translate(0, float64(m.Th))
			}
			opts.GeoM.Translate(
				float64(int(tx)*m.Tw-camera.Min.X),
				float64(int(ty)*m.Th-camera.Min.Y),
			)

			if sub != nil {
				screen.DrawImage(sub, &opts)
			}
		}
	}
}
