package xvec

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// ParseSVG parses an SVG document and converts supported drawing elements to
// xvec instructions scaled to the given canvas size. If width or height is 0,
// the SVG's own dimensions or viewBox are used instead.
//
// Supported elements: svg, g, path, rect, circle, ellipse, line, polyline,
// polygon. Unsupported elements are silently dropped. Supported attributes:
// fill, stroke, stroke-width, transform, d, points, cx, cy, r, rx, ry, x, y,
// width, height, x1, y1, x2, y2, viewBox, style.
//
// Colours support #RRGGBB, #RGB, #RRGGBBAA, and a set of named colours.
// The fill="none" and stroke="none" values are recognised.
func ParseSVG(r io.Reader, width, height float32) (*XVEC, error) {
	dec := xml.NewDecoder(r)

	x := &XVEC{Size: Size{W: 320, H: 240}, Antialias: true}

	var svgW, svgH float32
	var vb struct{ minX, minY, w, h float32 }
	var hasVB bool

	stack := []svgCtx{{fill: "black", stroke: "none", strokeWidth: "1", tr: affine{a: 1, d: 1}}}

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("svg: %w", err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			c := stack[len(stack)-1]
			for _, a := range t.Attr {
				switch a.Name.Local {
				case "fill":
					c.fill = a.Value
				case "stroke":
					c.stroke = a.Value
				case "stroke-width":
					c.strokeWidth = a.Value
				case "transform":
					c.tr = c.tr.chain(parseTransform(a.Value))
				}
			}

			switch t.Name.Local {
			case "svg":
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "width":
						svgW = parseLength(a.Value)
					case "height":
						svgH = parseLength(a.Value)
					case "viewBox":
						parts := strings.Fields(a.Value)
						if len(parts) == 4 {
							vb.minX = parseLength(parts[0])
							vb.minY = parseLength(parts[1])
							vb.w = parseLength(parts[2])
							vb.h = parseLength(parts[3])
							hasVB = true
						}
					}
				}

			case "g":
				// group — nothing to emit, children handled at deeper level

			case "rect":
				var rx, ry, rw, rh float32
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "x":
						rx = parseLength(a.Value)
					case "y":
						ry = parseLength(a.Value)
					case "width":
						rw = parseLength(a.Value)
					case "height":
						rh = parseLength(a.Value)
					}
				}
				emitShape(x, c, func() []Stepper {
					if c.tr.isIdentity() {
						return nil // signal: use simple primitive
					}
					x0, y0 := c.tr.apply(rx, ry)
					x1, y1 := c.tr.apply(rx+rw, ry)
					x2, y2 := c.tr.apply(rx+rw, ry+rh)
					x3, y3 := c.tr.apply(rx, ry+rh)
					return []Stepper{
						&MoveStep{X: x0, Y: y0},
						&LineStep{X: x1, Y: y1},
						&LineStep{X: x2, Y: y2},
						&LineStep{X: x3, Y: y3},
						&CloseStep{},
					}
				}, func() (Color, float32) {
					return parseColor(c.fill), parseLength(c.strokeWidth)
				}, func() (Color, float32) {
					return parseColor(c.stroke), parseLength(c.strokeWidth)
				}, func() (float32, float32, float32, float32) {
					return rx, ry, rw, rh
				})

			case "circle":
				var cx, cy, r float32
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "cx":
						cx = parseLength(a.Value)
					case "cy":
						cy = parseLength(a.Value)
					case "r":
						r = parseLength(a.Value)
					}
				}
				emitShape(x, c, func() []Stepper {
					if c.tr.isIdentity() {
						return nil
					}
					return circlePath(cx, cy, r, c.tr)
				}, func() (Color, float32) {
					return parseColor(c.fill), 0
				}, func() (Color, float32) {
					return parseColor(c.stroke), parseLength(c.strokeWidth)
				}, func() (float32, float32, float32, float32) {
					return cx - r, cy - r, 2 * r, 2 * r
				})

			case "ellipse":
				var cx, cy, rx, ry float32
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "cx":
						cx = parseLength(a.Value)
					case "cy":
						cy = parseLength(a.Value)
					case "rx":
						rx = parseLength(a.Value)
					case "ry":
						ry = parseLength(a.Value)
					}
				}
				emitShape(x, c, func() []Stepper {
					return ellipsePath(cx, cy, rx, ry, c.tr)
				}, func() (Color, float32) {
					return parseColor(c.fill), 0
				}, func() (Color, float32) {
					return parseColor(c.stroke), parseLength(c.strokeWidth)
				}, func() (float32, float32, float32, float32) {
					return cx - rx, cy - ry, 2 * rx, 2 * ry
				})

			case "line":
				var x1, y1, x2, y2 float32
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "x1":
						x1 = parseLength(a.Value)
					case "y1":
						y1 = parseLength(a.Value)
					case "x2":
						x2 = parseLength(a.Value)
					case "y2":
						y2 = parseLength(a.Value)
					}
				}
				emitShape(x, c, func() []Stepper {
					if c.tr.isIdentity() {
						return nil
					}
					ax, ay := c.tr.apply(x1, y1)
					bx, by := c.tr.apply(x2, y2)
					return []Stepper{
						&MoveStep{X: ax, Y: ay},
						&LineStep{X: bx, Y: by},
					}
				}, func() (Color, float32) {
					return Color{}, 0 // line has no fill
				}, func() (Color, float32) {
					return parseColor(c.stroke), parseLength(c.strokeWidth)
				}, func() (float32, float32, float32, float32) {
					return x1, y1, x2 - x1, y2 - y1
				})

			case "polyline", "polygon":
				var pts []float32
				for _, a := range t.Attr {
					if a.Name.Local == "points" {
						pts = parsePoints(a.Value)
					}
				}
				emitShape(x, c, func() []Stepper {
					return pointsToPath(pts, c.tr, t.Name.Local == "polygon")
				}, func() (Color, float32) {
					return parseColor(c.fill), 0
				}, func() (Color, float32) {
					return parseColor(c.stroke), parseLength(c.strokeWidth)
				}, nil)

			case "path":
				var d string
				for _, a := range t.Attr {
					if a.Name.Local == "d" {
						d = a.Value
					}
				}
				steps, err := parsePathData(d)
				if err != nil {
					return nil, fmt.Errorf("svg: path d: %w", err)
				}
				emitShape(x, c, func() []Stepper {
					if c.tr.isIdentity() {
						return steps
					}
					return transformSteps(steps, c.tr)
				}, func() (Color, float32) {
					return parseColor(c.fill), 0
				}, func() (Color, float32) {
					return parseColor(c.stroke), parseLength(c.strokeWidth)
				}, nil)
			}

			stack = append(stack, c)

		case xml.EndElement:
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
		}
	}

	// Apply viewBox/scaling to output size.
	outW, outH := width, height
	if outW == 0 {
		if hasVB {
			outW = vb.w
		} else {
			outW = svgW
		}
	}
	if outH == 0 {
		if hasVB {
			outH = vb.h
		} else {
			outH = svgH
		}
	}
	if outW == 0 {
		outW = 320
	}
	if outH == 0 {
		outH = 240
	}

	// Scale to output size.
	if hasVB {
		scaleX := outW / vb.w
		scaleY := outH / vb.h
		if scaleX != 1 || scaleY != 1 {
			scaleAll(x, scaleX, scaleY)
			translateAll(x, -vb.minX*scaleX, -vb.minY*scaleY)
		}
	} else if outW != svgW || outH != svgH {
		if svgW > 0 && svgH > 0 {
			scaleX := outW / svgW
			scaleY := outH / svgH
			scaleAll(x, scaleX, scaleY)
		}
	}

	x.Size = Size{W: outW, H: outH}
	return x, nil
}

