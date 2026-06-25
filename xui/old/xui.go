// Package xui is the xmas engine UI package.
// To keep everything relatively simple, there can only be a single active UI.
// However this UI can consist of multiple Controls.
// Only one Control is active at the time.
// Each Controls has an optional set of child Controls.
// Only one child Control per Control is active at one time.
// Each child Control needs to be fully contained in the parent Control
// and may not overflow it.
// Effectively this means the UI is "flat" apart from the Z ordering.
//
// Each Control has a Class that determines its behavior.
// Controls and Classes are separate, but can use embedding
// to extend each other in a double hierarchy.
package xui

import (
	"image"
	"image/color"
	"log/slog"
	"strings"
)

import (
	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

import (
	"github.com/xmasengine/xmas/xres/spleen8"
)

// TextInputField is an input field for IME text entry.
type TextInputField struct {
	textinput.Field
	Point
}

// Keymods are the current key modifers.
type KeyMods struct {
	Alt   bool
	Class bool
	Shift bool
	Meta  bool
}

// Rectangle is used for sizes and positions.
type Rectangle = image.Rectangle

// Point is used for position and offsets.
type Point = image.Point

// Color is a color.
type Color = color.Color

// Image is an image.Image
type Image = image.Image

// Surface is an ebiten.Image
type Surface = ebiten.Image

// RGBA is an RGBA color.
type RGBA = color.RGBA

// Face is a font face
type Face = text.Face

// DrawOptions are options for drawing an image.
type DrawOptions = ebiten.DrawImageOptions

var defaultFontFace = text.NewGoXFace(bitmapfont.Face)

func DrawText(dst *Surface, face Face, color color.RGBA, x, y int, str string) {
	opts := text.DrawOptions{}
	opts.LineSpacing = float64(LineHeight(face))
	opts.GeoM.Translate(float64(x), float64(y))
	opts.ColorScale.Scale(
		float32(color.R)/255.0,
		float32(color.G)/255.0,
		float32(color.B)/255.0,
		float32(color.A)/255.0,
	)
	text.Draw(dst, str, face, &opts)
}

func DrawTextLine(dst *Surface, face Face, color color.RGBA, x, y int, str string) {
	line, _, _ := strings.Cut(str, "\n")
	DrawText(dst, face, color, x, y, line)
}

func MeasureText(txt string, face Face, lineSpacingInPixels float64) (width, height float64) {
	return text.Measure(txt, face, lineSpacingInPixels)
}

func (s Style) MeasureText(txt string) Point {
	w, h := text.Measure(txt, s.Face, float64(LineHeight(s.Face)))
	return image.Pt(int(w), int(h))
}

func (s Style) DrawText(dst *Surface, at Point, txt string) {
	pt := at.Add(s.Margin)
	DrawText(dst, s.Face, s.Writing, pt.X, pt.Y, txt)
}

func (s Style) DrawTextLine(dst *Surface, at Point, txt string) {
	pt := at.Add(s.Margin)
	DrawTextLine(dst, s.Face, s.Writing, pt.X, pt.Y, txt)
}

func LineHeight(face Face) int {
	return int(face.Metrics().HAscent + face.Metrics().HDescent + face.Metrics().HLineGap)
}

func (s Style) LineHeight() int {
	if s.Face != nil {
		return LineHeight(s.Face)
	}
	return 8
}

// Root is the top level of the UI.
type Root struct {
	Control                           // Control root is also a Control
	NoTouchMouse    bool              // NoTouchMouse: set this to true to not translate touches to mouse events.
	TextInputFields []*TextInputField // Text input fields in use
	cx, cy          int
	chars           []rune
	keyMods         KeyMods // Current key KeyMods
	connected       []ebiten.GamepadID
	gamepads        []ebiten.GamepadID
	Hover           *Control     // Hover is the Control that is being hovered by the mouse.
	Focus           *Control     // Focus is the Control that has the input focus.
	Drag            *Control     // Drag is the Control that is being dragged by the mouse or touch.
	Mark            *Control     // Mark is the Control that has the joystick and arrow key marker.
	Default         EventHandler // Default event handler, used if none of the Controls accepts the event.
}

func NewRoot() *Root {
	res := &Root{}
	res.Default = Discard{}
	res.Class = NewRootClass(res)
	return res
}

// State is the state of a Control, or a requested state change.
type State struct {
	Focus bool
	Hover bool
	Pause bool
	Hide  bool
	Clip  bool
	Lock  bool
	Drag  bool
}

// Result is the result of an event handler
type Result struct {
	OK    bool
	State State // Reqquested state of the Control.
}

// A Renderer can render itself.
type Renderer interface {
	// Render renders the Control.
	// The root is passed for convenience, for example to
	// get fonts easily.
	Render(*Root, *Surface)
}

// A Class determines the behavior of a Control. It is a renderer and a listener.
type Class interface {
	Listener
	Renderer
}

// Discard is a handler that does nothing.
type Discard struct{}

func (Discard) HandleEvent(e Event) bool {
	return false // ignore event.
}

// Invisible is a Renderer that does nothing.
type Invisible struct{}

func (Invisible) Render(_ *Root, _ *Surface) {
}

// Control is the basic Control in the UI. Embed this to implement a Control.
// It can be the Root Control, a panel Control or a simple Control.
type Control struct {
	Class    Class     // A Control must embed a Class with the specific behavior.
	Layer    int       // Layer is the Z ordering of the Control.
	Bounds   Rectangle // Actual position and size of the Control.
	Extent   Rectangle // Extent is a rectangle with the desired size of the Control. It may be larger than the Bounds.
	Offset   Point     // Offset for scrolling
	Style    Style
	State    State
	Controls []*Control // Sub Controls of the Control if any.
}

// ControlClass is the basic class for a Control. Embed this to implement a class.
type ControlClass struct {
	BasicListener
}

func (w ControlClass) Render(r *Root, screen *Surface) {
	// draw nothing
}

func NewControlClass() *ControlClass {
	return &ControlClass{}
}

// Depth first traversal. Stops if non-nil is returned.
func (w *Control) EachControl(cb func(sub *Control) *Control) *Control {
	for i := len(w.Controls) - 1; i >= 0; i-- {
		sub := w.Controls[i]
		if res := sub.EachControl(cb); res != nil {
			return sub
		}
	}
	if res := cb(w); res != nil {
		return res
	}
	return nil
}

func (w *Control) FindTop(at Point) *Control {
	var top *Control
	for i := len(w.Controls) - 1; i >= 0; i-- {
		p := w.Controls[i]
		if !p.State.Hide {
			subTop := p.FindTop(at)
			if top == nil && subTop != nil {
				top = subTop
			} else if subTop != nil && top.Layer < subTop.Layer {
				top = subTop
			}
			if at.In(p.Bounds) {
				if top == nil {
					top = p
				} else if top.Layer < p.Layer {
					top = p
				}
			}

		}
	}
	return top
}

func (w *Control) Append(Controls ...*Control) {
	w.Controls = append(w.Controls, Controls...)
}

func (w *Control) Move(delta Point) *Control {
	w.Bounds = w.Bounds.Add(delta)
	w.MoveControls(delta)
	return w
}

func (w *Control) MoveControls(delta Point) *Control {
	for i := 0; i < len(w.Controls); i++ {
		sub := w.Controls[i]
		if sub.State.Lock {
			continue
		}
		sub.Move(delta)
	}
	return w
}

func (w *Control) MoveAll(delta Point) *Control {
	w.Bounds = w.Bounds.Add(delta)
	w.MoveAllControls(delta)
	return w
}

func (w *Control) MoveAllControls(delta Point) *Control {
	for i := 0; i < len(w.Controls); i++ {
		sub := w.Controls[i]
		sub.MoveAll(delta)
	}
	return w
}

// RenderControl renders the Controls inside this Control, not the Control itself.
func (w *Control) RenderControls(r *Root, screen *Surface) *Control {
	for i := 0; i < len(w.Controls); i++ {
		sub := w.Controls[i]
		sub.Class.Render(r, screen)
	}
	return w
}

const ControlScrollSlack = 2

func (p *Control) ScrollHorizontal(pos, low, high int) {
	if low == high {
		return
	}

	scrollRange := p.Extent.Dx() - ControlScrollSlack
	var noff Point
	noff.X = ((pos - low) * scrollRange) / (high - low)

	delta := p.Offset.Sub(noff)
	p.MoveControls(delta)
	p.Offset = noff
}

func (p *Control) ScrollVertical(pos, low, high int) {
	if low == high {
		return
	}

	scrollRange := p.Extent.Dy() - ControlScrollSlack
	var noff Point
	noff.Y = ((pos - low) * scrollRange) / (high - low)
	delta := p.Offset.Sub(noff)
	p.MoveControls(delta)
	p.Offset = noff
}

func (w *Control) AddControl(sub *Control) {
	w.Controls = append(w.Controls, sub)
}

func NewControl() *Control {
	res := &Control{}
	res.Class = NewControlClass()
	return res
}

func (r *Root) On(ev Eventer) bool {
	e := ev.Event()
	slog.Debug("Root.On ", "event", e)
	return e.Dispatch(r.Class)
}

func (r *Root) HandleEvent(e Event) bool {
	if r.Default != nil {
		return r.Default.HandleEvent(e)
	}
	dprintln("warning: Root.HandleEvent, event not handled: ")
	return false
}

// Update is called 60 times per second.
// Input should be checked during this function.
func (r *Root) Update() error {
	for _, gid := range r.gamepads {
		if inpututil.IsGamepadJustDisconnected(gid) {
			r.On(MakePadEvent(PadDetach, r, int(gid), 0, 0, nil))
		}
	}

	r.connected = inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, gid := range r.connected {
		r.On(MakePadEvent(PadAttach, r, int(gid), 0, 0, nil))
	}

	r.gamepads = r.gamepads[0:0]
	r.gamepads = ebiten.AppendGamepadIDs(r.gamepads)
	for _, gid := range r.gamepads {
		buttons := inpututil.AppendJustPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.On(MakePadEvent(PadPress, r, int(gid), int(button), 0, nil))
		}

		buttons = inpututil.AppendPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			dur := inpututil.GamepadButtonPressDuration(gid, button)
			r.On(MakePadEvent(PadHold, r, int(gid), int(button), dur, nil))
		}

		buttons = inpututil.AppendJustReleasedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.On(MakePadEvent(PadRelease, r, int(gid), int(button), 0, nil))
		}

		count := ebiten.GamepadAxisCount(gid)
		axes := make([]float64, count)
		moved := false
		for axis := 0; axis < count; axis++ {
			value := ebiten.GamepadAxisValue(gid, axis)
			axes[axis] = value
			moved = moved || ((value > 0.1) || (value < -0.1))
		}
		if (len(axes) > 0) && moved {
			r.On(MakePadEvent(PadMove, r, int(gid), 0, 0, axes))
		}
	}

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		r.On(MakeKeyEvent(KeyPress, r, -1, int(key), 0))
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		dur := inpututil.KeyPressDuration(key)
		r.On(MakeKeyEvent(KeyHold, r, -1, int(key), dur))
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		r.On(MakeKeyEvent(KeyRelease, r, -1, int(key), 0))
	}

	if len(r.chars) == 0 && cap(r.chars) == 0 {
		r.chars = make([]rune, 0, 32)
	} else {
		r.chars = r.chars[0:0]
	}

	r.chars = ebiten.AppendInputChars(r.chars)
	if len(r.chars) > 0 {
		slog.Debug("input chars", "chars", r.chars)
		r.On(MakeKeyEvent(KeyText, r, -1, 0, 0, r.chars...))
	}
	r.chars = r.chars[0:0]

	for id, field := range r.TextInputFields {
		if field.IsFocused() {
			handled, _ := field.HandleInput(field.X, field.Y)
			if handled {
				r.On(MakeKeyEvent(KeyText, r, id, 0, 0, []rune(field.Text())...))
			}
		}
	}

	touches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.On(MakeTouchEvent(TouchPress, r, int(touch), image.Pt(x, y), image.Point{}, 0))
	}

	touches = ebiten.AppendTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		px, py := inpututil.TouchPositionInPreviousTick(touch)
		dx, dy := x-px, y-py
		dur := inpututil.TouchPressDuration(touch)
		r.On(MakeTouchEvent(TouchHold, r, int(touch), image.Pt(x, y), image.Pt(dx, dy), dur))
	}

	touches = inpututil.AppendJustReleasedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.On(MakeTouchEvent(TouchRelease, r, int(touch), image.Pt(x, y), image.Point{}, 0))
	}

	x, y := ebiten.CursorPosition()
	dx, dy := x-r.cx, y-r.cy
	at := image.Pt(x, y)
	delta := image.Pt(dx, dy)

	for mb := ebiten.MouseButton(0); mb < ebiten.MouseButtonMax; mb++ {
		if inpututil.IsMouseButtonJustPressed(mb) {
			r.On(MakeMouseEvent(MousePress, r, int(mb), at, delta, 0))
		}
		if ebiten.IsMouseButtonPressed(mb) {
			dur := inpututil.MouseButtonPressDuration(mb)
			r.On(MakeMouseEvent(MouseHold, r, int(mb), at, delta, dur))
		}
		if inpututil.IsMouseButtonJustReleased(mb) {
			r.On(MakeMouseEvent(MouseRelease, r, int(mb), at, delta, 0))
		}
	}
	if dx != 0 || dy != 0 {
		r.On(MakeMouseEvent(MouseMove, r, -1, at, delta, 0))
	}
	r.cx = x
	r.cy = y

	wx, wy := ebiten.Wheel()
	if wx != 0 || wy != 0 {
		wheel := image.Pt(int(wx), int(wy))
		r.On(MakeMouseWheelEvent(MouseWheel, r, at, delta, wheel))
	}

	return nil
}

