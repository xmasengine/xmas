package xui

import "image"
import "slices"
import "fmt"
import "errors"
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
	Bounds Rectangle
	Style
	Done bool
}

func MakeLayer(bounds Rectangle) Layer {
	return Layer{Bounds: bounds, Style: DefaultStyle()}
}

func (m *Layer) Poll() Reply {
	res := m.UpdateKids()
	if res != Ignore {
		return res
	}
	return Ignore
}

func (m *Layer) Render(s *Surface) {
	m.Style.DrawBox(s, m.Bounds)
	m.DrawKids(s)
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
		if res == errors.Is(err, Termination) {
			m.Kids = slices.Delete(m.Kids, i, i+1)
		} else if errors.Is(err, LayerOK) {
			break // handled by toplevel
		} else {
			return err
		}
	}
	return nil
}

func (m *Layer) DrawKids(s *Surface) {
	for i := len(m.Kids) - 1; i >= 0; i-- {
		kid := m.Kids[i]
		if kid != nil {
			kid.Draw(s)
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

func Ask(bounds Rectangle, prompt, def string, on func(res string)) *Asker {
	return &Asker{Layer: MakeLayer(bounds), Prompt: prompt, On: on, Buf: []rune(def)}
}

func (a *Asker) Update() error {
	var keys []Key
	keys = inpututil.AppendJustPressedKeys(keys)
	for _, key := range keys {
		switch key {
		case ebiten.KeyEnter:
			a.On(string(a.Buf))
			return Termination
		case ebiten.KeyEscape:
			println("esc in ", a.Prompt)
			return Termination
		case ebiten.KeyBackspace:
			if len(a.Buf) > 0 {
				a.Buf = slices.Delete(a.Buf, len(a.Buf)-1, len(a.Buf))
			}
		}
	}

	var chars []rune
	chars = ebiten.AppendInputChars(chars)
	if len(chars) > 0 {
		a.Buf = append(a.Buf, chars...)

	}
	return LayerOK
}

func (a Asker) Draw(s *Surface) {
	a.Layer.Draw(s)
	ebitenutil.DebugPrintAt(s,
		fmt.Sprintf("%s>%s|", a.Prompt, string(a.Buf)),
		a.Bounds.Min.X, a.Bounds.Min.Y)
}

func Bounds(x, y, w, h int) Rectangle {
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
