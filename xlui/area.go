package xlui

import "slices"
import "strings"

import "github.com/xmasengine/xmas/xgal"

const AreaSizer = "WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW"
const AreaBlink = 60

func runeLinesToText(input [][]rune) string {
	lines := make([]string, len(input))
	for i, in := range input {
		if in == nil {
			lines[i] = ""
		} else {
			lines[i] = string(in)
		}
	}
	return strings.Join(lines, "\n")
}

func textToRuneLines(text string) [][]rune {
	lines := strings.Split(text, "\n")
	res := make([][]rune, len(lines))
	for i, line := range lines {
		res[i] = []rune(line)
	}
	return res
}

// NewArea returns a new multi line text area Control.
func NewArea(at xgal.Point, text string, lines int) *Control {
	// area variables
	var (
		blink  = true
		cursor xgal.Point
		input  [][]rune
	)

	area := NewControl(at)
	area.Text = text
	area.State.Clicked = false
	size := area.Style.Measure(AreaSizer)
	input = textToRuneLines(text)

	area.Bounds = xgal.Bound(area.Bounds.Min.X, area.Bounds.Min.Y, size.X, size.Y*lines)
	clip := xgal.Bound(
		area.Bounds.Min.X-area.Style.Margin.X,
		area.Bounds.Min.Y-area.Style.Margin.Y,
		area.Bounds.Dx()+area.Style.Margin.X,
		area.Bounds.Dy()+area.Style.Margin.Y,
	)

	area.Clip = &clip

	render := func(screen *xgal.Surface) {
		if area.Clip != nil {
			sub := screen.SubImage(*area.Clip)
			screen = sub.(*xgal.Surface)
		}

		delta := xgal.Pt(2, 0)
		style := area.Style
		if area.State.Hovered {
			style = style.Hovered()
		}
		if area.State.Focused {
			style = style.Focused()
		}

		style.DrawBox(screen, area.Bounds.Add(delta))
		// shift markers
		if len(input) > lines {
			if cursor.Y < len(input)-1 {
				xgal.Polyfill(screen, style.Fore,
					area.Bounds.Max.X-5, area.Bounds.Max.Y,
					area.Bounds.Max.X, area.Bounds.Max.Y-10,
					area.Bounds.Max.X-10, area.Bounds.Max.Y-10,
				)
			}
			if area.From.Y < 0 {
				xgal.Polyfill(screen, style.Fore,
					area.Bounds.Max.X-5, area.Bounds.Min.Y,
					area.Bounds.Max.X, area.Bounds.Min.Y+10,
					area.Bounds.Max.X-10, area.Bounds.Min.Y+10,
				)
			}
		}

		// don't shift the box, but do shift the text and cursor
		delta = delta.Add(area.From)
		style.Print(screen, area.Bounds.Min.Add(delta), area.Text)
		// Draw cursor if focused.
		if area.State.Focused {
			sz := area.Style.Measure(string(input[cursor.Y][:cursor.X]))
			box := area.Bounds
			min := box.Min.Add(delta)
			cx := min.X + style.Margin.X + sz.X
			cy := min.Y + style.Margin.Y + cursor.Y*size.Y
			ch := size.Y
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
			if cursor.X == 0 && cursor.Y > 0 {
				cursor.Y--
				cursor.X = len(input[cursor.Y])
			} else {
				cursor.X = max(0, cursor.X-1)
			}
		case xgal.KeyArrowRight:
			if cursor.X == len(input[cursor.Y]) && cursor.Y < (len(input)-1) {
				cursor.Y++
				cursor.X = 0
			} else {
				cursor.X = min(cursor.X+1, len(input[cursor.Y]))
			}
		case xgal.KeyArrowUp:
			cursor.Y = max(0, cursor.Y-1)
			cursor.X = min(cursor.X, len(input[cursor.Y]))
		case xgal.KeyArrowDown:
			cursor.Y = min(cursor.Y+1, len(input)-1)
			cursor.X = min(cursor.X, len(input[cursor.Y]))
		case xgal.KeyHome:
			cursor.X = 0
		case xgal.KeyEnd:
			cursor.X = len(input[cursor.Y])
		case xgal.KeyPageUp:
			cursor.Y = 0
		case xgal.KeyPageDown:
			cursor.Y = len(input) - 1
		case xgal.KeyNumpadEnter:
			if area.Class.Entry != nil {
				area.Text = runeLinesToText(input)
				return area.Class.Entry(area.Text)
			}
		case xgal.KeyBackspace:
			if cursor.X > 0 {
				input[cursor.Y] = slices.Delete(input[cursor.Y], cursor.X-1, cursor.X)
				cursor.X--
			} else if cursor.Y > 0 {
				cursor.X = len(input[cursor.Y-1])
				line := append(input[cursor.Y-1], input[cursor.Y]...)
				input = slices.Delete(input, cursor.Y-1, cursor.Y)
				cursor.Y--
				input[cursor.Y] = line
			}
		case xgal.KeyDelete:
			if cursor.X < len(input[cursor.Y]) {
				input[cursor.Y] = slices.Delete(input[cursor.Y], cursor.X, cursor.X+1)
			} else if cursor.Y < len(input)-1 {
				// cursor.X = len(input[cursor.Y+1])
				line := append(input[cursor.Y], input[cursor.Y+1]...)
				input = slices.Delete(input, cursor.Y, cursor.Y+1)
				input[cursor.Y] = line
			}
		case xgal.KeyEnter:
			line := input[cursor.Y]
			before, after := slices.Clone(line[:cursor.X]), slices.Clone(line[cursor.X:])
			input = slices.Insert(input, cursor.Y, before)
			cursor.Y++
			input[cursor.Y] = after
			cursor.X = 0
		default:
			return Ignore
		}
		blink = true
		area.Text = runeLinesToText(input)
		area.From = xgal.Pt(0, 0)
		if cursor.Y >= lines {
			// shift up
			area.From.Y = (lines - cursor.Y - 1) * size.Y
		}
		lineLen := len(AreaSizer)
		if cursor.X >= len(AreaSizer) {
			// shift left
			area.From.X = (lineLen - cursor.X - 1) * (size.X / len(AreaSizer))
		}

		return Accept
	}

	chars := func(chrs ...rune) Reply {
		if len(chrs) > 0 {
			input[cursor.Y] = slices.Insert(input[cursor.Y], cursor.X, chrs...)
			cursor.X += len(chrs)
			lineLen := len(AreaSizer)
			if cursor.X >= len(AreaSizer) {
				// shift left
				area.From.X = (lineLen - cursor.X - 1) * (size.X / len(AreaSizer))
			}
			area.Text = runeLinesToText(input)
			return Accept
		}
		return Ignore
	}

	tick := func(t int64) Reply {
		if (t % AreaBlink) == 0 {
			blink = !blink
		}
		return Ignore
	}

	area.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
		Tap:     tap,
		Chars:   chars,
		Tick:    tick,
	}
	return area
}

// Area adds a new test area to the layer.
func (l *Layer) Area(text string, lines int) *Control {
	at := l.Bounds.Min
	ctrl := NewArea(at, text, lines)
	return l.Append(ctrl)
}
