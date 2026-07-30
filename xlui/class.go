package xlui

import "github.com/xmasengine/xmas/xgal"

// Class of a Layer or Control contains custom handlers for the
// a particular kind of layer or control.
// If the handler is nil, this meand the type doesn4t support the operation
// and a default handler may be used.
type Class struct {
	Render  func(screen *xgal.Surface)
	Hover   func(at xgal.Point) Reply
	Click   func(at xgal.Point, button int) Reply
	Release func(at xgal.Point, button int) Reply
	Key     func(key int, duration int) Reply
	Tap     func(key int, mods Mods) Reply
	Lift    func(key int, mods Mods) Reply
	Entry   func(string) Reply
	Chars   func(chars ...rune) Reply
}

type ClickFunc func(at xgal.Point, button int) Reply

func (f ClickFunc) Link(linked ClickFunc) ClickFunc {
	return func(at xgal.Point, button int) Reply {
		linked(at, button)
		return f(at, button)
	}
}

func (c *Class) LinkClick(linked ClickFunc) {
	if c.Click == nil {
		c.Click = linked
	} else {
		cf := ClickFunc(c.Click)
		c.Click = cf.Link(linked)
	}
}
