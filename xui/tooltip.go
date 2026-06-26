package xui

import "github.com/xmasengine/xmas/xgal"

// TooltipLayer is a small text popup that follows the mouse cursor.
// It does not consume input events.
// The tooltip is invisible until the mouse enters the Trigger rectangle.
// Once visible it follows the cursor at Offset.
// Keep it in stored the parent's Kids.
type TooltipLayer struct {
	Text    string
	Style   Style
	Offset  xgal.Point // Offset from cursor, default (12, 12)
	Bounds  xgal.Rectangle
	Trigger xgal.Rectangle // Trigger area to show tooltip when mouse is inside.
	Visible bool
}

// Tooltip returns a new [TooltipLayer] that shows text when the mouse
// enters trigger.
func Tooltip(trigger xgal.Rectangle, text string) *TooltipLayer {
	return &TooltipLayer{
		Text:    text,
		Style:   DefaultStyle(),
		Offset:  xgal.Pt(12, 12),
		Trigger: trigger,
	}
}

var _ Widget = &TooltipLayer{}

func (t *TooltipLayer) Poll() Reply {
	pos := xgal.Mouse()
	t.Visible = pos.In(t.Trigger)
	if !t.Visible {
		return Ignore
	}
	sz := t.Style.MeasureText(t.Text)
	nw := sz.X + t.Style.Margin.X*2
	nh := sz.Y + t.Style.Margin.Y*2
	t.Bounds = xgal.Rect(pos.X+t.Offset.X, pos.Y+t.Offset.Y,
		pos.X+t.Offset.X+nw, pos.Y+t.Offset.Y+nh)
	return Ignore
}

func (t *TooltipLayer) Render(s *xgal.Surface) {
	if !t.Visible {
		return
	}
	t.Style.DrawBox(s, t.Bounds)
	t.Style.Ink(s, t.Bounds, t.Text)
}

func (t *TooltipLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	sz := t.Style.MeasureText(t.Text)
	nw := sz.X + t.Style.Margin.X*2
	nh := sz.Y + t.Style.Margin.Y*2
	t.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return t.Bounds
}

func (t *TooltipLayer) MoveBy(delta xgal.Point) {
	t.Bounds = t.Bounds.Add(delta)
	t.Trigger = t.Trigger.Add(delta)
}

// AddTooltip adds a [TooltipLayer] to this layer.
func (m *Layer) AddTooltip(trigger xgal.Rectangle, text string) *TooltipLayer {
	t := Tooltip(trigger, text)
	m.Add(t)
	return t
}
