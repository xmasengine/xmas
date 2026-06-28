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

	// Layout constants — everything is derived from these so all regions
	// scale together when windowWidth / windowHeight change.
	toolbarHeight = 28
	paletteHeight = 64
	statusHeight  = 20

	rightPanelX0  = windowWidth*3/4 + 2       // 482
	rightPanelX1  = windowWidth - 2             // 638
	listPanelY0   = toolbarHeight + 2           // 30
	listPanelY1   = windowHeight - paletteHeight - statusHeight - 46  // 350
	sliderPanelY0 = listPanelY1                 // 350
	sliderPanelY1 = sliderPanelY0 + 44          // 394
	sliderX0      = rightPanelX0 + 8            // 490
	sliderX1      = windowWidth - 10            // 630

	paletteY0  = windowHeight - paletteHeight - statusHeight  // 396
	paletteY1  = windowHeight - statusHeight                  // 460
	statusY0   = paletteY1                                    // 460

	canvasX1 = rightPanelX0 - 2  // 480
	canvasY1 = paletteY0         // 396

	helpPanelW   = 360
	helpPanelH   = 380
	helpPanelX0  = (windowWidth - helpPanelW) / 2   // 140
	helpPanelX1  = helpPanelX0 + helpPanelW          // 500
	helpPanelY0  = (windowHeight - helpPanelH) / 2   // 50
	helpPanelY1  = helpPanelY0 + helpPanelH          // 430
	helpTextX    = helpPanelX0 + 20                  // 160
	helpLineY0   = helpPanelY0 + 20                  // 70
	helpLineStep = 18

	messageW  = 320
	messageH  = 36
	messageX0 = (windowWidth - messageW) / 2   // 160
	messageX1 = messageX0 + messageW            // 480
	messageY0 = windowHeight - messageH - 4     // 440
	messageY1 = windowHeight - 4                // 476

	pathVertsLabelX = 460
	pathVertsLabelY = toolbarHeight + 6  // 34

	f1HintX = windowWidth - 60  // 580
)

type Tool int

const (
	ToolPick Tool = iota
	ToolCircle
	ToolDisk
	ToolRect
	ToolSlab
	ToolLine
	ToolStroke
	ToolFill
	toolCount
)

// UI colours
var (
	colBG         = xgal.Wash(25, 25, 45, 255)
	colCanvas     = xgal.Wash(20, 20, 35, 255)
	colHelpPanel  = xgal.Wash(20, 20, 40, 240)
	colOutline    = xgal.Wash(60, 60, 80, 255)
	colOutlineMsg = xgal.Wash(100, 100, 200, 255)
	colOutlineHlp = xgal.Wash(120, 120, 220, 255)
	colText       = xgal.Wash(200, 200, 220, 255)
	colTextDim    = xgal.Wash(120, 120, 160, 255)
	colWhite      = xgal.Wash(255, 255, 255, 255)
	colPreview    = xgal.Wash(255, 255, 255, 140)
	colPathLine   = xgal.Wash(100, 200, 255, 180)
	colPathVert   = xgal.Wash(100, 200, 255, 220)
	colPathFirst  = xgal.Wash(100, 255, 100, 220)
	colPathMouse  = xgal.Wash(100, 200, 255, 100)
	colOverlay    = xgal.Wash(0, 0, 0, 200)
	colSelect     = xgal.Wash(255, 230, 100, 220)
	colSelectFill = xgal.Wash(255, 230, 100, 60)
)

var toolNames = []string{"Pick", "Circle", "Disk", "Rect", "Slab", "Line", "Stroke", "Fill"}
var toolFKeys = []xgal.KeyCode{xgal.KeyF2, xgal.KeyF3, xgal.KeyF4, xgal.KeyF5, xgal.KeyF6, xgal.KeyF7, xgal.KeyF8, xgal.KeyF9}
var toolDigits = []xgal.KeyCode{xgal.KeyDigit1, xgal.KeyDigit2, xgal.KeyDigit3, xgal.KeyDigit4, xgal.KeyDigit5, xgal.KeyDigit6, xgal.KeyDigit7, xgal.KeyDigit8}

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
	{"F2–F9 / 1–7   Select tool", false},
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

	dragging bool
	dragLast xvec.Vertex
}

