package xui

import "image"
import "image/color"

type FrameClass struct {
	*Frame
	*BoxClass
}

func NewFrameClass(b *Frame) *FrameClass {
	res := &FrameClass{Frame: b}
	res.BoxClass = NewBoxClass(&b.Box)
	return res
}

// Frame is a frame that can display and image.
// If the Image is not nil, the Frame will be of the exact size of
// the image.
type Frame struct {
	Box
	Image  *Surface
	Extent Rectangle
}

func (f *Frame) Init(bounds Rectangle, img *Surface) *Frame {
	f.Image = img

	extent := bounds
	if f.Image != nil {
		wide, high := img.Size()
		extent.Max.X = extent.Min.X + wide
		extent.Max.Y = extent.Min.Y + high
	}
	f.Box.Init(bounds)
	f.Style.Fill = color.RGBA{0, 0, 0, 0}
	f.Style.Border = color.RGBA{0, 128, 128, 128}
	f.Extent = extent
	f.Class = NewFrameClass(f)
	return f
}

func NewFrame(bounds Rectangle, image *Surface) *Frame {
	p := &Frame{}
	return p.Init(bounds, image)
}

func (f *FrameClass) Render(r *Root, screen *Surface) {
	if f.Image != nil {
		clipped := screen.SubImage(f.Frame.Bounds).(*Surface)
		center := f.Frame.Bounds.Inset(f.Frame.Style.Margin.X)
		opts := DrawOptions{}
		opts.GeoM.Translate(
			float64(center.Min.X),
			float64(center.Min.Y),
		)
		clipped.DrawImage(f.Image, &opts)
	} else {
		f.BoxClass.Render(r, screen)
	}
	f.Frame.Style.DrawBox(screen, f.Frame.Bounds.Add(f.Frame.Style.Margin))
}

func (f *Frame) SetImage(img *Surface) {
	f.Image = img
	// Correct the size
	if f.Image != nil {
		wide, high := f.Image.Size()
		f.Bounds.Max = f.Bounds.Min.Add(image.Pt(wide, high))
	}
}

// Adds a Frame as a child of this widget.
func (w *Widget) AddFrame(bounds Rectangle, image *Surface) *Frame {
	frame := NewFrame(bounds, image)
	w.Widgets = append(w.Widgets, &frame.Widget)
	return frame
}

type CotClass struct {
	*Cot
	*BoxClass
}

func NewCotClass(c *Cot) *CotClass {
	res := &CotClass{Cot: c}
	res.BoxClass = NewBoxClass(&c.Box)
	return res
}

// Cot is a rectangular area that contains an image selection.
// Normally it is used by a selector.
type Cot struct {
	// Box inherited
	Box
	// Tile the Cot box is pointing at.
	Tile image.Point
	// Called on selection (click, enter key...)
	Select func(at image.Point)
}

func (c *Cot) Init(bounds Rectangle) *Cot {
	c.Box.Init(bounds)
	c.Style.Fill = color.RGBA{0, 0, 0, 0}
	c.Class = NewCotClass(c)
	return c
}

func NewCot(bounds Rectangle) *Cot {
	p := &Cot{}
	return p.Init(bounds)
}

// Call this after changing the "tile" the Cot is pointing to.
func (c *Cot) UpdateFromTile(min Point, margin Point) {
	mct := c.Tile
	size := c.Bounds.Size()
	rmin := image.Pt(mct.X*size.X, mct.Y*size.Y)
	rmin = rmin.Add(min).Add(image.Pt(margin.X*2, margin.Y*2))
	rmax := rmin.Add(size)
	c.Bounds.Min = rmin
	c.Bounds.Max = rmax
}

type ChooserClass struct {
	*Chooser
	*FrameClass
}

func NewChooserClass(c *Chooser) *ChooserClass {
	res := &ChooserClass{Chooser: c}
	res.FrameClass = NewFrameClass(&c.Frame)
	return res
}

// Chooser allows to choose a selection of a Frame using a Cot.
type Chooser struct {
	Frame        // Frame is the widget we extend.
	Hovered  Cot // Cot being hovered for the image selection.
	Selected Cot // Graphical tile cot box for the currently selected image part.
	// Callback if a selection of a Frame is chosen.
	// Which will be 0 for te right mouse, 1 for the left mouse button and 2 for the middle, etc.
	Select func(*Chooser)
	Result int
	Which  int
	At     Point
}

