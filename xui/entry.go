package xui

import "image"
import "slices"
import "strconv"
import "log/slog"

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Entry struct {
	Widget
	Change func(*Entry)
	cursor Point
	id     int
	Input  []rune
	Label  string
	Icon   *Icon
}

type EntryClass struct {
	*Entry
	*WidgetClass
}

func (e *Entry) SetLabel(label string) *Entry {
	e.Label = label
	return e
}

func (e *Entry) Text() string {
	return string(e.Input)
}

func (e *Entry) Init(bounds Rectangle, text string, ch func(*Entry)) *Entry {
	e.Change = ch
	e.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
	e.insertChars([]rune(text)...)
	e.Class = NewEntryClass(e)
	return e
}

func NewEntry(bounds Rectangle, text string, ch func(*Entry)) *Entry {
	e := &Entry{}
	return e.Init(bounds, text, ch)
}

func NewEntryClass(c *Entry) *EntryClass {
	res := &EntryClass{Entry: c}
	res.WidgetClass = NewWidgetClass()
	return res
}

func (ec EntryClass) Render(root *Root, screen *Surface) {
	e := ec.Entry
	box := e.Bounds

	clipped := screen.SubImage(e.Bounds.Inset(int(-e.Style.Stroke))).(*Surface)
	style := e.Style
	if e.State.Hover {
		style = HoverStyle()
	}
	if e.State.Focus {
		style = FocusStyle()
	}

	size := style.MeasureText(string(e.Input[:e.cursor.X]))
	if size.Y == 0 { // In case there is no text
		size.Y = e.Style.LineHeight()
	}
	cursorBox := box
	cursorBox.Min = cursorBox.Min.Add(image.Pt(size.X+e.Style.Margin.X, int(e.Style.Margin.Y/2)))
	cursorBox.Max = cursorBox.Min.Add(image.Pt(int(e.Style.Stroke), size.Y-int(e.Style.Margin.Y/2)))

	cursorColor := e.Style.Writing
	cursorColor.A = cursorColor.A / 4
	DrawLine(clipped, cursorBox, 1, cursorColor)
	style.DrawBox(screen, box)
	style.DrawText(clipped, box.Min, string(e.Input))

	if e.Icon != nil {
		ibox := image.Rect(box.Max.X-IconSize, box.Max.Y-IconSize, box.Max.X, box.Max.Y)
		e.Icon.Draw(clipped, ibox)
	}

	if e.Label != "" {
		// Draw the label to the right, outside the box, to simplify the label layout.
		at := image.Pt(box.Max.X+e.Style.Margin.X, box.Min.Y)
		style.DrawText(screen, at, e.Label)
	}
}

func (b *EntryClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *EntryClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}

func (b *EntryClass) OnActionFocus(e ActionEvent) bool {
	b.State.Focus = true
	dprintln("entry focused")
	return true
}

func (b *EntryClass) OnActionBlur(e ActionEvent) bool {
	b.State.Focus = false
	return true
}

func (e *EntryClass) MousePress(ev MouseEvent) bool {
	return true
}

func (e *EntryClass) MouseRelease(ev MouseEvent) bool {
	return true
}

func (e *EntryClass) OnKeyPress(ev KeyEvent) bool {
	slog.Debug("EntryClass.OnKeyPress", "ev", ev)

	switch ebiten.Key(ev.Code) {
	case ebiten.KeyF2:
		e.State.Focus = false
		ev.Root().Focus = nil
	case ebiten.KeyLeft:
		e.cursor.X = max(0, e.cursor.X-1)
	case ebiten.KeyRight:
		e.cursor.X = min(e.cursor.X+1, len(e.Input))
	case ebiten.KeyHome:
		e.cursor.X = 0
	case ebiten.KeyEnd:
		e.cursor.X = len(e.Input)
	case ebiten.KeyEnter:
		if e.Change != nil {
			e.Change(e.Entry)
		}
	case ebiten.KeyBackspace:
		if e.cursor.X > 0 {
			e.Input = slices.Delete(e.Input, e.cursor.X-1, e.cursor.X)
			e.cursor.X--
		}
	case ebiten.KeyDelete:
		if (e.cursor.X + 1) < len(e.Input) {
			e.Input = slices.Delete(e.Input, e.cursor.X, e.cursor.X+1)
		}
	/*
		case ebiten.KeyC:
			if ev.Modifiers.Control {
				WriteClipboard(ClipboardText, []byte(string(e.Input)))
				dprintln("Entry copy to clipboard.")

			}
		case ebiten.KeyV:
			if ev.Modifiers.Control {
				buf := ReadClipboard(ClipboardText)
				e.insertChars([]rune(string(buf)))
				dprintln("Entry read from clipboard.", string(buf))
			}
	*/
	default:
		return false
	}
	return true
}

func (e *EntryClass) OnKeyHold(ev KeyEvent) bool {
	if ev.Duration > 60 && (ev.Duration%30) == 0 {
		return e.OnKeyPress(ev)
	}
	return false
}

func (e *EntryClass) OnKeyText(ev KeyEvent) bool {
	slog.Debug("EntryClass.OnKeyText", "ev", ev)
	if len(ev.Chars) > 0 && (ev.ID == e.id || ev.ID < 0) {
		e.insertChars(ev.Chars...)
	}
	return true
}

func (e *Entry) insertChars(chars ...rune) {
	e.Input = slices.Insert(e.Input, e.cursor.X, chars...)
	e.cursor.X += len(chars)
}

func (p *Widget) AddEntry(bounds Rectangle, text string, ch func(*Entry)) *Entry {
	entry := NewEntry(bounds, text, ch)
	p.Widgets = append(p.Widgets, &entry.Widget)
	return entry
}

func (p *Widget) StringEntry(bounds Rectangle, bound *string) *Entry {
	if bound == nil {
		panic("Incorrect use of StringEntry")
	}
	ch := func(e *Entry) {
		*bound = string(e.Input)
	}
	entry := p.AddEntry(bounds, *bound, ch)
	return entry
}

func (p *Widget) IntEntry(bounds Rectangle, bound *int) *Entry {
	if bound == nil {
		panic("Incorrect use of IntEntry")
	}
	text := strconv.Itoa(*bound)
	ch := func(e *Entry) {
		*bound, _ = strconv.Atoi(string(e.Input))
	}
	entry := p.AddEntry(bounds, text, ch)
	return entry
}