func main() {
	file := flag.String("f", "", "xvec file to edit")
	flag.Parse()

	a := &App{
		doc: &xvec.XVEC{
			Size:      xvec.Size{W: 320, H: 240},
			Antialias: true,
		},
		tool:     ToolPick,
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
					a.pathSegSel = 1
					a.tool = Tool(t.Idx)
					a.pend = nil
				}
			}
		}
		t = xui.Toggle(xgal.Rect(i*btnW, 0, (i+1)*btnW, toolbarHeight), toolNames[i], toggled)
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
			if active && len(a.pathSteps) > 0 {
				a.pathSteps = append(a.pathSteps, xvec.Close())
				a.pathSegSel = 1
			}
		},
		func(active bool) {
			if active {
				a.pathFinish()
				a.toolSel = 0
				a.tool = Tool(a.toolSel)
			}
		},
	}
	for i, name := range segNames {
		// Draw over the normal toggles.
		t := xui.Toggle(xgal.Rect(i*btnW, 0, (i+1)*btnW, toolbarHeight), name, segFuncs[i])
		t.Style = xui.DefaultStyle()
		t.Group = &a.pathSegSel
		t.Idx = i
		t.Active = i == 1 // Line selected by default
		a.pathToggles = append(a.pathToggles, t)
	}
	a.pathSegSel = 1

	a.swSlider = xui.Slider(a.sliderBounds(), func(pos int) {
		a.defSW = float32(pos)
		if a.selInst >= 0 && a.selInst < len(a.doc.Instructions) {
			inst := a.doc.Instructions[a.selInst]
			if adj, ok := inst.(xvec.Adjuster); ok {
				adj.Adjust(xvec.Length(pos))
				a.dirty = true
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

func (a *App) toolbarBounds() xgal.Rectangle { return xgal.Rect(0, 0, windowWidth, toolbarHeight) }
func (a *App) canvasBounds() xgal.Rectangle  { return xgal.Rect(0, toolbarHeight, canvasX1, canvasY1) }
func (a *App) listBounds() xgal.Rectangle    { return xgal.Rect(rightPanelX0, listPanelY0, rightPanelX1, listPanelY1) }
func (a *App) paletteBounds() xgal.Rectangle { return xgal.Rect(0, paletteY0, windowWidth, paletteY1) }
func (a *App) statusBounds() xgal.Rectangle  { return xgal.Rect(0, statusY0, windowWidth, windowHeight) }
func (a *App) sliderBounds() xgal.Rectangle  { return xgal.Rect(sliderX0, sliderPanelY0, sliderX1, sliderPanelY1) }

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
				a.pathSegSel = 1
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

	// Rearrange selected instruction
	if a.selInst >= 0 && a.selInst < len(a.doc.Instructions) {
		if xgal.Tap(xgal.KeyPageUp) && a.selInst > 0 {
			a.doc.MoveUp(a.selInst)
			a.selInst--
			a.list.Selected = a.selInst
			a.syncList()
			a.dirty = true
		}
		if xgal.Tap(xgal.KeyPageDown) && a.selInst < len(a.doc.Instructions)-1 {
			a.doc.MoveDown(a.selInst)
			a.selInst++
			a.list.Selected = a.selInst
			a.syncList()
			a.dirty = true
		}
		if xgal.Tap(xgal.KeyHome) && a.selInst > 0 {
			a.doc.MoveToFront(a.selInst)
			a.selInst = 0
			a.list.Selected = 0
			a.syncList()
			a.dirty = true
		}
		if xgal.Tap(xgal.KeyEnd) && a.selInst < len(a.doc.Instructions)-1 {
			a.doc.MoveToBack(a.selInst)
			a.selInst = len(a.doc.Instructions) - 1
			a.list.Selected = a.selInst
			a.syncList()
			a.dirty = true
		}
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
	a.pollDrag()
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

func (a *App) pollDrag() {
	if a.tool != ToolPick || !a.dragging || a.selInst < 0 || a.selInst >= len(a.doc.Instructions) {
		a.dragging = false
		return
	}
	if xgal.Grip(xgal.MouseButtonLeft) {
		dx, dy := a.canvasDocXY()
		if dx != a.dragLast.X || dy != a.dragLast.Y {
			xvec.Move(a.doc.Instructions[a.selInst], dx-a.dragLast.X, dy-a.dragLast.Y)
			a.dragLast = xvec.V(dx, dy)
			a.dirty = true
		}
	} else {
		a.dragging = false
		a.syncList()
	}
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
	case ToolPick:
		inst := xvec.Pick(a.doc, dx, dy)
		if inst != nil {
			for i, v := range a.doc.Instructions {
				if v == inst {
					if i == a.selInst {
						a.dragging = true
						a.dragLast = xvec.V(dx, dy)
						return
					}
					a.selInst = i
					a.list.Selected = i
					a.syncList()
					if sw := xvec.StrokeWidth(inst); sw > 0 {
						a.defSW = sw
						a.swSlider.Pos = int(sw)
					}
					c := xvec.StrokeColor(inst)
					a.color = c
					a.palSel = palIndex(c)
					break
				}
			}
		} else {
			a.selInst = -1
			a.list.Selected = -1
			a.syncList()
		}
		return // skip dirty + syncList below

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
			if a.tool == ToolPick && a.selInst >= 0 && a.selInst < len(a.doc.Instructions) {
				if painter, ok := a.doc.Instructions[a.selInst].(xvec.Painter); ok {
					painter.Paint(a.color)
					a.dirty = true
					a.syncList()
				}
			}
		}
	}
}

func (a *App) pollList() {
	// Clicking on the slider must not also select a list item — the slider
	// area sits entirely inside the list bounds.
	if xgal.Click(xgal.MouseButtonLeft) && xgal.Mouse().In(a.sliderBounds()) {
		return
	}
	res := a.list.Poll()
	if res == xui.Accept && a.list.Selected >= 0 && a.list.Selected < len(a.doc.Instructions) {
		a.selInst = a.list.Selected
	}
}

func palIndex(c xvec.Color) int {
	q := func(v uint8) int {
		idx := (int(v) + 42) / 85
		if idx < 0 {
			return 0
		}
		if idx > 3 {
			return 3
		}
		return idx
	}
	return q(c.R)*16 + q(c.G)*4 + q(c.B)
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
	default:
		a.pathSteps = append(a.pathSteps, xvec.LineTo(dx, dy))
	}
}

func (a *App) pathFinish() {
	if len(a.pathSteps) == 0 {
		return
	}
	// Auto-close if the path doesn't end with a CloseStep.
	if _, ok := a.pathSteps[len(a.pathSteps)-1].(*xvec.CloseStep); !ok {
		a.pathSteps = append(a.pathSteps, xvec.Close())
	}
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
	a.pathSegSel = 1
	a.pend = nil
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
	xgal.Box(screen, xgal.Rect(0, 0, windowWidth, windowHeight), colBG)

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
	xgal.Box(screen, tb, colBG)
	for _, t := range a.toggles {
		t.Render(screen)
	}
}

func (a *App) drawCanvas(screen *xgal.Surface) {
	cv := a.canvasBounds()
	xgal.Box(screen, cv, colCanvas)
	xgal.Outline(screen, cv, 1, colOutline)

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

	a.drawSelection(screen, offX, offY, outW, outH, docW, docH)
}

func (a *App) drawSelection(screen *xgal.Surface, offX, offY, outW, outH, docW, docH float32) {
	if a.tool != ToolPick || a.selInst < 0 || a.selInst >= len(a.doc.Instructions) {
		return
	}
	sx := func(x float32) int { return int(x/docW*outW + offX) }
	sy := func(y float32) int { return int(y/docH*outH + offY) }

	switch v := a.doc.Instructions[a.selInst].(type) {
	case *xvec.CircleInstruction:
		r := int(float32(v.R)/docW*outW + 1.5)
		cx, cy := sx(v.C.X), sy(v.C.Y)
		xgal.Circle(screen, xgal.Pt(cx, cy), r, 1, colSelect)
		xgal.Box(screen, xgal.Rect(cx-1, cy-1, cx+2, cy+2), colSelect)

	case *xvec.DiskInstruction:
		r := int(float32(v.R)/docW*outW + 1.5)
		cx, cy := sx(v.C.X), sy(v.C.Y)
		xgal.Circle(screen, xgal.Pt(cx, cy), r, 1, colSelect)
		xgal.Box(screen, xgal.Rect(cx-1, cy-1, cx+2, cy+2), colSelect)

	case *xvec.RectInstruction:
		x1, y1 := sx(v.X), sy(v.Y)
		x2, y2 := sx(v.X+v.W), sy(v.Y+v.H)
		xgal.Outline(screen, xgal.Rect(x1, y1, x2, y2), 1, colSelect)
		a.drawHandles(screen, x1, y1, x2, y2)

	case *xvec.SlabInstruction:
		x1, y1 := sx(v.X), sy(v.Y)
		x2, y2 := sx(v.X+v.W), sy(v.Y+v.H)
		xgal.Outline(screen, xgal.Rect(x1, y1, x2, y2), 1, colSelect)
		a.drawHandles(screen, x1, y1, x2, y2)

	case *xvec.LineInstruction:
		x1, y1 := sx(v.X1), sy(v.Y1)
		x2, y2 := sx(v.X2), sy(v.Y2)
		xgal.Line(screen, x1, y1, x2, y2, 2, colSelect)
		xgal.Box(screen, xgal.Rect(x1-2, y1-2, x1+3, y1+3), colSelect)
		xgal.Box(screen, xgal.Rect(x2-2, y2-2, x2+3, y2+3), colSelect)

	case *xvec.FillInstruction:
		if x1, y1, x2, y2, ok := stepsBounds(v.Steps); ok {
			drawBounds(screen, sx(x1), sy(y1), sx(x2), sy(y2))
		}

	case *xvec.StrokeInstruction:
		if x1, y1, x2, y2, ok := stepsBounds(v.Steps); ok {
			drawBounds(screen, sx(x1), sy(y1), sx(x2), sy(y2))
		}
	}
}

func (a *App) drawHandles(screen *xgal.Surface, x1, y1, x2, y2 int) {
	for _, p := range []xgal.Point{
		xgal.Pt(x1, y1), xgal.Pt(x2, y1), xgal.Pt(x1, y2), xgal.Pt(x2, y2),
	} {
		xgal.Box(screen, xgal.Rect(p.X-2, p.Y-2, p.X+3, p.Y+3), colSelect)
	}
}

func drawBounds(screen *xgal.Surface, x1, y1, x2, y2 int) {
	if x2-x1 < 2 || y2-y1 < 2 {
		return
	}
	xgal.Outline(screen, xgal.Rect(x1, y1, x2, y2), 1, colSelect)
	for _, p := range []xgal.Point{
		xgal.Pt(x1, y1), xgal.Pt(x2, y1), xgal.Pt(x1, y2), xgal.Pt(x2, y2),
	} {
		xgal.Box(screen, xgal.Rect(p.X-2, p.Y-2, p.X+3, p.Y+3), colSelect)
	}
}

func stepsBounds(steps []xvec.Stepper) (xmin, ymin, xmax, ymax float32, ok bool) {
	var pts []xvec.Vertex
	for _, s := range steps {
		switch v := s.(type) {
		case *xvec.MoveStep:
			pts = append(pts, xvec.V(v.X, v.Y))
		case *xvec.LineStep:
			pts = append(pts, xvec.V(v.X, v.Y))
		case *xvec.QuadStep:
			pts = append(pts, xvec.V(v.X1, v.Y1), xvec.V(v.X2, v.Y2))
		case *xvec.CubicStep:
			pts = append(pts, xvec.V(v.X1, v.Y1), xvec.V(v.X2, v.Y2), xvec.V(v.X3, v.Y3))
		case *xvec.ArcStep:
			pts = append(pts,
				xvec.V(v.CX-v.R, v.CY), xvec.V(v.CX+v.R, v.CY),
				xvec.V(v.CX, v.CY-v.R), xvec.V(v.CX, v.CY+v.R))
		case *xvec.ArcToStep:
			pts = append(pts, xvec.V(v.X1, v.Y1), xvec.V(v.X2, v.Y2))
		}
	}
	for _, p := range pts {
		if !ok || p.X < xmin {
			xmin = p.X
		}
		if !ok || p.Y < ymin {
			ymin = p.Y
		}
		if !ok || p.X > xmax {
			xmax = p.X
		}
		if !ok || p.Y > ymax {
			ymax = p.Y
		}
		ok = true
	}
	return
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
	case ToolPick:
	case ToolLine:
		xgal.Line(screen, sx, sy, mx, my, 1, colPreview)

	case ToolCircle, ToolDisk:
		dx := float64(mx - sx)
		dy := float64(my - sy)
		r := int(math.Sqrt(dx*dx + dy*dy))
		if r < 1 {
			r = 1
		}
		xgal.Circle(screen, xgal.Pt(sx, sy), r, 1, colPreview)

	case ToolRect, ToolSlab:
		rx, ry := sx, sy
		rw, rh := mx, my
		if rw < rx {
			rx, rw = rw, rx
		}
		if rh < ry {
			ry, rh = rh, ry
		}
		xgal.Outline(screen, xgal.Rect(rx, ry, rw, rh), 1, colPreview)
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
			xgal.Line(prev, p0.X, p0.Y, p1.X, p1.Y, 1, colPathLine)
		}
	}
	// Vertex dots
	for i := range a.pathSteps {
		p := a.pathPoint(i, offX, offY, outW, outH, docW, docH)
		if p != nil {
			col := colPathVert
			if i == 0 {
				col = colPathFirst
			}
			xgal.Box(prev, xgal.Rect(p.X-3, p.Y-3, p.X+4, p.Y+4), col)
		}
	}
	// Preview line from last vertex to mouse
	if len(a.pathSteps) > 0 {
		last := a.pathPoint(len(a.pathSteps)-1, offX, offY, outW, outH, docW, docH)
		bounds := xgal.Rect(cv.Min.X, cv.Min.Y, cv.Min.X+int(outW), cv.Min.Y+int(outH))
		if last != nil && mp.In(bounds) {
			xgal.Line(prev, last.X, last.Y, mp.X, mp.Y, 1, colPathMouse)
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

	sb := xgal.Rect(0, toolbarHeight, windowWidth, toolbarHeight+26)
	xgal.Box(screen, sb, colBG)
	xgal.Outline(screen, sb, 1, colOutline)

	// Segment path toggles, including the fill and close buttons .
	for _, t := range a.pathToggles {
		t.Render(screen)
	}

	// Show vertex count in path
	n := len(a.pathSteps)
	status := fmt.Sprintf("Verts: %d  —  Click to add, Close to finish", n)
	xgal.Ink(screen, xgal.BuiltinFace, colText, pathVertsLabelX, pathVertsLabelY, status)
}

func (a *App) renderDoc() {
	w, h := int(a.doc.Size.W), int(a.doc.Size.H)
	if a.docSurf == nil || a.docSurf.Bounds().Dx() != w || a.docSurf.Bounds().Dy() != h {
		a.docSurf = xgal.NewSurface(w, h)
	}
	xgal.Clear(a.docSurf, colCanvas)
	a.doc.Draw(a.docSurf)
}

func (a *App) drawList(screen *xgal.Surface) {
	lb := a.list.Bounds
	xgal.Box(screen, lb, colBG)
	xgal.Outline(screen, lb, 1, colOutline)
	a.list.Render(screen)
}

func (a *App) drawStrokeSlider(screen *xgal.Surface) {
	// Background panel matching list style
	panel := xgal.Rect(rightPanelX0, sliderPanelY0, rightPanelX1, sliderPanelY1)
	xgal.Box(screen, panel, colBG)
	xgal.Outline(screen, panel, 1, colOutline)

	label := fmt.Sprintf("Stroke: %.0f", a.defSW)
	xgal.Ink(screen, xgal.BuiltinFace, colText, sliderX0, sliderPanelY0+6, label)

	a.swSlider.Render(screen)
}

func (a *App) drawPalette(screen *xgal.Surface) {
	pb := a.paletteBounds()
	xgal.Box(screen, pb, colBG)
	xgal.Outline(screen, pb, 1, colOutline)

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
			xgal.Outline(screen, swatch.Inset(1), 2, colWhite)
		}
	}
}

