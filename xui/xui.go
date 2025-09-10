// Package xui is the xmas engine UI package.
// To keep everything relatively simple, there can only be a single active UI.
// However this UI can consist of multiple Widgets.
// Only one Widget is active at the time.
// Each Widgets has an optional set of child widgets.
// Only one child widget per Widget is active at one time.
// Each child widget needs to be fully contained in the parent Widget
// and may not overflow it.
// Effectively this means the UI is "flat" apart from the Z ordering.
//
// Each Widget has a Class that determines its behavior.
// Widgets and Classes are separate, but can use embedding
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

// Style is the style of a Widget.
type Style struct {
	Fore    RGBA
	Border  RGBA
	Shadow  RGBA
	Fill    RGBA
	Writing RGBA
	Margin  Point
	Stroke  int
	Face    Face
}

var defaultFontFace = text.NewGoXFace(bitmapfont.Face)

func DrawTextLine(dst *Surface, face Face, color color.RGBA, x, y int, str string) {
	opts := text.DrawOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	opts.ColorScale.Scale(
		float32(color.R)/255.0,
		float32(color.G)/255.0,
		float32(color.B)/255.0,
		float32(color.A)/255.0,
	)
	text.Draw(dst, str, face, &opts)
}

func DrawText(dst *Surface, face Face, color color.RGBA, x, y int, str string) {
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		DrawTextLine(dst, face, color, x, y, line)
		y += LineHeight(face)
	}
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

func LineHeight(face Face) int {
	return int(face.Metrics().HAscent + face.Metrics().HDescent + face.Metrics().HLineGap)
}

// Root is the top level of the UI.
type Root struct {
	Widget                            // Widget root is also a widget
	NoTouchMouse    bool              // NoTouchMouse: set this to true to not translate touches to mouse events.
	TextInputFields []*TextInputField // Text input fields in use
	cx, cy          int
	chars           []rune
	keyMods         KeyMods // Current key KeyMods
	connected       []ebiten.GamepadID
	gamepads        []ebiten.GamepadID
	Focus           *Widget      // Focus is the Widget that has the input focus.
	Drag            *Widget      // Drag is the Widget that is being dragged by the mouse or touch.
	Mark            *Widget      // Mark is the Widget that has the joystick and arrow key marker.
	Default         EventHandler // Default event handler, used if none of the Widgets accepts the event.
}

func NewRoot() *Root {
	res := &Root{}
	res.Default = Discard{}
	res.Class = NewRootClass(res)
	return res
}

// State is the state of a Widget, or a requested state change.
type State struct {
	Focus bool
	Hover bool
	Pause bool
	Hide  bool
	Clip  bool
}

// Result is the result of an event handler
type Result struct {
	OK    bool
	State State // Reqquested state of the Widget.
}

// A Renderer can render itself.
type Renderer interface {
	// Render renders the Widget.
	// The root is passed for convenience, for example to
	// get fonts easily.
	Render(*Root, *Surface)
}

// A Class determines the behavior of a widget. It is a renderer and a listener.
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

// Widget is the basic widget in the UI. Embed this to implement a widget.
// It can be the Root widget, a panel widget or a simple widget.
type Widget struct {
	Class   Class     // A widget must embed a Class with the specific behavior.
	Layer   int       // Layer is the Z ordering of the widget.
	Bounds  Rectangle // Actual position and size of the widget.
	Size    Rectangle // Size is the desired size of the widget, may be bigger than Bounds.
	Style   Style
	State   State
	Widgets []*Widget // Sub widgets of the widget if any.
	Hover   *Widget   // Hover is the Widget that is being hovered by the mouse.
	Focus   *Widget   // Hover is the Widget that is being focused.
}

// WidgetClass is the basic class for a Widget. Embed this to implement a class.
type WidgetClass struct {
	BasicListener
}

func (w WidgetClass) Render(r *Root, screen *Surface) {
	// draw nothing
}

func NewWidgetClass() *WidgetClass {
	return &WidgetClass{}
}

func (w *Widget) FindTop(at Point) *Widget {
	var top *Widget
	for i := len(w.Widgets) - 1; i >= 0; i-- {
		p := w.Widgets[i]
		if at.In(p.Bounds) {
			if top == nil {
				top = p
			} else if top.Layer < p.Layer {
				top = p
			}
		}
	}
	if top != nil {
		sub := top.FindTop(at)
		if sub != nil {
			return sub
		}
	}
	return top
}

func (w *Widget) Append(widgets ...*Widget) {
	w.Widgets = append(w.Widgets, widgets...)
}

func NewWidget() *Widget {
	res := &Widget{}
	res.Class = NewWidgetClass()
	return res
}

