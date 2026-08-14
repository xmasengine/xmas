package xvec

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func mkcol(r, g, b, a uint8) Color { return Color{R: r, G: g, B: b, A: a} }

func equalColor(a, b Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && a.A == b.A
}

// roundtrip encodes x and decodes into a new XVEC.
func roundtrip(t *testing.T, x *XVEC) *XVEC {
	t.Helper()
	var buf bytes.Buffer
	if err := x.Encode(&buf); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	var x2 XVEC
	if err := x2.Decode(&buf); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	return &x2
}

func TestEmpty(t *testing.T) {
	x := &XVEC{Size: Size{100, 50}, Antialias: false}
	x2 := roundtrip(t, x)
	if x2.Size.W != 100 || x2.Size.H != 50 {
		t.Fatalf("size: got %v %v, want 100 50", x2.Size.W, x2.Size.H)
	}
	if x2.Antialias != false {
		t.Fatalf("antialias: got %v, want false", x2.Antialias)
	}
	if len(x2.Instructions) != 0 {
		t.Fatalf("instructions: got %d, want 0", len(x2.Instructions))
	}
}

func TestCircle(t *testing.T) {
	x := &XVEC{Size: Size{320, 240}, Antialias: true}
	x.Circle(160, 120, 50, 3, mkcol(255, 0, 0, 255))

	x2 := roundtrip(t, x)
	if len(x2.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x2.Instructions))
	}
	c, ok := x2.Instructions[0].(*CircleInstruction)
	if !ok {
		t.Fatalf("instruction type: got %T, want *CircleInstruction", x2.Instructions[0])
	}
	if c.C.X != 160 || c.C.Y != 120 {
		t.Fatalf("center: got %v, want (160,120)", c.C)
	}
	if float32(c.R) != 50 {
		t.Fatalf("radius: got %v, want 50", c.R)
	}
}

func TestDisk(t *testing.T) {
	x := &XVEC{Size: Size{80, 60}}
	x.Disk(40, 30, 20, mkcol(0, 255, 0, 128))

	x2 := roundtrip(t, x)
	d, ok := x2.Instructions[0].(*DiskInstruction)
	if !ok {
		t.Fatalf("got %T, want *DiskInstruction", x2.Instructions[0])
	}
	if d.C.X != 40 || d.C.Y != 30 {
		t.Fatalf("center: got %v, want (40,30)", d.C)
	}
	if float32(d.R) != 20 {
		t.Fatalf("radius: got %v, want 20", d.R)
	}
	if !equalColor(d.Color, mkcol(0, 255, 0, 128)) {
		t.Fatalf("color: got %v", d.Color)
	}
}

func TestRect(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}, Antialias: true}
	x.Rect(10, 20, 50, 30, 2, mkcol(128, 128, 128, 255))

	x2 := roundtrip(t, x)
	r, ok := x2.Instructions[0].(*RectInstruction)
	if !ok {
		t.Fatalf("got %T, want *RectInstruction", x2.Instructions[0])
	}
	if r.X != 10 || r.Y != 20 || r.W != 50 || r.H != 30 {
		t.Fatalf("rect: got %v,%v,%v,%v, want 10,20,50,30", r.X, r.Y, r.W, r.H)
	}
	if float32(r.Stroke) != 2 {
		t.Fatalf("stroke: got %v, want 2", r.Stroke)
	}
}

func TestSlab(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}}
	x.Slab(0, 0, 64, 64, mkcol(255, 255, 255, 255))

	x2 := roundtrip(t, x)
	fr, ok := x2.Instructions[0].(*SlabInstruction)
	if !ok {
		t.Fatalf("got %T, want *SlabInstruction", x2.Instructions[0])
	}
	if fr.X != 0 || fr.Y != 0 || fr.W != 64 || fr.H != 64 {
		t.Fatalf("slab: got %v,%v,%v,%v", fr.X, fr.Y, fr.W, fr.H)
	}
}

