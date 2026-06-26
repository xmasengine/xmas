// package xui implements simple layer based UI
package xui

import (
	"slices"

	"github.com/xmasengine/xmas/xgal"
)

type Reply int

const (
	Ignore Reply = iota
	Accept
	Raise
	Lower
	Finish
)

// Axis is the layout axis for child widgets in a container.
type Axis int

const (
	Vertical Axis = iota
	Horizontal
)

// A Widget is an element of an UI.
type Widget interface {
	Poll() Reply
	Render(screen *xgal.Surface)
	Place(bounds xgal.Rectangle) xgal.Rectangle
}

// Mover is an optional interface for widgets that can be repositioned
// without re-laying out their contents (e.g., scrolling or dragging).
type Mover interface {
	Widget
	MoveBy(delta xgal.Point)
}

// Layer is a container in the UI.
type Layer struct {
	Kids   []Widget
	Bounds xgal.Rectangle
	Style
	From xgal.Point
	Done bool
	Lock bool
	Drag bool
	Axis Axis // layout direction for child widgets
}

var _ Widget = &Layer{}

func MakeLayer(bounds xgal.Rectangle) Layer {
	return Layer{Bounds: bounds, Style: DefaultStyle(), Axis: Vertical}
}

func (m *Layer) Poll() Reply {
	res := m.PollKids()
	if res != Ignore {
		return res
	}
	return Ignore
}

func (m *Layer) Render(s *xgal.Surface) {
	m.Style.DrawBox(s, m.Bounds)
	m.RenderKids(s)
}

func (m *Layer) Add(g Widget) Widget {
	m.Kids = slices.Insert(m.Kids, 0, g)
	return g
}

func (m *Layer) PollKids() Reply {
	for i, kid := range m.Kids {
		if kid == nil {
			continue
		}
		res := kid.Poll()
		if res == Finish {
			m.Kids = slices.Delete(m.Kids, i, i+1)
		} else if res == Accept {
			break
		} else if res == Raise {
			if i < len(m.Kids)-1 {
				m.Kids[i], m.Kids[i+1] = m.Kids[i+1], m.Kids[i]
			}
			return res
		} else if res == Lower {
			if i > 0 {
				m.Kids[i], m.Kids[i-1] = m.Kids[i-1], m.Kids[i]
			}
			return res
		}
	}
	return Ignore
}

func (m *Layer) RenderKids(s *xgal.Surface) {
	for i := len(m.Kids) - 1; i >= 0; i-- {
		kid := m.Kids[i]
		if kid != nil {
			kid.Render(s)
		}
	}
}

// MoveBy moves all children relative to current position.
func (m *Layer) MoveBy(delta xgal.Point) {
	m.Bounds = m.Bounds.Add(delta)
	for _, kid := range m.Kids {
		if mv, ok := kid.(Mover); ok {
			mv.MoveBy(delta)
		}
	}
}

// Place lays out children depth-first. For Vertical (default), children stack
// downward; for Horizontal, they stack rightward. Margin is used as inner
// padding around children.
func (m *Layer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	target := m.Style.Inset(bounds)
	pad := m.Style.Margin

	total := xgal.Pt(0, 0)

	if m.Axis == Horizontal {
		for i := len(m.Kids) - 1; i >= 0; i-- {
			r := m.Kids[i].Place(target)
			kw, kh := r.Dx(), r.Dy()
			delta := xgal.Pt(kw, 0)
			target.Min = target.Min.Add(delta)
			total.X += kw
			if kh > total.Y {
				total.Y = kh
			}
		}
	} else {
		for i := len(m.Kids) - 1; i >= 0; i-- {
			r := m.Kids[i].Place(target)
			kw, kh := r.Dx(), r.Dy()
			delta := xgal.Pt(0, kh)
			target.Min = target.Min.Add(delta)
			total.Y += kh
			if kw > total.X {
				total.X = kw
			}
		}
	}
	total = total.Add(pad.Mul(2))

	m.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+total.X, bounds.Min.Y+total.Y)
	return m.Bounds
}

var _ Widget = &Layer{}
