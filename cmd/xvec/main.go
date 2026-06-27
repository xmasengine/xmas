package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"slices"

	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xui"
	"github.com/xmasengine/xmas/xvec"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

type Tool int

const (
	ToolCircle Tool = iota
	ToolDisk
	ToolRect
	ToolSlab
	ToolLine
	ToolStroke
	ToolFill
	toolCount
)

var toolNames = []string{"Circle", "Disk", "Rect", "Slab", "Line", "Stroke", "Fill"}
var toolFKeys = []xgal.KeyCode{xgal.KeyF2, xgal.KeyF3, xgal.KeyF4, xgal.KeyF5, xgal.KeyF6, xgal.KeyF7, xgal.KeyF8}
var toolDigits = []xgal.KeyCode{xgal.KeyDigit1, xgal.KeyDigit2, xgal.KeyDigit3, xgal.KeyDigit4, xgal.KeyDigit5, xgal.KeyDigit6, xgal.KeyDigit7}

func gencol(r, g, b int) xvec.Color {
	return xvec.Color{
		R: uint8(r * 85),
		G: uint8(g * 85),
		B: uint8(b * 85),
		A: 255,
	}
}

func genesisColors() [64]xvec.Color {
	var cols [64]xvec.Color
	for i := 0; i < 64; i++ {
		r := uint8((i / 16) * 85)
		g := uint8(((i % 16) / 4) * 85)
		b := uint8((i % 4) * 85)
		cols[i] = xvec.Color{R: r, G: g, B: b, A: 255}
	}
	return cols
}

// helpLines are the lines shown on the F1 overlay.
var helpLines = []struct {
	text string
	bold bool
}{
	{"xvec Editor Help", true},
	{"", false},
	{"F1            Toggle this help", false},
	{"F2–F7 / 1–6   Select tool", false},
	{"", false},
	{"Shapes:       First = start, second = size", false},
	{"Line:         Start → end point", false},
	{"Circle/Disk:  Center → radius point", false},
	{"Rect/Slab:    Corner → opposite corner", false},
	{"Stroke:       Click to add vertices, Close to finish", false},
	{"Fill:         Click to add vertices, Close to finish", false},
	{"", false},
	{"Palette:      Click to pick colour", false},
	{"Instr. list:  Click to select", false},
	{"Del/Backspace Delete selected shape", false},
	{"X             Clear all shapes", false},
	{"Ctrl+S        Save", false},
	{"Esc           Close help", false},
}

type App struct {
	doc      *xvec.XVEC
	docSurf  *xgal.Surface
	dirty    bool
	filename string
	showHelp bool

	tool    Tool
	color   xvec.Color
	selInst int

	palColors [64]xvec.Color
	palSel    int

	defSW float32

	pend *struct{ x, y float32 }

	list     *xui.ListLayer
	toolSel  int
	toggles  []*xui.ToggleLayer
	swSlider *xui.SliderLayer

	// Path editing
	pathSteps   []xvec.Stepper
	pathSegSel  int // 0=Move, 1=Line
	pathToggles []*xui.ToggleLayer

	msg      string
	msgTimer int
}

