package xdat

import (
	"encoding/xml"
	"io"
	"io/fs"
	"os"
)

import (
	"github.com/xmasengine/xmas/xgal"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Pose int

const PoseInterval = 2

type Player struct {
	Name           string        `xml:"name,attr"`           // Name of the character.
	Source         string        `xml:"src,attr"`            // Source file name to load the player's Texture from.
	PortraitSource string        `xml:"psrc,attr,omitempty"` // Source file name to load the player's Portrait from.
	Portrait       *xgal.Surface `xml:"-"`                   // The tile texture for this player if loaded.
	Texture        *xgal.Surface `xml:"-"`                   // The fortarit texture for this player if loaded.
	Sprite         *xgal.Surface `xml:"-"`                   // Sub graph to draw the player's sprite from.

	Bound xgal.Rectangle `xml:"-"` // Bound is the rendering bound (where to draw).
	Hit   xgal.Rectangle `xml:"-"` // Hit is the hit box (where the player "is").

	Direction Direction `xml:"-"`
	Pose      Pose      `xml:"-"`

	Depth uint16 `xml:"-"` // Depth is layer the player is "on".
}

func (p *Player) loadTexture(fsys fs.FS) error {
	if p.Source == "" {
		return nil
	}

	texture, err := xgal.Texture(fsys, p.Source)
	if err != nil {
		return err
	}
	if p.Texture != nil {
		p.Texture.Deallocate()
	}
	p.Texture = texture
	return nil
}

type World struct {
	XMLName   xml.Name  `xml:"world"`
	Name      string    `xml:"name,attr"`           // Name of the world.
	Start     string    `xml:"start,attr"`          // Start map.
	Copyright string    `xml:"copyright,omitempty"` // Copyright.
	License   string    `xml:"license,omitempty"`   // License.
	Credits   string    `xml:"credits,omitempty"`   // Credits.
	Players   []*Player `xml:"player"`              // Player character data.
}

func NewWorld(name string) *World {
	w := &World{}
	w.XMLName.Local = "world"
	w.Name = name
	return w
}

func (w World) SaveTo(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	enc.Indent("", " ")
	return enc.Encode(w)
}

func (w World) SaveFile(name string) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()
	return w.SaveTo(out)
}

func loadWorldFrom(rd io.Reader) (*World, error) {
	dec := xml.NewDecoder(rd)
	var world World
	err := dec.Decode(&world)
	return &world, err
}

func (w *World) loadTextures(fsys fs.FS) error {
	for _, player := range w.Players {
		err := player.loadTexture(fsys)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadWorld(fsys fs.FS, name string) (*World, error) {
	fin, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	world, err := loadWorldFrom(fin)
	if err != nil {
		return nil, err
	}
	err = world.loadTextures(fsys)
	if err != nil {
		return nil, err
	}

	if world.Players[0].Texture == nil {
		println("texture missing")
	}

	return world, err
}
