package xui

import (
	"math"

	"github.com/xmasengine/xmas/xgal"
)

// RingItem is a single entry in a [RingLayer].
type RingItem struct {
	Label   string
	Icon    Icon // optional icon drawn inside the item circle
	Action  func()
	SubRing *RingLayer // opened when confirmed, positioned at cursor
}

// RingLayer is a 16 bits era style ring menu. Items are arranged in a
// circle.  Press left/right to rotate the ring as items spin past a
// fixed cursor.  Confirm selects the item under the cursor.
// Cancel closes the ring.
type RingLayer struct {
	Bounds    xgal.Rectangle
	Style     Style
	Items     []RingItem
	Center    xgal.Point
	Radius    int
	Sel       int     // item index under the cursor
	cursorAng float64 // fixed cursor position in radians with default -PI/2 as the top

	Left    func() bool
	Right   func() bool
	Confirm func() bool
	Cancel  func() bool

	open    bool
	angle   float64    // open animation t going from 0 to one
	spin    int        // exact rotation step
	spinOff float64    // fractional offset for smooth animation
	sub     *RingLayer // active sub-ring, if any
}

// Ring creates a ring menu. Set Left/Right/Confirm/Cancel fields on
// the returned [RingLayer] to override the global [DefaultInput] bindings.
func Ring(cx, cy, radius int, items []RingItem) *RingLayer {
	return &RingLayer{
		Bounds:    xgal.Rect(cx-radius-24, cy-radius-24, cx+radius+24, cy+radius+24),
		Style:     DefaultStyle(),
		Items:     items,
		Center:    xgal.Pt(cx, cy),
		Radius:    radius,
		cursorAng: -math.Pi / 2, // top
	}
}

var _ Widget = &RingLayer{}

func (r *RingLayer) Poll() Reply {
	// delegate to active sub-ring
	if r.sub != nil {
		if r.sub.Poll() == Finish {
			r.sub = nil
		}
		return Accept
	}

	// animate opening
	if !r.open {
		r.angle += 0.15
		if r.angle >= 1 {
			r.angle = 1
			r.open = true
		}
		return Accept
	}

	// decay spin offset toward 0
	if r.spinOff != 0 {
		r.spinOff *= 0.8
		if math.Abs(r.spinOff) < 0.001 {
			r.spinOff = 0
		}
	}

	if input(r.Cancel, DefaultInput.Cancel) {
		return Finish
	}

	n := len(r.Items)
	if n == 0 {
		return Accept
	}

	if input(r.Left, DefaultInput.Left) {
		r.Sel = (r.Sel - 1 + n) % n
		r.spin++
		r.spinOff = -1
	}
	if input(r.Right, DefaultInput.Right) {
		r.Sel = (r.Sel + 1) % n
		r.spin--
		r.spinOff = 1
	}

	if input(r.Confirm, DefaultInput.Confirm) {
		if r.Sel >= 0 && r.Sel < n {
			item := r.Items[r.Sel]
			if item.SubRing != nil {
				cx := r.Center.X + int(float64(r.Radius)*math.Cos(r.cursorAng))
				cy := r.Center.Y + int(float64(r.Radius)*math.Sin(r.cursorAng))
				sub := item.SubRing
				sub.Center = xgal.Pt(cx, cy)
				sub.open = false
				sub.angle = 0
				sub.spin = 0
				sub.spinOff = 0
				sub.Sel = 0
				r.sub = sub
				return Accept
			}
			fn := item.Action
			if fn != nil {
				fn()
			}
		}
		return Finish
	}

	return Accept
}

func (r *RingLayer) Render(s *xgal.Surface) {
	n := len(r.Items)
	if n == 0 {
		return
	}

	a := r.angle
	if a <= 0 {
		return
	}
	rad := float64(r.Radius) * a
	spinAng := (float64(r.spin) + r.spinOff) * 2 * math.Pi / float64(n)

	// draw connecting ring segments
	for i := 0; i < n; i++ {
		ang1 := float64(i)*2*math.Pi/float64(n) + spinAng - math.Pi/2
		ang2 := float64(i+1)*2*math.Pi/float64(n) + spinAng - math.Pi/2
		x1 := r.Center.X + int(rad*math.Cos(ang1))
		y1 := r.Center.Y + int(rad*math.Sin(ang1))
		x2 := r.Center.X + int(rad*math.Cos(ang2))
		y2 := r.Center.Y + int(rad*math.Sin(ang2))
		xgal.Line(s, x1, y1, x2, y2, 2, r.Style.Border)
	}

	// draw items with label
	for i := 0; i < n; i++ {
		item := &r.Items[i]
		ang := float64(i)*2*math.Pi/float64(n) + spinAng - math.Pi/2
		ix := r.Center.X + int(rad*math.Cos(ang))
		iy := r.Center.Y + int(rad*math.Sin(ang))

		st := r.Style
		if i == r.Sel {
			st = st.FocusStyle()
		}

		// item circle
		st.DrawCircle(s, xgal.Pt(ix, iy), 14)

		// icon inside the circle
		if item.Icon.Image != nil {
			ib := item.Icon.Image.Bounds()
			xgal.Blit(s, item.Icon.Image,
				xgal.Rect(ix-ib.Dx()/2, iy-ib.Dy()/2, ix+ib.Dx()/2, iy+ib.Dy()/2),
				ib)
		}

		// label below the circle
		labelBounds := xgal.Rect(ix-14, iy+14, ix+14, iy+28)
		st.Ink(s, labelBounds, item.Label)

		// cursor highlight on the selected item
		if i == r.Sel {
			cx := r.Center.X + int(rad*math.Cos(r.cursorAng))
			cy := r.Center.Y + int(rad*math.Sin(r.cursorAng))
			xgal.Circle(s, xgal.Pt(cx, cy), 18, 3, xgal.Wash(255, 255, 100, 255))
		}
	}

	// render active sub-ring on top
	if r.sub != nil {
		r.sub.Render(s)
	}
}

func (r *RingLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	r.Bounds = bounds
	return r.Bounds
}

func (r *RingLayer) MoveBy(delta xgal.Point) {
	r.Bounds = r.Bounds.Add(delta)
	r.Center = r.Center.Add(delta)
	if r.sub != nil {
		r.sub.MoveBy(delta)
	}
}
