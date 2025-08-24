// Package xui is the xmas engine UI package.
// To keep everything relatively simple, there can only be a single active UI
// However this UI can consist of multiple panels.
// Only one panel is active at the time.
// Each panels has a set of widgets.
// Only one widget per panel is active at one time.
// Widget cannot contain any sub widgets.
// Each widget needs to be fully contained in the panel and may not overflow it.
// Effectively this means the UI is "flat" apart from the panels which have
// a Z ordering.
package xui

import (
	"image"
	"image/color"
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
	Alt     bool
	Control bool
	Shift   bool
	Meta    bool
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

// Style is the style of a Widget or Panel.
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
	Panels          []*Panel          // Panels of the UI.
	NoTouchMouse    bool              // NoTouchMouse: set this to true to not translate touches to mouse events.
	TextInputFields []*TextInputField // Text input fields in use
	cx, cy          int
	chars           []rune
	keyMods         KeyMods // Current key KeyMods
	connected       []ebiten.GamepadID
	gamepads        []ebiten.GamepadID
	Focus           *Panel  // Focus is the Panel that has the input focus.
	Hover           *Panel  // Hover is the Panel that is being hovered by the mouse.
	Drag            *Panel  // Drag is the panel that is being dragged by the mouse or touch.
	Mark            *Panel  // Mark is the panel that has the joystick and arrow key marker.
	Default         Handler // Default event handler, used if of the panels accepts the event.
}

func NewRoot() *Root {
	res := &Root{}
	res.Default = Discard{}
	return res
}

// Message is a kind of message that is sent to the UI components.
type Message int

const (
	NoMessage Message = iota
	PadDetach
	PadAttach
	PadPress
	PadHold
	PadRelease
	PadMove
	KeyPress
	KeyHold
	KeyRelease
	KeyText
	TouchPress
	TouchHold
	TouchRelease
	MousePress
	MouseRelease
	MouseHold
	MouseMove
	MouseWheel
	ActionFocus
	ActionBlur
	ActionHover
	ActionCrash
	ActionDrag
	ActionDrop
	ActionMark
	ActionClean
	LayoutGet
	LayoutSet
	LastMessage
)

// Event in an event that is sent to the UI components.
type Event struct {
	Msg      Message
	ID       int
	Button   int
	Duration int
	Axes     []float64
	Code     int
	Mods     KeyMods
	Chars    string
	At       Point
	Delta    Point
	Wheel    Point
	Bounds   Rectangle
}

// State is the state of a Widget or Panel, or a requested state change.
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
	State State // Reqquested state of the widget or panel.
}

// Handler can handle events events.
type Handler interface {
	// Handle should handle the event and return the result.
	// The root is passed for convenience, for example to
	// manipulate other panels easily.
	// A widget or panel will only receive events that it is intent to handle,
	// but it will not receive any mouse clicks ort ouches outside of its
	// bounds, unless if it is an active panel being dragged.
	Handle(*Root, Event) bool
}

// A Renderer can render itself.
type Renderer interface {
	// Render renders the widget or panel.
	// The root is passed for convenience, for example to
	// get fonts easily.
	Render(*Root, *Surface)
}

// A control is a renderer and a handler
type Control interface {
	Handler
	Renderer
}

type defaultHandler struct {
	norm Handler
	def  Handler
}

func (d defaultHandler) Handle(r *Root, e Event) bool {
	var res bool
	if d.norm != nil {
		res = d.norm.Handle(r, e)
	}
	if !res {
		if d.def != nil {
			res = d.def.Handle(r, e)
		}
	}
	return res
}

// HandleDefault tries to call normal first, and then def if normal
// returns false
func HandleDefault(norm, def Handler) Handler {
	return defaultHandler{norm: norm, def: def}
}

// Discard is a handler that does nothing.
type Discard struct{}

func (Discard) Handle(_ *Root, e Event) bool {
	return false // ignore event.
}

// Invisible is a Renderer that does nothing.
type Invisible struct{}

