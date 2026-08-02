package xlui

import "github.com/xmasengine/xmas/xgal"

const knobSize = 8

func (s *Control) slideTo(mouse xgal.Point) {
	track := s.sliderTrackSize()
	if track <= 0 {
		println("slider track size zero or negative")
		return
	}

	mousePos := mouse.X - s.Bounds.Min.X - s.Style.Margin.X
	if s.Orientation == Vertical {
		mousePos = mouse.Y - s.Bounds.Min.Y - s.Style.Margin.Y
	}
	p := mousePos * (s.High - s.Low) / track
	s.Value = min(max(p, s.Low), s.High)
	if s.Class.Slide != nil {
		s.Class.Slide(s.Value)
	}
}

func (s Control) knobPos() int {
	track := s.sliderTrackSize()
	if track <= 0 {
		println("slider track size zero or negative")
		return s.Style.Margin.X
	}
	p := ((s.Value - s.Low) * track) / (s.High - s.Low)
	return p
}

func (s Control) sliderTrackSize() int {
	if s.Orientation == Vertical {
		return s.Bounds.Dy() - knobSize
	}
	return s.Bounds.Dx() - knobSize
}

// NewSlider adds a new slider Control
func NewSlider(at xgal.Point, orientation Orientation, low, high, value int) *Control {
	slider := NewControl(at)
	slider.State.Clicked = false
	slider.Low = low
	slider.High = high
	slider.Value = value
	slider.Orientation = orientation

	size := xgal.Pt(knobSize+(high-low)*2, knobSize)
	if orientation == Vertical {
		size = xgal.Pt(knobSize, knobSize+(high-low)*2)
	}
	slider.Bounds = xgal.Bound(slider.Bounds.Min.X, slider.Bounds.Min.Y, size.X, size.Y)

	render := func(screen *xgal.Surface) {
		kb := slider.knobPos()
		knob := xgal.Bound(slider.Bounds.Min.X+kb, slider.Bounds.Min.Y, knobSize, knobSize)
		if orientation == Vertical {
			knob = xgal.Bound(slider.Bounds.Min.X, slider.Bounds.Min.Y+kb, knobSize, knobSize)
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
