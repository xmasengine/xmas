// xvec implements an extremely simple vector graphics file format that
// can be drawn directly using ebiten/v2/vector package and that is suited for
// low resolutions and pixel art.
// This means it supports the following features:
// - circle outlines and filled circles
// - rectangle outlines and filled rectangles
// - line strokes
// - filled and outline paths with lines, arcs, cubic and quad beziers.
// - optional antialiasing
//
// Text, fonts, and image resources are not supported.
// The xvec format is encoded from and to a simple text format
// that can be parsed using the text/Scanner package.
package xvec

import "image/color"
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/vector"

// Surface is a drawable image.
type Surface = ebiten.Image

// Color is an RGBA alpha pre multiplied color
type Color = color.RGBA

// Path is a vector path
type Path = vector.Path

type FillOptions = vector.FillOptions
type DrawOptions = vector.DrawPathOptions

// Vertex is a vertex in a vector path.
// In XVEC, all vertices and lengths are stored pre-multiplied with the
// size of the image and expressed in pixels.
type Vertex struct {
	X float32
	Y float32
}

// Length is a length in a vector path.
// In XVEC, all vertices and lengths are stored pre-multiplied with the
// size of the image and expressed in pixels.
type Length float32

// Size is a size in a vector path
type Size struct {
	W float32
	H float32
}

// XVEC is a simplified vector graphic.
type XVEC struct {
	Size         Size          // drawing size
	Antialias    bool          // draw using antialiasing or not
	Instructions []Instruction // drawing instructions
}

// Instruction are drawing instructions for XVEC
type Instruction interface {
	Draw(*Surface)
}

type CircleInstruction struct {
	C         Vertex
	R         Length
	Color     Color
	Stroke    Length
	Antialias bool
}

func Circle(cx, cy, r, stroke float32, col Color) *CircleInstruction {
	return &CircleInstruction{C: Vertex{X: cx, Y: cy}, R: Length(r), Stroke: Length(r)}
}

func (c CircleInstruction) Draw(s *Surface) {
	vector.StrokeCircle(s, c.C.X, c.C.Y, float32(c.R), float32(c.Stroke), c.Color, c.Antialias)
}

func (x *XVEC) Circle(cx, cy, r, stroke float32, col Color) *CircleInstruction {
	c := Circle(cx, cy, r, stroke, col)
	x.Instructions = append(x.Instructions, c)
	return c
}

type DiskInstruction struct {
	C         Vertex
	R         Length
	Color     Color
	Antialias bool
}

func Disk(cx, cy, r float32, col Color) *DiskInstruction {
	return &DiskInstruction{C: Vertex{X: cx, Y: cy}, R: Length(r), Color: col}
}

func (x *XVEC) Disk(cx, cy, r float32, col Color) *DiskInstruction {
	d := Disk(cx, cy, r, col)
	x.Instructions = append(x.Instructions, d)
	return d
}

func (d DiskInstruction) Draw(s *Surface) {
	vector.FillCircle(s, d.C.X, d.C.Y, float32(d.R), d.Color, d.Antialias)
}

type FillInstruction struct {
	Color Color
	Fill  *FillOptions
	Opts  *DrawOptions
	Steps []Stepper
}

func Fill() *FillInstruction {
	return &FillInstruction{}
}

func (f FillInstruction) Draw(s *Surface) {
	var path Path
	for _, step := range f.Steps {
		step.Step(&path)
	}
	vector.FillPath(s, &path, f.Fill, f.Opts)
}

type Stepper interface {
	Step(p *Path)
}

type ArcStep struct {
	C         Vertex
	R         Length
	Start     float32
	End       float32
	Direction vector.Direction
}

func (a ArcStep) Step(p Path) {
	p.Arc(a.C.X, a.C.Y, float32(a.R), a.Start, a.End, a.Direction)
}

