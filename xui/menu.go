package xui

import "image"
import "image/color"

// import "slices"

// A bar is a horizontal bar on top of the UI or a panel with menus in it.
type Bar struct {
	Box
	Click func(*Bar)
}

type BarClass struct {
	*Bar
	*BoxClass
}

func (b *Bar) Init(bounds Rectangle, ch func(*Bar)) *Bar {
	b.Click = ch
	b.Box.Init(bounds)
	b.Widget = Widget{Bounds: bounds, Style: BarStyle()}
	b.Style.Border = color.RGBA{64, 64, 64, 191}
	b.Class = NewBarClass(b)
	return b
}

func NewBar(bounds Rectangle, ch func(*Bar)) *Bar {
	e := &Bar{}
	return e.Init(bounds, ch)
}

func NewBarClass(c *Bar) *BarClass {
	res := &BarClass{Bar: c}
	res.BoxClass = NewBoxClass(&c.Box)
	return res
}

func (w *Widget) AddBar(bounds Rectangle, ch func(*Bar)) *Bar {
	bar := NewBar(bounds, ch)
	w.Widgets = append(w.Widgets, &bar.Widget)
	return bar
}

func (b *BarClass) OnActionDrag(e ActionEvent) bool {
	b.Bar.Move(e.Delta)
	return true
}

func (b *Bar) AddItem(bounds Rectangle, heading string, cb func(*Item)) *Item {
	wrap := func(i *Item) {
		if cb != nil {
			cb(i)
		}
		if i.Menu != nil && !i.Menu.State.Hide {
			i.Menu.State.Hide = true
		}
	}
	item := NewItem(bounds, heading, wrap)
	b.Widgets = append(b.Widgets, &item.Widget)
	return item
}

func (b *Bar) FitItem(heading string, cb func(*Item)) *Item {
	bounds := b.Bounds
	size := b.Style.MeasureText(heading)
	size = size.Add(b.Style.Margin.Mul(2))
	size.Y = bounds.Dy()
	r := Rectangle{bounds.Min, bounds.Min.Add(size)}
	if len(b.Widgets) > 0 {
		last := b.Widgets[len(b.Widgets)-1]
		pos := last.Bounds.Max
		pos.Y = b.Bounds.Min.Y
		r = Rectangle{Min: pos, Max: pos.Add(size)}
	}
	return b.AddItem(r, heading, cb)
}

func (b *Bar) FitItemWithMenu(heading string, cb func(*Item)) *Item {

	item := b.FitItem(heading, cb)
	bounds := item.Bounds
	delta := image.Pt(0, item.Bounds.Dy())
	bounds = item.Bounds.Add(delta)

	wrap := func(subm *Menu) {
		if cb != nil {
			cb(item)
		}
		if !subm.State.Hide {
			subm.State.Hide = true
		}
	}

	item.Menu = NewMenu(bounds, wrap)
	item.Menu.State.Hide = true // menu starts out hidden.
	item.Widgets = append(item.Widgets, &item.Menu.Widget)
	return item
}

func (b *BarClass) Render(r *Root, screen *Surface) {
	if b.Bar.State.Hide {
		return
	}

	b.BoxClass.Render(r, screen)
	b.Bar.RenderWidgets(r, screen)
}

type MenuNotifier interface {
	MenuNotify(item *Item, closed bool)
}

// Menu is a vertical box with items in it to select or toggle.
type Menu struct {
	Box
	Result int
	Click  func(*Menu)
}

type MenuClass struct {
	*Menu
	*BoxClass
}

func (m *Menu) Init(bounds Rectangle, ch func(*Menu)) *Menu {
	m.Click = ch
	m.Box.Init(bounds)
	m.Style.Border = color.RGBA{127, 127, 127, 191}
	m.Class = NewMenuClass(m)
	return m
}

func NewMenu(bounds Rectangle, ch func(*Menu)) *Menu {
	e := &Menu{}
	return e.Init(bounds, ch)
}

func NewMenuClass(c *Menu) *MenuClass {
	res := &MenuClass{Menu: c}
	res.BoxClass = NewBoxClass(&res.Menu.Box)
	return res
}

