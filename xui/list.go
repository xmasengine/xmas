package xui

import "github.com/xmasengine/xmas/xgal"

const ListItemHeight = 14

// ListLayer is a vertical list of selectable text items.
type ListLayer struct {
	Layer
	Items      []string
	OnSelect   func(index int)
	Selected   int // -1 for none
	ItemHeight int
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
	if l.OnSelect != nil {
		l.OnSelect(i)
	}
}

// AddList is a helper to add a [ListLayer] to a [Layer].
func (m *Layer) AddList(bounds xgal.Rectangle) *ListLayer {
	l := List(bounds)
	m.Add(l)
	return l
}

func (l *ListLayer) Poll() Reply {
	pos := xgal.Mouse()
	l.hoverIdx = -1

	for i := range l.Items {
		y := l.Bounds.Min.Y + i*l.ItemHeight
		itemBounds := xgal.Rect(l.Bounds.Min.X, y, l.Bounds.Max.X, y+l.ItemHeight)
		if pos.In(itemBounds) {
			l.hoverIdx = i
			if xgal.Click(xgal.MouseButtonLeft) {
				l.SelectItem(i)
				return Accept
			}
			return Accept
		}
	}
	return Ignore
}

func (l *ListLayer) Render(s *xgal.Surface) {
	l.Layer.Render(s)

	for i, text := range l.Items {
		y := l.Bounds.Min.Y + i*l.ItemHeight
		itemBounds := xgal.Rect(l.Bounds.Min.X, y, l.Bounds.Max.X, y+l.ItemHeight)

		st := l.Style
		if i == l.hoverIdx {
			st = st.HoverStyle()
		}
		if i == l.Selected {
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
