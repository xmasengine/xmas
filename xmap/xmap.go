package xmap

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

const version = 1

type Header struct {
	ID      [4]byte
	Version uint32
}

type Tile uint16

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
	cdata = bytes.Trim(bytes.Clone(cdata), " ")

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
	rd.FieldsPerRecord = -1

	t.Rows = make([]Row, LayerHeight)
	for i := 0; i < LayerHeight; i++ {
		record, err := rd.Read()
		if err != nil {
			return err
		}
		t.Rows[i] = make([]Tile, LayerWidth)
		for j, field := range record {
			v, err := strconv.Atoi(field)
			if err != nil {
				return errors.Join(errors.New("In "+string(in)), err)
			}
			t.Rows[i][j] = Tile(v)
		}
	}
	return nil
}

type Layer struct {
	Z     uint16 `xml:"z,attr"`
	W     uint16 `xml:"w,attr"`
	H     uint16 `xml:"h,attr"`
	Sheet uint16 `xml:"sheet,attr"`
	Tiles Tiles  `xml:"tiles"`
}

type Kind int16
type Lock int16
type Key int16
type Talk int16

type Thing struct {
	Name    string
	Kind    Kind
	Talk    uint16
	Z       uint16
	X       uint16
	Y       uint16
	W       uint16
	H       uint16
	Sheet   uint16
	Sprites [ThingSprites]uint16
}

type Zone struct {
	XMLName xml.Name `xml:"zone"`
	Name    string   `xml:"name,attr"`
	Layers  []Layer  `xml:"layer"`
}

func NewZone(name string) *Zone {
	z := &Zone{}
	z.XMLName.Local = "zone"
	z.Name = name
	z.Layers = make([]Layer, LayerCount)
	for l := 0; l < LayerCount; l++ {
		z.Layers[l].Tiles.Rows = make([]Row, LayerHeight)
		for r := 0; r < LayerHeight; r++ {
			z.Layers[l].Tiles.Rows[r] = make([]Tile, LayerWidth)
		}
	}
	return z
}

func (z Zone) SaveTo(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	enc.Indent("", " ")
	return enc.Encode(z)
}

func LoadFrom(rd io.Reader) (*Zone, error) {
	dec := xml.NewDecoder(rd)
	var zone Zone
	err := dec.Decode(&zone)
	return &zone, err
}
