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
