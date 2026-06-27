// Package xvec implements a vector graphics file format and renderer
// backed by ebitengine's ebiten/v2/vector package.
//
// The xvec text format stores a fixed-size canvas, an anti-aliasing flag,
// and a list of drawing instructions. See xvec/SPEC.md for the full format
// specification.
//
// Supported primitives:
//   - circle outlines (stroke) and filled circles (disk)
//   - rectangle outlines and filled rectangles
//   - line segments
//   - filled and stroked paths with lines, arcs, cubic and quad beziers
package xvec

import (
	"encoding"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Surface is a drawable image.
type Surface = ebiten.Image

// Color is an 8-bit RGBA color.
type Color = color.RGBA

// Path is a vector path.
type Path = vector.Path

// FillOptions is options for filling a path.
type FillOptions = vector.FillOptions

// StrokeOptions is options for stroking a path.
type StrokeOptions = vector.StrokeOptions

// DrawPathOptions is options for drawing a path.
type DrawPathOptions = vector.DrawPathOptions

// Vertex is a point in 2D space (float32).
type Vertex struct {
	X float32
	Y float32
}

// Length is a distance (float32).
type Length float32

// Size defines the drawing dimensions.
type Size struct {
	W float32
	H float32
}

// Compile-time interface checks.
var (
	_ encoding.TextMarshaler = (*CircleInstruction)(nil)
	_ encoding.TextMarshaler = (*DiskInstruction)(nil)
	_ encoding.TextMarshaler = (*RectInstruction)(nil)
	_ encoding.TextMarshaler = (*SlabInstruction)(nil)
	_ encoding.TextMarshaler = (*LineInstruction)(nil)
	_ encoding.TextMarshaler = (*FillInstruction)(nil)
	_ encoding.TextMarshaler = (*StrokeInstruction)(nil)
	_ encoding.TextMarshaler = (*MoveStep)(nil)
	_ encoding.TextMarshaler = (*LineStep)(nil)
	_ encoding.TextMarshaler = (*QuadStep)(nil)
	_ encoding.TextMarshaler = (*CubicStep)(nil)
	_ encoding.TextMarshaler = (*ArcStep)(nil)
	_ encoding.TextMarshaler = (*ArcToStep)(nil)
	_ encoding.TextMarshaler = (*CloseStep)(nil)
)

// XVEC is a complete vector graphic: a canvas size, antialias flag,
// and a list of drawing instructions.
type XVEC struct {
	Size         Size
	Antialias    bool
	Instructions []Instruction
}

// Instruction is a single drawing operation.
type Instruction interface {
	Draw(*Surface)
	encoding.TextMarshaler
}

// Stepper is a single step in a vector path (MoveTo, LineTo, etc.).
type Stepper interface {
	Step(p *Path)
	encoding.TextMarshaler
}

// CircleInstruction strokes a circle outline.
type CircleInstruction struct {
	C         Vertex
	R         Length
	Color     Color
	Stroke    Length
	Antialias bool
}

func (x *XVEC) Circle(cx, cy, r, stroke float32, col Color) *CircleInstruction {
	c := &CircleInstruction{
		C: Vertex{cx, cy}, R: Length(r), Stroke: Length(stroke),
		Color: col, Antialias: x.Antialias,
	}
	x.Instructions = append(x.Instructions, c)
	return c
}

func (c *CircleInstruction) Draw(s *Surface) {
	vector.StrokeCircle(s, c.C.X, c.C.Y, float32(c.R), float32(c.Stroke), c.Color, c.Antialias)
}

func (c *CircleInstruction) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("circle %s %s %s %s %s",
		ftos(c.C.X), ftos(c.C.Y), ftos(float32(c.R)), ftos(float32(c.Stroke)), coltos(c.Color))), nil
}

// DiskInstruction fills a circle.
type DiskInstruction struct {
	C         Vertex
	R         Length
	Color     Color
	Antialias bool
}

