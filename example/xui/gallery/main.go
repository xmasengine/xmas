package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xui"
)

type gallery struct {
	root       *xui.Layer
	pane       *xui.PaneLayer
	clickCount int
	btn        *xui.ButtonLayer
	entry      *xui.EntryLayer
	check      *xui.CheckboxLayer
	slider     *xui.SliderLayer
	list       *xui.ListLayer
	frameImg   *xgal.Surface
	tileImg    *xgal.Surface
	chooser    *xui.ChooserLayer
	tip        *xui.TooltipLayer
	status     *xui.LabelLayer
}

func (g *gallery) Update() error {
	g.root.Poll()

	// update status
	g.status.SetText(fmt.Sprintf("Count:%d  Check:%v  Slider:%d  Entry:%q",
		g.clickCount, g.check.Checked, g.slider.Pos, g.entry.Text()))

	// tooltip follows button bounds
	if g.btn != nil {
		g.tip.Trigger = g.btn.Bounds
	}

	return nil
}

func (g *gallery) Draw(screen *xgal.Surface) {
	xgal.Clear(screen, xgal.Wash(40, 40, 40, 255))

	// draw a tiled background grid
	for x := 0; x < 800; x += 16 {
		xgal.Line(screen, x, 0, x, 600, 1, xgal.Wash(50, 50, 50, 255))
	}
	for y := 0; y < 600; y += 16 {
		xgal.Line(screen, 0, y, 800, y, 1, xgal.Wash(50, 50, 50, 255))
	}

	g.root.Render(screen)
}

func (g *gallery) Layout(w, h int) (int, int) {
	g.root.Place(xgal.Rect(0, 0, w, h))
	return w, h
}

