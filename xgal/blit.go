package xgal

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// BlitOp specifies a flip or rotation for Blit.
type BlitOp int

const (
	FlipH  BlitOp = iota // Flip Horizontally
	FlipV                // Flip Vertically
	Rot90                // Rotate 90 degrees clockwise.
	Rot180               // Rotate 180 degrees clockwise.
	Rot270               // Rotate 270 degrees clockwise.
)

// Blit copies the source rectangle sr from src onto the destination rectangle
// dr of dst. Ops are flags: rotation is applied first, then flips.
func Blit(dst, src *Surface, dr, sr Rectangle, ops ...BlitOp) {
	sw, sh := float64(sr.Dx()), float64(sr.Dy())
	dw, dh := float64(dr.Dx()), float64(dr.Dy())

	sub := src.SubImage(sr).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	step := 0
	flipX, flipY := 1.0, 1.0

	for _, b := range ops {
		switch b {
		case FlipH:
			flipX = -flipX
		case FlipV:
			flipY = -flipY
		case Rot90:
			step = (step + 1) % 4
		case Rot180:
			step = (step + 2) % 4
		case Rot270:
			step = (step + 3) % 4
		}
	}

	// Effective dimensions after rotation
	w, h := sw, sh
	if step%2 == 1 {
		w, h = sh, sw
	}

	angle := float64(step) * math.Pi / 2

	// Rotation (applied first to the source)
	if angle != 0 {
		op.GeoM.Rotate(angle)
	}
	switch step {
	case 1:
		op.GeoM.Translate(w, 0)
	case 2:
		op.GeoM.Translate(sw, sh)
	case 3:
		op.GeoM.Translate(0, h)
	}

	// Flips (applied second, using effective rotated dimensions)
	if flipX < 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(w, 0)
	}
	if flipY < 0 {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, h)
	}

	// Scale and position
	op.GeoM.Scale(dw/w, dh/h)
	op.GeoM.Translate(float64(dr.Min.X), float64(dr.Min.Y))

	dst.DrawImage(sub, op)
}
