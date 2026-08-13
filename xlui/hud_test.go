package xlui

import "testing"

import "github.com/xmasengine/xmas/xgal"

// TestNewBarDegenerateHigh ensures a zero or negative High does not panic.
func TestNewBarDegenerateHigh(t *testing.T) {
	screen := xgal.Prepare(100, 20)
	for _, high := range []int{0, -5} {
		bar := NewBar(xgal.Pt(0, 0), Horizontal, 5, high)
		bar.Class.Render(screen) // used to panic on division by zero
		bar.High = 0
		bar.Class.Render(screen) // render must still be safe
	}
}

// TestNewSliderDegenerateRange ensures equal or inverted ranges do not panic
// and produce a usable control.
func TestNewSliderDegenerateRange(t *testing.T) {
	for _, c := range []struct {
		low, high int
	}{
		{10, 10},
		{10, 5},
	} {
		slider := NewSlider(xgal.Pt(0, 0), Horizontal, c.low, c.high, 7)
		if slider.High < slider.Low {
			t.Errorf("NewSlider(%d, %d): High %d < Low %d", c.low, c.high, slider.High, slider.Low)
		}
		if got := slider.knobPos(); got < 0 {
			t.Errorf("NewSlider(%d, %d): knobPos = %d, want >= 0", c.low, c.high, got)
		}
		slider.slideTo(xgal.Pt(0, 0)) // used to divide by zero
	}
}
