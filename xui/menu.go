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
	b.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
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

func (b *BarClass) OnActionDrag(e ActionEvent) bool {
	b.Bar.Move(e.Delta)
	return true
}

func (b *Bar) AddNewItem(bounds Rectangle, heading string) *Item {
	size := image.Pt(bounds.Dx(), b.Style.LineHeight())
	r := Rectangle{bounds.Min, bounds.Min.Add(size)}
	if len(b.Widgets) > 0 {
		last := b.Widgets[len(b.Widgets)-1]
		pos := last.Bounds.Max
		pos.Y = b.Bounds.Min.Y
		r = Rectangle{Min: pos, Max: pos.Add(size)}
	}
	item := NewItem(bounds, heading, nil)
	// XXX should use Move()
	item.Bounds = r
	b.Widgets = append(b.Widgets, &item.Widget)
	return item
}

/*
func (b *Bar) AddNewItemWithMenu(bounds Rectangle, heading string) *Item {
	item := b.AddNewItem(bounds, heading)
	item.Menu = NewMenu(bounds, "", item)
	delta := item.Min.Add(image.Pt(0, b.Bounds().Dy()))
	item.Menu.Move(delta)
	return item
}
*/

func (b *BarClass) Render(r *Root, screen *Surface) {
	if b.Bar.State.Hide {
		return
	}

	b.BoxClass.Render(r, screen)

	// RenderChildren(screen, b.Widgets...)
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

func NewMenu(bounds Rectangle, text string, ch func(*Menu)) *Menu {
	e := &Menu{}
	return e.Init(bounds, ch)
}

func NewMenuClass(c *Menu) *MenuClass {
	res := &MenuClass{Menu: c}
	res.BoxClass = NewBoxClass(&res.Menu.Box)
	return res
}

func (m *Menu) AddNewItem(bounds Rectangle, heading string, ch func(*Item)) *Item {
	size := image.Pt(bounds.Dx(), m.Style.LineHeight())
	r := Rectangle{m.Bounds.Min, m.Bounds.Min.Add(size)}
	if len(m.Widgets) > 0 {
		last := m.Widgets[len(m.Widgets)-1]
		pos := last.Bounds.Max
		pos.X = m.Bounds.Min.X
		r = Rectangle{Min: pos, Max: pos.Add(size)}
	}
	item := NewItem(bounds, heading, ch)
	item.Bounds = r
	m.Widgets = append(m.Widgets, &item.Widget)
	// XXX Should use Resize.
	m.Bounds.Max.Add(image.Pt(0, m.Style.LineHeight()))

	return item
}

func (m *Menu) AddNewItemWithMenu(bounds Rectangle, heading string, ch func(*Item), chm func(*Menu)) *Item {
	item := m.AddNewItem(bounds, heading, ch)
	item.Menu = NewMenu(bounds, heading, chm)

	delta := item.Bounds.Min.Add(image.Pt(m.Bounds.Dx(), m.Bounds.Dy()/2))
	item.Menu.Move(delta)
	return item
}

func (m *Menu) Move(delta Point) {
	m.Box.Move(delta)
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

	m.BoxClass.Render(r, screen)
	m.Menu.RenderWidgets(r, screen)
}

// Item is a  with an item to select or toggle or a menu inside it.
type Item struct {
	Label
	Menu   *Menu       // Optional sub-menu.
	Select func(*Item) // If not nil select will be called.
	Result int
}

type ItemClass struct {
	*Item
	*LabelClass
}

func (i *Item) Init(bounds Rectangle, text string, ch func(*Item)) *Item {
	i.Select = ch
	i.Label.Init(bounds, text)
	i.Style.Border = color.RGBA{127, 127, 127, 191}
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

	closed := true     // Normally close the menu on click.
	if i.Menu != nil { // If the item has a menu
		if i.Menu.State.Hide { // The sub menu it is hidden
			i.Menu.State.Hide = false // Show it
			closed = false            // We have to keep the menu open.
		} else { // The sub menu is visible
			i.Menu.State.Hide = true // Hide the submenu
			closed = true            // We have too close the whole menu.
			for _, sub := range i.Menu.Widgets {
				if iclass, ok := sub.Class.(*ItemClass); ok {
					item := iclass.Item
					if item.Menu != nil {
						item.Menu.State.Hide = true
					}
				}
			}
		}
	}
	// XXX probably not correct.
	i.Item.State.Hide = !closed

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
	l.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
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
