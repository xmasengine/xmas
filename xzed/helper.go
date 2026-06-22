package xzed

import "image"

import (
	"github.com/xmasengine/xmas/xui"
)

const helpText = `
F1:      toggle help
Ctrl-C:  toggle Picker
Ctrl-S:  save map
Ctrl-M:  load map
Alt-L:   add layer
Shift-L: layer up
Shift-K: layer down
`

// Helper is a widget that displays a help text.
type Helper struct {
	xui.Label
	// OnClick callback.
	OnClick func(*Helper)
}

type HelperClass struct {
	*xui.LabelClass
	*Helper
}

func NewHelperClass(b *Helper) *HelperClass {
	res := &HelperClass{Helper: b}
	res.LabelClass = xui.NewLabelClass(&b.Label)
	return res
}

func NewHelper(bounds image.Rectangle, text string, cb func(*Helper)) *Helper {
	h := &Helper{}
	h.Init(bounds, text, cb)
	return h
}

func (h *Helper) Init(bounds image.Rectangle, text string, cb func(*Helper)) *Helper {
	// Set up the tile box.
	h.Label.Init(bounds, text)
	h.Class = NewHelperClass(h)
	return h
}

func (t *HelperClass) Render(r *xui.Root, Surface *xui.Surface) {
	t.LabelClass.Render(r, Surface)
}

func AddHelper(w *xui.Widget, bounds image.Rectangle, title, text string, cb func(*Helper)) *Helper {
	h := NewHelper(bounds, text, cb)
	h.AddTitleBar(10, title)
	w.AddWidget(&h.Widget)
	return h
}
