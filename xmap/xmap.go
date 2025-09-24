// package xmap contains the map and world system of the xmas engine
// We use simple non namespaced XML as the disk format.
package xmap

import "encoding/xml"
import "io"
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

// Direction is a direction a player or mobile may be facing.
type Direction int

const (
	South Direction = 1 + iota
	East
	North
	West
	None = 0
)

// Action describes the action a player or mobiel sprite is performing.
type Action int

const (
	Stand Action = 1 + iota
	Walk
	Lift
	Carry
	Attack
	Hide = 0
)

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
	Index      int         `xml:"i,attr"`
	Animations []Animation `xml:"animation"`
}

type Row struct {
	Cells []Cell `xml:"c"`
}

type SrcAtlas struct {
	Src   string        `xml:"src,attr"`
	Atlas *ebiten.Image `xml:"-"`
}

type Layer struct {
	W       int      `xml:"w,attr"`
	H       int      `xml:"h,attr"`
	Tw      int      `xml:"tw,attr"`
	Th      int      `xml:"th,attr"`
	Objects []Object `xml:"object"`
	Rows    []Row    `xml:"row"`
	SrcAtlas
}

type Animation struct {
	Frames int `xml:"frames,attr"`
	Phase  int `xml:"-"`
	Tick   int `xml:"-"`
}

type Pose struct {
	Name      string `xml:"name,attr"`
	Action    `xml:"action,attr"`
	Direction `xml:"dir,attr"`
	X         int `xml:"x,attr"`
	Y         int `xml:"y,attr"`
	H         int `xml:"w,attr"`
	W         int `xml:"h,attr"`
	Animation
}

type Player struct {
	Tw    int       `xml:"tw,attr"`
	Th    int       `xml:"th,attr"`
	At    Point     `xml:"-"`
	Dir   Point     `xml:"-"`
	Index int       `xml:"-"`
	Poses []Pose    `xml:"pose"` // Poses the player has in the SrcAtlas.
	Pose  `xml:"-"` // curent pose
	SrcAtlas
}

const FramesPerSecond = 60

func (a *Animation) Update() {
	if a.Phase < 0 {
		println("phase problems")
		return
	}

	if a.Frames < 1 {
		println("frame problems")
		return
	}

	a.Tick++
	if a.Tick >= FramesPerSecond {
		a.Phase++
		if a.Phase >= a.Frames {
			a.Phase = 0
		}
		a.Tick = 0
	}
}

func (p *Player) Update() {
	if p.Action == Hide {
		return
	}
	p.Animation.Update()
}

func (p Player) RenderPose(screen *ebiten.Image, camera Rectangle, pos Pose) bool {
	if p.Atlas == nil {
		slog.Warn("Trying to render a player without an Atlas")
		return false
	}
	if pos.Action == Hide {
		return false
	}

	fx := pos.X
	fy := pos.Y
	if pos.Phase > 0 {
		fx += pos.Phase * pos.W
	}
	from := image.Rect(fx, fy, fx+pos.W, fy+pos.H)
	sub := p.Atlas.SubImage(from).(*Surface)
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(int(p.At.X)-camera.Min.X),
		float64(int(p.At.Y)-camera.Min.Y),
	)
	if sub != nil {
		screen.DrawImage(sub, &opts)
	}
	return true
}

func (p Player) Render(screen *ebiten.Image, camera Rectangle) {
	if p.Atlas == nil {
		slog.Warn("Trying to render a player without an Atlas")
		return
	}
	if p.RenderPose(screen, camera, p.Pose) {
		return
	}

	ab := p.Atlas.Bounds()

	id := p.Index
	idx := id % ab.Dx()
	idy := id / ab.Dx()
	fx := idx * p.Tw
	fy := idy * p.Th
	from := image.Rect(fx, fy, fx+p.Tw, fy+p.Th)
	sub := p.Atlas.SubImage(from).(*Surface)
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(int(p.At.X)-camera.Min.X),
		float64(int(p.At.Y)-camera.Min.Y),
	)
	if sub != nil {
		screen.DrawImage(sub, &opts)
	}
}

func (p Player) BestPose(dir Direction, act Action) Pose {
	for _, pose := range p.Poses {
		if pose.Direction == dir && pose.Action == act {
			return pose
		}
	}
	return p.Poses[0]
}

func (p *Player) AddPose(pose Pose) Pose {
	p.Poses = append(p.Poses, pose)
	return pose
}

func (p *Player) AddNewPose(act Action, x, y, w, h, frames int, dir Direction) Pose {
	pose := &Pose{X: x, Y: y, W: w, H: h, Action: act, Direction: dir}
	pose.Frames = frames
	return p.AddPose(*pose)
}

// AddNewPoses adds poses for 4 directions which havet o be on top of each
// other in the order south, east, north, west
func (p *Player) AddNewPoses(act Action, x, y, w, h, frames int) {
	for dir := South; dir <= West; dir++ {
		p.AddNewPose(act, x, y, w, h, frames, dir)
		y += h
	}
}

// Zone is a zone or level of the xmas engine.
type Zone struct {
	XMLName    xml.Name   `xml:"zone"`
	Name       string     `xml:"name,attr"`
	W          int        `xml:"w,attr"`
	H          int        `xml:"h,attr"`
	Script     string     `xml:"script"`
	Background Background `xml:"background"`
	Player     Player     `xml:"player"`
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

func (s *SrcAtlas) LoadSource(src string) error {
	img, err := xres.LoadImageFromFile(src)
	if err != nil {
		slog.Error("Could not load layer source", "err", err)
		return err
	}
	if s.Atlas != nil {
		s.Atlas.Deallocate()
		s.Atlas = nil
	}
	s.Atlas = img
	s.Src = src
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

func LoadZone(cb func() (io.ReadCloser, error)) (*Zone, error) {
	rd, err := cb()
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return ReadZone(rd)
}

func (z *Zone) Write(wr io.WriteCloser) error {
	enc := xml.NewEncoder(wr)
	enc.Indent("", "    ")
	defer enc.Close()
	err := enc.Encode(z)
	if err != nil {
		return err
	}
	return nil
}

func (z *Zone) Save(cb func() (io.WriteCloser, error)) error {
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
	z.Player.Render(screen, z.Camera)
}

// Alias the selectable file system callbacks.
var (
	FromName = xres.FromName
	FromRoot = xres.FromRoot
	FromFS   = xres.FromFS
	ToRoot   = xres.ToRoot
	ToName   = xres.ToName
)
