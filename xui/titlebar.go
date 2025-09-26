package xui

import "image"

const TitleBarSize = 8
const ResizerSize = 8

// TitleBarClass has to common methods for a TitleBar.
type TitleBarClass struct {
	*BoxClass
	*TitleBar
}

func NewTitleBarClass(s *TitleBar) *TitleBarClass {
	sc := &TitleBarClass{TitleBar: s}
	sc.BoxClass = NewBoxClass(&s.Box)
	return sc
}

// TitleBar is a bar on top of the widget that can be used to make another
// widget draggable and minimizeable.
type TitleBar struct {
	Box            // Inherit box.
	Text   string  // Text to display.
	For    *Widget // Widget that we are the header bar for.
	Button *Button // Button on the bar to hide the widget with. Set to nil to disable.
}

// OnActionDrag will be called when the drag operation begins.
func (t *TitleBarClass) OnActionDrag(e ActionEvent) bool {
	t.TitleBar.State.Drag = true
	if t.For != nil {
		t.For.MoveAll(e.Delta)
	}
	return true
}

// OnActionDrop will be called when the drag operation ends.
func (t *TitleBarClass) OnActionDrop(e ActionEvent) bool {
	t.TitleBar.State.Drag = false
	return true
}

func (t *TitleBar) Init(bounds Rectangle, header string, fw *Widget) *TitleBar {
	t.Box.Init(bounds)
	t.Text = header
	t.For = fw
	t.Style = DefaultStyle().WithTinyFont()
	t.Class = NewTitleBarClass(t)
	t.State.Lock = true

	bbox := bounds
	bbox.Max.X = bbox.Min.X + bounds.Dy()*3/4
	t.Button = t.AddButton(bbox, "X", func(b *Button) {
		dprintln("TitleBar button")
		if t.For != nil {
			dprintln("TitleBar for: ", t.For.State.Hide)
			t.For.State.Hide = !t.For.State.Hide
			t.State.Hide = false
		}
	})
	t.Button.Style = DefaultStyle().WithTinyFont()
	t.Button.State.Lock = true

	return t
}

func (t TitleBarClass) Render(r *Root, screen *Surface) {
	t.BoxClass.Render(r, screen)
	if t.Button != nil {
		t.Button.Class.Render(r, screen)
	}

	bounds := t.TitleBar.Bounds
	style := t.TitleBar.Style.ForState(t.TitleBar.State)

	at := bounds.Add(image.Pt(bounds.Dx()/4, 0))
	style.DrawText(screen, at.Min, t.Text)
}

func NewTitleBar(bounds Rectangle, header string, fw *Widget) *TitleBar {
	ph := &TitleBar{}
	return ph.Init(bounds, header, fw)
}

// AddNewTitleBar adds a new TitleBar to this widget on top of the widget.
// This will set up dragging  for the panel as well.
// Also sets up the panel header for use with panel.AddNode.
func (w *Widget) AddTitleBar(height int, header string) *TitleBar {
	bounds := w.Bounds
	bounds.Min.Y -= height
	bounds.Max.Y = bounds.Min.Y + height
	tb := NewTitleBar(bounds, header, w)
	w.Widgets = append(w.Widgets, &tb.Widget)
	return tb
}

/*
// NodeResizer is for nodes that can be resized.
type NodeResizer interface {
	Node
	NodeResize(delta Point)
}

type Resizer struct {
	// Wrap Gadget
	Gadget
	// NodeResizer is the node that will be resized.
	Resizer NodeResizer
	// Resizing is true if currently resizing.
	Resizing bool
}

func (c *Resizer) Init(wide, high int, resizer NodeResizer) *Resizer {
	c.Gadget.Init(wide, high)
	c.Resizer = resizer
	return c
}

func NewResizer(wide, high int, resizer NodeResizer) *Resizer {
	pc := &Resizer{}
	return pc.Init(wide, high, resizer)
}

func (p *Resizer) EventHandle(ctx *GOG, e Event) bool {
	return Dispatch(ctx, p)(e)
}

func (b *Resizer) MousePress(ctx *GOG, ev MouseEvent) bool {
	if !b.Accept(ev.Point) {
		return false
	}

	b.Resizing = true
	b.Style.Fill.R = 0
	b.Style.Fill.G = 255
	b.Style.Fill.B = 255
	return true
}

func (b *Resizer) MouseRelease(ctx *GOG, ev MouseEvent) bool {
	b.Style.Fill.R = 0
	b.Style.Fill.G = 0
	b.Style.Fill.B = 255
	b.Resizing = false
	if !b.Accept(ev.Point) {
		return false
	}
	return true
}

func (b *Resizer) MouseHold(ctx *GOG, ev MouseEvent) bool {
	if !b.Accept(ev.Point) && !b.Resizing {
		return false
	}

	b.Style.Fill.R = 0
	b.Style.Fill.G = 255
	b.Style.Fill.B = 0
	b.Move(ev.Delta)
	if b.Resizer != nil {
		b.Resizer.NodeResize(ev.Delta)
	}

	return true
}

func (b *Resizer) Resize(delta Point) {
	// Not resizable, but if we are asked to resize,
	// and not resizing ourselelves this is a sub-corner that we
	// have to move in stead.
	if !b.Resizing {
		b.Move(delta)
	}
}

// AddNewResizer adds a new Resizer as a Control of this panel with the given size and options.
// This will set up dragging for the panel as well.
func (p *Panel) AddNewResizer(wide, high int, resizer NodeResizer) *Resizer {
	corner := NewResizer(wide, high, resizer)
	p.Controls = append(p.Controls, corner)

	delta := p.Bounds().Min.Add(image.Pt(p.Bounds().Dx()-wide, p.Bounds().Dy()-high))
	corner.Move(delta)

	return corner
}
*/
