package xui

import "github.com/xmasengine/xmas/xgal"

const knobSize = 8

// SliderLayer is a draggable slider for selecting a value in a range.
type SliderLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Pos     int
	Low     int
	High    int
	OnSlide func(pos int)

	horizontal bool
	hover      bool
	dragging   bool
}

// Slider returns a new [SliderLayer] with the given bounds and callback.
// The slider is horizontal when wider than tall, vertical otherwise.
// Default range is 0–100.
func Slider(bounds xgal.Rectangle, onSlide func(int)) *SliderLayer {
	return &SliderLayer{
		Bounds:     bounds,
		Style:      DefaultStyle(),
		Low:        0,
		High:       100,
		OnSlide:    onSlide,
		horizontal: bounds.Dx() > bounds.Dy(),
	}
}

var _ Widget = &SliderLayer{}

func (s *SliderLayer) Poll() Reply {
	s.hover = xgal.Cursor().In(s.Bounds)

	if s.hover && xgal.Click(xgal.MouseButtonLeft) {
		s.dragging = true
		s.slideTo(xgal.Cursor())
		return Accept
	}

	if s.dragging {
		if xgal.Loose(xgal.MouseButtonLeft) {
			s.dragging = false
		} else {
			s.slideTo(xgal.Cursor())
		}
		return Accept
	}

	wx, wy := xgal.Wheel()
	if s.hover && (wx != 0 || wy != 0) {
		delta := int(wy)
		if s.horizontal {
			delta = int(wx)
		}
		s.Pos = clamp(s.Pos-delta, s.Low, s.High)
		if s.OnSlide != nil {
			s.OnSlide(s.Pos)
		}
		return Accept
	}

	if s.hover {
		return Accept
	}
	return Ignore
}

func (s *SliderLayer) slideTo(mouse xgal.Point) {
	track := s.trackSize()
	mousePos := mouse.X - s.Bounds.Min.X - s.Style.Margin.X
	if !s.horizontal {
		mousePos = mouse.Y - s.Bounds.Min.Y - s.Style.Margin.Y
	}
	p := mousePos * (s.High - s.Low) / track
	s.Pos = clamp(p, s.Low, s.High)
	if s.OnSlide != nil {
		s.OnSlide(s.Pos)
	}
}

func (s *SliderLayer) knobPos() int {
	track := s.trackSize()
	if track <= 0 {
		return s.Style.Margin.X
	}
	p := (s.Pos - s.Low) * track / (s.High - s.Low)
	return p
}

func (s *SliderLayer) trackSize() int {
	if s.horizontal {
		return s.Bounds.Dx() - 2*s.Style.Margin.X - knobSize
	}
	return s.Bounds.Dy() - 2*s.Style.Margin.Y - knobSize
}

func (s *SliderLayer) Render(dst *xgal.Surface) {
	style := s.Style
	if s.hover || s.dragging {
		style = style.HoverStyle()
	}
	if s.dragging {
		style = style.DragStyle()
	}

	style.DrawBox(dst, s.Bounds)

	track := s.trackSize()
	if track <= 0 {
		return
	}

	kp := s.knobPos()
	var kbox xgal.Rectangle
	if s.horizontal {
		kbox = xgal.Rect(
			s.Bounds.Min.X+s.Style.Margin.X+kp,
			s.Bounds.Min.Y,
			s.Bounds.Min.X+s.Style.Margin.X+kp+knobSize,
			s.Bounds.Max.Y,
		)
	} else {
		kbox = xgal.Rect(
			s.Bounds.Min.X,
			s.Bounds.Min.Y+s.Style.Margin.Y+kp,
			s.Bounds.Max.X,
			s.Bounds.Min.Y+s.Style.Margin.Y+kp+knobSize,
		)
	}

	ks := style.KnobStyle()
	ks.DrawBox(dst, kbox)
}

func (s *SliderLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	if s.horizontal {
		const maxWidth = 100
		// limit width
		nw := bounds.Dx()
		if nw > maxWidth {
			nw = maxWidth
		}
		nh := knobSize + s.Style.Margin.Y*2
		s.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
		return s.Bounds
	}
	// vertical
	const minHeight = 40
	nw := knobSize + s.Style.Margin.X*2
	nh := bounds.Dy() - s.Style.Margin.Y*2
	if nh < minHeight {
		nh = minHeight
	}
	s.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return s.Bounds
}

func (s *SliderLayer) MoveBy(delta xgal.Point) {
	s.Bounds = s.Bounds.Add(delta)
}

func (m *Layer) AddSlider(bounds xgal.Rectangle, onSlide func(int)) *SliderLayer {
	s := Slider(bounds, onSlide)
	m.Add(s)
	return s
}

func clamp(v, low, high int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}