func (x *XVEC) Disk(cx, cy, r float32, col Color) *DiskInstruction {
	d := &DiskInstruction{
		C: Vertex{cx, cy}, R: Length(r),
		Color: col, Antialias: x.Antialias,
	}
	x.Instructions = append(x.Instructions, d)
	return d
}

func (d *DiskInstruction) Draw(s *Surface) {
	vector.FillCircle(s, d.C.X, d.C.Y, float32(d.R), d.Color, d.Antialias)
}

func (d *DiskInstruction) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("disk %s %s %s %s",
		ftos(d.C.X), ftos(d.C.Y), ftos(float32(d.R)), coltos(d.Color))), nil
}

// RectInstruction strokes a rectangle outline.
type RectInstruction struct {
	X, Y, W, H float32
	Color      Color
	Stroke     Length
	Antialias  bool
}

func (x *XVEC) Rect(rx, ry, w, h, stroke float32, col Color) *RectInstruction {
	r := &RectInstruction{
		X: rx, Y: ry, W: w, H: h,
		Color: col, Stroke: Length(stroke), Antialias: x.Antialias,
	}
	x.Instructions = append(x.Instructions, r)
	return r
}

func (r *RectInstruction) Draw(s *Surface) {
	vector.StrokeRect(s, r.X, r.Y, r.W, r.H, float32(r.Stroke), r.Color, r.Antialias)
}

func (r *RectInstruction) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("rect %s %s %s %s %s %s",
		ftos(r.X), ftos(r.Y), ftos(r.W), ftos(r.H), ftos(float32(r.Stroke)), coltos(r.Color))), nil
}

// SlabInstruction fills a rectangle.
type SlabInstruction struct {
	X, Y, W, H float32
	Color      Color
	Antialias  bool
}

func (x *XVEC) Slab(rx, ry, w, h float32, col Color) *SlabInstruction {
	r := &SlabInstruction{
		X: rx, Y: ry, W: w, H: h,
		Color: col, Antialias: x.Antialias,
	}
	x.Instructions = append(x.Instructions, r)
	return r
}

func (r *SlabInstruction) Draw(s *Surface) {
	vector.FillRect(s, r.X, r.Y, r.W, r.H, r.Color, r.Antialias)
}

func (r *SlabInstruction) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("slab %s %s %s %s %s",
		ftos(r.X), ftos(r.Y), ftos(r.W), ftos(r.H), coltos(r.Color))), nil
}

// LineInstruction strokes a line segment.
type LineInstruction struct {
	X1, Y1, X2, Y2 float32
	Color          Color
	Stroke         Length
	Antialias      bool
}

func (x *XVEC) Line(x1, y1, x2, y2, stroke float32, col Color) *LineInstruction {
	l := &LineInstruction{
		X1: x1, Y1: y1, X2: x2, Y2: y2,
		Color: col, Stroke: Length(stroke), Antialias: x.Antialias,
	}
	x.Instructions = append(x.Instructions, l)
	return l
}

func (l *LineInstruction) Draw(s *Surface) {
	vector.StrokeLine(s, l.X1, l.Y1, l.X2, l.Y2, float32(l.Stroke), l.Color, l.Antialias)
}

func (l *LineInstruction) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("line %s %s %s %s %s %s",
		ftos(l.X1), ftos(l.Y1), ftos(l.X2), ftos(l.Y2), ftos(float32(l.Stroke)), coltos(l.Color))), nil
}

// FillInstruction fills a path built from steps.
type FillInstruction struct {
	Color     Color
	FillOpts  FillOptions
	DrawOpts  DrawPathOptions
	Steps     []Stepper
	Antialias bool
}

func (x *XVEC) Fill(col Color, steps ...Stepper) *FillInstruction {
	f := &FillInstruction{Color: col, Steps: steps, Antialias: x.Antialias}
	f.DrawOpts.AntiAlias = x.Antialias
	f.DrawOpts.ColorScale.Reset()
	x.Instructions = append(x.Instructions, f)
	return f
}

