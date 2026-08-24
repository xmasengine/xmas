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
	InputUp
	InputRight
	InputDown
	InputLeft
	InputAction
)

type Output func(screen *xgal.Surface)
type Runner func(input chan<- Input, output <-chan Output)

type Thread struct {
	input  chan<- Input
	output <-chan Output
	run    Runner
}

func (t *Thread) UpdateMode() (done bool, next Mode) {
	pressed := xgal.Keys()
	for _, k := range pressed {
		switch k {
		case xgal.KeyArrowUp:
			t.input <- InputUp
		case xgal.KeyArrowDown:
			t.input <- InputDown
		case xgal.KeyArrowLeft:
			t.input <- InputLeft
		case xgal.KeyArrowRight:
			t.input <- InputRight
		case xgal.KeyEnter:
			t.input <- InputAction
		}
	}
	return false, nil
}

func (t *Thread) DrawMode(screen *xgal.Surface) {
	for o := range t.output {
		if o == nil {
			break
		}
		o(screen)
	}
}