func (m *Menu) AddItem(bounds Rectangle, heading string, cb func(*Item)) *Item {
	wrap := func(i *Item) {
		if cb != nil {
			cb(i)
		}
		if i.Menu == nil || !i.Menu.State.Hide {
			m.State.Hide = true
		}
	}
	item := NewItem(bounds, heading, wrap)
	m.Widgets = append(m.Widgets, &item.Widget)
	m.Bounds = m.Bounds.Union(item.Bounds)
	return item
}

func (m *Menu) FitItem(heading string, ch func(*Item)) *Item {
	size := m.Style.MeasureText(heading)
	r := Rectangle{m.Bounds.Min, m.Bounds.Min.Add(size)}
	if len(m.Widgets) > 0 {
		last := m.Widgets[len(m.Widgets)-1]
		pos := last.Bounds.Max
		pos.X = m.Bounds.Min.X
		r = Rectangle{Min: pos, Max: pos.Add(size)}
	}
	item := m.AddItem(r, heading, ch)
	return item
}

func (m *Menu) FitItemWithMenu(heading string, cb func(*Item)) *Item {
	item := m.FitItem(heading, cb)
	bounds := item.Bounds
	delta := image.Pt(m.Bounds.Dx(), 0)
	bounds = item.Bounds.Add(delta)

	wrap := func(subm *Menu) {
		if cb != nil {
			cb(item)
		}
		if !subm.State.Hide {
			subm.State.Hide = true
		}
	}

	item.Menu = NewMenu(bounds, wrap)
	item.Menu.State.Hide = true // menu starts out hidden.
	item.Widgets = append(item.Widgets, &item.Menu.Widget)
	return item
}

func (m *Menu) HasVisibleSubMenus() bool {
	for _, sub := range m.Widgets {
		if iclass, ok := sub.Class.(*ItemClass); ok {
			item := iclass.Item
			if item.Menu != nil && !item.Menu.State.Hide {
				return true
			}
		}
	}

	return false
}

func (m *MenuClass) Render(r *Root, screen *Surface) {
	if m.Menu.State.Hide {
		return
	}
	style := m.Menu.Style
	if m.Menu.State.Hover {
		style = HoverStyle()
	}
	style.DrawBox(screen, m.Menu.Bounds)
	m.Menu.RenderWidgets(r, screen)
}

// Item is an item to select or toggle or with a folding menu in it.
type Item struct {
	Label
	Menu   *Menu       // Optional folding sub-menu.
	Select func(*Item) // If not nil select will be called.
	Result int
}

type ItemClass struct {
	*Item
	*LabelClass
}

func (i *Item) Init(bounds Rectangle, text string, ch func(*Item)) *Item {
	i.Select = ch
	println("item.Init", bounds.Dx(), bounds.Dy(), text)
	i.Label.Init(bounds, text)
	i.Style.Border = color.RGBA{127, 127, 127, 191}
	i.Style.Fill = color.RGBA{0, 127, 250, 191}
	i.Style.Writing = color.RGBA{255, 255, 200, 255}
	i.Style = i.Style.WithTinyFont()
	i.Class = NewItemClass(i)
	return i
}

func NewItem(bounds Rectangle, text string, ch func(*Item)) *Item {
	e := &Item{}
	return e.Init(bounds, text, ch)
}

func NewItemClass(c *Item) *ItemClass {
	res := &ItemClass{Item: c}
	res.LabelClass = NewLabelClass(&c.Label)
	return res
}

func (i *ItemClass) OnMousePress(e MouseEvent) bool {
	return true
}

func (i *ItemClass) OnMouseHold(e MouseEvent) bool {
	return true
}

func (i *ItemClass) OnMouseRelease(e MouseEvent) bool {
	if i.Select != nil {
		i.Select(i.Item)
	}

	if i.Menu != nil { // If the item has a menu
		if i.Menu.State.Hide { // The sub menu it is hidden
			i.Menu.State.Hide = false // Show it
		} else { // The sub menu is visible
			i.Menu.State.Hide = true // Hide the submenu
		}
	}

	return true
}