func main() {
	file := flag.String("f", "", "xvec file to edit")
	flag.Parse()

	a := &App{
		doc: &xvec.XVEC{
			Size:      xvec.Size{W: 320, H: 240},
			Antialias: true,
		},
		tool:     ToolCircle,
		selInst:  -1,
		filename: *file,
		defSW:    2,
	}

	// Load file if specified
	if a.filename != "" {
		f, err := os.Open(a.filename)
		if err == nil {
			a.doc = &xvec.XVEC{}
			if err := a.doc.Decode(f); err != nil {
				fmt.Fprintf(os.Stderr, "error loading %s: %v\n", a.filename, err)
				os.Exit(1)
			}
			f.Close()
		}
	}

	a.palColors = genesisColors()
	a.color = a.palColors[63]
	a.palSel = 63

	a.list = xui.List(a.listBounds())
	a.list.Selected = -1

	// Toolbar toggles
	btnW := windowWidth / int(toolCount)
	for i := range int(toolCount) {
		var t *xui.ToggleLayer
		toggled := func(active bool) {
			if active {
				if Tool(t.Idx) != a.tool {
					a.pathSteps = nil
					a.tool = Tool(t.Idx)
					a.pend = nil
				}
			}
		}
		t = xui.Toggle(xgal.Rect(i*btnW, 0, (i+1)*btnW, 28), toolNames[i], toggled)
		t.Style = xui.DefaultStyle()
		t.Group = &a.toolSel
		t.Idx = i
		t.Active = i == 0
		a.toggles = append(a.toggles, t)
	}

	// Path sub-toolbar toggles
	segNames := []string{"Move", "Line", "Close", "Done"}
	btnW = windowWidth / len(segNames)

	segFuncs := []func(bool){
		func(active bool) {},
		func(active bool) {},
		func(active bool) {
			if active {
				a.pathClose()
			}
		},
		func(active bool) {
			if active {
				a.toolSel = 0
				a.tool = Tool(a.toolSel)
				a.pend = nil
				a.pathSteps = nil
			}
		},
	}
	for i, name := range segNames {
		// Draw over the normal toggles.
		t := xui.Toggle(xgal.Rect(i*btnW, 0, (i+1)*btnW, 28), name, segFuncs[i])
		t.Style = xui.DefaultStyle()
		t.Group = &a.pathSegSel
		t.Idx = i
		t.Active = i == 1 // Line selected by default
		a.pathToggles = append(a.pathToggles, t)
	}
	a.pathSegSel = 1

	a.swSlider = xui.Slider(xgal.Rect(490, 358, 630, 380), func(pos int) {
		a.defSW = float32(pos)
		println("slider", pos, a.selInst)
		if a.selInst >= 0 && a.selInst < len(a.doc.Instructions) {
			inst := a.doc.Instructions[a.selInst]
			if adj, ok := inst.(xvec.Adjuster); ok {
				adj.Adjust(xvec.Length(pos))
			}
		}
	})
	a.swSlider.Low = 1
	a.swSlider.High = 20
	a.swSlider.Pos = 2

	title := "xvec editor"
	if a.filename != "" {
		title += " — " + a.filename
	}
	xgal.Cursor(true, xgal.Crosshair)
	xgal.Screen(windowWidth, windowHeight, title)
	xgal.Play(a)
}

func (a *App) toolbarBounds() xgal.Rectangle { return xgal.Rect(0, 0, windowWidth, 28) }
func (a *App) canvasBounds() xgal.Rectangle  { return xgal.Rect(0, 28, windowHeight, 396) }
func (a *App) listBounds() xgal.Rectangle    { return xgal.Rect(482, 30, 638, 394) }
func (a *App) paletteBounds() xgal.Rectangle { return xgal.Rect(0, 396, windowWidth, 460) }
func (a *App) statusBounds() xgal.Rectangle  { return xgal.Rect(0, 460, windowWidth, windowHeight) }

func ctrlHeld() bool {
	for _, k := range xgal.Keys() {
		if k == xgal.KeyControl || k == xgal.KeyControlLeft || k == xgal.KeyControlRight {
			return true
		}
	}
	return false
}

func (a *App) save() {
	if a.filename == "" {
		a.msg = "No filename set (run with -f <file>)"
		a.msgTimer = 180
		return
	}
	f, err := os.Create(a.filename)
	if err != nil {
		a.msg = fmt.Sprintf("Error saving: %v", err)
		a.msgTimer = 180
		return
	}
	if err := a.doc.Encode(f); err != nil {
		a.msg = fmt.Sprintf("Error saving: %v", err)
		a.msgTimer = 180
		f.Close()
		return
	}
	f.Close()
	a.msg = fmt.Sprintf("Saved %s", a.filename)
	a.msgTimer = 180
	a.dirty = false
}

func (a *App) load(path string) {
	f, err := os.Open(path)
	if err != nil {
		a.msg = fmt.Sprintf("Error loading: %v", err)
		a.msgTimer = 180
		return
	}
	defer f.Close()
	d := &xvec.XVEC{}
	if err := d.Decode(f); err != nil {
		a.msg = fmt.Sprintf("Error loading: %v", err)
		a.msgTimer = 180
		return
	}
	a.doc = d
	a.filename = path
	a.selInst = -1
	a.pend = nil
	a.dirty = true
	a.syncList()
	a.msg = fmt.Sprintf("Loaded %s", path)
	a.msgTimer = 180
}

