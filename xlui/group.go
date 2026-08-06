package xlui

// Group is a group of Controls.
type Group struct {
	Controls []*Control
}

// Stack is a stack of layers that belong together.
type Stack struct {
	Layers []*Layer
}