func (f *FillInstruction) Draw(s *Surface) {
	var path Path
	for _, step := range f.Steps {
		step.Step(&path)
	}
	fill := f.FillOpts
	opts := f.DrawOpts
	opts.ColorScale.ScaleWithColor(f.Color)

	vector.FillPath(s, &path, &fill, &opts)
}

func (f *FillInstruction) MarshalText() ([]byte, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "fill %s\n", coltos(f.Color))
	for _, step := range f.Steps {
		txt, err := step.MarshalText()
		if err != nil {
			return nil, err
		}
		b.WriteString("  ")
		b.Write(txt)
		b.WriteByte('\n')
	}
	b.WriteString("end")
	return []byte(b.String()), nil
}

// StrokeInstruction strokes a path built from steps.
type StrokeInstruction struct {
	Color      Color
	Width      Length
	StrokeOpts StrokeOptions
	DrawOpts   DrawPathOptions
	Steps      []Stepper
	Antialias  bool
}

func (x *XVEC) Stroke(width float32, col Color, steps ...Stepper) *StrokeInstruction {
	s := &StrokeInstruction{
		Color: col, Width: Length(width), Steps: steps, Antialias: x.Antialias,
	}
	s.DrawOpts.AntiAlias = x.Antialias
	s.DrawOpts.ColorScale.Reset()
	s.StrokeOpts.Width = float32(s.Width)
	x.Instructions = append(x.Instructions, s)
	return s
}

func (s *StrokeInstruction) Draw(dst *Surface) {
	var path Path
	for _, step := range s.Steps {
		step.Step(&path)
	}
	so := s.StrokeOpts
	so.Width = float32(s.Width)

	opts := s.DrawOpts
	opts.ColorScale.ScaleWithColor(s.Color)

	vector.StrokePath(dst, &path, &so, &opts)
}

func (s *StrokeInstruction) MarshalText() ([]byte, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "stroke %s %s\n", ftos(float32(s.Width)), coltos(s.Color))
	for _, step := range s.Steps {
		txt, err := step.MarshalText()
		if err != nil {
			return nil, err
		}
		b.WriteString("  ")
		b.Write(txt)
		b.WriteByte('\n')
	}
	b.WriteString("end")
	return []byte(b.String()), nil
}

// MoveStep moves the current point without drawing, starting a new sub-path.
type MoveStep struct {
	X, Y float32
}

// MoveTo returns a MoveStep that starts a new sub-path at (x, y).
func MoveTo(x, y float32) *MoveStep { return &MoveStep{X: x, Y: y} }

func (m *MoveStep) Step(p *Path) { p.MoveTo(m.X, m.Y) }

func (m *MoveStep) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("move %s %s", ftos(m.X), ftos(m.Y))), nil
}

// LineStep draws a straight line to (x, y).
type LineStep struct {
	X, Y float32
}

// LineTo returns a LineStep that draws a line to (x, y).
func LineTo(x, y float32) *LineStep { return &LineStep{X: x, Y: y} }

func (l *LineStep) Step(p *Path) { p.LineTo(l.X, l.Y) }

func (l *LineStep) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("line %s %s", ftos(l.X), ftos(l.Y))), nil
}

// QuadStep draws a quadratic Bézier curve to (x2, y2) with control point (x1, y1).
type QuadStep struct {
	X1, Y1, X2, Y2 float32
}

// QuadTo returns a QuadStep for a quadratic Bézier curve.
func QuadTo(x1, y1, x2, y2 float32) *QuadStep { return &QuadStep{X1: x1, Y1: y1, X2: x2, Y2: y2} }

func (q *QuadStep) Step(p *Path) { p.QuadTo(q.X1, q.Y1, q.X2, q.Y2) }

