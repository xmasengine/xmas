package xvec

import "math"

// Pick returns the topmost (last-drawn) instruction under point (x, y),
// or nil if nothing is under the point.
func Pick(xv *XVEC, x, y float32) Instruction {
	var best Instruction
	for _, inst := range xv.Instructions {
		if hit(inst, x, y) {
			best = inst
		}
	}
	return best
}

func hit(inst Instruction, x, y float32) bool {
	switch v := inst.(type) {
	case *CircleInstruction:
		return hitCircle(v, x, y)
	case *DiskInstruction:
		return hitDisk(v, x, y)
	case *RectInstruction:
		return hitRect(v, x, y)
	case *SlabInstruction:
		return hitSlab(v, x, y)
	case *LineInstruction:
		return hitLine(v, x, y)
	case *FillInstruction:
		return hitFill(v, x, y)
	case *StrokeInstruction:
		return hitStroke(v, x, y)
	}
	return false
}

func hitDisk(d *DiskInstruction, x, y float32) bool {
	dx := float64(x - d.C.X)
	dy := float64(y - d.C.Y)
	r := float64(d.R)
	return dx*dx+dy*dy <= r*r
}

func hitSlab(s *SlabInstruction, x, y float32) bool {
	return x >= s.X && x <= s.X+s.W && y >= s.Y && y <= s.Y+s.H
}

func hitCircle(c *CircleInstruction, x, y float32) bool {
	dx := float64(x - c.C.X)
	dy := float64(y - c.C.Y)
	dist := math.Sqrt(dx*dx + dy*dy)
	half := float64(c.Stroke) / 2
	return math.Abs(dist-float64(c.R)) <= half
}

func hitRect(r *RectInstruction, x, y float32) bool {
	if x < r.X || x > r.X+r.W || y < r.Y || y > r.Y+r.H {
		return false
	}
	half := float64(r.Stroke) / 2
	d := min4(
		float64(x-r.X),
		float64(r.X+r.W-x),
		float64(y-r.Y),
		float64(r.Y+r.H-y),
	)
	return d <= half
}

func hitLine(l *LineInstruction, x, y float32) bool {
	dx := float64(l.X2 - l.X1)
	dy := float64(l.Y2 - l.Y1)
	px := float64(x - l.X1)
	py := float64(y - l.Y1)

	t := (px*dx + py*dy) / (dx*dx + dy*dy)
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	cx := float64(l.X1) + t*dx
	cy := float64(l.Y1) + t*dy

	dist := math.Sqrt((float64(x)-cx)*(float64(x)-cx) + (float64(y)-cy)*(float64(y)-cy))
	return dist <= float64(l.Stroke)/2
}

func hitFill(f *FillInstruction, x, y float32) bool {
	polygons := stepsToPolygons(f.Steps)
	for _, pts := range polygons {
		if pointInPolygon(float64(x), float64(y), pts) {
			return true
		}
	}
	return false
}

func hitStroke(s *StrokeInstruction, x, y float32) bool {
	segs := stepsToSegments(s.Steps)
	half := float64(s.Stroke) / 2
	for _, seg := range segs {
		if pointToSegment(float64(x), float64(y), seg) <= half {
			return true
		}
	}
	return false
}

type fpoint struct{ x, y float64 }

type segment struct{ x1, y1, x2, y2 float64 }

func stepsToSegments(steps []Stepper) []segment {
	var segs []segment
	var cx, cy, sx, sy float64

	for _, s := range steps {
		switch v := s.(type) {
		case *MoveStep:
			cx = float64(v.X)
			cy = float64(v.Y)
			sx = cx
			sy = cy

		case *LineStep:
			segs = append(segs, segment{cx, cy, float64(v.X), float64(v.Y)})
			cx = float64(v.X)
			cy = float64(v.Y)

		case *CloseStep:
			segs = append(segs, segment{cx, cy, sx, sy})
			cx = sx
			cy = sy

		case *QuadStep:
			pts := evalQuad(cx, cy, float64(v.X1), float64(v.Y1), float64(v.X2), float64(v.Y2))
			for i := 1; i < len(pts); i++ {
				segs = append(segs, segment{pts[i-1].x, pts[i-1].y, pts[i].x, pts[i].y})
			}
			cx = float64(v.X2)
			cy = float64(v.Y2)

		case *CubicStep:
			pts := evalCubic(cx, cy, float64(v.X1), float64(v.Y1), float64(v.X2), float64(v.Y2), float64(v.X3), float64(v.Y3))
			for i := 1; i < len(pts); i++ {
				segs = append(segs, segment{pts[i-1].x, pts[i-1].y, pts[i].x, pts[i].y})
			}
			cx = float64(v.X3)
			cy = float64(v.Y3)

		case *ArcStep:
			pts := evalArc(float64(v.CX), float64(v.CY), float64(v.R), float64(v.Start), float64(v.End), v.Dir)
			for i := 1; i < len(pts); i++ {
				segs = append(segs, segment{pts[i-1].x, pts[i-1].y, pts[i].x, pts[i].y})
			}
			cx = pts[len(pts)-1].x
			cy = pts[len(pts)-1].y

		case *ArcToStep:
			// Approximate ArcTo as a straight line from current to end point.
			segs = append(segs, segment{cx, cy, float64(v.X2), float64(v.Y2)})
			cx = float64(v.X2)
			cy = float64(v.Y2)
		}
	}
	return segs
}

