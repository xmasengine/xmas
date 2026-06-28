package xui

import "github.com/xmasengine/xmas/xgal"

// ToggleLayer is a button that stays on or off. When Group is non-nil,
// exactly one toggle in the group is active at a time.
type ToggleLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Text    string
	Icon    Icon
	Active  bool
	Group   *int // shared selection index, nil for independent toggle
	Idx     int  // this toggle is active when Group != nil && *Group == Idx
	Toggled func(active bool)

	pressed bool
	hover   bool
	lastAct bool
}

// Toggle returns a new [ToggleLayer].
func Toggle(bounds xgal.Rectangle, text string, toggled func(active bool)) *ToggleLayer {
	return &ToggleLayer{
		Bounds:  bounds,
		Style:   DefaultStyle(),
		Text:    text,
		Toggled: toggled,
	}
}

var _ Widget = &ToggleLayer{}

func (t *ToggleLayer) Poll() Reply {
	if t.Group != nil {
		t.Active = *t.Group == t.Idx
	}

	t.hover = xgal.Mouse().In(t.Bounds)

	if t.hover && xgal.Click(xgal.MouseButtonLeft) {
		t.pressed = true
	}

	if t.pressed && xgal.Loose(xgal.MouseButtonLeft) {
		t.pressed = false
		if t.Group != nil {
			*t.Group = t.Idx
			t.Active = true
		} else {
			t.Active = !t.Active
		}
		if t.Toggled != nil && t.Active != t.lastAct {
			t.Toggled(t.Active)
		}
		t.lastAct = t.Active
		return Accept
	}

	if t.hover || t.pressed {
		return Accept
	}

	return Ignore
}

func (t *ToggleLayer) Render(s *xgal.Surface) {
	box := t.Bounds
	style := t.Style

	if t.Active {
		style = style.ActiveStyle()
	} else if t.hover {
		style = style.HoverStyle()
	}
	if t.pressed {
		box = box.Add(t.Style.Margin)
		if t.Active {
			style = style.ActiveStyle()
		} else {
			style = style.PressStyle()
		}
	}

	style.DrawBox(s, box)
	t.Icon.Blit(s, box.Min)
	style.Ink(s, t.Icon.TextBounds(box), t.Text)
}

func (t *ToggleLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	sz := t.Style.MeasureText(t.Text)
	nw := sz.X + t.Style.Margin.X*2 + t.Icon.Width()
	nh := sz.Y + t.Style.Margin.Y*2
	t.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return t.Bounds
}

func (t *ToggleLayer) MoveBy(delta xgal.Point) {
	t.Bounds = t.Bounds.Add(delta)
}

// AddToggle is a helper to add a [ToggleLayer] to a [Layer].
func (m *Layer) AddToggle(bounds xgal.Rectangle, text string, toggled func(active bool)) *ToggleLayer {
	t := Toggle(bounds, text, toggled)
	m.Add(t)
	return t
}
