package xui

import (
	"slices"

	"github.com/xmasengine/xmas/xgal"
)

// EntryLayer is a text input field.
type EntryLayer struct {
	Bounds xgal.Rectangle
	Style  Style
	Input  []rune
	Label  string
	Change func(string)
	cursor int
	focus  bool
	hover  bool
}

// Entry returns a new [EntryLayer] with the given bounds, initial text, and
// change callback. The callback is called when Enter is pressed.
func Entry(bounds xgal.Rectangle, text string, change func(string)) *EntryLayer {
	return &EntryLayer{
		Bounds: bounds,
		Style:  DefaultStyle(),
		Input:  []rune(text),
		Change: change,
	}
}

var _ Widget = &EntryLayer{}

func (e *EntryLayer) Text() string {
	return string(e.Input)
}

func (e *EntryLayer) Poll() Reply {
	e.hover = xgal.Mouse().In(e.Bounds)

	if xgal.Click(xgal.MouseButtonLeft) {
		e.focus = e.hover
		if e.focus {
			return Accept
		}
	}

	if !e.focus {
		return Ignore
	}

	for _, k := range xgal.Taps(nil) {
		switch k {
		case xgal.KeyArrowLeft:
			e.cursor = max(0, e.cursor-1)
		case xgal.KeyArrowRight:
			e.cursor = min(e.cursor+1, len(e.Input))
		case xgal.KeyHome:
			e.cursor = 0
		case xgal.KeyEnd:
			e.cursor = len(e.Input)
		case xgal.KeyEnter:
			if e.Change != nil {
				e.Change(string(e.Input))
			}
		case xgal.KeyBackspace:
			if e.cursor > 0 {
				e.Input = slices.Delete(e.Input, e.cursor-1, e.cursor)
				e.cursor--
			}
		case xgal.KeyDelete:
			if e.cursor < len(e.Input) {
				e.Input = slices.Delete(e.Input, e.cursor, e.cursor+1)
			}
		}
	}

	chars := xgal.Chars(nil)
	if len(chars) > 0 {
		e.Input = slices.Insert(e.Input, e.cursor, chars...)
		e.cursor += len(chars)
	}

	return Accept
}

func (e *EntryLayer) Render(s *xgal.Surface) {
	box := e.Bounds
	style := e.Style

	if e.focus {
		style = style.FocusStyle()
	} else if e.hover {
		style = style.HoverStyle()
	}

	style.DrawBox(s, box)

	txt := string(e.Input)
	style.Ink(s, box, txt)

	// Draw the cursor if focused.
	if e.focus {
		sz := style.MeasureText(txt[:e.cursor])
		cx := box.Min.X + style.Margin.X + sz.X
		cy := box.Min.Y + style.Margin.Y
		ch := box.Dy() - style.Margin.Y*2
		xgal.Line(s, cx, cy, cx, cy+ch, style.Stroke, style.Fore)
	}

	if e.Label != "" {
		at := xgal.Pt(box.Max.X+style.Margin.X, box.Min.Y)
		style.DrawText(s, at, e.Label)
	}
}

const minEntryW = 80
const minEntryH = 16

func (e *EntryLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	sz := e.Style.MeasureText(string(e.Input) + "  ")
	nw := sz.X + e.Style.Margin.X*2
	if nw < minEntryW {
		nw = minEntryW
	}
	nh := sz.Y + e.Style.Margin.Y*2
	if nh < minEntryH {
		nh = minEntryH
	}
	e.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return e.Bounds
}

func (e *EntryLayer) MoveBy(delta xgal.Point) {
	e.Bounds = e.Bounds.Add(delta)
}

func (m *Layer) AddEntry(bounds xgal.Rectangle, text string, change func(string)) *EntryLayer {
	e := Entry(bounds, text, change)
	m.Add(e)
	return e
}
