package tree

import "image"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/inpututil"
import "github.com/hajimehoshi/ebiten/v2/exp/textinput"

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

// Root is the root element of a UI.
// It is possble to use multiple roots however only one should be active
// at the same time.
type Root struct {
	Box

	TextInputFields []*TextInputField // Text input fields in use
	cx, cy          int
	chars           []rune
	KeyMods         KeyMods // Current key KeyMods
	connected       []ebiten.GamepadID
	gamepads        []ebiten.GamepadID
	NoTouchMouse    bool      // NoTouchMouse: set this to true to not translate touches to mouse events.
	Focus           Focusable // Element that has the input focus.
	Hover           Hoverable // Element that is being hovered by the mouse.
	Drag            Draggable // Element that is being dragged by the mouse.
	Mark            Markable  // Element that has the joystick and arrow key marker.
	Handler         Element   // Element that is the default event handler.
}

func (r *Root) Init(options ...Applier) *Root {
	r.Box.Init(nil, options...)
	for _, opt := range options {
		opt.Apply(r)
	}
	return r
}

func NewRoot(options ...Applier) *Root {
	r := Root{}
	return r.Init(options...)
}

func (r *Root) AppendTextInputField(x, y int) *TextInputField {
	field := &TextInputField{Point: image.Pt(x, y)}
	r.TextInputFields = append(r.TextInputFields, field)
	return field
}

/*


// KeyHandler is an Element that can handle a key.
type KeyHandler interface {
	Element
	KeyHandle(sym int, ch rune) Result
}

// MouseHandler is an Element that can handle mouse moves.
type MouseHandler interface {
	Element
	MouseHandle(delta Point) Result
}

// ClickHandler is an Element that can handle mouse clicks.
type ClickHandler interface {
	Element
	ClickHandle(delta Point, button int) Result
}

// PressHandler is an Element that can joypad button presses.
type PressHandler interface {
	Element
	PressHandle(button int) Result
}

// MoveHandler is an Element that can joypad axe motions.
type MoveHandler interface {
	Element
	MoveHandle(delta Point, axe int) Result
}


*/

func (r *Root) OnPadDetach(gid int) Result                 { return false }
func (r *Root) OnPadAttach(gid int) Result                 { return false }
func (r *Root) OnPadPress(gid, button int) Result          { return false }
func (r *Root) OnPadHold(gid, button, duration int) Result { return false }
func (r *Root) OnPadRelease(gid, button int) Result        { return false }
func (r *Root) OnPadMove(gid int, axes []float64) Result   { return false }

func (r *Root) OnKeyPress(key int) Result               { return false }
func (r *Root) OnKeyHold(key, duration int) Result      { return false }
func (r *Root) OnKeyRelease(key int) Result             { return false }
func (r *Root) OnInputText(id int, chars string) Result { return false }

func (r *Root) OnTouchPress(tid int, at Point) Result                     { return false }
func (r *Root) OnTouchHold(tid int, at, delta Point, duration int) Result { return false }
func (r *Root) OnTouchRelease(tid int, at Point) Result                   { return false }

func (r *Root) OnMousePress(at, delta Point, button int) Result {
	return false
}

func (r *Root) OnMouseHold(at, delta Point, button, duration int) Result { return false }
func (r *Root) OnMouseRelease(at, delta Point, button int) Result        { return false }

func (r *Root) OnMouseMove(at, delta Point) Result {
	var hover Hoverable
	var ok bool

	top := FindTop(at, &r.Box)
	if top == nil {
		return false
	}

	if hover, ok = top.(Hoverable); !ok {
		return false
	}

	if r.Hover != nil {
		r.Hover.OnUnhover(at)
	}
	hover.OnHover(at)
	r.Hover = hover
	return true
}

func (r *Root) OnMouseWheel(at, delta, wheel Point) Result { return false }

var _ UpdateHandler = &Root{}

// Update is called 60 times per second.
// Input should be checked during this function.
func (r *Root) Update() error {
	for _, gid := range r.gamepads {
		if inpututil.IsGamepadJustDisconnected(gid) {
			r.OnPadDetach(int(gid))
		}
	}

	r.connected = inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, gid := range r.connected {
		r.OnPadAttach(int(gid))
	}

	r.gamepads = r.gamepads[0:0]
	r.gamepads = ebiten.AppendGamepadIDs(r.gamepads)
	for _, gid := range r.gamepads {
		buttons := inpututil.AppendJustPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.OnPadPress(int(gid), int(button))
		}

		buttons = inpututil.AppendPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			dur := inpututil.GamepadButtonPressDuration(gid, button)
			r.OnPadHold(int(gid), int(button), dur)
		}

		buttons = inpututil.AppendJustReleasedGamepadButtons(gid, nil)
		for _, button := range buttons {
			r.OnPadRelease(int(gid), int(button))
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
			r.OnPadMove(int(gid), axes)
		}
	}

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		r.OnKeyPress(int(key))
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		dur := inpututil.KeyPressDuration(key)
		r.OnKeyHold(int(key), dur)
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		r.OnKeyRelease(int(key))
	}

	if len(r.chars) == 0 && cap(r.chars) == 0 {
		r.chars = make([]rune, 0, 32)
	} else {
		r.chars = r.chars[0:0]
	}

	r.chars = ebiten.AppendInputChars(r.chars)
	if len(r.chars) > 0 {
		r.OnInputText(-1, string(r.chars))
	}
	r.chars = r.chars[0:0]

	for id, field := range r.TextInputFields {
		if field.IsFocused() {
			handled, _ := field.HandleInput(field.X, field.Y)
			if handled {
				r.OnInputText(id, field.Text())
			}
		}
	}

	touches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.OnTouchPress(int(touch), image.Pt(x, y))
	}

	touches = ebiten.AppendTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		px, py := inpututil.TouchPositionInPreviousTick(touch)
		dx, dy := x-px, y-py
		dur := inpututil.TouchPressDuration(touch)
		r.OnTouchHold(int(touch), image.Pt(x, y), image.Pt(dx, dy), dur)
	}

	touches = inpututil.AppendJustReleasedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		r.OnTouchRelease(int(touch), image.Pt(x, y))
	}

	x, y := ebiten.CursorPosition()
	dx, dy := x-r.cx, y-r.cy
	at := image.Pt(x, y)
	delta := image.Pt(dx, dy)

	for mb := ebiten.MouseButton(0); mb < ebiten.MouseButtonMax; mb++ {
		if inpututil.IsMouseButtonJustPressed(mb) {
			r.OnMousePress(at, delta, int(mb))
		}
		if ebiten.IsMouseButtonPressed(mb) {
			dur := inpututil.MouseButtonPressDuration(mb)
			r.OnMouseHold(at, delta, dur, int(mb))
		}
		if inpututil.IsMouseButtonJustReleased(mb) {
			r.OnMouseRelease(at, delta, int(mb))
		}
	}
	if dx != 0 || dy != 0 {
		r.OnMouseMove(at, delta)
	}
	r.cx = x
	r.cy = y

	wx, wy := ebiten.Wheel()
	if wx != 0 || wy != 0 {
		wheel := image.Pt(int(wx), int(wy))
		r.OnMouseWheel(at, delta, wheel)
	}

	return nil
}

// Draw is called when the UI needs to be drawn
func (r *Root) Draw(screen *Surface) {
	l := r.Contain()
	for _, e := range l {
		state := e.State()
		if !state.Hide {
			e.Draw(screen)
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
