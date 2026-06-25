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