func (a *App) drawStatus(screen *xgal.Surface) {
	sb := a.statusBounds()
	xgal.Box(screen, sb, colBG)
	xgal.Outline(screen, sb, 1, colOutline)

	c := a.color
	fn := a.filename
	if fn == "" {
		fn = "(no file)"
	}
	text := fmt.Sprintf("Tool: %s  |  #%02x%02x%02x  |  %d inst  |  %s",
		toolNames[a.tool], c.R, c.G, c.B, len(a.doc.Instructions), fn)

	xgal.Ink(screen, xgal.BuiltinFace, colText, sb.Min.X+6, sb.Min.Y+4, text)

	// F1 hint on the right side
	xgal.Ink(screen, xgal.BuiltinFace, colTextDim, f1HintX, sb.Min.Y+4, "F1 help")
}

func (a *App) drawMessage(screen *xgal.Surface) {
	// Semi-transparent overlay across the status area
	msgBounds := xgal.Rect(messageX0, messageY0, messageX1, messageY1)
	xgal.Box(screen, msgBounds, colOverlay)
	xgal.Outline(screen, msgBounds, 1, colOutlineMsg)
	xgal.Ink(screen, xgal.BuiltinFace, colWhite,
		msgBounds.Min.X+8, msgBounds.Min.Y+6, a.msg)
}

func (a *App) drawHelpOverlay(screen *xgal.Surface) {
	// Dim the background
	xgal.Box(screen, xgal.Rect(0, 0, windowWidth, windowHeight), colOverlay)

	// Panel
	panel := xgal.Rect(helpPanelX0, helpPanelY0, helpPanelX1, helpPanelY1)
	xgal.Box(screen, panel, colHelpPanel)
	xgal.Outline(screen, panel, 1, colOutlineHlp)

	y := helpLineY0
	for _, ln := range helpLines {
		col := colText
		if ln.bold {
			col = colWhite
		}
		xgal.Ink(screen, xgal.BuiltinFace, col, helpTextX, y, ln.text)
		y += helpLineStep
	}
}

func (a *App) Layout(w, h int) (int, int) {
	return windowWidth, windowHeight
}