func (r *Root) On(e Event) bool {
	slog.Debug("Root.On ", "event", e)
	return e.Dispatch(r.Class)
}

func (r *Root) HandleEvent(e Event) bool {
	if r.Default != nil {
		return r.Default.HandleEvent(e)
	}
	println("warning: Root.HandleEvent, event not handled: ")
	return false
}

// Update is called 60 times per second.
// Input should be checked during this function.
func (r *Root) Update() error {
	for _, gid := range r.gamepads {
		if inpututil.IsGamepadJustDisconnected(gid) {
			r.On(Event{Msg: PadDetach, Pad: MakePadEvent(r, int(gid), 0, 0, nil)})
		}
	}

	r.connected = inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, gid := range r.connected {
		r.On(Event{Msg: PadAttach, Pad: MakePadEvent(r, int(gid), 0, 0, nil)})
	}

	r.gamepads = r.gamepads[0:0]
	r.gamepads = ebiten.AppendGamepadIDs(r.gamepads)
	for _, gid := range r.gamepads {
		buttons := inpututil.AppendJustPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.On(Event{Msg: PadPress, Pad: MakePadEvent(r, int(gid), int(button), 0, nil)})
		}

		buttons = inpututil.AppendPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			dur := inpututil.GamepadButtonPressDuration(gid, button)
			r.On(Event{Msg: PadHold, Pad: MakePadEvent(r, int(gid), int(button), dur, nil)})
		}

		buttons = inpututil.AppendJustReleasedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.On(Event{Msg: PadRelease, Pad: MakePadEvent(r, int(gid), int(button), 0, nil)})
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
			r.On(Event{Msg: PadMove, Pad: MakePadEvent(r, int(gid), 0, 0, axes)})
		}
	}

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		r.On(Event{Msg: KeyPress, Key: MakeKeyEvent(r, -1, int(key), 0, "")})
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		dur := inpututil.KeyPressDuration(key)
		r.On(Event{Msg: KeyHold, Key: MakeKeyEvent(r, -1, int(key), dur, "")})
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		r.On(Event{Msg: KeyRelease, Key: MakeKeyEvent(r, -1, int(key), 0, "")})
	}

	if len(r.chars) == 0 && cap(r.chars) == 0 {
		r.chars = make([]rune, 0, 32)
	} else {
		r.chars = r.chars[0:0]
	}

	r.chars = ebiten.AppendInputChars(r.chars)
	if len(r.chars) > 0 {
		r.On(Event{Msg: KeyText, Key: MakeKeyEvent(r, -1, 0, 0, string(r.chars))})
	}
	r.chars = r.chars[0:0]

	for id, field := range r.TextInputFields {
		if field.IsFocused() {
			handled, _ := field.HandleInput(field.X, field.Y)
			if handled {
				r.On(Event{Msg: KeyText, Key: MakeKeyEvent(r, id, 0, 0, field.Text())})
			}
		}
	}

	touches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.On(Event{Msg: TouchPress, Touch: MakeTouchEvent(r, int(touch), image.Pt(x, y), image.Point{}, 0)})
	}

	touches = ebiten.AppendTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		px, py := inpututil.TouchPositionInPreviousTick(touch)
		dx, dy := x-px, y-py
		dur := inpututil.TouchPressDuration(touch)
		r.On(Event{Msg: TouchHold, Touch: MakeTouchEvent(r, int(touch), image.Pt(x, y), image.Pt(dx, dy), dur)})
	}

	touches = inpututil.AppendJustReleasedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.On(Event{Msg: TouchRelease, Touch: MakeTouchEvent(r, int(touch), image.Pt(x, y), image.Point{}, 0)})
	}

	x, y := ebiten.CursorPosition()
	dx, dy := x-r.cx, y-r.cy
	at := image.Pt(x, y)
	delta := image.Pt(dx, dy)

	for mb := ebiten.MouseButton(0); mb < ebiten.MouseButtonMax; mb++ {
		if inpututil.IsMouseButtonJustPressed(mb) {
			r.On(Event{Msg: MousePress, Mouse: MakeMouseEvent(r, at, delta, 0)})
		}
		if ebiten.IsMouseButtonPressed(mb) {
			dur := inpututil.MouseButtonPressDuration(mb)
			r.On(Event{Msg: MouseHold, Mouse: MakeMouseEvent(r, at, delta, dur)})
		}
		if inpututil.IsMouseButtonJustReleased(mb) {
			r.On(Event{Msg: MouseRelease, Mouse: MakeMouseEvent(r, at, delta, 0)})
		}
	}
	if dx != 0 || dy != 0 {
		r.On(Event{Msg: MouseMove, Mouse: MakeMouseEvent(r, at, delta, 0)})
	}
	r.cx = x
	r.cy = y

	wx, wy := ebiten.Wheel()
	if wx != 0 || wy != 0 {
		wheel := image.Pt(int(wx), int(wy))
		r.On(Event{Msg: MouseWheel, Mouse: MakeMouseWheelEvent(r, at, delta, wheel)})
	}

	return nil
}