// svgCTX is the  SVG context
type svgCtx struct {
	fill, stroke, strokeWidth string
	tr                        affine
}

// emitShape emits fill/stroke instructions for a shape. If path returns nil,
// the shape supports simple primitives via rect.
func emitShape(
	x *XVEC, c svgCtx,
	path func() []Stepper,
	fillFn func() (Color, float32),
	strokeFn func() (Color, float32),
	rectFn func() (float32, float32, float32, float32),
) {
	fillCol, _ := fillFn()
	strokeCol, sw := strokeFn()

	steps := path()
	usePath := steps != nil

	if c.fill != "none" && c.fill != "" && c.fill != "transparent" {
		if usePath {
			x.Fill(fillCol, steps...)
		} else if rectFn != nil {
			rx, ry, rw, rh := rectFn()
			x.Slab(rx, ry, rw, rh, fillCol)
		}
	}

	if c.stroke != "none" && c.stroke != "" && c.stroke != "transparent" {
		if sw == 0 {
			sw = 1
		}
		if usePath {
			x.Stroke(sw, strokeCol, steps...)
		} else if rectFn != nil {
			rx, ry, rw, rh := rectFn()
			x.Rect(rx, ry, rw, rh, sw, strokeCol)
		}
	}
}

const bezierK float32 = 0.5522847498 // 4*(sqrt(2)-1)/3

