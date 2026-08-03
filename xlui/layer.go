package xlui

import (
	"github.com/xmasengine/xmas/xgal"
)

// Reply is the result of several event handlers.
// Event handlers must strictly observe the meaning of Reply.
// Otherwise the widgets, in particylar widget focus may malfunction.
type Reply int

const (
	Ignore  Reply = iota // Ignore: the widget ignored the input, other widgets *should* process it.
	Accept               // Accept: the widget accepted the input, other widgets *must not* process it.
	Raise                // Raise: the widget accepted and needs to be raised higher in the layer stack.
	Lower                // Lower: the widget accepted and needs to be lowered in the layer stack.
	Finish               // Finish: the widget is done processing and should be considered closed.
	Proceed              // Proceed: the widget accepted the input but other widgets should continue processing. Only for sub widgets.
)

// Orientation is the layout orientation for layers in a group.
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// Layer is a layer in the UI.
type Layer struct {
	// Class for custom or type specific behavior.
	Class Class

	Controls []*Control
	Bounds   xgal.Rectangle
	Clip     *xgal.Rectangle
	Style
	From        xgal.Point
	Lock        bool        // lock means the layer may not be dragged
	Orientation Orientation // layout orientation in the group
	State       State       // state of the layer

	Hovered *Control // Hovered is currently hovered control or none if nil.
	Clicked *Control // Clicked is currently clicked control or none if nil.
	Focused *Control // Clicked is currently focused control or none if nil.
}

func NewLayer(bounds xgal.Rectangle) *Layer {
	return &Layer{Bounds: bounds, Style: DefaultStyle(), Orientation: Horizontal}
}

func (l *Layer) SetFocus(c *Control) {
	if l.Focused != nil {
		l.Focused.State.Focused = false
	}
	l.Focused = c
	if l.Focused != nil {
		l.Focused.State.Focused = true
	}
}

func (l Layer) Render(s *xgal.Surface) {
	if l.Class.Render == nil {
		if l.State.Focused {
			l.Style.Focused().DrawBox(s, l.Bounds)
		} else {
			l.Style.DrawBox(s, l.Bounds)
		}
	} else {
		l.Class.Render(s)
	}

	for i := len(l.Controls) - 1; i >= 0; i-- {
		ctrl := l.Controls[i]
		ctrl.Render(s)
	}
}

// MoveBy moves all children relative to current position.
func (l *Layer) MoveBy(delta xgal.Point) {
	l.Bounds = l.Bounds.Add(delta)
	for i := 0; i < len(l.Controls); i++ {
		l.Controls[i].MoveBy(delta)
	}
}

// Appends adds a control to this layer and lays it out by a simple line algorithm.
func (l *Layer) Append(ctrl *Control) *Control {
	margin := l.Style.Margin
	if len(l.Controls) == 0 {
		ctrl.Bounds = xgal.Bound(ctrl.Bounds.Min.X+margin.X, ctrl.Bounds.Min.Y+margin.Y, ctrl.Bounds.Dx(), ctrl.Bounds.Dy())
	} else {
		last := l.Controls[len(l.Controls)-1]
		if l.Orientation == Horizontal && last.Bounds.Dx()+ctrl.Bounds.Dx() < l.Bounds.Dx() {
			// fits on the line
			ctrl.Bounds = xgal.Bound(last.Bounds.Max.X+margin.X, last.Bounds.Min.Y, ctrl.Bounds.Dx(), ctrl.Bounds.Dy())

		} else {
			ctrl.Bounds = xgal.Bound(ctrl.Bounds.Min.X+margin.X, last.Bounds.Max.Y+margin.Y, ctrl.Bounds.Dx(), ctrl.Bounds.Dy())
		}
	}
	l.Controls = append(l.Controls, ctrl)
	return ctrl
}

func (l *Layer) Label(text string) *Control {
	at := l.Bounds.Min
	ctrl := NewLabel(at, text)
	return l.Append(ctrl)
}

func (l *Layer) Button(text string) *Control {
	at := l.Bounds.Min
	ctrl := NewButton(at, text)
	return l.Append(ctrl)
}

// AllControls is an iterator of all controls which are not nil,
// and where check returns true from top to bottom.
func (l *Layer) AllControlsWhere(check func(*Control) bool) func(func(*Control) bool) {
	return func(yield func(*Control) bool) {
		for i := len(l.Controls) - 1; i >= 0; i-- {
			ctrl := l.Controls[i]
			if ctrl == nil {
				continue
			}
			if !check(ctrl) {
				continue
			}
			if !yield(ctrl) {
				break
			}
		}
	}
}

