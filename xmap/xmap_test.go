package xmap

import "testing"
import "bytes"
import "reflect"
import "os"

func TestRoundTrip(t *testing.T) {
	expect := Zone{}
	copy(expect.Name.Data[:], []rune("town"))
	expect.Name.Size = 4
	expect.Size = 4
	buf := &bytes.Buffer{}
	err := expect.SaveTo(buf)
	if err != nil {
		t.Fatalf("write error %s", err)
	}
	t.Logf("bytes: %v", buf.Bytes()[:300])
	buf2 := bytes.NewBuffer(buf.Bytes())
	observe, err := LoadFrom(buf2)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}
	if !reflect.DeepEqual(expect.Name, observe.Name) {
		t.Fatalf("not equal: %#v <-> %#v", expect.Name, observe.Name)
	}
	if expect.Size != observe.Size {
		t.Fatalf("not equal: %d <-> %d", expect.Size, observe.Size)
	}
}

func TestSaveTo(t *testing.T) {
	expect := Zone{}
	copy(expect.Name.Data[:], []rune("town"))
	expect.Name.Size = 4
	expect.Size = 4
	out, err := os.Create("test.xmas")
	if err != nil {
		t.Fatalf("file error %s", err)
	}
	defer out.Close()

	err = expect.SaveTo(out)
	if err != nil {
		t.Fatalf("write error %s", err)
	}
}
