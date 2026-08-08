// spleen8 is a tiny font embedded in go.
package spleen8

import _ "embed"

import "github.com/xmasengine/xmas/xgal"

//go:embed spleen8.bdf
var fontBuffer []byte

func Load() (xgal.Face, error) {
	font, err := xgal.EmbeddedFont(fontBuffer, 8, xgal.BDF)
	if err != nil {
		return nil, err
	}
	return font, nil
}

var Face xgal.Face

func init() {
	var err error
	Face, err = Load()
	if err != nil {
		panic(err)
	}
}