// ControlsAt is an iterator of all controls that at is inside of,
// from top to bottom.
func (l *Layer) ControlsAt(at xgal.Point) func(func(*Control) bool) {
	return l.AllControlsWhere(func(ctrl *Control) bool {
		return at.In(ctrl.Bounds)
	})
}

func (l *Layer) OnClick(at xgal.Point, button int) Reply {
	if l.Class.Click != nil {
		return l.Class.Click(at, button)
	}

	if l.Hovered != nil && at.In(l.Hovered.Bounds) {
		l.SetFocus(l.Hovered)
	}

	for ctrl := range l.ControlsAt(at) {
		if ctrl.Class.Click != nil {
			res := ctrl.Class.Click(at, button)
			if res != Ignore {
				ctrl.State.Clicked = true
				l.Clicked = ctrl
			}
			return res
		}
	}
	// If the layer was clicked, andn nothing else used the click, raise the layer.
	return Raise
}

func (l *Layer) OnHover(at xgal.Point) Reply {
	if l.Class.Hover != nil {
		return l.Class.Hover(at)
	}

	if l.Hovered != nil && !at.In(l.Hovered.Bounds) {
		l.Hovered.State.Hovered = false
		l.Hovered = nil
	}

	for ctrl := range l.ControlsAt(at) {
		if ctrl.Class.Hover != nil {
			res := ctrl.Class.Hover(at)
			if res != Ignore {
				ctrl.State.Hovered = true
				l.Hovered = ctrl
			}
			return res
		}
	}
	// If the layer was hovered, and nothing else used the click, raise the layer.
	return Raise
}

func (l *Layer) OnRelease(at xgal.Point, button int) Reply {
	if l.Class.Release != nil {
		return l.Class.Release(at, button)
	}
	if l.Clicked != nil {
		ctrl := l.Clicked
		if ctrl.Class.Release != nil {
			res := ctrl.Class.Release(at, button)
			if res != Ignore {
				ctrl.State.Clicked = false
				l.Clicked = nil
			}
			return res
		}
	}

	return Ignore
}

func (l *Layer) OnTap(key int, mods Mods) Reply {
	if l.Class.Tap != nil {
		return l.Class.Tap(key, mods)
	}

	if l.Focused != nil {
		if l.Focused.Class.Tap != nil {
			return l.Focused.Class.Tap(key, mods)
		}
	}

	for ctrl := range l.AllControlsWhere(func(ctrl *Control) bool {
		return ctrl.Class.Tap != nil
	}) {
		return ctrl.Class.Tap(key, mods)
	}
	return Ignore
}

func (l *Layer) OnKey(key int, dur int) Reply {
	if l.Class.Key != nil {
		return l.Class.Key(key, dur)
	}

	if l.Focused != nil {
		if l.Focused.Class.Key != nil {
			return l.Focused.Class.Key(key, dur)
		}
	}

	for ctrl := range l.AllControlsWhere(func(ctrl *Control) bool {
		return ctrl.Class.Key != nil
	}) {
		return ctrl.Class.Key(key, dur)
	}
	return Ignore
}

func (l *Layer) OnChars(chars ...rune) Reply {
	if l.Class.Key != nil {
		return l.Class.Chars(chars...)
	}

	if l.Focused != nil {
		if l.Focused.Class.Chars != nil {
			return l.Focused.Class.Chars(chars...)
		}
	}

	for ctrl := range l.AllControlsWhere(func(ctrl *Control) bool {
		return ctrl.Class.Chars != nil
	}) {
		return ctrl.Class.Chars(chars...)
	}
	return Ignore
}

func (l *Layer) OnLift(key int, mods Mods) Reply {
	if l.Class.Lift != nil {
		return l.Class.Lift(key, mods)
	}

	if l.Focused != nil {
		if l.Focused.Class.Lift != nil {
			return l.Focused.Class.Lift(key, mods)
		}
	}

	for ctrl := range l.AllControlsWhere(func(ctrl *Control) bool {
		return ctrl.Class.Lift != nil
	}) {
		return ctrl.Class.Lift(key, mods)
	}
	return Ignore
}

// NewAsker creates a new simple dialog pop up layer.
func NewAsker(bounds xgal.Rectangle, label, entry string, buttons ...string) *Layer {
	asker := NewLayer(bounds)
	asker.Label(label)
	asker.Orientation = Vertical
	e := asker.Entry(entry)

	e.Class.Entry = func(entry string) Reply {
		if asker.Class.Entry != nil {
			return asker.Class.Entry(entry)
		}
		return Ignore
	}
	// full width
	e.Bounds.Max.X = asker.Bounds.Max.X - asker.Style.Margin.X

	for i, buttonText := range buttons {
		if i > 0 {
			asker.Orientation = Horizontal
		}
		button := asker.Button(buttonText)
		button.Class.Click = func(at xgal.Point, mouseButton int) Reply {
			if asker.Class.Value != nil {
				asker.Class.Value(i)
			}
			// Call the entry if needed.
			if i > 0 && e.Text != "" && e.Class.Entry != nil {
				e.Class.Entry(e.Text)
			}
			return Finish
		}
	}

	return asker
}

