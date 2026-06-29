package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xui"
)

const (
	WindowW    = 640
	WindowH    = 480
	ToolbarH   = 28
	StatusH    = 20
	PalCell    = 15
	MaxPalRows = 3

	toggleW    = 54
	btnW       = 44
	sliderW    = WindowW - 8*toggleW - 2*btnW // 640 - 432 - 88 = 120
	numToggles = 8
)

type Tool int

const (
	ToolPencil Tool = iota
	ToolSelect
	ToolLine
	ToolStroke
	ToolFill
	ToolEyedropper
	ToolEraser
	ToolCircle
	ToolCount
)

var toolNames = []string{"Pencil", "Select", "Line", "Stroke", "Fill", "Eye", "Eraser", "Circle"}

type App struct {
	doc      *xgal.Paletted
	dirty    bool
	docSurf  *xgal.Surface
	filename string

	tool    Tool
	toolGrp *xui.ToggleGroupLayer

	copyBtn  *xui.ButtonLayer
	pasteBtn *xui.ButtonLayer

	fgIdx int
	bgIdx int

	brushSize int // radius in pixels, 0 = single pixel
	sizeSl    *xui.SliderLayer

	zoom int
	offX float64
	offY float64
	palH int

	sel xgal.Rectangle

	clip    *xgal.Paletted
	pasting bool
	pastePt xgal.Point

	mx, my         int
	lastMX, lastMY int

	kbdX, kbdY int
	kbdActive  bool
	kbdClick   bool
	kbdGrip    bool
	kbdLoose   bool

	drawing bool
	drawX0  int
	drawY0  int
	drawX1  int
	drawY1  int

	msg      string
	msgTimer int
	prevTool Tool

	showHelp bool

	ask *xui.AskLayer
}

func docPalette() xgal.Palette {
	// genesis palette, but 0 is transparent.
	pal := make(xgal.Palette, 65)
	pal[0] = xgal.Wash(0, 0, 0, 0)
	for j := 1; j < len(pal); j++ {
		i := j - 1
		r := uint8((i / 16) * 85)
		g := uint8(((i / 4) % 4) * 85)
		b := uint8((i % 4) * 85)
		pal[j] = xgal.Wash(r, g, b, 255)
	}
	return pal
}

func palIndex(c xgal.Color) int {
	return docPalette().Index(c)
}

func newDoc(w, h int) *xgal.Paletted {
	img := xgal.Express(xgal.Rect(0, 0, w, h), docPalette())
	for i := range img.Pix {
		img.Pix[i] = 0
	}
	return img
}

func (a *App) computePalH() {
	n := len(a.doc.Palette)
	perRow := (WindowW) / PalCell
	if perRow < 1 {
		perRow = 1
	}
	rows := (n + perRow - 1) / perRow
	if rows > MaxPalRows {
		rows = MaxPalRows
	}
	if rows < 1 {
		rows = 1
	}
	a.palH = rows*PalCell + 4
}

func (a *App) computeZoom() {
	b := a.doc.Bounds()
	iw, ih := b.Dx(), b.Dy()
	cvW := WindowW
	cvH := WindowH - ToolbarH - a.palH - StatusH
	zx := float64(cvW) / float64(iw)
	zy := float64(cvH) / float64(ih)
	z := zx
	if zy < z {
		z = zy
	}
	az := int(math.Floor(z))
	if az < 1 {
		az = 1
	}
	a.zoom = az
	a.recenter()
}

func (a *App) recenter() {
	b := a.doc.Bounds()
	iw, ih := b.Dx(), b.Dy()
	cvW := WindowW
	cvH := WindowH - ToolbarH - a.palH - StatusH
	a.offX = float64(cvW-iw*a.zoom) / 2
	a.offY = float64(ToolbarH) + float64(cvH-ih*a.zoom)/2
}

func (a *App) panToCursor(x, y int) {
	cvW := WindowW
	cvH := WindowH - ToolbarH - a.palH - StatusH
	sx := float64(x)*float64(a.zoom) + a.offX
	sy := float64(y)*float64(a.zoom) + a.offY
	if sx < 0 {
		a.offX -= sx
	} else if sx+float64(a.zoom) > float64(cvW) {
		a.offX -= (sx + float64(a.zoom)) - float64(cvW)
	}
	if sy < float64(ToolbarH) {
		a.offY -= sy - float64(ToolbarH)
	} else if sy+float64(a.zoom) > float64(cvH+ToolbarH) {
		a.offY -= (sy + float64(a.zoom)) - float64(cvH+ToolbarH)
	}
}

func (a *App) computeLayout() {
	a.computePalH()
	a.computeZoom()
}

func main() {
	flag.Parse()
	args := flag.Args()

	app := &App{
		doc:   newDoc(64, 64),
		tool:  ToolPencil,
		fgIdx: 1,
		bgIdx: 0,
		zoom:  4,
		palH:  2*PalCell + 4,
		msg:   "Pencil – click to draw",
	}

	if len(args) > 0 {
		app.filename = args[0]
		app.loadName(app.filename)
	}

	app.computeLayout()

	toggles := make([]*xui.ToggleLayer, ToolCount)
	for i := range toggles {
		idx := Tool(i)
		t := xui.Toggle(xgal.Rect(i*toggleW, 0, (i+1)*toggleW, ToolbarH), toolNames[i], nil)
		t.Toggled = func(active bool) {
			if active {
				app.setTool(Tool(idx))
			}
		}
		toggles[i] = t
	}
	app.toolGrp = xui.NewToggleGroup(toggles...)
	app.toolGrp.Active = int(app.tool)

	copyX := numToggles * toggleW
	app.copyBtn = xui.Button(xgal.Rect(copyX, 0, copyX+btnW, ToolbarH), "Copy", func() {
		if app.pasting {
			return
		}
		app.copySelection()
	})

	pasteX := copyX + btnW
	app.pasteBtn = xui.Button(xgal.Rect(pasteX, 0, pasteX+btnW, ToolbarH), "Paste", func() {
		app.pasteClipboard()
	})

	sliderX := pasteX + btnW
	app.sizeSl = xui.Slider(xgal.Rect(sliderX, 0, WindowW, ToolbarH), func(pos int) {
		app.brushSize = pos
	})
	app.sizeSl.Low = 1
	app.sizeSl.High = 10
	app.sizeSl.Pos = 1
	app.brushSize = 1
	app.sizeSl.Low = 0
	app.sizeSl.High = 10
	app.sizeSl.Pos = 0

	xgal.Screen(WindowW, WindowH, "xpix")
	xgal.Play(app)
}

