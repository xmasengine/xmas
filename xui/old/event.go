package xui

import "log/slog"

// Message is a kind of message that is sent to the listeners.
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
	Msg    Message
	Pad    PadEvent
	Touch  TouchEvent
	Key    KeyEvent
	Mouse  MouseEvent
	Action ActionEvent
	Layout LayoutEvent
}

type Eventer interface {
	Event() Event
}

func (e Event) Event() Event {
	return e
}

// Event can dispatch itself to a Listener.
func (e Event) Dispatch(l Listener) bool {
	slog.Debug("Event.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case PadDetach, PadAttach, PadPress, PadHold, PadRelease, PadMove:
		return e.Pad.Dispatch(l)

	case KeyPress, KeyHold, KeyRelease, KeyText:
		return e.Key.Dispatch(l)

	case TouchPress, TouchHold, TouchRelease:
		return e.Touch.Dispatch(l)

	case MousePress, MouseRelease, MouseHold, MouseMove, MouseWheel:
		return e.Mouse.Dispatch(l)

	case ActionFocus, ActionBlur, ActionHover, ActionCrash, ActionDrag, ActionDrop, ActionMark, ActionClean:
		return e.Action.Dispatch(l)

	case LayoutGet, LayoutSet:
		return e.Layout.Dispatch(l)

	default:
		return l.HandleEvent(e)
	}
}

type EventHandler interface {
	HandleEvent(e Event) bool
}

// Listener must implement all possible event handlers
// as well as EventHandler as a fallback.
type Listener interface {
	EventHandler
	PadHandler
	KeyHandler
	TouchHandler
	MouseHandler
	ActionHandler
	LayoutHandler
}

type PadHandler interface {
	OnPadDetach(PadEvent) bool
	OnPadAttach(PadEvent) bool
	OnPadPress(PadEvent) bool
	OnPadHold(PadEvent) bool
	OnPadRelease(PadEvent) bool
	OnPadMove(PadEvent) bool
}

type KeyHandler interface {
	OnKeyPress(KeyEvent) bool
	OnKeyHold(KeyEvent) bool
	OnKeyRelease(KeyEvent) bool
	OnKeyText(KeyEvent) bool
}

type TouchHandler interface {
	OnTouchPress(TouchEvent) bool
	OnTouchHold(TouchEvent) bool
	OnTouchRelease(TouchEvent) bool
}

type MouseHandler interface {
	OnMousePress(MouseEvent) bool
	OnMouseHold(MouseEvent) bool
	OnMouseRelease(MouseEvent) bool
	OnMouseMove(MouseEvent) bool
	OnMouseWheel(MouseEvent) bool
}

type ActionHandler interface {
	OnActionFocus(ActionEvent) bool
	OnActionBlur(ActionEvent) bool
	OnActionHover(ActionEvent) bool
	OnActionCrash(ActionEvent) bool
	OnActionDrag(ActionEvent) bool
	OnActionDrop(ActionEvent) bool
	OnActionMark(ActionEvent) bool
	OnActionClean(ActionEvent) bool
}

type LayoutHandler interface {
	OnLayoutGet(LayoutEvent) bool
	OnLayoutSet(LayoutEvent) bool
}

type BasicEvent struct {
	R   *Root
	Msg Message
}

func MakeBasicEvent(msg Message, r *Root) BasicEvent {
	return BasicEvent{Msg: msg, R: r}
}

func (e BasicEvent) Root() *Root {
	return e.R
}

type PadEvent struct {
	BasicEvent
	ID       int
	Button   int
	Duration int
	Axes     []float64
}

func MakePadEvent(msg Message, r *Root, id, button, duration int, axes []float64) PadEvent {
	return PadEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		ID:         id, Button: button, Axes: axes,
	}
}

// PadEvent can dispatch itself to a Padhandler.
func (e PadEvent) Dispatch(p PadHandler) bool {
	slog.Debug("Event.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case PadDetach:
		return p.OnPadDetach(e)
	case PadAttach:
		return p.OnPadAttach(e)
	case PadPress:
		return p.OnPadPress(e)
	case PadHold:
		return p.OnPadHold(e)
	case PadRelease:
		return p.OnPadRelease(e)
	case PadMove:
		return p.OnPadMove(e)
	default:
		return false
	}
}

// Event returns an event that wraps a pad event.
func (e PadEvent) Event() Event {
	return Event{Msg: e.Msg, Pad: e}
}

