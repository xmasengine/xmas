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
	return nil
}

func (a *App) Draw(screen *xgal.Surface) {
	xgal.Clear(screen, xgal.Paint(40, 80, 160, 255))
}

func (a *App) Layout(w, h int) (int, int) {
	return WindowW, WindowH
}

var _ xgal.Game = (*App)(nil)

func main() {
	app := &App{}
	xgal.Screen(WindowW*WindowScale, WindowH*WindowScale, "xpix")
	xgal.Play(app)
}