func ctrlHeld() bool {
	for _, k := range xgal.Keys() {
		if k == xgal.KeyControl || k == xgal.KeyControlLeft || k == xgal.KeyControlRight {
			return true
		}
	}
	return false
}

func (a *App) setTool(t Tool) {
	if a.tool == t {
		return
	}
	if t == ToolEyedropper {
		a.prevTool = a.tool
	}
	a.drawing = false
	a.tool = t
	a.toolGrp.Active = int(t)
}

func (a *App) loadName(name string) {
	src, err := xgal.Pixels(os.DirFS("."), name)
	if err != nil {
		a.setMsg("decode: " + err.Error())
		return
	}

	a.doc = xgal.Reduce(src, docPalette())
	a.dirty = true
	a.sel = xgal.Rectangle{}
	a.pasting = false
	a.clip = nil
	a.computeLayout()
	a.setMsg(fmt.Sprintf("loaded %s (%dx%d)", a.filename, a.doc.Bounds().Dx(), a.doc.Bounds().Dy()))
}

func (a *App) saveName(name string) string {
	a.filename = name
	ext := filepath.Ext(a.filename)
	if ext != ".png" && ext != ".gif" && ext != ".jpeg" {
		a.filename = a.filename + ".png"
	}

	err := xgal.Scribble(name, a.doc)
	if err != nil {
		return "save: " + err.Error()
	}

	return fmt.Sprintf("saved %s", a.filename)
}

func (a *App) resize(name string) {
	parts := strings.SplitN(name, "x", 2)
	if len(parts) != 2 {
		parts = strings.SplitN(name, "X", 2)
	}
	if len(parts) != 2 {
		a.setMsg("resize: expected WxH")
		return
	}
	w, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || w < 1 {
		a.setMsg("resize: invalid width")
		return
	}
	h, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || h < 1 {
		a.setMsg("resize: invalid height")
		return
	}

	b := a.doc.Bounds()
	pal := a.doc.Palette
	dst := xgal.Express(xgal.Rect(0, 0, w, h), pal)
	// Fill all with index 0 (transparent)
	for i := range dst.Pix {
		dst.Pix[i] = 0
	}
	// Copy overlapping region
	minX := min(b.Dx(), w)
	minY := min(b.Dy(), h)
	for y := 0; y < minY; y++ {
		srcOff := (y + b.Min.Y) * a.doc.Stride
		dstOff := y * dst.Stride
		for x := 0; x < minX; x++ {
			dst.Pix[dstOff+x] = a.doc.Pix[srcOff+x+b.Min.X]
		}
	}
	a.doc = dst
	a.sel = xgal.Rectangle{}
	a.pasting = false
	a.clip = nil
	a.dirty = true
	a.computeLayout()
	a.setMsg(fmt.Sprintf("resized to %dx%d", w, h))
}

func (a *App) setMsg(msg string) {
	a.msg = msg
	a.msgTimer = 180
}

func (a *App) statusText() string {
	var b string
	if a.pasting {
		b = "Paste – click to place, right-click to cancel"
	} else {
		b = toolNames[a.tool] + fmt.Sprintf(" | (%d,%d) | %dx | brush %d", a.mx, a.my, a.zoom, a.brushSize)
	}
	if a.msg != "" {
		b = b + "  |  " + a.msg
	}
	return b
}

func (a *App) screenToImg(sx, sy int) (int, int, bool) {
	ix := int((float64(sx) - a.offX) / float64(a.zoom))
	iy := int((float64(sy) - a.offY) / float64(a.zoom))
	b := a.doc.Bounds()
	if ix < b.Min.X || ix >= b.Max.X || iy < b.Min.Y || iy >= b.Max.Y {
		return 0, 0, false
	}
	return ix, iy, true
}

func (a *App) imgToScreen(ix, iy int) (int, int) {
	sx := int(float64(ix)*float64(a.zoom) + a.offX)
	sy := int(float64(iy)*float64(a.zoom) + a.offY)
	return sx, sy
}

func (a *App) rebuildSurface() {
	// A PalettedImage is also an Image.
	a.docSurf = xgal.Bake(a.doc)
	a.dirty = false
}

func (a *App) fillCircle(ix, iy int, idx uint8) {
	r := a.brushSize
	b := a.doc.Bounds()
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy > r*r+1 {
				continue
			}
			x, y := ix+dx, iy+dy
			if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
				a.doc.Pix[y*a.doc.Stride+x] = idx
			}
		}
	}
	a.dirty = true
}

