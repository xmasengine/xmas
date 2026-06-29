package xmap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var id = [4]byte{'x', 'm', 'a', 's'}

const version = 1

// str is a fixed size string with up to 256 characters.
type Str struct {
	Size uint16
	Data [255]rune
}

type Header struct {
	ID      [4]byte
	Version uint32
}

type Tile uint16

const LayerWidth = 64
const LayerHeight = 64
const ThingSprites = 16

type Row [LayerWidth]Tile
type Layer struct {
	Z     uint16
	W     uint16
	H     uint16
	Sheet uint16
	Rows  [LayerHeight]Row
}

type Kind int16
type Lock int16
type Key int16
type Talk int16

type Thing struct {
	Name    Str
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
	Name   Str
	X      uint16
	Y      uint16
	Size   uint16
	Layers [4]Layer
}

func (z Zone) SaveTo(wr io.Writer) error {
	head := Header{
		ID:      id,
		Version: version,
	}
	err := binary.Write(wr, binary.BigEndian, head)
	if err != nil {
		return err
	}
	err = binary.Write(wr, binary.BigEndian, z)
	if err != nil {
		return err
	}
	return nil
}

func LoadFrom(rd io.Reader) (*Zone, error) {
	head := Header{}
	err := binary.Read(rd, binary.BigEndian, &head)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(head.ID[:], id[:]) {
		return nil, errors.New("wrong ID in header")
	}
	if head.Version != version {
		return nil, errors.New("wrong version in header")
	}
	z := &Zone{}
	err = binary.Read(rd, binary.BigEndian, z)
	if err != nil {
		return nil, err
	}
	return z, nil
}
