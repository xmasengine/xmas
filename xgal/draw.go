package xgal

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Box draws a filled rectangle.
func Box(dst *Surface, r Rectangle, col RGBA) {
	vector.DrawFilledRect(
		dst, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
		col, false,
	)
}

// Outline draws the outline of a rectangle.
func Outline(dst *Surface, r Rectangle, stroke int, col RGBA) {
	vector.StrokeRect(
		dst, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
		float32(stroke), col, false,
	)
}

// Disk draws a filled circle.
func Disk(dst *Surface, c Point, r int, col RGBA) {
	vector.DrawFilledCircle(dst, float32(c.X), float32(c.Y),
		float32(r), col, false)
}

// Circle draws an empty circle.
func Circle(dst *Surface, c Point, r int, stroke int, col RGBA) {
	vector.StrokeCircle(dst, float32(c.X), float32(c.Y),
		float32(r), float32(stroke), col, false)
}

// Line draws a line from (x1, y1) to (x2, y2) with the given stroke width.
func Line(dst *Surface, x1, y1, x2, y2, stroke int, col RGBA) {
	vector.StrokeLine(dst, float32(x1), float32(y1), float32(x2), float32(y2),
		float32(stroke), col, false)
}

// Clear fills dst with a solid color.
func Clear(dst *Surface, color RGBA) {
	dst.Fill(color)
}

// Path is a vector path that can be stroked with [StrokePath].
// Use its MoveTo, LineTo, QuadTo, CubicTo, and Close methods to build it.
type Path = vector.Path

// Flood fills a vector path with the given color.
func Flood(dst *Surface, path *Path, col RGBA) {
	var opts vector.DrawPathOptions
	opts.ColorScale.ScaleWithColor(col)
	vector.FillPath(dst, path, &vector.FillOptions{}, &opts)
}

// Trace strokes a vector path with the given stroke width and color.
func Trace(dst *Surface, path *Path, width int, col RGBA) {
	var so vector.StrokeOptions
	so.Width = float32(width)
	var opts vector.DrawPathOptions
	opts.ColorScale.ScaleWithColor(col)
	vector.StrokePath(dst, path, &so, &opts)
}

// Quad draws a quadratic Bézier curve from (x0,y0) to (x2,y2)
// with control point (x1,y1).
func Quad(dst *Surface, x0, y0, x1, y1, x2, y2, strokeWidth int, col RGBA) {
	var p vector.Path
	p.MoveTo(float32(x0), float32(y0))
	p.QuadTo(float32(x1), float32(y1), float32(x2), float32(y2))
	var so vector.StrokeOptions
	so.Width = float32(strokeWidth)
	var opts vector.DrawPathOptions
	opts.ColorScale.ScaleWithColor(col)
	vector.StrokePath(dst, &p, &so, &opts)
}

// Cubic draws a cubic Bézier curve from (x0,y0) to (x3,y3)
// with control points (x1,y1) and (x2,y2).
func Cubic(dst *Surface, x0, y0, x1, y1, x2, y2, x3, y3, strokeWidth int, col RGBA) {
	var p vector.Path
	p.MoveTo(float32(x0), float32(y0))
	p.CubicTo(float32(x1), float32(y1), float32(x2), float32(y2), float32(x3), float32(y3))
	var so vector.StrokeOptions
	so.Width = float32(strokeWidth)
	var opts vector.DrawPathOptions
	opts.ColorScale.ScaleWithColor(col)
	vector.StrokePath(dst, &p, &so, &opts)
}

// Andrew draws St Andrews cross, or an X shape.
func Andreas(surface *Surface, r Rectangle, thick int, col RGBA) {
	Line(surface, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, thick, col)
	r.Min.X, r.Max.X = r.Max.X, r.Min.X
	Line(surface, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, thick, col)
}
