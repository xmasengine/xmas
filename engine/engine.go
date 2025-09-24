package engine

import (
	"fmt"
	"image"
	"log/slog"
	"strings"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

import (
	"github.com/xmasengine/xmas/xmap"
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
	Zone       *xmap.Zone
	Debug      bool
}

func New(sw, sh int) *Engine {
	engine := &Engine{ScreenSize: image.Point{X: sw, Y: sh}, Msg: "!"}
	engine.At = image.Rect(0, 0, ViewWidth, ViewHeight)
	engine.Pressed = make([]ebiten.Key, 16)
	engine.Root = xui.NewRoot()
	engine.testZone()
	return engine
}

func (engine *Engine) testZone() {
	zone := xmap.NewZone("forest", 64, 64)
	layer := &zone.Layers[0]
	err := layer.LoadSource("pack/image/gfx/overworld.png")
	if err != nil {
		slog.Error("LoadSource", "file", "pack/image/gfx/overworld.png")
	}
	layer.FillIndex(image.Rect(0, 0, 63, 63), 0)
	engine.Zone = zone
}

func (engine *Engine) testUI() {
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
	// box1.AddSlider(image.Rect(130, 40, 140, 140), nil, func(s *xui.Slider) { lab1.SetText("Slide!"); println("slider clicked", s.Pos) })
	box1.AddVerticalScroller(func(s *xui.Slider) { lab1.SetText("vScroll!"); println("vscroll clicked", s.Pos) })

	box2 := engine.Root.AddBox(image.Rect(210, 40, 430, 170))
	box2.AddCheckbox(image.Rect(220, 50, 380, 70), "Check", func(b *xui.Checkbox) { lab1.SetText("Check!"); println("checkbox clicked") })
	chooser := box2.AddChooser(image.Rect(220, 70, 380, 120), img, image.Pt(16, 16), func(c *xui.Chooser) {
		lab1.SetText("Chooser!")
		atx := c.Selected.Bounds.Min.X
		aty := c.Selected.Bounds.Min.Y
		println("chooser clicked", atx, aty)
	})
	vs := chooser.AddVerticalScroller(func(s *xui.Slider) {
		lab1.SetText("cvScroll!")

		// XXX this makes the image in the frame scroll
		// but it should be handled in xui.
		widget := &chooser.Frame
		scrollRange := widget.Extent.Dy() - xui.WidgetScrollSlack
		var noff xui.Point
		noff.Y = ((s.Pos - s.Low) * scrollRange) / (s.High - s.Low)
		widget.Offset = noff
		println("chooser vscroll clicked", s.Pos, noff.Y)
	})
	vs.Layer = chooser.Layer + 100
	box2.AddEntry(image.Rect(220, 130, 380, 150), "Entry", func(b *xui.Entry) { lab1.SetText(b.Text()); println("entry changed") })
	// box2.AddSlider(image.Rect(220, 155, 380, 165), nil, func(s *xui.Slider) { lab1.SetText("hSlide!"); println("hslider clicked", s.Pos) })
	box2.AddHorizontalScroller(func(s *xui.Slider) { lab1.SetText("hScroll!"); println("hscroll clicked", s.Pos) })
}

func (g *Engine) Update() error {
	if g.Root != nil {
		g.Root.Update()
	}

	g.Pressed = g.Pressed[:0]
	g.Pressed = inpututil.AppendPressedKeys(g.Pressed)
	var delta image.Point
	for _, k := range g.Pressed {
		switch k {
		case ebiten.KeyUp:
			delta.Y = -1
		case ebiten.KeyDown:
			delta.Y = 1
		case ebiten.KeyLeft:
			delta.X = -1
		case ebiten.KeyRight:
			delta.X = 1
		case ebiten.KeyF:
			g.Debug = !g.Debug
		default:
		}
	}

	if g.Zone != nil {
		g.Zone.Camera = g.Zone.Camera.Add(delta)
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
	if g.Root != nil {
		g.Root.Draw(screen)
	}
	if g.Zone != nil {
		g.Zone.Draw(screen)
	}
	if g.Debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
	}
}

func (g *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ViewWidth, ViewHeight
}
