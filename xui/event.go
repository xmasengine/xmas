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
	impl, ok := c.(PadHandler)
	if !ok {
		return c.HandleEvent(e)
	}
	switch e.Msg {
	case PadDetach:
		return impl.OnPadDetach(e)
	case PadAttach:
		return impl.OnPadAttach(e)
	case PadPress:
		return impl.OnPadPress(e)
	case PadHold:
		return impl.OnPadHold(e)
	case PadRelease:
		return impl.OnPadRelease(e)
	case PadMove:
		return impl.OnPadMove(e)
	}
	return c.HandleEvent(e)
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

func MakeKeyEvent(r *Root, msg Message, id, code, duration int, chars string) KeyEvent {
	return KeyEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		ID:         id, Code: code, Chars: chars,
	}
}

func (e KeyEvent) Dispatch(c EventHandler) bool {
	impl, ok := c.(KeyHandler)
	if !ok {
		return c.HandleEvent(e)
	}
	switch e.Msg {
	case KeyPress:
		return impl.OnKeyPress(e)
	case KeyHold:
		return impl.OnKeyHold(e)
	case KeyRelease:
		return impl.OnKeyRelease(e)
	case KeyText:
		return impl.OnKeyText(e)
	}
	return c.HandleEvent(e)
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

func MakeTouchEvent(r *Root, msg Message, id int, at, delta Point, duration int) TouchEvent {
	return TouchEvent{
		BasicEvent: MakeBasicEvent(r, msg),
		ID:         id, At: at, Delta: delta, Duration: duration,
	}
}

func (e TouchEvent) Dispatch(c EventHandler) bool {
	impl, ok := c.(TouchHandler)
	if !ok {
		return c.HandleEvent(e)
	}
	switch e.Msg {
	case TouchPress:
		return impl.OnTouchPress(e)
	case TouchHold:
		return impl.OnTouchHold(e)
	case TouchRelease:
		return impl.OnTouchRelease(e)
	}
	return c.HandleEvent(e)
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

	impl, ok := c.(MouseHandler)
	if !ok {
		return c.HandleEvent(e)
	}

	switch e.Msg {
	case MousePress:
		return impl.OnMousePress(e)
	case MouseRelease:
		return impl.OnMouseRelease(e)
	case MouseHold:
		return impl.OnMouseHold(e)
	case MouseMove:
		println("MouseEvent.Dispatch: MouseMove ok")
		return impl.OnMouseMove(e)
	case MouseWheel:
		return impl.OnMouseWheel(e)
	}
	println("MouseEvent.Dispatch uses the default handler")
	return c.HandleEvent(e)
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

func (e ActionEvent) Dispatch(c EventHandler) bool {
	impl, ok := c.(ActionHandler)
	if !ok {
		return c.HandleEvent(e)
	}

	switch e.Msg {
	case ActionBlur:
		return impl.OnActionBlur(e)
	case ActionHover:
		return impl.OnActionHover(e)
	case ActionCrash:
		return impl.OnActionCrash(e)
	case ActionDrag:
		return impl.OnActionDrag(e)
	case ActionDrop:
		return impl.OnActionDrop(e)
	case ActionMark:
		return impl.OnActionMark(e)
	case ActionClean:
		return impl.OnActionClean(e)
	}
	return c.HandleEvent(e)
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
	impl, ok := c.(LayoutHandler)
	if !ok {
		return c.HandleEvent(e)
	}
	switch e.Msg {
	case LayoutGet:
		return impl.OnLayoutGet(e)
	case LayoutSet:
		return impl.OnLayoutSet(e)
	}
	return c.HandleEvent(e)
}

type BasicLayoutHandler struct {
	*Widget
}

var _ LayoutHandler = BasicLayoutHandler{}

func (BasicLayoutHandler) OnLayoutGet(e LayoutEvent) bool { return false }
func (BasicLayoutHandler) OnLayoutSet(e LayoutEvent) bool { return false }
