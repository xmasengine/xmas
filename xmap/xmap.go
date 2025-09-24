// package xmap contains the map and world system of the xmas engine
// We use simple non namespaced XML as the disk format.
package xmap

import "encoding/xml"
import "io"
import "os"
import "io/fs"
import "image"
import "image/color"
import "log/slog"

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

import (
	"github.com/xmasengine/xmas/xres"
)

// Rectangle is used for sizes and positions.
type Rectangle = image.Rectangle

// Point is used for position and offsets.
type Point = image.Point

// Color is a color.
type Color = color.Color

// Image is an image.Image
type Image = image.Image

// Surface is an ebiten.Image
type Surface = ebiten.Image

// RGBA is an RGBA color.
type RGBA = color.RGBA

// Face is a font face
type Face = text.Face

type Background struct {
	Name string `xml:"name,attr"`
}

type Foe struct {
	Name string `xml:"name,attr"`
	X    int    `xml:"x,attr"`
	Y    int    `xml:"y,attr"`
}

type Hidden struct {
	Name string `xml:"name,attr"`
	X    int    `xml:"x,attr"`
	Y    int    `xml:"y,attr"`
}

type Object struct {
	Foe    *Foe    `xml:"foe"`
	Hidden *Hidden `xml:"hidden"`
}

type Cell struct {
	Index int `xml:"i,attr"`
}

type Row struct {
	Cells []Cell `xml:"c"`
}

type Layer struct {
	W       int           `xml:"w,attr"`
	H       int           `xml:"h,attr"`
	Src     string        `xml:"src,attr"`
	Tw      int           `xml:"tw,attr"`
	Th      int           `xml:"th,attr"`
	Objects []Object      `xml:"object"`
	Rows    []Row         `xml:"row"`
	Atlas   *ebiten.Image `xml:"-"`
}

// Zone is a zone or level of the xmas engine.
type Zone struct {
	XMLName    xml.Name   `xml:"zone"`
	Name       string     `xml:"name,attr"`
	W          int        `xml:"w,attr"`
	H          int        `xml:"h,attr"`
	Script     string     `xml:"script"`
	Background Background `xml:"background"`
	Layers     []Layer    `xml:"layer"`
	Camera     Rectangle  `xml:"-"`
}

func MakeRow(w int) Row {
	r := Row{}
	r.Cells = make([]Cell, w)
	return r
}

func MakeLayer(w, h int) Layer {
	l := Layer{W: w, H: h, Tw: 16, Th: 16}
	for y := 0; y < h; y++ {
		r := MakeRow(w)
		l.Rows = append(l.Rows, r)
	}
	return l
}

func (l *Layer) LoadSource(src string) error {
	img, err := xres.LoadImageFromFile(src)
	if err != nil {
		slog.Error("Could not load layer source", "err", err)
		return err
	}
	if l.Atlas != nil {
		l.Atlas.Deallocate()
		l.Atlas = nil
	}
	l.Atlas = img
	l.Src = src
	return nil
}

func (l *Layer) SetlIndex(p Point, idx int) {
	if p.X < 0 || p.Y < 0 || p.Y >= len(l.Rows) || p.X >= len(l.Rows[p.Y].Cells) {
		return
	}
	l.Rows[p.Y].Cells[p.X].Index = idx
}

func (l *Layer) FillIndex(r Rectangle, idx int) {
	r = r.Canon()
	if r.Min.X < 0 || r.Min.Y < 0 ||
		r.Max.Y >= len(l.Rows) || r.Max.X >= len(l.Rows[r.Max.Y].Cells) {
		// XXX the guard condition is a bit weak.
		return
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			l.Rows[y].Cells[x].Index = idx
		}
	}
}

func (l *Layer) Render(screen *ebiten.Image, camera Rectangle) {
	if l.Atlas == nil {
		slog.Warn("Trying to render a layer without an Atlas")
		return
	}
	ab := l.Atlas.Bounds()

	starty := camera.Min.Y / l.Th
	if starty < 0 {
		starty = 0
	}
	endy := min(camera.Max.Y/l.Th, len(l.Rows)-1)

	// This draws the whole layer. Only draw visible part using a camera.
	for ty := starty; ty < endy; ty++ {
		row := l.Rows[ty]

		startx := max(camera.Min.X/l.Tw, 0)
		endx := min(camera.Max.X/l.Tw, len(row.Cells)-1)
		for tx := startx; tx < endx; tx++ {
			cell := row.Cells[tx]
			id := cell.Index
			idx := id % ab.Dx()
			idy := id / ab.Dx()
			fx := idx * l.Tw
			fy := idy * l.Th

			from := image.Rect(fx, fy, fx+l.Tw, fy+l.Th)
			sub := l.Atlas.SubImage(from).(*Surface)

			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(
				float64(int(tx)*l.Tw-camera.Min.X),
				float64(int(ty)*l.Th-camera.Min.Y),
			)
			if sub != nil {
				screen.DrawImage(sub, &opts)
			}
		}
	}
}

func (z *Zone) AddLayer(l Layer) *Zone {
	z.Layers = append(z.Layers, l)
	return z
}

func NewZone(name string, w, h int) *Zone {
	l := MakeLayer(w, h)
	cam := image.Rect(0, 0, 16*32, 16*32)
	z := &Zone{W: w, H: h, Name: name, Camera: cam}
	return z.AddLayer(l)
}

func ReadZone(rd io.Reader) (*Zone, error) {
	z := &Zone{}
	dec := xml.NewDecoder(rd)
	err := dec.Decode(z)
	if err != nil {
		return nil, err
	}
	return z, nil
}

func LoadZone(cb func() (io.Reader, error)) (*Zone, error) {
	rd, err := cb()
	if err != nil {
		return nil, err
	}
	return ReadZone(rd)
}

func (z *Zone) Write(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	enc.Indent("", "    ")
	defer enc.Close()
	err := enc.Encode(z)
	if err != nil {
		return err
	}
	return nil
}

func (z *Zone) Save(cb func() (io.Writer, error)) error {
	wr, err := cb()
	if err != nil {
		return err
	}
	return z.Write(wr)
}

func (z *Zone) Draw(screen *Surface) {
	for _, layer := range z.Layers {
		layer.Render(screen, z.Camera)
	}
}

func FromName(name string) func() (io.Reader, error) {
	return func() (io.Reader, error) {
		return os.Open(name)
	}
}

func FromRoot(root *os.Root, name string) func() (io.Reader, error) {
	return func() (io.Reader, error) {
		return root.Open(name)
	}
}

func FromFS(sys fs.FS, name string) func() (io.Reader, error) {
	return func() (io.Reader, error) {
		return sys.Open(name)
	}
}

func ToRoot(root *os.Root, name string) func() (io.Writer, error) {
	return func() (io.Writer, error) {
		return root.Create(name)
	}
}

func ToName(name string) func() (io.Writer, error) {
	return func() (io.Writer, error) {
		return os.Create(name)
	}
}
