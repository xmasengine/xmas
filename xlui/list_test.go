package xlui

import "testing"

import "github.com/xmasengine/xmas/xgal"

func testList() *ListLayer {
	list := NewList(xgal.Rect(0, 0, 100, ListItemHeight*5+10), "a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	return list
}

// TestNewListDefaults ensures a new list has sane defaults.
func TestNewListDefaults(t *testing.T) {
	list := NewList(xgal.Rect(0, 0, 100, 100))
	if list.Selected != -1 {
		t.Errorf("Selected = %d, want -1", list.Selected)
	}
	if list.ItemHeight != ListItemHeight {
		t.Errorf("ItemHeight = %d, want %d", list.ItemHeight, ListItemHeight)
	}
	if list.Offset != 0 {
		t.Errorf("Offset = %d, want 0", list.Offset)
	}
}

// TestListSelectOutOfRange ensures out of range selections are ignored.
func TestListSelectOutOfRange(t *testing.T) {
	list := testList()
	list.Select(-1)
	if list.Selected != -1 {
		t.Errorf("Select(-1): Selected = %d, want -1", list.Selected)
	}
	list.Select(10)
	if list.Selected != -1 {
		t.Errorf("Select(10): Selected = %d, want -1", list.Selected)
	}
}

// TestListSelectCallsValue ensures Value fires on selection.
func TestListSelectCallsValue(t *testing.T) {
	list := testList()
	called := -1
	list.Class.Value = func(i int) Reply { called = i; return Accept }
	list.Select(3)
	if called != 3 {
		t.Errorf("Value = %d, want 3", called)
	}
	if list.Selected != 3 {
		t.Errorf("Selected = %d, want 3", list.Selected)
	}
}

// TestListEnsureVisibleScrolls ensures the selection is scrolled into view.
func TestListEnsureVisibleScrolls(t *testing.T) {
	list := testList() // 10 items, limit 5
	list.Select(7)
	if list.Selected != 7 {
		t.Fatalf("Selected = %d, want 7", list.Selected)
	}
	if list.Offset != 3 {
		t.Errorf("Offset = %d, want 3", list.Offset)
	}
	list.Select(1)
	if list.Offset != 1 {
		t.Errorf("Offset = %d, want 1", list.Offset)
	}
}

// TestListClampOffsetBounds ensures the offset never exceeds the visible range.
func TestListClampOffsetBounds(t *testing.T) {
	list := testList() // 10 items, limit 5, max offset 5
	for _, off := range []int{-3, 8} {
		list.Offset = off
		list.clampOffset()
		if list.Offset < 0 || list.Offset > 5 {
			t.Errorf("clampOffset(%d): Offset = %d, want in [0, 5]", off, list.Offset)
		}
	}
}

// TestListClickSelects ensures clicking an item selects it and raises the
// layer so it goes on top and becomes draggable.
func TestListClickSelects(t *testing.T) {
	list := testList()
	res := list.click(xgal.Pt(50, ListItemHeight+2), int(xgal.MouseButtonLeft))
	if res != Raise {
		t.Errorf("click returned %v, want Raise", res)
	}
	if list.Selected != 1 {
		t.Errorf("Selected = %d, want 1", list.Selected)
	}
}

// TestListClickEmptyArea ensures clicking empty space does not select but is
// still accepted (and raises) so the click is not passed to layers below.
func TestListClickEmptyArea(t *testing.T) {
	list := testList()
	res := list.click(xgal.Pt(50, ListItemHeight*5+8), int(xgal.MouseButtonLeft))
	if res != Raise {
		t.Errorf("click returned %v, want Raise", res)
	}
	if list.Selected != -1 {
		t.Errorf("Selected = %d, want -1", list.Selected)
	}
}

// TestListClickOutside ensures clicks outside the list are ignored.
func TestListClickOutside(t *testing.T) {
	list := testList()
	res := list.click(xgal.Pt(150, 10), int(xgal.MouseButtonLeft))
	if res != Ignore {
		t.Errorf("click returned %v, want Ignore", res)
	}
}

// TestListClickRaisesAndDrags ensures the list behaves like any other layer:
// clicking it makes it draggable and raises it to the top of the layer stack.
func TestListClickRaisesAndDrags(t *testing.T) {
	var ui UI
	base := ui.Layer(xgal.Rect(0, 0, 50, 50))
	list := ui.List(xgal.Rect(0, 0, 100, 100), "a", "b", "c")

	ui.Click(xgal.Pt(5, ListItemHeight+2), int(xgal.MouseButtonLeft))

	if list.Selected != 1 {
		t.Errorf("Selected = %d, want 1", list.Selected)
	}
	if ui.Dragged != &list.Layer {
		t.Fatal("list should be draggable after click")
	}
	if ui.Focused != &list.Layer {
		t.Fatal("list is not focused after click")
	}

	// Add a layer on top, then click the list again; it must move to the top.
	ui.Layer(xgal.Rect(50, 50, 100, 100))
	ui.Click(xgal.Pt(10, 10), int(xgal.MouseButtonLeft))
	if got, want := ui.LayerIndex(&list.Layer), len(ui.Layers)-1; got != want {
		t.Fatalf("list index after re-click = %d, want %d (top)", got, want)
	}
	if got := ui.LayerIndex(base); got != 0 {
		t.Fatalf("base layer index = %d, want 0", got)
	}
}

// TestListWheelScrolls ensures the wheel scrolls and clamps.
func TestListWheelScrolls(t *testing.T) {
	list := testList() // max offset 5
	list.Offset = 2
	res := list.wheel(xgal.Pt(50, 50), 1)
	if res != Accept {
		t.Errorf("wheel returned %v, want Accept", res)
	}
	if list.Offset != 1 {
		t.Errorf("Offset = %d, want 1", list.Offset)
	}
	list.wheel(xgal.Pt(50, 50), 100)
	if list.Offset != 0 {
		t.Errorf("Offset = %d, want 0", list.Offset)
	}
	list.wheel(xgal.Pt(50, 50), -100)
	if list.Offset != 5 {
		t.Errorf("Offset = %d, want 5", list.Offset)
	}
}

// TestListTapArrows ensures arrow keys navigate the selection.
func TestListTapArrows(t *testing.T) {
	list := testList()
	list.Select(4)
	if res := list.tap(int(xgal.KeyArrowDown), Mods{}); res != Accept {
		t.Errorf("arrow down returned %v, want Accept", res)
	}
	if list.Selected != 5 {
		t.Errorf("Selected = %d, want 5", list.Selected)
	}
	if res := list.tap(int(xgal.KeyArrowUp), Mods{}); res != Accept {
		t.Errorf("arrow up returned %v, want Accept", res)
	}
	if list.Selected != 4 {
		t.Errorf("Selected = %d, want 4", list.Selected)
	}
	list.Select(0)
	if res := list.tap(int(xgal.KeyArrowUp), Mods{}); res != Ignore {
		t.Errorf("arrow up at top returned %v, want Ignore", res)
	}
	list.Select(9)
	if res := list.tap(int(xgal.KeyArrowDown), Mods{}); res != Ignore {
		t.Errorf("arrow down at bottom returned %v, want Ignore", res)
	}
}
