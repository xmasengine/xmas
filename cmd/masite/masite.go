// masite is a Master System Tilemap Editor
// It uses ebitengine and a very simplified widget system.
// The maps are drastically simplified.
// A map can only have one layer and one tile image,
// however the editor can set all extended flags for the SMS.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func errExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func main() {
	var err error
	name := ""
	from := ""
	w := 32
	h := 24
	scale := 2
	flag.StringVar(&name, "m", "", "map input file name")
	flag.IntVar(&w, "w", h, "map width for new map")
	flag.IntVar(&h, "h", h, "map height for new map")
	flag.IntVar(&scale, "S", scale, "ui scale factor")
	flag.StringVar(&from, "f", "", "tile source for new map")

	flag.Parse()

	var tm *Map

	if from != "" {
		tm, err = NewMap(w, h, from)
		errExit(err)
	} else {
		tm, err = LoadMap(name)
		errExit(err)
	}

	sw, sh := ebiten.Monitor().Size()
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowTitle("mashite")
	edit := NewEditor(tm, name, sw, sh, scale)
	if err := ebiten.RunGame(edit); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}