func (q *QuadStep) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("quad %s %s %s %s", ftos(q.X1), ftos(q.Y1), ftos(q.X2), ftos(q.Y2))), nil
}

// CubicStep draws a cubic Bézier curve to (x3, y3) with control points (x1, y1) and (x2, y2).
type CubicStep struct {
	X1, Y1, X2, Y2, X3, Y3 float32
}

// CubicTo returns a CubicStep for a cubic Bézier curve.
func CubicTo(x1, y1, x2, y2, x3, y3 float32) *CubicStep {
	return &CubicStep{X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3}
}

func (c *CubicStep) Step(p *Path) { p.CubicTo(c.X1, c.Y1, c.X2, c.Y2, c.X3, c.Y3) }

func (c *CubicStep) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("cubic %s %s %s %s %s %s",
		ftos(c.X1), ftos(c.Y1), ftos(c.X2), ftos(c.Y2), ftos(c.X3), ftos(c.Y3))), nil
}

// ArcStep draws a circular arc. Start and End are in radians.
type ArcStep struct {
	CX, CY, R float32
	Start     float32
	End       float32
	Dir       Direction
}

// Direction is the sweep direction of an arc.
type Direction = vector.Direction

const (
	Clockwise        Direction = vector.Clockwise
	CounterClockwise Direction = vector.CounterClockwise
)

// Arc returns an ArcStep drawing an arc centred at (cx, cy) with radius r,
// from start to end (both in radians), in the given direction.
func Arc(cx, cy, r, start, end float32, dir Direction) *ArcStep {
	return &ArcStep{CX: cx, CY: cy, R: r, Start: start, End: end, Dir: dir}
}

func (a *ArcStep) Step(p *Path) { p.Arc(a.CX, a.CY, a.R, a.Start, a.End, a.Dir) }

func (a *ArcStep) MarshalText() ([]byte, error) {
	ds := "C"
	if a.Dir == CounterClockwise {
		ds = "CC"
	}
	return []byte(fmt.Sprintf("arc %s %s %s %s %s %s",
		ftos(a.CX), ftos(a.CY), ftos(a.R), ftos(a.Start), ftos(a.End), ds)), nil
}

// ArcToStep draws a circular arc to (x2, y2) with turning point (x1, y1) and radius r.
type ArcToStep struct {
	X1, Y1, X2, Y2, R float32
}

// ArcTo returns an ArcToStep drawing an arc from the current point to (x2, y2).
func ArcTo(x1, y1, x2, y2, r float32) *ArcToStep {
	return &ArcToStep{X1: x1, Y1: y1, X2: x2, Y2: y2, R: r}
}

func (a *ArcToStep) Step(p *Path) { p.ArcTo(a.X1, a.Y1, a.X2, a.Y2, a.R) }

func (a *ArcToStep) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("arcto %s %s %s %s %s",
		ftos(a.X1), ftos(a.Y1), ftos(a.X2), ftos(a.Y2), ftos(a.R))), nil
}

// CloseStep closes the current sub-path by drawing a line back to its start point.
type CloseStep struct{}

// Close returns a CloseStep that closes the current sub-path.
func Close() *CloseStep { return &CloseStep{} }

func (c *CloseStep) Step(p *Path) { p.Close() }

func (c *CloseStep) MarshalText() ([]byte, error) {
	return []byte("close"), nil
}

// Encode writes the xvec text format to w.
func (x *XVEC) Encode(w io.Writer) error {
	if _, err := fmt.Fprint(w, "xvec 1\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "size %s %s\n", ftos(x.Size.W), ftos(x.Size.H)); err != nil {
		return err
	}
	aa := "false"
	if x.Antialias {
		aa = "true"
	}
	if _, err := fmt.Fprintf(w, "antialias %s\n", aa); err != nil {
		return err
	}
	for _, inst := range x.Instructions {
		txt, err := inst.MarshalText()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%s\n", txt); err != nil {
			return err
		}
	}
	return nil
}