func TestLine(t *testing.T) {
	x := &XVEC{Size: Size{200, 200}}
	x.Line(10, 10, 190, 190, 1, mkcol(0, 0, 255, 255))

	x2 := roundtrip(t, x)
	l, ok := x2.Instructions[0].(*LineInstruction)
	if !ok {
		t.Fatalf("got %T, want *LineInstruction", x2.Instructions[0])
	}
	if l.X1 != 10 || l.Y1 != 10 || l.X2 != 190 || l.Y2 != 190 {
		t.Fatalf("line: got %v,%v,%v,%v", l.X1, l.Y1, l.X2, l.Y2)
	}
}

func TestFillPath(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	steps := []Stepper{
		MoveTo(10, 10),
		LineTo(90, 10),
		LineTo(50, 80),
		Close(),
	}
	x.Fill(mkcol(255, 0, 0, 200), steps...)

	x2 := roundtrip(t, x)
	f, ok := x2.Instructions[0].(*FillInstruction)
	if !ok {
		t.Fatalf("got %T, want *FillInstruction", x2.Instructions[0])
	}
	if !equalColor(f.Color, mkcol(255, 0, 0, 200)) {
		t.Fatalf("color: got %v", f.Color)
	}
	if len(f.Steps) != 4 {
		t.Fatalf("steps: got %d, want 4", len(f.Steps))
	}
}

func TestStrokePath(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	steps := []Stepper{
		MoveTo(10, 50),
		QuadTo(50, 10, 90, 50),
		CubicTo(90, 90, 50, 90, 10, 50),
		Close(),
	}
	x.Stroke(3, mkcol(0, 255, 0, 255), steps...)

	x2 := roundtrip(t, x)
	s, ok := x2.Instructions[0].(*StrokeInstruction)
	if !ok {
		t.Fatalf("got %T, want *StrokeInstruction", x2.Instructions[0])
	}
	if float32(s.Stroke) != 3 {
		t.Fatalf("stroke: got %v, want 3", s.Stroke)
	}
	if len(s.Steps) != 4 {
		t.Fatalf("steps: got %d, want 4", len(s.Steps))
	}
}

func TestArc(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	steps := []Stepper{
		MoveTo(50, 10),
		Arc(50, 50, 40, 0, 3.14159, Clockwise),
		ArcTo(50, 90, 10, 50, 20),
		Close(),
	}
	x.Stroke(1, mkcol(0, 0, 0, 255), steps...)

	x2 := roundtrip(t, x)
	s, ok := x2.Instructions[0].(*StrokeInstruction)
	if !ok {
		t.Fatalf("got %T, want *StrokeInstruction", x2.Instructions[0])
	}
	if len(s.Steps) != 4 {
		t.Fatalf("got %d steps, want 4", len(s.Steps))
	}
	_, ok = s.Steps[1].(*ArcStep)
	if !ok {
		t.Fatalf("step[1]: got %T, want *ArcStep", s.Steps[1])
	}
	_, ok = s.Steps[2].(*ArcToStep)
	if !ok {
		t.Fatalf("step[2]: got %T, want *ArcToStep", s.Steps[2])
	}
}

func TestAllPrimitives(t *testing.T) {
	x := &XVEC{Size: Size{400, 300}, Antialias: true}
	red := mkcol(255, 0, 0, 255)
	green := mkcol(0, 255, 0, 255)
	blue := mkcol(0, 0, 255, 255)

	x.Circle(200, 150, 80, 2, red)
	x.Disk(200, 150, 40, green)
	x.Rect(10, 10, 100, 50, 1, blue)
	x.Slab(120, 10, 50, 50, red)
	x.Line(0, 0, 400, 300, 3, green)
	x.Fill(blue, MoveTo(300, 200), LineTo(350, 200), LineTo(325, 250), Close())
	x.Stroke(2, red, MoveTo(50, 200), CubicTo(100, 150, 150, 250, 200, 200), Close())

	x2 := roundtrip(t, x)
	if len(x2.Instructions) != 7 {
		t.Fatalf("got %d instructions, want 7", len(x2.Instructions))
	}
	if x2.Antialias != true {
		t.Fatalf("antialias: got false, want true")
	}
	if x2.Size.W != 400 || x2.Size.H != 300 {
		t.Fatalf("size: got %v,%v, want 400,300", x2.Size.W, x2.Size.H)
	}
}

