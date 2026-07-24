// package xlui implements simple layer based UI
// xlui consists of 3 levels: the UI, the layers in the UI and the controls
// in the layers.
// To simplify event handling, the UI controls the layers, and the layers
// manage the controls.
// There are specific handlers but these only get called if needed.
// A control is a leaf, and cannot contain sub controls.
// This means that for complex widgets like, for example a list, this will
// be implemented as a layer, which can be dragged, focused, and
// manipulated, like any other layer on the screen.
// This is somewhat unusual but it drasically simplifies event handling
// and the reduced the complexity of the UI.
// Furthermore the UI is a vertical stack, and the topmost Layer than can
// accept input must accept it.
package xlui

import "slices"

import "github.com/xmasengine/xmas/xgal"

// UI is the single user interface, at least for one window.
type UI struct {
	Layers []*Layer // Layers in botttom to top order.
	Groups []Group
}

// xlui is the the global UI
var xlui UI

func (u *UI) Add(l *Layer) *Layer {
	u.Layers = slices.Insert(u.Layers, 0, l)
	return l
}

type handler[T any] func(u *UI, l *Layer, t T) Reply
type getHandler[T any] func(l *Layer) handler[T]

func handleFor[T any](u *UI, gh getHandler[T], t T) Reply {
	for i := len(u.Layers) - 1; 1 >= 0; i-- {
		layer := u.Layers[i]
		if layer == nil {
			continue
		}
		handler := gh(layer)
		if handler == nil {
			continue
		}
		res := handler(u, layer, t)
		if res == Finish {
			u.Layers = slices.Delete(u.Layers, i, i+1)
		} else if res == Accept {
			break
		} else if res == Raise {
			if i < len(u.Layers)-1 {
				u.Layers[i], u.Layers[i+1] = u.Layers[i+1], u.Layers[i]
			}
			return res
		} else if res == Lower {
			if i > 0 {
				u.Layers[i], u.Layers[i-1] = u.Layers[i-1], u.Layers[i]
			}
			return res
		}
	}
	return Ignore
}

func (u *UI) Poll() Reply {
	for mb := xgal.MouseButton(0); mb < xgal.MouseButtonMax; mb++ {
		if xgal.Click(mb) {
			return u.Click(xgal.Cursor(), int(mb))
		}
	}
	return Ignore
}

func (u *UI) onReply(i int, res Reply) Reply {
	if res == Finish {
		u.Layers = slices.Delete(u.Layers, i, i+1)
	} else if res == Accept {
		return res
	} else if res == Raise {
		if i < len(u.Layers)-1 {
			u.Layers[i], u.Layers[i+1] = u.Layers[i+1], u.Layers[i]
		}
		return res
	} else if res == Lower {
		if i > 0 {
			u.Layers[i], u.Layers[i-1] = u.Layers[i-1], u.Layers[i]
		}
		return res
	}
	return Ignore
}

func (u *UI) Click(at xgal.Point, button int) Reply {
	for i := len(u.Layers) - 1; i >= 0; i-- {
		layer := u.Layers[i]
		if layer == nil {
			continue
		}
		if !at.In(layer.Bounds) {
			continue
		}
		res := layer.Click(at, button)
		return u.onReply(i, res)
	}
	return Ignore
}

func (u *UI) Render(s *xgal.Surface) {
	for i := 0; i < len(u.Layers); i++ {
		layer := u.Layers[i]
		if layer != nil {
			layer.Render(s)
		}
	}
}

func (u *UI) Layer(bounds xgal.Rectangle) *Layer {
	layer := NewLayer(bounds)
	u.Layers = append(u.Layers, layer)
	return layer
}
