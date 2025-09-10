package xui

type Button struct {
	Widget
	Text    string
	Clicked func(*Button)
	pressed bool
	Result  int // May be set freely except on dialog Buttons.
}

type ButtonClass struct {
	*Button
	*WidgetClass
}

func NewButtonClass(b *Button) *ButtonClass {
	res := &ButtonClass{Button: b}
	res.WidgetClass = NewWidgetClass()
	return res
}

func (b ButtonClass) Render(r *Root, screen *Surface) {
	box := b.Bounds
	style := b.Style

	if b.pressed {
		box = box.Add(b.Style.Margin)
		style = PressStyle()
	} else if b.State.Hover {
		style = HoverStyle()
	}

	at := box.Min

	style.DrawBox(screen, box)
	style.DrawText(screen, at, b.Text)
}

func (b *ButtonClass) OnActionHover(e ActionEvent) bool {
	b.State.Hover = true
	return true
}

func (b *ButtonClass) OnActionCrash(e ActionEvent) bool {
	b.State.Hover = false
	return true
}

func (b *ButtonClass) OnMousePress(e MouseEvent) bool {
	b.pressed = true
	return true
}

func (b *ButtonClass) OnMouseRelease(e MouseEvent) bool {
	b.pressed = false
	if b.Clicked != nil {
		b.Clicked(b.Button)
	}
	return true
}

func (b *Button) SetText(text string) {
	b.Text = text
}

func (b *Button) Init(bounds Rectangle, text string, cl func(*Button)) *Button {
	b.Text = text
	b.Clicked = cl
	b.Widget = Widget{Bounds: bounds, Style: DefaultStyle()}
	b.Class = NewButtonClass(b)
	return b
}

func NewButton(bounds Rectangle, text string, cl func(*Button)) *Button {
	ini := &Button{}
	return ini.Init(bounds, text, cl)
}

func (p *Widget) AddButton(bounds Rectangle, text string, cl func(*Button)) *Button {
	b := NewButton(bounds, text, cl)
	p.Widgets = append(p.Widgets, &b.Widget)
	return b
}
