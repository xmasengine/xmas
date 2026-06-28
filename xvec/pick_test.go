package xvec

import "testing"

func TestPickDisk(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Disk(50, 50, 20, mkcol(255, 0, 0, 255))

	if got := Pick(x, 50, 50); got == nil {
		t.Fatal("Pick at center: got nil, want *DiskInstruction")
	}
	if got := Pick(x, 60, 50); got == nil {
		t.Fatal("Pick inside disk: got nil")
	}
	if got := Pick(x, 71, 50); got != nil {
		t.Fatal("Pick outside disk: got instruction, want nil")
	}
	if got := Pick(x, 70, 50); got == nil {
		t.Fatal("Pick on edge of disk: got nil, want instruction")
	}
}

func TestPickCircle(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Circle(50, 50, 30, 4, mkcol(255, 0, 0, 255))

	if Pick(x, 50, 20) == nil {
		t.Fatal("Pick on circle stroke top: got nil, want instruction")
	}
	if Pick(x, 50, 15) != nil {
		t.Fatal("Pick inside circle hole: got instruction, want nil")
	}
	if Pick(x, 50, 85) != nil {
		t.Fatal("Pick far outside circle: got instruction, want nil")
	}
	if Pick(x, 50, 50) != nil {
		t.Fatal("Pick at circle center (hole): got instruction, want nil")
	}
	if Pick(x, 80, 50) == nil {
		t.Fatal("Pick on circle stroke right edge: got nil, want instruction")
	}
}

func TestPickSlab(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Slab(10, 20, 50, 30, mkcol(0, 255, 0, 255))

	if Pick(x, 35, 35) == nil {
		t.Fatal("Pick inside slab: got nil, want instruction")
	}
	if Pick(x, 10, 20) == nil {
		t.Fatal("Pick at slab corner: got nil, want instruction")
	}
	if Pick(x, 61, 50) != nil {
		t.Fatal("Pick outside slab right: got instruction, want nil")
	}
	if Pick(x, 5, 35) != nil {
		t.Fatal("Pick outside slab left: got instruction, want nil")
	}
}

func TestPickRect(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Rect(10, 10, 60, 40, 4, mkcol(0, 0, 255, 255))

	if Pick(x, 10, 10) == nil {
		t.Fatal("Pick on rect corner: got nil, want instruction")
	}
	if Pick(x, 40, 10) == nil {
		t.Fatal("Pick on rect top edge: got nil, want instruction")
	}
	if Pick(x, 40, 13) != nil {
		t.Fatal("Pick inside rect hollow: got instruction, want nil")
	}
	if Pick(x, 5, 30) != nil {
		t.Fatal("Pick outside rect: got instruction, want nil")
	}
	// just inside left edge but past the stroke width
	if Pick(x, 13, 30) != nil {
		t.Fatal("Pick inside hollow area just past stroke: got instruction, want nil")
	}
}

func TestPickLine(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Line(10, 10, 90, 90, 6, mkcol(255, 0, 0, 255))

	if Pick(x, 50, 50) == nil {
		t.Fatal("Pick on line midpoint: got nil, want instruction")
	}
	if Pick(x, 10, 10) == nil {
		t.Fatal("Pick at line start: got nil, want instruction")
	}
	if Pick(x, 90, 90) == nil {
		t.Fatal("Pick at line end: got nil, want instruction")
	}
	// Perpendicular distance sqrt(2) from line should be within stroke 6/2=3
	if Pick(x, 51, 49) == nil {
		t.Fatal("Pick near line (1px off): got nil, want instruction")
	}
	if Pick(x, 55, 45) != nil {
		t.Fatal("Pick far from line: got instruction, want nil")
	}
}

func TestPickFillTriangle(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Fill(mkcol(255, 0, 0, 200),
		MoveTo(10, 10),
		LineTo(90, 10),
		LineTo(50, 80),
		Close(),
	)

	if Pick(x, 50, 30) == nil {
		t.Fatal("Pick inside triangle: got nil, want instruction")
	}
	if Pick(x, 5, 5) != nil {
		t.Fatal("Pick outside triangle top-left: got instruction, want nil")
	}
	if Pick(x, 50, 81) != nil {
		t.Fatal("Pick below triangle: got instruction, want nil")
	}
	if Pick(x, 50, 10) == nil {
		t.Fatal("Pick on triangle top edge: got nil, want instruction")
	}
}

