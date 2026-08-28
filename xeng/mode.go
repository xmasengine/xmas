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

type Output func(screen *xgal.Surface) (done bool, next Mode)
type Runner func(output chan<- Output, input <-chan Input)

type Thread struct {
	input  chan Input
	output chan Output
	run    Runner
	done   bool
	next   Mode
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
	return t.done, t.next
}

func (t *Thread) DrawMode(screen *xgal.Surface) {
	for o := range t.output {
		if o == nil {
			break
		}
		d, n := o(screen)
		if d {
			t.done = d
			t.next = n
		}
	}
}

var _ Mode = &Thread{}

func NewThread(run Runner) *Thread {
	t := &Thread{}
	t.run = run
	return t
}

// Start a thread. No effect if already started.
func (t *Thread) Start() {
	if t.input == nil && t.output == nil {
		t.input = make(chan Input)
		t.output = make(chan Output)
		go t.run(t.output, t.input)
	}
}

// Stop a thread. No effect if already stopped or not started.
func (t *Thread) Stop() {
	if t.input != nil && t.output != nil {
		t.input <- InputStop
		close(t.input)
		close(t.output)
		t.input = nil
		t.output = nil
		t.done = true
	}
}