func circlePath(cx, cy, r float32, tr affine) []Stepper {
	return ellipsePath(cx, cy, r, r, tr)
}

func ellipsePath(cx, cy, rx, ry float32, tr affine) []Stepper {
	var steps []Stepper
	emit := func(x, y float32) {
		x, y = tr.apply(x, y)
		if len(steps) == 0 {
			steps = append(steps, &MoveStep{X: x, Y: y})
		}
	}
	// Four 90-degree cubic bezier segments.
	// Starting at angle 0 (rightmost point), going counter-clockwise.
	for i := 0; i < 4; i++ {
		a1 := float64(i) * math.Pi / 2
		a2 := float64(i+1) * math.Pi / 2
		cos1, sin1 := float32(math.Cos(a1)), float32(math.Sin(a1))
		cos2, sin2 := float32(math.Cos(a2)), float32(math.Sin(a2))

		x0 := cx + rx*cos1
		y0 := cy + ry*sin1
		x1 := cx + rx*(cos1-bezierK*sin1)
		y1 := cy + ry*(sin1+bezierK*cos1)
		x2 := cx + rx*(cos2+bezierK*sin2)
		y2 := cy + ry*(sin2-bezierK*cos2)
		x3 := cx + rx*cos2
		y3 := cy + ry*sin2

		if i == 0 {
			emit(x0, y0)
		}
		x1, y1 = tr.apply(x1, y1)
		x2, y2 = tr.apply(x2, y2)
		x3, y3 = tr.apply(x3, y3)
		steps = append(steps, &CubicStep{X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3})
	}
	steps = append(steps, &CloseStep{})
	return steps
}

func pointsToPath(pts []float32, tr affine, closePath bool) []Stepper {
	if len(pts) < 2 {
		return nil
	}
	steps := make([]Stepper, 0, len(pts)/2+1)
	for i := 0; i+1 < len(pts); i += 2 {
		x, y := tr.apply(pts[i], pts[i+1])
		if i == 0 {
			steps = append(steps, &MoveStep{X: x, Y: y})
		} else {
			steps = append(steps, &LineStep{X: x, Y: y})
		}
	}
	if closePath {
		steps = append(steps, &CloseStep{})
	}
	return steps
}

// svgArcToCubics converts an SVG elliptical arc to cubic bezier segments.
// The arc goes from (x1,y1) to (x2,y2) with the given radii and rotation.
// Returns at most one CubicStep per 90-degree arc segment.
func svgArcToCubics(x1, y1, rx, ry, rot float32, large, sweep int, x2, y2 float32) []*CubicStep {
	if rx == 0 || ry == 0 || (x1 == x2 && y1 == y2) {
		return nil
	}

	cosR := float32(math.Cos(float64(rot)))
	sinR := float32(math.Sin(float64(rot)))

	// Step 1: Compute (x1', y1') in rotated coordinate system.
	dx := (x1 - x2) / 2
	dy := (y1 - y2) / 2
	x1p := cosR*dx + sinR*dy
	y1p := -sinR*dx + cosR*dy

	// Ensure radii are large enough.
	l := x1p*x1p/(rx*rx) + y1p*y1p/(ry*ry)
	if l > 1 {
		s := float32(math.Sqrt(float64(l)))
		rx *= s
		ry *= s
	}

	// Step 2: Compute center (cx', cy') in rotated system.
	denom := rx*rx*y1p*y1p + ry*ry*x1p*x1p
	if denom == 0 {
		return nil
	}
	num := rx*rx*ry*ry - denom
	if num < 0 {
		num = 0
	}
	var factor float32
	if num > 0 {
		factor = float32(math.Sqrt(float64(num) / float64(denom)))
	}
	if large == sweep {
		factor = -factor
	}
	cxp := factor * rx * y1p / ry
	cyp := -factor * ry * x1p / rx

	// Step 3: Compute center (cx, cy) in original coordinates.
	cx := cosR*cxp - sinR*cyp + (x1+x2)/2
	cy := sinR*cxp + cosR*cyp + (y1+y2)/2

	// Step 4: Compute start and delta angles.
	ux := (x1p - cxp) / rx
	uy := (y1p - cyp) / ry
	vx := (-x1p - cxp) / rx
	vy := (-y1p - cyp) / ry

	theta1 := angle(ux, uy)
	dTheta := angle(vx, vy) - theta1

	if sweep == 0 && dTheta > 0 {
		dTheta -= 2 * math.Pi
	} else if sweep == 1 && dTheta < 0 {
		dTheta += 2 * math.Pi
	}

	// Step 5: Split into segments of at most 90 degrees.
	n := int(math.Ceil(float64(math.Abs(float64(dTheta)) / (math.Pi / 2))))
	seg := dTheta / float32(n)

	var steps []*CubicStep
	for i := 0; i < n; i++ {
		η1 := theta1 + float32(i)*seg
		η2 := theta1 + float32(i+1)*seg
		step := arcSegment(cx, cy, rx, ry, cosR, sinR, η1, η2)
		steps = append(steps, step)
	}
	return steps
}

