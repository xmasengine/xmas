package main

import (
	"io/fs"
	"os"
)

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
	FS    fs.FS
	Image *xgal.Surface
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
	var err error
	app := &App{}
	wd, _ := os.Getwd()
	app.FS = os.DirFS(wd)

	layer := app.Layer(xgal.Bound(10, 10, WindowW-10*2, 40))
	layer.Label("hello")
	layer.Button("OK")
	layer.Orientation = xlui.Vertical // Set to vertical.
	done := layer.Button("Done")      // will go below.
	done.Class.Click = func(at xgal.Point, button int) xlui.Reply {
		println("Click main done, finish", button)
		return xlui.Finish
	}

	layer2 := app.Layer(xgal.Bound(15, 15, WindowW-10*2, 40))
	layer2.Label("hello 2")
	entry := layer2.Entry("")
	entry.Class.Entry = func(string) xlui.Reply {
		println("Accept: ", entry.Text)
		return xlui.Accept
	}
	entry2 := layer2.Entry("")
	entry2.Class.Entry = func(string) xlui.Reply {
		println("Accept: ", entry2.Text)
		return xlui.Accept
	}

	layer2.CheckboxWithLabel(false, "Check")
	layer2.Checkbox(true)
	g := &xlui.Group{}

	layer2.Orientation = xlui.Vertical
	layer2.Toggle("Foo", g)
	layer2.Orientation = xlui.Horizontal

	layer2.Toggle("Bar", g)
	layer2.Toggle("Baz", g)
	layer2.Toggle("Quux", nil)

	app.Image, err = xgal.Texture(app.FS, "pack/tile/tile_0002.png")
	if err != nil {
		layer3 := app.Layer(xgal.Bound(10, 10, WindowW-10*2, 40))
		layer3.Label("Error:" + err.Error())
	} else {
		layer3 := app.Layer(xgal.Bound(10, 50, app.Image.Bounds().Dx()+10, app.Image.Bounds().Dy()+10))
		tileSize := xgal.Pt(8, 8)
		layer3.Chooser(app.Image, tileSize)
	}

	layer4 := app.Layer(xgal.Bound(10, 50, WindowW-10*2, 40))
	slider := layer4.Slider(xlui.Vertical, 0, 10, 2)
	slider.Class.Slide = func(value int) xlui.Reply {
		println("Slider: ", value)
		return xlui.Accept
	}
	slider2 := layer4.Slider(xlui.Horizontal, 0, 10, 2)
	slider2.Class.Slide = func(value int) xlui.Reply {
		println("Slider: ", value)
		return xlui.Accept
	}

	xgal.Screen(WindowW*WindowScale, WindowH*WindowScale, "xpix")
	xgal.Play(app)
}