func TestPickFillMultiSubpath(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Fill(mkcol(0, 0, 255, 255),
		MoveTo(10, 10), LineTo(40, 10), LineTo(25, 40), Close(),
		MoveTo(60, 60), LineTo(90, 60), LineTo(75, 90), Close(),
	)

	if Pick(x, 25, 25) == nil {
		t.Fatal("Pick inside first triangle: got nil, want instruction")
	}
	if Pick(x, 75, 75) == nil {
		t.Fatal("Pick inside second triangle: got nil, want instruction")
	}
	if Pick(x, 50, 50) != nil {
		t.Fatal("Pick in gap between triangles: got instruction, want nil")
	}
}

func TestPickStrokeQuad(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Stroke(4, mkcol(0, 255, 0, 255),
		MoveTo(10, 50),
		QuadTo(50, 10, 90, 50),
	)

	if Pick(x, 50, 29) == nil {
		t.Fatal("Pick near quad apex: got nil, want instruction")
	}
	if Pick(x, 10, 50) == nil {
		t.Fatal("Pick at quad start: got nil, want instruction")
	}
	if Pick(x, 90, 50) == nil {
		t.Fatal("Pick at quad end: got nil, want instruction")
	}
	if Pick(x, 50, 5) != nil {
		t.Fatal("Pick far from quad: got instruction, want nil")
	}
}

func TestPickStrokeCubic(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Stroke(3, mkcol(0, 0, 255, 255),
		MoveTo(10, 10),
		CubicTo(30, 90, 70, 90, 90, 10),
	)

	if Pick(x, 50, 69) == nil {
		t.Fatal("Pick near cubic curve: got nil, want instruction")
	}
	if Pick(x, 10, 10) == nil {
		t.Fatal("Pick at cubic start: got nil, want instruction")
	}
	if Pick(x, 90, 10) == nil {
		t.Fatal("Pick at cubic end: got nil, want instruction")
	}
	if Pick(x, 50, 80) != nil {
		t.Fatal("Pick far from cubic: got instruction, want nil")
	}
}

func TestPickStrokeArc(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Stroke(3, mkcol(255, 0, 0, 255),
		MoveTo(80, 50),
		Arc(50, 50, 30, 0, 3.14159, CounterClockwise),
	)

	if Pick(x, 71, 71) == nil {
		t.Fatal("Pick near arc midpoint: got nil, want instruction")
	}
	if Pick(x, 20, 50) == nil {
		t.Fatal("Pick near arc left endpoint: got nil, want instruction")
	}
}

func TestPickStrokeArcTo(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Stroke(4, mkcol(0, 255, 0, 255),
		MoveTo(10, 50),
		LineTo(50, 10),
		ArcTo(50, 10, 90, 50, 20),
	)

	if Pick(x, 50, 10) == nil {
		t.Fatal("Pick at arcto vertex: got nil, want instruction")
	}
	if Pick(x, 90, 50) == nil {
		t.Fatal("Pick at arcto end: got nil, want instruction")
	}
}

func TestPickStrokeClose(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Stroke(4, mkcol(255, 0, 0, 255),
		MoveTo(10, 10),
		LineTo(90, 10),
		LineTo(90, 90),
		Close(),
	)

	if Pick(x, 50, 10) == nil {
		t.Fatal("Pick on top edge: got nil, want instruction")
	}
	if Pick(x, 10, 10) == nil {
		t.Fatal("Pick at close point (start): got nil, want instruction")
	}
}

func TestPickTopmost(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Slab(0, 0, 100, 100, mkcol(255, 0, 0, 255))
	x.Disk(50, 50, 20, mkcol(0, 255, 0, 255))

	got := Pick(x, 50, 50)
	if got == nil {
		t.Fatal("Pick overlapping: got nil, want *DiskInstruction (topmost)")
	}
	if _, ok := got.(*DiskInstruction); !ok {
		t.Fatalf("Pick overlapping: got %T, want *DiskInstruction (topmost)", got)
	}
}

