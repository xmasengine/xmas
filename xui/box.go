package xui

type Box struct {
	Widget
}

func (w *Widget) AddBox(bounds Rectangle) *Box {
	box := NewBox(bounds)
	w.Widgets = append(w.Widgets, &box.Widget)
	return box
}

func (b *Box) Init(bounds Rectangle) *Box {
	b.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
	b.Class = NewBoxClass(b)
	return b
}

func NewBox(bounds Rectangle) *Box {
	box := &Box{}
	return box.Init(bounds)
}

type BoxClass struct {
	*Box
	*WidgetClass
}

func NewBoxClass(b *Box) *BoxClass {
	res := &BoxClass{Box: b}
	res.WidgetClass = NewWidgetClass()
	return res
}

// Render is called when the element needs to be drawn.
func (bc BoxClass) Render(r *Root, screen *Surface) {
	b := bc.Box
	style := b.Style
	if b.State.Hover {
		style = HoverStyle()
	}
	style.DrawBox(screen, b.Bounds)
	for _, w := range b.Widgets {
		if !w.State.Hide {
			w.Class.Render(r, screen)
		}
	}
}

func (b *BoxClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *BoxClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}
