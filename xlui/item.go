package xlui

import "github.com/xmasengine/xmas/xgal"

// NewItem returns a Control that can serve as an item in an item or
// configuration menu.
func NewItem(at xgal.Point, icon *xgal.Image, text string, value int, options ...string) *Control {
	item := NewControlWithText(at, DefaultStyle(), text)
	item.Value = value
	return item
}
