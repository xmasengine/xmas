package xlui

import "github.com/xmasengine/xmas/xgal"
import "github.com/xmasengine/xmas/xvec"
import "github.com/xmasengine/xmas/xres/fontres"

import "os"
import "log/slog"

var (
	TinyFace   = fontres.TinyFace
	SmallFace  = fontres.SmallFace
	NormalFace = xgal.BuiltinFace
)

var DefaultFace = SmallFace

func (s Style) Measure(txt string) xgal.Point {
	w, h := xgal.Measure(txt, s.Face, float64(xgal.Stride(s.Face)))
	return xgal.Pt(int(w), int(h))
}

func (s Style) Stride() int {
	return xgal.Stride(s.Face)
}

func (s Style) Print(dst *xgal.Surface, at xgal.Point, txt string) {
	pt := at.Add(s.Margin).Add(s.Offset)
	xgal.Print(dst, s.Face, s.Fore, pt.X, pt.Y, txt)
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
	Shade  xgal.Point // Shade is the direction of the shadow
	Gloom  int        // Gloom is the width of the shadow
	Offset xgal.Point // Offset for buttons, etc.

	Vec   *xvec.XVEC
	Frame *xgal.Surface
}

var DefaultFS = os.DirFS("pack/image/ui")

const DefaultFrameName = "frame.xvec"

var DefaultVec *xvec.XVEC
var DefaultFrame *xgal.Surface

func init() {
	var err error
	DefaultVec, err = xvec.ParseFS(DefaultFS, DefaultFrameName)
	if err != nil {
		slog.Error("Could not load frame.", "err", err)
	} else {
		DefaultFrame = xgal.Prepare(int(DefaultVec.Size.W), int(DefaultVec.Size.H))
		DefaultVec.Draw(DefaultFrame)
	}
}

func DefaultStyle() Style {
	s := Style{}
	s.Fore = xgal.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	s.Border = xgal.RGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xff}
	s.Shadow = xgal.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	s.Fill = xgal.RGBA{R: 0x00, G: 0x00, B: 0x55, A: 0xaa}
	s.Stroke = 1
	s.Gloom = 10
	s.Margin = xgal.Pt(2, 2)
	s.Shade = xgal.Pt(1, 1)
	s.Face = DefaultFace
	s.Frame = nil
	return s
}

func (s Style) DrawRect(dst *xgal.Surface, r xgal.Rectangle) {
	xgal.Outline(dst, r, int(s.Stroke), s.Border)
}

func (s Style) DrawPlainBox(dst *xgal.Surface, r xgal.Rectangle) {
	xgal.Box(dst, r.Add(s.Offset), s.Fill)
}

const OffNine = 5

func (s Style) DrawBox(dst *xgal.Surface, r xgal.Rectangle) {
	xgal.Box(dst, r.Add(s.Offset), s.Fill)

	if s.Frame != nil {
		offset := r.Min.Add(s.Offset)
		offsetX, offsetY := float32(offset.X), float32(offset.Y)
		dstW, dstH := (r.Dx()), (r.Dy())
		xgal.NineSlice(dst, s.Frame, offsetX, offsetY, dstW, dstH, OffNine, OffNine, OffNine, OffNine)
	}

	if s.Gloom > 0 {
		lo := r.Min.Add(s.Offset).Add(s.Shade)
		hi := r.Max.Add(s.Offset).Add(s.Shade)
		xgal.Line(dst, hi.X, lo.Y, hi.X, lo.Y, s.Gloom, s.Shadow)
		xgal.Line(dst, lo.X, hi.Y, lo.X, hi.Y, s.Gloom, s.Shadow)
	}

	if s.Stroke > 0 {
		xgal.Outline(dst, r.Add(s.Offset), s.Stroke, s.Border)
	}
}

