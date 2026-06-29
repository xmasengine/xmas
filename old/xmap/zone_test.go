package xmap

import "os"
import "testing"
import "path/filepath"

func TestNewZone(t *testing.T) {
	z := NewZone("Forest", 12, 15)

	if z.Name != "Forest" {
		t.Fatalf("z.Name")
	}
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

func TestZoneSaveToName(t *testing.T) {
	z := NewZone("Forest", 12, 15)
	name := filepath.Join(t.TempDir(), "zone.xmas.xml")
	err := z.Save(ToName(name))
	t.Logf("%s", name)
	if err != nil {
		t.Fatalf("z.Save, %s", err)
	}
}

func TestZoneSaveToRoot(t *testing.T) {
	z := NewZone("Forest", 12, 15)
	name := "zone.xmas.xml"
	root, err := os.OpenRoot(t.TempDir())
	if err != nil {
		t.Fatalf("os.OpenRoot, %s", err)
	}
	err = z.Save(ToRoot(root, name))
	t.Logf("%s", name)
	if err != nil {
		t.Fatalf("z.Save, %s", err)
	}
}

func TestLoadZone(t *testing.T) {
	z1 := NewZone("Forest", 12, 15)
	name := filepath.Join(t.TempDir(), "zone.xmas.xml")
	err := z1.Save(ToName(name))
	if err != nil {
		t.Fatalf("z.Save, %s", err)
	}
	t.Logf("%s", name)
	z2, err := LoadZone(FromName(name))
	if err != nil {
		t.Fatalf("LoadZone, %s", err)
	}

	if z1.Name != z2.Name {
		t.Fatalf("z.Name")
	}

	if z1.W != z2.W {
		t.Fatalf("z.W")
	}
	if z1.H != z2.H {
		t.Fatalf("z.J")
	}

	if len(z1.Layers) != len(z2.Layers) {
		t.Fatalf("z.Layers")
	}

	if len(z1.Layers[0].Rows) != len(z2.Layers[0].Rows) {
		t.Fatalf("z.Layers.Rows")
	}

	if len(z1.Layers[0].Rows[0].Cells) != len(z1.Layers[0].Rows[0].Cells) {
		t.Fatalf("z.Layers.Rows.Cells")
	}

}