const MenuClosedValue = -1

// NewMenuWithValueOffset creates a simple vertical menu with buttons.
// The value of the buttons if offset with offset.
// Use Class.Value as a callback get the clicked button index.
func NewMenuWithValueOffset(bounds xgal.Rectangle, offset int, options ...string) *Layer {
	menu := NewLayer(bounds)
	menu.Style = DefaultStyle()
	menu.Orientation = Vertical
	width := menu.Bounds.Dx()
	height := menu.Bounds.Dy() + menu.Style.Margin.Y
	for i, option := range options {
		button := menu.Button(option)
		button.Style = MenuStyle()
		button.Value = i + offset
		button.Class.Click = func(at xgal.Point, mouseButton int) Reply {
			if mouseButton == int(xgal.MouseButtonRight) {
				if menu.Class.Value != nil {
					menu.Class.Value(MenuClosedValue)
				}
				return Finish
			}

			// Call the entry if needed.
			if menu.Class.Entry != nil {
				menu.Class.Entry(option)
			}

			if menu.Class.Value != nil {
				menu.Class.Value(button.Value)
			}
			return Finish
		}
		buttonWidth := button.Bounds.Dx() + button.Style.Margin.X*2
		if buttonWidth > width {
			width = buttonWidth
		}
		height += button.Bounds.Dy() + menu.Style.Margin.Y
	}
	menu.Bounds = xgal.Bound(menu.Bounds.Min.X, menu.Bounds.Min.Y, width, height)
	return menu
}

// NewMenu creates a simple vertical menu with buttons.
// The value of the buttons will be 0..(len options)-1.
// Use Class.Value as a callback get the clicked button index.
func NewMenu(bounds xgal.Rectangle, options ...string) *Layer {
	return NewMenuWithValueOffset(bounds, 0, options...)
}

type SubMenuOption struct {
	Name    string   // Name of the menu and top button
	Options []string // Options in the menu.
	Value   int      // Value offset.
}

func SubMenuWithOffset(name string, offset int, options ...string) SubMenuOption {
	return SubMenuOption{
		Name:    name,
		Options: options,
		Value:   offset,
	}
}

func SubMenu(name string, options ...string) SubMenuOption {
	return SubMenuWithOffset(name, 0, options...)
}

const SubMenuDefaultOffset = 100

// NewMenuBar creates a horizontal menu bar with vertical pop up menus.
func NewMenuBar(ui *UI, bounds xgal.Rectangle, subs ...SubMenuOption) *Layer {
	bar := NewLayer(bounds)
	bar.Orientation = Horizontal
	group := &Group{}
	height := bar.Bounds.Dy()
	var menu *Layer

	for si, sub := range subs {
		toggle := bar.Toggle(sub.Name, group)
		offset := sub.Value
		if offset == 0 {
			offset = SubMenuDefaultOffset * si
		}

		toggle.Class.LinkClick(func(at xgal.Point, button int) Reply {
			if menu != nil {
				ui.CloseLayer(menu)
				menu = nil
			}
			menu = ui.MenuWithValueOffset(xgal.Bound(toggle.Bounds.Min.X, toggle.Bounds.Max.Y, 0, 0), offset, sub.Options...)
			menu.Class.LinkValue(func(value int) Reply {
				toggle.State.Clicked = false
				menu = nil
				if value != MenuClosedValue && bar.Class.Value != nil {
					bar.Class.Value(value)
				}
				return Finish
			})
			menu.Class.Entry = func(entry string) Reply {
				if bar.Class.Entry != nil {
					bar.Class.Entry(entry)
				}
				return Accept
			}
			return Accept
		})
		toggleHeight := toggle.Bounds.Dy() + toggle.Style.Margin.Y*2
		if toggleHeight > height {
			height = toggleHeight
		}
	}
	width := bar.Bounds.Dx()
	bar.Bounds = xgal.Bound(bar.Bounds.Min.X, bar.Bounds.Min.Y, width, height)

	return bar
}

func (u *UI) MenuBar(bounds xgal.Rectangle, subs ...SubMenuOption) *Layer {
	layer := NewMenuBar(u, bounds, subs...)
	u.Layers = append(u.Layers, layer)
	return layer
}
