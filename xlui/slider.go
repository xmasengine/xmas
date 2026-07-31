package xlui

import "github.com/xmasengine/xmas/xgal"

const knobSize = 8

func (s *Control) slideTo(mouse xgal.Point) {
	track := s.sliderTrackSize()
	if track <= 0 {
		track = s.Style.Margin.X
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
		return s.Style.Margin.X
	}
	p := (s.Value - s.Low) * track / (s.High - s.Low)
	return p
}

func (s Control) sliderTrackSize() int {
	if s.Orientation == Horizontal {
		return s.Bounds.Dx() - 2*s.Style.Margin.X - knobSize
	}
	return s.Bounds.Dy() - 2*s.Style.Margin.Y - knobSize
}

// NewSlider adds a new slider Control
func NewSlider(at xgal.Point, orientation Orientation, low, high int) *Control {
	slider := NewControl(at)
	slider.State.Clicked = false
	size := xgal.Pt(knobSize, high-low)
	if orientation == Horizontal {
		size = xgal.Pt(knobSize+high-low, knobSize)
	}
	slider.Bounds = xgal.Bound(slider.Bounds.Min.X, slider.Bounds.Min.Y, size.X, size.Y)

	render := func(screen *xgal.Surface) {

		style := slider.Style
		if slider.State.Hovered {
			style = style.Hovered()
		}
		if slider.State.Focused {
			style = style.Focused()
		}

		style.DrawBox(screen, slider.Bounds)

		track := slider.sliderTrackSize()
		if track <= 0 {
			return
		}

		kp := slider.knobPos()
		var kbox xgal.Rectangle
		if slider.Orientation == Horizontal {
			kbox = xgal.Rect(
				slider.Bounds.Min.X+slider.Style.Margin.X+kp,
				slider.Bounds.Min.Y,
				slider.Bounds.Min.X+slider.Style.Margin.X+kp+knobSize,
				slider.Bounds.Max.Y,
			)
		} else {
			kbox = xgal.Rect(
				slider.Bounds.Min.X,
				slider.Bounds.Min.Y+slider.Style.Margin.Y+kp,
				slider.Bounds.Max.X,
				slider.Bounds.Min.Y+slider.Style.Margin.Y+kp+knobSize,
			)
		}

		ks := style.KnobStyle()
		ks.DrawBox(screen, kbox)
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
func (l *Layer) Slider(orientation Orientation, low, high int) *Control {
	at := l.Bounds.Min
	ctrl := NewSlider(at, orientation, low, high)
	return l.Append(ctrl)
}
