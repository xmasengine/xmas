package engine

import (
	"image"
	"strings"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
}

func New(sw, sh int) *Engine {
	Engine := &Engine{ScreenSize: image.Point{X: sw, Y: sh}, Msg: "!"}
	Engine.At = image.Rect(0, 0, ViewWidth, ViewHeight)
	Engine.Pressed = make([]ebiten.Key, 16)
	return Engine
}

func (g *Engine) Update() error {
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
}

func (g *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ViewWidth, ViewHeight
}