func stepsToPolygons(steps []Stepper) [][]fpoint {
	var polys [][]fpoint
	var poly []fpoint
	var cx, cy, sx, sy float64

	flush := func() {
		if len(poly) > 0 {
			polys = append(polys, poly)
			poly = nil
		}
	}

	for _, s := range steps {
		switch v := s.(type) {
		case *MoveStep:
			flush()
			poly = append(poly, fpoint{float64(v.X), float64(v.Y)})
			cx = float64(v.X)
			cy = float64(v.Y)
			sx = cx
			sy = cy

		case *LineStep:
			poly = append(poly, fpoint{float64(v.X), float64(v.Y)})
			cx = float64(v.X)
			cy = float64(v.Y)

		case *CloseStep:
			if cx != sx || cy != sy {
				poly = append(poly, fpoint{sx, sy})
			}
			cx = sx
			cy = sy

		case *QuadStep:
			pts := evalQuad(cx, cy, float64(v.X1), float64(v.Y1), float64(v.X2), float64(v.Y2))
			for i := 1; i < len(pts); i++ {
				poly = append(poly, pts[i])
			}
			cx = float64(v.X2)
			cy = float64(v.Y2)

		case *CubicStep:
			pts := evalCubic(cx, cy, float64(v.X1), float64(v.Y1), float64(v.X2), float64(v.Y2), float64(v.X3), float64(v.Y3))
			for i := 1; i < len(pts); i++ {
				poly = append(poly, pts[i])
			}
			cx = float64(v.X3)
			cy = float64(v.Y3)

		case *ArcStep:
			pts := evalArc(float64(v.CX), float64(v.CY), float64(v.R), float64(v.Start), float64(v.End), v.Dir)
			for i := 1; i < len(pts); i++ {
				poly = append(poly, pts[i])
			}
			cx = pts[len(pts)-1].x
			cy = pts[len(pts)-1].y

		case *ArcToStep:
			poly = append(poly, fpoint{float64(v.X2), float64(v.Y2)})
			cx = float64(v.X2)
			cy = float64(v.Y2)
		}
	}
	flush()
	return polys
}

// evalQuad samples a quadratic Bézier into approxLineCount line segments.
func evalQuad(p0x, p0y, p1x, p1y, p2x, p2y float64) []fpoint {
	n := 16
	var pts []fpoint
	for i := 0; i <= n; i++ {
		t := float64(i) / float64(n)
		mt := 1 - t
		x := mt*mt*p0x + 2*mt*t*p1x + t*t*p2x
		y := mt*mt*p0y + 2*mt*t*p1y + t*t*p2y
		pts = append(pts, fpoint{x, y})
	}
	return pts
}

// evalCubic samples a cubic Bézier into approxLineCount line segments.
func evalCubic(p0x, p0y, p1x, p1y, p2x, p2y, p3x, p3y float64) []fpoint {
	n := 16
	var pts []fpoint
	for i := 0; i <= n; i++ {
		t := float64(i) / float64(n)
		mt := 1 - t
		x := mt*mt*mt*p0x + 3*mt*mt*t*p1x + 3*mt*t*t*p2x + t*t*t*p3x
		y := mt*mt*mt*p0y + 3*mt*mt*t*p1y + 3*mt*t*t*p2y + t*t*t*p3y
		pts = append(pts, fpoint{x, y})
	}
	return pts
}

