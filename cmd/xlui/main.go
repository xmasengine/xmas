package main

import (
	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xlui"
)

const (
	WindowW     = 240
	WindowH     = 192
	WindowScale = 3
)

type App struct {
	xlui.UI
}

func (a *App) Update() error {
	a.UI.Poll()
	return nil
}

func (a *App) Draw(screen *xgal.Surface) {
	xgal.Clear(screen, xgal.Paint(40, 80, 160, 255))
	a.UI.Render(screen)
}

func (a *App) Layout(w, h int) (int, int) {
	return WindowW, WindowH
}

var _ xgal.Game = (*App)(nil)

func main() {
	app := &App{}
	layer := app.Layer(xgal.Bound(10, 10, WindowW-10*2, 32))
	layer.Label("hello")
	layer.Button("OK")
	layer.Orientation = xlui.Vertical // Set to vertical.
	done := layer.Button("Done")      // will go below.
	done.Class.Click = func(at xgal.Point, button int) xlui.Reply {
		println("Click main done, finish", button)
		return xlui.Finish
	}
	xgal.Screen(WindowW*WindowScale, WindowH*WindowScale, "xpix")
	xgal.Play(app)
}