type RootClass struct {
	*Root
	*WidgetClass
}

func NewRootClass(r *Root) *RootClass {
	res := &RootClass{Root: r}
	res.WidgetClass = NewWidgetClass()
	return res
}

func (r *RootClass) OnMouseMove(e MouseEvent) bool {
	w := r.Root
	hover := w.FindTop(e.At)

	if w.Hover != nil && w.Hover != hover {
		Event{Msg: ActionCrash, Action: MakeActionEvent(e.Root(), e.At, image.Point{})}.Dispatch(w.Hover.Class)
	}

	w.Hover = hover
	if w.Hover != nil {
		return Event{Msg: ActionHover, Action: MakeActionEvent(e.Root(), e.At, image.Point{})}.Dispatch(w.Hover.Class)
	}
	return false
}

func (r *RootClass) OnMousePress(e MouseEvent) bool {
	w := r.Root
	top := w.FindTop(e.At)

	if w.Focus != nil && w.Focus != top {
		Event{Msg: ActionBlur, Action: MakeActionEvent(e.Root(), e.At, image.Point{})}.Dispatch(w.Focus.Class)
	}

	if w.Focus != top {
		w.Focus = top
		Event{Msg: ActionFocus, Action: MakeActionEvent(e.Root(), e.At, image.Point{})}.Dispatch(w.Focus.Class)
	}

	if w.Focus != nil {
		return Event{Msg: MousePress, Mouse: e}.Dispatch(w.Focus.Class)
	}
	return false
}

func (r *RootClass) OnMouseRelease(e MouseEvent) bool {
	w := r.Root
	if w.Focus != nil {
		return Event{Msg: MouseRelease, Mouse: e}.Dispatch(w.Focus.Class)
	}
	return false
}

func (r *RootClass) Render(_ *Root, screen *Surface) {
	for _, p := range r.Root.Widgets {
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

func DefaultStyle() Style {
	s := Style{}
	s.Border = color.RGBA{50, 50, 50, 245}
	s.Writing = color.RGBA{245, 245, 245, 245}
	s.Shadow = color.RGBA{15, 15, 15, 191}
	s.Fill = color.RGBA{0, 0, 245, 245}
	s.Face = defaultFontFace
	s.Stroke = 1
	s.Margin = image.Pt(2, 2)
	return s
}

func FocusStyle() Style {
	s := DefaultStyle()
	s.Border = color.RGBA{240, 240, 240, 245}
	s.Writing = color.RGBA{245, 245, 245, 245}
	s.Fill = color.RGBA{128, 128, 245, 245}
	return s
}

func HoverStyle() Style {
	s := DefaultStyle()
	s.Border = color.RGBA{240, 240, 50, 250}
	return s
}

func PressStyle() Style {
	s := DefaultStyle()
	s.Fill = color.RGBA{0, 45, 245, 245}
	return s
}

func FillRect(Surface *Surface, r Rectangle, col color.RGBA) {
	vector.DrawFilledRect(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
		col, false,
	)
}

func DrawRect(Surface *Surface, r Rectangle, thick int, col color.RGBA) {
	vector.StrokeRect(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
		float32(thick), col, false,
	)
}

// DrawsLine draws a line on the diagonal of the Rectangle r.
func DrawLine(Surface *Surface, r Rectangle, thick int, col color.RGBA) {
	vector.StrokeLine(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Max.X), float32(r.Max.Y),
		float32(thick), col, false,
	)
}

func (s Style) DrawRect(Surface *Surface, r Rectangle) {
	DrawRect(Surface, r, int(s.Stroke), s.Border)
}

func (s Style) DrawBox(Surface *Surface, r Rectangle) {
	if s.Shadow.A != 0 {
		shadow := s.Shadow
		shadow.A = (shadow.A / 2) + 1 // make half transparent
		right := image.Rect(r.Max.X+1, r.Min.Y+1, r.Max.X+1, r.Max.Y+1)
		DrawLine(Surface, right, 1, shadow)
		bottom := image.Rect(r.Min.X+1, r.Max.Y+1, r.Max.X+1, r.Max.Y+1)
		DrawLine(Surface, bottom, 1, shadow)
	}

	vector.DrawFilledRect(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()), s.Fill, false,
	)

	if s.Stroke > 0 {
		vector.StrokeRect(
			Surface, float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()),
			float32(s.Stroke), s.Border, false,
		)
	}
}

