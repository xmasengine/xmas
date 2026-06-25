package xgal

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MouseButton is a mouse button.
type MouseButton = ebiten.MouseButton

const (
	// MouseButtonLeft is the left mouse button.
	MouseButtonLeft MouseButton = ebiten.MouseButtonLeft
	// MouseButtonRight is the right mouse button.
	MouseButtonRight MouseButton = ebiten.MouseButtonRight
	// MouseButtonMiddle is the middle mouse button.
	MouseButtonMiddle MouseButton = ebiten.MouseButtonMiddle
)

// Mouse returns the current mouse cursor position.
func Mouse() Point {
	x, y := ebiten.CursorPosition()
	return image.Pt(x, y)
}

// Click reports whether the mouse button was just pressed.
func Click(button MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

// Grip reports whether the mouse button is currently held.
func Grip(button MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

// Loose reports whether the mouse button was just released.
func Loose(button MouseButton) bool {
	return inpututil.IsMouseButtonJustReleased(button)
}

// Wheel returns the scroll wheel movement since the last frame.
// Positive Y scrolls toward the user (down), positive X scrolls right.
func Wheel() (xoff, yoff float64) {
	return ebiten.Wheel()
}
