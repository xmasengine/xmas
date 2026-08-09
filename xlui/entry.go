package xlui

import "slices"

import "github.com/xmasengine/xmas/xgal"

const EntrySizer = "WWWWWWWW"
const EntryBlink = 60

// NewEntry returns a new text entry Control.
func NewEntry(at xgal.Point, text string) *Control {
	entry := NewControl(at)
	entry.Text = text
	entry.State.Clicked = false
	size := entry.Style.Measure(EntrySizer)
	entry.Bounds = xgal.Bound(entry.Bounds.Min.X, entry.Bounds.Min.Y, size.X, size.Y)
	// entry variables
	var (
		blink  = true
		cursor int
		input  []rune
	)

	if len(entry.Text) > 0 {
		input = []rune(entry.Text)
		cursor = len(input)
	}

	render := func(screen *xgal.Surface) {
		delta := xgal.Pt(2, 0)
		style := entry.Style
		if entry.State.Hovered {
			style = style.Hovered()
		}
		if entry.State.Focused {
			style = style.Focused()
		}

		style.DrawBox(screen, entry.Bounds.Add(delta))
		style.Print(screen, entry.Bounds.Min.Add(delta), entry.Text)
		// Draw cursor if focused.
		if entry.State.Focused {
			sz := entry.Style.Measure(entry.Text[:cursor])
			box := entry.Bounds
			min := box.Min.Add(delta)
			cx := min.X + style.Margin.X + sz.X
			cy := min.Y + style.Margin.Y
			ch := box.Dy() - style.Margin.Y*2
			if blink {
				xgal.Line(screen, cx, cy, cx, cy+ch, style.Stroke, style.Fore)
			}
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

	chars := func(chrs ...rune) Reply {
		if len(chrs) > 0 {
			input = slices.Insert(input, cursor, chrs...)
			cursor += len(chrs)
			entry.Text = string(input)
			return Accept
		}
		return Ignore
	}

	tick := func(t int64) Reply {
		if (t % EntryBlink) == 0 {
			blink = !blink
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
		Tick:    tick,
	}
	return entry
}

// Entry adds a new test entry to the layer.
func (l *Layer) Entry(text string) *Control {
	at := l.Bounds.Min
	ctrl := NewEntry(at, text)
	return l.Append(ctrl)
}