func TestDecodeRaw(t *testing.T) {
	src := `xvec 1
size 100 100
antialias true
circle 50 50 30 2 #ff0000ff
disk 50 50 10 #00ff00ff
rect 10 10 80 30 1 #0000ffff
slab 20 20 20 20 #808080ff
line 0 0 100 100 1 #ffffffff
fill #c83232ff
  move 10 10
  line 90 10
  quad 90 50 50 90
  close
end
stroke 3 #000000ff
  move 10 50
  cubic 50 10 90 90 50 50
  arc 50 50 30 0 3.14159 C
  arcto 50 90 10 50 20
  close
end
`
	var x XVEC
	if err := x.Decode(strings.NewReader(src)); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(x.Instructions) != 7 {
		t.Fatalf("got %d instructions, want 7", len(x.Instructions))
	}
	// Verify roundtrip preserves the same number of instructions
	var buf bytes.Buffer
	if err := x.Encode(&buf); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	var x2 XVEC
	x2.Decode(&buf)
	if len(x2.Instructions) != 7 {
		t.Fatalf("second roundtrip: got %d instructions, want 7", len(x2.Instructions))
	}
}

func TestDefaults(t *testing.T) {
	src := "xvec 1\nsize 640 480\n"
	var x XVEC
	if err := x.Decode(strings.NewReader(src)); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if x.Size.W != 640 || x.Size.H != 480 {
		t.Fatalf("size: got %v,%v, want 640,480", x.Size.W, x.Size.H)
	}
	if x.Antialias != true {
		t.Fatalf("antialias: got false, want true")
	}
	if len(x.Instructions) != 0 {
		t.Fatalf("instructions: got %d, want 0", len(x.Instructions))
	}
}

func TestRoundtripColor(t *testing.T) {
	colors := []Color{
		mkcol(0, 0, 0, 0),
		mkcol(255, 255, 255, 255),
		mkcol(128, 64, 32, 16),
		mkcol(1, 2, 3, 4),
	}
	for _, c := range colors {
		x := &XVEC{Size: Size{10, 10}}
		x.Disk(5, 5, 2, c)
		x2 := roundtrip(t, x)
		d := x2.Instructions[0].(*DiskInstruction)
		if !equalColor(d.Color, c) {
			t.Fatalf("color roundtrip: got %v, want %v", d.Color, c)
		}
	}
}

func TestMarshalTextInterface(t *testing.T) {
	var inst Instruction = &CircleInstruction{}
	if _, err := inst.MarshalText(); err != nil {
		t.Fatalf("MarshalText via interface: %v", err)
	}
	var step Stepper = &MoveStep{}
	if _, err := step.MarshalText(); err != nil {
		t.Fatalf("MarshalText via interface: %v", err)
	}
}

func TestErrorGarbage(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize abc def"))
	if err == nil {
		t.Fatal("expected error for garbage input")
	}
}

func TestErrorTruncated(t *testing.T) {
	// Missing arguments after a keyword.
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\ncircle 10 10"))
	if err == nil {
		t.Fatal("expected error for truncated input")
	}
}

func TestErrorUnknownKeyword(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\nfoobar 1 2 3\ncircle 5 5 2 1 #ff0000ff"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(x.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x.Instructions))
	}
}

func TestLineComments(t *testing.T) {
	src := "xvec 1\nsize 100 100\n// this is a comment\ncircle 50 50 30 2 #ff0000ff\n// another comment\n"
	var x XVEC
	if err := x.Decode(strings.NewReader(src)); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(x.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x.Instructions))
	}
}

func TestBadVersion(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 2\nsize 10 10"))
	if err == nil {
		t.Fatal("expected error for bad version")
	}
}

func TestBadColorNoHash(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\ncircle 5 5 2 1 ff0000ff"))
	if err == nil {
		t.Fatal("expected error for color without #")
	}
}

func TestFillMissingClose(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\nfill #ff0000ff\n  move 0 0\n  line 10 0\n  line 5 10\nend"))
	if err == nil {
		t.Fatal("expected error for fill missing close")
	}
}

