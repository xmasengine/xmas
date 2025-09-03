package engine

import (
	"image"
	"strings"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

import (
	"github.com/xmasengine/xmas/xui"
)

// const ViewWidth = 320 // 2

const ViewWidth = 426  // *2
const ViewHeight = 240 // 2

// const ViewHeight = 240 * 2

type Engine struct {
	Msg        string
	Pressed    []ebiten.Key
	Script     strings.Reader
	DebugRow   int
	ScreenSize image.Point
	At         image.Rectangle
	Root       *xui.Root
}

func New(sw, sh int) *Engine {
	engine := &Engine{ScreenSize: image.Point{X: sw, Y: sh}, Msg: "!"}
	engine.At = image.Rect(0, 0, ViewWidth, ViewHeight)
	engine.Pressed = make([]ebiten.Key, 16)
	engine.Root = xui.NewRoot()
	box1 := xui.NewBox(image.Rect(20, 30, 200, 150))
	box2 := xui.NewBox(image.Rect(70, 90, 150, 100))
	box1.AddButton(image.Rect(25, 130, 125, 147), "Button")
	engine.Root.Append(box1, box2)
	return engine
}

func (g *Engine) Update() error {
	if g.Root != nil {
		g.Root.Update()
	}

	g.Pressed = g.Pressed[:0]
	g.Pressed = inpututil.AppendPressedKeys(g.Pressed)
	for _, k := range g.Pressed {
		switch k {
		case ebiten.KeyUp:
		case ebiten.KeyDown:
		case ebiten.KeyLeft:
		case ebiten.KeyRight:
		default:
		}
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return ebiten.Termination
	default:
	}

	return nil
}

const tileDebug = false

func (g *Engine) Draw(screen *ebiten.Image) {
	g.Root.Draw(screen)
}

func (g *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ViewWidth, ViewHeight
}