func (Invisible) Render(_ *Root, _ *Surface) {
}

// Mapper is a handler that maps events to individual handlers.
type Mapper struct {
	Handlers [LastMessage]func(*Root, Event) bool
	Renderer func(*Root, *Surface)
}

func (m Mapper) Handle(r *Root, e Event) bool {
	if e.Msg <= NoMessage {
		return false
	}
	if e.Msg >= LastMessage {
		return false
	}
	handler := m.Handlers[e.Msg]
	if handler == nil {
		return false
	}

	return handler(r, e)
}

func (m Mapper) Render(r *Root, s *Surface) {
	if m.Renderer != nil {
		m.Renderer(r, s)
	}
}

func NewMapper(r func(*Root, *Surface)) *Mapper {
	m := &Mapper{}
	m.Renderer = r
	return m
}

func (m *Mapper) Add(e Message, h func(*Root, Event) bool) *Mapper {
	if e <= NoMessage {
		panic("Mapper.Add Message out of range")
	}
	if e >= LastMessage {
		panic("Mapper.Add Message out of range")
	}

	m.Handlers[e] = h
	return m
}

type Dispatcher struct {
	Target any
}

func (d Dispatcher) Handle(r *Root, e Event) bool {
	return Dispatch(d.Target, r, e)
}

func Dispatch(d any, r *Root, e Event) bool {
	switch e.Msg {
	case PadDetach:
		if impl, ok := d.(interface{ OnPadDetach(r *Root, e Event) bool }); ok {
			return impl.OnPadDetach(r, e)
		}
	case PadAttach:
		if impl, ok := d.(interface{ OnPadAttach(r *Root, e Event) bool }); ok {
			impl.OnPadAttach(r, e)
		}
	case PadPress:
		if impl, ok := d.(interface{ OnPadPress(r *Root, e Event) bool }); ok {
			impl.OnPadPress(r, e)
		}
	case PadHold:
		if impl, ok := d.(interface{ OnPadHold(r *Root, e Event) bool }); ok {
			impl.OnPadHold(r, e)
		}
	case PadRelease:
		if impl, ok := d.(interface{ OnPadRelease(r *Root, e Event) bool }); ok {
			impl.OnPadRelease(r, e)
		}
	case PadMove:
		if impl, ok := d.(interface{ OnPadMove(r *Root, e Event) bool }); ok {
			impl.OnPadMove(r, e)
		}
	case KeyPress:
		if impl, ok := d.(interface{ OnKeyPress(r *Root, e Event) bool }); ok {
			impl.OnKeyPress(r, e)
		}
	case KeyHold:
		if impl, ok := d.(interface{ OnKeyHold(r *Root, e Event) bool }); ok {
			impl.OnKeyHold(r, e)
		}
	case KeyRelease:
		if impl, ok := d.(interface{ OnKeyRelease(r *Root, e Event) bool }); ok {
			impl.OnKeyRelease(r, e)
		}
	case KeyText:
		if impl, ok := d.(interface{ OnKeyText(r *Root, e Event) bool }); ok {
			impl.OnKeyText(r, e)
		}
	case TouchPress:
		if impl, ok := d.(interface{ OnTouchPress(r *Root, e Event) bool }); ok {
			impl.OnTouchPress(r, e)
		}
	case TouchHold:
		if impl, ok := d.(interface{ OnTouchHold(r *Root, e Event) bool }); ok {
			impl.OnTouchHold(r, e)
		}
	case TouchRelease:
		if impl, ok := d.(interface{ OnTouchRelease(r *Root, e Event) bool }); ok {
			impl.OnTouchRelease(r, e)
		}
	case MousePress:
		if impl, ok := d.(interface{ OnMousePress(r *Root, e Event) bool }); ok {
			impl.OnMousePress(r, e)
		}
	case MouseRelease:
		if impl, ok := d.(interface{ OnMouseRelease(r *Root, e Event) bool }); ok {
			impl.OnMouseRelease(r, e)
		}
	case MouseHold:
		if impl, ok := d.(interface{ OnMouseHold(r *Root, e Event) bool }); ok {
			impl.OnMouseHold(r, e)
		}
	case MouseMove:
		if impl, ok := d.(interface{ OnMouseMove(r *Root, e Event) bool }); ok {
			impl.OnMouseMove(r, e)
		}
	case MouseWheel:
		if impl, ok := d.(interface{ OnMouseWheel(r *Root, e Event) bool }); ok {
			impl.OnMouseWheel(r, e)
		}
	case ActionFocus:
		if impl, ok := d.(interface{ OnActionFocus(r *Root, e Event) bool }); ok {
			impl.OnActionFocus(r, e)
		}
	case ActionBlur:
		if impl, ok := d.(interface{ OnActionBlur(r *Root, e Event) bool }); ok {
			impl.OnActionBlur(r, e)
		}
	case ActionHover:
		if impl, ok := d.(interface{ OnActionHover(r *Root, e Event) bool }); ok {
			impl.OnActionHover(r, e)
		}
	case ActionCrash:
		if impl, ok := d.(interface{ OnActionCrash(r *Root, e Event) bool }); ok {
			impl.OnActionCrash(r, e)
		}
	case ActionDrag:
		if impl, ok := d.(interface{ OnActionDrag(r *Root, e Event) bool }); ok {
			impl.OnActionDrag(r, e)
		}
	case ActionDrop:
		if impl, ok := d.(interface{ OnActionDrop(r *Root, e Event) bool }); ok {
			impl.OnActionDrop(r, e)
		}
	case ActionMark:
		if impl, ok := d.(interface{ OnActionMark(r *Root, e Event) bool }); ok {
			impl.OnActionMark(r, e)
		}
	case ActionClean:
		if impl, ok := d.(interface{ OnActionClean(r *Root, e Event) bool }); ok {
			impl.OnActionClean(r, e)
		}
	case LayoutGet:
		if impl, ok := d.(interface{ OnLayoutGet(r *Root, e Event) bool }); ok {
			impl.OnLayoutGet(r, e)
		}
	case LayoutSet:
		if impl, ok := d.(interface{ OnLayoutSet(r *Root, e Event) bool }); ok {
			impl.OnLayoutSet(r, e)
		}
	default:
		panic("Unknown message")
	}
	return false
}

