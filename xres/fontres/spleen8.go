// fontres contains several tiny fonts embedded in go.
package fontres

import _ "embed"

import "github.com/xmasengine/xmas/xgal"

//go:embed spleen8.bdf
var spleen8Buffer []byte

//go:embed f8x13.bdf
var f8x13Buffer []byte

//go:embed f6x10.bdf
var f6x10Buffer []byte

func Load(buf []byte) (xgal.Face, error) {
	font, err := xgal.EmbeddedFont(buf, 8, xgal.BDF)
	if err != nil {
		return nil, err
	}
	return font, nil
}

var TinyFace xgal.Face
var SmallFace xgal.Face
var MediumFace xgal.Face

func init() {
	var err error
	TinyFace, err = Load(spleen8Buffer)
	if err != nil {
		panic(err)
	}
	SmallFace, err = Load(f6x10Buffer)
	if err != nil {
		panic(err)
	}
	MediumFace, err = Load(f8x13Buffer)
	if err != nil {
		panic(err)
	}
}
