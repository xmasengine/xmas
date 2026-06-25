package xgal

import (
	"bytes"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/zachomedia/go-bdf"
)

// Face is a font face.
type Face = text.Face

// BuiltinFace is the built-in 7×14 pixel monochrome bitmap font face.
var BuiltinFace Face = text.NewGoXFace(bitmapfont.Face)

// Typeface loads a font file from fsys as a [Face].
// Supported formats: BDF (.bdf), TrueType (.ttf), OpenType (.otf).
// A point size may be provided for TTF/OTF fonts by passing it as an optional
// argument; the default is 12. BDF fonts ignore the size.
func Typeface(fsys fs.FS, name string, size ...float64) (Face, error) {
	pt := 12.0
	if len(size) > 0 {
		pt = size[0]
	}

	buf, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, err
	}

	switch ext := strings.ToLower(filepath.Ext(name)); ext {
	case ".bdf":
		parsed, err := bdf.Parse(buf)
		if err != nil {
			return nil, err
		}
		return text.NewGoXFace(parsed.NewFace()), nil
	case ".ttf", ".otf":
		src, err := text.NewGoTextFaceSource(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		return &text.GoTextFace{Source: src, Size: pt}, nil
	default:
		if b, err := bdf.Parse(buf); err == nil {
			return text.NewGoXFace(b.NewFace()), nil
		}
		src, err := text.NewGoTextFaceSource(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		return &text.GoTextFace{Source: src, Size: pt}, nil
	}
}

// Ink draws str onto dst at (x, y) using face and color.
func Ink(dst *Surface, face Face, color RGBA, x, y int, str string) {
	if face == nil {
		face = BuiltinFace
	}
	opts := text.DrawOptions{}
	opts.LineSpacing = float64(Stride(face))
	opts.GeoM.Translate(float64(x), float64(y))
	opts.ColorScale.Scale(
		float32(color.R)/255,
		float32(color.G)/255,
		float32(color.B)/255,
		float32(color.A)/255,
	)
	text.Draw(dst, str, face, &opts)
}

// Measure returns the width and height of str when drawn with face at the
// given stride between the lines.
func Measure(str string, face Face, stride float64) (width, height float64) {
	return text.Measure(str, face, stride)
}

// Stride returns the recommended line height in pixels for face.
func Stride(face Face) int {
	m := face.Metrics()
	return int(m.HAscent + m.HDescent + m.HLineGap)
}