type RootClass struct {
	*Root
	*ControlClass
}

func NewRootClass(r *Root) *RootClass {
	res := &RootClass{Root: r}
	res.ControlClass = NewControlClass()
	return res
}

func (r *Root) SetFocus(w *Control, at Point) *Control {
	old := r.Focus
	if r.Focus != nil && r.Focus != w {
		MakeActionEvent(ActionBlur, r, at, image.Point{}).Dispatch(r.Focus.Class)
	}

	if r.Focus != w {
		r.Focus = w
		if r.Focus != nil {
			MakeActionEvent(ActionFocus, r, at, image.Point{}).Dispatch(r.Focus.Class)
		}
	}
	return old
}

func (r *Root) SetHover(w *Control, at Point) *Control {
	old := r.Hover

	if r.Hover != nil && r.Hover != w {
		MakeActionEvent(ActionCrash, r, at, image.Point{}).Dispatch(r.Hover.Class)
	}

	r.Hover = w
	if r.Hover != nil {
		MakeActionEvent(ActionHover, r, at, image.Point{}).Dispatch(r.Hover.Class)
	}
	return old
}

func (r *Root) SetDrag(w *Control, at, delta Point) *Control {
	old := r.Drag

	if r.Drag != nil && r.Drag != w {
		MakeActionEvent(ActionDrop, r, at, delta).Dispatch(r.Drag.Class)
	}

	r.Drag = w
	if r.Drag != nil {
		MakeActionEvent(ActionDrag, r, at, delta).Dispatch(r.Drag.Class)
	}
	return old
}

