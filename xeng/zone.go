package xeng

import (
	"github.com/xmasengine/xmas/xdat"
	"github.com/xmasengine/xmas/xgal"
)

type Layer struct {
	Source *xgal.Surface
}

type Zone struct {
	// data of the zone
	*xdat.Zone
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

func (z *Zone) RenderLayer(screen *xgal.Surface, camera xgal.Rectangle, m xdat.Layer, index int) {
	if m.Texture == nil {
		// Can't draw if there is no texture loaded.
		return
	}

	starty := camera.Min.Y / int(m.TileHeight)
	if starty < 0 {
		starty = 0
	}
	endy := min(camera.Max.Y/int(m.TileHeight), len(m.Tiles.Rows)-1)

	// This draws the whole layer. Only draw visible part using a camera.
	for ty := starty; ty < endy; ty++ {
		row := m.Tiles.Rows[ty]

		startx := max(camera.Min.X/int(m.TileWidth), 0)
		endx := min(1+camera.Max.X/int(m.TileWidth), len(row)-1)
		for tx := startx; tx < endx; tx++ {
			cell := row[tx]
			if cell == 0 && index > 0 {
				continue // 0 is empty when not level 0
			}
			idx := cell.X()
			idy := cell.Y()
			fx := int(idx) * int(m.TileWidth)
			fy := int(idy) * int(m.TileHeight)

			from := xgal.Rect(fx, fy, fx+int(m.TileWidth), fy+int(m.TileHeight))
			sub := m.Texture.SubImage(from).(*xgal.Surface)
			opts := xgal.BlitOpts{}

			if cell.Has(xdat.FlagHorizontal) {
				opts.FlipH = true
			}
			if cell.Has(xdat.FlagVertical) {
				opts.FlipV = true
			}
			if cell.Has(xdat.FlagRotate90) {
				opts.Rot = xgal.Rot90
			}
			if cell.Has(xdat.FlagRotate180) {
				opts.Rot = xgal.Rot180
			}
			if cell.Has(xdat.FlagRotate270) {
				opts.Rot = xgal.Rot270
			}

			atx := int(tx)*int(m.TileWidth) - camera.Min.X
			aty := int(ty)*int(m.TileHeight) - camera.Min.Y
			to := xgal.Rect(atx, aty, atx+int(m.TileWidth), aty+int(m.TileHeight))
			xgal.Blit(screen, sub, to, sub.Bounds(), opts)
		}
	}
	// m.RenderPresences(screen, camera, layer)
}

func (z *Zone) Render(screen *xgal.Surface, camera xgal.Rectangle) {
	for i, layer := range z.Layers {
		z.RenderLayer(screen, camera, layer, i)
	}
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