// angle returns the angle of (x,y) from the positive x-axis, in [-π, π].
func angle(x, y float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

// arcSegment produces one cubic bezier for an arc from η1 to η2 on an
// ellipse centred at (cx,cy) with radii (rx,ry), rotated by cosR/sinR.
func arcSegment(cx, cy, rx, ry, cosR, sinR, η1, η2 float32) *CubicStep {
	cos1 := float32(math.Cos(float64(η1)))
	sin1 := float32(math.Sin(float64(η1)))
	cos2 := float32(math.Cos(float64(η2)))
	sin2 := float32(math.Sin(float64(η2)))

	Δ := η2 - η1
	α := float32(4.0 / 3.0 * math.Tan(float64(Δ)/4))

	// Control points in ellipse-local space (centered at origin).
	lx0 := rx * cos1
	ly0 := ry * sin1
	lx1 := rx * (cos1 - α*sin1)
	ly1 := ry * (sin1 + α*cos1)
	lx2 := rx * (cos2 + α*sin2)
	ly2 := ry * (sin2 - α*cos2)
	lx3 := rx * cos2
	ly3 := ry * sin2

	// Rotate by rot and translate to center.
	rot := func(x, y float32) (float32, float32) {
		return x*cosR - y*sinR + cx, x*sinR + y*cosR + cy
	}
	_, _ = lx0, ly0 // starting point (implicit, matches previous end)
	x1, y1 := rot(lx1, ly1)
	x2, y2 := rot(lx2, ly2)
	x3, y3 := rot(lx3, ly3)
	return &CubicStep{X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3}
}

func transformSteps(steps []Stepper, tr affine) []Stepper {
	out := make([]Stepper, len(steps))
	for i, s := range steps {
		switch v := s.(type) {
		case *MoveStep:
			x, y := tr.apply(v.X, v.Y)
			out[i] = &MoveStep{X: x, Y: y}
		case *LineStep:
			x, y := tr.apply(v.X, v.Y)
			out[i] = &LineStep{X: x, Y: y}
		case *QuadStep:
			x1, y1 := tr.apply(v.X1, v.Y1)
			x2, y2 := tr.apply(v.X2, v.Y2)
			out[i] = &QuadStep{X1: x1, Y1: y1, X2: x2, Y2: y2}
		case *CubicStep:
			x1, y1 := tr.apply(v.X1, v.Y1)
			x2, y2 := tr.apply(v.X2, v.Y2)
			x3, y3 := tr.apply(v.X3, v.Y3)
			out[i] = &CubicStep{X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3}
		case *ArcStep:
			cx, cy := tr.apply(v.CX, v.CY)
			// Scale radius by average of transform scale factors.
			sx := float32(math.Sqrt(float64(tr.a*tr.a + tr.b*tr.b)))
			sy := float32(math.Sqrt(float64(tr.c*tr.c + tr.d*tr.d)))
			r := v.R * (sx + sy) / 2
			out[i] = &ArcStep{CX: cx, CY: cy, R: r, Start: v.Start, End: v.End, Dir: v.Dir}
		case *CloseStep:
			out[i] = &CloseStep{}
		default:
			out[i] = s
		}
	}
	return out
}

func scaleAll(x *XVEC, sx, sy float32) {
	tr := affine{a: sx, d: sy}
	for _, inst := range x.Instructions {
		switch v := inst.(type) {
		case *CircleInstruction:
			v.C.X *= sx
			v.C.Y *= sy
			v.R = Length(float32(v.R) * (sx + sy) / 2)
			v.Stroke = Length(float32(v.Stroke) * (sx + sy) / 2)
		case *DiskInstruction:
			v.C.X *= sx
			v.C.Y *= sy
			v.R = Length(float32(v.R) * (sx + sy) / 2)
		case *RectInstruction:
			v.X *= sx
			v.Y *= sy
			v.W *= sx
			v.H *= sy
			v.Stroke = Length(float32(v.Stroke) * (sx + sy) / 2)
		case *SlabInstruction:
			v.X *= sx
			v.Y *= sy
			v.W *= sx
			v.H *= sy
		case *LineInstruction:
			v.X1 *= sx
			v.Y1 *= sy
			v.X2 *= sx
			v.Y2 *= sy
			v.Stroke = Length(float32(v.Stroke) * (sx + sy) / 2)
		case *FillInstruction:
			v.Steps = transformSteps(v.Steps, tr)
		case *StrokeInstruction:
			v.Stroke = Length(float32(v.Stroke) * (sx + sy) / 2)
			v.Steps = transformSteps(v.Steps, tr)
		}
	}
}

func translateAll(x *XVEC, dx, dy float32) {
	tr := affine{e: dx, f: dy}
	for _, inst := range x.Instructions {
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
			v.Steps = transformSteps(v.Steps, tr)
		case *StrokeInstruction:
			v.Steps = transformSteps(v.Steps, tr)
		}
	}
}

