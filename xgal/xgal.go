// Package xgal is a wrapper around the low level ebitengine game library
// and its supporting libraries, as well as the image and image/color library.
// While ebitengine works well, it is a bit of a hassle to import everything
// separately.
// All types used from other packages are aliased so this package is standalone.
package xgal

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Color is a color.
type Color = color.Color

// RGBA is an 8-bit-per-channel RGBA color.
type RGBA = color.RGBA

// Rectangle is a 2D rectangle defined by two Points.
type Rectangle = image.Rectangle

// Point is a 2D point with integer coordinates.
type Point = image.Point

// Surface is an off-screen image that can be drawn to.
type Surface = ebiten.Image

// DrawOptions specifies options for drawing a Surface onto another Surface.
type DrawOptions = ebiten.DrawImageOptions