func (r *RootClass) OnMouseMove(e MouseEvent) bool {
	w := r.Root
	hover := w.FindTop(e.At)

	if w.Drag != nil && w.Drag == hover {
		MakeActionEvent(ActionDrag, e.Root(), e.At, e.Delta).Dispatch(w.Drag.Class)
	}

	if w.Hover != nil && w.Hover != hover {
		MakeActionEvent(ActionCrash, e.Root(), e.At, image.Point{}).Dispatch(w.Hover.Class)
	}

	w.Hover = hover
	if w.Hover != nil {
		MakeActionEvent(ActionHover, e.Root(), e.At, image.Point{}).Dispatch(w.Hover.Class)
	}
	return true
}

func (r *RootClass) OnMousePress(e MouseEvent) bool {
	w := r.Root
	top := w.FindTop(e.At)
	r.SetFocus(top, e.At)
	r.SetDrag(top, e.At, e.Delta)

	if r.Hover != nil {
		e.Dispatch(r.Hover.Class)
	}

	if r.Focus != nil && r.Focus != r.Hover {
		e.Dispatch(r.Focus.Class)
	}
	return false
}

func (r *RootClass) OnMouseRelease(e MouseEvent) bool {
	if r.Hover != nil {
		e.Dispatch(r.Hover.Class)
	}

	if r.Focus != nil && r.Focus != r.Hover {
		e.Dispatch(r.Focus.Class)
	}

	r.SetDrag(nil, e.At, e.Delta)
	return false
}

