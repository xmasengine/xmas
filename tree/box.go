package tree

type Box struct {
	Widget
	elements    List
	orientation Orientation
}

type Orientation int

const (
	Vertical   Orientation = 0
	Horizontal Orientation = 0
)

func (o Orientation) Apply(e Element) {
	if w, ok := e.(*Box); ok {
		w.orientation = o
	}
}

func (o Orientation) Shift(at Point, size Point) Point {
	if o == Vertical {
		at.Y += size.Y
		return at
	}
	if o == Horizontal {
		at.X += size.X
		return at
	}
	return at
}

func (b Box) Contain() List {
	return b.elements
}

func (b *Box) Init(r *Root, options ...Applier) *Box {
	b.Widget.Init(r, options...)
	for _, opt := range options {
		opt.Apply(b)
	}
	return b
}

func NewBox(r *Root, options ...Applier) *Box {
	b := &Box{}
	return b.Init(r, options...)
}

func (l List) Apply(e Element) {
	if w, ok := e.(*Box); ok {
		w.elements = l
	}
}

// Drab is called when the element needs to be drawn
func (b Box) Draw(screen *Surface) {
	b.Widget.Draw(screen)
	l := b.Contain()
	for _, e := range l {
		state := e.State()
		if !state.Hide {
			e.Draw(screen)
		}
	}
}

// Place places the Box at the given bounds.
// Also moves the contained widgets.
func (b *Box) Place(bounds Bounds) (size Point) {
	delta := b.Bounds().Min.Sub(bounds.Min)
	b.bounds = bounds
	l := b.Contain()

	for _, e := range l {
		state := e.State()
		if state.Hide {
			continue
		}
		eb := e.Bounds()
		eb.Rectangle = eb.Add(delta)
		e.Place(eb)
	}

	return b.bounds.Size()
}

// Bounds are the actual absolute visual bounds of the element,
// in screen coordinates, as should be used for layout,
// ignoring any popups or overflows.
func (b Box) Bounds() Bounds {
	return b.bounds
}

// State returns the state of the element.
func (b Box) State() State {
	return b.state
}

// Modify sets the state of the element.
func (b *Box) Modify(state State) {
	b.state = state
}