func (s Style) DrawCircleInBox(Surface *xgal.Surface, box xgal.Rectangle) {
	box = box.Add(s.Offset)
	r := box.Dx()
	if box.Dy() < r {
		r = box.Dy()
	}
	r = r / 2
	c := xgal.Pt((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2)
	s.DrawCircle(Surface, c, r)
}

func (s Style) DrawCircle(dst *xgal.Surface, c xgal.Point, r int) {
	c = c.Add(s.Offset)
	if r < 0 {
		r = 1
	}
	xgal.Disk(dst, c, r, s.Fill)

	if s.Stroke > 0 {
		xgal.Circle(dst, c, r, s.Stroke, s.Border)
	}
}

func (s Style) DrawX(dst *xgal.Surface, bounds xgal.Rectangle) {
	bounds = bounds.Add(s.Offset)
	if s.Stroke > 0 {
		xgal.Andreas(dst, bounds, s.Stroke, s.Border)
	}
}

func (s Style) Ink(dst *xgal.Surface, bounds xgal.Rectangle, text string) {
	bounds = bounds.Add(s.Offset)
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

func (s Style) Focused() Style {
	s.Border = xgal.Paint(240, 240, 240, 245)
	return s
}

func (s Style) Clicked() Style {
	s.Border = xgal.Paint(0xff, 0xff, 0xff, 0xff)
	s.Fill = xgal.Paint(0x55, 0x55, 0xff, 0xff)
	s.Gloom = 2
	s.Shade = xgal.Pt(-1, -1)
	s.Offset = xgal.Pt(0, 1)
	return s
}

func (s Style) Hovered() Style {
	s.Border = xgal.Paint(240, 240, 50, 250)
	return s
}

func FocusStyle() Style {
	s := DefaultStyle()
	s.Border = xgal.Paint(240, 140, 40, 245)
	s.Fill = xgal.Paint(128, 128, 200, 240)
	return s
}

func HoverStyle() Style {
	s := DefaultStyle()
	s.Border = xgal.Paint(240, 240, 50, 250)
	return s
}

func PressStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.Paint(15, 45, 200, 240)
	return s
}

func ButtonStyle() Style {
	s := DefaultStyle()
	s.Margin = xgal.Pt(2, 0)
	s.Fill = xgal.Paint(0x00, 0x00, 0xaa, 0xaa)
	s.Shadow = xgal.Paint(0x00, 0x00, 0x00, 0xff)
	return s
}

func BarStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.Paint(45, 45, 200, 250)
	return s
}

func CheckStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.Paint(245, 245, 245, 250)
	return s
}

func (s Style) HoverStyle() Style {
	s.Border = xgal.Paint(200, 200, 45, 250)
	return s
}

func (s Style) FocusStyle() Style {
	s.Border = xgal.Paint(240, 140, 40, 245)
	s.Fill = xgal.Paint(128, 128, 200, 245)
	return s
}

func (s Style) PressStyle() Style {
	s.Fill = xgal.Paint(15, 45, 200, 240)
	return s
}

func (s Style) DragStyle() Style {
	s.Fill = xgal.Paint(15, 128, 200, 240)
	return s
}

func (s Style) BarStyle() Style {
	s.Fill = xgal.Paint(45, 45, 245, 250)
	return s
}

func (s Style) CheckStyle() Style {
	s.Fill = xgal.Paint(245, 245, 245, 250)
	return s
}

func (s Style) ActiveStyle() Style {
	s.Fill = xgal.Paint(60, 60, 140, 255)
	s.Border = xgal.Paint(120, 120, 220, 255)
	return s
}

func (s Style) KnobStyle() Style {
	s.Fill = xgal.Paint(245, 245, 245, 250)
	return s
}

func MenuStyle() Style {
	s := DefaultStyle()
	s.Margin = xgal.Pt(2, 0)
	s.Gloom = 0
	s.Stroke = 0
	return s
}

func HUDBarStyle() Style {
	s := DefaultStyle()
	s.Fill = xgal.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	return s
}
