package xgal

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TouchID identifies a single touch input.
type TouchID = ebiten.TouchID

// Flick reports whether the specific touch was just pressed this frame.
func Flick(id TouchID) bool {
	for _, tid := range inpututil.AppendJustPressedTouchIDs(nil) {
		if tid == id {
			return true
		}
	}
	return false
}

// Flicks returns all touches that were just pressed this frame.
// If buf is provided, results are appended to it.
func Flicks(buf ...[]TouchID) []TouchID {
	var b []TouchID
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustPressedTouchIDs(b)
}

// Touches returns all currently active touches.
// If buf is provided, results are appended to it.
func Touches(buf ...[]TouchID) []TouchID {
	var b []TouchID
	if len(buf) > 0 {
		b = buf[0]
	}
	return ebiten.AppendTouchIDs(b)
}

// Drops returns all touches that were just released this frame.
// If buf is provided, results are appended to it.
func Drops(buf ...[]TouchID) []TouchID {
	var b []TouchID
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustReleasedTouchIDs(b)
}

// Touch returns the current position of the touch.
func Touch(id TouchID) Point {
	x, y := ebiten.TouchPosition(id)
	return image.Pt(x, y)
}

// LastTouch returns the previous frame position of the touch.
func LastTouch(id TouchID) Point {
	x, y := inpututil.TouchPositionInPreviousTick(id)
	return image.Pt(x, y)
}