func TestPickNone(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Disk(10, 10, 5, mkcol(255, 0, 0, 255))

	if got := Pick(x, 90, 90); got != nil {
		t.Fatal("Pick far away: got instruction, want nil")
	}
}

func TestPickEmpty(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	if got := Pick(x, 50, 50); got != nil {
		t.Fatal("Pick empty xvec: got instruction, want nil")
	}
}

func TestStrokeColor(t *testing.T) {
	red := mkcol(255, 0, 0, 255)
	green := mkcol(0, 255, 0, 200)
	blue := mkcol(0, 0, 255, 128)
	white := mkcol(255, 255, 255, 255)

	tests := []struct {
		inst    Instruction
		want    Color
	}{
		{&CircleInstruction{Color: red}, red},
		{&DiskInstruction{Color: green}, green},
		{&RectInstruction{Color: blue}, blue},
		{&SlabInstruction{Color: white}, white},
		{&LineInstruction{Color: red}, red},
		{&FillInstruction{Color: green}, green},
		{&StrokeInstruction{Color: blue}, blue},
	}
	for _, tc := range tests {
		if got := StrokeColor(tc.inst); got != tc.want {
			t.Fatalf("StrokeColor(%T) = %v, want %v", tc.inst, got, tc.want)
		}
	}
}

func TestStrokeWidth(t *testing.T) {
	tests := []struct {
		inst Instruction
		want float32
	}{
		{&CircleInstruction{Stroke: 3}, 3},
		{&RectInstruction{Stroke: 5}, 5},
		{&LineInstruction{Stroke: 2}, 2},
		{&StrokeInstruction{Stroke: 4}, 4},
		{&DiskInstruction{}, 0},
		{&SlabInstruction{}, 0},
		{&FillInstruction{}, 0},
	}
	for _, tc := range tests {
		if got := StrokeWidth(tc.inst); got != tc.want {
			t.Fatalf("StrokeWidth(%T) = %v, want %v", tc.inst, got, tc.want)
		}
	}
}

func TestMoveCircle(t *testing.T) {
	c := &CircleInstruction{C: V(10, 20), R: 5, Stroke: 2}
	Move(c, 3, 7)
	if c.C.X != 13 || c.C.Y != 27 {
		t.Fatalf("CircleInstruction after Move: C = (%v, %v), want (13, 27)", c.C.X, c.C.Y)
	}
	if c.R != 5 || c.Stroke != 2 {
		t.Fatal("Move changed fields other than C")
	}
}

func TestMoveDisk(t *testing.T) {
	d := &DiskInstruction{C: V(10, 20), R: 5}
	Move(d, -3, 7)
	if d.C.X != 7 || d.C.Y != 27 {
		t.Fatalf("DiskInstruction after Move: C = (%v, %v), want (7, 27)", d.C.X, d.C.Y)
	}
}

func TestMoveRect(t *testing.T) {
	r := &RectInstruction{X: 10, Y: 20, W: 30, H: 40, Stroke: 2}
	Move(r, 5, -5)
	if r.X != 15 || r.Y != 15 {
		t.Fatalf("RectInstruction after Move: (X, Y) = (%v, %v), want (15, 15)", r.X, r.Y)
	}
	if r.W != 30 || r.H != 40 || r.Stroke != 2 {
		t.Fatal("Move changed fields other than X, Y")
	}
}

func TestMoveSlab(t *testing.T) {
	s := &SlabInstruction{X: 10, Y: 20, W: 30, H: 40}
	Move(s, 5, -5)
	if s.X != 15 || s.Y != 15 {
		t.Fatalf("SlabInstruction after Move: (X, Y) = (%v, %v), want (15, 15)", s.X, s.Y)
	}
}

func TestMoveLine(t *testing.T) {
	l := &LineInstruction{X1: 0, Y1: 0, X2: 10, Y2: 20, Stroke: 3}
	Move(l, 5, -5)
	if l.X1 != 5 || l.Y1 != -5 || l.X2 != 15 || l.Y2 != 15 {
		t.Fatalf("LineInstruction after Move: (%v,%v)-(%v,%v), want (5,-5)-(15,15)",
			l.X1, l.Y1, l.X2, l.Y2)
	}
}

