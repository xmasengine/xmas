package tree

type PadHandler interface {
	OnPadDetach(gid int) Result
	OnPadAttach(gid int) Result
	OnPadPress(gid, button int) Result
	OnPadHold(gid, button, duration int) Result
	OnPadRelease(gid, button int) Result
	OnPadMove(gid int, axes []float64) Result
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

type Focusable interface {
	Element
	OnFocus(at Point) Result
	OnBlur(at Point) Result
}

type Hoverable interface {
	Element
	OnHover(at Point) Result
	OnUnhover(at Point) Result
}

type Draggable interface {
	OnDrag(at, delta Point) Result
	OnDrop(at, delta Point) Result
}

type Markable interface {
	OnMark(at Point) Result
	OnUnmark(at, delta Point) Result
}

type UpdateHandler interface {
	PadHandler
	KeyHandler
	TouchHandler
	MouseHandler
}

type Event struct {
	Msg Message
	Pad struct {
		ID       int
		Button   int
		Duration int
		Axes     []float64
	}
	Key struct {
		ID       int
		Code     int
		Duration int
		Chars    string
	}
	Touch struct {
		ID       int
		At       Point
		Delta    Point
		Duration int
	}
	Mouse struct {
		ID       int
		At       Point
		Delta    Point
		Duration int
		Wheel    Point
	}
	Action struct {
		At    Point
		Delta Point
	}
	Layout struct {
		Bounds Rectangle
	}
}

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
)
