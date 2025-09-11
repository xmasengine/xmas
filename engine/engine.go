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
	box1 := engine.Root.AddBox(image.Rect(20, 30, 200, 150))
	box2 := engine.Root.AddBox(image.Rect(210, 40, 410, 160))
	lab1 := box1.AddLabel(image.Rect(25, 30, 125, 47), "Label")
	box1.AddButton(image.Rect(25, 130, 125, 147), "Button", func(b *xui.Button) { lab1.SetText("Click!"); println("button clicked") })
	box2.AddCheckbox(image.Rect(220, 50, 380, 70), "Check", func(b *xui.Checkbox) { lab1.SetText("Check!"); println("checkbox clicked") })
	box2.AddEntry(image.Rect(220, 90, 380, 120), "Entry", func(b *xui.Entry) { lab1.SetText(b.Text()); println("entry changed") })
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