func TestStrokeMissingClose(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\nstroke 1 #ff0000ff\n  move 0 0\n  line 10 0\nend"))
	if err == nil {
		t.Fatal("expected error for stroke missing close")
	}
}

func TestMultiSubPathFill(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	blue := mkcol(0, 0, 255, 255)
	x.Fill(blue,
		MoveTo(10, 10), LineTo(90, 10), LineTo(50, 40), Close(),
		MoveTo(10, 60), LineTo(90, 60), LineTo(50, 90), Close(),
	)

	x2 := roundtrip(t, x)
	if len(x2.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x2.Instructions))
	}
	f, ok := x2.Instructions[0].(*FillInstruction)
	if !ok {
		t.Fatalf("got %T, want *FillInstruction", x2.Instructions[0])
	}
	if len(f.Steps) != 8 {
		t.Fatalf("got %d steps, want 8", len(f.Steps))
	}
}

func TestMultiSubPathStroke(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	red := mkcol(255, 0, 0, 255)
	x.Stroke(1, red,
		MoveTo(10, 10), LineTo(90, 10), Close(),
		MoveTo(10, 50), LineTo(90, 50), Close(),
	)

	x2 := roundtrip(t, x)
	if len(x2.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x2.Instructions))
	}
	s, ok := x2.Instructions[0].(*StrokeInstruction)
	if !ok {
		t.Fatalf("got %T, want *StrokeInstruction", x2.Instructions[0])
	}
	if len(s.Steps) != 6 {
		t.Fatalf("got %d steps, want 6", len(s.Steps))
	}
}

func TestFillEmptyAllowed(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\nfill #ff0000ff\n  end"))
	if err != nil {
		t.Fatalf("unexpected error for empty fill: %v", err)
	}
}

func TestBadColorTooFewHexDigits(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\ncircle 5 5 2 1 #fff"))
	if err == nil {
		t.Fatal("expected error for too few hex digits")
	}
}

func TestBadColorTooManyHexDigits(t *testing.T) {
	var x XVEC
	err := x.Decode(strings.NewReader("xvec 1\nsize 10 10\ncircle 5 5 2 1 #ff0000ffff"))
	if err == nil {
		t.Fatal("expected error for too many hex digits")
	}
}

const testSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
  <rect x="10" y="10" width="80" height="80" fill="#ff0000ff" stroke="#0000ffff" stroke-width="2"/>
  <circle cx="50" cy="50" r="30" fill="none" stroke="#00ff00ff" stroke-width="3"/>
  <path d="M10 10 L90 10 L50 90 Z" fill="#ffff00ff"/>
</svg>`

func TestParseSVG(t *testing.T) {
	x, err := ParseSVG(strings.NewReader(testSVG), 100, 100)
	if err != nil {
		t.Fatalf("ParseSVG: %v", err)
	}
	if x.Size.W != 100 || x.Size.H != 100 {
		t.Fatalf("size: got %v,%v, want 100,100", x.Size.W, x.Size.H)
	}
	// Should have 4 instructions: fill rect, stroke rect, stroke circle, fill path
	if len(x.Instructions) != 4 {
		t.Fatalf("got %d instructions, want 4", len(x.Instructions))
	}
}

func TestParseSVGScale(t *testing.T) {
	x, err := ParseSVG(strings.NewReader(testSVG), 200, 0)
	if err != nil {
		t.Fatalf("ParseSVG: %v", err)
	}
	if x.Size.W != 200 || x.Size.H != 100 {
		t.Fatalf("size: got %v,%v, want 200,100", x.Size.W, x.Size.H)
	}
}

func TestParseSVGArc(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
  <path d="M10 50 A40 30 0 0 1 90 50" stroke="#ff0000ff" stroke-width="2" fill="none"/>
</svg>`
	x, err := ParseSVG(strings.NewReader(svg), 0, 0)
	if err != nil {
		t.Fatalf("ParseSVG: %v", err)
	}
	if len(x.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x.Instructions))
	}
	s, ok := x.Instructions[0].(*StrokeInstruction)
	if !ok {
		t.Fatalf("got %T, want *StrokeInstruction", x.Instructions[0])
	}
	// M10 50 + A arc → Move + at least 1 Cubic = at least 2 steps.
	if len(s.Steps) < 2 {
		t.Fatalf("got %d steps, want >= 2", len(s.Steps))
	}
	// First step should be a Move.
	if _, ok := s.Steps[0].(*MoveStep); !ok {
		t.Fatalf("first step should be MoveStep, got %T", s.Steps[0])
	}
	// At least one cubic bezier step.
	hasCubic := false
	for _, step := range s.Steps {
		if _, ok := step.(*CubicStep); ok {
			hasCubic = true
			break
		}
	}
	if !hasCubic {
		t.Fatal("expected at least one CubicStep in arc conversion")
	}
}

