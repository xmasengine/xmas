package xlui

import "github.com/xmasengine/xmas/xgal"

// State is the state of the control or layer.
type State struct {
	Hovered bool
	Clicked bool
	Focused bool
}

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
	State       State       // State of the control
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
	button.State.Clicked = false

	render := func(screen *xgal.Surface) {
		delta := xgal.Pt(0, 0)
		if button.State.Clicked {
			delta = xgal.Pt(2, 2)
		}
		style := button.Style
		if button.State.Hovered {
			style = style.Hovered()
		}
		style.DrawBox(screen, button.Bounds.Add(delta))
		style.Print(screen, button.Bounds.Min.Add(delta), button.Text)
	}

	click := func(at xgal.Point, which int) Reply {
		return Accept
	}

	release := func(at xgal.Point, which int) Reply {
		return Accept
	}

	hover := func(at xgal.Point) Reply {
		return Accept
	}

	button.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
	}
	return button
}
