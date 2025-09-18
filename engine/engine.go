package engine

import (
	"image"
	"log/slog"
	"strings"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

import (
	"github.com/xmasengine/xmas/xres"
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
	img, ierr := xres.LoadImageFromFile("pack/tile/tile_0001.png")
	if ierr != nil {
		slog.Error("LoadImageFromFile", "file", "pack/tile/tile_0001.png")
	}
	box1 := engine.Root.AddBox(image.Rect(20, 30, 200, 150))
	lab1 := box1.AddLabel(image.Rect(25, 100, 125, 120), "Label")

	bar1 := box1.AddBar(image.Rect(25, 35, 125, 50), func(b *xui.Bar) { lab1.SetText("Bar!"); println("bar clicked") })
	_ = bar1
	hello := bar1.FitItemWithMenu("hello", func(b *xui.Item) { lab1.SetText("hello"); println("bar item hello clicked") })
	menu := hello.Menu
	menu.FitItem("sub1", func(b *xui.Item) { lab1.SetText("sub1"); println("bar item hello > sub1 clicked") })
	sub2 := menu.FitItemWithMenu("sub2", func(b *xui.Item) { lab1.SetText("sub2"); println("bar item hello > sub2 clicked") })

	subMenu := sub2.Menu
	subMenu.FitItem("subsub1", func(b *xui.Item) { lab1.SetText("subsub1"); println("bar item hello > subsub1 clicked") })
	subMenu.FitItem("subsub2", func(b *xui.Item) { lab1.SetText("subsub2"); println("bar item hello > subsub2 clicked") })

	bar1.FitItem("world", func(b *xui.Item) { lab1.SetText("world"); println("bar item world clicked") })
	box1.AddButton(image.Rect(25, 130, 125, 147), "Button", func(b *xui.Button) { lab1.SetText("Click!"); println("button clicked") })
	box1.AddSlider(image.Rect(130, 40, 140, 140), nil, func(s *xui.Slider) { lab1.SetText("Slide!"); println("slider clicked", s.Pos) })

	box2 := engine.Root.AddBox(image.Rect(210, 40, 430, 170))
	box2.AddCheckbox(image.Rect(220, 50, 380, 70), "Check", func(b *xui.Checkbox) { lab1.SetText("Check!"); println("checkbox clicked") })
	box2.AddChooser(image.Rect(220, 70, 380, 120), img, image.Pt(16, 16), func(c *xui.Chooser) {
		lab1.SetText("Chooser!")
		atx := c.Selected.Bounds.Min.X
		aty := c.Selected.Bounds.Min.Y
		println("chooser clicked", atx, aty)
	})
	box2.AddEntry(image.Rect(220, 130, 380, 150), "Entry", func(b *xui.Entry) { lab1.SetText(b.Text()); println("entry changed") })
	box2.AddSlider(image.Rect(220, 155, 380, 165), nil, func(s *xui.Slider) { lab1.SetText("hSlide!"); println("hslider clicked", s.Pos) })
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