// Widget is a widget in the UI. It is part of a panel.
type Widget struct {
	Control           // A widget must be a Control.
	Bounds  Rectangle // Actual position and size of the widget.
	Size    Rectangle // Size is the desired size of the widget, may be bigger than Bounds.
	Style   Style
	State   State
}

// Panel is a panel in the UI it is a part of the UI that
// responds to input events.
type Panel struct {
	Control           // A panel must be a Control.
	Bounds  Rectangle // Actual position of the panel.
	Size    Rectangle // Size is the desired size of the panel, may be bigger than Bounds, and offset from it for scrolling.
	Style   Style
	State   State
	Widgets []*Widget // Widgets of the panel.
}

func (r *Root) FindTop(at Point) *Panel {
	for i := len(r.Panels) - 1; i >= 0; i-- {
		p := r.Panels[i]
		if at.In(p.Bounds) {
			return p
		}
	}
	return nil
}

func (r *Root) HandleMouseMove(e Event) bool {
	hover := r.FindTop(e.At)

	if r.Hover != nil && r.Hover != hover {
		r.Hover.Handle(r, Event{Msg: ActionCrash, At: e.At})
	}

	r.Hover = hover
	if r.Hover != nil {
		return r.Hover.Handle(r, Event{Msg: ActionHover, At: e.At})
	}
	return false
}

