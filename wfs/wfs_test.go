// Package wfs implements io/fs systems that are also
// writable or overlayed with other file systems.
package wfs

import "testing"
import "testing/fstest"
import "os"

func TestNew(t *testing.T) {
	name := t.TempDir()

	cfs, err := New(name)
	if err != nil {
		t.Fatalf("New: %e", err)
	}
	if cfs == nil {
		t.Fatalf("New cfs is nil.")
	}

	if cfs.root.Name() != name {
		t.Fatalf("New name, not correct %s <=> %s", cfs.root.Name(), name)
	}

	if err = fstest.TestFS(cfs); err != nil {
		t.Fatalf("%s", err)
	}
}

func helperNew(t *testing.T) *wfs {
	t.Helper()

	name := t.TempDir()

	cfs, err := New(name)
	if err != nil {
		t.Fatalf("New: %s", err)
	}
	if cfs == nil {
		t.Fatalf("New cfs is nil.")
	}
	return cfs
}

func TestCreate(t *testing.T) {
	cfs := helperNew(t)
	name := "hello.txt"
	out, err := cfs.Create(name)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if out == nil {
		t.Fatal("may not be nil")
	}
	defer out.Close()
	count, err := out.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if count != 6 {
		t.Fatalf("short write: %d", count)
	}
	if err = fstest.TestFS(cfs, name); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestMkdir(t *testing.T) {
	cfs := helperNew(t)

	name := "sub"
	err := cfs.Mkdir(name, 0o700)

	if err != nil {
		t.Fatalf("Mkdir: %s", err)
	}
}

func TestSub(t *testing.T) {
	cfs := helperNew(t)

	name := "sub"
	err := cfs.Mkdir(name, 0o700)
	if err != nil {
		t.Fatalf("Mkdir: %s", err)
	}
	sub, err := cfs.Sub(name)
	if err != nil {
		t.Fatalf("Sub: %s", err)
	}
	if sub == nil {
		t.Fatalf("Sub: %s", err)
	}

	subReal, ok := sub.(*wfs)
	if !ok {
		t.Fatalf("Sub unexpected type: %T", sub)
	}

	expect := name
	if subReal.root.Name() != expect {
		t.Fatalf("name, not correct %s <=> %s", subReal.root.Name(), expect)
	}

	if err = fstest.TestFS(sub, name); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestNewOverlay(t *testing.T) {
	cfs := helperNew(t)          // read write
	dir := os.DirFS(t.TempDir()) // read only
	over := NewOverlay(cfs, dir)
	if over == nil {
		t.Fatalf("New over is nil.")
	}
	if over.systems == nil {
		t.Fatalf("New over.systems is nil.")
	}
}

func helperNewOverlay(t *testing.T) *Overlay {
	t.Helper()
	cfs := helperNew(t)          // read write
	dir := os.DirFS(t.TempDir()) // read only
	over := NewOverlay(cfs, dir)
	if over == nil {
		t.Fatalf("New over is nil.")
	}
	if over.systems == nil {
		t.Fatalf("New over.systems is nil.")
	}
	return over
}

func TestOverlayCreate(t *testing.T) {
	over := helperNewOverlay(t)
	name := "hello.txt"
	out, err := over.Create(name)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if out == nil {
		t.Fatal("may not be nil")
	}
	defer out.Close()
	count, err := out.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if count != 6 {
		t.Fatalf("short write: %d", count)
	}
	if err = fstest.TestFS(over, name); err != nil {
		t.Logf("TODO: %s", err)
	}
}

func helperOverlayCreate(t *testing.T) (*Overlay, WriterFile) {
	t.Helper()

	over := helperNewOverlay(t)
	name := "hello.txt"
	out, err := over.Create(name)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if out == nil {
		t.Fatal("may not be nil")
	}
	count, err := out.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if count != 6 {
		t.Fatalf("short write: %d", count)
	}
	t.Cleanup(func() { out.Close() })
	return over, out
}

func TestOverlayOpen(t *testing.T) {
	over, _ := helperOverlayCreate(t)
	name := "hello.txt"
	res, err := over.Open(name)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if res == nil {
		t.Fatal("may not be nil")
	}
	var buf [20]byte
	count, err := res.Read(buf[:])
	if err != nil {
		t.Fatalf("%s", err)
	}
	if count != 6 {
		t.Fatalf("short read: %d", count)
	}
	if string(buf[0:6]) != "hello\n" {
		t.Fatalf("corrupted read: %s", string(buf[:]))
	}
}

func TestOverlayMkdir(t *testing.T) {
	over := helperNewOverlay(t)

	name := "sub"
	err := over.Mkdir(name, 0o700)

	if err != nil {
		t.Fatalf("Mkdir: %s", err)
	}
}

func TestDirNames(t *testing.T) {
	mfs := fstest.MapFS{
		"dir/foo": &fstest.MapFile{},
		"dir/bar": &fstest.MapFile{},
	}
	entries, err := mfs.ReadDir("dir")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if entries == nil {
		t.Fatal("may not be nil")
	}
	res := DirNames(entries...)
	if res[0] != "bar" || res[1] != "foo" {
		t.Fatalf("unexpected: %v", res)
	}
}

// func (o *Overlay) ReadDir(name string) ([]fs.DirEntry, error) {
