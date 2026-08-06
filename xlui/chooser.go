package xlui

import "github.com/xmasengine/xmas/xgal"

func chooserMouseTile(c *Control, image *xgal.Surface, tileSize, pos xgal.Point) xgal.Point {
	mx := (pos.X - c.Bounds.Min.X) / tileSize.X
	my := (pos.Y - c.Bounds.Min.Y) / tileSize.Y
	if mx < 0 {
		mx = 0
	}
	if my < 0 {
		my = 0
	}
	if image != nil {
		w := image.Bounds().Dx() / tileSize.X
		h := image.Bounds().Dy() / tileSize.Y
		if mx >= w {
			mx = w - 1
		}
		if my >= h {
			my = h - 1
		}
	}
	return xgal.Pt(mx, my)
}

// NewChooser adds a new image chooser Control
func NewChooser(at xgal.Point, image *xgal.Surface, tileSize xgal.Point) *Control {
	chooser := NewControl(at)
	chooser.State.Clicked = false
	size := image.Bounds()
	chooser.Bounds = xgal.Bound(chooser.Bounds.Min.X, chooser.Bounds.Min.Y, size.Dx(), size.Dy())

	var focused xgal.Point
	var hovered xgal.Point

	set := func(args ...any) error {
		if len(args) != 1 {
			return Error("Chooser: need 1 argument")
		}
		if surf, ok := args[0].(*xgal.Surface); ok {
			image = surf
			chooser.Bounds = xgal.Bound(chooser.Bounds.Min.X, chooser.Bounds.Min.Y, size.Dx(), size.Dy())
			return nil
		} else {
			return Error("Chooser: not a Surface")
		}
	}

	render := func(screen *xgal.Surface) {
		if image != nil {
			xgal.Blit(screen, image, chooser.Bounds, image.Bounds())
		}

		if chooser.State.Hovered {
			hb := xgal.Rect(
				chooser.Bounds.Min.X+hovered.X*tileSize.X,
				chooser.Bounds.Min.Y+hovered.Y*tileSize.Y,
				chooser.Bounds.Min.X+(hovered.X+1)*tileSize.X,
				chooser.Bounds.Min.Y+(hovered.Y+1)*tileSize.Y,
			)
			xgal.Outline(screen, hb, 2, chooser.Style.Hovered().Border)
		}

		if chooser.State.Focused {
			hb := xgal.Rect(
				chooser.Bounds.Min.X+focused.X*tileSize.X,
				chooser.Bounds.Min.Y+focused.Y*tileSize.Y,
				chooser.Bounds.Min.X+(focused.X+1)*tileSize.X,
				chooser.Bounds.Min.Y+(focused.Y+1)*tileSize.Y,
			)
			xgal.Outline(screen, hb, 2, chooser.Style.Focused().Border)
		}
	}

	click := func(at xgal.Point, which int) Reply {
		focused = chooserMouseTile(chooser, image, tileSize, at)
		return Accept
	}

	release := func(at xgal.Point, which int) Reply {
		return Accept
	}

	hover := func(at xgal.Point) Reply {
		hovered = chooserMouseTile(chooser, image, tileSize, at)
		return Accept
	}

	chooser.Class = Class{
		Render:  render,
		Click:   click,
		Release: release,
		Hover:   hover,
		Set:     set,
	}

	return chooser
}

// Chooser adds a new image chooser to the layer.
func (l *Layer) Chooser(image *xgal.Surface, tileSize xgal.Point) *Control {
	at := l.Bounds.Min
	ctrl := NewChooser(at, image, tileSize)
	return l.Append(ctrl)
}