func (r *Root) On(e Event) bool {
	switch e.Msg {
	case PadDetach:
		return r.Default.Handle(r, e)
	case PadAttach:
		return r.Default.Handle(r, e)
	case PadPress:
		return r.Default.Handle(r, e)
	case PadHold:
		return r.Default.Handle(r, e)
	case PadRelease:
		return r.Default.Handle(r, e)
	case PadMove:
		return r.Default.Handle(r, e)
	case KeyPress:
		return r.Default.Handle(r, e)
	case KeyHold:
		return r.Default.Handle(r, e)
	case KeyRelease:
		return r.Default.Handle(r, e)
	case KeyText:
		return r.Default.Handle(r, e)
	case TouchPress:
		return r.Default.Handle(r, e)
	case TouchHold:
		return r.Default.Handle(r, e)
	case TouchRelease:
		return r.Default.Handle(r, e)
	case MousePress:
		return r.Default.Handle(r, e)
	case MouseMove:
		return r.HandleMouseMove(e)
	case MouseRelease:
		return r.Default.Handle(r, e)
	case MouseHold:
		return r.Default.Handle(r, e)
	case MouseWheel:
		return r.Default.Handle(r, e)
	case ActionFocus:
		return r.Default.Handle(r, e)
	case ActionBlur:
		return r.Default.Handle(r, e)
	case ActionHover:
		return r.Default.Handle(r, e)
	case ActionCrash:
		return r.Default.Handle(r, e)
	case ActionDrag:
		return r.Default.Handle(r, e)
	case ActionDrop:
		return r.Default.Handle(r, e)
	case ActionMark:
		return r.Default.Handle(r, e)
	case ActionClean:
		return r.Default.Handle(r, e)
	case LayoutGet:
		return r.Default.Handle(r, e)
	case LayoutSet:
		return r.Default.Handle(r, e)
	default:
		panic("Unknown event message")
	}
}

// Update is called 60 times per second.
// Input should be checked during this function.
func (r *Root) Update() error {
	for _, gid := range r.gamepads {
		if inpututil.IsGamepadJustDisconnected(gid) {
			r.On(Event{Msg: PadDetach, ID: int(gid)})
		}
	}

	r.connected = inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, gid := range r.connected {
		r.On(Event{Msg: PadAttach, ID: int(gid)})
	}

	r.gamepads = r.gamepads[0:0]
	r.gamepads = ebiten.AppendGamepadIDs(r.gamepads)
	for _, gid := range r.gamepads {
		buttons := inpututil.AppendJustPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.On(Event{Msg: PadPress, ID: int(gid), Button: int(button)})
		}

		buttons = inpututil.AppendPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			dur := inpututil.GamepadButtonPressDuration(gid, button)
			r.On(Event{Msg: PadHold, ID: int(gid), Button: int(button), Duration: dur})
		}

		buttons = inpututil.AppendJustReleasedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.On(Event{Msg: PadRelease, ID: int(gid), Button: int(button)})
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
			r.On(Event{Msg: PadRelease, ID: int(gid), Axes: axes})
		}
	}

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		r.On(Event{Msg: KeyPress, Code: int(key)})
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		dur := inpututil.KeyPressDuration(key)
		r.On(Event{Msg: KeyHold, Code: int(key), Duration: dur})
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		r.On(Event{Msg: KeyRelease, Code: int(key)})
	}

	if len(r.chars) == 0 && cap(r.chars) == 0 {
		r.chars = make([]rune, 0, 32)
	} else {
		r.chars = r.chars[0:0]
	}

	r.chars = ebiten.AppendInputChars(r.chars)
	if len(r.chars) > 0 {
		r.On(Event{Msg: KeyText, ID: -1, Chars: string(r.chars)})
	}
	r.chars = r.chars[0:0]

	for id, field := range r.TextInputFields {
		if field.IsFocused() {
			handled, _ := field.HandleInput(field.X, field.Y)
			if handled {
				r.On(Event{Msg: KeyText, ID: id, Chars: field.Text()})
			}
		}
	}

	touches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.On(Event{Msg: TouchPress, ID: int(touch), At: image.Pt(x, y)})
	}

	touches = ebiten.AppendTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		px, py := inpututil.TouchPositionInPreviousTick(touch)
		dx, dy := x-px, y-py
		dur := inpututil.TouchPressDuration(touch)
		r.On(Event{Msg: TouchHold, ID: int(touch), At: image.Pt(x, y), Delta: image.Pt(dx, dy), Duration: dur})
	}

	touches = inpututil.AppendJustReleasedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.On(Event{Msg: TouchRelease, ID: int(touch), At: image.Pt(x, y)})
	}

	x, y := ebiten.CursorPosition()
	dx, dy := x-r.cx, y-r.cy
	at := image.Pt(x, y)
	delta := image.Pt(dx, dy)

	for mb := ebiten.MouseButton(0); mb < ebiten.MouseButtonMax; mb++ {
		if inpututil.IsMouseButtonJustPressed(mb) {
			r.On(Event{Msg: MousePress, At: at, Delta: delta})
		}
		if ebiten.IsMouseButtonPressed(mb) {
			dur := inpututil.MouseButtonPressDuration(mb)
			r.On(Event{Msg: MouseHold, At: at, Delta: delta, Duration: dur})
		}
		if inpututil.IsMouseButtonJustReleased(mb) {
			r.On(Event{Msg: MouseRelease, At: at, Delta: delta})
		}
	}
	if dx != 0 || dy != 0 {
		r.On(Event{Msg: MouseMove, At: at, Delta: delta})
	}
	r.cx = x
	r.cy = y

	wx, wy := ebiten.Wheel()
	if wx != 0 || wy != 0 {
		wheel := image.Pt(int(wx), int(wy))
		r.On(Event{Msg: MouseWheel, At: at, Delta: delta, Wheel: wheel})
	}

	return nil
}