func (a *App) Update() error {
	// Tool hotkeys: 1–6 and F2–F7
	for i := range int(toolCount) {
		if xgal.Tap(toolDigits[i]) || xgal.Tap(toolFKeys[i]) {
			if a.tool != Tool(i) {
				a.pathSteps = nil
				a.toolSel = i
				a.tool = Tool(i)
				a.pend = nil
			}
		}
	}

	// Delete selected instruction
	if (xgal.Tap(xgal.KeyDelete) || xgal.Tap(xgal.KeyBackspace)) && a.selInst >= 0 && a.selInst < len(a.doc.Instructions) {
		a.doc.Instructions = slices.Delete(a.doc.Instructions, a.selInst, a.selInst+1)
		a.selInst = -1
		a.dirty = true
		a.syncList()
	}

	// Clear all: X
	if xgal.Tap(xgal.KeyX) && len(a.doc.Instructions) > 0 {
		a.doc.Instructions = nil
		a.selInst = -1
		a.dirty = true
		a.syncList()
	}

	// Save: Ctrl+S
	if xgal.Tap(xgal.KeyS) && ctrlHeld() {
		a.save()
	}

	// Help toggle: F1 / Esc closes
	if xgal.Tap(xgal.KeyF1) {
		a.showHelp = !a.showHelp
	}
	if a.showHelp && xgal.Tap(xgal.KeyEscape) {
		a.showHelp = false
	}

	// Skip UI interaction while help is open
	if a.showHelp {
		return nil
	}

	// Poll path tools when path mode active
	if a.tool == ToolStroke || a.tool == ToolFill {
		for _, t := range a.pathToggles {
			t.Poll()
		}
	} else {
		// Poll toolbar toggles
		for _, t := range a.toggles {
			t.Poll()
		}
	}

	a.pollCanvas()
	a.pollPalette()
	a.pollList()
	a.swSlider.Poll()

	// Message timer
	if a.msgTimer > 0 {
		a.msgTimer--
	}

	return nil
}

func (a *App) canvasDocXY() (float32, float32) {
	cv := a.canvasBounds()
	mx, my := xgal.Mouse().X, xgal.Mouse().Y
	docW := float32(a.doc.Size.W)
	docH := float32(a.doc.Size.H)
	cvW := float32(cv.Dx())
	cvH := float32(cv.Dy())
	dx := (float32(mx-cv.Min.X) / cvW) * docW
	dy := (float32(my-cv.Min.Y) / cvH) * docH
	return dx, dy
}

