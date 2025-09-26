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
	"github.com/xmasengine/xmas/xlog"
	"github.com/xmasengine/xmas/xmap"
	"github.com/xmasengine/xmas/xres"
	"github.com/xmasengine/xmas/xui"
)

// const ViewWidth = 320 // 2

const ViewWidth = 426  // *2
const ViewHeight = 240 // 2

// const ViewHeight = 240 * 2

type Engine struct {
	Log        xlog.Log
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
	engine.Log.Hide = true
	engine.testZone()
	engine.testUI()
	return engine
}

func dprintln(msg string, vars ...any) {
	slog.Info(msg, "vars", vars)
}

func (engine *Engine) testZone() {
	zone := xmap.NewZone("forest", 64, 64)
	layer := &zone.Layers[0]
	err := layer.LoadSource("pack/image/gfx/overworld.png")
	if err != nil {
		slog.Error("LoadSource", "file", "pack/image/gfx/overworld.png")
	}
	layer.FillIndex(image.Rect(0, 0, 63, 63), 0)
	player := &zone.Player
	err = player.LoadSource("pack/image/gfx/character.png")
	if err != nil {
		slog.Error("LoadSource", "file", "pack/image/gfx/character.png")
	}
	player.Tw = 16
	player.Th = 32
	player.At = image.Pt(160, 160)

	player.AddNewPoses(xmap.Stand, 0, 0, 16, 32, 1)
	player.AddNewPoses(xmap.Walk, 16, 0, 16, 32, 3)
	zone.Player = *player

	engine.Zone = zone
}

func (engine *Engine) testUI() {
	img, ierr := xres.LoadImageFromFile("pack/tile/tile_0001.png")
	if ierr != nil {
		slog.Error("LoadImageFromFile", "file", "pack/tile/tile_0001.png")
	}
	box1 := engine.Root.AddBox(image.Rect(20, 30, 200, 150))
	lab1 := box1.AddLabel(image.Rect(25, 100, 125, 120), "Label")
	lab1.AddTitleBar(10, "Drag Me")

	bar1 := box1.AddBar(image.Rect(25, 35, 125, 50), func(b *xui.Bar) { lab1.SetText("Bar!"); dprintln("bar clicked") })
	_ = bar1
	hello := bar1.FitItemWithMenu("hello", func(b *xui.Item) { lab1.SetText("hello"); dprintln("bar item hello clicked") })
	menu := hello.Menu
	menu.FitItem("sub1", func(b *xui.Item) { lab1.SetText("sub1"); dprintln("bar item hello > sub1 clicked") })
	sub2 := menu.FitItemWithMenu("sub2", func(b *xui.Item) { lab1.SetText("sub2"); dprintln("bar item hello > sub2 clicked") })

	subMenu := sub2.Menu
	subMenu.FitItem("subsub1", func(b *xui.Item) { lab1.SetText("subsub1"); dprintln("bar item hello > subsub1 clicked") })
	subMenu.FitItem("subsub2", func(b *xui.Item) { lab1.SetText("subsub2"); dprintln("bar item hello > subsub2 clicked") })

	bar1.FitItem("world", func(b *xui.Item) { lab1.SetText("world"); dprintln("bar item world clicked") })
	box1.AddButton(image.Rect(25, 130, 125, 147), "Button", func(b *xui.Button) { lab1.SetText("Click!"); dprintln("button clicked") })
	// box1.AddSlider(image.Rect(130, 40, 140, 140), nil, func(s *xui.Slider) { lab1.SetText("Slide!"); dprintln("slider clicked", s.Pos) })
	box1.AddVerticalScroller(func(s *xui.Slider) { lab1.SetText("vScroll!"); dprintln("vscroll clicked", s.Pos) })

	box2 := engine.Root.AddBox(image.Rect(210, 40, 430, 170))
	box2.AddCheckbox(image.Rect(220, 50, 380, 70), "Check", func(b *xui.Checkbox) { lab1.SetText("Check!"); dprintln("checkbox clicked") })
	chooser := box2.AddChooser(image.Rect(220, 70, 380, 120), img, image.Pt(16, 16), func(c *xui.Chooser) {
		lab1.SetText("Chooser!")
		atx := c.Selected.Bounds.Min.X
		aty := c.Selected.Bounds.Min.Y
		dprintln("chooser clicked", atx, aty)
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
		dprintln("chooser vscroll clicked", s.Pos, noff.Y)
	})
	vs.Layer = chooser.Layer + 100
	box2.AddEntry(image.Rect(220, 130, 380, 150), "Entry", func(b *xui.Entry) { lab1.SetText(b.Text()); dprintln("entry changed") })
	// box2.AddSlider(image.Rect(220, 155, 380, 165), nil, func(s *xui.Slider) { lab1.SetText("hSlide!"); dprintln("hslider clicked", s.Pos) })
	box2.AddHorizontalScroller(func(s *xui.Slider) { lab1.SetText("hScroll!"); dprintln("hscroll clicked", s.Pos) })
	// Add a title bar for dragging
	box2.AddTitleBar(10, "Box 2")
}

func (g *Engine) Update() error {
	g.Log.Update()

	if g.Root != nil {
		g.Root.Update()
	}

	g.Pressed = g.Pressed[:0]
	g.Pressed = inpututil.AppendPressedKeys(g.Pressed)
	var delta image.Point
	var mdelta image.Point
	act := xmap.Stand
	var dir xmap.Direction
	if g.Zone != nil {
		dir = g.Zone.Player.Direction
	}
	for _, k := range g.Pressed {
		switch k {
		case ebiten.KeyUp:
			delta.Y = -1
			dir = xmap.North
			act = xmap.Walk
		case ebiten.KeyDown:
			delta.Y = 1
			dir = xmap.South
			act = xmap.Walk
		case ebiten.KeyLeft:
			delta.X = -1
			dir = xmap.West
			act = xmap.Walk
		case ebiten.KeyRight:
			delta.X = 1
			dir = xmap.East
			act = xmap.Walk
		case ebiten.KeyPageUp:
			mdelta.Y = -1
		case ebiten.KeyPageDown:
			mdelta.Y = 1
		case ebiten.KeyHome:
			mdelta.X = -1
		case ebiten.KeyEnd:
			mdelta.X = 1
		case ebiten.KeyF:
			g.Debug = !g.Debug
		default:
		}
	}

	if g.Zone != nil {
		g.Zone.Camera = g.Zone.Camera.Add(mdelta)
		g.Zone.Player.At = g.Zone.Player.At.Add(delta)
		pose := g.Zone.Player.BestPose(dir, act)
		g.Zone.Player.Pose = pose
		g.Zone.Player.Update()
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
	if g.Zone != nil {
		g.Zone.Draw(screen)
		if g.Debug {
			pose := g.Zone.Player.Pose
			ebitenutil.DebugPrint(screen, fmt.Sprintf("pose: %d %d %d %d %d",
				pose.Direction, pose.Action, pose.Phase, pose.Frames, pose.Tick))
		}
	}
	if g.Root != nil {
		g.Root.Draw(screen)
	}
	if g.Debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n%f\n", ebiten.ActualFPS()))
	}
	g.Log.Draw(screen)
}

func (g *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Log.Layout(ViewWidth, ViewHeight)
	return ViewWidth, ViewHeight
}
