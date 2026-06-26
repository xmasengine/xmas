package xui

import "github.com/xmasengine/xmas/xgal"

// TabItem is a single row inside a [ScreenTab].
type TabItem struct {
	Icon     Icon   // optional icon drawn left of the label
	Label    string // display label
	Value    string // current value shown next to the label
	Activate func() // called on confirm (for usable items)

	// For settings: toggle or multi-choice.  If both are set Options wins.
	Bool    *bool    // if non-nil, toggled by left/right
	Options []string // if non-empty, cycled by left/right; Value is ignored
	OptIdx  int      // current option index
}

// ScreenTab groups related items under one tab header.
type ScreenTab struct {
	Label string
	Items []TabItem
}

// ScreenLayer is a full-screen overlay with tabbed panels for status,
// items, settings, etc.
//
// Navigation:
//   - left/right  – switch tabs, or change item value when an item is selected
//   - up/down     – select item within the active tab
//   - confirm     – activate the selected item (calls [TabItem.Activate])
//   - cancel      – close the screen
type ScreenLayer struct {
	Bounds  xgal.Rectangle
	Style   Style
	Tabs    []ScreenTab
	SelTab  int
	SelItem int // selected item index within the tab, -1 = none

	Close   func() bool
	NextTab func() bool
	PrevTab func() bool
	Confirm func() bool
	Up      func() bool
	Down    func() bool
	Left    func() bool
	Right   func() bool
}

// Screen creates a full-screen menu. Set callback fields on the
// returned [ScreenLayer] to override the global [DefaultInput] bindings.
func Screen(screenW, screenH int, tabs []ScreenTab) *ScreenLayer {
	return &ScreenLayer{
		Bounds:  xgal.Rect(0, 0, screenW, screenH),
		Style:   DefaultStyle(),
		Tabs:    tabs,
		SelItem: 0,
	}
}

var _ Widget = &ScreenLayer{}

func (g *ScreenLayer) Poll() Reply {
	if input(g.Close, DefaultInput.Cancel) {
		return Finish
	}

	if input(g.NextTab, DefaultInput.NextTab) {
		g.SelTab = (g.SelTab + 1) % max(len(g.Tabs), 1)
		g.SelItem = 0
	}
	if input(g.PrevTab, DefaultInput.PrevTab) {
		g.SelTab = (g.SelTab - 1 + len(g.Tabs)) % max(len(g.Tabs), 1)
		g.SelItem = 0
	}

	tab := &g.Tabs[g.SelTab]

	// navigate items
	if input(g.Up, DefaultInput.Up) && len(tab.Items) > 0 {
		g.SelItem = (g.SelItem - 1 + len(tab.Items)) % len(tab.Items)
	}
	if input(g.Down, DefaultInput.Down) && len(tab.Items) > 0 {
		g.SelItem = (g.SelItem + 1) % len(tab.Items)
	}

	// modify selected item (left/right on the same key as PrevTab/NextTab —
	// item change takes priority when an adjustable item is selected).
	if g.SelItem >= 0 && g.SelItem < len(tab.Items) {
		item := &tab.Items[g.SelItem]
		adjustable := len(item.Options) > 0 || item.Bool != nil

		if adjustable {
			if input(g.Left, DefaultInput.Left) {
				if len(item.Options) > 0 {
					item.OptIdx = (item.OptIdx - 1 + len(item.Options)) % len(item.Options)
				} else if item.Bool != nil {
					*item.Bool = !*item.Bool
				}
			}
			if input(g.Right, DefaultInput.Right) {
				if len(item.Options) > 0 {
					item.OptIdx = (item.OptIdx + 1) % len(item.Options)
				} else if item.Bool != nil {
					*item.Bool = !*item.Bool
				}
			}
		}
		if input(g.Confirm, DefaultInput.Confirm) && item.Activate != nil {
			item.Activate()
		}
	}

	return Accept
}

func (g *ScreenLayer) Render(s *xgal.Surface) {
	// dim background
	xgal.Box(s, g.Bounds, xgal.Wash(0, 0, 0, 200))

	// tab bar
	tabH := 24
	tx := g.Bounds.Min.X
	for i, tab := range g.Tabs {
		sz := g.Style.MeasureText(tab.Label)
		tw := sz.X + g.Style.Margin.X*4
		tabBounds := xgal.Rect(tx, g.Bounds.Min.Y, tx+tw, g.Bounds.Min.Y+tabH)
		st := g.Style
		if i == g.SelTab {
			st = st.FocusStyle()
		}
		st.DrawBox(s, tabBounds)
		st.Ink(s, tabBounds, tab.Label)
		tx += tw
	}

	// content area
	contentBounds := xgal.Rect(g.Bounds.Min.X, g.Bounds.Min.Y+tabH,
		g.Bounds.Max.X, g.Bounds.Max.Y)
	g.renderItems(s, contentBounds)
}

func (g *ScreenLayer) renderItems(s *xgal.Surface, bounds xgal.Rectangle) {
	if g.SelTab < 0 || g.SelTab >= len(g.Tabs) {
		return
	}
	tab := &g.Tabs[g.SelTab]

	itemH := 20
	y := bounds.Min.Y + g.Style.Margin.Y
	for i := range tab.Items {
		item := &tab.Items[i]
		row := xgal.Rect(bounds.Min.X+g.Style.Margin.X, y,
			bounds.Max.X-g.Style.Margin.X, y+itemH)

		st := g.Style
		if i == g.SelItem {
			st = st.FocusStyle()
		}
		st.DrawBox(s, row)

		item.Icon.Blit(s, row.Min)
		labelBounds := item.Icon.TextBounds(row)
		st.Ink(s, labelBounds, item.Label)

		// value / toggle / option
		valX := bounds.Max.X - g.Style.Margin.X - 100
		valBounds := xgal.Rect(valX, y, bounds.Max.X-g.Style.Margin.X, y+itemH)
		val := item.Value
		if len(item.Options) > 0 && item.OptIdx < len(item.Options) {
			val = item.Options[item.OptIdx]
		} else if item.Bool != nil {
			if *item.Bool {
				val = "ON"
			} else {
				val = "OFF"
			}
		}
		if val != "" {
			xgal.Ink(s, g.Style.Face, st.Fore, valBounds.Min.X, valBounds.Min.Y, val)
		}

		y += itemH
	}
}

func (g *ScreenLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	g.Bounds = bounds
	return g.Bounds
}

func (g *ScreenLayer) MoveBy(delta xgal.Point) {
	g.Bounds = g.Bounds.Add(delta)
}
