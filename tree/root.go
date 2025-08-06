package tree

// Root is the root element of a UI.
// It is possble to use multiple roots however only one should be active
// at the same time.
type Root struct {
	Box
}

func (r *Root) Init(options ...Applier) *Root {
	r.Box.Init(nil, options...)
	for _, opt := range options {
		opt.Apply(r)
	}
	return r
}

func NewRoot(options ...Applier) *Root {
	r := Root{}
	return r.Init(options...)
}

// Update is called 60 times per second.
// Input should be checked during this function.
func (r *Root) Update() error {
	return nil
}

// Draw is called when the UI needs to be drawn
func (r *Root) Draw(screen *Surface) {

}

// Layout is called when the contents of the element need to be laid out.
// The element should accept that the available size is less than its
// real size and draw it appropiately, such as scrolling.
// The returned elementWidth and elementHeight must be smaller than
// or equal to the available width.
func (r *Root) Layout(availableWidth, availableHeight int) (elementWidth, elementHeight int) {

	return availableWidth, availableHeight
}
