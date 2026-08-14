package xlui

import "github.com/xmasengine/xmas/xgal"

const ListItemHeight = 14

// ListLayer is a vertical, scrollable list of selectable text items.
// It is a layer based widget so it can be dragged, focused and closed
// like any other layer.
type ListLayer struct {
	Layer

	Items      []string
	Selected   int // -1 for none.
	ItemHeight int
	Offset     int // First visible item index.

	limit    int // Max visible items.
	hoverIdx int
}

// NewList returns a new [ListLayer].
func NewList(bounds xgal.Rectangle, items ...string) *ListLayer {
	limit := bounds.Dy() / ListItemHeight
	if limit > len(items) {
		limit = len(items)
	}
	list := &ListLayer{
		Layer:      *NewLayer(bounds),
		Selected:   -1,
		ItemHeight: ListItemHeight,
		hoverIdx:   -1,
		Items:      items,
		limit:      limit,
	}
	list.Class = Class{
		Render: list.render,
		Hover:  list.hover,
		Click:  list.click,
		Wheel:  list.wheel,
		Tap:    list.tap,
	}
	return list
}

// List adds a new [ListLayer] to the UI.
func (u *UI) List(bounds xgal.Rectangle, items ...string) *ListLayer {
	list := NewList(bounds, items...)
	u.Append(&list.Layer)
	return list
}

// Select programmatically selects the i-th item
// and calls Class.Value if available.
func (l *ListLayer) Select(i int) {
	if i < 0 || i >= len(l.Items) {
		return
	}
	l.Selected = i
	l.EnsureVisible()
	if l.Class.Value != nil {
		l.Class.Value(i)
	}
}

// scrollVisible scrolls the list so the selected item is in view.
func (l *ListLayer) EnsureVisible() {
	if l.Selected < 0 || l.limit <= 0 || len(l.Items) <= l.limit {
		l.Offset = 0
		return
	}
	if l.Selected < l.Offset {
		l.Offset = l.Selected
	}
	if l.Selected >= l.Offset+l.limit {
		l.Offset = l.Selected - l.limit + 1
	}
	l.clampOffset()
}

func (l *ListLayer) clampOffset() {
	max := len(l.Items) - l.limit
	if max < 0 {
		max = 0
	}
	if l.Offset < 0 {
		l.Offset = 0
	} else if l.Offset > max {
		l.Offset = max
	}
}

// itemIndex returns the index of the item under at, or -1.
func (l *ListLayer) itemIndex(at xgal.Point) int {
	if l.ItemHeight <= 0 || !at.In(l.Bounds) {
		return -1
	}
	row := (at.Y - l.Bounds.Min.Y) / l.ItemHeight
	if row < 0 || row >= l.limit {
		return -1
	}
	return l.Offset + row
}

func (l *ListLayer) render(s *xgal.Surface) {
	l.clampOffset()

	style := l.Style
	if l.State.Focused {
		style = style.Focused()
	}
	style.DrawBox(s, l.Bounds)

	for i := 0; i < l.limit; i++ {
		idx := l.Offset + i
		if idx < 0 || idx >= len(l.Items) {
			continue
		}
		y := l.Bounds.Min.Y + i*l.ItemHeight
		item := xgal.Rect(l.Bounds.Min.X, y, l.Bounds.Max.X, y+l.ItemHeight)

		istyle := l.Style
		if idx == l.hoverIdx {
			istyle = istyle.HoverStyle()
		}
		if idx == l.Selected {
			istyle = istyle.Focused()
		}
		istyle.DrawBox(s, item)
		istyle.Ink(s, item, l.Items[idx])
	}

	// shift markers
	if l.Offset > 0 {
		xgal.Polyfill(s, style.Fore,
			l.Bounds.Max.X-5, l.Bounds.Min.Y,
			l.Bounds.Max.X, l.Bounds.Min.Y+10,
			l.Bounds.Max.X-10, l.Bounds.Min.Y+10,
		)
	} else if len(l.Items) > l.limit {
		xgal.Polyfill(s, style.Fore,
			l.Bounds.Max.X-5, l.Bounds.Max.Y,
			l.Bounds.Max.X, l.Bounds.Max.Y-10,
			l.Bounds.Max.X-10, l.Bounds.Max.Y-10,
		)
	}
}

func (l *ListLayer) hover(at xgal.Point) Reply {
	l.hoverIdx = l.itemIndex(at)
	if at.In(l.Bounds) {
		return Accept
	}
	return Ignore
}

func (l *ListLayer) click(at xgal.Point, which int) Reply {
	if which != int(xgal.MouseButtonLeft) {
		return Ignore
	}
	if !at.In(l.Bounds) {
		return Ignore
	}
	if idx := l.itemIndex(at); idx >= 0 {
		l.Select(idx)
	}
	// If a click is accepted, raise the List to the top of the layer stack.
	return Raise
}

func (l *ListLayer) wheel(at xgal.Point, delta int) Reply {
	if !at.In(l.Bounds) {
		return Ignore
	}
	l.Offset -= delta
	l.clampOffset()
	return Accept
}

func (l *ListLayer) tap(key int, mods Mods) Reply {
	switch xgal.KeyCode(key) {
	case xgal.KeyArrowUp:
		if l.Selected > 0 {
			l.Select(l.Selected - 1)
			return Accept
		}
	case xgal.KeyArrowDown:
		if l.Selected < len(l.Items)-1 {
			l.Select(l.Selected + 1)
			return Accept
		}
	}
	return Ignore
}
