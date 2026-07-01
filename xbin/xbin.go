// package xbin implements an efficient but somewhat flexible
// binary formatting package
package xbin

import "encoding/binary"
import "io"
import "bytes"

const IDLength = 8

type ID [IDLength]byte

func (i ID) String() string {
	idx := bytes.IndexByte(i[:], 0)
	if idx >= 0 {
		return string(i[:idx])
	}
	return string(i[:])
}

type Block struct {
	ID     ID
	Data   []byte
	Blocks []Block
}

func (b Block) String() string {
	s := "<" + b.ID.String() + ">\n"
	s += string(b.Data)
	for _, sub := range b.Blocks {
		s += sub.String()
	}
	s += "\n</" + b.ID.String() + ">\n"
	return s
}

func Make(id string, data []byte, blocks ...Block) Block {
	res := Block{}
	copy(res.ID[:], []byte(id))
	res.Data = data
	res.Blocks = blocks
	return res
}

func (b *Block) Append(c Block) int {
	b.Blocks = append(b.Blocks, c)
	return len(b.Blocks) - 1
}

func (b *Block) Add(id string, data []byte, blocks ...Block) int {
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

func (b Block) Encode(wr io.Writer) (err error) {
	size := uint32(len(b.Data))
	count := uint32(len(b.Blocks))
	bwr := binWriter{wr: wr, ByteOrder: binary.BigEndian}

	bwr.Write(b.ID[:])
	bwr.Write(size)
	bwr.Write(count)
	bwr.Write(b.Data)
	if bwr.Err != nil {
		return bwr.Err
	}

	for _, block := range b.Blocks {
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

func (b *Block) Decode(rd io.Reader) (err error) {
	var size, count uint32
	var id ID

	bwr := binReader{rd: rd, ByteOrder: binary.BigEndian}

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
	var blocks []Block
	if count > 0 {

		blocks = make([]Block, count)

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
	b.Blocks = blocks

	return nil
}
