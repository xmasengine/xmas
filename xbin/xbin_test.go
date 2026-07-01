package xbin

import "testing"
import "bytes"
import "reflect"

func TestBlockString(t *testing.T) {
	b := Make("hello", []byte("world"))
	exp := "<hello>\nworld\n</hello>\n"
	obs := b.String()
	if exp != obs {
		t.Fatalf("\nexp: %s\nobs: %s\n%#v\n%#v", exp, obs, []byte(exp), []byte(obs))
	}
}

func TestEncode(t *testing.T) {
	buf := &bytes.Buffer{}
	b := Make("hello", []byte("world"))
	err := b.Encode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	exp := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x77, 0x6f, 0x72, 0x6c, 0x64}
	obs := buf.Bytes()
	if !reflect.DeepEqual(exp, obs) {
		t.Fatalf("\nexp: %s\nobs: %s\n%#v\n%#v", exp, obs, []byte(exp), []byte(obs))
	}
}

func TestDecode(t *testing.T) {
	data := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x77, 0x6f, 0x72, 0x6c, 0x64}
	buf := bytes.NewBuffer(data)
	var obs Block
	err := obs.Decode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	exp := Make("hello", []byte("world"))
	if !reflect.DeepEqual(exp, obs) {
		t.Fatalf("\nexp: %s\nobs: %s\n%#v\n%#v", exp, obs, exp, obs)
	}
}

func TestRoundTrip(t *testing.T) {
	exp := Make("hello", []byte("world"))
	exp.Add("foo", []byte("bar"))

	buf := &bytes.Buffer{}
	err := exp.Encode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	t.Logf("%#v\n", buf.Bytes())

	var obs Block
	err = obs.Decode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if !reflect.DeepEqual(exp, obs) {
		t.Fatalf("\nexp: %s\nobs: %s\n%#v\n%#v\n", exp, obs, exp, obs)
	}
}