func (a *App) pollCanvas() {
	if !xgal.Click(xgal.MouseButtonLeft) {
		return
	}
	cv := a.canvasBounds()
	mx, my := xgal.Mouse().X, xgal.Mouse().Y
	if mx < cv.Min.X || mx >= cv.Max.X || my < cv.Min.Y || my >= cv.Max.Y {
		return
	}

	docW := float32(a.doc.Size.W)
	docH := float32(a.doc.Size.H)
	cvW := float32(cv.Dx())
	cvH := float32(cv.Dy())

	dx := (float32(mx-cv.Min.X) / cvW) * docW
	dy := (float32(my-cv.Min.Y) / cvH) * docH

	if a.tool == ToolStroke || a.tool == ToolFill {
		a.pollCanvasPath(dx, dy)
		return
	}

	switch a.tool {
	case ToolLine:
		if a.pend == nil {
			a.pend = &struct{ x, y float32 }{dx, dy}
			return
		}
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.LineInstruction{X1: a.pend.x, Y1: a.pend.y,
				X2: dx, Y2: dy,
				Stroke: xvec.Length(a.defSW), Color: a.color, Antialias: true})
		a.pend = nil

	case ToolCircle:
		if a.pend == nil {
			a.pend = &struct{ x, y float32 }{dx, dy}
			return
		}
		dx2 := float64(dx - a.pend.x)
		dy2 := float64(dy - a.pend.y)
		r := float32(math.Sqrt(dx2*dx2 + dy2*dy2))
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.CircleInstruction{C: xvec.V(a.pend.x, a.pend.y), R: xvec.Length(r),
				Stroke: xvec.Length(a.defSW), Color: a.color, Antialias: true})
		a.pend = nil

	case ToolDisk:
		if a.pend == nil {
			a.pend = &struct{ x, y float32 }{dx, dy}
			return
		}
		dx2 := float64(dx - a.pend.x)
		dy2 := float64(dy - a.pend.y)
		r := float32(math.Sqrt(dx2*dx2 + dy2*dy2))
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.DiskInstruction{C: xvec.V(a.pend.x, a.pend.y), R: xvec.Length(r),
				Color: a.color, Antialias: true})
		a.pend = nil

	case ToolRect:
		if a.pend == nil {
			a.pend = &struct{ x, y float32 }{dx, dy}
			return
		}
		x := min(a.pend.x, dx)
		y := min(a.pend.y, dy)
		w := max(a.pend.x, dx) - x
		h := max(a.pend.y, dy) - y
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.RectInstruction{X: x, Y: y, W: w, H: h,
				Stroke: xvec.Length(a.defSW), Color: a.color, Antialias: true})
		a.pend = nil

	case ToolSlab:
		if a.pend == nil {
			a.pend = &struct{ x, y float32 }{dx, dy}
			return
		}
		x := min(a.pend.x, dx)
		y := min(a.pend.y, dy)
		w := max(a.pend.x, dx) - x
		h := max(a.pend.y, dy) - y
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.SlabInstruction{X: x, Y: y, W: w, H: h,
				Color: a.color, Antialias: true})
		a.pend = nil
	}

	a.dirty = true
	a.syncList()
}

func (a *App) pollPalette() {
	if !xgal.Click(xgal.MouseButtonLeft) {
		return
	}
	pb := a.paletteBounds()
	mx, my := xgal.Mouse().X, xgal.Mouse().Y
	if mx < pb.Min.X || mx >= pb.Max.X || my < pb.Min.Y || my >= pb.Max.Y {
		return
	}
	cols := 16
	cw := pb.Dx() / cols
	ch := pb.Dy() / 4
	c := (mx - pb.Min.X) / cw
	r := (my - pb.Min.Y) / ch
	if c >= 0 && c < cols && r >= 0 && r < 4 {
		idx := r*cols + c
		if idx < 64 {
			a.color = a.palColors[idx]
			a.palSel = idx
		}
	}
}

func (a *App) pollList() {
	res := a.list.Poll()
	if res == xui.Accept && a.list.Selected >= 0 && a.list.Selected < len(a.doc.Instructions) {
		a.selInst = a.list.Selected
	}
}

func (a *App) syncList() {
	a.list.Items = make([]string, len(a.doc.Instructions))
	for i, inst := range a.doc.Instructions {
		a.list.Items[i] = a.instLabel(i, inst)
	}
	a.list.Selected = a.selInst
}

func (a *App) pollCanvasPath(dx, dy float32) {
	if len(a.pathSteps) == 0 {
		a.pathSteps = append(a.pathSteps, xvec.MoveTo(dx, dy))
		return
	}
	switch a.pathSegSel {
	case 0:
		a.pathSteps = append(a.pathSteps, xvec.MoveTo(dx, dy))
	case 1:
		a.pathSteps = append(a.pathSteps, xvec.LineTo(dx, dy))
	}
}

func (a *App) pathClose() {
	if len(a.pathSteps) == 0 {
		return
	}
	a.pathSteps = append(a.pathSteps, xvec.Close())
	if a.tool == ToolStroke {
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.StrokeInstruction{
				Color: a.color, Stroke: xvec.Length(a.defSW),
				Steps: a.pathSteps, Antialias: true,
			})
	} else {
		a.doc.Instructions = append(a.doc.Instructions,
			&xvec.FillInstruction{
				Color: a.color, Steps: a.pathSteps, Antialias: true,
			})
	}
	a.pathSteps = nil
	a.dirty = true
	a.syncList()
}

