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

import "github.com/xmasengine/xmas/xgal"

// xlui is the the global UI
var xlui UI

// Mods are the key mods
type Mods struct {
	Shift   bool
	Alt     bool
	Control bool
	Meta    bool
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func Add(l *Layer) *Layer {
	return xlui.Add(l)
}

func Asker(bounds xgal.Rectangle, label, entry string, buttons ...string) *Layer {
	return xlui.Asker(bounds, label, entry, buttons...)
}

func Chars(chars ...rune) Reply {
	return xlui.Chars(chars...)
}

func Click(at xgal.Point, button int) Reply {
	return xlui.Click(at, button)
}

func CloseLayer(layer *Layer) {
	xlui.CloseLayer(layer)
}

func Hover(at xgal.Point) Reply {
	return xlui.Hover(at)
}

func Key(key xgal.KeyCode, duration int) Reply {
	return xlui.Key(key, duration)
}

func AddLayer(bounds xgal.Rectangle) *Layer {
	return xlui.Layer(bounds)
}

func LayerIndex(search *Layer) int {
	return xlui.LayerIndex(search)
}

func LayersAt(at xgal.Point) func(func(int, *Layer) bool) {
	return xlui.LayersAt(at)
}

func Lift(key xgal.KeyCode, mods Mods) Reply {
	return xlui.Lift(key, mods)
}

func Menu(bounds xgal.Rectangle, options ...string) *Layer {
	return xlui.Menu(bounds, options...)
}

func MenuBar(bounds xgal.Rectangle, subs ...SubMenuOption) *Layer {
	return xlui.MenuBar(bounds, subs...)
}

func MenuWithValueOffset(bounds xgal.Rectangle, offset int, options ...string) *Layer {
	return xlui.MenuWithValueOffset(bounds, offset)
}

func Poll() Reply {
	return xlui.Poll()
}

func Release(at xgal.Point, button int) Reply {
	return xlui.Release(at, button)
}

func Render(s *xgal.Surface) {
	xlui.Render(s)
}

func SetFocus(l *Layer) {
	xlui.SetFocus(l)
}

func Tap(key xgal.KeyCode, mods Mods) Reply {
	return xlui.Tap(key, mods)
}

func Tick(tick int64) Reply {
	return xlui.Tick(tick)
}

func Wheel(at xgal.Point, delta int) Reply {
	return xlui.Wheel(at, delta)
}
