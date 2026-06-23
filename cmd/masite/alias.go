package main

import "github.com/hajimehoshi/ebiten/v2"
import "image"
import "image/color"
import "errors"

// A few type aliases for convenience
type (
	Color     = color.Color
	RGBA      = color.RGBA
	Image     = image.Image
	Surface   = ebiten.Image
	Game      = ebiten.Game
	Rectangle = image.Rectangle
	Point     = image.Point
	Key       = ebiten.Key
)

var (
	Termination = ebiten.Termination
	MidgetOK    = errors.New("OK")
)
