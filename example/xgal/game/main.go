package main

import (
	"fmt"
	"os"

	"github.com/xmasengine/xmas/xgal"
)

type game struct{ failed, started bool }

func (g *game) Update() error {
	if !g.started {
		g.started = true
		return nil
	}
	return xgal.Quit
}

func (g *game) Draw(screen *xgal.Surface) {
	src := xgal.Prepare(8, 8)
	dst := xgal.Prepare(8, 8)

	red := xgal.RGBA{R: 255, A: 255}

	// Mark pixel (1,1) red.
	src.Set(1, 1, red)

	// identity
	xgal.Blit(dst, src, xgal.Rect(0, 0, 8, 8), xgal.Rect(0, 0, 8, 8))
	if _, _, _, a := dst.At(1, 1).RGBA(); a == 0 {
		fmt.Println("FAIL: identity")
		g.failed = true
		return
	}

	// FlipH: pixel center (1.5,1.5) → 8-1.5 = 6.5 → pixel (6,1)
	dst.Clear()
	xgal.Blit(dst, src, xgal.Rect(0, 0, 8, 8), xgal.Rect(0, 0, 8, 8), xgal.FlipH)
	if _, _, _, a := dst.At(6, 1).RGBA(); a == 0 {
		fmt.Println("FAIL: FlipH")
		g.failed = true
		return
	}

	// FlipV: (1.5,1.5) → (1.5, 6.5) → pixel (1,6)
	dst.Clear()
	xgal.Blit(dst, src, xgal.Rect(0, 0, 8, 8), xgal.Rect(0, 0, 8, 8), xgal.FlipV)
	if _, _, _, a := dst.At(1, 6).RGBA(); a == 0 {
		fmt.Println("FAIL: FlipV")
		g.failed = true
		return
	}

	// Rot180: (1.5,1.5) → (6.5,6.5) → pixel (6,6)
	dst.Clear()
	xgal.Blit(dst, src, xgal.Rect(0, 0, 8, 8), xgal.Rect(0, 0, 8, 8), xgal.Rot180)
	if _, _, _, a := dst.At(6, 6).RGBA(); a == 0 {
		fmt.Println("FAIL: Rot180")
		g.failed = true
		return
	}

	// Sub-rect: blit src[2,2,6,6] into dst[0,0,4,4]
	dst.Clear()
	src.Set(2, 2, xgal.RGBA{G: 255, A: 255})
	xgal.Blit(dst, src, xgal.Rect(0, 0, 4, 4), xgal.Rect(2, 2, 6, 6))
	_, gs, _, _ := dst.At(0, 0).RGBA()
	if gs == 0 {
		fmt.Println("FAIL: sub-rect")
		g.failed = true
		return
	}

	// --- Text ---

	if s := xgal.Stride(xgal.BuiltinFace); s <= 0 {
		fmt.Println("FAIL: Stride")
		g.failed = true
		return
	}

	w, h := xgal.Measure("X", xgal.BuiltinFace, 16)
	if w <= 0 || h <= 0 {
		fmt.Println("FAIL: Measure")
		g.failed = true
		return
	}

	txt := xgal.Prepare(50, 30)
	xgal.Ink(txt, xgal.BuiltinFace, red, 0, 0, "X")
	inkOk := false
	for y := 0; y < 30; y++ {
		for x := 0; x < 50; x++ {
			if _, _, _, a := txt.At(x, y).RGBA(); a > 0 {
				inkOk = true
				break
			}
		}
	}
	if !inkOk {
		fmt.Println("FAIL: Ink")
		g.failed = true
		return
	}

	fmt.Println("PASS")
}

func (g *game) Layout(w, h int) (int, int) { return 8, 8 }

func main() {
	g := &game{}
	xgal.Screen(8, 8, "xgaltest")
	if err := xgal.Play(g); err != nil {
		fmt.Fprintf(os.Stderr, "RUNTIME: %v\n", err)
		os.Exit(2)
	}
	if g.failed {
		os.Exit(1)
	}
}
