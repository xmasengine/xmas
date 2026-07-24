package xlui

import "github.com/xmasengine/xmas/xgal"

// Class of a Layer or Control contains custom handlers for the
// a particular kind of layer or control.
// If the handler is nil, this meand the type doesn4t support the operation
// and a default handler may be used.
type Class struct {
	Render func(screen *xgal.Surface)
	Click  func(at xgal.Point, button int) Reply
	Key    func(key int, mod int) Reply
}
