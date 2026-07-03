package xeng

import (
	"github.com/xmasengine/xmas/xdat"
	"github.com/xmasengine/xmas/xgal"
)

type Zone struct {
	// static data of the zone
	xdat.Zone
}

/*
func (m *Zone) RenderPresences(screen *xgal.Surface, camera xgal.Rectangle) {

	starty := camera.Min.Y / m.Th
	if starty < 0 {
		starty = 0
	}

	for _, presence := range m.Presences {

		atx := presence.X - camera.Min.X
		aty := presence.Y - camera.Min.Y

		if m.Sprites.Surface == nil || m.Flags {
			to := Bounds(atx, aty, m.Tw, m.Th)
			FillRect(screen, to, presenceColor)
			// draw colored rectangle if sprites are not available
			// or if flags mode is set
		}

		if m.Sprites.Surface == nil {
			continue
		}

		aty = aty - presence.Height + FeetHeight
		// "shift up so the "feet" stand on the position of the presence.
		ab := m.Sprites.Surface.Bounds()
		tilew := ab.Dx() / m.Tw
		id := presence.Offset
		idx := id % tilew
		idy := id / tilew
		fx := idx * m.Tw
		fy := idy * m.Th
		from := image.Rect(fx, fy, fx+presence.Width, fy+presence.Height)
		sub := m.Sprites.Surface.SubImage(from).(*Surface)
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(atx), float64(aty))
		if sub != nil {
			screen.DrawImage(sub, &opts)
		}
	}
}
*/

func (z *Zone) RenderLayer(screen *xgal.Surface, camera xgal.Rectangle, idx int) {
	ab := m.Surface.Bounds()
	if idx < 0 || idx > len(z.Layers) {
		return
	}
	m := z.Layers[idx]

	starty := camera.Min.Y / m.Th
	if starty < 0 {
		starty = 0
	}
	endy := min(camera.Max.Y/m.Th, len(m.Rows)-1)

	// This draws the whole layer. Only draw visible part using a camera.
	for ty := starty; ty < endy; ty++ {
		row := m.Rows[ty]

		startx := max(camera.Min.X/m.Tw, 0)
		endx := min(camera.Max.X/m.Tw, len(row.Cells)-1)
		for tx := startx; tx < endx; tx++ {
			cell := row.Cells[tx]
			id := int(cell.Index)
			if cell.Flag&FlagExtended != 0 {
				id += 255
			}
			tilew := ab.Dx() / m.Tw
			idx := id % tilew
			idy := id / tilew
			fx := idx * m.Tw
			fy := idy * m.Th

			from := image.Rect(fx, fy, fx+m.Tw, fy+m.Th)
			sub := m.Surface.SubImage(from).(*Surface)
			opts := ebiten.DrawImageOptions{}
			if cell.Flag&FlagHorizontalFlip != 0 {
				opts.GeoM.Scale(-1, 1)
				opts.GeoM.Translate(float64(m.Tw), 0)
			}
			if cell.Flag&FlagVerticalFlip != 0 {
				opts.GeoM.Scale(1, -1)
				opts.GeoM.Translate(0, float64(m.Th))
			}

			atx := int(tx)*m.Tw - camera.Min.X
			aty := int(ty)*m.Th - camera.Min.Y

			opts.GeoM.Translate(float64(atx), float64(aty))

			if sub != nil {
				screen.DrawImage(sub, &opts)
			}

			to := Bounds(atx, aty, m.Tw, m.Th)
			if m.Flags {
				cell.Flag.Render(screen, to)
			}
		}
	}
	m.RenderPresences(screen, camera)
}

/*
func (m *Zone) FloodFill(atTile Point, cell Cell) {
	now := m.Get(atTile)
	if now.Index == cell.Index && now.Flag == cell.Flag {
		return // already ok
	}
	if !m.Inside(atTile) {
		return
	}

	m.Put(atTile, cell)
	// the floodfill is recursive but the maps are small so
	// it should not cause problems.
	for dx := -1; dx <= 1; dx++ {
		at2 := atTile
		at2.X += dx
		now2 := m.Get(at2)
		if now2.Index == now.Index {
			m.FloodFill(at2, cell)
		}
	}
	for dy := -1; dy <= 1; dy++ {
		at2 := atTile
		at2.Y += dy
		now2 := m.Get(at2)
		if now2.Index == now.Index {
			m.FloodFill(at2, cell)
		}
	}
}
*/
