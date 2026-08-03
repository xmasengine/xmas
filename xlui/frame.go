package xlui

import "github.com/xmasengine/xmas/xgal"

// NewFrame returns a new Frame, that is a control thqt displays an image.
func NewFrame(at xgal.Point, img *xgal.Surface) *Control {
	frame := NewControl(at)
	if img != nil {
		size := img.Bounds()
		frame.Bounds = xgal.Bound(frame.Bounds.Min.X, frame.Bounds.Min.Y, size.Dx(), size.Dy())
	}

	frame.Class.Render = func(screen *xgal.Surface) {
		if img != nil {
			xgal.Blit(screen, img, frame.Bounds, img.Bounds())
		}
	}
	return frame
}

func (l *Layer) Frame(img *xgal.Surface) *Control {
	at := l.Bounds.Min
	frame := NewFrame(at, img)
	l.Append(frame)
	return frame
}