func (a *App) instLabel(i int, inst xvec.Instruction) string {
	switch v := inst.(type) {
	case *xvec.CircleInstruction:
		return fmt.Sprintf("%3d Circle  (%.0f,%.0f) r=%.0f", i, v.C.X, v.C.Y, float32(v.R))
	case *xvec.DiskInstruction:
		return fmt.Sprintf("%3d Disk    (%.0f,%.0f) r=%.0f", i, v.C.X, v.C.Y, float32(v.R))
	case *xvec.RectInstruction:
		return fmt.Sprintf("%3d Rect    (%.0f,%.0f) %.0f×%.0f", i, v.X, v.Y, v.W, v.H)
	case *xvec.SlabInstruction:
		return fmt.Sprintf("%3d Slab    (%.0f,%.0f) %.0f×%.0f", i, v.X, v.Y, v.W, v.H)
	case *xvec.LineInstruction:
		return fmt.Sprintf("%3d Line    (%.0f,%.0f)→(%.0f,%.0f)", i, v.X1, v.Y1, v.X2, v.Y2)
	case *xvec.FillInstruction:
		return fmt.Sprintf("%3d Fill    %d steps", i, len(v.Steps))
	case *xvec.StrokeInstruction:
		return fmt.Sprintf("%3d Stroke  %d steps w=%.0f", i, len(v.Steps), float32(v.Stroke))
	default:
		return fmt.Sprintf("%3d ?", i)
	}
}

func (a *App) Draw(screen *xgal.Surface) {
	xgal.Box(screen, xgal.Rect(0, 0, windowWidth, windowHeight), xgal.Wash(40, 40, 60, 255))

	a.drawToolbar(screen)
	if a.tool == ToolStroke || a.tool == ToolFill {
		a.drawPathSubToolbar(screen)
	}

	a.drawCanvas(screen)
	a.drawList(screen)
	a.drawStrokeSlider(screen)
	a.drawPalette(screen)
	a.drawStatus(screen)

	if a.showHelp {
		a.drawHelpOverlay(screen)
	} else if a.msgTimer > 0 {
		a.drawMessage(screen)
	}
}

func (a *App) drawToolbar(screen *xgal.Surface) {
	tb := a.toolbarBounds()
	xgal.Box(screen, tb, xgal.Wash(25, 25, 45, 255))
	for _, t := range a.toggles {
		t.Render(screen)
	}
}

func (a *App) drawCanvas(screen *xgal.Surface) {
	cv := a.canvasBounds()
	xgal.Box(screen, cv, xgal.Wash(20, 20, 35, 255))
	xgal.Outline(screen, cv, 1, xgal.Wash(80, 80, 110, 255))

	if a.dirty || a.docSurf == nil {
		a.renderDoc()
		a.dirty = false
	}
	if a.docSurf == nil {
		return
	}

	docW := float32(a.doc.Size.W)
	docH := float32(a.doc.Size.H)
	cvW := float32(cv.Dx())
	cvH := float32(cv.Dy())

	s := cvW / docW
	if docH*s > cvH {
		s = cvH / docH
	}
	outW := docW * s
	outH := docH * s
	offX := float32(cv.Min.X) + (float32(cv.Dx())-outW)/2
	offY := float32(cv.Min.Y) + (float32(cv.Dy())-outH)/2

	xgal.Blit(screen, a.docSurf,
		xgal.Rect(int(offX), int(offY), int(offX+outW), int(offY+outH)),
		a.docSurf.Bounds())

	// Preview pending shapes
	if a.pend != nil {
		a.drawPreviews(screen, cv)
	}

	// Path preview
	if a.tool == ToolStroke || a.tool == ToolFill {
		a.drawPathPreview(screen, cv, offX, offY, outW, outH, docW, docH)
	}
}

