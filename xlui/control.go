package xlui

import "github.com/xmasengine/xmas/xgal"

// Control is a possibly interactieve part of the UI inside a layer.
type Control struct {
	// Class for custom or type specific behavior.
	Class
	// Data
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

type LabelControl struct {
	*Control
	Text string
}

func (l *LabelControl) Render(screen *xgal.Surface) {
	l.Control.Style.Print(screen, l.Control.Bounds.Min, l.Text)
}

func (l *LabelControl) Class() Class {
	return Class{
		Render: l.Render,
	}
}

func NewLabelControl(at xgal.Point, text string, style Style) *LabelControl {
	ctrl := NewControl(at)
	ctrl.Style = style
	size := ctrl.Style.Measure(text)
	ctrl.Bounds = xgal.Bound(at.X, at.Y, size.X, size.Y)
	label := &LabelControl{Text: text, Control: ctrl}
	ctrl.Class = label.Class()
	return label
}

func NewLabel(at xgal.Point, text string) *Control {
	label := NewLabelControl(at, text, DefaultStyle())
	return label.Control
}

type ButtonControl struct {
	*LabelControl
}

func (l *ButtonControl) Render(screen *xgal.Surface) {
	l.Control.Style.DrawBox(screen, l.Bounds)
	l.Control.Style.Print(screen, l.Control.Bounds.Min, l.Text)
}

func (l *ButtonControl) Class() Class {
	return Class{
		Render: l.Render,
	}
}

func NewButtonControl(at xgal.Point, text string) *ButtonControl {
	label := NewLabelControl(at, text, ButtonStyle())
	button := &ButtonControl{LabelControl: label}
	label.Control.Bounds = button.Control.Bounds.Add(label.Style.Margin)
	label.Control.Class = button.Class()
	return button
}

func NewButton(at xgal.Point, text string) *Control {
	button := NewButtonControl(at, text)
	return button.Control
}