// ── 2D affine transform ──

type affine struct {
	a, b, c, d, e, f float32
}

func (tr affine) isIdentity() bool {
	return tr.a == 1 && tr.b == 0 && tr.c == 0 && tr.d == 1 && tr.e == 0 && tr.f == 0
}

func (tr affine) apply(x, y float32) (float32, float32) {
	return tr.a*x + tr.c*y + tr.e, tr.b*x + tr.d*y + tr.f
}

func (tr affine) chain(other affine) affine {
	return affine{
		a: tr.a*other.a + tr.c*other.b,
		b: tr.b*other.a + tr.d*other.b,
		c: tr.a*other.c + tr.c*other.d,
		d: tr.b*other.c + tr.d*other.d,
		e: tr.a*other.e + tr.c*other.f + tr.e,
		f: tr.b*other.e + tr.d*other.f + tr.f,
	}
}

// ── Colour parsing ──

var namedColours = map[string]Color{
	"black":       {R: 0, G: 0, B: 0, A: 255},
	"white":       {R: 255, G: 255, B: 255, A: 255},
	"red":         {R: 255, G: 0, B: 0, A: 255},
	"green":       {R: 0, G: 128, B: 0, A: 255},
	"blue":        {R: 0, G: 0, B: 255, A: 255},
	"yellow":      {R: 255, G: 255, B: 0, A: 255},
	"cyan":        {R: 0, G: 255, B: 255, A: 255},
	"magenta":     {R: 255, G: 0, B: 255, A: 255},
	"orange":      {R: 255, G: 165, B: 0, A: 255},
	"purple":      {R: 128, G: 0, B: 128, A: 255},
	"gray":        {R: 128, G: 128, B: 128, A: 255},
	"grey":        {R: 128, G: 128, B: 128, A: 255},
	"silver":      {R: 192, G: 192, B: 192, A: 255},
	"maroon":      {R: 128, G: 0, B: 0, A: 255},
	"navy":        {R: 0, G: 0, B: 128, A: 255},
	"olive":       {R: 128, G: 128, B: 0, A: 255},
	"teal":        {R: 0, G: 128, B: 128, A: 255},
	"aqua":        {R: 0, G: 255, B: 255, A: 255},
	"fuchsia":     {R: 255, G: 0, B: 255, A: 255},
	"lime":        {R: 0, G: 255, B: 0, A: 255},
	"transparent": {R: 0, G: 0, B: 0, A: 0},
}

