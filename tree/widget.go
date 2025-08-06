package tree

// Widget is a basic widget. It implements element.
type Widget struct {
	state  State
	bounds Bounds
	style  Style
	root   *Root
}

func (w *Widget) Init(r *Root, options ...Applier) *Widget {
	w.root = r
	if r != nil {
		w.style = r.style
	} else {
		w.style = DefaultStyle()
	}
	for _, opt := range options {
		opt.Apply(w)
	}
	return w
}

func NewWidget(r *Root, options ...Applier) *Widget {
	w := &Widget{}
	return w.Init(r, options...)
}

func (s Style) Apply(e Element) {
	if w, ok := e.(*Widget); ok {
		w.style = s
	}
}

func (b Bounds) Apply(e Element) {
	if w, ok := e.(*Widget); ok {
		w.bounds = b
	}
}

var _ Applier = Style{}

// Draw is called when the element needs to be drawn
func (w Widget) Draw(screen *Surface) {
	w.style.DrawBox(screen, w.bounds.Rectangle)
}

// Place places the widget at the given bounds.
func (w *Widget) Place(bounds Bounds) (size Point) {
	w.bounds = bounds
	return w.bounds.Size()
}

// Bounds are the actual absolute visual bounds of the element,
// in screen coordinates, as should be used for layout,
// ignoring any popups or overflows.
func (w Widget) Bounds() Bounds {
	return w.bounds
}

// State returns the state of the element.
func (w Widget) State() State {
	return w.state
}

// Modify sets the state of the element.
func (w *Widget) Modify(state State) {
	w.state = state
}