func TestMoveUp(t *testing.T) {
	x := &XVEC{}
	x.Disk(10, 10, 5, mkcol(1, 0, 0, 255))
	x.Disk(20, 20, 5, mkcol(2, 0, 0, 255))
	x.Disk(30, 30, 5, mkcol(3, 0, 0, 255))

	x.MoveUp(1)
	if x.Instructions[0].(*DiskInstruction).Color.R != 2 ||
		x.Instructions[1].(*DiskInstruction).Color.R != 1 {
		t.Fatal("MoveUp(1): expected 2, 1, 3")
	}

	x.MoveUp(1)
	if x.Instructions[0].(*DiskInstruction).Color.R != 1 ||
		x.Instructions[1].(*DiskInstruction).Color.R != 2 ||
		x.Instructions[2].(*DiskInstruction).Color.R != 3 {
		t.Fatal("MoveUp(1) again: expected 1, 2, 3")
	}
}

func TestMoveDown(t *testing.T) {
	x := &XVEC{}
	x.Disk(10, 10, 5, mkcol(1, 0, 0, 255))
	x.Disk(20, 20, 5, mkcol(2, 0, 0, 255))
	x.Disk(30, 30, 5, mkcol(3, 0, 0, 255))

	x.MoveDown(1)
	if x.Instructions[0].(*DiskInstruction).Color.R != 1 ||
		x.Instructions[1].(*DiskInstruction).Color.R != 3 ||
		x.Instructions[2].(*DiskInstruction).Color.R != 2 {
		t.Fatal("MoveDown(1): expected 1, 3, 2")
	}
}

func TestMoveToFront(t *testing.T) {
	x := &XVEC{}
	x.Disk(10, 10, 5, mkcol(1, 0, 0, 255))
	x.Disk(20, 20, 5, mkcol(2, 0, 0, 255))
	x.Disk(30, 30, 5, mkcol(3, 0, 0, 255))

	x.MoveToFront(2)
	if x.Instructions[0].(*DiskInstruction).Color.R != 3 ||
		x.Instructions[1].(*DiskInstruction).Color.R != 1 ||
		x.Instructions[2].(*DiskInstruction).Color.R != 2 {
		t.Fatal("MoveToFront(2): expected 3, 1, 2")
	}

	x.MoveToFront(1)
	if x.Instructions[0].(*DiskInstruction).Color.R != 1 ||
		x.Instructions[1].(*DiskInstruction).Color.R != 3 ||
		x.Instructions[2].(*DiskInstruction).Color.R != 2 {
		t.Fatal("MoveToFront(1) then: expected 1, 3, 2")
	}
}

func TestMoveToBack(t *testing.T) {
	x := &XVEC{}
	x.Disk(10, 10, 5, mkcol(1, 0, 0, 255))
	x.Disk(20, 20, 5, mkcol(2, 0, 0, 255))
	x.Disk(30, 30, 5, mkcol(3, 0, 0, 255))

	x.MoveToBack(0)
	if x.Instructions[0].(*DiskInstruction).Color.R != 2 ||
		x.Instructions[1].(*DiskInstruction).Color.R != 3 ||
		x.Instructions[2].(*DiskInstruction).Color.R != 1 {
		t.Fatal("MoveToBack(0): expected 2, 3, 1")
	}
}

