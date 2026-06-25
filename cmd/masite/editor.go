package main

import (
	"fmt"
)

import (
	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xui"
)

type Editor struct {
	Name   string
	Map    *Map
	Camera xgal.Rectangle
	Hover  xgal.Point
	Tile   xgal.Point // Tile we are hovering
	Cell   Cell
	Scale  int
	Error  error
	Midget xui.Layer // Child widget layers
}

func (e Editor) Draw(screen *xgal.Surface) {
	if e.Map != nil {
		e.Map.Render(screen, e.Camera)
	}
	if e.Error != nil {
		xgal.Debug(screen, fmt.Sprintf("Error: %s", e.Error),
			e.Map.Width*e.Map.Tw, 10,
		)
	} else {
		xgal.Debug(screen, fmt.Sprintf("%s: (%d,%d): %d %d",
			e.Name, e.Hover.X, e.Hover.Y, e.Cell.Index, e.Cell.Flag,
		), e.Map.Width*e.Map.Tw, 10)
	}
	e.Midget.Render(screen)
}

func (e Editor) Layout(w, h int) (rw, th int) {
	e.Midget.Place(w, h)
	return e.Camera.Dx() / e.Scale, e.Camera.Dy() / e.Scale
}

const HELP = `HELP:
	Pause: Exit witout saving.
	F1: This help.
	F2: Save map in mashite format.
	Mouse Wheel: Select tile index.
	Esc: Cancel dialogs.
`

func (e *Editor) Update() error {
	var err error
	e.Hover = xgal.Mouse()
	e.Tile = e.Map.ToTile(e.Hover, e.Camera)

	_, wheel := xgal.Wheel()
	if wheel > 0 {
		e.Cell.Index++
	} else if wheel < 0 {
		e.Cell.Index = max(0, e.Cell.Index-1)
	}

	res := e.Midget.Poll()
	if res != xui.Accept {
		switch {
		case xgal.Tap(xgal.KeyPause):
			if len(e.Midget.Kids) < 1 {
				e.Midget.YesNo(50, 50, 250, 100, "Quit", "Y",
					func(resp bool) {
						e.Midget.Done = resp
					},
				)
			}
		case xgal.Tap(xgal.KeyH):
			e.Cell.Flag ^= FlagHorizontalFlip
		case xgal.Tap(xgal.KeyV):
			e.Cell.Flag ^= FlagVerticalFlip
		case xgal.Tap(xgal.KeyF10):
			e.Error = nil
		case xgal.Tap(xgal.KeyF1):
			e.Midget.Ask(100, 0, 250, 200, HELP, "", func(name string) {})
		case xgal.Tap(xgal.KeyF2):
			e.Midget.Ask(50, 50, 250, 100, "Save As", e.Name,
				func(name string) {
					e.Error = e.Map.Save(name)
					if e.Error == nil {
						e.Name = name
					}
				},
			)
		case xgal.Click(xgal.MouseButtonLeft):
			e.Map.Put(e.Tile, e.Cell)
		default:
		}
	}

	if e.Midget.Done {
		return Termination
	}
	return err
}

func NewEditor(tm *Map, name string, w, h, scale int) *Editor {
	return &Editor{Map: tm, Name: name, Camera: xgal.Rect(0, 0, w, h),
		Scale:  scale,
		Midget: xui.MakeLayer(xgal.Rect(0, 0, 0, 0)),
	}
}
