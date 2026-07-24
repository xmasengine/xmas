// Package xgal, short for Xmasengine GAme Library,
// provides types and functions for writing 2D games in Go.
// It covers graphics, audio, video, input, text, and fonts in a
// short-named, standalone API. All external types are aliased so
// users never need to import the underlying libraries directly.
package xgal

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Color is any color value. Use [RGBA] to construct concrete colors.
type Color = color.Color

// RGBA is an 8-bit-per-channel RGBA color.
type RGBA = color.RGBA

// Palette is a set of colors, typically for indexed images.
type Palette = color.Palette

// Paletted is an image with a palette.
type Paletted = image.Paletted

// PalettedImage is an image interface with a palette.
type PalettedImage = image.PalettedImage

// Rectangle is a 2D rectangle defined by two [Point] values.
// Use [Rect] to construct one.
type Rectangle = image.Rectangle

// Point is a 2D point with integer coordinates. Use [Pt] to construct one.
type Point = image.Point

// Surface is a drawable image. The game screen is a *Surface,
// and loaded textures are *Surfaces. They can be drawn onto each
// other with [Blit], [Blend], [Scale], or the vector drawing functions.
type Surface = ebiten.Image

// Image is an loaded bit map image.
type Image = image.Image

// DrawOptions specifies how one [Surface] is drawn onto another.
// Key fields:
//
//	GeoM      — affine transform (translate, scale, rotate)
//	Filter    — [Nearest], [Linear], or [Pixelated]
//	Blend     — [BlendNormal], [BlendCopy], etc.
//	ColorScale— multiply color and alpha
type DrawOptions = ebiten.DrawImageOptions

// Pt returns a [Point] with the given coordinates.
func Pt(x, y int) Point { return image.Pt(x, y) }

// Rect returns a [Rectangle] from (x0,y0) to (x1,y1).
func Rect(x0, y0, x1, y1 int) Rectangle { return image.Rect(x0, y0, x1, y1) }

// Bound returns a [Rectangle] from (x0,y0) to (x0+w,y0+h).
func Bound(x0, y0, w, h int) Rectangle { return image.Rect(x0, y0, x0+w, y0+h) }

// Paint returns a color with the given channel values. Prefer this to
// RGBA{R, G, B, A} to avoid vet warnings about unkeyed fields.
func Paint(r, g, b, a uint8) RGBA { return RGBA{r, g, b, a} }

// Wash returns a color with the given channel values. Prefer this to
// RGBA{R, G, B, A} to avoid vet warnings about unkeyed fields.
func Wash(r, g, b, a uint8) RGBA { return RGBA{r, g, b, a} }

// Tint returns an opaque color with the given channel values.
func Tint(r, g, b uint8) RGBA { return RGBA{r, g, b, 255} }

// Common RGBA colors.
var (
	Black       = RGBA{A: 255}
	White       = RGBA{R: 255, G: 255, B: 255, A: 255}
	Transparent = RGBA{}
	Opaque      = RGBA{A: 255}
)

// Recolor converts a coler to an RGBA
func Recolor(c Color) RGBA {
	r16, g16, b16, a16 := c.RGBA()

	r8 := uint8(r16 >> 8)
	g8 := uint8(g16 >> 8)
	b8 := uint8(b16 >> 8)
	a8 := uint8(a16 >> 8)
	return Paint(r8, g8, b8, a8)
}
