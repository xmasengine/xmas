package xui

import (
	"time"

	"github.com/xmasengine/xmas/xgal"
)

// PopupLayer is a text banner that appears when entering a map popup.
// It fades in, holds, then fades out and auto-finishes.
type PopupLayer struct {
	Text   string
	Style  Style
	Bounds xgal.Rectangle

	start   time.Time
	fadeIn  time.Duration
	hold    time.Duration
	fadeOut time.Duration
}

const (
	popupFadeIn  = 500 * time.Millisecond
	popupHold    = 1500 * time.Millisecond
	popupFadeOut = 500 * time.Millisecond
)

// Popup creates a popup popup centered on the screen. Call
// [Layer.Add] to show it; it auto-finishes after the animation cycle.
func Popup(screenW, screenH int, text string) *PopupLayer {
	sz := DefaultStyle().MeasureText(text)
	nw := sz.X + 40
	nh := sz.Y + 16
	x := (screenW - nw) / 2
	y := screenH / 3
	return &PopupLayer{
		Text:    text,
		Style:   DefaultStyle(),
		Bounds:  xgal.Rect(x, y, x+nw, y+nh),
		start:   time.Now(),
		fadeIn:  popupFadeIn,
		hold:    popupHold,
		fadeOut: popupFadeOut,
	}
}

var _ Widget = &PopupLayer{}

func (z *PopupLayer) Poll() Reply {
	elapsed := time.Since(z.start)
	total := z.fadeIn + z.hold + z.fadeOut
	if elapsed >= total {
		return Finish
	}
	return Accept
}

func (z *PopupLayer) Render(s *xgal.Surface) {
	elapsed := time.Since(z.start)
	alpha := uint8(255)

	switch {
	case elapsed < z.fadeIn:
		t := float64(elapsed) / float64(z.fadeIn)
		alpha = uint8(t * 255)
	case elapsed < z.fadeIn+z.hold:
		alpha = 255
	default:
		t := float64(elapsed-z.fadeIn-z.hold) / float64(z.fadeOut)
		alpha = uint8((1 - t) * 255)
		if alpha > 255 {
			alpha = 0
		}
	}

	st := z.Style
	st.Fill.A = alpha
	st.Border.A = alpha
	st.Fore.A = alpha
	st.Shadow.A = alpha / 2

	st.DrawBox(s, z.Bounds)
	st.Ink(s, z.Bounds, z.Text)
}

func (z *PopupLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	return z.Bounds
}

func (z *PopupLayer) MoveBy(delta xgal.Point) {
	z.Bounds = z.Bounds.Add(delta)
}
