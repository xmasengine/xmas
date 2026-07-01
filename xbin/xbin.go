// package xbin implements an efficient but somewhat flexible
// binary formatting package based on binary encodable tress with an ID.
package xbin

import "encoding/binary"
import "io"
import "bytes"

// IDLength is the fixed length of the ID.
const IDLength = 8

// ByteOrder for xbin is always BigEndian.
var ByteOrder = binary.BigEndian

type ID [IDLength]byte

func (i ID) String() string {
	idx := bytes.IndexByte(i[:], 0)
	if idx >= 0 {
		return string(i[:idx])
	}
	return string(i[:])
}

func MakeID(s string) ID {
	var res ID
	copy(res[:], []byte(s))
	return res
}

// Tree allows for flexible binary data format that can contain tagged
// sub trees and data on each tree node.
type Tree struct {
	ID    ID
	Data  []byte
	Trees []Tree
}

func (b Tree) String() string {
	s := "<" + b.ID.String() + ">\n"
	s += string(b.Data)
	for _, sub := range b.Trees {
		s += sub.String()
	}
	s += "\n</" + b.ID.String() + ">\n"
	return s
}

func Make(id string, data []byte, blocks ...Tree) Tree {
	res := Tree{}
	copy(res.ID[:], []byte(id))
	res.Data = data
	res.Trees = blocks
	return res
}

func (b *Tree) Append(c Tree) int {
	b.Trees = append(b.Trees, c)
	return len(b.Trees) - 1
}

func (b *Tree) Add(id string, data []byte, blocks ...Tree) int {
	block := Make(id, data, blocks...)
	return b.Append(block)
}

type binWriter struct {
	wr        io.Writer
	Err       error
	ByteOrder binary.ByteOrder
}

func (b *binWriter) Write(v any) error {
	if b.Err == nil {
		b.Err = binary.Write(b.wr, b.ByteOrder, v)
	}
	return b.Err
}

func (b Tree) Encode(wr io.Writer) (err error) {
	size := uint32(len(b.Data))
	count := uint32(len(b.Trees))
	bwr := binWriter{wr: wr, ByteOrder: ByteOrder}

	bwr.Write(b.ID[:])
	bwr.Write(size)
	bwr.Write(count)
	bwr.Write(b.Data)
	if bwr.Err != nil {
		return bwr.Err
	}

	for _, block := range b.Trees {
		err = block.Encode(wr)
		if err != nil {
			return err
		}
	}

	return nil
}

type binReader struct {
	rd        io.Reader
	Err       error
	ByteOrder binary.ByteOrder
}

func (b *binReader) Read(v any) error {
	if b.Err == nil {
		b.Err = binary.Read(b.rd, b.ByteOrder, v)
	}
	return b.Err
}

func (b *Tree) Decode(rd io.Reader) (err error) {
	var size, count uint32
	var id ID

	bwr := binReader{rd: rd, ByteOrder: ByteOrder}

	bwr.Read(id[:])
	bwr.Read(&size)
	bwr.Read(&count)
	var data []byte
	if size > 0 {
		data = make([]byte, size)
		bwr.Read(data)
	}
	if bwr.Err != nil {
		return bwr.Err
	}
	var blocks []Tree
	if count > 0 {

		blocks = make([]Tree, count)

		for i, block := range blocks {
			err = block.Decode(rd)
			if err != nil {
				return err
			}
			blocks[i] = block
		}
	}
	b.ID = id
	b.Data = data
	b.Trees = blocks

	return nil
}

func (b *Tree) EncodeData(data any) error {
	var err error
	b.Data, err = binary.Append(b.Data, ByteOrder, data)
	return err
}

func (b *Tree) DecodeData(data any) error {
	var err error
	_, err = binary.Decode(b.Data, ByteOrder, data)
	return err
}

// Finds the ID in a block in a breadth first search.
// Also considers the tag of the block itself.
func (b Tree) FindID(id ID) (Tree, bool) {
	var res Tree
	found := make(chan Tree, len(b.Trees)*2)
	found <- b
	for len(found) > 0 {
		elt := <-found
		if elt.ID == id {
			return elt, true
		}
		for _, sub := range elt.Trees {
			if sub.ID == id {
				return sub, true
			}
			found <- sub
		}
	}
	return res, false
}
