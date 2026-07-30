package xlui

import "github.com/xmasengine/xmas/xgal"

const CheckboxSize = 8
const CheckboxSizer = "X"

func (l *Layer) Checkbox(checked bool) *Control {
	at := l.Bounds.Min
	ctrl := NewCheckbox(at, checked)
	return l.Append(ctrl)
}

func (l *Layer) CheckboxWithLabel(checked bool, text string) (*Control, *Control) {
	checkbox := l.Checkbox(checked)
	label := l.Label(text)
	label.Class.Click = func(at xgal.Point, button int) Reply {
		return checkbox.Class.Click(at, button)
	}

	return checkbox, label
}

func NewCheckbox(at xgal.Point, checked bool) *Control {
	checkbox := NewControl(at)
	checkbox.State.Clicked = false
	checkbox.Checked = checked
	size := checkbox.Style.Measure(CheckboxSizer)
	size.X = max(size.X, CheckboxSize)
	size.Y = max(size.Y, CheckboxSize)

	checkbox.Bounds = xgal.Bound(checkbox.Bounds.Min.X, checkbox.Bounds.Min.Y, size.X, size.Y)

	render := func(screen *xgal.Surface) {
		style := checkbox.Style
		if checkbox.State.Hovered {
			style = style.Hovered()
		}
		box := checkbox.Bounds
		cy := box.Min.Y + (box.Dy()-CheckboxSize)/2
		ibox := xgal.Rect(box.Min.X+style.Margin.X, cy, box.Min.X+style.Margin.X+CheckboxSize, cy+CheckboxSize)

		cstyle := style.CheckStyle()
		if checkbox.Checked {
			cstyle.DrawBox(screen, ibox)
		} else {
			style.DrawRect(screen, ibox)
		}
	}

	click := func(at xgal.Point, which int) Reply {
		checkbox.Checked = !checkbox.Checked
		return Accept
	}

	release := func(at xgal.Point, which int) Reply {
		return Accept
	}

	hover := func(at xgal.Point) Reply {
		return Accept
	}

	checkbox.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
	}
	return checkbox
}