func (a *App) bresenhamBrush(x0, y0, x1, y1 int, idx uint8) {
	dx := x1 - x0
	dy := y1 - y0
	adx := dx
	if adx < 0 {
		adx = -adx
	}
	ady := dy
	if ady < 0 {
		ady = -ady
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := adx - ady
	for {
		a.fillCircle(x0, y0, idx)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -ady {
			err -= ady
			x0 += sx
		}
		if e2 < adx {
			err += adx
			y0 += sy
		}
	}
}

func (a *App) setPixel(ix, iy int, idx int) {
	a.fillCircle(ix, iy, uint8(idx))
}

func (a *App) getPixel(ix, iy int) uint8 {
	b := a.doc.Bounds()
	if ix < b.Min.X || ix >= b.Max.X || iy < b.Min.Y || iy >= b.Max.Y {
		return 0
	}
	return a.doc.Pix[iy*a.doc.Stride+ix]
}

func (a *App) floodFill(ix, iy int, fillIdx int) {
	b := a.doc.Bounds()
	target := a.doc.ColorIndexAt(ix, iy)
	if target == uint8(fillIdx) {
		return
	}
	w, h := b.Dx(), b.Dy()
	fill := uint8(fillIdx)
	type pt struct{ x, y int }
	stack := []pt{{ix, iy}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if p.x < 0 || p.x >= w || p.y < 0 || p.y >= h {
			continue
		}
		if a.doc.Pix[p.y*a.doc.Stride+p.x] != target {
			continue
		}
		a.doc.Pix[p.y*a.doc.Stride+p.x] = fill
		a.dirty = true
		stack = append(stack, pt{p.x - 1, p.y}, pt{p.x + 1, p.y}, pt{p.x, p.y - 1}, pt{p.x, p.y + 1})
	}
}

func (a *App) strokeRect(x0, y0, x1, y1 int, idx int) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	b := a.doc.Bounds()
	for x := x0; x <= x1; x++ {
		if x >= b.Min.X && x < b.Max.X {
			if y0 >= b.Min.Y && y0 < b.Max.Y {
				a.doc.Pix[y0*a.doc.Stride+x] = uint8(idx)
			}
			if y1 >= b.Min.Y && y1 < b.Max.Y && y1 != y0 {
				a.doc.Pix[y1*a.doc.Stride+x] = uint8(idx)
			}
		}
	}
	for y := y0 + 1; y < y1; y++ {
		if y >= b.Min.Y && y < b.Max.Y {
			if x0 >= b.Min.X && x0 < b.Max.X {
				a.doc.Pix[y*a.doc.Stride+x0] = uint8(idx)
			}
			if x1 >= b.Min.X && x1 < b.Max.X && x1 != x0 {
				a.doc.Pix[y*a.doc.Stride+x1] = uint8(idx)
			}
		}
	}
	a.dirty = true
}

func (a *App) fillRect(x0, y0, x1, y1 int, idx int) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	b := a.doc.Bounds()
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
				a.doc.Pix[y*a.doc.Stride+x] = uint8(idx)
			}
		}
	}
	a.dirty = true
}

func (a *App) strokeCircle(x0, y0, x1, y1 int, idx int) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	cx := (x0 + x1) / 2
	cy := (y0 + y1) / 2
	rx := (x1 - x0) / 2
	ry := (y1 - y0) / 2
	if rx < 0 {
		rx = 0
	}
	if ry < 0 {
		ry = 0
	}
	// Use radius as max of rx, ry for a circle that fits the rect
	r := rx
	if ry > r {
		r = ry
	}
	if r < 1 {
		r = 1
	}
	b := a.doc.Bounds()
	set := func(x, y int) {
		if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
			a.doc.Pix[y*a.doc.Stride+x] = uint8(idx)
		}
	}
	x, y := 0, r
	d := 3 - 2*r
	for x <= y {
		set(cx+x, cy+y)
		set(cx-x, cy+y)
		set(cx+x, cy-y)
		set(cx-x, cy-y)
		set(cx+y, cy+x)
		set(cx-y, cy+x)
		set(cx+y, cy-x)
		set(cx-y, cy-x)
		if d < 0 {
			d += 4*x + 6
		} else {
			y--
			d += 4*(x-y) + 10
		}
		x++
	}
	a.dirty = true
}

