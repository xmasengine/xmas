package xui

import "image"

type SliderClass struct {
	*Slider
	*BoxClass
}

func NewSliderClass(s *Slider) *SliderClass {
	sc := &SliderClass{Slider: s}
	sc.BoxClass = NewBoxClass(&s.Box)
	return sc
}

type HorizontalSliderClass struct {
	*Slider
	*SliderClass
}

func NewHorizontalSliderClass(s *Slider) *HorizontalSliderClass {
	sc := &HorizontalSliderClass{Slider: s}
	sc.SliderClass = NewSliderClass(s)
	return sc
}

type VerticalSliderClass struct {
	*Slider
	*SliderClass
}

func NewVerticalSliderClass(s *Slider) *VerticalSliderClass {
	sc := &VerticalSliderClass{Slider: s}
	sc.SliderClass = NewSliderClass(s)
	return sc
}

// SliderSpecial is an additional interface of special methods that
// that a SliderClass has to implement.
type SliderSpecial interface {
	SliderUpdate() // SliderUpdate is an update callback.
}

// Slider is a Slider gadget, vertical by default.
type Slider struct {
	Box
	Pos    int
	Low    int
	High   int
	Knob   image.Point
	Radius int
	Scroll func(*Slider)

	Scrolled      *Widget // Widget that will be scrolled if not nil.
	SliderSpecial         // Update Callback
}

func (s *Slider) SetValue(pos int, lowHigh ...int) {
	if len(lowHigh) > 0 {
		s.Low = lowHigh[0]
	}
	if len(lowHigh) > 1 {
		s.High = lowHigh[1]
	}
	if pos >= s.Low && pos <= s.High {
		s.Pos = pos
	}
	s.SliderSpecial.SliderUpdate()
}

// Init initializes a slider.
func (s *Slider) Init(bounds Rectangle, scrolled *Widget, cb func(*Slider)) *Slider {
	s.Box.Init(bounds)
	if horizontal := s.Bounds.Dx() > s.Bounds.Dy(); horizontal {
		cl := NewHorizontalSliderClass(s)
		s.Class = cl
		s.SliderSpecial = cl
	} else {
		cl := NewVerticalSliderClass(s)
		s.Class = cl
		s.SliderSpecial = cl
	}

	s.Low = 0
	s.Pos = 0
	s.High = 100
	s.Radius = min(s.Bounds.Dx(), s.Bounds.Dy()) / 2
	s.Knob = bounds.Min
	s.Scrolled = scrolled
	s.Scroll = cb
	s.SliderSpecial.SliderUpdate()
	return s
}

// NewSilder return a Slider that will also scroll scrolled if not nil.
func NewSlider(bounds Rectangle, scrolled *Widget, cb func(*Slider)) *Slider {
	s := &Slider{}
	return s.Init(bounds, scrolled, cb)
}

func (s *HorizontalSliderClass) SliderUpdate() {
	delta := image.Point{}

	s.Knob.Y = s.Bounds.Min.Y + s.Bounds.Dy()/2
	s.Knob.X = s.Bounds.Min.X + s.Style.Margin.X
	dx := ((s.Pos - s.Low) * (s.Bounds.Dx() - 2*s.Style.Margin.X)) / s.High
	s.Knob.X += dx
	delta.X = dx // XXX needs scalingfunc (s *HorizontalSliderClass) SliderUpdate() {

	if s.Scrolled != nil {
		s.Scrolled.Move(delta)
	}
	if s.Scroll != nil {
		s.Scroll(s.Slider)
	}
}

func (s *VerticalSliderClass) SliderUpdate() {
	delta := image.Point{}
	s.Knob.X = s.Bounds.Min.X + s.Bounds.Dx()/2
	s.Knob.Y = s.Bounds.Min.Y + s.Style.Margin.Y
	dy := ((s.Pos - s.Low) * (s.Bounds.Dy() - 2*s.Style.Margin.Y)) / s.High
	s.Knob.Y += dy
	delta.Y = dy // XXX needs scaling

	if s.Scrolled != nil {
		s.Scrolled.Move(delta)
	}
	if s.Scroll != nil {
		s.Scroll(s.Slider)
	}
}

func (s *HorizontalSliderClass) OnMousePress(ev MouseEvent) bool {
	dX := ev.At.X - s.Slider.Bounds.Min.X - s.Slider.Style.Margin.X
	hX := s.Slider.Bounds.Dx() - 2*s.Slider.Style.Margin.X
	s.Pos = dX * (s.High - s.Low) / hX
	s.SliderUpdate()
	return true
}

func (s *VerticalSliderClass) OnMousePress(ev MouseEvent) bool {
	dY := ev.At.Y - s.Slider.Bounds.Min.Y - s.Slider.Style.Margin.Y
	hY := s.Slider.Bounds.Dy() - 2*s.Slider.Style.Margin.Y
	s.Pos = dY * (s.High - s.Low) / hY
	s.SliderUpdate()
	return true
}

func (s *SliderClass) OnMouseRelease(ev MouseEvent) bool {
	return false
}

func (s *SliderClass) OnMouseHold(ev MouseEvent) bool {
	return s.OnMousePress(ev)
}

func (s *SliderClass) OnMouseWheel(ev MouseEvent) bool {
	s.Pos -= ev.Wheel.Y
	// clamp
	s.Pos = max(min(s.Pos, s.High), s.Low)
	s.SliderSpecial.SliderUpdate()
	return true
}

func (s *SliderClass) Render(r *Root, screen *Surface) {
	s.Slider.Style.ForState(s.Slider.State).DrawBox(screen, s.Slider.Bounds)
	s.Slider.Style.KnobStyle().DrawCircle(screen, s.Knob, s.Radius)
}

// AddSlider adds a Slider as a Control of this widget.
// This will set up automatic scrolling for the panel as well.
func (w *Widget) AddSlider(bounds Rectangle, scrolled *Widget, cb func(*Slider)) *Slider {
	slider := NewSlider(bounds, scrolled, cb)
	w.Widgets = append(w.Widgets, &slider.Widget)
	return slider
}

// AddScroller  adds a Slider as a Control of this widget.
// This will set up automatic scrolling for the widget as well.
func (w *Widget) AddScroller(bounds Rectangle, cb func(*Slider)) *Slider {
	return w.AddSlider(bounds, w, cb)
}
