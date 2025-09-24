// package xmap contains the map and world system of the xmas engine
// We use simple non namespaced XML as the disk format.
package xmap

import "encoding/xml"
import "io"
import "os"
import "io/fs"

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
	W       int      `xml:"w,attr"`
	H       int      `xml:"h,attr"`
	Src     string   `xml:"src,attr"`
	Tw      int      `xml:"tw,attr"`
	Th      int      `xml:"th,attr"`
	Objects []Object `xml:"object"`
	Rows    []Row    `xml:"row"`
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
}

func MakeRow(w int) Row {
	r := Row{}
	r.Cells = make([]Cell, w)
	return r
}

func MakeLayer(w, h int) Layer {
	l := Layer{W: w, H: h, Tw: 32, Th: 32}
	for y := 0; y < h; y++ {
		r := MakeRow(w)
		l.Rows = append(l.Rows, r)
	}
	return l
}

func NewZone(name string, w, h int) *Zone {
	l := MakeLayer(w, h)
	return &Zone{Layers: []Layer{l}, W: w, H: h, Name: name}
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
