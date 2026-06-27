package xgal

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game is the interface a user implements. It has three methods:
//
//	Update() error  — called 60 times per second; return [Quit] to exit
//	Draw(screen *Surface) — called after each Update
//	Layout(w, h int) (int, int) — reports the game's logical size
type Game = ebiten.Game

// Quit is the sentinel error returned from [Game.Update] to terminate the game.
var Quit = ebiten.Termination

// Play runs the game loop. Call it with a [Game] instance created by the user.
func Play(game Game) error {
	return ebiten.RunGame(game)
}

// MonitorType holds info about a display monitor.
type MonitorType = ebiten.MonitorType

// Monitor returns the current display monitor. It has a Size() method.
func Monitor() *MonitorType {
	return ebiten.Monitor()
}

// Screen sets the window size and title.
// If w or h is 0 or negative, that dimension is set to the monitor size.
func Screen(w, h int, title string) {
	if w <= 0 || h <= 0 {
		if m := ebiten.Monitor(); m != nil {
			mw, mh := m.Size()
			if w <= 0 {
				w = mw
			}
			if h <= 0 {
				h = mh
			}
		}
	}
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle(title)
}

// NewSurface creates a new off-screen [Surface] with the given dimensions.
func NewSurface(w, h int) *Surface {
	return ebiten.NewImage(w, h)
}

// CursorShape selects a system cursor icon. Pass one to [Cursor].
type CursorShape = ebiten.CursorShapeType

const (
	// Arrow is the standard pointer arrow.
	Arrow CursorShape = ebiten.CursorShapeDefault
	// Hand is a pointing hand, useful for clickable items.
	Hand CursorShape = ebiten.CursorShapePointer
	// Crosshair is a crosshair reticle.
	Crosshair CursorShape = ebiten.CursorShapeCrosshair
	// IBeam is the vertical line cursor used for text selection.
	IBeam CursorShape = ebiten.CursorShapeText
)

// Cursor shows or hides the system cursor and optionally sets its shape.
func Cursor(show bool, shape ...CursorShape) {
	if show {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}
	if len(shape) > 0 {
		ebiten.SetCursorShape(shape[0])
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

// Pace sets the target ticks per second.
func Pace(tps int) {
	ebiten.SetTPS(tps)
}

// Stretch enables or disables window resizing by the user.
func Stretch(on bool) {
	if on {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	} else {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	}
}

// Keep controls whether the screen is kept between frames.
func Keep(on bool) {
	ebiten.SetScreenClearedEveryFrame(!on)
}

// Focus reports whether the window has focus.
func Focus() bool {
	return ebiten.IsFocused()
}

// Decorate sets whether the window has a border and optionally sets the
// window icon. The icon images should be different sizes (e.g. 16×16, 32×32).
func Decorate(bordered bool, icons ...image.Image) {
	ebiten.SetWindowDecorated(bordered)
	if len(icons) > 0 {
		ebiten.SetWindowIcon(icons)
	}
}

// Pixel sets up the window for a pixel-art game after calling [Screen].
// It enables resizing and sets the minimum window size to the game resolution.
func Pixel(w, h int) {
	ebiten.SetWindowSizeLimits(w, h, -1, -1)
	Stretch(true)
}

// Debug print a debug text at the given location.
func Debug(surface *Surface, str string, x, y int) {
	ebitenutil.DebugPrintAt(surface, str, x, y)
}