// Init initialize the Chooser.
func (c *Chooser) Init(bounds Rectangle, img *Surface, cotSize Point, cb func(*Chooser)) *Chooser {
	c.Frame.Init(bounds, img)
	cotBounds := bounds
	cotBounds.Max.X = cotBounds.Min.X + cotSize.X
	cotBounds.Max.Y = cotBounds.Min.Y + cotSize.Y
	c.Hovered.Init(cotBounds)
	c.Hovered.Style.Border = color.RGBA{255, 255, 255, 255}
	c.Hovered.State.Hide = false

	// Set up a Cot which is a box that marks chosen tile in the map.
	c.Selected.Init(cotBounds)
	c.Selected.Style.Border = color.RGBA{255, 255, 0, 255}
	c.Selected.State.Hide = true
	c.Select = cb
	c.updateCot()
	c.updateSelectedCot()

	c.Class = NewChooserClass(c)
	return c
}

// New creates a Chooser. The size should be that of the image.
// The editor itself is invisible but it has activatable child widgets.
func NewChooser(bounds Rectangle, image *Surface, cotSize Point, cb func(*Chooser)) *Chooser {
	c := &Chooser{}
	return c.Init(bounds, image, cotSize, cb)
}

func (c *Chooser) updateMouse(at Point) bool {
	// New mouse tile location.
	mx := (at.X - c.Frame.Bounds.Min.X - c.Frame.Style.Margin.X*2) / c.Hovered.Bounds.Dx()
	my := (at.Y - c.Frame.Bounds.Min.Y - c.Frame.Style.Margin.Y*2) / c.Hovered.Bounds.Dy()

	c.Hovered.Tile.X = max(0, mx)
	c.Hovered.Tile.Y = max(0, my)
	if c.Image != nil {
		w, h := c.Image.Size()
		c.Hovered.Tile.X = min(c.Hovered.Tile.X, int(w))
		c.Hovered.Tile.Y = min(c.Hovered.Tile.Y, int(h))
	}

	return true
}

func (c *Chooser) updateCot() bool {
	c.Hovered.UpdateFromTile(c.Frame.Bounds.Min, c.Frame.Style.Margin)
	return true
}

func (c *Chooser) selectCot(which int) bool {
	c.Selected.Tile = c.Hovered.Tile
	if c.Select != nil {
		c.Which = which
		c.Select(c)
		c.Which = -1
	}
	c.Selected.State.Hide = false

	return c.updateSelectedCot()
}

func (c *Chooser) updateSelectedCot() bool {
	c.Selected.UpdateFromTile(c.Frame.Bounds.Min, c.Frame.Style.Margin)
	return true
}

func (c *ChooserClass) OnMouseMove(ev MouseEvent) bool {
	c.updateMouse(ev.At)
	c.updateCot()
	return true
}

func (c *ChooserClass) OnMousePress(ev MouseEvent) bool {
	c.selectCot(ev.Button)
	return true
}

func (c *ChooserClass) OnMouseHold(ev MouseEvent) bool {
	return true
}

func (c *ChooserClass) OnActionHover(e ActionEvent) bool {
	c.Chooser.State.Hover = true
	c.Hovered.State.Hide = false
	c.updateMouse(e.At)
	c.updateCot()
	return true
}

func (c *ChooserClass) OnActionCrash(e ActionEvent) bool {
	c.Chooser.State.Hover = false
	c.Hovered.State.Hide = true
	return true
}

// Move moves the Chooser.
func (p *Chooser) Move(delta Point) {
	p.Frame.Move(delta)
	p.Hovered.Move(delta)
	p.Selected.Move(delta)
}

// Draw is the gadget draw handler.
func (c *ChooserClass) Render(r *Root, screen *Surface) {
	c.FrameClass.Render(r, screen)
	if !c.Hovered.State.Hide {
		c.Hovered.Class.Render(r, screen)
	}
	if !c.Selected.State.Hide {
		c.Selected.Class.Render(r, screen)
	}
}

func (p *Chooser) SetImage(img *Surface) {
	p.Frame.SetImage(img)
	p.updateCot()
}

// AddChooser a new image Chooser as a child of this widget
func (w *Widget) AddChooser(
	bounds Rectangle, image *Surface, cotSize Point, cb func(*Chooser),
) *Chooser {
	chooser := NewChooser(bounds, image, cotSize, cb)
	w.Widgets = append(w.Widgets, &chooser.Widget)
	return chooser
}
