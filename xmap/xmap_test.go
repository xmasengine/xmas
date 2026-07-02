package xmap

import "testing"
import "bytes"
import "reflect"
import "os"

func TestRoundTrip(t *testing.T) {
	expect := NewZone("town")
	buf := &bytes.Buffer{}
	err := expect.SaveTo(buf)
	if err != nil {
		t.Fatalf("write error %s", err)
	}
	t.Logf("string: %v", buf.String())
	buf2 := bytes.NewBuffer(buf.Bytes())
	observe, err := LoadFrom(buf2)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}
	if !reflect.DeepEqual(*expect, observe) {
		t.Fatalf("not equal: %#v <-> %#v", expect, observe)
	}
}

func TestSaveTo(t *testing.T) {
	expect := NewZone("town")
	buf := &bytes.Buffer{}
	err := expect.SaveTo(buf)
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
