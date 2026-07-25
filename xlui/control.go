package xlui

import "github.com/xmasengine/xmas/xgal"

// Control is a possibly interactieve part of the UI inside a layer.
type Control struct {
	// Class for custom or type specific behavior.
	Class
	// Data
	Text   string // For use by text controls.
	Bounds xgal.Rectangle
	Clip   *xgal.Rectangle
	Style
	From        xgal.Point
	Orientation Orientation // layout orientation in the layer
}

func (c Control) Render(s *xgal.Surface) {
	if c.Class.Render == nil {
		c.Style.DrawBox(s, c.Bounds)
		return
	}
	c.Class.Render(s)
}

// MoveBy moves the control.
func (c *Control) MoveBy(delta xgal.Point) {
	c.Bounds = c.Bounds.Add(delta)
}

const (
	ControlWidth  = 20
	ControlHeight = 10
)

func NewControl(at xgal.Point) *Control {
	return &Control{Bounds: xgal.Bound(at.X, at.Y, ControlWidth, ControlHeight), Style: DefaultStyle(), Orientation: Vertical}
}

func NewControlWithText(at xgal.Point, style Style, text string) *Control {
	ctrl := NewControl(at)
	ctrl.Style = ButtonStyle()
	ctrl.Text = text
	if ctrl.Text != "" {
		size := ctrl.Style.Measure(ctrl.Text)
		ctrl.Bounds = xgal.Bound(at.X, at.Y, size.X, size.Y)
	}
	return ctrl
}

func NewLabel(at xgal.Point, text string) *Control {
	label := NewControlWithText(at, ButtonStyle(), text)
	render := func(screen *xgal.Surface) {
		label.Style.Print(screen, label.Bounds.Min, label.Text)
	}
	label.Class = Class{
		Render: render,
	}
	return label
}

func NewButton(at xgal.Point, text string) *Control {
	button := NewControlWithText(at, ButtonStyle(), text)
	clicked := false

	render := func(screen *xgal.Surface) {
		delta := xgal.Pt(0, 0)
		if clicked {
			delta = xgal.Pt(2, 2)
		}
		button.Style.DrawBox(screen, button.Bounds.Add(delta))
		button.Style.Print(screen, button.Bounds.Min.Add(delta), button.Text)
	}

	click := func(at xgal.Point, which int) Reply {
		clicked = true
		println("Click", button.Text, clicked, which)
		return Accept
	}

	button.Class = Class{
		Render: render,
		Click:  click,
	}
	return button
}