// Decode parses the xvec text format from r using text/scanner.
func (x *XVEC) Decode(r io.Reader) error {
	x.Size = Size{320, 240}
	x.Antialias = true
	x.Instructions = nil

	p := &scannerParser{}
	p.s.Init(r)
	p.s.Mode = scanner.ScanInts | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanComments | scanner.SkipComments

	var curFill *FillInstruction
	var curStroke *StrokeInstruction

	tok := p.s.Scan()
	for tok != scanner.EOF {
		if p.err != nil {
			return p.err
		}
		kw := p.s.TokenText()
		inPath := curFill != nil || curStroke != nil
		switch kw {
		case "xvec":
			v := p.raw()
			if v != "1" {
				return fmt.Errorf("unsupported xvec version %q", v)
			}

		case "size":
			x.Size.W = p.float()
			x.Size.H = p.float()

		case "antialias":
			x.Antialias = p.ident() == "true"

		case "circle":
			c := &CircleInstruction{C: V(p.float(), p.float()), R: Length(p.float()), Stroke: Length(p.float()), Color: p.color(), Antialias: x.Antialias}
			x.Instructions = append(x.Instructions, c)

		case "disk":
			d := &DiskInstruction{C: V(p.float(), p.float()), R: Length(p.float()), Color: p.color(), Antialias: x.Antialias}
			x.Instructions = append(x.Instructions, d)

		case "rect":
			r := &RectInstruction{X: p.float(), Y: p.float(), W: p.float(), H: p.float(), Stroke: Length(p.float()), Color: p.color(), Antialias: x.Antialias}
			x.Instructions = append(x.Instructions, r)

		case "slab":
			fr := &SlabInstruction{X: p.float(), Y: p.float(), W: p.float(), H: p.float(), Color: p.color(), Antialias: x.Antialias}
			x.Instructions = append(x.Instructions, fr)

		case "line":
			if inPath {
				addStep(curFill, curStroke, &LineStep{X: p.float(), Y: p.float()})
			} else {
				l := &LineInstruction{X1: p.float(), Y1: p.float(), X2: p.float(), Y2: p.float(), Stroke: Length(p.float()), Color: p.color(), Antialias: x.Antialias}
				x.Instructions = append(x.Instructions, l)
			}

		case "fill":
			curFill = &FillInstruction{Color: p.color(), Steps: nil, Antialias: x.Antialias}
			if x.Antialias {
				curFill.DrawOpts.AntiAlias = true
			}

		case "stroke":
			w := p.float()
			curStroke = &StrokeInstruction{Color: p.color(), Width: Length(w), Steps: nil, Antialias: x.Antialias}
			if x.Antialias {
				curStroke.DrawOpts.AntiAlias = true
			}

		case "end":
			if curFill != nil {
				if len(curFill.Steps) > 0 {
					if _, ok := curFill.Steps[len(curFill.Steps)-1].(*CloseStep); !ok {
						return fmt.Errorf("fill path must end with close")
					}
				}
				x.Instructions = append(x.Instructions, curFill)
				curFill = nil
			} else if curStroke != nil {
				if len(curStroke.Steps) > 0 {
					if _, ok := curStroke.Steps[len(curStroke.Steps)-1].(*CloseStep); !ok {
						return fmt.Errorf("stroke path must end with close")
					}
				}
				x.Instructions = append(x.Instructions, curStroke)
				curStroke = nil
			}

		case "move":
			addStep(curFill, curStroke, &MoveStep{X: p.float(), Y: p.float()})

		case "quad":
			addStep(curFill, curStroke, &QuadStep{X1: p.float(), Y1: p.float(), X2: p.float(), Y2: p.float()})

		case "cubic":
			addStep(curFill, curStroke, &CubicStep{X1: p.float(), Y1: p.float(), X2: p.float(), Y2: p.float(), X3: p.float(), Y3: p.float()})

		case "arc":
			cx := p.float()
			cy := p.float()
			r := p.float()
			start := p.float()
			end := p.float()
			ds := CounterClockwise
			if p.ident() == "C" {
				ds = Clockwise
			}
			addStep(curFill, curStroke, &ArcStep{CX: cx, CY: cy, R: r, Start: start, End: end, Dir: ds})

		case "arcto":
			addStep(curFill, curStroke, &ArcToStep{X1: p.float(), Y1: p.float(), X2: p.float(), Y2: p.float(), R: p.float()})

		case "close":
			addStep(curFill, curStroke, &CloseStep{})
		}
		if p.err != nil {
			return p.err
		}
		tok = p.s.Scan()
	}

	return nil
}