type BasicPadHandler struct {
	*Widget
}

func (BasicPadHandler) OnPadAttach(e PadEvent) bool {
	return false
}

func (BasicPadHandler) OnPadDetach(e PadEvent) bool {
	return false
}

func (BasicPadHandler) OnPadPress(e PadEvent) bool {
	return false
}

func (BasicPadHandler) OnPadHold(e PadEvent) bool {
	return false
}

func (BasicPadHandler) OnPadRelease(e PadEvent) bool {
	return false
}

func (BasicPadHandler) OnPadMove(e PadEvent) bool {
	return false
}

var _ PadHandler = BasicPadHandler{}

type KeyEvent struct {
	BasicEvent
	ID       int
	Code     int
	Duration int
	Chars    []rune
}

func MakeKeyEvent(msg Message, r *Root, id, code, duration int, chars ...rune) KeyEvent {
	return KeyEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		ID:         id,
		Code:       code,
		Chars:      chars,
	}
}

// KeyEvent can dispatch itself to a Keyhandler.
func (e KeyEvent) Dispatch(p KeyHandler) bool {
	slog.Debug("KeyEvent.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case KeyPress:
		return p.OnKeyPress(e)
	case KeyHold:
		return p.OnKeyHold(e)
	case KeyRelease:
		return p.OnKeyRelease(e)
	case KeyText:
		return p.OnKeyText(e)
	default:
		panic("Incorrect KeyEvent message")
	}
}

// Event returns an event that wraps a Key event.
func (e KeyEvent) Event() Event {
	return Event{Msg: e.Msg, Key: e}
}

type BasicKeyHandler struct {
	*Widget
}

func (BasicKeyHandler) OnKeyPress(e KeyEvent) bool {
	return false
}

func (BasicKeyHandler) OnKeyHold(e KeyEvent) bool {
	return false
}

func (BasicKeyHandler) OnKeyRelease(e KeyEvent) bool {
	return false
}

func (BasicKeyHandler) OnKeyText(e KeyEvent) bool {
	return false
}

var _ KeyHandler = BasicKeyHandler{}

type TouchEvent struct {
	BasicEvent
	ID       int
	At       Point
	Delta    Point
	Duration int
}

func MakeTouchEvent(msg Message, r *Root, id int, at, delta Point, duration int) TouchEvent {
	return TouchEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		ID:         id, At: at, Delta: delta, Duration: duration,
	}
}

// TouchEvent can dispatch itself to a Touchhandler.
func (e TouchEvent) Dispatch(p TouchHandler) bool {
	slog.Debug("TouchEvent.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case TouchPress:
		return p.OnTouchPress(e)
	case TouchHold:
		return p.OnTouchHold(e)
	case TouchRelease:
		return p.OnTouchRelease(e)
	default:
		return false
	}
}

// Event returns an event that wraps a Touch event.
func (e TouchEvent) Event() Event {
	return Event{Msg: e.Msg, Touch: e}
}

type BasicTouchHandler struct {
	*Widget
}

func (BasicTouchHandler) OnTouchPress(e TouchEvent) bool {
	return false
}

func (BasicTouchHandler) OnTouchHold(e TouchEvent) bool {
	return false
}

func (BasicTouchHandler) OnTouchRelease(e TouchEvent) bool {
	return false
}

var _ TouchHandler = BasicTouchHandler{}

type MouseEvent struct {
	BasicEvent
	Button   int
	At       Point
	Delta    Point
	Duration int
	Wheel    Point
}

func MakeMouseEvent(msg Message, r *Root, button int, at, delta Point, duration int) MouseEvent {
	return MouseEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		At:         at, Delta: delta, Duration: duration,
	}
}

func MakeMouseWheelEvent(msg Message, r *Root, at, delta, wheel Point) MouseEvent {
	return MouseEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		At:         at, Delta: delta, Wheel: wheel,
	}
}

// MouseEvent can dispatch itself to a Mousehandler.
func (e MouseEvent) Dispatch(p MouseHandler) bool {
	slog.Debug("Event.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case MousePress:
		return p.OnMousePress(e)
	case MouseHold:
		return p.OnMouseHold(e)
	case MouseRelease:
		return p.OnMouseRelease(e)
	case MouseMove:
		return p.OnMouseMove(e)
	case MouseWheel:
		return p.OnMouseWheel(e)
	default:
		return false
	}
}

