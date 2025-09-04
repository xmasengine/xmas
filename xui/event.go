package xui

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

// Event has a Message and can dispatch itself, it also has a reference to the root.
// The dispatch function is normally of the signature On<Message>(r *Root, e <EventType>).
type Event interface {
	Message() Message
	Dispatch(EventHandler) bool
	Root() *Root
}

type EventHandler interface {
	HandleEvent(e Event) bool
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
	OnKeyPress(key int) Result
	OnKeyHold(key, duration int) Result
	OnKeyRelease(key int) Result
	OnInputText(id int, chars string) Result
}

type TouchHandler interface {
	OnTouchPress(tid int, at Point) Result
	OnTouchHold(tid int, at, delta Point, duration int) Result
	OnTouchRelease(tid int, at Point) Result
}

type MouseHandler interface {
	OnMousePress(at, delta Point, button int) Result
	OnMouseHold(at, delta Point, button, duration int) Result
	OnMouseRelease(at, delta Point, button int) Result
	OnMouseMove(at, delta Point) Result
	OnMouseWheel(at, delta, wheel Point) Result
}

type ActionHandler interface {
	OnFocus(at Point) Result
	OnBlur(at Point) Result
	OnHover(at Point) Result
	OnUnhover(at Point) Result
	OnDrag(at, delta Point) Result
	OnDrop(at, delta Point) Result
	OnMark(at Point) Result
	OnUnmark(at, delta Point) Result
}

type BasicEvent struct {
	R   *Root
	Msg Message
}

func MakeBasicEvent(r *Root, msg Message) BasicEvent {
	return BasicEvent{R: r, Msg: msg}
}

func (e BasicEvent) Message() Message {
	return e.Msg
}

func (e BasicEvent) Root() *Root {
	return e.R
}

func (e BasicEvent) Dispatch(c EventHandler) bool {
	return false
}

type PadEvent struct {
	BasicEvent
	ID       int
	Button   int
	Duration int
	Axes     []float64
}

func MakePadEvent(r *Root, msg Message, id, button, duration int, axes []float64) PadEvent {
	return PadEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		ID:         id, Button: button, Axes: axes,
	}
}

func (e PadEvent) Dispatch(c EventHandler) bool {
	switch e.Msg {
	case PadDetach:
		if impl, ok := c.(interface {
			OnPadDetach(e PadEvent) bool
		}); ok {
			return impl.OnPadDetach(e)
		}
	case PadAttach:
		if impl, ok := c.(interface {
			OnPadAttach(e PadEvent) bool
		}); ok {
			return impl.OnPadAttach(e)
		}
	case PadPress:
		if impl, ok := c.(interface {
			OnPadPress(e PadEvent) bool
		}); ok {
			return impl.OnPadPress(e)
		}
	case PadHold:
		if impl, ok := c.(interface {
			OnPadHold(e PadEvent) bool
		}); ok {
			return impl.OnPadHold(e)
		}
	case PadRelease:
		if impl, ok := c.(interface {
			OnPadRelease(e PadEvent) bool
		}); ok {
			return impl.OnPadRelease(e)
		}
	case PadMove:
		if impl, ok := c.(interface {
			OnPadMove(e PadEvent) bool
		}); ok {
			return impl.OnPadMove(e)
		}
	}
	return c.HandleEvent(e)
}

type KeyEvent struct {
	BasicEvent
	ID       int
	Code     int
	Duration int
	Chars    string
}

func MakeKeyEvent(r *Root, msg Message, id, code, duration int, chars string) KeyEvent {
	return KeyEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		ID:         id, Code: code, Chars: chars,
	}
}

func (e KeyEvent) Dispatch(c EventHandler) bool {
	switch e.Msg {
	case KeyPress:
		if impl, ok := c.(interface {
			OnKeyPress(e KeyEvent) bool
		}); ok {
			return impl.OnKeyPress(e)
		}
	case KeyHold:
		if impl, ok := c.(interface {
			OnKeyHold(e KeyEvent) bool
		}); ok {
			return impl.OnKeyHold(e)
		}
	case KeyRelease:
		if impl, ok := c.(interface {
			OnKeyRelease(e KeyEvent) bool
		}); ok {
			return impl.OnKeyRelease(e)
		}
	case KeyText:
		if impl, ok := c.(interface {
			OnKeyText(e KeyEvent) bool
		}); ok {
			return impl.OnKeyText(e)
		}
	}
	return c.HandleEvent(e)
}

type TouchEvent struct {
	BasicEvent
	ID       int
	At       Point
	Delta    Point
	Duration int
}

func MakeTouchEvent(r *Root, msg Message, id int, at, delta Point, duration int) TouchEvent {
	return TouchEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		ID:         id, At: at, Delta: delta, Duration: duration,
	}
}

