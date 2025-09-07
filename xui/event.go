package xui

import "log/slog"

// Event can dispatch itself to a dispatcher, it also has a
// reference to the root.
type Event interface {
	Dispatch(Dispatcher) bool
	Root() *Root
}

// Dispatcher is normally one of the one method handlers below.
type Dispatcher any

type EventHandler interface {
	HandleEvent(e Event) bool
}

type PadDetachHandler interface{ OnPadDetach(PadDetachEvent) bool }
type PadAttachHandler interface{ OnPadAttach(PadAttachEvent) bool }
type PadPressHandler interface{ OnPadPress(PadPressEvent) bool }
type PadHoldHandler interface{ OnPadHold(PadHoldEvent) bool }
type PadReleaseHandler interface{ OnPadRelease(PadReleaseEvent) bool }
type PadMoveHandler interface{ OnPadMove(PadMoveEvent) bool }
type KeyPressHandler interface{ OnKeyPress(KeyPressEvent) bool }
type KeyHoldHandler interface{ OnKeyHold(KeyHoldEvent) bool }
type KeyReleaseHandler interface{ OnKeyRelease(KeyReleaseEvent) bool }
type KeyTextHandler interface{ OnKeyText(KeyTextEvent) bool }
type TouchPressHandler interface{ OnTouchPress(TouchPressEvent) bool }
type TouchHoldHandler interface{ OnTouchHold(TouchHoldEvent) bool }
type TouchReleaseHandler interface{ OnTouchRelease(TouchReleaseEvent) bool }
type MousePressHandler interface{ OnMousePress(MousePressEvent) bool }
type MouseHoldHandler interface{ OnMouseHold(MouseHoldEvent) bool }
type MouseReleaseHandler interface{ OnMouseRelease(MouseReleaseEvent) bool }
type MouseMoveHandler interface{ OnMouseMove(MouseMoveEvent) bool }
type MouseWheelHandler interface{ OnMouseWheel(MouseWheelEvent) bool }
type ActionFocusHandler interface{ OnActionFocus(ActionFocusEvent) bool }
type ActionBlurHandler interface{ OnActionBlur(ActionBlurEvent) bool }
type ActionHoverHandler interface{ OnActionHover(ActionHoverEvent) bool }
type ActionCrashHandler interface{ OnActionCrash(ActionCrashEvent) bool }
type ActionDragHandler interface{ OnActionDrag(ActionDragEvent) bool }
type ActionDropHandler interface{ OnActionDrop(ActionDropEvent) bool }
type ActionMarkHandler interface{ OnActionMark(ActionMarkEvent) bool }
type ActionCleanHandler interface{ OnActionClean(ActionCleanEvent) bool }
type LayoutGetHandler interface{ OnLayoutGet(LayoutGetEvent) bool }
type LayoutSetHandler interface{ OnLayoutSet(LayoutSetEvent) bool }

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

func (e BasicEvent) Dispatch(c Dispatcher) bool {
	return false
}

type PadDetachEvent struct{ PadEvent }
type PadAttachEvent struct{ PadEvent }
type PadPressEvent struct{ PadEvent }
type PadHoldEvent struct{ PadEvent }
type PadReleaseEvent struct{ PadEvent }
type PadMoveEvent struct{ PadEvent }
type KeyPressEvent struct{ KeyEvent }
type KeyHoldEvent struct{ KeyEvent }
type KeyReleaseEvent struct{ KeyEvent }
type KeyTextEvent struct{ KeyEvent }
type TouchPressEvent struct{ TouchEvent }
type TouchHoldEvent struct{ TouchEvent }
type TouchReleaseEvent struct{ TouchEvent }
type MousePressEvent struct{ MouseEvent }
type MouseReleaseEvent struct{ MouseEvent }
type MouseHoldEvent struct{ MouseEvent }
type MouseMoveEvent struct{ MouseEvent }
type MouseWheelEvent struct{ MouseEvent }
type ActionFocusEvent struct{ ActionEvent }
type ActionBlurEvent struct{ ActionEvent }
type ActionHoverEvent struct{ ActionEvent }
type ActionCrashEvent struct{ ActionEvent }
type ActionDragEvent struct{ ActionEvent }
type ActionDropEvent struct{ ActionEvent }
type ActionMarkEvent struct{ ActionEvent }
type ActionCleanEvent struct{ ActionEvent }
type LayoutGetEvent struct{ LayoutEvent }
type LayoutSetEvent struct{ LayoutEvent }