// Draw is called when the UI needs to be drawn in game.ui
func (r *Root) Draw(screen *Surface) {
	for _, p := range r.Panels {
		if !p.State.Hide {
			p.Render(r, screen)
		}
	}
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

func NewBox(bounds Rectangle) *Panel {
	p := &Panel{Bounds: bounds, Style: DefaultStyle()}
	p.Control = &box{Panel: p}
	return p
}

type box struct {
	*Panel
}

// Render is called when the element needs to be drawn
func (b box) Render(r *Root, screen *Surface) {
	style := b.Style
	if b.State.Hover {
		style = HoverStyle()
	}
	style.DrawBox(screen, b.Bounds)
	for _, w := range b.Widgets {
		if !w.State.Hide {
			w.Render(r, screen)
		}
	}
}

func (b *box) Handle(r *Root, e Event) bool {
	switch e.Msg {
	case ActionHover:
		b.State.Hover = true
		return true
	case ActionCrash:
		b.State.Hover = false
		return true
	default:
		return false
	}
}

func NewButton(bounds Rectangle, text string) *Widget {
	b := &Widget{Bounds: bounds, Style: DefaultStyle()}
	b.Control = &button{Widget: b, Text: text}
	return b
}

type button struct {
	*Widget
	Text    string
	Clicked func(*button)
	pressed bool
	Result  int // May be set freely except on dialog buttons.
}

func (b button) Render(r *Root, screen *Surface) {
	box := b.Bounds

	if b.pressed {
		box = box.Add(b.Style.Margin)
	}

	at := box.Min

	b.Style.DrawBox(screen, box)
	b.Style.DrawText(screen, at, b.Text)
}

func (b *button) Handle(r *Root, e Event) bool {
	return Dispatch(b, r, e)
}

func (b *button) OnMousePress(r *Root, ev Event) bool {
	b.pressed = true
	return true
}

func (b *button) OnMouseRelease(r *Root, ev Event) bool {
	b.pressed = false
	if b.Clicked != nil {
		b.Clicked(b)
	}
	return true
}

func (p *Panel) AddButton(bounds Rectangle, text string) *Widget {
	b := NewButton(bounds, text)
	p.Widgets = append(p.Widgets, b)
	return b
}