func (a *App) Update() error {
	// Ask dialog takes priority over all other input
	if a.ask != nil {
		if xgal.Tap(xgal.KeyEscape) {
			a.ask = nil
			return nil
		}
		if a.ask.Poll() == xui.Finish {
			a.ask = nil
		}
		return nil
	}

	if xgal.Tap(xgal.KeyQ) || xgal.Tap(xgal.KeyEscape) {
		if a.showHelp {
			a.showHelp = false
			return nil
		}
		return xgal.Quit
	}

	if a.showHelp {
		if xgal.Tap(xgal.KeyF1) || xgal.Click(xgal.MouseButtonLeft) || xgal.Tap(xgal.KeyEscape) {
			a.showHelp = false
		}
		return nil
	}

	if xgal.Tap(xgal.KeyF1) {
		a.showHelp = true
		return nil
	}

	if xgal.Tap(xgal.KeyS) && ctrlHeld() {
		name := a.filename
		if name == "" {
			name = "untitled.png"
		}
		dw := xui.DefaultStyle().MeasureText("Save as:").X + xui.DefaultStyle().Margin.X*4
		if dw < 300 {
			dw = 300
		}
		dh := (xui.DefaultStyle().MeasureText("X").Y*2 + xui.DefaultStyle().Margin.Y*6) * 2
		bounds := xgal.Rect(WindowW/2-dw/2, WindowH/2-dh/2, WindowW/2+dw/2, WindowH/2+dh/2)
		a.ask = xui.AskEntry(bounds, "Save as:", name, func(s string) {
			a.setMsg(a.saveName(s))
		}, "Save", "Cancel")
		return nil
	}
	if xgal.Tap(xgal.KeyO) && ctrlHeld() {
		dw := xui.DefaultStyle().MeasureText("Open:").X + xui.DefaultStyle().Margin.X*4
		if dw < 300 {
			dw = 300
		}
		dh := (xui.DefaultStyle().MeasureText("X").Y*2 + xui.DefaultStyle().Margin.Y*6) * 2
		bounds := xgal.Rect(WindowW/2-dw/2, WindowH/2-dh/2, WindowW/2+dw/2, WindowH/2+dh/2)
		a.ask = xui.AskEntry(bounds, "Open:", "", func(s string) {
			a.loadName(s)
		}, "Open", "Cancel")
		return nil
	}
	if xgal.Tap(xgal.KeyR) && ctrlHeld() {
		cur := fmt.Sprintf("%dx%d", a.doc.Bounds().Dx(), a.doc.Bounds().Dy())
		dw := xui.DefaultStyle().MeasureText("Resize:").X + xui.DefaultStyle().Margin.X*4
		if dw < 300 {
			dw = 300
		}
		dh := (xui.DefaultStyle().MeasureText("X").Y*2 + xui.DefaultStyle().Margin.Y*6) * 2
		bounds := xgal.Rect(WindowW/2-dw/2, WindowH/2-dh/2, WindowW/2+dw/2, WindowH/2+dh/2)
		a.ask = xui.AskEntry(bounds, "Resize (WxH):", cur, a.resize, "Resize", "Cancel")
		return nil
	}

	if a.msgTimer > 0 {
		a.msgTimer--
		if a.msgTimer == 0 {
			a.msg = ""
		}
	}

	mx, my := xgal.Mouse().X, xgal.Mouse().Y
	ix, iy, inBounds := a.screenToImg(mx, my)
	a.mx, a.my = ix, iy

	// Zoom
	if xgal.Tap(xgal.KeyMinus) || xgal.Tap(xgal.KeyNumpadSubtract) {
		if a.zoom > 1 {
			a.zoom--
			a.recenter()
		}
	}
	if xgal.Tap(xgal.KeyEqual) || xgal.Tap(xgal.KeyNumpadAdd) {
		if a.zoom < 64 {
			a.zoom++
			a.recenter()
		}
	}
	if xgal.Tap(xgal.KeyDigit0) || xgal.Tap(xgal.KeyNumpad0) {
		a.computeZoom()
	}

	// Pan with PageUp/PageDown/Home/End
	cvW := WindowW
	cvH := WindowH - ToolbarH - a.palH - StatusH
	panStep := cvH / 2
	if xgal.Tap(xgal.KeyPageUp) {
		a.offY += float64(panStep)
	}
	if xgal.Tap(xgal.KeyPageDown) {
		a.offY -= float64(panStep)
	}
	if xgal.Tap(xgal.KeyHome) {
		a.offX = 0
		a.offY = float64(ToolbarH)
	}
	if xgal.Tap(xgal.KeyEnd) {
		b := a.doc.Bounds()
		a.offX = float64(cvW - b.Dx()*a.zoom)
		a.offY = float64(ToolbarH + cvH - b.Dy()*a.zoom)
	}

	// F-key tool switching (F1 is help, handled above)
	if !a.pasting {
		if xgal.Tap(xgal.KeyF2) {
			a.setTool(ToolPencil)
		} else if xgal.Tap(xgal.KeyF3) {
			a.setTool(ToolSelect)
		} else if xgal.Tap(xgal.KeyF4) {
			a.setTool(ToolLine)
		} else if xgal.Tap(xgal.KeyF5) {
			a.setTool(ToolStroke)
		} else if xgal.Tap(xgal.KeyF6) {
			a.setTool(ToolFill)
		} else if xgal.Tap(xgal.KeyF7) {
			a.setTool(ToolEyedropper)
		} else if xgal.Tap(xgal.KeyF8) {
			a.setTool(ToolEraser)
		} else if xgal.Tap(xgal.KeyF9) {
			a.setTool(ToolCircle)
		}
		// Digit shortcuts
		if xgal.Tap(xgal.KeyDigit1) {
			a.setTool(ToolPencil)
		}
		if xgal.Tap(xgal.KeyDigit2) {
			a.setTool(ToolSelect)
		}
		if xgal.Tap(xgal.KeyDigit3) {
			a.setTool(ToolLine)
		}
		if xgal.Tap(xgal.KeyDigit4) {
			a.setTool(ToolStroke)
		}
		if xgal.Tap(xgal.KeyDigit5) {
			a.setTool(ToolFill)
		}
		if xgal.Tap(xgal.KeyDigit6) {
			a.setTool(ToolEyedropper)
		}
		if xgal.Tap(xgal.KeyDigit7) {
			a.setTool(ToolEraser)
		}
		if xgal.Tap(xgal.KeyDigit8) {
			a.setTool(ToolCircle)
		}
	}

	// Palette input
	if !inBounds {
		a.pollPalette(mx, my)
	}

	// UI widgets
	a.toolGrp.Poll()
	a.copyBtn.Poll()
	a.pasteBtn.Poll()
	a.sizeSl.Poll()

	// Copy / Paste actions
	if xgal.Tap(xgal.KeyC) && ctrlHeld() {
		a.copySelection()
	}
	if xgal.Tap(xgal.KeyV) && ctrlHeld() {
		a.pasteClipboard()
	}
	if xgal.Tap(xgal.KeyX) && ctrlHeld() {
		a.cutSelection()
	}
	if xgal.Tap(xgal.KeyA) && ctrlHeld() {
		a.selectAll()
	}
	if xgal.Tap(xgal.KeyDelete) || xgal.Tap(xgal.KeyBackspace) {
		a.clearSelection()
	}

	if a.pasting {
		a.pollPaste(mx, my)
		return nil
	}

	// Keyboard cursor: arrow keys move cursor, space = left click
	curX, curY := ix, iy
	if a.kbdActive {
		curX, curY = a.kbdX, a.kbdY
	}
	dx, dy := 0, 0
	if xgal.Tap(xgal.KeyArrowLeft) {
		dx = -1
	}
	if xgal.Tap(xgal.KeyArrowRight) {
		dx = 1
	}
	if xgal.Tap(xgal.KeyArrowUp) {
		dy = -1
	}
	if xgal.Tap(xgal.KeyArrowDown) {
		dy = 1
	}
	if dx != 0 || dy != 0 {
		a.kbdActive = true
		curX += dx
		curY += dy
		b := a.doc.Bounds()
		if curX < b.Min.X {
			curX = b.Min.X
		}
		if curX >= b.Max.X {
			curX = b.Max.X - 1
		}
		if curY < b.Min.Y {
			curY = b.Min.Y
		}
		if curY >= b.Max.Y {
			curY = b.Max.Y - 1
		}
		a.kbdX, a.kbdY = curX, curY
		// Pan to keep cursor visible
		a.panToCursor(curX, curY)
	}

	// Keyboard space simulation
	a.kbdClick = false
	a.kbdGrip = false
	a.kbdLoose = false
	if a.kbdActive {
		if xgal.Tap(xgal.KeySpace) {
			a.kbdClick = true
		}
		if xgal.Key(xgal.KeySpace) {
			a.kbdGrip = true
		}
		if xgal.Lift(xgal.KeySpace) {
			a.kbdLoose = true
		}
	}
	// Mouse movement deactivates keyboard cursor
	if mx != a.lastMX || my != a.lastMY {
		a.kbdActive = false
	}
	a.lastMX, a.lastMY = mx, my

	a.mx, a.my = curX, curY

	// Tool-specific input
	if curX >= 0 && curY >= 0 {
		b := a.doc.Bounds()
		if curX >= b.Min.X && curX < b.Max.X && curY >= b.Min.Y && curY < b.Max.Y {
			switch a.tool {
			case ToolPencil:
				a.pollPencil(curX, curY)
			case ToolSelect:
				a.pollSelect(curX, curY)
			case ToolLine:
				a.pollLine(curX, curY)
			case ToolStroke:
				a.pollStroke(curX, curY)
			case ToolFill:
				a.pollFill(curX, curY)
			case ToolEyedropper:
				a.pollEyedropper(curX, curY)
			case ToolEraser:
				a.pollEraser(curX, curY)
			case ToolCircle:
				a.pollCircle(curX, curY)
			}
		}
	}

	if a.tool != ToolCircle && (xgal.Loose(xgal.MouseButtonLeft) || xgal.Loose(xgal.MouseButtonRight) || a.kbdLoose) {
		a.drawing = false
	}

	return nil
}

