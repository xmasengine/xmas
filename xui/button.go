package xui

import "github.com/xmasengine/xmas/xgal"

// ButtonLayer is a clickable button.
type ButtonLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Text    string
	Icon    Icon // optional icon, drawn left of text
	Clicked func()
	pressed bool
	hover   bool
}

// Button returns a new [ButtonLayer] with the given bounds, text, and click handler.
func Button(bounds xgal.Rectangle, text string, clicked func()) *ButtonLayer {
	return &ButtonLayer{
		Bounds:  bounds,
		Style:   DefaultStyle(),
		Text:    text,
		Clicked: clicked,
	}
}

var _ Widget = &ButtonLayer{}

func (b *ButtonLayer) Poll() Reply {
	b.hover = xgal.Mouse().In(b.Bounds)
	if !b.hover {
		if xgal.Loose(xgal.MouseButtonLeft) {
			b.pressed = false
		}
		return Ignore
	}

	if xgal.Click(xgal.MouseButtonLeft) {
		b.pressed = true
		return Accept
	}

	if xgal.Loose(xgal.MouseButtonLeft) {
		b.pressed = false
		if b.Clicked != nil {
			b.Clicked()
		}
		return Accept
	}

	return Ignore
}

func (b *ButtonLayer) Render(s *xgal.Surface) {
	box := b.Bounds
	style := b.Style

	if b.hover {
		style = style.HoverStyle()
	}
	if b.pressed {
		box = box.Add(b.Style.Margin)
		style = style.PressStyle()
	}

	style.DrawBox(s, box)
	b.Icon.Blit(s, box.Min)
	style.Ink(s, b.Icon.TextBounds(box), b.Text)
}

func (b *ButtonLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	sz := b.Style.MeasureText(b.Text)
	nw := sz.X + b.Style.Margin.X*2 + b.Icon.Width()
	nh := sz.Y + b.Style.Margin.Y*2
	b.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return b.Bounds
}

func (b *ButtonLayer) MoveBy(delta xgal.Point) {
	b.Bounds = b.Bounds.Add(delta)
}

func (m *Layer) AddButton(bounds xgal.Rectangle, text string, clicked func()) *ButtonLayer {
	b := Button(bounds, text, clicked)
	m.Add(b)
	return b
}