func parseColor(s string) Color {
	if s == "" {
		return Color{R: 0, G: 0, B: 0, A: 255}
	}
	if s[0] == '#' {
		var hex string
		switch len(s) {
		case 4: // #RGB
			hex = fmt.Sprintf("%c%c%c%c%c%c", s[1], s[1], s[2], s[2], s[3], s[3])
		case 7: // #RRGGBB
			hex = s[1:] + "ff"
		case 9: // #RRGGBBAA
			hex = s[1:]
		default:
			return Color{R: 0, G: 0, B: 0, A: 255}
		}
		v, err := strconv.ParseUint(hex, 16, 64)
		if err != nil {
			return Color{R: 0, G: 0, B: 0, A: 255}
		}
		return Color{
			R: uint8(v >> 24),
			G: uint8(v >> 16),
			B: uint8(v >> 8),
			A: uint8(v),
		}
	}
	if c, ok := namedColours[strings.ToLower(s)]; ok {
		return c
	}
	return Color{R: 0, G: 0, B: 0, A: 255}
}

// ── Transform parsing ──

func parseTransform(s string) affine {
	if s == "" {
		return affine{}
	}
	tr := affine{a: 1, d: 1}
	s = strings.TrimSpace(s)
	if s == "" {
		return tr
	}
	// Consume each function(...)
	for {
		s = strings.TrimSpace(s)
		if s == "" {
			break
		}
		idx := strings.IndexByte(s, '(')
		if idx < 0 {
			break
		}
		fn := strings.TrimSpace(s[:idx])
		s = s[idx+1:]
		idx = strings.IndexByte(s, ')')
		if idx < 0 {
			break
		}
		args := s[:idx]
		s = s[idx+1:]

		nums := parseNumbers(args)
		switch fn {
		case "translate":
			var tx, ty float32
			if len(nums) > 0 {
				tx = nums[0]
			}
			if len(nums) > 1 {
				ty = nums[1]
			}
			tr = tr.chain(affine{a: 1, d: 1, e: tx, f: ty})
		case "scale":
			var sx, sy float32
			if len(nums) > 0 {
				sx = nums[0]
			}
			if len(nums) > 1 {
				sy = nums[1]
			} else {
				sy = sx
			}
			tr = tr.chain(affine{a: sx, d: sy})
		case "rotate":
			if len(nums) == 1 {
				a := nums[0] * math.Pi / 180
				cos := float32(math.Cos(float64(a)))
				sin := float32(math.Sin(float64(a)))
				tr = tr.chain(affine{a: cos, b: sin, c: -sin, d: cos})
			} else if len(nums) == 3 {
				a := nums[0] * math.Pi / 180
				cx, cy := nums[1], nums[2]
				cos := float32(math.Cos(float64(a)))
				sin := float32(math.Sin(float64(a)))
				t1 := affine{a: 1, d: 1, e: cx, f: cy}
				r := affine{a: cos, b: sin, c: -sin, d: cos}
				t2 := affine{a: 1, d: 1, e: -cx, f: -cy}
				tr = tr.chain(t1.chain(r.chain(t2)))
			}
		case "skewX":
			if len(nums) > 0 {
				a := nums[0] * math.Pi / 180
				tr = tr.chain(affine{a: 1, d: 1, c: float32(math.Tan(float64(a)))})
			}
		case "skewY":
			if len(nums) > 0 {
				a := nums[0] * math.Pi / 180
				tr = tr.chain(affine{a: 1, d: 1, b: float32(math.Tan(float64(a)))})
			}
		case "matrix":
			if len(nums) >= 6 {
				tr = tr.chain(affine{
					a: nums[0], b: nums[1], c: nums[2],
					d: nums[3], e: nums[4], f: nums[5],
				})
			}
		}
	}
	return tr
}

// ── SVG path data parser ──