func (a *App) copySelection() {
	if a.sel.Empty() {
		a.setMsg("no selection")
		return
	}
	b := a.sel
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return
	}
	pal := make(xgal.Palette, len(a.doc.Palette))
	copy(pal, a.doc.Palette)
	a.clip = xgal.Express(xgal.Rect(0, 0, w, h), pal)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx := x + b.Min.X
			sy := y + b.Min.Y
			a.clip.Pix[y*a.clip.Stride+x] = a.doc.Pix[sy*a.doc.Stride+sx]
		}
	}
	a.setMsg(fmt.Sprintf("copied %dx%d", w, h))
}

func (a *App) cutSelection() {
	if a.sel.Empty() {
		return
	}
	a.copySelection()
	a.clearSelection()
}

func (a *App) pasteClipboard() {
	if a.clip == nil {
		a.setMsg("nothing to paste")
		return
	}
	a.pasting = true
	a.pastePt = xgal.Pt(a.mx, a.my)
}

func (a *App) placePaste() {
	if a.clip == nil || !a.pasting {
		return
	}
	b := a.clip.Bounds()
	offX := a.pastePt.X - b.Dx()/2
	offY := a.pastePt.Y - b.Dy()/2
	docB := a.doc.Bounds()
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			dx := x + offX
			dy := y + offY
			if dx >= docB.Min.X && dx < docB.Max.X && dy >= docB.Min.Y && dy < docB.Max.Y {
				a.doc.Pix[dy*a.doc.Stride+dx] = a.clip.Pix[y*a.clip.Stride+x]
			}
		}
	}
	a.dirty = true
	a.pasting = false
	a.setMsg("paste placed")
}

func (a *App) pollPaste(mx, my int) {
	ix, iy, _ := a.screenToImg(mx, my)
	a.pastePt = xgal.Pt(ix, iy)
	if a.clip != nil {
		b := a.clip.Bounds()
		docB := a.doc.Bounds()
		if a.pastePt.X < docB.Min.X+b.Dx()/2 {
			a.pastePt.X = docB.Min.X + b.Dx()/2
		}
		if a.pastePt.X > docB.Max.X-b.Dx()/2 {
			a.pastePt.X = docB.Max.X - b.Dx()/2
		}
		if a.pastePt.Y < docB.Min.Y+b.Dy()/2 {
			a.pastePt.Y = docB.Min.Y + b.Dy()/2
		}
		if a.pastePt.Y > docB.Max.Y-b.Dy()/2 {
			a.pastePt.Y = docB.Max.Y - b.Dy()/2
		}
	}
	if xgal.Click(xgal.MouseButtonLeft) {
		a.placePaste()
	}
	if xgal.Click(xgal.MouseButtonRight) {
		a.pasting = false
		a.setMsg("paste cancelled")
	}
}

func (a *App) pollPencil(ix, iy int) {
	if xgal.Click(xgal.MouseButtonLeft) || a.kbdClick {
		a.drawing = true
		a.drawX0, a.drawY0 = ix, iy
		a.drawX1, a.drawY1 = ix, iy
		a.setPixel(ix, iy, a.fgIdx)
	}
	if xgal.Click(xgal.MouseButtonRight) {
		a.drawing = true
		a.drawX0, a.drawY0 = ix, iy
		a.drawX1, a.drawY1 = ix, iy
		a.setPixel(ix, iy, a.bgIdx)
	}
	if a.drawing && (xgal.Grip(xgal.MouseButtonLeft) || a.kbdGrip) {
		if ix != a.drawX1 || iy != a.drawY1 {
			a.bresenhamBrush(a.drawX1, a.drawY1, ix, iy, uint8(a.fgIdx))
			a.drawX1, a.drawY1 = ix, iy
		}
	}
	if a.drawing && xgal.Grip(xgal.MouseButtonRight) {
		if ix != a.drawX1 || iy != a.drawY1 {
			a.bresenhamBrush(a.drawX1, a.drawY1, ix, iy, uint8(a.bgIdx))
			a.drawX1, a.drawY1 = ix, iy
		}
	}
	if a.drawing && (xgal.Loose(xgal.MouseButtonLeft) || xgal.Loose(xgal.MouseButtonRight) || a.kbdLoose) {
		a.drawing = false
	}
}

func (a *App) pollEraser(ix, iy int) {
	if xgal.Click(xgal.MouseButtonLeft) || xgal.Click(xgal.MouseButtonRight) || a.kbdClick {
		a.drawing = true
		a.drawX0, a.drawY0 = ix, iy
		a.drawX1, a.drawY1 = ix, iy
		a.setPixel(ix, iy, 0)
	}
	if a.drawing && (xgal.Grip(xgal.MouseButtonLeft) || xgal.Grip(xgal.MouseButtonRight) || a.kbdGrip) {
		if ix != a.drawX1 || iy != a.drawY1 {
			a.bresenhamBrush(a.drawX1, a.drawY1, ix, iy, 0)
			a.drawX1, a.drawY1 = ix, iy
		}
	}
	if a.drawing && (xgal.Loose(xgal.MouseButtonLeft) || xgal.Loose(xgal.MouseButtonRight) || a.kbdLoose) {
		a.drawing = false
	}
}

