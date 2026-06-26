package xui

import (
	"github.com/xmasengine/xmas/xgal"
)

// MenuBarLayer is a horizontal menu bar, typically used inside a [PaneLayer].
type MenuBarLayer struct {
	Layer
}

// MenuBar creates a horizontal menu bar with the given bounds.
func MenuBar(bounds xgal.Rectangle) *MenuBarLayer {
	m := &MenuBarLayer{}
	m.Layer = MakeLayer(bounds)
	m.Layer.Axis = Horizontal
	m.Layer.Bounds = bounds
	return m
}

var _ Widget = &MenuBarLayer{}

// AddItem adds a top-level menu item to the bar with explicit bounds.
func (m *MenuBarLayer) AddItem(text string, click func()) *MenuItemLayer {
	bounds := m.Bounds
	if len(m.Kids) > 0 {
		last := m.Kids[len(m.Kids)-1].(*MenuItemLayer)
		bounds = xgal.Rect(last.Bounds.Max.X, bounds.Min.Y,
			last.Bounds.Max.X, bounds.Max.Y)
	}
	item := MenuItem(bounds, text, click)
	m.Add(item)
	return item
}

// FitItem adds a menu item sized to fit the text.
func (m *MenuBarLayer) FitItem(text string, click func()) *MenuItemLayer {
	sz := m.Style.MeasureText(text)
	x := m.Bounds.Min.X
	if len(m.Kids) > 0 {
		last := m.Kids[len(m.Kids)-1].(*MenuItemLayer)
		x = last.Bounds.Max.X
	}
	bounds := xgal.Rect(x, m.Bounds.Min.Y, x+sz.X+m.Style.Margin.X*2, m.Bounds.Max.Y)
	item := MenuItem(bounds, text, click)
	m.Add(item)
	return item
}

// MenuItemLayer is a clickable item in a [MenuBarLayer] or [MenuLayer],
// optionally with a dropdown submenu.
type MenuItemLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Text    string
	Icon    Icon // optional icon, drawn left of text
	Click   func()
	Submenu *MenuLayer
	leaf    bool // true when this item is in a dropdown and should return Finish on click
	hover   bool
}

// MenuItem returns a new [MenuItemLayer].
func MenuItem(bounds xgal.Rectangle, text string, click func()) *MenuItemLayer {
	return &MenuItemLayer{
		Bounds: bounds,
		Style:  DefaultStyle(),
		Text:   text,
		Click:  click,
	}
}

var _ Widget = &MenuItemLayer{}

func (i *MenuItemLayer) Poll() Reply {
	i.hover = xgal.Mouse().In(i.Bounds)

	if i.Submenu != nil && !i.Submenu.hidden {
		i.Submenu.Bounds.Min = i.Bounds.Min.Add(xgal.Pt(0, i.Bounds.Dy()))
		i.Submenu.Bounds.Max = i.Submenu.Bounds.Min.Add(
			xgal.Pt(i.Bounds.Dx()*3, i.Submenu.ItemHeight()*len(i.Submenu.Kids)))

		res := i.Submenu.Poll()
		if res != Ignore {
			if res == Finish {
				i.Submenu.hidden = true
				return Accept // don't propagate Finish — would delete this item
			}
			return res
		}
		i.hover = true
		return Accept
	}

	if i.hover && xgal.Click(xgal.MouseButtonLeft) {
		if i.Click != nil {
			i.Click()
		}
		if i.Submenu != nil {
			i.Submenu.hidden = !i.Submenu.hidden
			return Accept
		}
		return i.reply()
	}

	if i.hover {
		return Accept
	}

	return Ignore
}

func (i *MenuItemLayer) reply() Reply {
	if i.leaf {
		return Finish
	}
	return Accept
}

func (i *MenuItemLayer) Render(s *xgal.Surface) {
	style := i.Style
	if i.hover {
		style = style.HoverStyle()
	}
	style.DrawBox(s, i.Bounds)

	i.Icon.Blit(s, i.Bounds.Min)
	style.Ink(s, i.Icon.TextBounds(i.Bounds), i.Text)

	if i.Submenu != nil && !i.Submenu.hidden {
		i.Submenu.Render(s)
	}
}

func (i *MenuItemLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	sz := i.Style.MeasureText(i.Text)
	nw := sz.X + i.Style.Margin.X*2 + i.Icon.Width()
	nh := sz.Y + i.Style.Margin.Y*2
	i.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return i.Bounds
}

func (i *MenuItemLayer) MoveBy(delta xgal.Point) {
	i.Bounds = i.Bounds.Add(delta)
	if i.Submenu != nil {
		i.Submenu.MoveBy(delta)
	}
}

