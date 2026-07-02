// package xbin implements an efficient but somewhat flexible
// binary big endian formatting package based on binary encodable trees
// with an ID.
package xbin

import "encoding/binary"
import "io"
import "bytes"

// IDLength is the fixed length of the ID.
const IDLength = 8

// ByteOrder for xbin is always BigEndian.
var ByteOrder = binary.BigEndian

// ID is the identfier of a Tree.
type ID [IDLength]byte

// String implements the stringer interrface for ID. NUL bytes will be truncated.
func (i ID) String() string {
	idx := bytes.IndexByte(i[:], 0)
	if idx >= 0 {
		return string(i[:idx])
	}
	return string(i[:])
}

// MakeID makes an ID from a string. If too sort it will contain NUL bytes.
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

// String implements the striger interface for tree.
func (b Tree) String() string {
	s := "<" + b.ID.String() + ">\n"
	s += string(b.Data)
	for _, sub := range b.Trees {
		s += sub.String()
	}
	s += "\n</" + b.ID.String() + ">\n"
	return s
}

// Make constructs a tree with the given ID, data and sub trees.
func Make(id string, data []byte, blocks ...Tree) Tree {
	res := Tree{}
	res.ID = MakeID(id)
	res.Data = data
	res.Trees = blocks
	return res
}

// Append adds a tree to this tree. Returns the index where it is stored.
func (b *Tree) Append(c Tree) int {
	b.Trees = append(b.Trees, c)
	return len(b.Trees) - 1
}

// Add makes a new tre and append it. Returns the index where it is stored.
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

// Encode encodes the tree to an io.Writer.
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

// Decode decodes the tree from an io.Writer.
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

// EncodeData encodes fixed size data as per binary.Append to this tree's Data.
func (b *Tree) EncodeData(data any) error {
	var err error
	b.Data, err = binary.Append(b.Data, ByteOrder, data)
	return err
}

// DecodeData decodes fixed size data as per binary.Decode to this tree's Data.
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
