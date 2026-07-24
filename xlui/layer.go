package xlui

import (
	"github.com/xmasengine/xmas/xgal"
)

// Reply is the result of several event handlers.
// Event handlers must strictly observe the meaning of Reply.
// Otherwise the widgets, in particylar widget focus may malfunction.
type Reply int

const (
	Ignore  Reply = iota // Ignore: the widget ignored the input, other widgets *should* process it.
	Accept               // Accept: the widget accepted the input, other widgets *must not* process it.
	Raise                // Raise: the widget accepted and needs to be raised higher in the layer stack.
	Lower                // Lower: the widget accepted and needs to be lowered in the layer stack.
	Finish               // Finish: the widget is done processing and should be considered closed.
	Proceed              // Proceed: the widget accepted the input but other widgets should continue processing. Only for sub widgets.
)

// Orientation is the layout orientation for layers in a group.
type Orientation int

const (
	Vertical Orientation = iota
	Horizontal
)

// Layer is a layer in the UI.
type Layer struct {
	// Class for custom or type specific behavior.
	Class

	Controls []*Control
	Bounds   xgal.Rectangle
	Clip     *xgal.Rectangle
	Style
	From        xgal.Point
	Done        bool
	Lock        bool
	Drag        bool
	Orientation Orientation // layout orientation in the group
}

func NewLayer(bounds xgal.Rectangle) *Layer {
	return &Layer{Bounds: bounds, Style: DefaultStyle(), Orientation: Horizontal}
}

func (l Layer) Render(s *xgal.Surface) {
	if l.Class.Render == nil {
		l.Style.DrawBox(s, l.Bounds)
	} else {
		l.Class.Render(s)
	}

	for i := len(l.Controls) - 1; i >= 0; i-- {
		ctrl := l.Controls[i]
		ctrl.Render(s)
	}
}

// MoveBy moves all children relative to current position.
func (l *Layer) MoveBy(delta xgal.Point) {
	l.Bounds = l.Bounds.Add(delta)
	for i := 0; i < len(l.Controls); i++ {
		l.Controls[i].MoveBy(delta)
	}
}

// Appends adds a control to this layer and lays it out by a simple line algorithm.
func (l *Layer) Append(ctrl *Control) *Control {
	if len(l.Controls) == 0 {
		ctrl.Bounds = xgal.Bound(ctrl.Bounds.Min.X+2, ctrl.Bounds.Min.Y+2, ctrl.Bounds.Dx(), ctrl.Bounds.Dy())
	} else {
		last := l.Controls[len(l.Controls)-1]
		if l.Orientation == Horizontal && last.Bounds.Dx()+ctrl.Bounds.Dx() < l.Bounds.Dx() {
			// fits on the line
			ctrl.Bounds = xgal.Bound(last.Bounds.Max.X+2, last.Bounds.Min.Y, ctrl.Bounds.Dx(), ctrl.Bounds.Dy())

		} else {
			ctrl.Bounds = xgal.Bound(ctrl.Bounds.Min.X+2, last.Bounds.Max.Y+2, ctrl.Bounds.Dx(), ctrl.Bounds.Dy())
		}
	}
	l.Controls = append(l.Controls, ctrl)
	return ctrl
}

func (l *Layer) Label(text string) *Control {
	at := l.Bounds.Min
	ctrl := NewLabel(at, text)
	return l.Append(ctrl)
}

func (l *Layer) Button(text string) *Control {
	at := l.Bounds.Min
	ctrl := NewButton(at, text)
	return l.Append(ctrl)
}

func (l *Layer) Click(at xgal.Point, button int) Reply {
	if l.Class.Click != nil {
		return l.Class.Click(at, button)
	}

	for i := len(l.Controls) - 1; i >= 0; i-- {
		ctrl := l.Controls[i]
		if ctrl == nil {
			continue
		}
		if !at.In(ctrl.Bounds) {
			continue
		}
		if ctrl.Class.Click != nil {
			res := ctrl.Class.Click(at, button)
			return res
		}
	}
	return Ignore
}
