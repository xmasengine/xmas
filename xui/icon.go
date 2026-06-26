package xui

import "github.com/xmasengine/xmas/xgal"

// Icon is a small surface drawn left of text in widgets like Button,
// Label, MenuItem, etc. The zero value is safe to use.
// In that case Width returns 0 and Blit is a no-op.
type Icon struct {
	Image *xgal.Surface
	gap   int
}

const iconGap = 2

// Width returns the icon width plus the gap, or 0 if empty.
func (ic Icon) Width() int {
	if ic.Image == nil {
		return 0
	}
	return ic.Image.Bounds().Dx() + iconGap
}

// Blit draws the icon at pt on s, if the icon's Image is non-nil.
func (ic Icon) Blit(s *xgal.Surface, pt xgal.Point) {
	if ic.Image == nil {
		return
	}
	ib := ic.Image.Bounds()
	xgal.Blit(s, ic.Image, xgal.Rect(pt.X, pt.Y, pt.X+ib.Dx(), pt.Y+ib.Dy()), ib)
}

// TextBounds returns a rectangle shifted right by the icon width.
// Useful for computing the text area after the icon.
func (ic Icon) TextBounds(bounds xgal.Rectangle) xgal.Rectangle {
	if ic.Image == nil {
		return bounds
	}
	off := ic.Image.Bounds().Dx() + iconGap
	return xgal.Rect(bounds.Min.X+off, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
}