func TestMoveEdgeCases(t *testing.T) {
	x := &XVEC{} // empty
	x.MoveUp(0)
	x.MoveDown(0)
	x.MoveToFront(0)
	x.MoveToBack(0)
	// no panic, no change

	x.Disk(10, 10, 5, mkcol(1, 0, 0, 255)) // single instruction
	x.MoveUp(0)
	x.MoveDown(0)
	x.MoveToFront(0)
	x.MoveToBack(0)
	if len(x.Instructions) != 1 {
		t.Fatal("single instruction should remain")
	}

	// out-of-range indices
	x.MoveUp(-1)
	x.MoveDown(5)
	x.MoveToFront(-1)
	x.MoveToBack(5)
	if len(x.Instructions) != 1 {
		t.Fatal("out-of-range indices should be no-ops")
	}
}

func TestDrawTransformed(t *testing.T) {
	x := &XVEC{Size: Size{10, 10}, Antialias: false}
	x.Slab(0, 0, 10, 10, mkcol(255, 0, 0, 255))

	dst := ebiten.NewImage(40, 40)
	var g ebiten.GeoM
	g.Scale(2, 2)
	g.Translate(5, 5)
	x.DrawTransformed(dst, g)
	if x.cached == nil {
		t.Fatal("DrawTransformed did not cache the rasterization")
	}
	if w, h := x.cached.Bounds().Dx(), x.cached.Bounds().Dy(); w != 10 || h != 10 {
		t.Fatalf("cached rasterization is %dx%d, want 10x10", w, h)
	}
	first := x.cached
	x.DrawScaled(dst, V(1, 1), Size{20, 20})
	if x.cached != first {
		t.Fatal("DrawScaled re-rasterized instead of reusing the cache")
	}

	var r ebiten.GeoM
	r.Translate(-5, -5)
	r.Rotate(math.Pi / 2)
	r.Translate(5, 5)
	x.DrawTransformed(dst, r) // rotation must not panic
}

func TestDefUse(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}}
	d := x.Def("sprite")
	d.Slab(0, 0, 16, 16, mkcol(255, 0, 0, 255))
	u := x.Use("sprite", V(10, 10), Size{2, 2}, math.Pi/2)
	if u.Vec != d {
		t.Fatal("Use did not resolve to the def")
	}

	x2 := roundtrip(t, x)
	if len(x2.Defs) != 1 {
		t.Fatalf("got %d defs, want 1", len(x2.Defs))
	}
	if x2.Defs[0].Name != "sprite" {
		t.Fatalf("def name: got %q, want sprite", x2.Defs[0].Name)
	}
	if len(x2.Defs[0].Vec.Instructions) != 1 {
		t.Fatalf("def instructions: got %d, want 1", len(x2.Defs[0].Vec.Instructions))
	}
	if len(x2.Instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(x2.Instructions))
	}
	u2, ok := x2.Instructions[0].(*UseInstruction)
	if !ok {
		t.Fatalf("got %T, want *UseInstruction", x2.Instructions[0])
	}
	if u2.Name != "sprite" {
		t.Fatalf("use name: got %q", u2.Name)
	}
	if u2.At.X != 10 || u2.At.Y != 10 {
		t.Fatalf("use at: got %v", u2.At)
	}
	if u2.Scale.W != 2 || u2.Scale.H != 2 {
		t.Fatalf("use scale: got %v", u2.Scale)
	}
	if math.Abs(float64(u2.Rotate)-math.Pi/2) > 1e-6 {
		t.Fatalf("use rotate: got %v, want %v", u2.Rotate, math.Pi/2)
	}
	if u2.Vec == nil {
		t.Fatal("use did not resolve its def after roundtrip")
	}
}

func TestUseNoClauses(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}}
	x.Def("sprite").Disk(32, 32, 10, mkcol(0, 255, 0, 255))
	x.Use("sprite", Vertex{}, Size{}, 0)

	x2 := roundtrip(t, x)
	u, ok := x2.Instructions[0].(*UseInstruction)
	if !ok {
		t.Fatalf("got %T, want *UseInstruction", x2.Instructions[0])
	}
	if u.At != (Vertex{}) || u.Scale != (Size{}) || u.Rotate != 0 {
		t.Fatalf("expected zero clauses, got at=%v scale=%v rotate=%v", u.At, u.Scale, u.Rotate)
	}
}

