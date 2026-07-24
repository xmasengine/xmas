package xlui

import "github.com/xmasengine/xmas/xgal"

// Control is a possibly interactieve part of the UI inside a layer.
type Control struct {
	Bounds xgal.Rectangle
	Clip   *xgal.Rectangle
	Style
	From        xgal.Point
	Orientation Orientation // layout orientation in the layer

	// Flexible handlers.
	OnRender func(s *xgal.Surface, c Control)
}

func (c Control) Render(s *xgal.Surface) {
	if c.OnRender == nil {
		c.Style.DrawBox(s, c.Bounds)
		return
	}
	c.OnRender(s, c)
}

// MoveBy moves the control.
func (c *Control) MoveBy(delta xgal.Point) {
	c.Bounds = c.Bounds.Add(delta)
}