func parsePathData(d string) ([]Stepper, error) {
	p := pathParser{data: d}
	var steps []Stepper
	var cx, cy float32   // current point
	var pcx, pcy float32 // previous control point (for S/s, T/t)
	relative := false

	for {
		cmd := p.nextCmd()
		if cmd == 0 {
			break
		}
		relative = cmd >= 'a' && cmd <= 'z'
		abs := cmd
		if relative {
			abs = cmd - 32
		}
		switch abs {
		case 'M':
			for {
				x, y := p.nextNum(), p.nextNum()
				if p.err != nil {
					return nil, fmt.Errorf("path: M needs coordinates")
				}
				if relative {
					x += cx
					y += cy
				}
				cx, cy = x, y
				if len(steps) == 0 || !isMove(steps[len(steps)-1]) {
					steps = append(steps, &MoveStep{X: cx, Y: cy})
					pcx, pcy = cx, cy
				}
				if !p.more() {
					break
				}
				// Remaining pairs are implicit L/l.
				for p.more() {
					x, y := p.nextNum(), p.nextNum()
					if p.err != nil {
						break
					}
					if relative {
						x += cx
						y += cy
					}
					cx, cy = x, y
					steps = append(steps, &LineStep{X: cx, Y: cy})
				}
				break
			}

		case 'L':
			for p.more() {
				x, y := p.nextNum(), p.nextNum()
				if p.err != nil {
					break
				}
				if relative {
					cx += x
					cy += y
				} else {
					cx, cy = x, y
				}
				steps = append(steps, &LineStep{X: cx, Y: cy})
			}
			if p.err != nil {
				return nil, fmt.Errorf("path: L needs coordinates")
			}
			pcx, pcy = cx, cy

		case 'H':
			for p.more() {
				v := p.nextNum()
				if p.err != nil {
					break
				}
				if relative {
					cx += v
				} else {
					cx = v
				}
				steps = append(steps, &LineStep{X: cx, Y: cy})
			}
			if p.err != nil {
				return nil, fmt.Errorf("path: H needs a number")
			}

		case 'V':
			for p.more() {
				v := p.nextNum()
				if p.err != nil {
					break
				}
				if relative {
					cy += v
				} else {
					cy = v
				}
				steps = append(steps, &LineStep{X: cx, Y: cy})
			}
			if p.err != nil {
				return nil, fmt.Errorf("path: V needs a number")
			}

		case 'C':
			for {
				x1, y1 := p.nextNum(), p.nextNum()
				x2, y2 := p.nextNum(), p.nextNum()
				x3, y3 := p.nextNum(), p.nextNum()
				if p.err != nil {
					return nil, fmt.Errorf("path: C needs 6 coordinates")
				}
				if relative {
					x1 += cx
					y1 += cy
					x2 += cx
					y2 += cy
					x3 += cx
					y3 += cy
				}
				steps = append(steps, &CubicStep{X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3})
				pcx, pcy = x2, y2
				cx, cy = x3, y3
				if !p.more() {
					break
				}
			}

		case 'S':
			for {
				x2, y2 := p.nextNum(), p.nextNum()
				x3, y3 := p.nextNum(), p.nextNum()
				if p.err != nil {
					return nil, fmt.Errorf("path: S needs 4 coordinates")
				}
				if relative {
					x2 += cx
					y2 += cy
					x3 += cx
					y3 += cy
				}
				x1, y1 := reflectControl(cx, cy, pcx, pcy)
				steps = append(steps, &CubicStep{X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3})
				pcx, pcy = x2, y2
				cx, cy = x3, y3
				if !p.more() {
					break
				}
			}

		case 'Q':
			for {
				x1, y1 := p.nextNum(), p.nextNum()
				x2, y2 := p.nextNum(), p.nextNum()
				if p.err != nil {
					return nil, fmt.Errorf("path: Q needs 4 coordinates")
				}
				if relative {
					x1 += cx
					y1 += cy
					x2 += cx
					y2 += cy
				}
				steps = append(steps, &QuadStep{X1: x1, Y1: y1, X2: x2, Y2: y2})
				pcx, pcy = x1, y1
				cx, cy = x2, y2
				if !p.more() {
					break
				}
			}

		case 'T':
			for {
				x2, y2 := p.nextNum(), p.nextNum()
				if p.err != nil {
					return nil, fmt.Errorf("path: T needs 2 coordinates")
				}
				if relative {
					x2 += cx
					y2 += cy
				}
				x1, y1 := reflectControl(cx, cy, pcx, pcy)
				steps = append(steps, &QuadStep{X1: x1, Y1: y1, X2: x2, Y2: y2})
				pcx, pcy = x1, y1
				cx, cy = x2, y2
				if !p.more() {
					break
				}
			}

		case 'A':
			for {
				rx := p.nextNum()
				ry := p.nextNum()
				xRot := p.nextNum()
				large := int(p.nextNum())
				sweep := int(p.nextNum())
				x2 := p.nextNum()
				y2 := p.nextNum()
				if p.err != nil {
					return nil, fmt.Errorf("path: A needs 7 parameters")
				}
				if relative {
					x2 += cx
					y2 += cy
				}
				cubics := svgArcToCubics(cx, cy, rx, ry, xRot*math.Pi/180, large, sweep, x2, y2)
				for _, c := range cubics {
					steps = append(steps, c)
				}
				cx, cy = x2, y2
				pcx, pcy = cx, cy
				if !p.more() {
					break
				}
			}

		case 'Z':
			steps = append(steps, &CloseStep{})
			// Reset to last M position.
			for i := len(steps) - 2; i >= 0; i-- {
				if m, ok := steps[i].(*MoveStep); ok {
					cx, cy = m.X, m.Y
					break
				}
			}
		}
	}
	return steps, p.err
}

