package xlui

import "github.com/xmasengine/xmas/xgal"
import "log/slog"

const defaultDiameter = 8

func (s *Control) slideTo(mouse xgal.Point) {
	track := s.sliderTrackSize()
	if track <= 0 {
		slog.Error("slider track size zero or negative")
		return
	}

	mousePos := mouse.X - s.Bounds.Min.X - s.Style.Margin.X
	if s.Orientation == Vertical {
		mousePos = mouse.Y - s.Bounds.Min.Y - s.Style.Margin.Y
	}
	p := mousePos * (s.High - s.Low) / track
	s.Value = min(max(p, s.Low), s.High)
	if s.Class.Value != nil {
		s.Class.Value(s.Value)
	}
}

func (s Control) knobPos() int {
	span := s.High - s.Low
	if span <= 0 {
		return s.Style.Margin.X
	}
	track := s.sliderTrackSize()
	if track <= 0 {
		slog.Error("slider track size zero or negative")
		return s.Style.Margin.X
	}
	p := ((s.Value - s.Low) * track) / span
	return p
}

func (s Control) sliderTrackSize() int {
	k := s.Style.Diameter
	if s.Orientation == Vertical {
		return s.Bounds.Dy() - k
	}
	return s.Bounds.Dx() - k
}

// NewSlider adds a new slider Control
func NewSlider(at xgal.Point, orientation Orientation, low, high, value int) *Control {
	slider := NewControl(at)
	slider.State.Clicked = false
	if high < low {
		high = low
	}
	slider.Low = low
	slider.High = high
	slider.Value = min(max(value, low), high)
	slider.Orientation = orientation

	span := high - low
	if span < 1 {
		span = 1
	}
	k := slider.Style.Diameter
	size := xgal.Pt(k+span*2, k)
	if orientation == Vertical {
		size = xgal.Pt(k, k+span*2)
	}
	slider.Bounds = xgal.Bound(slider.Bounds.Min.X, slider.Bounds.Min.Y, size.X, size.Y)

	render := func(screen *xgal.Surface) {
		kb := slider.knobPos()
		k := slider.Style.Diameter
		knob := xgal.Bound(slider.Bounds.Min.X+kb, slider.Bounds.Min.Y, k, k)
		if orientation == Vertical {
			knob = xgal.Bound(slider.Bounds.Min.X, slider.Bounds.Min.Y+kb, k, k)
		}

		style := slider.Style
		if slider.State.Hovered {
			style = style.Hovered()
		}
		if slider.State.Focused {
			style = style.Focused()
		}

		style.DrawBox(screen, slider.Bounds)
		ks := style.KnobStyle()
		ks.DrawCircleInBox(screen, knob)
	}

	click := func(at xgal.Point, which int) Reply {
		slider.slideTo(xgal.Cursor())
		return Accept
	}

	wheel := func(at xgal.Point, delta int) Reply {
		slider.Value = min(max(slider.Value+delta, slider.Low), slider.High)
		return Accept
	}

	release := func(at xgal.Point, which int) Reply {
		return Accept
	}

	hover := func(at xgal.Point) Reply {
		return Accept
	}

	slider.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Wheel:   wheel,
		Hover:   hover,
	}

	return slider
}

// Slider adds a new image chooser to the layer.
// The orientation determines if it is vertical or horizontal.
func (l *Layer) Slider(orientation Orientation, low, high, value int) *Control {
	at := l.Bounds.Min
	ctrl := NewSlider(at, orientation, low, high, value)
	return l.Append(ctrl)
}
