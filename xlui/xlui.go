// package xlui implements simple layer based UI
// xlui consists of 3 levels: the UI, the layers in the UI and the controls
// in the layers.
//
// To simplify event handling, the UI controls the layers, and the layers
// manage the controls. The state of the elemets is managed by the container.
// Class is used to customize the layers and controls, but
// the custonization is limited to the strict neccesary.
//
// There are specific handlers in Class but these only get called if needed.
//
// A control is a leaf, and cannot contain sub controls.
// This means that for complex widgets like, for example a list, this will
// be implemented as a layer, which can be dragged, focused, and
// manipulated, like any other layer on the screen.
//
// This is somewhat unusual but it drasically simplifies event handling
// and the reduced the complexity of the UI.
// Furthermore the UI is a vertical stack, and the topmost Layer than can
// accept input must accept it.
package xlui

import "slices"

import "github.com/xmasengine/xmas/xgal"

// Mods are the key mods
type Mods struct {
	Shift   bool
	Alt     bool
	Control bool
	Meta    bool
}

// UI is the single user interface, at least for one window.
type UI struct {
	Layers     []*Layer // Layers in botttom to top order.
	Focused    *Layer   // Layer that is currently focused.
	Dragged    *Layer   // Layer that is currently being dragged.
	LastCursor xgal.Point
	Mods       Mods
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
	u.Tick(xgal.Tick())

	kr := u.pollKeys()
	if kr != Ignore {
		return kr
	}

	for mb := xgal.MouseButton(0); mb < xgal.MouseButtonMax; mb++ {
		cursor := xgal.Cursor()
		if xgal.Click(mb) {
			return u.Click(cursor, int(mb))
		}
		if xgal.Release(mb) {
			return u.Release(cursor, int(mb))
		}
		if u.LastCursor != cursor {
			delta := cursor.Sub(u.LastCursor)
			u.LastCursor = cursor
			if u.Dragged != nil {
				u.Dragged.MoveBy(delta)
			} else {
				u.Hover(cursor)
			}
		}
	}
	return Ignore
}

func (u *UI) pollKeys() Reply {
	accepted := false
	keys := xgal.Taps(nil)
	for _, key := range keys {
		switch xgal.KeyCode(key) {
		case xgal.KeyAlt:
			u.Mods.Alt = true
		case xgal.KeyShift:
			u.Mods.Shift = true
		case xgal.KeyControl:
			u.Mods.Control = true
		case xgal.KeyMeta:
			u.Mods.Meta = true
		}
		res := u.Tap(key, u.Mods)
		if res == Accept {
			accepted = true
		} else if res != Ignore {
			u.onReply(u.focusOrTopIndex(), res)
		}
	}

	keys = xgal.Keys(nil)
	for _, key := range keys {
		dur := xgal.Tapped(key)
		res := u.Key(key, dur)
		if res == Accept {
			accepted = true
		}
	}

	keys = xgal.Lifts(nil)
	for _, key := range keys {
		switch xgal.KeyCode(key) {
		case xgal.KeyAlt:
			u.Mods.Alt = false
		case xgal.KeyShift:
			u.Mods.Shift = false
		case xgal.KeyControl:
			u.Mods.Control = false
		case xgal.KeyMeta:
			u.Mods.Meta = false
		}
		res := u.Lift(key, u.Mods)
		if res == Accept {
			accepted = true
		}
	}

	chars := xgal.Chars(nil)
	res := u.Chars(chars...)
	if res == Accept {
		accepted = true
	}

	if accepted {
		return Accept
	}
	return Ignore
}

func (u *UI) focusOrTopIndex() int {
	l := u.focusOrTop()
	for i := 0; i < len(u.Layers); i++ {
		if u.Layers[i] == l {
			return i
		}
	}
	return -1
}

func (u *UI) focusOrTop() *Layer {
	if u.Focused != nil {
		return u.Focused
	}

	if len(u.Layers) < 1 {
		return nil
	}

	return u.Layers[len(u.Layers)-1]
}

func (u *UI) Tap(key xgal.KeyCode, mods Mods) Reply {
	top := u.focusOrTop()
	if top != nil {
		return top.OnTap(int(key), mods)
	}
	return Ignore
}

func (u *UI) Key(key xgal.KeyCode, duration int) Reply {
	top := u.focusOrTop()
	if top != nil {
		return top.OnKey(int(key), duration)
	}
	return Ignore
}

func (u *UI) Lift(key xgal.KeyCode, mods Mods) Reply {
	top := u.focusOrTop()
	if top != nil {
		return top.OnLift(int(key), mods)
	}
	return Ignore
}