func (a *App) pollSelect(ix, iy int) {
	if xgal.Grip(xgal.MouseButtonLeft) || a.kbdGrip {
		if !a.drawing {
			a.drawing = true
			a.drawX0, a.drawY0 = ix, iy
		}
		a.drawX1, a.drawY1 = ix, iy
	}
	if a.drawing && (xgal.Loose(xgal.MouseButtonLeft) || a.kbdLoose) {
		x0, y0 := a.drawX0, a.drawY0
		x1, y1 := ix, iy
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		a.sel = xgal.Rect(x0, y0, x1+1, y1+1)
		a.drawing = false
	}
}

func (a *App) pollLine(ix, iy int) {
	if xgal.Grip(xgal.MouseButtonLeft) || a.kbdGrip {
		if !a.drawing {
			a.drawing = true
			a.drawX0, a.drawY0 = ix, iy
		}
		a.drawX1, a.drawY1 = ix, iy
	}
	if a.drawing && (xgal.Loose(xgal.MouseButtonLeft) || a.kbdLoose) {
		a.bresenhamBrush(a.drawX0, a.drawY0, ix, iy, uint8(a.fgIdx))
		a.drawing = false
	}
	if xgal.Grip(xgal.MouseButtonRight) {
		if !a.drawing {
			a.drawing = true
			a.drawX0, a.drawY0 = ix, iy
		}
		a.drawX1, a.drawY1 = ix, iy
	}
	if a.drawing && xgal.Loose(xgal.MouseButtonRight) {
		a.bresenhamBrush(a.drawX0, a.drawY0, ix, iy, uint8(a.bgIdx))
		a.drawing = false
	}
}

func (a *App) pollStroke(ix, iy int) {
	if !a.sel.Empty() {
		if xgal.Click(xgal.MouseButtonLeft) || a.kbdClick {
			b := a.sel
			a.strokeRect(b.Min.X, b.Min.Y, b.Max.X-1, b.Max.Y-1, a.fgIdx)
		}
		if xgal.Click(xgal.MouseButtonRight) {
			b := a.sel
			a.strokeRect(b.Min.X, b.Min.Y, b.Max.X-1, b.Max.Y-1, a.bgIdx)
		}
		return
	}
	if xgal.Grip(xgal.MouseButtonLeft) || a.kbdGrip {
		if !a.drawing {
			a.drawing = true
			a.drawX0, a.drawY0 = ix, iy
		}
		a.drawX1, a.drawY1 = ix, iy
	}
	if a.drawing && (xgal.Loose(xgal.MouseButtonLeft) || a.kbdLoose) {
		a.strokeRect(a.drawX0, a.drawY0, ix, iy, a.fgIdx)
		a.drawing = false
	}
	if xgal.Grip(xgal.MouseButtonRight) {
		if !a.drawing {
			a.drawing = true
			a.drawX0, a.drawY0 = ix, iy
		}
		a.drawX1, a.drawY1 = ix, iy
	}
	if a.drawing && xgal.Loose(xgal.MouseButtonRight) {
		a.strokeRect(a.drawX0, a.drawY0, ix, iy, a.bgIdx)
		a.drawing = false
	}
}

func (a *App) pollCircle(ix, iy int) {
	if xgal.Click(xgal.MouseButtonLeft) || a.kbdClick {
		if !a.drawing {
			a.drawing = true
			a.drawX0, a.drawY0 = ix, iy
		} else {
			dx := ix - a.drawX0
			dy := iy - a.drawY0
			r := int(math.Sqrt(float64(dx*dx + dy*dy)))
			if r < 1 {
				r = 1
			}
			a.strokeCircle(a.drawX0-r, a.drawY0-r, a.drawX0+r, a.drawY0+r, a.fgIdx)
			a.drawing = false
		}
	}
	if xgal.Click(xgal.MouseButtonRight) {
		if a.drawing {
			a.drawing = false
		}
	}
	if a.drawing {
		a.drawX1, a.drawY1 = ix, iy
	}
}

func (a *App) pollFill(ix, iy int) {
	if !a.sel.Empty() {
		if xgal.Click(xgal.MouseButtonLeft) || a.kbdClick {
			a.fillRect(a.sel.Min.X, a.sel.Min.Y, a.sel.Max.X-1, a.sel.Max.Y-1, a.fgIdx)
		}
		if xgal.Click(xgal.MouseButtonRight) {
			a.fillRect(a.sel.Min.X, a.sel.Min.Y, a.sel.Max.X-1, a.sel.Max.Y-1, a.bgIdx)
		}
	} else {
		if xgal.Click(xgal.MouseButtonLeft) || a.kbdClick {
			a.floodFill(ix, iy, a.fgIdx)
		}
		if xgal.Click(xgal.MouseButtonRight) {
			a.floodFill(ix, iy, a.bgIdx)
		}
	}
}

func (a *App) pollEyedropper(ix, iy int) {
	if xgal.Click(xgal.MouseButtonLeft) || a.kbdClick {
		a.fgIdx = int(a.getPixel(ix, iy))
		a.setTool(a.prevTool)
	}
	if xgal.Click(xgal.MouseButtonRight) {
		a.bgIdx = int(a.getPixel(ix, iy))
		a.setTool(a.prevTool)
	}
}

func (a *App) selectAll() {
	b := a.doc.Bounds()
	a.sel = b
	a.setMsg("selected all")
}

func (a *App) clearSelection() {
	if a.sel.Empty() {
		return
	}
	b := a.sel
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			a.doc.Pix[y*a.doc.Stride+x] = 0
		}
	}
	a.dirty = true
	a.setMsg("cleared selection")
}

func (a *App) pollPalette(mx, my int) {
	perRow := (WindowW) / PalCell
	if perRow < 1 {
		perRow = 1
	}
	palTop := WindowH - StatusH - a.palH + 2

	for i := 1; i < len(a.doc.Palette); i++ {
		row := i / perRow
		if row >= MaxPalRows {
			break
		}
		colIdx := i % perRow
		px := colIdx * PalCell
		py := palTop + row*PalCell
		rect := xgal.Rect(px, py, px+PalCell, py+PalCell)
		if my >= rect.Min.Y && my < rect.Max.Y && mx >= rect.Min.X && mx < rect.Max.X {
			if xgal.Click(xgal.MouseButtonLeft) {
				a.fgIdx = i
			}
			if xgal.Click(xgal.MouseButtonRight) {
				a.bgIdx = i
			}
		}
	}
}

