// spleen8 is a tiny font embedded in go.
package spleen8

import _ "embed"

import "github.com/zachomedia/go-bdf"
import "golang.org/x/image/font"
import "github.com/hajimehoshi/ebiten/v2/text/v2"

//go:embed spleen8.bdf
var fontBuffer []byte

func Load() (font.Face, error) {
	font, err := bdf.Parse(fontBuffer)
	if err != nil {
		return nil, err
	}
	return font.NewFace(), nil
}

var Face font.Face
var XFace text.Face

func init() {
	var err error
	Face, err = Load()
	if err != nil {
		panic(err)
	}
	XFace = text.NewGoXFace(Face)
}