func TestMoveFillSteps(t *testing.T) {
	f := &FillInstruction{Steps: []Stepper{
		&MoveStep{X: 10, Y: 10},
		&LineStep{X: 90, Y: 10},
		&LineStep{X: 50, Y: 80},
		&CloseStep{},
	}}
	Move(f, 5, 10)
	if f.Steps[0].(*MoveStep).X != 15 || f.Steps[0].(*MoveStep).Y != 20 {
		t.Fatal("MoveStep not moved")
	}
	if f.Steps[1].(*LineStep).X != 95 || f.Steps[1].(*LineStep).Y != 20 {
		t.Fatal("LineStep not moved")
	}
}

func TestMoveStrokeSteps(t *testing.T) {
	s := &StrokeInstruction{Steps: []Stepper{
		&MoveStep{X: 10, Y: 10},
		&QuadStep{X1: 30, Y1: 90, X2: 50, Y2: 10},
		&CubicStep{X1: 50, Y1: 10, X2: 70, Y2: 90, X3: 90, Y3: 10},
		&ArcStep{CX: 50, CY: 50, R: 30, Start: 0, End: 3.14159, Dir: CounterClockwise},
		&ArcToStep{X1: 50, Y1: 10, X2: 90, Y2: 50, R: 20},
		&LineStep{X: 0, Y: 0},
		&CloseStep{},
	}}
	Move(s, 10, -10)
	steps := s.Steps

	m := steps[0].(*MoveStep)
	if m.X != 20 || m.Y != 0 {
		t.Fatalf("MoveStep: got (%v,%v), want (20,0)", m.X, m.Y)
	}
	q := steps[1].(*QuadStep)
	if q.X1 != 40 || q.Y1 != 80 || q.X2 != 60 || q.Y2 != 0 {
		t.Fatalf("QuadStep: got (%v,%v)-(%v,%v), want (40,80)-(60,0)", q.X1, q.Y1, q.X2, q.Y2)
	}
	c := steps[2].(*CubicStep)
	if c.X1 != 60 || c.Y1 != 0 || c.X2 != 80 || c.Y2 != 80 || c.X3 != 100 || c.Y3 != 0 {
		t.Fatalf("CubicStep: got (%v,%v)-(%v,%v)-(%v,%v), want (60,0)-(80,80)-(100,0)",
			c.X1, c.Y1, c.X2, c.Y2, c.X3, c.Y3)
	}
	a := steps[3].(*ArcStep)
	if a.CX != 60 || a.CY != 40 {
		t.Fatalf("ArcStep: got CX,CY = (%v,%v), want (60,40)", a.CX, a.CY)
	}
	if a.R != 30 || a.Start != 0 || a.End != 3.14159 || a.Dir != CounterClockwise {
		t.Fatal("ArcStep changed non-position fields")
	}
	at := steps[4].(*ArcToStep)
	if at.X1 != 60 || at.Y1 != 0 || at.X2 != 100 || at.Y2 != 40 {
		t.Fatalf("ArcToStep: got (%v,%v)-(%v,%v), want (60,0)-(100,40)", at.X1, at.Y1, at.X2, at.Y2)
	}
	l := steps[5].(*LineStep)
	if l.X != 10 || l.Y != -10 {
		t.Fatalf("LineStep: got (%v,%v), want (10,-10)", l.X, l.Y)
	}
}

func TestMoveNil(t *testing.T) {
	Move(nil, 1, 1) // must not panic
}

func TestMoveUnsupported(t *testing.T) {
	Move(&struct{ Instruction }{}, 1, 1) // must not panic
}

func TestMoveThenPick(t *testing.T) {
	x := &XVEC{Size: Size{100, 100}}
	x.Disk(50, 50, 10, mkcol(255, 0, 0, 255))

	got := Pick(x, 55, 50)
	if got == nil {
		t.Fatal("Pick before Move: got nil, want instruction")
	}
	Move(got, 20, 0)
	if Pick(x, 55, 50) != nil {
		t.Fatal("Pick after Move at old position: got instruction, want nil")
	}
	if Pick(x, 75, 50) == nil {
		t.Fatal("Pick after Move at new position: got nil, want instruction")
	}
}
