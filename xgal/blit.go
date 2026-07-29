package xgal

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// BlitOpts specifies the options such as flip or rotation for Blit.
type BlitOpts struct {
	FlipH bool
	FlipV bool
	Rot   Rot
}

// Rot is the stepwise rotation in 90 degree steps.
type Rot uint8

const (
	Rot0   Rot = iota // Do not roate.
	Rot90             // Rotate 90 degrees clockwise.
	Rot180            // Rotate 180 degrees clockwise.
	Rot270            // Rotate 270 degrees clockwise.
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

func (opts BlitOpts) toDrawImageOptions(dr, sr Rectangle) *ebiten.DrawImageOptions {
	sw, sh := float64(sr.Dx()), float64(sr.Dy())
	dw, dh := float64(dr.Dx()), float64(dr.Dy())

	op := &ebiten.DrawImageOptions{}

	flipX, flipY := 1.0, 1.0

	if opts.FlipH {
		flipX = -flipX
	}
	if opts.FlipV {
		flipY = -flipY
	}
	step := int(opts.Rot) % 4

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
	if dw != w && dh != h {
		op.GeoM.Scale(dw/w, dh/h)
	}
	op.GeoM.Translate(float64(dr.Min.X), float64(dr.Min.Y))

	return op
}

func drawImageOptions(dr, sr Rectangle) *ebiten.DrawImageOptions {
	sw, sh := float64(sr.Dx()), float64(sr.Dy())
	dw, dh := float64(dr.Dx()), float64(dr.Dy())

	op := &ebiten.DrawImageOptions{}

	// Scale and position
	if dw != sw && dh != sh {
		op.GeoM.Scale(dw/sw, dh/sh)
	}
	op.GeoM.Translate(float64(dr.Min.X), float64(dr.Min.Y))

	return op
}

// Blit copies the source rectangle sr from src onto the destination rectangle
// dr of dst. Ops are flags: rotation is applied first, then flips.
func Blit(dst, src *Surface, dr, sr Rectangle, ops ...BlitOpts) {
	sub := src.SubImage(sr).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	if len(ops) >= 1 {
		op = ops[0].toDrawImageOptions(dr, sr)
	} else {
		op = drawImageOptions(dr, sr)
	}

	dst.DrawImage(sub, op)
}

// Blend copies sr from src onto dr of dst with the given blend mode.
// Ops are the same rotation/flip flags as Blit.
func Blend(dst, src *Surface, dr, sr Rectangle, mode BlendMode, ops ...BlitOpts) {
	sub := src.SubImage(sr).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	if len(ops) >= 1 {
		op = ops[0].toDrawImageOptions(dr, sr)
	}

	op.Blend = mode
	dst.DrawImage(sub, op)
}

// Scale draws src onto dst scaled by sx and sy.
func Scale(dst, src *Surface, sx, sy float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sx, sy)
	dst.DrawImage(src, op)
}

// Zoom draws src onto dst, scaled uniformely and translated
func Zoom(dst, src *Surface, x, y, zoom float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(zoom), float64(zoom))
	op.GeoM.Translate(x, y)
	op.Filter = ebiten.FilterNearest
	dst.DrawImage(src, op)
}

// NineSlice draws a 9-slice image with one draw call, and custom positioning.
// src: The source texture
// offsetX, offsetY: The screen position where the 9-slice should be drawn
// dstW, dstH: The target dimensions on the screen
// left, right, top, bottom: The inset margins of the 9-slice in pixels
func NineSlice(screen *ebiten.Image, src *ebiten.Image, offsetX, offsetY float32, dstW, dstH int, left, right, top, bottom int) {
	w, h := src.Bounds().Dx(), src.Bounds().Dy()

	// 1. Define coordinate arrays for Source (Texture Space)
	srcX := [4]float32{0, float32(left), float32(w - right), float32(w)}
	srcY := [4]float32{0, float32(top), float32(h - bottom), float32(h)}

	// 2. Define coordinate arrays for Destination (Screen Space) with translation applied
	dstX := [4]float32{offsetX, offsetX + float32(left), offsetX + float32(dstW-right), offsetX + float32(dstW)}
	dstY := [4]float32{offsetY, offsetY + float32(top), offsetY + float32(dstH-bottom), offsetY + float32(dstH)}

	// 3. Allocate vertex array on the stack
	var vertices [16]ebiten.Vertex

	// 4. Populate vertices
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			vertices[y*4+x] = ebiten.Vertex{
				DstX:   dstX[x],
				DstY:   dstY[y],
				SrcX:   srcX[x],
				SrcY:   srcY[y],
				ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
			}
		}
	}

	// 5. Allocate index array on the stack
	var indices [54]uint16
	idx := 0

	// 6. Generate indices linking the grid into triangles
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			topLeft := uint16(y*4 + x)
			topRight := topLeft + 1
			bottomLeft := topLeft + 4
			bottomRight := bottomLeft + 1

			indices[idx] = topLeft
			indices[idx+1] = topRight
			indices[idx+2] = bottomLeft

			indices[idx+3] = topRight
			indices[idx+4] = bottomRight
			indices[idx+5] = bottomLeft
			idx += 6
		}
	}

	// 7. Execute draw call
	op := &ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices[:], indices[:], src, op)
}