// evalArc samples an arc into approxLineCount line segments.
func evalArc(cx, cy, r, start, end float64, dir Direction) []fpoint {
	n := 20
	if dir == CounterClockwise {
		if end < start {
			end += 2 * math.Pi
		}
	} else {
		if start < end {
			start += 2 * math.Pi
		}
	}
	var pts []fpoint
	for i := 0; i <= n; i++ {
		t := float64(i) / float64(n)
		a := start + (end-start)*t
		pts = append(pts, fpoint{
			cx + r*math.Cos(a),
			cy + r*math.Sin(a),
		})
	}
	return pts
}

// pointToSegment returns the perpendicular distance from (px, py) to the
// line segment defined by seg.
func pointToSegment(px, py float64, seg segment) float64 {
	dx := seg.x2 - seg.x1
	dy := seg.y2 - seg.y1
	ex := px - seg.x1
	ey := py - seg.y1

	t := (ex*dx + ey*dy) / (dx*dx + dy*dy)
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	cx := seg.x1 + t*dx
	cy := seg.y1 + t*dy

	d := (px-cx)*(px-cx) + (py-cy)*(py-cy)
	if d < 0 {
		return 0
	}
	return math.Sqrt(d)
}

// pointInPolygon returns true if (x, y) is inside the polygon using the ray
// casting algorithm.
func pointInPolygon(x, y float64, poly []fpoint) bool {
	n := len(poly)
	inside := false
	j := n - 1
	for i := 0; i < n; i++ {
		if ((poly[i].y > y) != (poly[j].y > y)) &&
			(x < (poly[j].x-poly[i].x)*(y-poly[i].y)/(poly[j].y-poly[i].y)+poly[i].x) {
			inside = !inside
		}
		j = i
	}
	return inside
}

// Move translates an instruction by the offset (dx, dy). All coordinates
// belonging to the instruction are shifted by the given amount.
func Move(inst Instruction, dx, dy float32) {
	switch v := inst.(type) {
	case *CircleInstruction:
		v.C.X += dx
		v.C.Y += dy
	case *DiskInstruction:
		v.C.X += dx
		v.C.Y += dy
	case *RectInstruction:
		v.X += dx
		v.Y += dy
	case *SlabInstruction:
		v.X += dx
		v.Y += dy
	case *LineInstruction:
		v.X1 += dx
		v.Y1 += dy
		v.X2 += dx
		v.Y2 += dy
	case *FillInstruction:
		moveSteps(v.Steps, dx, dy)
	case *StrokeInstruction:
		moveSteps(v.Steps, dx, dy)
	}
}

func moveSteps(steps []Stepper, dx, dy float32) {
	for _, s := range steps {
		switch v := s.(type) {
		case *MoveStep:
			v.X += dx
			v.Y += dy
		case *LineStep:
			v.X += dx
			v.Y += dy
		case *QuadStep:
			v.X1 += dx
			v.Y1 += dy
			v.X2 += dx
			v.Y2 += dy
		case *CubicStep:
			v.X1 += dx
			v.Y1 += dy
			v.X2 += dx
			v.Y2 += dy
			v.X3 += dx
			v.Y3 += dy
		case *ArcStep:
			v.CX += dx
			v.CY += dy
		case *ArcToStep:
			v.X1 += dx
			v.Y1 += dy
			v.X2 += dx
			v.Y2 += dy
		}
	}
}

// StrokeColor returns the color of the instruction. If the instruction has no
// color field the zero color is returned.
func StrokeColor(inst Instruction) Color {
	switch v := inst.(type) {
	case *CircleInstruction:
		return v.Color
	case *DiskInstruction:
		return v.Color
	case *RectInstruction:
		return v.Color
	case *SlabInstruction:
		return v.Color
	case *LineInstruction:
		return v.Color
	case *FillInstruction:
		return v.Color
	case *StrokeInstruction:
		return v.Color
	}
	return Color{}
}

// StrokeWidth returns the stroke width of the instruction, or 0 if it has no stroke.
func StrokeWidth(inst Instruction) float32 {
	switch v := inst.(type) {
	case *CircleInstruction:
		return float32(v.Stroke)
	case *RectInstruction:
		return float32(v.Stroke)
	case *LineInstruction:
		return float32(v.Stroke)
	case *StrokeInstruction:
		return float32(v.Stroke)
	}
	return 0
}

func min4(a, b, c, d float64) float64 {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	if d < m {
		m = d
	}
	return m
}
