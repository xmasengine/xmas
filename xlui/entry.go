package xlui

import "slices"

import "github.com/xmasengine/xmas/xgal"

func NewEntry(at xgal.Point, text string) *Control {
	entry := NewControl(at)
	entry.State.Clicked = false

	var cursor int
	var input []rune

	render := func(screen *xgal.Surface) {
		delta := xgal.Pt(0, 0)
		if entry.State.Clicked {
			delta = xgal.Pt(2, 2)
		}
		style := entry.Style
		if entry.State.Hovered {
			style = style.Hovered()
		}
		style.DrawBox(screen, entry.Bounds.Add(delta))
		style.Print(screen, entry.Bounds.Min.Add(delta), entry.Text)
		// Draw cursor if focused.
		if entry.State.Focused {
			sz := entry.Style.Measure(entry.Text[:cursor])
			box := entry.Bounds
			cx := box.Min.X + style.Margin.X + sz.X
			cy := box.Min.Y + style.Margin.Y
			ch := box.Dy() - style.Margin.Y*2
			xgal.Line(screen, cx, cy, cx, cy+ch, style.Stroke, style.Fore)
		}
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

	tap := func(key int, mods Mods) Reply {
		switch xgal.KeyCode(key) {
		case xgal.KeyArrowLeft:
			cursor = max(0, cursor-1)
		case xgal.KeyArrowRight:
			cursor = min(cursor+1, len(input))
		case xgal.KeyHome:
			cursor = 0
		case xgal.KeyEnd:
			cursor = len(input)
		case xgal.KeyEnter:
			if entry.Class.Entry != nil {
				return entry.Class.Entry(string(input))
			}
		case xgal.KeyBackspace:
			if cursor > 0 {
				input = slices.Delete(input, cursor-1, cursor)
				cursor--
			}
		case xgal.KeyDelete:
			if cursor < len(input) {
				input = slices.Delete(input, cursor, cursor+1)
			}
		default:
			return Ignore
		}
		entry.Text = string(input)
		return Accept
	}

	chars := func(chars ...rune) Reply {
		if len(chars) > 0 {
			input = slices.Insert(input, cursor, chars...)
			cursor += len(chars)
			entry.Text = string(input)
			return Accept
		}
		return Ignore
	}

	entry.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
		Tap:     tap,
		Chars:   chars,
	}
	return entry
}

func (l *Layer) Entry(text string) *Control {
	at := l.Bounds.Min
	ctrl := NewEntry(at, text)
	return l.Append(ctrl)
}