func (e PadDetachEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(PadDetachHandler); ok {
		return impl.OnPadDetach(e)
	}
	return false
}

func (e PadAttachEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(PadAttachHandler); ok {
		return impl.OnPadAttach(e)
	}
	return false
}

func (e PadPressEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(PadPressHandler); ok {
		return impl.OnPadPress(e)
	}
	return false
}

func (e PadHoldEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(PadHoldHandler); ok {
		return impl.OnPadHold(e)
	}
	return false
}

func (e PadReleaseEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(PadReleaseHandler); ok {
		return impl.OnPadRelease(e)
	}
	return false
}

func (e PadMoveEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(PadMoveHandler); ok {
		return impl.OnPadMove(e)
	}
	return false
}

func (e KeyPressEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(KeyPressHandler); ok {
		return impl.OnKeyPress(e)
	}
	return false
}

func (e KeyHoldEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(KeyHoldHandler); ok {
		return impl.OnKeyHold(e)
	}
	return false
}

func (e KeyReleaseEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(KeyReleaseHandler); ok {
		return impl.OnKeyRelease(e)
	}
	return false
}

func (e KeyTextEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(KeyTextHandler); ok {
		return impl.OnKeyText(e)
	}
	return false
}

func (e TouchPressEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(TouchPressHandler); ok {
		return impl.OnTouchPress(e)
	}
	return false
}

func (e TouchHoldEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(TouchHoldHandler); ok {
		return impl.OnTouchHold(e)
	}
	return false
}

func (e TouchReleaseEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(TouchReleaseHandler); ok {
		return impl.OnTouchRelease(e)
	}
	return false
}

func (e MousePressEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(MousePressHandler); ok {
		return impl.OnMousePress(e)
	}
	return false
}

func (e MouseReleaseEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(MouseReleaseHandler); ok {
		return impl.OnMouseRelease(e)
	}
	return false
}

func (e MouseHoldEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(MouseHoldHandler); ok {
		return impl.OnMouseHold(e)
	}
	return false
}

func (e MouseMoveEvent) Dispatch(d Dispatcher) bool {
	slog.Info("MouseMoveEvent", "d", d)
	if impl, ok := d.(MouseMoveHandler); ok {
		slog.Info("MouseMoveEvent ok")
		return impl.OnMouseMove(e)
	}
	return false
}

func (e MouseWheelEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(MouseWheelHandler); ok {
		return impl.OnMouseWheel(e)
	}
	return false
}

func (e ActionFocusEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionFocusHandler); ok {
		return impl.OnActionFocus(e)
	}
	return false
}

func (e ActionBlurEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionBlurHandler); ok {
		return impl.OnActionBlur(e)
	}
	return false
}

func (e ActionHoverEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionHoverHandler); ok {
		return impl.OnActionHover(e)
	}
	return false
}

func (e ActionCrashEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionCrashHandler); ok {
		return impl.OnActionCrash(e)
	}
	return false
}

func (e ActionDragEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionDragHandler); ok {
		return impl.OnActionDrag(e)
	}
	return false
}

func (e ActionDropEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionDropHandler); ok {
		return impl.OnActionDrop(e)
	}
	return false
}

func (e ActionMarkEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionMarkHandler); ok {
		return impl.OnActionMark(e)
	}
	return false
}

func (e ActionCleanEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(ActionCleanHandler); ok {
		return impl.OnActionClean(e)
	}
	return false
}

func (e LayoutGetEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(LayoutGetHandler); ok {
		return impl.OnLayoutGet(e)
	}
	return false
}

func (e LayoutSetEvent) Dispatch(d Dispatcher) bool {
	if impl, ok := d.(LayoutSetHandler); ok {
		return impl.OnLayoutSet(e)
	}
	return false
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
