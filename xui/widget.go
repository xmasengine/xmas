package xui

import "image"
import "slices"
import "fmt"
import "github.com/xmasengine/xmas/xgal"

type Reply int

const (
	Ignore Reply = iota
	Accept
	Raise
	Lower
	Finish
)

type Widget interface {
	Poll() Reply
	Render(screen *xgal.Surface)
	Place(w, h int) (myw, myh int)
}

// Layer is a layer in the UI.
// Layers may have child layers.
// It simply implements the Widget interface
type Layer struct {
	Kids   []Widget
	Bounds xgal.Rectangle
	Style
	Done bool
}

func MakeLayer(bounds xgal.Rectangle) Layer {
	return Layer{Bounds: bounds, Style: DefaultStyle()}
}

func (m *Layer) Poll() Reply {
	res := m.PollKids()
	if res != Ignore {
		return res
	}
	return Ignore
}

func (m *Layer) Render(s *xgal.Surface) {
	m.Style.DrawBox(s, m.Bounds)
	m.RenderKids(s)
}

func (m *Layer) Add(g Widget) Widget {
	m.Kids = slices.Insert(m.Kids, 0, g)
	return g
}

func (m *Layer) PollKids() Reply {
	for i, kid := range m.Kids {
		if kid == nil {
			continue
		}
		res := kid.Poll()
		if res == Finish {
			m.Kids = slices.Delete(m.Kids, i, i+1)
		} else if res == Accept {
			break // handled by toplevel
		} else {
			return res
		}
	}
	return Ignore
}

func (m *Layer) RenderKids(s *xgal.Surface) {
	for i := len(m.Kids) - 1; i >= 0; i-- {
		kid := m.Kids[i]
		if kid != nil {
			kid.Render(s)
		}
	}
}

func (m *Layer) Place(w, h int) (rw, rh int) {
	m.Bounds.Max = m.Bounds.Min.Add(image.Pt(w, h))
	return m.Bounds.Dx(), m.Bounds.Dy()
}

var _ Widget = &Layer{}

type Asker struct {
	Layer
	Prompt string
	Buf    []rune
	On     func(string)
	Cursor int
}

func Ask(bounds xgal.Rectangle, prompt, def string, on func(res string)) *Asker {
	return &Asker{Layer: MakeLayer(bounds), Prompt: prompt, On: on, Buf: []rune(def)}
}

func (a *Asker) Poll() Reply {
	var keys []xgal.KeyCode
	keys = xgal.Taps(keys)
	for _, key := range keys {
		switch key {
		case xgal.KeyEnter:
			a.On(string(a.Buf))
			return Finish
		case xgal.KeyEscape:
			println("esc in ", a.Prompt)
			return Finish
		case xgal.KeyBackspace:
			if len(a.Buf) > 0 {
				a.Buf = slices.Delete(a.Buf, len(a.Buf)-1, len(a.Buf))
			}
		}
	}

	var chars []rune
	chars = xgal.Chars(chars)
	if len(chars) > 0 {
		a.Buf = append(a.Buf, chars...)

	}
	return Accept
}

func (a Asker) Draw(s *xgal.Surface) {
	a.Layer.Render(s)
	xgal.Ink(s, a.Style.Face, a.Style.Fore,
		a.Bounds.Min.X, a.Bounds.Min.Y,
		fmt.Sprintf("%s>%s|", a.Prompt, string(a.Buf)))
}

func Bounds(x, y, w, h int) xgal.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

func (m *Layer) Ask(x, y, w, h int, prompt, def string, on func(res string)) *Asker {
	ask := Ask(Bounds(x, y, w, h), prompt, def, on)
	m.Add(ask)
	return ask
}

func (m *Layer) YesNo(x, y, w, h int, prompt, def string, on func(res bool)) *Asker {
	wrap := func(sres string) {
		on(sres == def)
	}
	ask := Ask(Bounds(x, y, w, h), prompt, def, wrap)
	m.Add(ask)
	return ask
}
