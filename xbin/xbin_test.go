package xbin

import "testing"
import "bytes"
import "reflect"

func TestTreeString(t *testing.T) {
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
	var obs Tree
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

	var obs Tree
	err = obs.Decode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if !reflect.DeepEqual(exp, obs) {
		t.Fatalf("\nexp: %s\nobs: %s\n%#v\n%#v\n", exp, obs, exp, obs)
	}
}

func TestEncodeData(t *testing.T) {
	type ted struct {
		X uint16
		Y uint16
		T uint8
	}
	te := ted{X: 1, Y: 2, T: 't'}

	block := Make("hello", nil)
	block.Add("foo", []byte("bar"))
	err := block.EncodeData(te)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	buf := &bytes.Buffer{}
	err = block.Encode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	exp := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x2, 0x74, 0x66, 0x6f, 0x6f, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x62, 0x61, 0x72}
	obs := buf.Bytes()
	if !reflect.DeepEqual(exp, obs) {
		t.Fatalf("\nexp: %s\nobs: %s\n%#v\n%#v", exp, obs, []byte(exp), []byte(obs))
	}
}

func TestDecodeData(t *testing.T) {
	type ted struct {
		X uint16
		Y uint16
		T uint8
	}
	exp := ted{X: 1, Y: 2, T: 't'}
	buf := bytes.NewBuffer([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x2, 0x74, 0x66, 0x6f, 0x6f, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x62, 0x61, 0x72})

	var obs ted

	var block Tree
	err := block.Decode(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	err = block.DecodeData(&obs)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if !reflect.DeepEqual(exp, obs) {
		t.Fatalf("\nexp: %#v\nobs: %#v\n", exp, obs)
	}
}

func TestFindID(t *testing.T) {
	subsub1 := Make("hello", []byte("bad"))
	subsub2 := Make("world", []byte("bad"))
	sub1 := Make("world", []byte("bad"), subsub1, subsub2)
	sub2 := Make("hello", []byte("ok"))
	block := Make("/", []byte("bad"), sub1, sub2)

	id := MakeID("hello")
	obsb, ok := block.FindID(id)
	if !ok {
		t.Fatalf("not found: %s", id)
	}
	exp := "ok"
	obs := string(obsb.Data)

	if exp != obs {
		t.Fatalf("\nexp: %#v\nobs: %#v\n", exp, obs)
	}
}
