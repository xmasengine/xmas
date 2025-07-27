package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xmasengine/xmas/engine"
)

func main() {
	prof := ""
	pmem := ""
	flag.StringVar(&prof, "P", "", "pprof profile file")
	flag.StringVar(&pmem, "M", "", "memory profile file")
	flag.Parse()

	if prof != "" {
		pout, err := os.Create(prof)
		if err != nil {
			os.Exit(1)
		}
		pprof.StartCPUProfile(pout)
		defer pprof.StopCPUProfile()
	}

	if pmem != "" {
		mout, err := os.Create(pmem)
		if err != nil {
			os.Exit(1)
		}
		defer func() {
			prof := pprof.Lookup("allocs")
			prof.WriteTo(mout, 1)
			mout.Close()
		}()
	}

	sw, sh := ebiten.Monitor().Size()
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowTitle("xmas: Xmas Game Engine.")
	en := engine.New(sw, sh)
	if err := ebiten.RunGame(en); err != nil {
		log.Fatal(err)
	}
}
