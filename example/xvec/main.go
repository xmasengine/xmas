package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xvec"
)

var (
	fileFlag     = flag.String("f", "", "path to .xvec or .svg file")
	widthFlag    = flag.Float64("w", 0, "output canvas width")
	heightFlag   = flag.Float64("h", 0, "output canvas height")
	outFlag      = flag.String("o", "", "write xvec output to file")
)

func main() {
	flag.Parse()

	var x *xvec.XVEC
	var err error

	if *fileFlag != "" {
		f, openErr := os.Open(*fileFlag)
		if openErr != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", *fileFlag, openErr)
			os.Exit(1)
		}
		defer f.Close()

		ext := strings.ToLower(filepath.Ext(*fileFlag))
		if ext == ".svg" {
			x, err = xvec.ParseSVG(f, float32(*widthFlag), float32(*heightFlag))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error converting SVG: %v\n", err)
				os.Exit(1)
			}
		} else {
			x = &xvec.XVEC{}
			if err := x.Decode(f); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", *fileFlag, err)
				os.Exit(1)
			}
		}
	} else {
		x = demo()
	}

	w, h := int(x.Size.W), int(x.Size.H)
	if w == 0 {
		if *widthFlag > 0 {
			w = int(*widthFlag)
		} else {
			w = 320
		}
	}
	if h == 0 {
		if *heightFlag > 0 {
			h = int(*heightFlag)
		} else {
			h = 240
		}
	}

	if *outFlag != "" {
		f, err := os.Create(*outFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating %s: %v\n", *outFlag, err)
			os.Exit(1)
		}
		if err := x.Encode(f); err != nil {
			fmt.Fprintf(os.Stderr, "error encoding: %v\n", err)
			os.Exit(1)
		}
		f.Close()
		os.Exit(0)
	}

	g := &game{x: *x}
	xgal.Screen(w, h, "xvec viewer")
	if err := xgal.Play(g); err != nil {
		fmt.Fprintf(os.Stderr, "RUNTIME: %v\n", err)
		os.Exit(2)
	}
}

type game struct {
	x xvec.XVEC
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *xgal.Surface) {
	surf := xgal.Prepare(int(g.x.Size.W), int(g.x.Size.H))
	g.x.Draw(surf)
	xgal.Blit(screen, surf, surf.Bounds(), surf.Bounds())
}

func (g *game) Layout(w, h int) (int, int) {
	return int(g.x.Size.W), int(g.x.Size.H)
}

func mkcol(r, g, b, a uint8) xvec.Color {
	return xvec.Color{R: r, G: g, B: b, A: a}
}

func demo() *xvec.XVEC {
	x := &xvec.XVEC{Size: xvec.Size{W: 160, H: 120}}
	red := mkcol(255, 0, 0, 255)
	green := mkcol(0, 255, 0, 255)
	blue := mkcol(0, 0, 255, 255)
	white := mkcol(255, 255, 255, 255)

	x.Slab(0, 0, 160, 120, mkcol(30, 30, 50, 255))
	x.Circle(80, 60, 50, 2, white)
	x.Disk(80, 60, 20, red)
	x.Rect(10, 10, 60, 40, 1, green)
	x.Line(0, 0, 160, 120, 1, blue)
	x.Fill(mkcol(0, 200, 200, 100),
		xvec.MoveTo(80, 20),
		xvec.LineTo(140, 60),
		xvec.LineTo(80, 100),
		xvec.LineTo(20, 60),
		xvec.Close(),
	)
	x.Stroke(1, mkcol(255, 255, 0, 255),
		xvec.MoveTo(40, 30),
		xvec.CubicTo(120, 10, 120, 110, 40, 90),
	)
	return x
}
