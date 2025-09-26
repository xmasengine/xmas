package xui

type LabelClass struct {
	*Label
	*WidgetClass
}

func NewLabelClass(b *Label) *LabelClass {
	res := &LabelClass{Label: b}
	res.WidgetClass = NewWidgetClass()
	return res
}

type Label struct {
	Widget
	Text    string
	pressed bool
}

func (b LabelClass) Render(r *Root, screen *Surface) {
	box := b.Bounds
	style := b.Style

	if b.State.Hover {
		style = HoverStyle().WithTinyFont()
	}

	at := box.Min

	style.DrawBox(screen, box)
	style.DrawText(screen, at, b.Text)

	b.Widget.RenderWidgets(r, screen)
}

func (b *LabelClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *LabelClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}

func (l *Label) SetText(text string) {
	l.Text = text
}

func (p *Widget) AddLabel(bounds Rectangle, text string) *Label {
	b := NewLabel(bounds, text)
	p.Widgets = append(p.Widgets, &b.Widget)
	return b
}

func (l *Label) Init(bounds Rectangle, text string) *Label {
	l.Text = text
	l.Widget = Widget{Bounds: bounds, Style: DefaultStyle().WithTinyFont()}
	l.Class = NewLabelClass(l)
	return l
}

func NewLabel(bounds Rectangle, text string) *Label {
	res := &Label{}
	return res.Init(bounds, text)
}
