package xui

import "github.com/xmasengine/xmas/xgal"

// ChooserLayer is a tile-selection widget. It displays an image and lets the
// user pick a tile from it. Hovered and selected tiles are highlighted.
type ChooserLayer struct {
	Layer
	Image    *xgal.Surface
	TileSize xgal.Point
	Hovered  xgal.Point
	Selected xgal.Point
	Select   func(x, y int)
}

// Chooser returns a new [ChooserLayer].
func Chooser(bounds xgal.Rectangle, img *xgal.Surface, tileSize xgal.Point, selectFn func(x, y int)) *ChooserLayer {
	c := &ChooserLayer{
		Image:    img,
		TileSize: tileSize,
		Select:   selectFn,
	}
	c.Layer = MakeLayer(bounds)
	return c
}

var _ Widget = &ChooserLayer{}

func (c *ChooserLayer) Poll() Reply {
	pos := xgal.Mouse()
	if !pos.In(c.Bounds) {
		return Ignore
	}

	mx := (pos.X - c.Bounds.Min.X) / c.TileSize.X
	my := (pos.Y - c.Bounds.Min.Y) / c.TileSize.Y
	if mx < 0 {
		mx = 0
	}
	if my < 0 {
		my = 0
	}
	if c.Image != nil {
		w := c.Image.Bounds().Dx() / c.TileSize.X
		h := c.Image.Bounds().Dy() / c.TileSize.Y
		if mx >= w {
			mx = w - 1
		}
		if my >= h {
			my = h - 1
		}
	}

	c.Hovered = xgal.Pt(mx, my)

	if xgal.Click(xgal.MouseButtonLeft) || xgal.Click(xgal.MouseButtonRight) {
		c.Selected = c.Hovered
		if c.Select != nil {
			c.Select(mx, my)
		}
		return Accept
	}

	return Accept
}

func (c *ChooserLayer) Render(s *xgal.Surface) {
	c.Layer.Render(s)

	if c.Image != nil {
		xgal.Blit(s, c.Image, c.Bounds, c.Image.Bounds())
	}

	hb := xgal.Rect(
		c.Bounds.Min.X+c.Hovered.X*c.TileSize.X,
		c.Bounds.Min.Y+c.Hovered.Y*c.TileSize.Y,
		c.Bounds.Min.X+(c.Hovered.X+1)*c.TileSize.X,
		c.Bounds.Min.Y+(c.Hovered.Y+1)*c.TileSize.Y,
	)
	xgal.Outline(s, hb, 2, xgal.Wash(255, 255, 255, 255))

	sb := xgal.Rect(
		c.Bounds.Min.X+c.Selected.X*c.TileSize.X,
		c.Bounds.Min.Y+c.Selected.Y*c.TileSize.Y,
		c.Bounds.Min.X+(c.Selected.X+1)*c.TileSize.X,
		c.Bounds.Min.Y+(c.Selected.Y+1)*c.TileSize.Y,
	)
	xgal.Outline(s, sb, 2, xgal.Wash(255, 255, 0, 255))
}

func (c *ChooserLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	nw, nh := bounds.Dx(), bounds.Dy()
	if c.Image != nil {
		nw = c.Image.Bounds().Dx()
		nh = c.Image.Bounds().Dy()
	}
	c.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return c.Bounds
}

func (m *Layer) AddChooser(bounds xgal.Rectangle, img *xgal.Surface, tileSize xgal.Point, selectFn func(x, y int)) *ChooserLayer {
	c := Chooser(bounds, img, tileSize, selectFn)
	m.Add(c)
	return c
}
