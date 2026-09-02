package xeng

import (
	"github.com/xmasengine/xmas/xdat"
	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xres/fontres"
)

type Intro struct {
	World *xdat.World
	anim  int
}

func (i *Intro) UpdateMode() (done bool, next Mode) {
	pressed := xgal.Keys()
	for _, k := range pressed {
		switch k {
		case xgal.KeyEscape, xgal.KeyEnter:
			return true, nil
		}
	}
	i.anim++
	return false, nil
}

func (i Intro) DrawMode(screen *xgal.Surface) {
	white := xgal.Paint(255, 255, 255, 255)
	black := xgal.Paint(0, 0, 0, 255)
	xgal.Clear(screen, black)
	cx := 120 - (i.anim * 2)
	if cx < 40 {
		cx = 40
	}
	cy := 40

	xgal.Line(screen, cx, cy, cx, cy+120, 12, white)
	xgal.Line(screen, cx-40, cy+40, cx+40, cy+40, 12, white)

	if i.World == nil {
		return
	}

	tx := cx + 20
	ty := cy + 60
	font := fontres.MediumFace
	_, h := xgal.Measure(i.World.Name, font, float64(xgal.Stride(font)))
	xgal.Print(screen, font, white, tx, ty, i.World.Name)
	ty += int(h)
	xgal.Print(screen, font, white, tx, ty, i.World.Copyright)
	ty += int(h)
	xgal.Print(screen, font, white, tx, ty, "Made with XMASEngine")
	ty += int(h)
	xgal.Print(screen, font, white, tx, ty, "Press key or button...")
}
