package xui

import "github.com/xmasengine/xmas/xgal"

// AskLayer is a simple dialog with a prompt and one or more text buttons
// and an optional Entry.
// It returns Finish on Poll once a button is clicked, and sets Result to
// the index of the clicked button (0‑based).
type AskLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Prompt  string
	Buttons []string
	Entry   *EntryLayer
	Result  int // -1 while open
	hover   int // index of hovered button, -1 for none
}

// Ask returns a new [AskLayer]. The caller must add it to a container
// (e.g. via [Layer.Add]). After Poll returns Finish, read Result to see which
// button was pressed.
func Ask(bounds xgal.Rectangle, prompt string, buttons ...string) *AskLayer {
	return &AskLayer{
		Bounds:  bounds,
		Style:   DefaultStyle(),
		Prompt:  prompt,
		Buttons: buttons,
		Result:  -1,
		hover:   -1,
	}
}

// AskEntry returns a new [AskLayer] with an entry.  The caller must add it to
// a container (e.g. via [Layer.Add]). After Poll returns Finish, read Result
// to see which button was pressed.
func AskEntry(bounds xgal.Rectangle, prompt string, entry string, enter func(string), buttons ...string) *AskLayer {
	ls := DefaultStyle().Stride()
	entryBounds := xgal.Rect(bounds.Min.X, bounds.Min.Y+ls, bounds.Max.X, bounds.Min.Y+ls*2)
	entryWidget := Entry(entryBounds, entry, enter)
	return &AskLayer{
		Bounds:  bounds,
		Style:   DefaultStyle(),
		Prompt:  prompt,
		Buttons: buttons,
		Entry:   entryWidget,
		Result:  -1,
		hover:   -1,
	}
}

var _ Widget = &AskLayer{}

func (a *AskLayer) Poll() Reply {
	a.hover = -1
	pos := xgal.Mouse()
	if !pos.In(a.Bounds) {
		return Ignore
	}

	if a.Entry != nil {
		res := a.Entry.Poll()
		if res != Ignore {
			return res
		}
	}

	pad := a.Style.Margin
	bw := (a.Bounds.Dx() - pad.X*(len(a.Buttons)+1)) / max(len(a.Buttons), 1)
	bh := a.Style.MeasureText("X").Y + pad.Y*2
	by := a.Bounds.Max.Y - pad.Y - bh

	for i := range a.Buttons {
		bx := a.Bounds.Min.X + pad.X + i*(bw+pad.X)
		bb := xgal.Rect(bx, by, bx+bw, by+bh)
		if pos.In(bb) {
			a.hover = i
			if xgal.Click(xgal.MouseButtonLeft) {
				a.Result = i
				return Finish
			}
			return Accept
		}
	}
	return Accept
}

func (a *AskLayer) Render(s *xgal.Surface) {
	a.Style.DrawBox(s, a.Bounds)
	a.Style.Ink(s, a.Bounds, a.Prompt)
	if a.Entry != nil {
		a.Entry.Render(s)
	}
	a.drawButtons(s)
}

func (a *AskLayer) drawButtons(s *xgal.Surface) {
	pad := a.Style.Margin
	bw := (a.Bounds.Dx() - pad.X*(len(a.Buttons)+1)) / max(len(a.Buttons), 1)
	bh := a.Style.MeasureText("X").Y + pad.Y*2
	by := a.Bounds.Max.Y - pad.Y - bh

	for i, label := range a.Buttons {
		bx := a.Bounds.Min.X + pad.X + i*(bw+pad.X)
		bb := xgal.Rect(bx, by, bx+bw, by+bh)

		st := a.Style
		if i == a.hover {
			st = st.HoverStyle()
		}
		st.DrawBox(s, bb)
		st.Ink(s, bb, label)
	}
}

func (a *AskLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	a.Bounds = bounds
	return a.Bounds
}

func (a *AskLayer) MoveBy(delta xgal.Point) {
	a.Bounds = a.Bounds.Add(delta)
}

// AddAsk adds an [AskLayer] to this layer as a modal dialog.
func (m *Layer) AddAsk(prompt string, buttons ...string) *AskLayer {
	pos := xgal.Mouse()
	sz := m.Style.MeasureText(prompt)
	dw := sz.X + m.Style.Margin.X*4
	dh := sz.Y + m.Style.Margin.Y*6 + sz.Y // text + button row
	bounds := xgal.Rect(pos.X, pos.Y, pos.X+dw, pos.Y+dh)
	d := Ask(bounds, prompt, buttons...)
	m.Add(d)
	return d
}

// AddAskEntry adds an [AskLayer] with Entry to this layer as a modal dialog.
func (m *Layer) AddAskEntry(prompt, entry string, enter func(string), buttons ...string) *AskLayer {
	pos := xgal.Mouse()
	sz := m.Style.MeasureText(prompt)
	dw := sz.X + m.Style.Margin.X*4
	dh := sz.Y + m.Style.Margin.Y*6 + sz.Y // text + button row
	bounds := xgal.Rect(pos.X, pos.Y, pos.X+dw, pos.Y+dh)
	d := AskEntry(bounds, prompt, entry, enter, buttons...)
	m.Add(d)
	return d
}
