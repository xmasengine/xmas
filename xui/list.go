package xui

import "github.com/xmasengine/xmas/xgal"

const ListItemHeight = 14

// ListLayer is a vertical list of selectable text items with optional scrolling.
type ListLayer struct {
	Layer
	Items      []string
	OnSelect   func(index int)
	Selected   int // -1 for none
	ItemHeight int
	Limit      int // max visible items, 0 = show all
	Offset     int // first visible item index
	hoverIdx   int
}

// List returns a new [ListLayer].
func List(bounds xgal.Rectangle) *ListLayer {
	l := &ListLayer{Selected: -1, hoverIdx: -1, ItemHeight: ListItemHeight}
	l.Layer = MakeLayer(bounds)
	return l
}

var _ Widget = &ListLayer{}

// AddItem appends an auto-positioned item.
func (l *ListLayer) AddItem(text string) {
	l.Items = append(l.Items, text)
}

// SelectItem programmatically selects the i-th item.
func (l *ListLayer) SelectItem(i int) {
	if i < 0 || i >= len(l.Items) {
		return
	}
	l.Selected = i
	l.ensureVisible()
	if l.OnSelect != nil {
		l.OnSelect(i)
	}
}

// EnsureVisible scrolls the list so the selected item is in view.
func (l *ListLayer) EnsureVisible() {
	if l.Selected < 0 || l.Limit <= 0 || len(l.Items) <= l.Limit {
		l.Offset = 0
		return
	}
	if l.Selected < l.Offset {
		l.Offset = l.Selected
	}
	if l.Selected >= l.Offset+l.Limit {
		l.Offset = l.Selected - l.Limit + 1
	}
	maxOff := len(l.Items) - l.Limit
	if l.Offset > maxOff {
		l.Offset = maxOff
	}
}

func (l *ListLayer) ensureVisible() { l.EnsureVisible() }

func (l *ListLayer) visibleCount() int {
	n := len(l.Items)
	if l.Limit > 0 && n > l.Limit {
		return l.Limit
	}
	return n
}

// AddList is a helper to add a [ListLayer] to a [Layer].
func (m *Layer) AddList(bounds xgal.Rectangle) *ListLayer {
	l := List(bounds)
	m.Add(l)
	return l
}

func (l *ListLayer) clampOffset() {
	maxOff := len(l.Items) - l.visibleCount()
	if maxOff < 0 {
		maxOff = 0
	}
	if l.Offset < 0 {
		l.Offset = 0
	} else if l.Offset > maxOff {
		l.Offset = maxOff
	}
}

func (l *ListLayer) Poll() Reply {
	pos := xgal.Cursor()
	l.hoverIdx = -1

	// Clamp offset in case items changed externally
	l.clampOffset()

	// Mouse wheel scrolling
	if pos.In(l.Bounds) && l.Limit > 0 {
		_, wy := xgal.Wheel()
		if wy != 0 {
			l.Offset -= int(wy)
			l.clampOffset()
		}
	}

	// Arrow key navigation when hovering
	if pos.In(l.Bounds) && l.Selected >= 0 {
		if xgal.Tap(xgal.KeyArrowDown) && l.Selected < len(l.Items)-1 {
			l.SelectItem(l.Selected + 1)
			return Accept
		}
		if xgal.Tap(xgal.KeyArrowUp) && l.Selected > 0 {
			l.SelectItem(l.Selected - 1)
			return Accept
		}
	}

	// Mouse click on visible items
	vc := l.visibleCount()
	for i := 0; i < vc; i++ {
		idx := l.Offset + i
		if idx < 0 || idx >= len(l.Items) {
			continue
		}
		y := l.Bounds.Min.Y + i*l.ItemHeight
		itemBounds := xgal.Rect(l.Bounds.Min.X, y, l.Bounds.Max.X, y+l.ItemHeight)
		if pos.In(itemBounds) {
			l.hoverIdx = idx
			if xgal.Click(xgal.MouseButtonLeft) {
				l.SelectItem(idx)
				return Accept
			}
			return Accept
		}
	}
	return Ignore
}

func (l *ListLayer) Render(s *xgal.Surface) {
	l.Layer.Render(s)

	vc := l.visibleCount()
	for i := 0; i < vc; i++ {
		idx := l.Offset + i
		if idx < 0 || idx >= len(l.Items) {
			continue
		}
		text := l.Items[idx]
		y := l.Bounds.Min.Y + i*l.ItemHeight
		itemBounds := xgal.Rect(l.Bounds.Min.X, y, l.Bounds.Max.X, y+l.ItemHeight)

		st := l.Style
		if idx == l.hoverIdx {
			st = st.HoverStyle()
		}
		if idx == l.Selected {
			st.Stroke = 2
		}

		st.DrawBox(s, itemBounds)
		st.Ink(s, itemBounds, text)
	}
}

func (l *ListLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	nw := minTotalW
	for _, item := range l.Items {
		sz := l.Style.MeasureText(item)
		if w := sz.X + l.Style.Margin.X*2; w > nw {
			nw = w
		}
	}
	nh := len(l.Items) * l.ItemHeight
	l.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+nw, bounds.Min.Y+nh)
	return l.Bounds
}