/*
func FillCircle(dst *ebiten.Image, cx, cy, r float32, clr color.Color, antialias bool)
func FillPath(dst *ebiten.Image, path *Path, fillOptions *FillOptions, ...)
func FillRect(dst *ebiten.Image, x, y, width, height float32, clr color.Color, ...)
func StrokeCircle(dst *ebiten.Image, cx, cy, r float32, strokeWidth float32, clr color.Color, ...)
func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float32, strokeWidth float32, ...)
func StrokePath(dst *ebiten.Image, path *Path, strokeOptions *StrokeOptions, ...)
func StrokeRect(dst *ebiten.Image, x, y, width, height float32, strokeWidth float32, ...)
type AddPathOptions
type AddStrokeOptions
type Direction
type DrawPathOptions
type FillOptions
type FillRule
type LineCap
type LineJoin
type Path
func (p *Path) AddPath(src *Path, options *AddPathOptions)
func (p *Path) AddStroke(src *Path, options *AddStrokeOptions)
func (p *Path) AppendVerticesAndIndicesForFilling(vertices []ebiten.Vertex, indices []uint16) ([]ebiten.Vertex, []uint16)deprecated
func (p *Path) AppendVerticesAndIndicesForStroke(vertices []ebiten.Vertex, indices []uint16, op *StrokeOptions) ([]ebiten.Vertex, []uint16)deprecated
func (p *Path) Arc(x, y, radius, startAngle, endAngle float32, dir Direction)
func (p *Path) ArcTo(x1, y1, x2, y2, radius float32)
func (p *Path) Bounds() image.Rectangle
func (p *Path) Close()
func (p *Path) CubicTo(x1, y1, x2, y2, x3, y3 float32)
func (p *Path) LineTo(x, y float32)
func (p *Path) MoveTo(x, y float32)
func (p *Path) QuadTo(x1, y1, x2, y2 float32)
func (p *Path) Reset()
type StrokeOptions


func FillCircle(dst *ebiten.Image, cx, cy, r float32, clr color.Color, antialias bool)
func FillPath(dst *ebiten.Image, path *Path, fillOptions *FillOptions, ...)
func FillRect(dst *ebiten.Image, x, y, width, height float32, clr color.Color, ...)
func StrokeCircle(dst *ebiten.Image, cx, cy, r float32, strokeWidth float32, clr color.Color, ...)
func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float32, strokeWidth float32, ...)
func StrokePath(dst *ebiten.Image, path *Path, strokeOptions *StrokeOptions, ...)
func StrokeRect(dst *ebiten.Image, x, y, width, height float32, strokeWidth float32, ...)


type AddPathOptions
type AddStrokeOptions
type Direction
type DrawPathOptions
type FillOptions
type FillRule
type LineCap
type LineJoin
type Path
func (p *Path) AddPath(src *Path, options *AddPathOptions)
func (p *Path) AddStroke(src *Path, options *AddStrokeOptions)
func (p *Path) AppendVerticesAndIndicesForFilling(vertices []ebiten.Vertex, indices []uint16) ([]ebiten.Vertex, []uint16)deprecated
func (p *Path) AppendVerticesAndIndicesForStroke(vertices []ebiten.Vertex, indices []uint16, op *StrokeOptions) ([]ebiten.Vertex, []uint16)deprecated
func (p *Path) Arc(x, y, radius, startAngle, endAngle float32, dir Direction)
func (p *Path) ArcTo(x1, y1, x2, y2, radius float32)
func (p *Path) Bounds() image.Rectangle
func (p *Path) Close()
func (p *Path) CubicTo(x1, y1, x2, y2, x3, y3 float32)
func (p *Path) LineTo(x, y float32)
func (p *Path) MoveTo(x, y float32)
func (p *Path) QuadTo(x1, y1, x2, y2 float32)
func (p *Path) Reset()
type StrokeOptions

bjorn@ancelot~/src/xmas$ go doc -all vector
package vector // import "github.com/hajimehoshi/ebiten/v2/vector"

Package vector provides functions for vector graphics rendering.

This package is under experiments and the API might be changed with breaking
backward compatibility.

FUNCTIONS

func FillCircle(dst *ebiten.Image, cx, cy, r float32, clr color.Color, antialias bool)
    FillCircle fills a circle with the specified center position (cx, cy),
    the radius (r), width and color.

func FillPath(dst *ebiten.Image, path *Path, fillOptions *FillOptions, drawPathOptions *DrawPathOptions)
    FillPath fills the specified path with the specified options.

func FillRect(dst *ebiten.Image, x, y, width, height float32, clr color.Color, antialias bool)
    FillRect fills a rectangle with the specified width and color.

func StrokeCircle(dst *ebiten.Image, cx, cy, r float32, strokeWidth float32, clr color.Color, antialias bool)
    StrokeCircle strokes a circle with the specified center position (cx, cy),
    the radius (r), width and color.

func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float32, strokeWidth float32, clr color.Color, antialias bool)
    StrokeLine strokes a line (x0, y0)-(x1, y1) with the specified width and
    color.

func StrokePath(dst *ebiten.Image, path *Path, strokeOptions *StrokeOptions, drawPathOptions *DrawPathOptions)
    StrokePath strokes the specified path with the specified options.

func StrokeRect(dst *ebiten.Image, x, y, width, height float32, strokeWidth float32, clr color.Color, antialias bool)
    StrokeRect strokes a rectangle with the specified width and color.


TYPES

type AddPathOptions struct {
	// GeoM is a geometry matrix to apply to the path.
	//
	// The default (zero) value is an identity matrix.
	GeoM ebiten.GeoM
}
    AddPathOptions is options for Path.AddPath.

type AddStrokeOptions struct {
	// StrokeOptions is options for the stroke.
	StrokeOptions

	// GeoM is a geometry matrix to apply to the path.
	//
	// The default (zero) value is an identity matrix.
	GeoM ebiten.GeoM
}
    AddStrokeOptions is options for Path.AddStroke.

type Direction int
    Direction represents clockwise or counterclockwise.

const (
	Clockwise Direction = iota
	CounterClockwise
)
type DrawPathOptions struct {
	// AntiAlias is whether the path is drawn with anti-aliasing.
	// The default (zero) value is false.
	AntiAlias bool

	// ColorScale is the color scale to apply to the path.
	// The default (zero) value is identity, which is (1, 1, 1, 1) (white).
	ColorScale ebiten.ColorScale

	// Blend is the blend mode to apply to the path.
	// The default (zero) value is ebiten.BlendSourceOver.
	Blend ebiten.Blend
}
    DrawPathOptions is options to draw a path.

type FillOptions struct {
	// FillRule is the rule whether an overlapped region is rendered or not.
	// The default (zero) value is FillRuleNonZero.
	FillRule FillRule
}
    FillOptions is options to fill a path.

type FillRule int
    FillRule is the rule whether an overlapped region is rendered or not.

const (
	// FillRuleNonZero means that triangles are rendered based on the non-zero rule.
	// If and only if the number of overlaps is not 0, the region is rendered.
	FillRuleNonZero FillRule = iota

	// FillRuleEvenOdd means that triangles are rendered based on the even-odd rule.
	// If and only if the number of overlaps is odd, the region is rendered.
	FillRuleEvenOdd
)
type LineCap int
    LineCap represents the way in which how the ends of the stroke are rendered.

const (
	LineCapButt LineCap = iota
	LineCapRound
	LineCapSquare
)
type LineJoin int
    LineJoin represents the way in which how two segments are joined.

const (
	LineJoinMiter LineJoin = iota
	LineJoinBevel
	LineJoinRound
)
type Path struct {
	// Has unexported fields.
}
    Path represents a collection of vector graphics operations.

func (p *Path) AddPath(src *Path, options *AddPathOptions)
    AddPath adds the given path src to this path p as a sub-path.

func (p *Path) AddStroke(src *Path, options *AddStrokeOptions)
    AddStroke adds a stroke path to the path p.

    The added stroke path must be rendered with FileRuleNonZero.

func (p *Path) AppendVerticesAndIndicesForFilling(vertices []ebiten.Vertex, indices []uint16) ([]ebiten.Vertex, []uint16)
    AppendVerticesAndIndicesForFilling appends vertices and indices to fill this
    path and returns them.

    AppendVerticesAndIndicesForFilling works in a similar way
    to the built-in append function. If the arguments are nils,
    AppendVerticesAndIndicesForFilling returns new slices.

    The returned vertice's SrcX and SrcY are 0, and ColorR, ColorG, ColorB,
    and ColorA are 1.

    The returned values are intended to be passed to DrawTriangles or
    DrawTrianglesShader with FileRuleNonZero or FillRuleEvenOdd in order to
    render a complex polygon like a concave polygon, a polygon with holes,
    or a self-intersecting polygon.

    The returned vertices and indices should be rendered with a solid
    (non-transparent) color with the default Blend (source-over). Otherwise,
    there is no guarantee about the rendering result.

    Deprecated: as of v2.9. Use FillPath instead.

func (p *Path) AppendVerticesAndIndicesForStroke(vertices []ebiten.Vertex, indices []uint16, op *StrokeOptions) ([]ebiten.Vertex, []uint16)
    AppendVerticesAndIndicesForStroke appends vertices and indices to render
    a stroke of this path and returns them. AppendVerticesAndIndicesForStroke
    works in a similar way to the built-in append function. If the arguments are
    nils, AppendVerticesAndIndicesForStroke returns new slices.

    The returned vertice's SrcX and SrcY are 0, and ColorR, ColorG, ColorB,
    and ColorA are 1.

    The returned values are intended to be passed to DrawTriangles
    or DrawTrianglesShader with a solid (non-transparent) color with
    FillRuleFillAll or FillRuleNonZero, not FileRuleEvenOdd.

    Deprecated: as of v2.9. Use StrokePath or Path.AddStroke instead.

func (p *Path) Arc(x, y, radius, startAngle, endAngle float32, dir Direction)
    Arc adds an arc to the path. (x, y) is the center of the arc.

func (p *Path) ArcTo(x1, y1, x2, y2, radius float32)
    ArcTo adds an arc curve to the path. (x1, y1) is the first control point,
    and (x2, y2) is the second control point.

func (p *Path) Bounds() image.Rectangle
    Bounds returns the minimum bounding rectangle of the path.

func (p *Path) Close()
    Close adds a new line from the last position of the current sub-path to
    the first position of the current sub-path, and marks the current sub-path
    closed. Following operations for this path will start with a new sub-path.

func (p *Path) CubicTo(x1, y1, x2, y2, x3, y3 float32)
    CubicTo adds a cubic Bézier curve to the path. (x1, y1) and (x2, y2) are the
    control points, and (x3, y3) is the destination.

func (p *Path) LineTo(x, y float32)
    LineTo adds a line segment to the path, which starts from the last position
    of the current sub-path and ends to the given position (x, y). If p doesn't
    have any sub-paths or the last sub-path is closed, LineTo sets (x, y) as the
    start position of a new sub-path.

func (p *Path) MoveTo(x, y float32)
    MoveTo starts a new sub-path with the given position (x, y) without adding a
    sub-path,

func (p *Path) QuadTo(x1, y1, x2, y2 float32)
    QuadTo adds a quadratic Bézier curve to the path. (x1, y1) is the control
    point, and (x2, y2) is the destination.

func (p *Path) Reset()
    Reset resets the path. Reset doesn't release the allocated memory so that
    the memory can be reused.

type StrokeOptions struct {
	// Width is the stroke width in pixels.
	//
	// The default (zero) value is 0.
	Width float32

	// LineCap is the way in which how the ends of the stroke are rendered.
	// Line caps are not rendered when the sub-path is marked as closed.
	//
	// The default (zero) value is [LineCapButt].
	LineCap LineCap

	// LineJoin is the way in which how two segments are joined.
	//
	// The default (zero) value is [LineJoinMiter].
	LineJoin LineJoin

	// MiterLimit is the miter limit for [LineJoinMiter].
	// For details, see https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-miterlimit.
	//
	// The default (zero) value is 0.
	MiterLimit float32
}
    StrokeOptions is options to render a stroke.


*/
