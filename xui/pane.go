package xui

import "github.com/xmasengine/xmas/xgal"

const ScrollSpeed = 16

// Caption is the caption of a Pane
type Caption struct {
	Text   string
	Bounds xgal.Rectangle
	Close  xgal.Rectangle
	Style  Style
}

const CaptionHeight = 16
const CaptionCloseMargin = 2
const CaptionCloseSize = 14

func NewCaption(bounds xgal.Rectangle, text string) *Caption {
	c := &Caption{}
	c.Style = DefaultStyle()
	c.Style.Margin.Y = 0
	c.Text = text
	c.Bounds = bounds
	c.Bounds.Max.Y = bounds.Min.Y + CaptionHeight
	c.Close = c.Bounds
	c.Close.Min.X = c.Close.Max.X - CaptionCloseSize
	c.Close = c.Close.Inset(CaptionCloseMargin)
	return c
}

func (c *Caption) Render(s *xgal.Surface) {
	if c.Text != "" {
		c.Style.DrawBox(s, c.Bounds)
		c.Style.DrawBox(s, c.Close)
		c.Style.DrawX(s, c.Close.Inset(CaptionCloseMargin))
		c.Style.Ink(s, c.Bounds, c.Text)
	}
}

// Drag allows the Pane to be dragged if not locked.
func (c *Caption) Drag(p *PaneLayer) Reply {
	// When dragging the mouse may go out of the caption.
	// Nevertheless allow it to keep on dragging until released.
	if p.Drag {
		if xgal.Loose(xgal.MouseButtonLeft) {
			p.Drag = false
		} else {
			now := xgal.Cursor()
			delta := now.Sub(p.From)
			if delta.Eq(xgal.Point{}) {
				return Ignore
			}
			p.MoveBy(delta)
			p.From = now
			return Accept
		}
	}

	if !xgal.Cursor().In(c.Bounds) || p.Lock {
		p.Drag = false
		return Ignore
	}

	if xgal.Click(xgal.MouseButtonLeft) {
		p.Drag = true
		p.From = xgal.Cursor()
	}

	return Ignore
}

// Poll allows the caption to be interacted with.
func (c *Caption) Poll() Reply {
	if !xgal.Cursor().In(c.Close) {
		return Ignore
	}
	if xgal.Click(xgal.MouseButtonLeft) {
		return Finish
	}
	return Ignore
}

type PaneLayer struct {
	Layer
	Caption  *Caption
	ScrollY  int
	contentH int
}

func Pane(bounds xgal.Rectangle, heading string) *PaneLayer {
	p := &PaneLayer{}
	p.Layer = MakeLayer(bounds)
	p.Caption = NewCaption(p.Layer.Bounds, heading)
	return p
}

func (p *PaneLayer) Poll() Reply {
	res := p.Caption.Poll()
	if res != Ignore {
		return res
	}
	res = p.Caption.Drag(p)
	if res != Ignore {
		return res
	}

	_, wy := xgal.Wheel()
	if wy != 0 && xgal.Cursor().In(p.Bounds) {
		p.scroll(int(wy) * ScrollSpeed)
		return Accept
	}

	return p.PollKids()
}

func (p *PaneLayer) Render(s *xgal.Surface) {
	clipped := s.SubImage(p.Bounds).(*xgal.Surface)
	p.Layer.Render(clipped)
	if p.Caption != nil {
		p.Caption.Render(s)
	}
	p.renderPopups(s)
}

// renderPopups draws open submenus on the main surface (outside the clip)
// so they float above the pane and other siblings.
func (p *PaneLayer) renderPopups(s *xgal.Surface) {
	for _, kid := range p.Kids {
		if mb, ok := kid.(*MenuBarLayer); ok {
			for _, mkid := range mb.Kids {
				if mi, ok := mkid.(*MenuItemLayer); ok && mi.Submenu != nil && !mi.Submenu.hidden {
					mi.Submenu.Render(s)
				}
			}
		}
	}
}

var _ Widget = &PaneLayer{}

const minTotalW = 40
const minTotalH = 40

func (p *PaneLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	// Preserve the pane's current position (may have been set by dragging).
	pos := p.Bounds.Min
	cy := pos.Y
	if p.Caption != nil {
		cy = pos.Y + CaptionHeight
	}
	cb := xgal.Rect(pos.X, cy, bounds.Max.X, bounds.Max.Y)
	r := p.Layer.Place(cb)
	kw, kh := r.Dx(), r.Dy()

	totalW := bounds.Dx()
	if kw > totalW {
		totalW = kw
	}
	if totalW < minTotalW {
		totalW = minTotalW
	}
	totalH := (cy - pos.Y) + kh
	if totalH < minTotalH {
		totalH = minTotalH
	}

	if p.Caption != nil {
		p.Caption.Bounds = xgal.Rect(pos.X, pos.Y, pos.X+totalW, pos.Y+CaptionHeight)
		p.Caption.Close = p.Caption.Bounds
		p.Caption.Close.Min.X = p.Caption.Close.Max.X - CaptionCloseSize
		p.Caption.Close = p.Caption.Close.Inset(CaptionCloseMargin)
	}

	p.Bounds = xgal.Rect(pos.X, pos.Y, pos.X+totalW, pos.Y+totalH)
	p.contentH = kh
	// Re-apply accumulated scroll offset (Place resets children to natural
	// positions, so we need to move them again).
	if p.ScrollY != 0 {
		for _, kid := range p.Kids {
			if mv, ok := kid.(Mover); ok {
				mv.MoveBy(xgal.Pt(0, -p.ScrollY))
			}
		}
	}
	return p.Bounds
}

func (p *PaneLayer) MoveBy(delta xgal.Point) {
	p.Layer.MoveBy(delta)
	if p.Caption != nil {
		p.Caption.Bounds = p.Caption.Bounds.Add(delta)
		p.Caption.Close = p.Caption.Close.Add(delta)
	}
}

func (p *PaneLayer) scroll(dy int) {
	paneH := p.Bounds.Dy()
	if p.Caption != nil {
		paneH -= p.Caption.Bounds.Dy()
	}
	maxScroll := p.contentH - paneH
	if maxScroll < 0 {
		maxScroll = 0
	}
	newScroll := clamp(p.ScrollY+dy, 0, maxScroll)
	delta := p.ScrollY - newScroll
	if delta != 0 {
		p.ScrollY = newScroll
		for _, kid := range p.Kids {
			if mv, ok := kid.(Mover); ok {
				mv.MoveBy(xgal.Pt(0, delta))
			}
		}
	}
}