func TestDefUseRaw(t *testing.T) {
	src := `xvec 1
size 100 100
antialias true
def sprite
  slab 0 0 16 16 #ff0000ff
  disk 8 8 4 #00ff00ff
end
use sprite at 20 30
use sprite at 0 0 scale 2 2 rotate 45
`
	var x XVEC
	if err := x.Decode(strings.NewReader(src)); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(x.Defs) != 1 {
		t.Fatalf("got %d defs, want 1", len(x.Defs))
	}
	if len(x.Defs[0].Vec.Instructions) != 2 {
		t.Fatalf("def instructions: got %d, want 2", len(x.Defs[0].Vec.Instructions))
	}
	if len(x.Instructions) != 2 {
		t.Fatalf("got %d instructions, want 2", len(x.Instructions))
	}
	u0, ok := x.Instructions[0].(*UseInstruction)
	if !ok {
		t.Fatalf("got %T, want *UseInstruction", x.Instructions[0])
	}
	if u0.At.X != 20 || u0.At.Y != 30 || u0.Scale != (Size{}) || u0.Rotate != 0 {
		t.Fatalf("use 0: at=%v scale=%v rotate=%v", u0.At, u0.Scale, u0.Rotate)
	}
	u1 := x.Instructions[1].(*UseInstruction)
	if u1.Scale.W != 2 || u1.Scale.H != 2 {
		t.Fatalf("use 1 scale: got %v", u1.Scale)
	}
	if math.Abs(float64(u1.Rotate)-math.Pi/4) > 1e-6 {
		t.Fatalf("use 1 rotate: got %v, want %v", u1.Rotate, math.Pi/4)
	}
}

func TestDefErrors(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{"undefined use", "xvec 1\nsize 10 10\nuse ghost\n"},
		{"nested def", "xvec 1\nsize 10 10\ndef a\ndef b\nend\nend\n"},
		{"empty def", "xvec 1\nsize 10 10\ndef a\nend\n"},
		{"bare end", "xvec 1\nsize 10 10\nend\n"},
		{"at without use", "xvec 1\nsize 10 10\nat 1 2\n"},
		{"scale without use", "xvec 1\nsize 10 10\nscale 2 2\n"},
		{"rotate without use", "xvec 1\nsize 10 10\nrotate 45\n"},
		{"zero scale x", "xvec 1\nsize 10 10\ndef a\nslab 0 0 1 1 #ff0000ff\nend\nuse a scale 0 2\n"},
		{"zero scale y", "xvec 1\nsize 10 10\ndef a\nslab 0 0 1 1 #ff0000ff\nend\nuse a scale 2 0\n"},
		{"size after instruction", "xvec 1\ncircle 5 5 2 1 #ff0000ff\nsize 10 10\n"},
		{"antialias after instruction", "xvec 1\nsize 10 10\ndisk 5 5 2 #ff0000ff\nantialias false\n"},
		{"size after def", "xvec 1\ndef a\nslab 0 0 1 1 #ff0000ff\nend\nsize 10 10\n"},
		{"size after use", "xvec 1\nsize 10 10\ndef a\nslab 0 0 1 1 #ff0000ff\nend\nuse a\nsize 20 20\n"},
		{"size in def after instruction", "xvec 1\nsize 10 10\ndef a\nslab 0 0 1 1 #ff0000ff\nsize 5 5\nend\n"},
		{"antialias in def after instruction", "xvec 1\nsize 10 10\ndef a\nslab 0 0 1 1 #ff0000ff\nantialias false\nend\n"},
		{"size inside path", "xvec 1\nfill #ff0000ff\n  size 5 5\nend\n"},
	}
	for _, c := range cases {
		var x XVEC
		if err := x.Decode(strings.NewReader(c.src)); err == nil {
			t.Fatalf("%s: expected error", c.name)
		}
	}
}

