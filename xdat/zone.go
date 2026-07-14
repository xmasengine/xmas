// package xdat implements the engine's data structures and
// saving and loading of these structures, including resouurces such as
// tile or sprite images. It does not display or run them.
package xdat

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"os"
	"slices"
	"strconv"
)

import (
	"github.com/xmasengine/xmas/xgal"
)

const version = 1

type Header struct {
	ID      [4]byte
	Version uint32
}

type Flag uint32

const (
	FlagSolid Flag = 1 << (16 + iota)
	FlagSpecial
	FlagHarm
	FlagBless
	FlagHorizontal
	FlagVertical
	FlagRotate90
	FlagRotate180
	FlagRotate270
)

// Tile consists of upper 16 biths with the flags,
// middle 8 bits with the tile Y coordinate in tiles in the texture
// and low 8 bith with the x coordinate in tiles in the texture
type Tile uint32

func (t Tile) X() uint8 {
	return uint8(t & 0xff)
}

func (t Tile) Y() uint8 {
	return uint8((t >> 8) & 0xff)
}

func (t Tile) Has(f Flag) bool {
	return (Flag(t) & f) == f
}

func (t *Tile) Set(f Flag) Tile {
	(*t) = Tile(f) | *t
	return *t
}

func (t *Tile) Toggle(f Flag) Tile {
	(*t) = Tile(f) ^ *t
	return *t
}

const LayerCount = 4
const LayerWidth = 64
const LayerHeight = 64
const ThingSprites = 16

type Row []Tile
type Tiles struct {
	Rows []Row
}

func (t *Tiles) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tok, err := d.Token()
	if err != nil {
		return err
	}

	cdata, ok := tok.(xml.CharData)
	if !ok {
		return errors.New("expected character data")
	}

	// need to clone as per the contract of Token()
	// We also strip off surrounding spaces and a single beautifying
	// leading newline.
	cdata = bytes.TrimPrefix(bytes.Trim(bytes.Clone(cdata), " "), []byte{'\n'})

	tok, err = d.Token()
	if err != nil {
		return err
	}

	end, ok := tok.(xml.EndElement)
	if !ok || end.Name != start.Name {
		return errors.New("expected end of tag " + start.Name.Local)
	}

	err = t.UnmarshalText([]byte(cdata))
	if err != nil {
		return err
	}

	return nil
}

func (t Tiles) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	text, err := t.MarshalText()
	if err != nil {
		return err
	}
	// This newline is purely for aesthetic reasons.
	text = slices.Insert(text, 0, '\n')
	cdata := xml.CharData(text)
	end := xml.EndElement{Name: start.Name}

	tokens := []xml.Token{start, cdata, end}

	for _, tok := range tokens {
		err = e.EncodeToken(tok)
		if err != nil {
			return err
		}

	}
	return nil
}

var _ xml.Marshaler = &Tiles{}
var _ xml.Unmarshaler = &Tiles{}

func (t Tiles) MarshalText() ([]byte, error) {
	buf := &bytes.Buffer{}
	wr := csv.NewWriter(buf)
	wr.Comma = ' '

	for _, row := range t.Rows {
		record := make([]string, len(row))
		for j, cell := range row {
			record[j] = strconv.Itoa(int(cell))
		}
		wr.Write(record)
	}
	wr.Flush()
	return buf.Bytes(), nil
}

func (t *Tiles) UnmarshalText(in []byte) error {
	buf := bytes.NewBuffer(in)
	rd := csv.NewReader(buf)
	rd.Comma = ' '
	rd.Comment = '#'
	rd.TrimLeadingSpace = true
	rd.FieldsPerRecord = 0 // all rows must be equally long

	for {
		record, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		row := make([]Tile, len(record))
		for j, field := range record {
			v, err := strconv.Atoi(field)
			if err != nil {
				return errors.Join(errors.New("In "+string(in)), err)
			}
			row[j] = Tile(v)
		}
		t.Rows = append(t.Rows, row)
	}
	return nil
}

func (t Tiles) Contains(tx, ty int) bool {
	if tx < 0 || ty < 0 {
		return false
	}

	if ty >= len(t.Rows) {
		return false
	}
	if tx >= len(t.Rows[ty]) {
		return false
	}
	return true
}

func (t *Tiles) Set(at xgal.Point, cell Tile) bool {
	if !t.Contains(at.X, at.Y) {
		return false
	}
	t.Rows[at.Y][at.X] = cell
	return true
}

func (t Tiles) Get(at xgal.Point) Tile {
	if !t.Contains(at.X, at.Y) {
		return 0
	}
	return t.Rows[at.Y][at.X]
}

type Layer struct {
	Depth      int           `xml:"z,attr"`   // Depth is the depth position of the layer
	Width      int           `xml:"w,attr"`   // Width is the width expressed in tiles.
	Height     int           `xml:"h,attr"`   // Height is the height expressed in tiles.
	TileWidth  int           `xml:"tw,attr"`  // TileWidth is the width of the tiles in this layer.
	TileHeight int           `xml:"th,attr"`  // TileHeight is the height of the thiles in this layer.
	Source     string        `xml:"src,attr"` // Source file name to load the Leyare's Texture from.
	Tiles      Tiles         `xml:"tiles"`    // Tiles
	Texture    *xgal.Surface `xml:"-"`        // The tile texture for this layer if loaded.
}

// MakeLayer makes a layer with the default size and tile size.
func MakeLayer() Layer {
	return MakeLayerWith(LayerWidth, LayerHeight, 8, 8)
}