func (a *App) Draw(screen *xgal.Surface) {
	xgal.Clear(screen, xgal.Wash(60, 60, 60, 255))

	// Draw doc image
	if a.dirty || a.docSurf == nil {
		a.rebuildSurface()
	}
	if a.docSurf != nil {
		xgal.Zoom(screen, a.docSurf, a.offX, a.offY, float64(a.zoom))
	}

	// Pixel grid at high zoom
	if a.zoom >= 4 {
		b := a.doc.Bounds()
		gridCol := xgal.Wash(100, 100, 100, 80)
		for y := 0; y < b.Dy(); y++ {
			sy := int(float64(y)*float64(a.zoom) + a.offY)
			xgal.Box(screen, xgal.Rect(int(a.offX), sy, int(a.offX)+b.Dx()*a.zoom, sy+1), gridCol)
		}
		for x := 0; x < b.Dx(); x++ {
			sx := int(float64(x)*float64(a.zoom) + a.offX)
			xgal.Box(screen, xgal.Rect(sx, int(a.offY), sx+1, int(a.offY)+b.Dy()*a.zoom), gridCol)
		}
	}

	// Selection outline
	if !a.sel.Empty() {
		sx0, sy0 := a.imgToScreen(a.sel.Min.X, a.sel.Min.Y)
		sx1, sy1 := a.imgToScreen(a.sel.Max.X, a.sel.Max.Y)
		col := xgal.Wash(255, 255, 255, 200)
		xgal.Box(screen, xgal.Rect(sx0, sy0, sx1, sy0+1), col)
		xgal.Box(screen, xgal.Rect(sx0, sy1, sx1, sy1+1), col)
		xgal.Box(screen, xgal.Rect(sx0, sy0, sx0+1, sy1), col)
		xgal.Box(screen, xgal.Rect(sx1, sy0, sx1+1, sy1), col)
	}

	// Paste preview
	if a.pasting && a.clip != nil {
		cb := a.clip.Bounds()
		offX := a.pastePt.X - cb.Dx()/2
		offY := a.pastePt.Y - cb.Dy()/2
		for y := 0; y < cb.Dy(); y++ {
			for x := 0; x < cb.Dx(); x++ {
				sx, sy := a.imgToScreen(offX+x, offY+y)
				idx := a.clip.Pix[y*a.clip.Stride+x]
				col := a.clip.Palette[idx]
				if _, _, _, a_ := col.RGBA(); a_ > 0 {
					palColor := xgal.Recolor(col)
					xgal.Box(screen, xgal.Rect(sx, sy, sx+a.zoom, sy+a.zoom), palColor)
				}
			}
		}
		psx, psy := a.imgToScreen(offX, offY)
		peX := psx + cb.Dx()*a.zoom
		peY := psy + cb.Dy()*a.zoom
		outline := xgal.Wash(255, 255, 0, 220)
		xgal.Box(screen, xgal.Rect(psx, psy, peX, psy+1), outline)
		xgal.Box(screen, xgal.Rect(psx, peY, peX, peY+1), outline)
		xgal.Box(screen, xgal.Rect(psx, psy, psx+1, peY), outline)
		xgal.Box(screen, xgal.Rect(peX, psy, peX+1, peY), outline)
	}

	// Drawing preview (select/line/rect)
	if a.drawing {
		switch a.tool {
		case ToolLine:
			sx0, sy0 := a.imgToScreen(a.drawX0, a.drawY0)
			sx1, sy1 := a.imgToScreen(a.drawX1, a.drawY1)
			previewCol := xgal.Wash(255, 255, 255, 180)
			xgal.Box(screen, xgal.Rect(sx0-2, sy0-2, sx0+2, sy0+2), previewCol)
			xgal.Box(screen, xgal.Rect(sx1-2, sy1-2, sx1+2, sy1+2), previewCol)
			xgal.Line(screen, sx0+a.zoom/2, sy0+a.zoom/2, sx1+a.zoom/2, sy1+a.zoom/2, 1, previewCol)
		case ToolSelect:
			sx0, sy0 := a.imgToScreen(a.drawX0, a.drawY0)
			sx1, sy1 := a.imgToScreen(a.drawX1, a.drawY1)
			if sx0 > sx1 {
				sx0, sx1 = sx1, sx0
			}
			if sy0 > sy1 {
				sy0, sy1 = sy1, sy0
			}
			previewCol := xgal.Wash(255, 255, 255, 160)
			xgal.Box(screen, xgal.Rect(sx0, sy0, sx1, sy0+1), previewCol)
			xgal.Box(screen, xgal.Rect(sx0, sy1, sx1, sy1+1), previewCol)
			xgal.Box(screen, xgal.Rect(sx0, sy0, sx0+1, sy1), previewCol)
			xgal.Box(screen, xgal.Rect(sx1, sy0, sx1+1, sy1), previewCol)
		case ToolStroke:
			sx0, sy0 := a.imgToScreen(a.drawX0, a.drawY0)
			sx1, sy1 := a.imgToScreen(a.drawX1, a.drawY1)
			if sx0 > sx1 {
				sx0, sx1 = sx1, sx0
			}
			if sy0 > sy1 {
				sy0, sy1 = sy1, sy0
			}
			previewCol := xgal.Wash(255, 255, 255, 120)
			xgal.Box(screen, xgal.Rect(sx0, sy0, sx1, sy0+1), previewCol)
			xgal.Box(screen, xgal.Rect(sx0, sy1, sx1, sy1+1), previewCol)
			xgal.Box(screen, xgal.Rect(sx0, sy0, sx0+1, sy1), previewCol)
			xgal.Box(screen, xgal.Rect(sx1, sy0, sx1+1, sy1), previewCol)
		case ToolCircle:
			if a.drawing {
				sx0, sy0 := a.imgToScreen(a.drawX0, a.drawY0)
				sx1, sy1 := a.imgToScreen(a.drawX1, a.drawY1)
				dx := sx1 - sx0
				if dx < 0 {
					dx = -dx
				}
				dy := sy1 - sy0
				if dy < 0 {
					dy = -dy
				}
				sr := (dx + dy) / 2
				if sr < 2 {
					sr = 2
				}
				previewCol := xgal.Wash(255, 255, 255, 120)
				xgal.Circle(screen, xgal.Pt(sx0, sy0), sr, 1, previewCol)
			}
		}
	}

	// Pixel cursor for all tools
	if a.mx >= 0 && a.my >= 0 {
		b := a.doc.Bounds()
		if a.mx >= b.Min.X && a.mx < b.Max.X && a.my >= b.Min.Y && a.my < b.Max.Y {
			px := int(float64(a.mx)*float64(a.zoom) + a.offX)
			py := int(float64(a.my)*float64(a.zoom) + a.offY)
			cursorCol := xgal.Recolor(a.doc.Palette[a.fgIdx])
			xgal.Box(screen, xgal.Rect(px, py, px+a.zoom, py+a.zoom), cursorCol)
		}
	}

	// Brush preview circle for pencil/eraser
	if a.mx >= 0 && a.my >= 0 && (a.tool == ToolPencil || a.tool == ToolEraser) {
		b := a.doc.Bounds()
		if a.mx >= b.Min.X && a.mx < b.Max.X && a.my >= b.Min.Y && a.my < b.Max.Y {
			if a.brushSize > 0 || a.zoom >= 2 {
				cx := int(float64(a.mx)*float64(a.zoom)+a.offX) + a.zoom/2
				cy := int(float64(a.my)*float64(a.zoom)+a.offY) + a.zoom/2
				rad := a.brushSize*a.zoom + a.zoom/2
				if rad < a.zoom/2 {
					rad = a.zoom/2 + 1
				}
				col := xgal.Wash(255, 255, 255, 200)
				if a.tool == ToolEraser {
					col = xgal.Wash(255, 100, 100, 200)
				}
				xgal.Circle(screen, xgal.Pt(cx, cy), rad, 1, col)
			}
		}
	}

	// Toolbar background
	xgal.Box(screen, xgal.Rect(0, 0, WindowW, ToolbarH), xgal.Wash(40, 40, 40, 255))
	a.toolGrp.Render(screen)
	a.copyBtn.Render(screen)
	a.pasteBtn.Render(screen)

	// Brush size slider label
	slX := numToggles*toggleW + 2*btnW
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(160, 160, 160, 255), slX, 8, fmt.Sprintf("Sz:%d", a.brushSize))
	a.sizeSl.Render(screen)

	// Palette area
	a.drawPalette(screen)

	// Status bar
	xgal.Box(screen, xgal.Rect(0, WindowH-StatusH, WindowW, WindowH), xgal.Wash(30, 30, 30, 255))
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(200, 200, 200, 255), 4, WindowH-StatusH+4, a.statusText())

	// Ask dialog
	if a.ask != nil {
		a.ask.Render(screen)
	}

	// Help overlay
	if a.showHelp {
		drawHelp(screen)
	}
}