func (s Style) DrawCircleInBox(Surface *Surface, box Rectangle) {
	r := box.Dx()
	if box.Dy() < r {
		r = box.Dy()
	}
	r = r / 2
	c := image.Pt((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2)
	s.DrawCircle(Surface, c, r)
}

func (s Style) DrawCircle(Surface *Surface, c Point, r int) {
	if r < 0 {
		r = 1
	}
	vector.DrawFilledCircle(Surface, float32(c.X), float32(c.Y),
		float32(r), s.Fill, false)

	if s.Stroke > 0 {
		vector.StrokeCircle(
			Surface, float32(c.X), float32(c.Y),
			float32(r), float32(s.Stroke), s.Border, false,
		)
	}
}

func (w *Widget) AddBox(bounds Rectangle) *Box {
	box := NewBox(bounds)
	w.Widgets = append(w.Widgets, &box.Widget)
	return box
}

func NewBox(bounds Rectangle) *Box {
	box := &Box{}
	box.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
	box.Class = NewBoxClass(box)
	return box
}

type Box struct {
	Widget
}

type BoxClass struct {
	*Box
	*WidgetClass
}

func NewBoxClass(b *Box) *BoxClass {
	res := &BoxClass{Box: b}
	res.WidgetClass = NewWidgetClass()
	return res
}

// Render is called when the element needs to be drawn.
func (bc BoxClass) Render(r *Root, screen *Surface) {
	b := bc.Box
	style := b.Style
	if b.State.Hover {
		style = HoverStyle()
	}
	style.DrawBox(screen, b.Bounds)
	for _, w := range b.Widgets {
		if !w.State.Hide {
			w.Class.Render(r, screen)
		}
	}
}

func (b *BoxClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *BoxClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}

type LabelClass struct {
	*Label
	*WidgetClass
}

func NewLabelClass(b *Label) *LabelClass {
	res := &LabelClass{Label: b}
	res.WidgetClass = NewWidgetClass()
	return res
}

type Label struct {
	Widget
	Text    string
	pressed bool
}

func (b LabelClass) Render(r *Root, screen *Surface) {
	box := b.Bounds
	style := b.Style

	if b.State.Hover {
		style = HoverStyle()
	}

	at := box.Min

	style.DrawBox(screen, box)
	style.DrawText(screen, at, b.Text)
}

func (b *LabelClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *LabelClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}

func (l *Label) SetText(text string) {
	l.Text = text
}

func (p *Widget) AddLabel(bounds Rectangle, text string) *Label {
	b := NewLabel(bounds, text)
	p.Widgets = append(p.Widgets, &b.Widget)
	return b
}

func NewLabel(bounds Rectangle, text string) *Label {
	res := &Label{Text: text}
	res.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
	res.Class = NewLabelClass(res)
	return res
}

type Button struct {
	Widget
	Text    string
	Clicked func(*Button)
	pressed bool
	Result  int // May be set freely except on dialog Buttons.
}

type ButtonClass struct {
	*Button
	*WidgetClass
}

func NewButtonClass(b *Button) *ButtonClass {
	res := &ButtonClass{Button: b}
	res.WidgetClass = NewWidgetClass()
	return res
}

func (b ButtonClass) Render(r *Root, screen *Surface) {
	box := b.Bounds
	style := b.Style

	if b.pressed {
		box = box.Add(b.Style.Margin)
		style = PressStyle()
	} else if b.State.Hover {
		style = HoverStyle()
	}

	at := box.Min

	style.DrawBox(screen, box)
	style.DrawText(screen, at, b.Text)
}

func (b *ButtonClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *ButtonClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}

func (b *ButtonClass) OnMousePress(e MouseEvent) bool {
	b.pressed = true
	return true
}

func (b *ButtonClass) OnMouseRelease(e MouseEvent) bool {
	b.pressed = false
	if b.Clicked != nil {
		b.Clicked(b.Button)
	}
	return true
}

func (b *Button) SetText(text string) {
	b.Text = text
}

func NewButton(bounds Rectangle, text string, cl func(*Button)) *Button {
	b := &Button{Text: text, Clicked: cl}
	b.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
	b.Class = NewButtonClass(b)
	return b
}

func (p *Widget) AddButton(bounds Rectangle, text string, cl func(*Button)) *Button {
	b := NewButton(bounds, text, cl)
	p.Widgets = append(p.Widgets, &b.Widget)
	return b
}