// Draw renders all instructions onto s.
func (x *XVEC) Draw(s *Surface) {
	for _, inst := range x.Instructions {
		inst.Draw(s)
	}
}

// V is a shorthand for Vertex{X: x, Y: y}.
func V(x, y float32) Vertex { return Vertex{X: x, Y: y} }

func addStep(fill *FillInstruction, stroke *StrokeInstruction, step Stepper) {
	if fill != nil {
		fill.Steps = append(fill.Steps, step)
	}
	if stroke != nil {
		stroke.Steps = append(stroke.Steps, step)
	}
}

func ftos(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

func coltos(c Color) string {
	v := uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
	return fmt.Sprintf("#%08x", v)
}

// scannerParser wraps text/scanner.Scanner with error propagation.
type scannerParser struct {
	s   scanner.Scanner
	err error
}

// raw advances to the next token and returns its text. Sets error on EOF.
func (p *scannerParser) raw() string {
	if p.err != nil {
		return ""
	}
	tok := p.s.Scan()
	if tok == scanner.EOF {
		p.err = fmt.Errorf("unexpected end of input")
		return ""
	}
	return p.s.TokenText()
}

// float reads the next token as a float32. Sets error on parse failure.
func (p *scannerParser) float() float32 {
	txt := p.raw()
	if p.err != nil {
		return 0
	}
	f, err := strconv.ParseFloat(txt, 32)
	if err != nil {
		p.err = fmt.Errorf("expected number, got %q", txt)
		return 0
	}
	return float32(f)
}

// ident reads the next token as a string identifier.
func (p *scannerParser) ident() string {
	txt := p.raw()
	if p.err != nil {
		return ""
	}
	return txt
}

// color reads a #RRGGBBAA hex color.
func (p *scannerParser) color() Color {
	tok := p.raw()
	if tok != "#" {
		p.err = fmt.Errorf("expected '#' for color, got %q", tok)
		return Color{}
	}
	return p.readHex()
}

// readHex reads exactly eight hex digits and returns the RGBA colour.
func (p *scannerParser) readHex() Color {
	if p.err != nil {
		return Color{}
	}
	var buf [8]byte
	for i := range buf {
		ch := p.s.Next()
		if ch < 0 || !isHex(byte(ch)) {
			p.err = fmt.Errorf("expected 8 hex digits for colour, got %d", i)
			return Color{}
		}
		buf[i] = byte(ch)
	}
	// Reject if more hex digits follow — colour must be exactly #RRGGBBAA.
	if ch := p.s.Peek(); ch >= 0 && isHex(byte(ch)) {
		p.err = fmt.Errorf("colour must be exactly 8 hex digits")
		return Color{}
	}
	v, err := strconv.ParseUint(string(buf[:]), 16, 64)
	if err != nil {
		p.err = fmt.Errorf("invalid colour hex %q", string(buf[:]))
		return Color{}
	}
	return Color{
		R: uint8(v >> 24),
		G: uint8(v >> 16),
		B: uint8(v >> 8),
		A: uint8(v),
	}
}

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}