func (a *App) drawPreviews(screen *xgal.Surface, cv xgal.Rectangle) {
	docW := float32(a.doc.Size.W)
	docH := float32(a.doc.Size.H)
	cvW := float32(cv.Dx())
	cvH := float32(cv.Dy())

	s := cvW / docW
	if docH*s > cvH {
		s = cvH / docH
	}
	outW := docW * s
	outH := docH * s
	offX := float32(cv.Min.X) + (float32(cv.Dx())-outW)/2
	offY := float32(cv.Min.Y) + (float32(cv.Dy())-outH)/2

	mx, my := xgal.Mouse().X, xgal.Mouse().Y
	sx := int(a.pend.x/docW*outW + offX)
	sy := int(a.pend.y/docH*outH + offY)

	switch a.tool {
	case ToolLine:
		xgal.Line(screen, sx, sy, mx, my, 1, xgal.Wash(255, 255, 255, 140))

	case ToolCircle, ToolDisk:
		dx := float64(mx - sx)
		dy := float64(my - sy)
		r := int(math.Sqrt(dx*dx + dy*dy))
		if r < 1 {
			r = 1
		}
		xgal.Circle(screen, xgal.Pt(sx, sy), r, 1, xgal.Wash(255, 255, 255, 140))

	case ToolRect, ToolSlab:
		rx, ry := sx, sy
		rw, rh := mx, my
		if rw < rx {
			rx, rw = rw, rx
		}
		if rh < ry {
			ry, rh = rh, ry
		}
		xgal.Outline(screen, xgal.Rect(rx, ry, rw, rh), 1, xgal.Wash(255, 255, 255, 140))
	}
}

func (a *App) drawPathPreview(screen *xgal.Surface, cv xgal.Rectangle, offX, offY, outW, outH, docW, docH float32) {
	prev := screen
	mp := xgal.Mouse()
	// Draw lines connecting path vertices
	for i := 1; i < len(a.pathSteps); i++ {
		p0 := a.pathPoint(i-1, offX, offY, outW, outH, docW, docH)
		p1 := a.pathPoint(i, offX, offY, outW, outH, docW, docH)
		if p0 != nil && p1 != nil {
			xgal.Line(prev, p0.X, p0.Y, p1.X, p1.Y, 1, xgal.Wash(100, 200, 255, 180))
		}
	}
	// Vertex dots
	for i := range a.pathSteps {
		p := a.pathPoint(i, offX, offY, outW, outH, docW, docH)
		if p != nil {
			col := xgal.Wash(100, 200, 255, 220)
			if i == 0 {
				col = xgal.Wash(100, 255, 100, 220)
			}
			xgal.Box(prev, xgal.Rect(p.X-3, p.Y-3, p.X+4, p.Y+4), col)
		}
	}
	// Preview line from last vertex to mouse
	if len(a.pathSteps) > 0 {
		last := a.pathPoint(len(a.pathSteps)-1, offX, offY, outW, outH, docW, docH)
		bounds := xgal.Rect(cv.Min.X, cv.Min.Y, cv.Min.X+int(outW), cv.Min.Y+int(outH))
		if last != nil && mp.In(bounds) {
			xgal.Line(prev, last.X, last.Y, mp.X, mp.Y, 1, xgal.Wash(100, 200, 255, 100))
		}
	}
}

func (a *App) pathPoint(idx int, offX, offY, outW, outH, docW, docH float32) *xgal.Point {
	s := a.pathSteps[idx]
	var x, y float32
	switch v := s.(type) {
	case *xvec.MoveStep:
		x, y = v.X, v.Y
	case *xvec.LineStep:
		x, y = v.X, v.Y
	default:
		return nil
	}
	px := int(x/docW*outW + offX)
	py := int(y/docH*outH + offY)
	return &xgal.Point{X: px, Y: py}
}

func (a *App) drawPathSubToolbar(screen *xgal.Surface) {

	sb := xgal.Rect(0, 28, windowWidth, 54)
	xgal.Box(screen, sb, xgal.Wash(25, 25, 45, 255))
	xgal.Outline(screen, sb, 1, xgal.Wash(60, 60, 80, 255))

	// Segment path toggles, including the fill and close buttons .
	for _, t := range a.pathToggles {
		t.Render(screen)
	}

	// Show vertex count in path
	n := len(a.pathSteps)
	status := fmt.Sprintf("Verts: %d  —  Click to add, Close to finish", n)
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(180, 180, 200, 255), 460, 34, status)
}

