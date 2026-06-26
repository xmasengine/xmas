package xui

import "github.com/xmasengine/xmas/xgal"

const CheckboxSize = 8

// CheckboxLayer is a toggleable checkbox with a label.
type CheckboxLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Text    string
	Checked bool
	OnCheck func(bool)
	hover   bool
	pressed bool
}

// Checkbox returns a new [CheckboxLayer] with the given bounds, text, and
// toggle callback.
func Checkbox(bounds xgal.Rectangle, text string, onCheck func(bool)) *CheckboxLayer {
	return &CheckboxLayer{
		Bounds:  bounds,
		Style:   DefaultStyle(),
		Text:    text,
		OnCheck: onCheck,
	}
}

var _ Widget = &CheckboxLayer{}

func (c *CheckboxLayer) Poll() Reply {
	c.hover = xgal.Mouse().In(c.Bounds)

	if c.hover && xgal.Click(xgal.MouseButtonLeft) {
		c.pressed = true
		return Accept
	}

	if c.pressed {
		if xgal.Loose(xgal.MouseButtonLeft) {
			c.pressed = false
			if c.hover {
				c.Checked = !c.Checked
				if c.OnCheck != nil {
					c.OnCheck(c.Checked)
				}
			}
		}
		return Accept
	}

	if c.hover {
		return Accept
	}

	return Ignore
}

func (c *CheckboxLayer) Render(s *xgal.Surface) {
	box := c.Bounds
	style := c.Style

	if c.hover {
		style = style.HoverStyle()
	}

	style.DrawBox(s, box)

	cy := box.Min.Y + (box.Dy()-CheckboxSize)/2
	ibox := xgal.Rect(box.Min.X+style.Margin.X, cy, box.Min.X+style.Margin.X+CheckboxSize, cy+CheckboxSize)

	cstyle := style.CheckStyle()
	if c.Checked {
		cstyle.DrawBox(s, ibox)
	} else {
		cstyle.DrawRect(s, ibox)
	}

	at := xgal.Pt(ibox.Max.X+style.Margin.X, box.Min.Y)
	style.DrawText(s, at, c.Text)
}

func (c *CheckboxLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	tsz := c.Style.MeasureText(c.Text)
	nw := c.Style.Margin.X + CheckboxSize + c.Style.Margin.X + tsz.X + c.Style.Margin.X
	nh := tsz.Y + c.Style.Margin.Y*2
	if nh < CheckboxSize+c.Style.Margin.Y*2 {
		nh = CheckboxSize + c.Style.Margin.Y*2
	}
	c.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return c.Bounds
}

func (c *CheckboxLayer) MoveBy(delta xgal.Point) {
	c.Bounds = c.Bounds.Add(delta)
}

func (m *Layer) AddCheckbox(bounds xgal.Rectangle, text string, onCheck func(bool)) *CheckboxLayer {
	c := Checkbox(bounds, text, onCheck)
	m.Add(c)
	return c
}