func (i *ItemClass) Render(r *Root, screen *Surface) {
	if i.Item.State.Hide {
		return
	}

	i.LabelClass.Render(r, screen)

	if i.Menu != nil {
		i.Menu.Class.Render(r, screen)
	}
}

func (i *Item) Move(delta Point) {
	i.Label.Move(delta)
	if i.Menu != nil {
		i.Menu.Move(delta)
	}
}

// List is a vertical box with items in it to select.
type List struct {
	Box
	Select func(*Item)
	Result int
}

type ListClass struct {
	*List
	*BoxClass
}

func (l *List) Init(bounds Rectangle, text string, ch func(*Item)) *List {
	l.Select = ch
	l.Box.Init(bounds)
	l.Widget = Widget{Bounds: bounds, Style: BarStyle()}
	l.Style.Border = color.RGBA{127, 127, 127, 191}
	l.Class = NewListClass(l)
	return l
}

func NewList(bounds Rectangle, text string, ch func(*Item)) *List {
	e := &List{}
	return e.Init(bounds, text, ch)
}

func NewListClass(c *List) *ListClass {
	res := &ListClass{List: c}
	res.BoxClass = NewBoxClass(&c.Box)
	return res
}

func (m *List) AddNewItem(bounds Rectangle, heading string) *Item {
	box := m.Bounds
	delta := box.Min
	if len(m.Widgets) > 0 {
		last := m.Widgets[len(m.Widgets)-1]
		delta = last.Bounds.Max
		delta.X = box.Min.X
	}
	item := NewItem(bounds, heading, nil)
	item.Move(delta)
	m.Widgets = append(m.Widgets, &item.Widget)
	return item
}

func (p *ListClass) OnMouseWheel(ev MouseEvent) bool {
	// pos := p.Panel.verticalScrollPos(PanelScrollRange)
	// pos -= ev.Wheel.Y * PanelScrollSpeed

	// p.Panel.scrollVertical(pos, PanelScrollRange)

	return true
}

func (m *List) SelectItem(it *Item) {
	m.Result = it.Result
	if m.Select != nil {
		m.Select(it)
	}
	for _, sub := range m.Widgets {
		if iclass, ok := sub.Class.(*ItemClass); ok {
			item := iclass.Item
			if item.Result == m.Result {
				item.Style.Stroke = 1
			} else {
				item.Style.Stroke = 0
			}
		}
	}
}

func (l *ListClass) Render(r *Root, screen *Surface) {
	l.Render(r, screen)
}

const ListItemHeight = 10

func NewListWithItems(bounds Rectangle, names ...string) *List {
	list := NewList(bounds, "", nil)

	for i, name := range names {
		item := list.AddNewItem(bounds, name)
		item.Result = i
		item.Select = list.SelectItem
		bounds = bounds.Add(image.Pt(0, ListItemHeight))
	}
	// list.Clip(NodeClip(list)) // List panel will clip based on its own bounds.
	return list
}

func (w *Widget) AddNewListWithItems(bounds Rectangle, names ...string) *List {
	if w.Bounds.Dx() > w.Style.Margin.X*2 {
		bounds = bounds.Add(image.Pt(0, -w.Style.Margin.X*2))
	}
	list := NewListWithItems(bounds, names...)
	w.Widgets = append(w.Widgets, &list.Widget)
	return list
}

/*
func (w *Widget) StringListWithItems(bounds Rectangle, bind *string, names ...string) *List {
	if bind == nil {
		panic("Incorrect use of StringListWithItems")
	}

	list := NewListWithItems(bounds, names...)
	w.Widgets = append(w.Widgets, &list.Widget)

	idx := slices.Index(names, *bind)
	if idx >= 0 {
		list.SelectItem(list.Items[idx])
	}

	list.Select = func(l *Item) {
		if l.Result >= 0 && l.Result < len(names) {
			*bind = names[l.Result]
		}
	}

	if len(names)*ListItemHeight > h {
		list.AddNewScrollBar(4, h)
	}

	return list
}
*/
