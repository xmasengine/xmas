package xui

import "github.com/xmasengine/xmas/xgal"

// LabelLayer is a text label.
type LabelLayer struct {
	Bounds xgal.Rectangle
	Style  Style
	Text   string
	Icon   Icon // optional icon, drawn left of text
	hover  bool
}

// Label returns a new [LabelLayer] with the given bounds and text.
func Label(bounds xgal.Rectangle, text string) *LabelLayer {
	return &LabelLayer{
		Bounds: bounds,
		Style:  DefaultStyle(),
		Text:   text,
	}
}

var _ Widget = &LabelLayer{}

func (l *LabelLayer) Poll() Reply {
	if xgal.Cursor().In(l.Bounds) {
		l.hover = true
		return Accept
	}
	l.hover = false
	return Ignore
}

func (l *LabelLayer) Render(s *xgal.Surface) {
	style := l.Style
	if l.hover {
		style = style.HoverStyle()
	}
	style.DrawBox(s, l.Bounds)
	l.Icon.Blit(s, l.Bounds.Min)
	style.Ink(s, l.Icon.TextBounds(l.Bounds), l.Text)
}

func (l *LabelLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	sz := l.Style.MeasureText(l.Text)
	nw := sz.X + l.Style.Margin.X*2 + l.Icon.Width()
	nh := sz.Y + l.Style.Margin.Y*2
	l.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return l.Bounds
}

func (l *LabelLayer) MoveBy(delta xgal.Point) {
	l.Bounds = l.Bounds.Add(delta)
}

func (l *LabelLayer) SetText(text string) {
	l.Text = text
}

func (m *Layer) AddLabel(bounds xgal.Rectangle, text string) *LabelLayer {
	l := Label(bounds, text)
	m.Add(l)
	return l
}
