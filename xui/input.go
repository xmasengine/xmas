package xui

import "github.com/xmasengine/xmas/xgal"

// InputKeys holds optional direction/action callbacks for widgets like
// [Ring] and [Screen].  Each field may be nil; when nil the corresponding
// field from [DefaultInput] is tried.  Set [DefaultInput] to configure
// global key bindings, then override per-widget by setting fields
// directly.
type InputKeys struct {
	Left    func() bool
	Right   func() bool
	Up      func() bool
	Down    func() bool
	Confirm func() bool
	Cancel  func() bool
	NextTab func() bool
	PrevTab func() bool
}

// DefaultInput is the global fallback for widget input callbacks.
// Modify it at init time to change bindings for all widgets.
var DefaultInput = InputKeys{
	Left:    TapAny(xgal.KeyArrowLeft, xgal.KeyA),
	Right:   TapAny(xgal.KeyArrowRight, xgal.KeyD),
	Up:      TapAny(xgal.KeyArrowUp, xgal.KeyW),
	Down:    TapAny(xgal.KeyArrowDown, xgal.KeyS),
	Confirm: TapAny(xgal.KeyEnter, xgal.KeySpace),
	Cancel:  TapAny(xgal.KeyEscape),
	NextTab: TapAny(xgal.KeyQ),
	PrevTab: TapAny(xgal.KeyE),
}

// TapAny returns a callback that returns true when any of the given key
// codes is tapped this frame.
func TapAny(keys ...xgal.KeyCode) func() bool {
	return func() bool {
		for _, k := range keys {
			if xgal.Tap(k) {
				return true
			}
		}
		return false
	}
}

// input fires cb if non-nil, otherwise falls back to default.
func input(cb, fallback func() bool) bool {
	if cb != nil {
		return cb()
	}
	if fallback != nil {
		return fallback()
	}
	return false
}
