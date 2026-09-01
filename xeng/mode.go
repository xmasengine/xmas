package xeng

/*
import (
	"fmt"
	"image"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"
)

import (
	"github.com/xmasengine/xmas/xdat"

	"github.com/xmasengine/xmas/xlog"
	// "github.com/xmasengine/xmas/xres"
	"github.com/xmasengine/xmas/xlui"
	"github.com/xmasengine/xmas/xzed"
)
*/

import (
	"github.com/xmasengine/xmas/xgal"
)

type Mode interface {
	UpdateMode() (done bool, next Mode)
	DrawMode(screen *xgal.Surface)
}

type Input int

const (
	InputNone Input = iota
	InputStop
	InputUp
	InputRight
	InputDown
	InputLeft
	InputAction
)

func GetInput() []Input {
	input := []Input{}
	pressed := xgal.Keys()
	for _, k := range pressed {
		switch k {
		case xgal.KeyArrowUp:
			input = append(input, InputUp)
		case xgal.KeyArrowDown:
			input = append(input, InputDown)
		case xgal.KeyArrowLeft:
			input = append(input, InputLeft)
		case xgal.KeyArrowRight:
			input = append(input, InputRight)
		case xgal.KeyEnter:
			input = append(input, InputAction)
		}
	}
	return input
}
