package xdat

import "testing"
import "bytes"

import "github.com/d4l3k/messagediff"

func TestRoundTrip(t *testing.T) {
	expect := NewZone("town")
	buf := &bytes.Buffer{}
	err := expect.SaveTo(buf)
	if err != nil {
		t.Fatalf("write error %s", err)
	}
	buf2 := bytes.NewBuffer(buf.Bytes())
	observe, err := LoadFrom(buf2)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}

	if diff, ok := messagediff.PrettyDiff(expect, observe); !ok {
		t.Fatalf("\ndiff: %s\n", diff)
	}
}

func TestSaveTo(t *testing.T) {
	zone := NewZone("town")
	buf := &bytes.Buffer{}
	err := zone.SaveTo(buf)
	if err != nil {
		t.Fatalf("write error %s", err)
	}
	// We check the contents with TestRoundTrip

}
