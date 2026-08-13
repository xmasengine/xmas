package xlui

import (
	"testing"

	"github.com/xmasengine/xmas/xgal"
)

func testLayer() *Layer {
	return NewLayer(xgal.Rect(0, 0, 50, 50))
}

// TestCloseLayerByIndexClearsDragged ensures the removed layer is no longer
// dragged.
func TestCloseLayerByIndexClearsDragged(t *testing.T) {
	u := &UI{}
	l := testLayer()
	u.Append(l)
	u.Dragged = l
	u.closeLayerByIndex(0)
	if u.Dragged != nil {
		t.Fatal("closeLayerByIndex left u.Dragged pointing at a removed layer")
	}
	if len(u.Layers) != 0 {
		t.Fatalf("layers = %d, want 0", len(u.Layers))
	}
}

// TestCloseLayerByIndexKeepsOtherDragged ensures closing a layer does not
// clear the drag of a different layer.
func TestCloseLayerByIndexKeepsOtherDragged(t *testing.T) {
	u := &UI{}
	keep := testLayer()
	drag := testLayer()
	u.Append(keep)
	u.Append(drag)
	u.Dragged = drag
	u.closeLayerByIndex(0) // closes keep
	if u.Dragged != drag {
		t.Fatal("closeLayerByIndex cleared u.Dragged of a different layer")
	}
	if len(u.Layers) != 1 {
		t.Fatalf("layers = %d, want 1", len(u.Layers))
	}
}