func (e TouchEvent) Dispatch(c EventHandler) bool {
	switch e.Msg {
	case TouchPress:
		if impl, ok := c.(interface {
			OnTouchPress(e TouchEvent) bool
		}); ok {
			return impl.OnTouchPress(e)
		}
	case TouchHold:
		if impl, ok := c.(interface {
			OnTouchHold(e TouchEvent) bool
		}); ok {
			return impl.OnTouchHold(e)
		}
	case TouchRelease:
		if impl, ok := c.(interface {
			OnTouchRelease(e TouchEvent) bool
		}); ok {
			return impl.OnTouchRelease(e)
		}
	}
	return c.HandleEvent(e)
}

type MouseEvent struct {
	BasicEvent
	At       Point
	Delta    Point
	Duration int
	Wheel    Point
}

func MakeMouseEvent(r *Root, msg Message, at, delta Point, duration int) MouseEvent {
	return MouseEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		At:         at, Delta: delta, Duration: duration,
	}
}

func MakeMouseWheelEvent(r *Root, msg Message, at, delta, wheel Point) MouseEvent {
	return MouseEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		At:         at, Delta: delta, Wheel: wheel,
	}
}

func (e MouseEvent) Dispatch(c EventHandler) bool {
	println("MouseEvent.Dispatch")

	switch e.Msg {
	case MousePress:
		if impl, ok := c.(interface {
			OnMousePress(e MouseEvent) bool
		}); ok {
			return impl.OnMousePress(e)
		}
	case MouseRelease:
		if impl, ok := c.(interface {
			OnMouseRelease(e MouseEvent) bool
		}); ok {
			return impl.OnMouseRelease(e)
		}
	case MouseHold:
		if impl, ok := c.(interface {
			OnMouseHold(e MouseEvent) bool
		}); ok {
			return impl.OnMouseHold(e)
		}
	case MouseMove:
		if impl, ok := c.(interface {
			OnMouseMove(e MouseEvent) bool
		}); ok {
			println("MouseEvent.Dispatch: MouseMove ok")
			return impl.OnMouseMove(e)
		}
	case MouseWheel:
		if impl, ok := c.(interface {
			OnMouseWheel(e MouseEvent) bool
		}); ok {
			return impl.OnMouseWheel(e)
		}
	}
	println("MouseEvent.Dispatch uses the default handler")
	return c.HandleEvent(e)
}

type ActionEvent struct {
	BasicEvent
	At    Point
	Delta Point
}

func (e ActionEvent) Dispatch(c EventHandler) bool {
	switch e.Msg {
	case ActionBlur:
		if impl, ok := c.(interface{ OnActionBlur(e ActionEvent) bool }); ok {
			return impl.OnActionBlur(e)
		}
	case ActionHover:
		if impl, ok := c.(interface{ OnActionHover(e ActionEvent) bool }); ok {
			return impl.OnActionHover(e)
		}
	case ActionCrash:
		if impl, ok := c.(interface{ OnActionCrash(e ActionEvent) bool }); ok {
			return impl.OnActionCrash(e)
		}
	case ActionDrag:
		if impl, ok := c.(interface{ OnActionDrag(e ActionEvent) bool }); ok {
			return impl.OnActionDrag(e)
		}
	case ActionDrop:
		if impl, ok := c.(interface{ OnActionDrop(e ActionEvent) bool }); ok {
			return impl.OnActionDrop(e)
		}
	case ActionMark:
		if impl, ok := c.(interface{ OnActionMark(e ActionEvent) bool }); ok {
			return impl.OnActionMark(e)
		}
	case ActionClean:
		if impl, ok := c.(interface{ OnActionClean(e ActionEvent) bool }); ok {
			return impl.OnActionClean(e)
		}
	}
	return c.HandleEvent(e)
}

func MakeActionEvent(r *Root, msg Message, at, delta Point) ActionEvent {
	return ActionEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		At:         at, Delta: delta,
	}
}

type LayoutEvent struct {
	BasicEvent
	Bounds Rectangle
}

func MakeLayoutEvent(r *Root, msg Message, id int, bounds Rectangle) LayoutEvent {
	return LayoutEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		Bounds:     bounds,
	}
}

func (e LayoutEvent) Dispatch(c EventHandler) bool {
	switch e.Msg {
	case LayoutGet:
		if impl, ok := c.(interface{ OnLayoutGet(e LayoutEvent) bool }); ok {
			return impl.OnLayoutGet(e)
		}
	case LayoutSet:
		if impl, ok := c.(interface{ OnLayoutSet(e LayoutEvent) bool }); ok {
			return impl.OnLayoutSet(e)
		}
	}
	return c.HandleEvent(e)
}
