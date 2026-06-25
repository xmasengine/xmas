package xgal

import "github.com/hajimehoshi/ebiten/v2"

// Game is the game interface.
type Game = ebiten.Game

// Quit is the sentinel error returned from [Game.Update] to terminate the game.
var Quit = ebiten.Termination

// Play runs the game. Shorthand for [ebiten.RunGame].
func Play(game Game) error {
	return ebiten.RunGame(game)
}

// Screen sets the window size and title.
func Screen(w, h int, title string) {
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle(title)
}

// NewSurface creates a new off-screen [Surface] with the given dimensions.
func NewSurface(w, h int) *Surface {
	return ebiten.NewImage(w, h)
}

// Cursor shows or hides the system cursor.
func Cursor(show bool) {
	if show {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}
}

// Grab locks the cursor to the window when lock is true.
func Grab(lock bool) {
	if lock {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
}

// Expand toggles fullscreen mode.
func Expand(on bool) {
	ebiten.SetFullscreen(on)
}

// Tick is the current tick count.
var Tick = ebiten.Tick

// FPS returns the current frames per second.
func FPS() float64 {
	return ebiten.ActualFPS()
}

// TPS returns the current ticks per second.
func TPS() float64 {
	return ebiten.ActualTPS()
}
