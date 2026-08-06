package xlui

import "github.com/xmasengine/xmas/xgal"
import "fmt"

// NewStatistic returns a Control to display statistics.
// Call the Set function pointer in the Class to update the stats.
func NewStatistic(at xgal.Point, form string, args ...any) *Control {
	statistic := NewControlWithText(at, ButtonStyle(), fmt.Sprintf(form, args...))
	render := func(screen *xgal.Surface) {
		statistic.Style.Print(screen, statistic.Bounds.Min, statistic.Text)
	}
	set := func(args ...any) error {
		statistic.Text = fmt.Sprintf(form, args...)
		return nil
	}
	statistic.Class = Class{
		Render: render,
		Set:    set,
	}
	return statistic
}

// Statistic adds a new statistic display to the layer.
func (l *Layer) Statistic(form string, args ...any) *Control {
	at := l.Bounds.Min
	ctrl := NewStatistic(at, form, args...)
	return l.Append(ctrl)
}

// NewBar adds a new bar Control for bar display of statistics.
// Set bar.Value and Bar.High to update this.
func NewBar(at xgal.Point, orientation Orientation, value, high int) *Control {
	bar := NewControl(at)
	bar.Style = HUDBarStyle()
	bar.State.Clicked = false
	bar.High = high
	bar.Value = value
	bar.Orientation = orientation

	size := xgal.Pt(bar.High, knobSize)
	if orientation == Vertical {
		size = xgal.Pt(knobSize, bar.High)
	}
	bar.Bounds = xgal.Bound(bar.Bounds.Min.X, bar.Bounds.Min.Y, size.X, size.Y)

	render := func(screen *xgal.Surface) {
		bl := size.Mul(bar.Value).Div(bar.High)
		filling := xgal.Bound(bar.Bounds.Min.X, bar.Bounds.Min.Y, bl.X, knobSize)
		if orientation == Vertical {
			filling = xgal.Bound(bar.Bounds.Min.X, bar.Bounds.Min.Y, knobSize, bl.Y)
		}

		style := bar.Style
		if bar.State.Hovered {
			style = style.Hovered()
		}
		if bar.State.Focused {
			style = style.Focused()
		}

		style.DrawPlainBox(screen, filling)
		style.DrawRect(screen, bar.Bounds)
	}

	bar.Class.Render = render
	return bar
}

// Bar adds a new status bar to the layer.
func (l *Layer) Bar(orientation Orientation, value, high int) *Control {
	at := l.Bounds.Min
	ctrl := NewBar(at, orientation, value, high)
	return l.Append(ctrl)
}
