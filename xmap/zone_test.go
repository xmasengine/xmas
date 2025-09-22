package xmap

import "testing"
import "path/filepath"

func TestNewZone(t *testing.T) {
	z := NewZone(12, 15)
	if z.W != 12 {
		t.Fatalf("z.W")
	}
	if z.H != 15 {
		t.Fatalf("z.J")
	}
	if len(z.Layers) != 1 {
		t.Fatalf("z.Layers")
	}

	if len(z.Layers[0].Rows) != 15 {
		t.Fatalf("z.Layers.Rows")
	}

	if len(z.Layers[0].Rows[0].Cells) != 12 {
		t.Fatalf("z.Layers.Rows.Cells")
	}
}

func TestZoneSave(t *testing.T) {
	z := NewZone(12, 15)
	name := filepath.Join(t.TempDir(), "zone.xmas.xml")
	err := z.Save(name)
	t.Logf("%s", name)
	if err != nil {
		t.Fatalf("z.Save, %s", err)
	}
}