// AddDropdown lazily creates a submenu and returns it.
func (i *MenuItemLayer) AddDropdown() *MenuLayer {
	if i.Submenu == nil {
		i.Submenu = Menu(xgal.Rect(
			i.Bounds.Min.X, i.Bounds.Min.Y+i.Bounds.Dy(),
			i.Bounds.Min.X+i.Bounds.Dx()*3, i.Bounds.Min.Y+i.Bounds.Dy()))
	}
	return i.Submenu
}

// MenuLayer is a vertical dropdown menu.
type MenuLayer struct {
	Bounds xgal.Rectangle
	Style  Style
	Kids   []Widget
	hidden bool
	itemH  int // cached item height
}

// Menu returns a new [MenuLayer] at the given position (hidden by default).
func Menu(bounds xgal.Rectangle) *MenuLayer {
	return &MenuLayer{
		Bounds: bounds,
		Style:  DefaultStyle(),
		hidden: true,
	}
}

var _ Widget = &MenuLayer{}

func (m *MenuLayer) Poll() Reply {
	if m.hidden {
		return Ignore
	}
	for i := len(m.Kids) - 1; i >= 0; i-- {
		kid := m.Kids[i]
		res := kid.Poll()
		if res != Ignore {
			return res
		}
	}
	return Ignore
}

func (m *MenuLayer) Render(s *xgal.Surface) {
	if m.hidden {
		return
	}
	m.Style.DrawBox(s, m.Bounds)
	for i := len(m.Kids) - 1; i >= 0; i-- {
		m.Kids[i].Render(s)
	}
}

func (m *MenuLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	cy := bounds.Min.Y
	maxW := 0
	for _, kid := range m.Kids {
		r := kid.Place(xgal.Rect(bounds.Min.X, cy, bounds.Max.X, bounds.Max.Y))
		kw, kh := r.Dx(), r.Dy()
		cy += kh
		if kw > maxW {
			maxW = kw
		}
	}
	totalH := cy - bounds.Min.Y
	m.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+maxW, bounds.Min.Y+totalH)
	return m.Bounds
}

func (m *MenuLayer) MoveBy(delta xgal.Point) {
	m.Bounds = m.Bounds.Add(delta)
	for _, kid := range m.Kids {
		if mv, ok := kid.(Mover); ok {
			mv.MoveBy(delta)
		}
	}
}

func (m *MenuLayer) ItemHeight() int {
	if m.itemH > 0 {
		return m.itemH
	}
	for _, kid := range m.Kids {
		if mi, ok := kid.(*MenuItemLayer); ok {
			h := mi.Bounds.Dy()
			if h > m.itemH {
				m.itemH = h
			}
		}
	}
	if m.itemH == 0 {
		m.itemH = CaptionHeight
	}
	return m.itemH
}

// AddItem adds a dropdown item to this menu.
func (m *MenuLayer) AddItem(text string, click func()) *MenuItemLayer {
	bounds := m.Bounds
	if len(m.Kids) > 0 {
		last := m.Kids[len(m.Kids)-1].(*MenuItemLayer)
		bounds = xgal.Rect(bounds.Min.X, last.Bounds.Max.Y,
			bounds.Max.X, last.Bounds.Max.Y+bounds.Dy())
	}
	item := MenuItem(bounds, text, click)
	item.Style = DefaultStyle()
	item.leaf = true
	m.Kids = append(m.Kids, item)
	return item
}

// FitItem adds a dropdown item sized to fit the text.
func (m *MenuLayer) FitItem(text string, click func()) *MenuItemLayer {
	sz := m.Style.MeasureText(text)
	y := m.Bounds.Min.Y
	if len(m.Kids) > 0 {
		last := m.Kids[len(m.Kids)-1].(*MenuItemLayer)
		y = last.Bounds.Max.Y
	}
	bounds := xgal.Rect(m.Bounds.Min.X, y, m.Bounds.Max.X, y+sz.Y+m.Style.Margin.Y*2)
	item := MenuItem(bounds, text, click)
	item.Style = DefaultStyle()
	item.leaf = true
	m.Kids = append(m.Kids, item)
	return item
}

// PaneLayer menu helpers

// AddMenuBar adds a [MenuBarLayer] at the top of the pane, below the caption.
func (p *PaneLayer) AddMenuBar() *MenuBarLayer {
	y := p.Bounds.Min.Y
	if p.Caption != nil {
		y = p.Caption.Bounds.Max.Y
	}
	bounds := xgal.Rect(p.Bounds.Min.X, y, p.Bounds.Max.X, y+CaptionHeight)
	mb := MenuBar(bounds)
	p.Add(mb)
	return mb
}