func main() {
	g := &gallery{}

	// ── self-drawn 16×16 icons ──
	iconStar := xgal.Prepare(16, 16)
	xgal.Box(iconStar, xgal.Rect(0, 0, 16, 16), xgal.Wash(0, 0, 0, 0))
	xgal.Ink(iconStar, xgal.BuiltinFace, xgal.Wash(255, 220, 0, 255), 2, 0, "*")

	iconDoc := xgal.Prepare(16, 16)
	xgal.Box(iconDoc, xgal.Rect(0, 0, 16, 16), xgal.Wash(0, 0, 0, 0))
	xgal.Box(iconDoc, xgal.Rect(2, 1, 14, 15), xgal.Wash(220, 220, 200, 255))
	xgal.Outline(iconDoc, xgal.Rect(2, 1, 14, 15), 1, xgal.Wash(100, 100, 80, 255))
	xgal.Line(iconDoc, 5, 5, 11, 5, 1, xgal.Wash(80, 80, 60, 255))
	xgal.Line(iconDoc, 5, 8, 11, 8, 1, xgal.Wash(80, 80, 60, 255))
	xgal.Line(iconDoc, 5, 11, 9, 11, 1, xgal.Wash(80, 80, 60, 255))

	iconHelp := xgal.Prepare(16, 16)
	xgal.Box(iconHelp, xgal.Rect(0, 0, 16, 16), xgal.Wash(0, 0, 0, 0))
	xgal.Circle(iconHelp, xgal.Pt(8, 7), 6, 1, xgal.Wash(200, 200, 100, 255))
	xgal.Ink(iconHelp, xgal.BuiltinFace, xgal.Wash(200, 200, 100, 255), 6, 3, "?")

	iconHand := xgal.Prepare(16, 16)
	xgal.Box(iconHand, xgal.Rect(0, 0, 16, 16), xgal.Wash(0, 0, 0, 0))
	xgal.Disk(iconHand, xgal.Pt(8, 8), 6, xgal.Wash(200, 200, 220, 255))
	xgal.Circle(iconHand, xgal.Pt(8, 8), 6, 1, xgal.Wash(100, 100, 120, 255))
	for _, p := range []xgal.Point{xgal.Pt(6, 5), xgal.Pt(7, 5), xgal.Pt(9, 5), xgal.Pt(10, 5), xgal.Pt(7, 11), xgal.Pt(8, 11), xgal.Pt(9, 11)} {
		xgal.Box(iconHand, xgal.Rect(p.X, p.Y, p.X+1, p.Y+1), xgal.Wash(100, 100, 120, 255))
	}

	g.root = &xui.Layer{}
	g.root.Style = xui.DefaultStyle()
	g.root.Axis = xui.Vertical

	// ── Title ──
	title := g.root.AddLabel(xgal.Rect(0, 0, 200, 22), "xui Gallery")
	title.Icon = xui.Icon{Image: iconStar}

	// ── Pane (draggable, scrollable, with menu) ──
	g.pane = xui.Pane(xgal.Rect(20, 30, 500, 420), "Demo Pane")
	g.pane.Style.Fill = xgal.Wash(20, 20, 80, 220)

	// menu bar
	bar := g.pane.AddMenuBar()
	fileMenu := bar.FitItem("File", nil)
	fileMenu.Icon = xui.Icon{Image: iconDoc}
	fileMenu.AddDropdown().FitItem("New", func() {
		g.root.AddAsk("Start over?", "OK")
	})
	fileMenu.AddDropdown().FitItem("Quit", func() {
		os.Exit(0)
	})
	bar.FitItem("Help", func() {
		g.root.AddAsk("xui Gallery — interactive widget demo", "OK")
	}).Icon = xui.Icon{Image: iconHelp}

	// ── content inside pane ──

	// Button
	g.btn = g.pane.AddButton(xgal.Rect(10, 10, 130, 24), "Click me (0)", func() {
		g.clickCount++
		g.btn.Text = "Click me (" + strconv.Itoa(g.clickCount) + ")"
	})
	g.btn.Icon = xui.Icon{Image: iconHand}

	// Entry
	g.entry = g.pane.AddEntry(xgal.Rect(10, 40, 200, 22), "", func(val string) {
		println("entry:", val)
	})

	// Checkbox
	g.check = g.pane.AddCheckbox(xgal.Rect(10, 70, 150, 18), "Enable option", nil)

	// Slider
	g.slider = g.pane.AddSlider(xgal.Rect(10, 100, 200, 16), func(pos int) {
		println("slider:", pos)
	})

	// List
	g.list = g.pane.AddList(xgal.Rect(10, 130, 200, 80))
	g.list.AddItem("Apple")
	g.list.AddItem("Banana")
	g.list.AddItem("Cherry")
	g.list.OnSelect = func(i int) {
		println("selected:", g.list.Items[i])
	}

	// Frame (hand-drawn image)
	g.frameImg = xgal.Prepare(48, 48)
	xgal.Box(g.frameImg, xgal.Rect(0, 0, 48, 48), xgal.Wash(0, 80, 0, 255))
	xgal.Disk(g.frameImg, xgal.Pt(24, 24), 16, xgal.Wash(255, 200, 0, 255))
	xgal.Circle(g.frameImg, xgal.Pt(24, 24), 16, 2, xgal.Wash(200, 100, 0, 255))
	xgal.Line(g.frameImg, 12, 32, 36, 32, 3, xgal.Wash(200, 50, 50, 255))
	frameLayer := g.pane.AddFrame(xgal.Rect(220, 10, 280, 58), g.frameImg)
	frameLayer.Style.Border = xgal.Wash(0, 200, 200, 200)

	// Chooser (tilesheet with colored squares)
	g.tileImg = xgal.Prepare(64, 64)
	colors := []xgal.RGBA{
		xgal.Wash(255, 0, 0, 255), xgal.Wash(0, 255, 0, 255), xgal.Wash(0, 0, 255, 255),
		xgal.Wash(255, 255, 0, 255), xgal.Wash(255, 0, 255, 255), xgal.Wash(0, 255, 255, 255),
		xgal.Wash(200, 100, 0, 255), xgal.Wash(100, 200, 0, 255), xgal.Wash(200, 200, 200, 255),
		xgal.Wash(255, 128, 0, 255), xgal.Wash(128, 0, 255, 255), xgal.Wash(0, 200, 128, 255),
		xgal.Wash(255, 200, 200, 255), xgal.Wash(200, 255, 200, 255), xgal.Wash(200, 200, 255, 255),
		xgal.Wash(255, 255, 255, 255),
	}
	for i, c := range colors {
		tx := (i % 4) * 16
		ty := (i / 4) * 16
		xgal.Box(g.tileImg, xgal.Rect(tx, ty, tx+16, ty+16), c)
		xgal.Outline(g.tileImg, xgal.Rect(tx, ty, tx+16, ty+16), 1, xgal.Wash(0, 0, 0, 255))
	}
	g.chooser = g.pane.AddChooser(xgal.Rect(220, 60, 284, 124), g.tileImg,
		xgal.Pt(16, 16), func(x, y int) {
			println("tile:", x, y)
		})

	// add pane to root
	g.root.Add(g.pane)

	// ── Status bar ──
	g.status = g.root.AddLabel(xgal.Rect(0, 460, 500, 18), "Ready")

	// ── Tooltip ──
	g.tip = xui.Tooltip(xgal.Rect(0, 0, 0, 0), "Click the button!")
	g.root.Add(g.tip)

	// ── Run ──
	xgal.Screen(800, 600, "xui Gallery")
	xgal.Decorate(true)
	xgal.Stretch(true)
	if err := xgal.Play(g); err != nil {
		fmt.Fprintf(os.Stderr, "RUNTIME: %v\n", err)
		os.Exit(2)
	}
}
