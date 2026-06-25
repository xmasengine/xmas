package xui

import "github.com/xmasengine/xmas/xgal"

func (s Style) MeasureText(txt string) Point {
	w, h := text.Measure(txt, s.Face, float64(LineHeight(s.Face)))
	return image.Pt(int(w), int(h))
}

func (s Style) DrawText(dst *Surface, at Point, txt string) {
	pt := at.Add(s.Margin)
	DrawText(dst, s.Face, s.Writing, pt.X, pt.Y, txt)
}

func (s Style) DrawTextLine(dst *Surface, at Point, txt string) {
	pt := at.Add(s.Margin)
	DrawTextLine(dst, s.Face, s.Writing, pt.X, pt.Y, txt)
}

// Style is a style for a widget
type Style struct {
	Fore   RGBA
	Border RGBA
	Shadow RGBA
	Filled RGBA
	Stroke int
	Margin xgal.Point
	Face   xgal.Face
}

func DefaultStyle() Style {
	s := Style{}
	s.Border = RGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xff}
	s.Shadow = RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa}
	s.Filled = RGBA{R: 0x00, G: 0x00, B: 0x55, A: 0xaa}
	s.Stroke = 1
	s.Margin = image.Pt(2, 2)
	s.Face = xgal.DefaultFace
	return s
}

func (s Style) DrawRect(Surface *Surface, r Rectangle) {
	xgal.Outline(Surface, r, int(s.Stroke), s.Border)
}

func (s Style) DrawBox(Surface *Surface, r Rectangle) {
	if s.Shadow.A != 0 {
		shadow := s.Shadow
		shadow.A = (shadow.A / 2) + 1 // make half transparent
		right := image.Rect(r.Max.X+1, r.Min.Y+1, r.Max.X+1, r.Max.Y+1)
		xgal.Line(Surface, right, 1, shadow)
		bottom := image.Rect(r.Min.X+1, r.Max.Y+1, r.Max.X+1, r.Max.Y+1)
		xgal.Line(Surface, bottom, 1, shadow)
	}

	xgal.Box(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()), s.Filled, false,
	)

	if s.Stroke > 0 {
		xgal.Outline(
			Surface, float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()),
			float32(s.Stroke), s.Border, false,
		)
	}
}

func (s Style) DrawCircleInBox(Surface *Surface, box Rectangle) {
	r := box.Dx()
	if box.Dy() < r {
		r = box.Dy()
	}
	r = r / 2
	c := image.Pt((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2)
	s.DrawCircle(Surface, c, r)
}

func (s Style) DrawCircle(Surface *Surface, c Point, r int) {
	if r < 0 {
		r = 1
	}
	vector.DrawFilledCircle(Surface, float32(c.X), float32(c.Y),
		float32(r), s.Filled, false)

	if s.Stroke > 0 {
		vector.StrokeCircle(
			Surface, float32(c.X), float32(c.Y),
			float32(r), float32(s.Stroke), s.Border, false,
		)
	}
}

func FocusStyle() Style {
	s := DefaultStyle()
	s.Border = color.RGBA{240, 140, 40, 245}
	s.Writing = color.RGBA{245, 245, 245, 245}
	s.Fill = color.RGBA{128, 128, 200, 240}
	return s
}

func HoverStyle() Style {
	s := DefaultStyle()
	s.Border = color.RGBA{240, 240, 50, 250}
	return s
}

func PressStyle() Style {
	s := DefaultStyle()
	s.Fill = color.RGBA{15, 45, 200, 240}
	return s
}

func BarStyle() Style {
	s := DefaultStyle().WithTinyFont()
	s.Fill = color.RGBA{45, 45, 200, 250}
	return s
}

func CheckStyle() Style {
	s := DefaultStyle()
	s.Fill = color.RGBA{245, 245, 245, 250}
	return s
}

func (s Style) HoverStyle() Style {
	s.Border = color.RGBA{200, 200, 45, 250}
	return s
}

func (s Style) FocusStyle() Style {
	s.Border = color.RGBA{240, 140, 40, 245}
	s.Writing = color.RGBA{245, 245, 245, 245}
	s.Fill = color.RGBA{128, 128, 200, 245}
	return s
}

func (s Style) PressStyle() Style {
	s.Fill = color.RGBA{15, 45, 200, 240}
	return s
}

func (s Style) DragStyle() Style {
	s.Fill = color.RGBA{15, 128, 200, 240}
	return s
}

func (s Style) ForState(state State) Style {
	if state.Focus {
		return s.FocusStyle()
	}
	if state.Hover {
		return s.HoverStyle()
	}
	if state.Drag {
		return s.DragStyle()
	}
	return s
}

func (s Style) BarStyle() Style {
	s = s.WithTinyFont()
	s.Fill = color.RGBA{45, 45, 245, 250}
	return s
}

func (s Style) CheckStyle() Style {
	s.Fill = color.RGBA{245, 245, 245, 250}
	return s
}

func (s Style) KnobStyle() Style {
	s.Fill = color.RGBA{245, 245, 245, 250}
	return s
}

func (s Style) WithTinyFont() Style {
	s.Face = spleen8.XFace
	return s
}
