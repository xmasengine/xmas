package xlui

import "github.com/xmasengine/xmas/xgal"

func (l *Layer) Toggle(text string, group *Group) *Control {
	at := l.Bounds.Min
	ctrl := NewToggle(at, text, group)
	if group != nil {
		group.Controls = append(group.Controls, ctrl)
	}
	return l.Append(ctrl)
}

// Toggle is a toggle that stays ON or OFF.
// If grouped only one of the group can be ON.
func NewToggle(at xgal.Point, text string, group *Group) *Control {
	toggle := NewControlWithText(at, ButtonStyle(), text)
	toggle.State.Clicked = false
	toggle.Checked = false

	render := func(screen *xgal.Surface) {
		delta := xgal.Pt(0, 0)
		if toggle.Checked {
			delta = xgal.Pt(2, 2)
		}
		style := toggle.Style
		if toggle.State.Hovered {
			style = style.Hovered()
		}
		style.DrawBox(screen, toggle.Bounds.Add(delta))
		style.Print(screen, toggle.Bounds.Min.Add(delta), toggle.Text)
	}

	click := func(at xgal.Point, which int) Reply {
		if group != nil {
			for _, ctrl := range group.Controls {
				ctrl.Checked = ctrl == toggle
			}
		} else {
			toggle.Checked = !toggle.Checked
		}
		return Accept
	}

	release := func(at xgal.Point, which int) Reply {
		return Accept
	}

	hover := func(at xgal.Point) Reply {
		return Accept
	}

	toggle.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
	}
	return toggle
}
