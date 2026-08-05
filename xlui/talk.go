package xlui

import "github.com/xmasengine/xmas/xgal"

const TalkSizer = "WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW"
const TalkTick = 10

// NewTalk returns a new animated multi line text display Control.
func NewTalk(at xgal.Point, text string, lines int) *Control {
	// talk variables
	var (
		blink  = true
		cursor xgal.Point
		output [][]rune
		reveal int
	)

	talk := NewControl(at)
	talk.Text = text
	talk.State.Clicked = false
	size := talk.Style.Measure(TalkSizer)
	output = textToRuneLines(text)

	talk.Bounds = xgal.Bound(talk.Bounds.Min.X, talk.Bounds.Min.Y, size.X, size.Y*lines)
	clip := xgal.Bound(
		talk.Bounds.Min.X-talk.Style.Margin.X,
		talk.Bounds.Min.Y-talk.Style.Margin.Y,
		talk.Bounds.Dx()+talk.Style.Margin.X,
		talk.Bounds.Dy()+talk.Style.Margin.Y,
	)

	talk.Clip = &clip

	render := func(screen *xgal.Surface) {
		if talk.Clip != nil {
			sub := screen.SubImage(*talk.Clip)
			screen = sub.(*xgal.Surface)
		}

		delta := xgal.Pt(2, 0)
		style := talk.Style
		if talk.State.Hovered {
			style = style.Hovered()
		}
		if talk.State.Focused {
			style = style.Focused()
		}

		style.DrawBox(screen, talk.Bounds.Add(delta))
		// shift markers
		if len(output) > lines {
			if cursor.Y < len(output)-1 {
				xgal.Polyfill(screen, style.Fore,
					talk.Bounds.Max.X-5, talk.Bounds.Max.Y,
					talk.Bounds.Max.X, talk.Bounds.Max.Y-10,
					talk.Bounds.Max.X-10, talk.Bounds.Max.Y-10,
				)
			}
			if talk.From.Y < 0 {
				xgal.Polyfill(screen, style.Fore,
					talk.Bounds.Max.X-5, talk.Bounds.Min.Y,
					talk.Bounds.Max.X, talk.Bounds.Min.Y+10,
					talk.Bounds.Max.X-10, talk.Bounds.Min.Y+10,
				)
			}
		}

		// don't shift the box, but do shift the text and cursor
		delta = delta.Add(talk.From)
		style.Print(screen, talk.Bounds.Min.Add(delta), talk.Text[:reveal])
		// Draw cursor
		sz := talk.Style.Measure(string(output[cursor.Y][:cursor.X]))
		box := talk.Bounds
		min := box.Min.Add(delta)
		cx := min.X + style.Margin.X + sz.X
		cy := min.Y + style.Margin.Y + cursor.Y*size.Y
		ch := size.Y
		if blink {
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
		case xgal.KeyArrowUp:
			cursor.Y = max(0, cursor.Y-1)
			cursor.X = min(cursor.X, len(output[cursor.Y]))
		case xgal.KeyArrowDown:
			cursor.Y = min(cursor.Y+1, len(output)-1)
			cursor.X = min(cursor.X, len(output[cursor.Y]))
		case xgal.KeyHome:
			cursor.X = 0
		case xgal.KeyEnd:
			cursor.X = len(output[cursor.Y])
		case xgal.KeyPageUp:
			cursor.Y = 0
		case xgal.KeyPageDown:
			cursor.Y = len(output) - 1
		case xgal.KeyNumpadEnter, xgal.KeyEnter:
			if talk.Class.Entry != nil {
				talk.Text = runeLinesToText(output)
				return talk.Class.Entry(talk.Text)
			}
		default:
			return Ignore
		}
		blink = true
		talk.From = xgal.Pt(0, 0)
		if cursor.Y >= lines {
			// shift up
			talk.From.Y = (lines - cursor.Y - 1) * size.Y
		}
		lineLen := len(AreaSizer)
		if cursor.X >= len(AreaSizer) {
			// shift left
			talk.From.X = (lineLen - cursor.X - 1) * (size.X / len(AreaSizer))
		}

		return Accept
	}

	tick := func(t int64) Reply {
		if (t%TalkTick) == 0 && (reveal < len(talk.Text)) {
			if cursor.X == len(output[cursor.Y]) && cursor.Y < (len(output)-1) {
				cursor.Y++
				cursor.X = 0
			} else {
				cursor.X = min(cursor.X+1, len(output[cursor.Y]))
			}
			if cursor.Y >= lines {
				// shift up
				talk.From.Y = (lines - cursor.Y - 1) * size.Y
			}
			lineLen := len(AreaSizer)
			if cursor.X >= len(AreaSizer) {
				// shift left
				talk.From.X = (lineLen - cursor.X - 1) * (size.X / len(AreaSizer))
			}
			reveal++
		}
		return Ignore
	}

	talk.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
		Tap:     tap,
		Tick:    tick,
	}
	return talk
}

// Talk adds a new test talk to the layer.
func (l *Layer) Talk(text string, lines int) *Control {
	at := l.Bounds.Min
	ctrl := NewTalk(at, text, lines)
	return l.Append(ctrl)
}