func (a *App) renderDoc() {
	w, h := int(a.doc.Size.W), int(a.doc.Size.H)
	if a.docSurf == nil || a.docSurf.Bounds().Dx() != w || a.docSurf.Bounds().Dy() != h {
		a.docSurf = xgal.NewSurface(w, h)
	}
	xgal.Clear(a.docSurf, xgal.Wash(20, 20, 35, 255))
	a.doc.Draw(a.docSurf)
}

func (a *App) drawList(screen *xgal.Surface) {
	lb := a.list.Bounds
	xgal.Box(screen, lb, xgal.Wash(30, 30, 50, 255))
	xgal.Outline(screen, lb, 1, xgal.Wash(60, 60, 80, 255))
	a.list.Render(screen)
}

func (a *App) drawStrokeSlider(screen *xgal.Surface) {
	// Background panel matching list style
	panel := xgal.Rect(482, 350, 638, 394)
	xgal.Box(screen, panel, xgal.Wash(30, 30, 50, 255))
	xgal.Outline(screen, panel, 1, xgal.Wash(60, 60, 80, 255))

	label := fmt.Sprintf("Stroke: %.0f", a.defSW)
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(200, 200, 220, 255), 490, 356, label)

	a.swSlider.Render(screen)
}

func (a *App) drawPalette(screen *xgal.Surface) {
	pb := a.paletteBounds()
	xgal.Box(screen, pb, xgal.Wash(25, 25, 45, 255))
	xgal.Outline(screen, pb, 1, xgal.Wash(60, 60, 80, 255))

	cols := 16
	rows := 4
	cw := pb.Dx() / cols
	ch := pb.Dy() / rows

	for i, col := range a.palColors {
		r := i / cols
		c := i % cols
		swatch := xgal.Rect(pb.Min.X+c*cw, pb.Min.Y+r*ch, pb.Min.X+(c+1)*cw, pb.Min.Y+(r+1)*ch)
		xgal.Box(screen, swatch.Inset(1), col)

		if i == a.palSel {
			xgal.Outline(screen, swatch.Inset(1), 2, xgal.Wash(255, 255, 255, 255))
		}
	}
}

func (a *App) drawStatus(screen *xgal.Surface) {
	sb := a.statusBounds()
	xgal.Box(screen, sb, xgal.Wash(25, 25, 45, 255))
	xgal.Outline(screen, sb, 1, xgal.Wash(60, 60, 80, 255))

	c := a.color
	fn := a.filename
	if fn == "" {
		fn = "(no file)"
	}
	text := fmt.Sprintf("Tool: %s  |  #%02x%02x%02x  |  %d inst  |  %s",
		toolNames[a.tool], c.R, c.G, c.B, len(a.doc.Instructions), fn)

	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(200, 200, 220, 255), sb.Min.X+6, sb.Min.Y+4, text)

	// F1 hint on the right side
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(120, 120, 160, 255), 580, sb.Min.Y+4, "F1 help")
}

func (a *App) drawMessage(screen *xgal.Surface) {
	// Semi-transparent overlay across the status area
	msgBounds := xgal.Rect(160, 440, windowHeight, 476)
	xgal.Box(screen, msgBounds, xgal.Wash(0, 0, 0, 200))
	xgal.Outline(screen, msgBounds, 1, xgal.Wash(100, 100, 200, 255))
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(255, 255, 255, 255),
		msgBounds.Min.X+8, msgBounds.Min.Y+6, a.msg)
}

func (a *App) drawHelpOverlay(screen *xgal.Surface) {
	// Dim the background
	xgal.Box(screen, xgal.Rect(0, 0, windowWidth, windowHeight), xgal.Wash(0, 0, 0, 200))

	// Panel
	panel := xgal.Rect(140, 50, 500, 430)
	xgal.Box(screen, panel, xgal.Wash(20, 20, 40, 240))
	xgal.Outline(screen, panel, 1, xgal.Wash(120, 120, 220, 255))

	y := 70
	for _, ln := range helpLines {
		col := xgal.Wash(200, 200, 220, 255)
		if ln.bold {
			col = xgal.Wash(255, 255, 255, 255)
		}
		xgal.Ink(screen, xgal.BuiltinFace, col, 160, y, ln.text)
		y += 18
	}
}

func (a *App) Layout(w, h int) (int, int) {
	return windowWidth, windowHeight
}
