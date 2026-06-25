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

// Filter controls how pixel colors are sampled when drawing a [Surface]
// at a non-1:1 scale.
type Filter = ebiten.Filter

const (
	// Nearest uses the color of the nearest pixel. Produces crisp edges for
	// pixel art at integer scales, but may look jagged at non-integer scales.
	Nearest Filter = ebiten.FilterNearest
	// Linear interpolates between neighbouring pixels. Produces smooth,
	// slightly blurred edges — good for natural images.
	Linear Filter = ebiten.FilterLinear
	// Pixelated is like [Nearest] but stays crisp even at non-integer scales.
	// Best for pixel art games where the window can be any size.
	Pixelated Filter = ebiten.FilterPixelated
)

// BlendMode is a blend mode for [Blend].
type BlendMode = ebiten.Blend

var (
	// BlendNormal is standard alpha blending (source over destination).
	BlendNormal BlendMode = ebiten.BlendSourceOver
	// BlendCopy overwrites the destination with the source.
	BlendCopy BlendMode = ebiten.BlendCopy
	// BlendAdd is additive blending (source added to destination).
	BlendAdd BlendMode = ebiten.BlendLighter
	// BlendErase clears the destination.
	BlendErase BlendMode = ebiten.BlendClear
)

func blitOpts(dr, sr Rectangle, ops []BlitOp) *ebiten.DrawImageOptions {
	sw, sh := float64(sr.Dx()), float64(sr.Dy())
	dw, dh := float64(dr.Dx()), float64(dr.Dy())

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

	return op
}

// Blit copies the source rectangle sr from src onto the destination rectangle
// dr of dst. Ops are flags: rotation is applied first, then flips.
func Blit(dst, src *Surface, dr, sr Rectangle, ops ...BlitOp) {
	sub := src.SubImage(sr).(*ebiten.Image)
	op := blitOpts(dr, sr, ops)
	dst.DrawImage(sub, op)
}

// Blend copies sr from src onto dr of dst with the given blend mode.
// Ops are the same rotation/flip flags as Blit.
func Blend(dst, src *Surface, dr, sr Rectangle, mode BlendMode, ops ...BlitOp) {
	sub := src.SubImage(sr).(*ebiten.Image)
	op := blitOpts(dr, sr, ops)
	op.Blend = mode
	dst.DrawImage(sub, op)
}

// Scale draws src onto dst scaled by sx and sy.
func Scale(dst, src *Surface, sx, sy float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sx, sy)
	dst.DrawImage(src, op)
}
