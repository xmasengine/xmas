package xui

import "image"

// SliderClass has to common methods for a slider.
type SliderClass struct {
	*Slider
	*BoxClass
}

func NewSliderClass(s *Slider) *SliderClass {
	sc := &SliderClass{Slider: s}
	sc.BoxClass = NewBoxClass(&s.Box)
	return sc
}

// HorizontalSliderClass embeds SliderClass and implements a horizontal slider.
type HorizontalSliderClass struct {
	*Slider
	*SliderClass
}

func NewHorizontalSliderClass(s *Slider) *HorizontalSliderClass {
	sc := &HorizontalSliderClass{Slider: s}
	sc.SliderClass = NewSliderClass(s)
	return sc
}

// VerticalSliderClass embeds SliderClass and implements a vertical slider.
type VerticalSliderClass struct {
	*Slider
	*SliderClass
}

func NewVerticalSliderClass(s *Slider) *VerticalSliderClass {
	sc := &VerticalSliderClass{Slider: s}
	sc.SliderClass = NewSliderClass(s)
	return sc
}

// KnobClass has to common methods for a Knob.
type KnobClass struct {
	*Knob
	*WidgetClass
}

func NewKnobClass(k *Knob) *KnobClass {
	kc := &KnobClass{Knob: k}
	kc.WidgetClass = NewWidgetClass()
	return kc
}

func (s *KnobClass) Render(r *Root, screen *Surface) {
	s.Knob.Style.KnobStyle().DrawCircle(screen, s.Knob.Bounds.Min, s.Knob.Radius)
}

// Knob is the knob of the slider.
type Knob struct {
	Widget
	Radius   int
	OnScroll func(*Knob)
}

// Init initializes a Knob.
func (k *Knob) Init(bounds Rectangle, cb func(*Knob)) *Knob {
	k.Bounds = bounds
	k.OnScroll = cb
	k.Radius = min(bounds.Dx()/2, bounds.Dy()/2)
	k.Class = NewKnobClass(k)
	return k
}

// NewKnob return a Knob for a slider.
func NewKnob(bounds Rectangle, cb func(*Knob)) *Knob {
	s := &Knob{}
	return s.Init(bounds, cb)
}

// AddKnob adds a know to this widget.
func (w *Widget) AddKnob(bounds Rectangle, cb func(*Knob)) *Knob {
	knob := NewKnob(bounds, cb)
	w.Widgets = append(w.Widgets, &knob.Widget)
	return knob
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
	Knob   *Knob //
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
	s.Scrolled = scrolled
	s.Scroll = cb
	knob := s.AddKnob(bounds, nil)
	s.Knob = knob
	s.SliderSpecial.SliderUpdate()
	return s
}

// NewSlider return a Slider that will also scroll scrolled if not nil.
func NewSlider(bounds Rectangle, scrolled *Widget, cb func(*Slider)) *Slider {
	s := &Slider{}
	return s.Init(bounds, scrolled, cb)
}

func (s *HorizontalSliderClass) SliderUpdate() {
	delta := image.Point{}

	s.Knob.Bounds.Min.Y = s.Bounds.Min.Y + s.Bounds.Dy()/2
	s.Knob.Bounds.Min.X = s.Bounds.Min.X + s.Style.Margin.X
	dx := ((s.Pos - s.Low) * (s.Bounds.Dx() - 2*s.Style.Margin.X)) / s.High
	s.Knob.Bounds.Min.X += dx
	delta.X = dx

	if s.Scrolled != nil {
		s.Scrolled.ScrollHorizontal(s.Pos, s.Low, s.High)
	}
	if s.Scroll != nil {
		s.Scroll(s.Slider)
	}
}

func (s *VerticalSliderClass) SliderUpdate() {
	delta := image.Point{}
	s.Knob.Bounds.Min.X = s.Bounds.Min.X + s.Bounds.Dx()/2
	s.Knob.Bounds.Min.Y = s.Bounds.Min.Y + s.Style.Margin.Y
	dy := ((s.Pos - s.Low) * (s.Bounds.Dy() - 2*s.Style.Margin.Y)) / s.High
	s.Knob.Bounds.Min.Y += dy
	delta.Y = dy

	if s.Scrolled != nil {
		s.Scrolled.ScrollVertical(s.Pos, s.Low, s.High)
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
	s.Slider.Knob.Class.Render(r, screen)
}

// AddSlider adds a Slider as a Control of this widget.
// This will set up automatic scrolling for the panel as well.
func (w *Widget) AddSlider(bounds Rectangle, scrolled *Widget, cb func(*Slider)) *Slider {
	slider := NewSlider(bounds, scrolled, cb)
	w.Widgets = append(w.Widgets, &slider.Widget)
	return slider
}

var DefaultScrollerSize Point = image.Pt(8, 8)

// AddVerticalScroller adds a Slider as a Control of this widget.
// It will be locked on the right.
// This will set up automatic scrolling for the widget as well.
func (w *Widget) AddVerticalScroller(cb func(*Slider)) *Slider {
	bounds := w.Bounds
	bounds.Min.X = bounds.Max.X - DefaultScrollerSize.X
	s := w.AddSlider(bounds, w, cb)
	s.State.Lock = true
	return s
}

// AddHorizontalScroller adds a Slider as a Control of this widget.
// It will be locked on the bottom.
// This will set up automatic scrolling for the widget as well.
func (w *Widget) AddHorizontalScroller(cb func(*Slider)) *Slider {
	bounds := w.Bounds
	bounds.Min.Y = bounds.Max.Y - DefaultScrollerSize.Y
	s := w.AddSlider(bounds, w, cb)
	s.State.Lock = true
	return s
}