func drawHelp(screen *xgal.Surface) {
	// Dim background
	xgal.Box(screen, xgal.Rect(0, 0, WindowW, WindowH), xgal.Wash(0, 0, 0, 180))

	lines := []string{
		"  xpix – Pixel Art Editor  ",
		"",
		"TOOLS",
		"  F2  Pencil      F3  Select",
		"  F4  Line        F5  Stroke",
		"  F6  Fill        F7  Eyedropper",
		"  F8  Eraser      F9  Circle",
		"  1-8  same tools via digit keys",
		"",
		"ACTIONS",
		"  Ctrl+C  Copy selection",
		"  Ctrl+V  Paste",
		"  Ctrl+X  Cut selection",
		"  Ctrl+A  Select all",
		"  Del     Clear selection",
		"  Ctrl+S  Save",
		"  Ctrl+O  Open",
		"  Ctrl+R  Resize canvas",
		"",
		"VIEW",
		"  +/-     Zoom in/out",
		"  0       Auto-fit zoom",
		"  Q/Esc   Quit",
	}

	face := xgal.BuiltinFace
	lineH := 14
	totalH := len(lines) * lineH
	x0 := 20
	y0 := (WindowH-totalH)/2 - 20
	if y0 < 10 {
		y0 = 10
	}

	xgal.Box(screen, xgal.Rect(x0-10, y0-10, WindowW-x0+10, y0+totalH+10), xgal.Wash(30, 30, 30, 240))
	xgal.Outline(screen, xgal.Rect(x0-10, y0-10, WindowW-x0+10, y0+totalH+10), 1, xgal.Wash(200, 200, 200, 255))

	for i, line := range lines {
		xgal.Ink(screen, face, xgal.Wash(220, 220, 220, 255), x0, y0+i*lineH, line)
	}
}

func (a *App) drawPalette(screen *xgal.Surface) {
	perRow := (WindowW) / PalCell
	if perRow < 1 {
		perRow = 1
	}
	palTop := WindowH - StatusH - a.palH + 2

	xgal.Box(screen, xgal.Rect(0, palTop-2, WindowW, WindowH-StatusH),
		xgal.Wash(45, 45, 45, 255))

	for i := 1; i < len(a.doc.Palette); i++ {
		row := i / perRow
		if row >= MaxPalRows {
			break
		}
		colIdx := i % perRow
		px := colIdx * PalCell
		py := palTop + row*PalCell
		palColor := xgal.Recolor(a.doc.Palette[i])

		xgal.Box(screen, xgal.Rect(px+1, py+1, px+PalCell-1, py+PalCell-1), palColor)
		if i == a.fgIdx {
			xgal.Outline(screen, xgal.Rect(px, py, px+PalCell, py+PalCell), 1, xgal.Wash(255, 255, 255, 255))
		} else if i == a.bgIdx {
			xgal.Outline(screen, xgal.Rect(px, py, px+PalCell, py+PalCell), 1, xgal.Wash(180, 180, 180, 255))
		}

		if i == a.bgIdx && i == a.fgIdx {
			xgal.Box(screen, xgal.Rect(px+1, py+1, px+PalCell-1, py+PalCell/2), xgal.Wash(255, 255, 255, 255))
		}
	}
}

func (a *App) Layout(w, h int) (int, int) {
	return WindowW, WindowH
}

var _ xgal.Game = (*App)(nil)
