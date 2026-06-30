// package xbin implements an efficient but somewhat flexible
// binary formatting package
package xbin

import "encoding/binary"
import "io"

const IDLength = 8

type ID [IDLength]byte

type Block struct {
	ID     ID
	Data   []byte
	Blocks []Block
}

func try(err error) func(error) {
	return func(check error) {
		if check != nil {
			err = check
		}
	}
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
	bwr.Read(size)
	bwr.Read(count)
	if size > 0 {
		data := make([]byte, size)
		bwr.Read(data)
	}
	if bwr.Err != nil {
		return bwr.Err
	}
	blocks := make([]Block, count)

	for _, block := range blocks {
		err = block.Decode(rd)
		if err != nil {
			return err
		}
	}

	return nil
}