// MakeLayerWith makes a layer with the given parameters.
func MakeLayerWith(w, h, tw, th int) Layer {
	l := Layer{}
	l.Width = w
	l.Height = h
	l.TileWidth = tw
	l.TileHeight = th

	l.Tiles.Rows = make([]Row, l.Height)
	for r := 0; r < l.Height; r++ {
		l.Tiles.Rows[r] = make([]Tile, l.Width)
	}
	return l
}

func (l *Layer) SetSource(fsys fs.FS, src string) error {
	texture, err := xgal.Texture(fsys, src)
	if err != nil {
		return err
	}
	l.Texture = texture
	l.Source = src
	return nil
}

func (l *Layer) loadTexture(fsys fs.FS) error {
	if l.Source == "" {
		return nil
	}

	texture, err := xgal.Texture(fsys, l.Source)
	if err != nil {
		return err
	}
	if l.Texture != nil {
		l.Texture.Deallocate()
	}
	l.Texture = texture
	return nil
}

func (l *Layer) Contains(tx, ty int) bool {
	return l.Tiles.Contains(tx, ty)
}

func (l *Layer) Set(at xgal.Point, cell Tile) bool {
	return l.Tiles.Set(at, cell)
}

func (l Layer) Get(at xgal.Point) Tile {
	return l.Tiles.Get(at)
}

func (l *Layer) FloodFill(at xgal.Point, cell Tile) {
	now := l.Get(at)
	if now == cell {
		return // already ok
	}
	if !l.Contains(at.X, at.Y) {
		return
	}

	l.Set(at, cell)
	// this floodfill is recursive but the maps are small so
	// it should not cause problems.
	for dx := -1; dx <= 1; dx++ {
		at2 := at
		at2.X += dx
		now2 := l.Get(at2)
		if now2 == now {
			l.FloodFill(at2, cell)
		}
	}
	for dy := -1; dy <= 1; dy++ {
		at2 := at
		at2.Y += dy
		now2 := l.Get(at2)
		if now2 == now {
			l.FloodFill(at2, cell)
		}
	}
}

func (l *Layer) ToTile(at xgal.Point, camera xgal.Rectangle) xgal.Point {
	off := at.Sub(camera.Min)
	return xgal.Pt(off.X/l.TileWidth, off.Y/l.TileHeight)
}

type Kind int16
type Lock int16
type Key int16

type Thing struct {
	Name       string
	Kind       Kind
	Talk       string
	Sprites    [ThingSprites]uint16
	Depth      uint16        `xml:"z,attr"`   // Depth is the depth position of the layer
	Width      uint16        `xml:"w,attr"`   // Width is the width expressed in tiles.
	Height     uint16        `xml:"h,attr"`   // Height is the height expressed in tiles.
	TileWidth  uint16        `xml:"tw,attr"`  // TileWidth is the width of the tiles in this layer.
	TileHeight uint16        `xml:"th,attr"`  // TileHeight is the height of the thiles in this layer.
	Source     string        `xml:"src,attr"` // Source file name to load the Leyare's Texture from.
	Texture    *xgal.Surface `xml:"-"`        // The tile texture for this Thing if loaded.

}

type Zone struct {
	XMLName xml.Name `xml:"zone"`
	Name    string   `xml:"name,attr"`
	Layers  []Layer  `xml:"layer"`
	Talks   []Talk   `xml:"talk"`
}

func NewZone(name string) *Zone {
	z := &Zone{}
	z.XMLName.Local = "zone"
	z.Name = name
	for l := 0; l < LayerCount; l++ {
		layer := MakeLayer()
		z.Layers = append(z.Layers, layer)
	}
	return z
}

func (z Zone) SaveTo(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	enc.Indent("", " ")
	return enc.Encode(z)
}

func (z Zone) SaveFile(name string) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()
	return z.SaveTo(out)
}

func LoadFrom(rd io.Reader) (*Zone, error) {
	dec := xml.NewDecoder(rd)
	var zone Zone
	err := dec.Decode(&zone)
	return &zone, err
}

func (z *Zone) loadLayerTextures(fsys fs.FS) error {
	for _, layer := range z.Layers {
		err := layer.loadTexture(fsys)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadZone(fsys fs.FS, name string) (*Zone, error) {
	fin, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	zone, err := LoadFrom(fin)
	if err != nil {
		return nil, err
	}
	err = zone.loadLayerTextures(fsys)
	if err != nil {
		return nil, err
	}
	return zone, err
}

// Talk is a dialog
type Talk struct {
	Name  string    `xml:"name,attr"` // identifying name
	Speak []Speaker `xml:"speak`
}

type Speaker interface {
	Speak() string
}

// Say is a single speech expression or question
type Say struct {
	When string `xml:"expr,attr,omitempty"` // expression with condition
	Who  string `xml:"who,attr"`            // who is speaking
	Say  string `xml:"say`
}

func (s Say) Speak() string {
	return s.Say
}

// Ask is a speech question with multiple answers
type Ask struct {
	When    string `xml:"expr,attr,omitempty"` // expression with condition
	Who     string `xml:"name,attr"`           // who is speaking
	Ask     string `xml:"ask`
	Replies []Reply
}

func (a Ask) Speak() string {
	return a.Ask
}

type Reply struct {
	When  string `xml:"expr,attr,omitempty"` // expression with condition of reply
	Expr  string `xml:"expr,attr,omitempty"` // expression with value of reply
	Reply string `xml:"reply`
}

func (r Reply) Speak() string {
	return r.Reply
}

// If can be used for simple scripting with expressions.
type If struct {
	Expr string `xml:"expr,attr"`
	Then any    `xml:"expr,attr"`
}

// On can be used for simple event scripting with expressions.
type On struct {
	Expr string `xml:"expr,attr"`
	Body any
}

// Expr can replace itself with its expression value.
type Expr struct {
	Expr string `xml:"expr,attr"`
}
