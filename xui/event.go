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

// Event can dispatch itself to a Listener.
func (e Event) Dispatch(l Listener) bool {
	slog.Info("Event.Dispatch", "Msg", e.Msg)
	switch e.Msg {
	case PadDetach:
		return l.OnPadDetach(e.Pad)

	case PadAttach:
		return l.OnPadAttach(e.Pad)

	case PadPress:
		return l.OnPadPress(e.Pad)

	case PadHold:
		return l.OnPadHold(e.Pad)

	case PadRelease:
		return l.OnPadRelease(e.Pad)

	case PadMove:
		return l.OnPadMove(e.Pad)

	case KeyPress:
		return l.OnKeyPress(e.Key)

	case KeyHold:
		return l.OnKeyHold(e.Key)

	case KeyRelease:
		return l.OnKeyRelease(e.Key)

	case KeyText:
		return l.OnKeyText(e.Key)

	case TouchPress:
		return l.OnTouchPress(e.Touch)

	case TouchHold:
		return l.OnTouchHold(e.Touch)

	case TouchRelease:
		return l.OnTouchRelease(e.Touch)

	case MousePress:
		return l.OnMousePress(e.Mouse)

	case MouseRelease:
		return l.OnMouseRelease(e.Mouse)

	case MouseHold:
		return l.OnMouseHold(e.Mouse)

	case MouseMove:
		return l.OnMouseMove(e.Mouse)

	case MouseWheel:
		return l.OnMouseWheel(e.Mouse)

	case ActionFocus:
		return l.OnActionFocus(e.Action)

	case ActionBlur:
		return l.OnActionBlur(e.Action)

	case ActionHover:
		return l.OnActionHover(e.Action)

	case ActionCrash:
		return l.OnActionCrash(e.Action)

	case ActionDrag:
		return l.OnActionDrag(e.Action)

	case ActionDrop:
		return l.OnActionDrop(e.Action)

	case ActionMark:
		return l.OnActionMark(e.Action)

	case ActionClean:
		return l.OnActionClean(e.Action)
	case LayoutGet:
		return l.OnLayoutGet(e.Layout)
	case LayoutSet:
		return l.OnLayoutSet(e.Layout)
	default:
		return l.HandleEvent(e)
	}
	return false
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
	R *Root
}

func MakeBasicEvent(r *Root) BasicEvent {
	return BasicEvent{R: r}
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

func MakePadEvent(r *Root, id, button, duration int, axes []float64) PadEvent {
	return PadEvent{
		BasicEvent: MakeBasicEvent(r),
		ID:         id, Button: button, Axes: axes,
	}
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
	Chars    string
}

func MakeKeyEvent(r *Root, id, code, duration int, chars string) KeyEvent {
	return KeyEvent{
		BasicEvent: MakeBasicEvent(r),
		ID:         id, Code: code, Chars: chars,
	}
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

func MakeTouchEvent(r *Root, id int, at, delta Point, duration int) TouchEvent {
	return TouchEvent{
		BasicEvent: MakeBasicEvent(r),
		ID:         id, At: at, Delta: delta, Duration: duration,
	}
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
	At       Point
	Delta    Point
	Duration int
	Wheel    Point
}

func MakeMouseEvent(r *Root, at, delta Point, duration int) MouseEvent {
	return MouseEvent{
		BasicEvent: MakeBasicEvent(r),
		At:         at, Delta: delta, Duration: duration,
	}
}

func MakeMouseWheelEvent(r *Root, at, delta, wheel Point) MouseEvent {
	return MouseEvent{
		BasicEvent: MakeBasicEvent(r),
		At:         at, Delta: delta, Wheel: wheel,
	}
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

func MakeActionEvent(r *Root, at, delta Point) ActionEvent {
	return ActionEvent{
		BasicEvent: MakeBasicEvent(r),
		At:         at, Delta: delta,
	}
}

type LayoutEvent struct {
	BasicEvent
	Bounds Rectangle
}

func MakeLayoutEvent(r *Root, id int, bounds Rectangle) LayoutEvent {
	return LayoutEvent{
		BasicEvent: MakeBasicEvent(r),
		Bounds:     bounds,
	}
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
	BasicPadHandler
	BasicKeyHandler
	BasicTouchHandler
	BasicMouseHandler
	BasicActionHandler
	BasicLayoutHandler
	BasicEventHandler
}

var _ Listener = BasicListener{}
