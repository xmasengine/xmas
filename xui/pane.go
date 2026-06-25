package xui

import "github.com/xmasengine/xmas/xgal"

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
	if !xgal.Mouse().In(c.Bounds) || p.Lock {
		p.Drag = false
		return Ignore
	}

	if xgal.Click(xgal.MouseButtonLeft) {
		p.Drag = true
		p.From = xgal.Mouse()
	}
	if p.Drag {
		if xgal.Loose(xgal.MouseButtonLeft) {
			p.Drag = false
		} else {
			now := xgal.Mouse()
			delta := now.Sub(p.From)
			if delta.Eq(xgal.Point{}) {
				return Ignore
			}
			p.Bounds = p.Bounds.Add(delta)
			c.Bounds = c.Bounds.Add(delta)
			c.Close = c.Close.Add(delta)
			p.From = now
			return Accept
		}
	}
	return Ignore
}

// Pol allows the caption to be interacted with.
func (c *Caption) Poll() Reply {
	if !xgal.Mouse().In(c.Close) {
		return Ignore
	}
	if xgal.Click(xgal.MouseButtonLeft) {
		return Finish
	}
	return Ignore
}

type PaneLayer struct {
	Layer
	Caption *Caption
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
	return p.PollKids()
}

func (p *PaneLayer) Render(dst *xgal.Surface) {
	p.Layer.Render(dst)
	p.Caption.Render(dst)
}

var _ Widget = &PaneLayer{}