// Event returns an event that wraps a Mouse event.
func (e MouseEvent) Event() Event {
	return Event{Msg: e.Msg, Mouse: e}
}

type BasicMouseHandler struct {
	*Widget
}

func (BasicMouseHandler) OnMouseAttach(e MouseEvent) bool {
	return false
}

func (BasicMouseHandler) OnMouseDetach(e MouseEvent) bool {
	return false
}

func (BasicMouseHandler) OnMousePress(e MouseEvent) bool {
	return false
}

func (BasicMouseHandler) OnMouseHold(e MouseEvent) bool {
	return false
}

func (BasicMouseHandler) OnMouseRelease(e MouseEvent) bool {
	return false
}

func (BasicMouseHandler) OnMouseMove(e MouseEvent) bool {
	return false
}

func (BasicMouseHandler) OnMouseWheel(e MouseEvent) bool {
	return false
}

var _ MouseHandler = BasicMouseHandler{}

type ActionEvent struct {
	BasicEvent
	At    Point
	Delta Point
}

// Event can dispatch itself to a Listener.
func (e ActionEvent) Dispatch(l ActionHandler) bool {
	slog.Debug("ActionEvent.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case ActionFocus:
		return l.OnActionFocus(e)

	case ActionBlur:
		return l.OnActionBlur(e)

	case ActionHover:
		return l.OnActionHover(e)

	case ActionCrash:
		return l.OnActionCrash(e)

	case ActionDrag:
		return l.OnActionDrag(e)

	case ActionDrop:
		return l.OnActionDrop(e)

	case ActionMark:
		return l.OnActionMark(e)
	default:
		return false
	}
	return false
}

// Event returns an event that wraps a Action event.
func (e ActionEvent) Event() Event {
	return Event{Msg: e.Msg, Action: e}
}

type BasicActionHandler struct {
	*Widget
}

func (BasicActionHandler) OnActionFocus(e ActionEvent) bool { return false }
func (BasicActionHandler) OnActionBlur(e ActionEvent) bool  { return false }
func (BasicActionHandler) OnActionHover(e ActionEvent) bool { return false }
func (BasicActionHandler) OnActionCrash(e ActionEvent) bool { return false }
func (BasicActionHandler) OnActionDrag(e ActionEvent) bool  { return false }
func (BasicActionHandler) OnActionDrop(e ActionEvent) bool  { return false }
func (BasicActionHandler) OnActionMark(e ActionEvent) bool  { return false }
func (BasicActionHandler) OnActionClean(e ActionEvent) bool { return false }

var _ ActionHandler = BasicActionHandler{}

func MakeActionEvent(msg Message, r *Root, at, delta Point) ActionEvent {
	return ActionEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		At:         at, Delta: delta,
	}
}

type LayoutEvent struct {
	BasicEvent
	Bounds Rectangle
}

func MakeLayoutEvent(msg Message, r *Root, id int, bounds Rectangle) LayoutEvent {
	return LayoutEvent{
		BasicEvent: MakeBasicEvent(msg, r),
		Bounds:     bounds,
	}
}

// Event can dispatch itself to a Listener.
func (e LayoutEvent) Dispatch(l LayoutHandler) bool {
	slog.Debug("LayoutEvent.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case LayoutGet:
		return l.OnLayoutGet(e)
	case LayoutSet:
		return l.OnLayoutSet(e)
	default:
		return false
	}
}

// Event returns an event that wraps a Layout event.
func (e LayoutEvent) Event() Event {
	return Event{Msg: e.Msg, Layout: e}
}

type BasicLayoutHandler struct {
	*Widget
}

var _ LayoutHandler = BasicLayoutHandler{}

func (BasicLayoutHandler) OnLayoutGet(e LayoutEvent) bool { return false }
func (BasicLayoutHandler) OnLayoutSet(e LayoutEvent) bool { return false }

type BasicEventHandler struct {
	*Widget
}

func (BasicEventHandler) HandleEvent(e Event) bool { return false }

var _ EventHandler = BasicEventHandler{}

type BasicListener struct {
	Widget *Widget
	BasicPadHandler
	BasicKeyHandler
	BasicTouchHandler
	BasicMouseHandler
	BasicActionHandler
	BasicLayoutHandler
	BasicEventHandler
}

var _ Listener = BasicListener{}