func (u *UI) Chars(chars ...rune) Reply {
	top := u.focusOrTop()
	if top != nil {
		return top.OnChars(chars...)
	}
	return Ignore
}

func (u *UI) Tick(tick int64) Reply {
	// Tick is passed to all layers since it is used for animation.
	res := Ignore
	for i, layer := range u.Layers {
		sres := layer.OnTick(tick)
		u.onReply(i, sres)
		if sres > res {
			res = sres
		}
	}
	return res
}

func (u *UI) SetFocus(l *Layer) {
	if u.Focused != nil {
		u.Focused.State.Focused = false
	}
	u.Focused = l
	if u.Focused != nil {
		u.Focused.State.Focused = true
	}
}

func (u *UI) LayerIndex(search *Layer) int {
	for i, layer := range u.Layers {
		if layer == search {
			return i
		}
	}
	return -1
}

func (u *UI) CloseLayer(layer *Layer) {
	u.closeLayerByIndex(u.LayerIndex(layer))
}

func (u *UI) indexOOB(i int) bool {
	return (i < 0) || (i >= len(u.Layers))
}

func (u *UI) closeLayerByIndex(i int) {
	if u.indexOOB(i) {
		return
	}
	del := u.Layers[i]
	if u.Focused == del {
		u.SetFocus(nil)
	}
	u.Layers = slices.Delete(u.Layers, i, i+1)
}

func (u *UI) setFocusByIndex(i int) {
	u.SetFocus(u.Layers[i])
}

func (u *UI) raiseLayerByIndex(i int) {
	if i < len(u.Layers)-1 {
		u.Layers[i], u.Layers[i+1] = u.Layers[i+1], u.Layers[i]
		u.SetFocus(u.Layers[i+1])
	} else {
		u.SetFocus(u.Layers[i])
	}
}

func (u *UI) lowerLayerByIndex(i int) {
	if i > 0 {
		u.Layers[i], u.Layers[i-1] = u.Layers[i-1], u.Layers[i]
		u.SetFocus(u.Layers[i])
	}
}

func (u *UI) onReply(i int, res Reply) Reply {
	if res == Finish {
		u.closeLayerByIndex(i)
	} else if res == Accept {
		u.setFocusByIndex(i)
	} else if res == Raise {
		u.raiseLayerByIndex(i)
	} else if res == Lower {
		u.lowerLayerByIndex(i)
	}
	return Ignore
}

// LayersAt is an iterator of all controls that at is inside of,
// from top to bottom.
func (u *UI) LayersAt(at xgal.Point) func(func(int, *Layer) bool) {
	return func(yield func(int, *Layer) bool) {
		for i := len(u.Layers) - 1; i >= 0; i-- {
			layer := u.Layers[i]
			if layer == nil {
				continue
			}
			if !at.In(layer.Bounds) {
				continue
			}
			if !yield(i, layer) {
				break
			}
		}
	}
}

func (u *UI) Release(at xgal.Point, button int) Reply {
	if u.Dragged != nil {
		u.Dragged = nil
	}

	if u.Focused != nil {
		res := u.Focused.OnRelease(at, button)
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

		res := layer.OnClick(at, button)
		if res != Accept {
			if u.Dragged == nil && !layer.Lock {
				u.Dragged = layer
			}
		}

		if res != Ignore {
			return u.onReply(i, res)
		}
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

func (u *UI) Hover(at xgal.Point) Reply {
	for _, layer := range u.LayersAt(at) {
		res := layer.OnHover(at)
		if res != Ignore {
			return res
		}
	}
	return Ignore
}

func (u *UI) Layer(bounds xgal.Rectangle) *Layer {
	layer := NewLayer(bounds)
	u.Layers = append(u.Layers, layer)
	return layer
}

func (u *UI) Asker(bounds xgal.Rectangle, label, entry string, buttons ...string) *Layer {
	layer := NewAsker(bounds, label, entry, buttons...)
	u.Layers = append(u.Layers, layer)
	return layer
}

func (u *UI) Menu(bounds xgal.Rectangle, options ...string) *Layer {
	layer := NewMenu(bounds, options...)
	u.Layers = append(u.Layers, layer)
	return layer
}

func (u *UI) MenuWithValueOffset(bounds xgal.Rectangle, offset int, options ...string) *Layer {
	layer := NewMenuWithValueOffset(bounds, offset, options...)
	u.Layers = append(u.Layers, layer)
	return layer
}
