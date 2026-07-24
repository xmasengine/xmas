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
	Groups   []Group
	Controls []Control
	Bounds   xgal.Rectangle
	Clip     *xgal.Rectangle
	Style
	From        xgal.Point
	Done        bool
	Lock        bool
	Drag        bool
	Orientation Orientation // layout orientation in the group

	// Flexible handlers
	OnRender func(s *xgal.Surface, l Layer)
}

func MakeLayer(bounds xgal.Rectangle) Layer {
	return Layer{Bounds: bounds, Style: DefaultStyle(), Orientation: Vertical}
}

func (l Layer) Render(s *xgal.Surface) {
	l.Style.DrawBox(s, l.Bounds)
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
