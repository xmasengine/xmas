package xui

import "github.com/xmasengine/xmas/xgal"

func (s Style) MeasureText(txt string) xgal.Point {
	w, h := xgal.Measure(txt, s.Face, float64(xgal.Stride(s.Face)))
	return xgal.Pt(int(w), int(h))
}

func (s Style) DrawText(dst *xgal.Surface, at xgal.Point, txt string) {
	pt := at.Add(s.Margin)
	xgal.Ink(dst, s.Face, s.Fore, pt.X, pt.Y, txt)
}

func (s Style) DrawTextLine(dst *xgal.Surface, at xgal.Point, txt string) {
	pt := at.Add(s.Margin)
	xgal.Ink(dst, s.Face, s.Fore, pt.X, pt.Y, txt)
}

// Style is a style for a widget
type Style struct {
	Fore   xgal.RGBA
	Border xgal.RGBA
	Shadow xgal.RGBA
	Fill   xgal.RGBA
	Stroke int
	Margin xgal.Point
	Face   xgal.Face
}

func DefaultStyle() Style {
	s := Style{}
	s.Fore = xgal.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	s.Border = xgal.RGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xff}
	s.Shadow = xgal.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa}
	s.Fill = xgal.RGBA{R: 0x00, G: 0x00, B: 0x55, A: 0xaa}
	s.Stroke = 1
	s.Margin = xgal.Pt(2, 2)
	s.Face = xgal.BuiltinFace
	return s
}

func (s Style) DrawRect(dst *xgal.Surface, r xgal.Rectangle) {
	xgal.Outline(dst, r, int(s.Stroke), s.Border)
}

func (s Style) DrawBox(dst *xgal.Surface, r xgal.Rectangle) {
	if s.Shadow.A != 0 {
		shadow := s.Shadow
		shadow.A = (shadow.A / 2) + 1 // make half transparent
		xgal.Line(dst, r.Max.X+1, r.Min.Y+1, r.Max.X+1, r.Max.Y+1, 1, shadow)
		xgal.Line(dst, r.Min.X+1, r.Max.Y+1, r.Max.X+1, r.Max.Y+1, 1, shadow)
	}

	xgal.Box(dst, r, s.Fill)

	if s.Stroke > 0 {
		xgal.Outline(dst, r, s.Stroke, s.Border)
	}
}

func (s Style) DrawCircleInBox(Surface *xgal.Surface, box xgal.Rectangle) {
	r := box.Dx()
	if box.Dy() < r {
		r = box.Dy()
	}
	r = r / 2
	c := xgal.Pt((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2)
	s.DrawCircle(Surface, c, r)
}

func (s Style) DrawCircle(dst *xgal.Surface, c xgal.Point, r int) {
	if r < 0 {
		r = 1
	}
	xgal.Disk(dst, c, r, s.Fill)

	if s.Stroke > 0 {
		xgal.Circle(dst, c, r, s.Stroke, s.Border)
	}
}

func (s Style) DrawX(dst *xgal.Surface, bounds xgal.Rectangle) {
	if s.Stroke > 0 {
		xgal.Andreas(dst, bounds, s.Stroke, s.Border)
	}
}

func (s Style) Ink(dst *xgal.Surface, bounds xgal.Rectangle, text string) {
	xgal.Ink(dst, s.Face, s.Fore, bounds.Min.X+s.Margin.X, bounds.Min.Y+s.Margin.Y, text)
}

// Inset shrinks the given rectangle by the style's margins and returns it.
func (s Style) Inset(bounds xgal.Rectangle) xgal.Rectangle {
	margin := s.Margin
	xmin := bounds.Min.X + margin.X
	ymin := bounds.Min.Y + margin.Y
	xmax := bounds.Max.X - margin.X
	ymax := bounds.Max.Y - margin.Y
	return xgal.Rect(xmin, ymin, xmax, ymax)
}

func FocusStyle() Style {
	s := DefaultStyle()
	s.Border = xgal.Wash(240, 140, 40, 245)
	s.Fill = xgal.Wash(128, 128, 200, 240)
	return s
}

func HoverStyle() Style {
	s := DefaultStyle()
	s.Border = xgal.Wash(240, 240, 50, 250)
	return s
}

func PressStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.Wash(15, 45, 200, 240)
	return s
}

func BarStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.Wash(45, 45, 200, 250)
	return s
}

func CheckStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.Wash(245, 245, 245, 250)
	return s
}

func (s Style) HoverStyle() Style {
	s.Border = xgal.Wash(200, 200, 45, 250)
	return s
}

func (s Style) FocusStyle() Style {
	s.Border = xgal.Wash(240, 140, 40, 245)
	s.Fill = xgal.Wash(128, 128, 200, 245)
	return s
}

func (s Style) PressStyle() Style {
	s.Fill = xgal.Wash(15, 45, 200, 240)
	return s
}

func (s Style) DragStyle() Style {
	s.Fill = xgal.Wash(15, 128, 200, 240)
	return s
}

func (s Style) BarStyle() Style {
	s.Fill = xgal.Wash(45, 45, 245, 250)
	return s
}

func (s Style) CheckStyle() Style {
	s.Fill = xgal.Wash(245, 245, 245, 250)
	return s
}

func (s Style) ActiveStyle() Style {
	s.Fill = xgal.Wash(60, 60, 140, 255)
	s.Border = xgal.Wash(120, 120, 220, 255)
	return s
}

func (s Style) KnobStyle() Style {
	s.Fill = xgal.Wash(245, 245, 245, 250)
	return s
}
