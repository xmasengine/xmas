// package xlui implements simple layer based UI
// xlui consists of 3 levels: the UI, the layers in the UI and the controls
// in the layers.
//
// To simplify event handling, the UI controls the layers, and the layers
// manage the controls. The state of the elements is managed by the container.
// Class is used to customize the layers and controls, but
// the customization is limited to the strict necessary.
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

import "strconv"
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

func Append(l *Layer) *Layer {
	return xlui.Append(l)
}

func Asker(bounds xgal.Rectangle, label, entry string, buttons ...string) *Layer {
	return xlui.Asker(bounds, label, entry, buttons...)
}

func Ask(x, y, w, h int, label, entry string, handler func(name string) bool) *Layer {
	ask := xlui.Asker(xgal.Bound(x, y, w, h), label, entry, " X ", " V ")
	ask.Class.Entry = func(v string) Reply {
		if handler(v) {
			return Finish
		} else {
			return Accept
		}
	}
	return ask
}

func Complain(x, y, w, h int, err error) *Layer {
	complain := xlui.Complain(xgal.Bound(x, y, w, h), err)
	return complain
}

func Display(x, y, w, h int, text string) *Layer {
	display := xlui.Display(xgal.Bound(x, y, w, h), text)
	return display
}

func AskString(x, y, w, h int, prompt string, str *string) *Layer {
	on := func(sres string) bool {
		*str = sres
		return true
	}
	return Ask(x, y, w, h, prompt, *str, on)
}

func AskInt(x, y, w, h int, prompt string, i *int) *Layer {
	on := func(sres string) bool {
		res, err := strconv.Atoi(sres)
		if err == nil {
			*i = res
			return true
		} else {
			Complain(x+20, y+20, w, h, err)
			return false
		}
	}
	return Ask(x, y, w, h, prompt, strconv.Itoa(*i), on)
}

func AskByte(x, y, w, h int, prompt string, i *byte) *Layer {
	on := func(sres string) bool {
		res, err := strconv.Atoi(sres)
		if err == nil {
			*i = byte(res)
			return true
		} else {
			Complain(x+20, y+20, w, h, err)
			return false
		}
	}
	return Ask(x, y, w, h, prompt, strconv.Itoa(int(*i)), on)
}

func Dialog(x, y, w, h int, prompt string, handler func(v int) bool, buttons ...string) *Layer {
	value := func(v int) Reply {
		if handler(v) {
			return Finish
		} else {
			return Accept
		}
	}

	dialog := xlui.Dialog(xgal.Bound(x, y, w, h), prompt, buttons...)
	dialog.Class.Value = value
	return dialog
}

func DialogBool(x, y, w, h int, prompt string, handler func(b bool) bool, tr, fa string) *Layer {
	value := func(v int) bool {
		return handler(v == 0)
	}
	return Dialog(x, y, w, h, prompt, value, tr, fa)
}

func Choose(x, y, tw, th int, title string, texture *xgal.Surface, handler func(x, y int) bool) *Layer {
	return xlui.Choose(x, y, tw, th, title, texture, handler)
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
	return xlui.MenuWithValueOffset(bounds, offset, options...)
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