func TestUseDraw(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}}
	x.Def("sprite").Slab(0, 0, 16, 16, mkcol(255, 0, 0, 255))
	x.Use("sprite", V(10, 10), Size{}, 0)

	dst := ebiten.NewImage(100, 100)
	x.Draw(dst) // must not panic
}

func TestUseDrawTransform(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}}
	x.Def("sprite").Disk(32, 32, 8, mkcol(255, 0, 0, 255))
	x.Use("sprite", V(0, 0), Size{2, 2}, math.Pi/2)
	x.Use("sprite", V(5, 5), Size{}, 0)

	dst := ebiten.NewImage(200, 200)
	x.Draw(dst) // must not panic
}

func TestDefRecursion(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}}
	d := x.Def("loop")
	d.Use("loop", V(0, 0), Size{}, 0)
	dst := ebiten.NewImage(64, 64)
	x.Draw(dst) // must terminate without panic
}

func TestDefOwnSettings(t *testing.T) {
	x := &XVEC{Size: Size{64, 64}, Antialias: true}
	d := x.Def("sprite")
	d.Size = Size{16, 16}
	d.Antialias = false
	d.Slab(0, 0, 16, 16, mkcol(255, 0, 0, 255))
	x.Use("sprite", V(0, 0), Size{}, 0)

	x2 := roundtrip(t, x)
	if len(x2.Defs) != 1 {
		t.Fatalf("got %d defs, want 1", len(x2.Defs))
	}
	if x2.Defs[0].Vec.Size.W != 16 || x2.Defs[0].Vec.Size.H != 16 {
		t.Fatalf("def size: got %v, want 16x16", x2.Defs[0].Vec.Size)
	}
	if x2.Defs[0].Vec.Antialias != false {
		t.Fatalf("def antialias: got %v, want false", x2.Defs[0].Vec.Antialias)
	}
	if x2.Size.W != 64 || x2.Size.H != 64 || x2.Antialias != true {
		t.Fatalf("root settings changed: size=%v antialias=%v", x2.Size, x2.Antialias)
	}
}

func TestDefInheritSettings(t *testing.T) {
	src := `xvec 1
size 100 100
antialias false
def spr
  disk 50 50 10 #ff0000ff
end
`
	var x XVEC
	if err := x.Decode(strings.NewReader(src)); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	d := x.Defs[0].Vec
	if d.Size.W != 100 || d.Size.H != 100 {
		t.Fatalf("def inherited size: got %v, want 100x100", d.Size)
	}
	if d.Antialias != false {
		t.Fatalf("def inherited antialias: got %v, want false", d.Antialias)
	}
	if disk, ok := d.Instructions[0].(*DiskInstruction); !ok || disk.Antialias != false {
		t.Fatalf("def instruction should inherit antialias=false")
	}
}

func TestDefOwnSettingsRaw(t *testing.T) {
	src := `xvec 1
size 100 100
antialias true
def spr
  size 20 20
  antialias false
  slab 0 0 20 20 #ff0000ff
end
use spr at 0 0
`
	var x XVEC
	if err := x.Decode(strings.NewReader(src)); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	d := x.Defs[0].Vec
	if d.Size.W != 20 || d.Size.H != 20 {
		t.Fatalf("def size: got %v, want 20x20", d.Size)
	}
	if d.Antialias != false {
		t.Fatalf("def antialias: got %v, want false", d.Antialias)
	}
	if x.Size.W != 100 || x.Size.H != 100 || x.Antialias != true {
		t.Fatalf("root settings changed: size=%v antialias=%v", x.Size, x.Antialias)
	}
	if slab, ok := d.Instructions[0].(*SlabInstruction); !ok || slab.Antialias != false {
		t.Fatalf("def instruction should use the def's antialias=false")
	}

	// roundtrip preserves the per-def settings
	x2 := roundtrip(t, &x)
	if x2.Defs[0].Vec.Size.W != 20 || x2.Defs[0].Vec.Size.H != 20 {
		t.Fatalf("roundtrip def size: got %v, want 20x20", x2.Defs[0].Vec.Size)
	}
	if x2.Defs[0].Vec.Antialias != false {
		t.Fatalf("roundtrip def antialias: got %v, want false", x2.Defs[0].Vec.Antialias)
	}
}
