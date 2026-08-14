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
	"io/fs"
	"math"
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

func (v Vertex) Add(u Vertex) Vertex {
	return Vertex{X: v.X + u.X, Y: v.Y + u.Y}
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
	_ encoding.TextMarshaler = (*UseInstruction)(nil)
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
// a list of drawing instructions, and named sub-graphics.
type XVEC struct {
	Size         Size
	Antialias    bool
	Instructions []Instruction
	Defs         []Def // named sub-graphics, in definition order
	cached       *Surface
	drawDepth    int // recursion guard for cyclic defs
}

// Def is a named sub-graphic that can be drawn with the use instruction.
type Def struct {
	Name string
	Vec  *XVEC
}

// findDef returns the named sub-graphic, or nil if it does not exist.
func (x *XVEC) findDef(name string) *XVEC {
	for _, d := range x.Defs {
		if d.Name == name {
			return d.Vec
		}
	}
	return nil
}

// Def returns the named sub-graphic, creating and registering it if needed.
func (x *XVEC) Def(name string) *XVEC {
	if vec := x.findDef(name); vec != nil {
		return vec
	}
	vec := &XVEC{Size: x.Size, Antialias: x.Antialias}
	x.Defs = append(x.Defs, Def{Name: name, Vec: vec})
	return vec
}

// Use adds a use instruction that draws the named sub-graphic, optionally
// scaled by scale, rotated by rotate radians, and translated to at.
// The sub-graphic must have been defined before the use.
func (x *XVEC) Use(name string, at Vertex, scale Size, rotate float32) *UseInstruction {
	u := &UseInstruction{Name: name, At: at, Scale: scale, Rotate: rotate, Vec: x.findDef(name)}
	x.Instructions = append(x.Instructions, u)
	return u
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

// Adjuster is an instruction that has a stroke size that can be adjusted.
type Adjuster interface {
	Adjust(stroke Length)
}

// Painter is an instruction that has a color that can be painted to change it.
type Painter interface {
	Paint(color Color)
}

// CircleInstruction strokes a circle outline.
type CircleInstruction struct {
	C         Vertex
	R         Length
	Color     Color
	Stroke    Length
	Antialias bool
}

func (c *CircleInstruction) Paint(color Color) {
	c.Color = color
}

func (c *CircleInstruction) Adjust(stroke Length) {
	c.Stroke = stroke
}

func (x *XVEC) Circle(cx, cy, r, stroke float32, col Color) *CircleInstruction {
	c := &CircleInstruction{
		C: Vertex{cx, cy}, R: Length(r), Stroke: Length(stroke),
		Color: col, Antialias: x.Antialias,
	}
	x.Instructions = append(x.Instructions, c)
	return c
}

func (c CircleInstruction) Draw(s *Surface) {
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

func (c *DiskInstruction) Paint(color Color) {
	c.Color = color
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

func (r *RectInstruction) Paint(color Color) {
	r.Color = color
}

func (r *RectInstruction) Adjust(stroke Length) {
	r.Stroke = stroke
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

func (s *SlabInstruction) Paint(color Color) {
	s.Color = color
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

func (r SlabInstruction) Render(s *Surface, at Vertex, scale Size) {
	x := r.X + at.X
	y := r.Y + at.Y
	w := r.W * scale.W
	h := r.H * scale.H
	vector.FillRect(s, x, y, w, h, r.Color, r.Antialias)
}

// LineInstruction strokes a line segment.
type LineInstruction struct {
	X1, Y1, X2, Y2 float32
	Color          Color
	Stroke         Length
	Antialias      bool
}

func (l *LineInstruction) Adjust(stroke Length) {
	l.Stroke = stroke
}

func (l *LineInstruction) Paint(color Color) {
	l.Color = color
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

// UseInstruction draws a previously defined sub-graphic, optionally scaled,
// rotated, and offset.
type UseInstruction struct {
	Name   string
	At     Vertex
	Scale  Size
	Rotate float32 // radians
	Vec    *XVEC   // resolved sub-graphic
}

func (u *UseInstruction) Draw(s *Surface) {
	if u.Vec == nil {
		return
	}
	var geom ebiten.GeoM
	geom.Translate(float64(u.At.X), float64(u.At.Y))
	if u.Rotate != 0 {
		geom.Rotate(float64(u.Rotate))
	}
	if u.Scale != (Size{}) {
		geom.Scale(float64(u.Scale.W), float64(u.Scale.H))
	}
	u.Vec.DrawTransformed(s, geom)
}

func (u *UseInstruction) MarshalText() ([]byte, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "use %s", u.Name)
	if u.At != (Vertex{}) {
		fmt.Fprintf(&b, " at %s %s", ftos(u.At.X), ftos(u.At.Y))
	}
	if u.Scale != (Size{}) {
		fmt.Fprintf(&b, " scale %s %s", ftos(u.Scale.W), ftos(u.Scale.H))
	}
	if u.Rotate != 0 {
		fmt.Fprintf(&b, " rotate %s", ftos(u.Rotate*180/math.Pi))
	}
	return []byte(b.String()), nil
}

// FillInstruction fills a path built from steps.
type FillInstruction struct {
	Color     Color
	FillOpts  FillOptions
	DrawOpts  DrawPathOptions
	Steps     []Stepper
	Antialias bool
}

func (f *FillInstruction) Paint(color Color) {
	f.Color = color
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
	Stroke     Length
	StrokeOpts StrokeOptions
	DrawOpts   DrawPathOptions
	Steps      []Stepper
	Antialias  bool
}

func (s *StrokeInstruction) Paint(color Color) {
	s.Color = color
}

func (s *StrokeInstruction) Adjust(stroke Length) {
	s.Stroke = stroke
	s.StrokeOpts.Width = float32(s.Stroke)
}

func (x *XVEC) Stroke(stroke float32, col Color, steps ...Stepper) *StrokeInstruction {
	s := &StrokeInstruction{
		Color: col, Stroke: Length(stroke), Steps: steps, Antialias: x.Antialias,
	}
	s.DrawOpts.AntiAlias = x.Antialias
	s.DrawOpts.ColorScale.Reset()
	s.StrokeOpts.Width = float32(s.Stroke)
	x.Instructions = append(x.Instructions, s)
	return s
}

func (s *StrokeInstruction) Draw(dst *Surface) {
	var path Path
	for _, step := range s.Steps {
		step.Step(&path)
	}
	so := s.StrokeOpts
	so.Width = float32(s.Stroke)

	opts := s.DrawOpts
	opts.ColorScale.ScaleWithColor(s.Color)

	vector.StrokePath(dst, &path, &so, &opts)
}

func (s *StrokeInstruction) MarshalText() ([]byte, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "stroke %s %s\n", ftos(float32(s.Stroke)), coltos(s.Color))
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
	for _, d := range x.Defs {
		if _, err := fmt.Fprintf(w, "def %s\n", d.Name); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "size %s %s\n", ftos(d.Vec.Size.W), ftos(d.Vec.Size.H)); err != nil {
			return err
		}
		aa := "false"
		if d.Vec.Antialias {
			aa = "true"
		}
		if _, err := fmt.Fprintf(w, "antialias %s\n", aa); err != nil {
			return err
		}
		if err := encodeInstructions(w, d.Vec.Instructions); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, "end\n"); err != nil {
			return err
		}
	}
	return encodeInstructions(w, x.Instructions)
}

func encodeInstructions(w io.Writer, instructions []Instruction) error {
	for _, inst := range instructions {
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
	x.Defs = nil

	p := &scannerParser{}
	p.s.Init(r)
	p.s.Mode = scanner.ScanInts | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanComments | scanner.SkipComments

	var curFill *FillInstruction
	var curStroke *StrokeInstruction
	var curDef *XVEC
	var curUse *UseInstruction

	addInst := func(inst Instruction) {
		if curDef != nil {
			curDef.Instructions = append(curDef.Instructions, inst)
		} else {
			x.Instructions = append(x.Instructions, inst)
		}
	}
	currentAA := func() bool {
		if curDef != nil {
			return curDef.Antialias
		}
		return x.Antialias
	}
	seenInstruction := func() bool {
		if curFill != nil || curStroke != nil {
			return true
		}
		if curDef != nil {
			return len(curDef.Instructions) > 0
		}
		return len(x.Instructions) > 0 || len(x.Defs) > 0
	}

	tok := p.s.Scan()
	for tok != scanner.EOF {
		if p.err != nil {
			return p.err
		}
		kw := p.s.TokenText()
		inPath := curFill != nil || curStroke != nil
		if kw != "at" && kw != "scale" && kw != "rotate" {
			curUse = nil
		}
		switch kw {
		case "xvec":
			v := p.raw()
			if v != "1" {
				return fmt.Errorf("unsupported xvec version %q", v)
			}

		case "size":
			if seenInstruction() {
				return fmt.Errorf("size must come before all drawing instructions")
			}
			target := x
			if curDef != nil {
				target = curDef
			}
			target.Size.W = p.float()
			target.Size.H = p.float()

		case "antialias":
			if seenInstruction() {
				return fmt.Errorf("antialias must come before all drawing instructions")
			}
			target := x
			if curDef != nil {
				target = curDef
			}
			target.Antialias = p.ident() == "true"

		case "def":
			if curDef != nil {
				return fmt.Errorf("def cannot be nested")
			}
			if inPath {
				return fmt.Errorf("def not allowed inside a path")
			}
			name := p.ident()
			def := &XVEC{Size: x.Size, Antialias: x.Antialias}
			x.Defs = append(x.Defs, Def{Name: name, Vec: def})
			curDef = def

		case "use":
			if inPath {
				return fmt.Errorf("use not allowed inside a path")
			}
			name := p.ident()
			u := &UseInstruction{Name: name, Vec: x.findDef(name)}
			if u.Vec == nil {
				return fmt.Errorf("use %q: undefined def", name)
			}
			addInst(u)
			curUse = u

		case "at", "scale", "rotate":
			if curUse == nil {
				return fmt.Errorf("%s only allowed after use", kw)
			}
			switch kw {
			case "at":
				curUse.At = V(p.float(), p.float())
			case "scale":
				curUse.Scale = Size{p.float(), p.float()}
				if curUse.Scale.W == 0 || curUse.Scale.H == 0 {
					return fmt.Errorf("scale components must be non-zero")
				}
			case "rotate":
				curUse.Rotate = float32(p.float()) * math.Pi / 180
			}

		case "rule":
			if curFill == nil {
				return fmt.Errorf("rule only allowed inside a fill")
			}
			switch p.ident() {
			case "evenodd":
				curFill.FillOpts.FillRule = vector.FillRuleEvenOdd
			case "nonzero":
				curFill.FillOpts.FillRule = vector.FillRuleNonZero
			default:
				return fmt.Errorf("unknown fill rule")
			}

		case "cap":
			if curStroke == nil {
				return fmt.Errorf("cap only allowed inside a stroke")
			}
			switch p.ident() {
			case "butt":
				curStroke.StrokeOpts.LineCap = vector.LineCapButt
			case "round":
				curStroke.StrokeOpts.LineCap = vector.LineCapRound
			case "square":
				curStroke.StrokeOpts.LineCap = vector.LineCapSquare
			default:
				return fmt.Errorf("unknown cap type")
			}

		case "join":
			if curStroke == nil {
				return fmt.Errorf("join only allowed inside a stroke")
			}
			switch p.ident() {
			case "miter":
				curStroke.StrokeOpts.LineJoin = vector.LineJoinMiter
			case "round":
				curStroke.StrokeOpts.LineJoin = vector.LineJoinRound
			case "bevel":
				curStroke.StrokeOpts.LineJoin = vector.LineJoinBevel
			default:
				return fmt.Errorf("unknown cap type")
			}

		case "circle":
			c := &CircleInstruction{C: V(p.float(), p.float()), R: Length(p.float()), Stroke: Length(p.float()), Color: p.color(), Antialias: currentAA()}
			addInst(c)

		case "disk":
			d := &DiskInstruction{C: V(p.float(), p.float()), R: Length(p.float()), Color: p.color(), Antialias: currentAA()}
			addInst(d)

		case "rect":
			r := &RectInstruction{X: p.float(), Y: p.float(), W: p.float(), H: p.float(), Stroke: Length(p.float()), Color: p.color(), Antialias: currentAA()}
			addInst(r)

		case "slab":
			fr := &SlabInstruction{X: p.float(), Y: p.float(), W: p.float(), H: p.float(), Color: p.color(), Antialias: currentAA()}
			addInst(fr)

		case "line":
			if inPath {
				addStep(curFill, curStroke, &LineStep{X: p.float(), Y: p.float()})
			} else {
				l := &LineInstruction{X1: p.float(), Y1: p.float(), X2: p.float(), Y2: p.float(), Stroke: Length(p.float()), Color: p.color(), Antialias: currentAA()}
				addInst(l)
			}

		case "fill":
			curFill = &FillInstruction{Color: p.color(), Steps: nil, Antialias: currentAA()}
			if currentAA() {
				curFill.DrawOpts.AntiAlias = true
			}

		case "stroke":
			w := p.float()
			curStroke = &StrokeInstruction{Color: p.color(), Stroke: Length(w), Steps: nil, Antialias: currentAA()}
			if currentAA() {
				curStroke.DrawOpts.AntiAlias = true
			}

		case "end":
			if curFill != nil {
				if len(curFill.Steps) > 0 {
					if _, ok := curFill.Steps[len(curFill.Steps)-1].(*CloseStep); !ok {
						return fmt.Errorf("fill path must end with close")
					}
				}
				addInst(curFill)
				curFill = nil
			} else if curStroke != nil {
				if len(curStroke.Steps) > 0 {
					if _, ok := curStroke.Steps[len(curStroke.Steps)-1].(*CloseStep); !ok {
						return fmt.Errorf("stroke path must end with close")
					}
				}
				addInst(curStroke)
				curStroke = nil
			} else if curDef != nil {
				if len(curDef.Instructions) == 0 {
					return fmt.Errorf("def block must not be empty")
				}
				curDef = nil
			} else {
				return fmt.Errorf("end without an open fill, stroke, or def")
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

// Draw renders all instructions onto s at their natural size and position.
func (x *XVEC) Draw(s *Surface) {
	for _, inst := range x.Instructions {
		inst.Draw(s)
	}
}

// maxDefDepth limits def/use nesting during drawing; exceeding it aborts the
// draw to avoid unbounded recursion on cyclic defs.
const maxDefDepth = 16

// DrawTransformed renders x onto dst, transformed by geom. The graphic is
// rasterized once at its natural size and cached; subsequent calls reuse the
// rasterization. Callers MUST NOT modify x after the first transformed draw.
func (x *XVEC) DrawTransformed(dst *Surface, geom ebiten.GeoM) {
	if x.drawDepth > maxDefDepth {
		return
	}
	x.drawDepth++
	defer func() { x.drawDepth-- }()
	if x.cached == nil {
		w, h := int(x.Size.W), int(x.Size.H)
		if w < 1 {
			w = 1
		}
		if h < 1 {
			h = 1
		}
		x.cached = ebiten.NewImage(w, h)
		x.Draw(x.cached)
	}
	opts := &ebiten.DrawImageOptions{GeoM: geom}
	dst.DrawImage(x.cached, opts)
}

// DrawScaled renders x onto dst at position at, scaled so that the canvas
// fits the given size.
func (x *XVEC) DrawScaled(dst *Surface, at Vertex, size Size) {
	var geom ebiten.GeoM
	geom.Translate(float64(at.X), float64(at.Y))
	if x.Size.W > 0 && x.Size.H > 0 && size.W > 0 && size.H > 0 {
		geom.Scale(float64(size.W)/float64(x.Size.W), float64(size.H)/float64(x.Size.H))
	}
	x.DrawTransformed(dst, geom)
}

// V is a shorthand for Vertex{X: x, Y: y}.
func V(x, y float32) Vertex { return Vertex{X: x, Y: y} }

// MoveUp swaps the instruction at index i with the one before it.
// It is a no-op if i <= 0.
func (x *XVEC) MoveUp(i int) {
	if i <= 0 || i >= len(x.Instructions) {
		return
	}
	x.Instructions[i], x.Instructions[i-1] = x.Instructions[i-1], x.Instructions[i]
}

// MoveDown swaps the instruction at index i with the one after it.
// It is a no-op if i >= len(Instructions)-1.
func (x *XVEC) MoveDown(i int) {
	if i < 0 || i >= len(x.Instructions)-1 {
		return
	}
	x.Instructions[i], x.Instructions[i+1] = x.Instructions[i+1], x.Instructions[i]
}

// MoveToFront moves the instruction at index i to index 0
// (the bottom of the draw order). It is a no-op if i <= 0.
func (x *XVEC) MoveToFront(i int) {
	if i <= 0 || i >= len(x.Instructions) {
		return
	}
	inst := x.Instructions[i]
	x.Instructions = append(x.Instructions[:i], x.Instructions[i+1:]...)
	x.Instructions = append([]Instruction{inst}, x.Instructions...)
}

// MoveToBack moves the instruction at index i to the last index
// (the top of the draw order). It is a no-op if i >= len(Instructions)-1.
func (x *XVEC) MoveToBack(i int) {
	if i < 0 || i >= len(x.Instructions)-1 {
		return
	}
	inst := x.Instructions[i]
	x.Instructions = append(x.Instructions[:i], x.Instructions[i+1:]...)
	x.Instructions = append(x.Instructions, inst)
}

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

func Parse(rd io.Reader) (*XVEC, error) {
	var vec XVEC
	err := vec.Decode(rd)
	if err != nil {
		return nil, err
	}
	return &vec, nil
}

func ParseFS(fs fs.FS, name string) (*XVEC, error) {
	fin, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	return Parse(fin)
}