func isMove(s Stepper) bool {
	switch s.(type) {
	case *MoveStep, *CloseStep:
		return true
	}
	return false
}

type pathParser struct {
	data string
	pos  int
	err  error
}

func (p *pathParser) nextCmd() byte {
	for p.pos < len(p.data) {
		c := p.data[p.pos]
		switch {
		case c >= 'A' && c <= 'Z':
			p.pos++
			return c
		case c >= 'a' && c <= 'z':
			p.pos++
			return c
		default:
			p.pos++
		}
	}
	return 0
}

func (p *pathParser) nextNum() float32 {
	// Skip whitespace and commas.
	for p.pos < len(p.data) {
		c := p.data[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == ',' {
			p.pos++
			continue
		}
		break
	}
	if p.pos >= len(p.data) {
		p.err = fmt.Errorf("expected number")
		return 0
	}
	start := p.pos
	// Optional sign.
	if p.pos < len(p.data) && (p.data[p.pos] == '-' || p.data[p.pos] == '+') {
		p.pos++
	}
	// Digits.
	hasDigits := false
	for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
		p.pos++
		hasDigits = true
	}
	// Optional decimal point and more digits.
	if p.pos < len(p.data) && p.data[p.pos] == '.' {
		p.pos++
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
			hasDigits = true
		}
	}
	if !hasDigits {
		p.err = fmt.Errorf("expected number at pos %d", start)
		return 0
	}
	// Optional exponent.
	if p.pos < len(p.data) && (p.data[p.pos] == 'e' || p.data[p.pos] == 'E') {
		p.pos++
		if p.pos < len(p.data) && (p.data[p.pos] == '-' || p.data[p.pos] == '+') {
			p.pos++
		}
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	}
	v, err := strconv.ParseFloat(p.data[start:p.pos], 32)
	if err != nil {
		p.err = fmt.Errorf("bad number %q", p.data[start:p.pos])
		return 0
	}
	return float32(v)
}

// more reports whether the next non-whitespace character is a number start.
func (p *pathParser) more() bool {
	for i := p.pos; i < len(p.data); i++ {
		c := p.data[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == ',' {
			continue
		}
		if c == '-' || c == '+' || c == '.' || (c >= '0' && c <= '9') {
			return true
		}
		return false
	}
	return false
}

func reflectControl(cx, cy, pcx, pcy float32) (float32, float32) {
	return cx + (cx - pcx), cy + (cy - pcy)
}

// ── Points attribute parser ──

func parsePoints(s string) []float32 {
	return parseNumbers(s)
}

// parseNumbers reads a sequence of floats separated by whitespace/commas.
func parseNumbers(s string) []float32 {
	var out []float32
	for i := 0; i < len(s); {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == ',' {
			i++
			continue
		}
		if c == '-' || c == '+' || c == '.' || (c >= '0' && c <= '9') {
			start := i
			if i < len(s) && (s[i] == '-' || s[i] == '+') {
				i++
			}
			for i < len(s) && s[i] >= '0' && s[i] <= '9' {
				i++
			}
			if i < len(s) && s[i] == '.' {
				i++
				for i < len(s) && s[i] >= '0' && s[i] <= '9' {
					i++
				}
			}
			if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
				i++
				if i < len(s) && (s[i] == '-' || s[i] == '+') {
					i++
				}
				for i < len(s) && s[i] >= '0' && s[i] <= '9' {
					i++
				}
			}
			v, _ := strconv.ParseFloat(s[start:i], 32)
			out = append(out, float32(v))
			continue
		}
		i++
	}
	return out
}

// ── Length parsing ──

func parseLength(s string) float32 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	// Strip trailing unit suffix.
	for i, c := range s {
		if c == 'e' || c == 'E' {
			continue // part of exponent
		}
		if c < '0' || c > '9' {
			if c == '.' || c == '-' || c == '+' {
				continue
			}
			if i > 0 && (s[i-1] >= '0' && s[i-1] <= '9') {
				v, err := strconv.ParseFloat(s[:i], 32)
				if err == nil {
					return float32(v)
				}
				return 0
			}
		}
	}
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(v)
}