func (r *RootClass) OnMouseWheel(e MouseEvent) bool {
	if r.Hover != nil {
		e.Dispatch(r.Hover.Class)
	}

	if r.Focus != nil && r.Focus != r.Hover {
		e.Dispatch(r.Focus.Class)
	}
	return false
}

func (r *RootClass) OnKeyPress(e KeyEvent) bool {
	w := r.Root
	if w.Focus != nil {
		slog.Debug("RootClass.OnKeyPress", "e", e)
		return e.Dispatch(w.Focus.Class)
	}
	return false
}

func (r *RootClass) OnKeyHold(e KeyEvent) bool {
	w := r.Root
	if w.Focus != nil {
		return e.Dispatch(w.Focus.Class)
	}
	return false
}

func (r *RootClass) OnKeyText(e KeyEvent) bool {
	w := r.Root
	if w.Focus != nil {
		return e.Dispatch(w.Focus.Class)
	}
	return false
}

func (r *RootClass) Render(_ *Root, screen *Surface) {
	for _, p := range r.Root.Controls {
		if !p.State.Hide {
			p.Class.Render(r.Root, screen)
		}
	}
}

// Draw is called when the UI needs to be drawn in game.ui
func (r *Root) Draw(screen *Surface) {
	r.Class.Render(r, screen)
}

// Layout is called when the contents of the element need to be laid out.
// The element should accept that the available size is less than its
// real size and draw it appropiately, such as scrolling.
// The returned elementWidth and elementHeight must be smaller than
// or equal to the available width.
func (r *Root) Layout(availableWidth, availableHeight int) (elementWidth, elementHeight int) {

	return availableWidth, availableHeight
}

func dprintln(msg string, vars ...any) {
	slog.Info(msg, "vars", vars)
}
