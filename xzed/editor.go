package xzed

import (
	"fmt"
	"image"
	"io/fs"
	"os"
)

import (
	"github.com/xmasengine/xmas/xdat"
	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xui"
)

type Editor struct {
	Layer         xui.Layer // Editor is a widget layer
	Name          string
	Zone          *xdat.Zone
	Camera        image.Rectangle
	Hover         image.Point
	Tile          image.Point // Tile we are hovering
	Cell          xdat.Tile
	Depth         int
	Scale         int
	Error         error
	Message       string
	TileWatcher   *Watcher
	SpriteWatcher *Watcher
	MessageTicks  int
	Fsys          fs.FS
	// Presence      Presence
	// Backup
	// Commander *Tila
}

// var _ xui.Widget = &Editor{}

func (e Editor) ActiveLayer() *xdat.Layer {
	if e.Zone == nil {
		return nil
	}
	if e.Depth < 0 || e.Depth >= len(e.Zone.Layers) {
		return nil
	}
	return &e.Zone.Layers[e.Depth]
}

func (e Editor) Render(screen *xgal.Surface) {
	m := e.ActiveLayer()
	style := e.Layer.Style
	if m != nil {
		if e.Tile.In(image.Rect(0, 0, m.Width-1, m.Height-1)) {
			cr := xgal.Bound(e.Tile.X*m.TileWidth, e.Tile.Y*m.TileHeight,
				m.TileWidth, m.TileHeight).Add(e.Camera.Min)
			style.DrawRect(screen, cr)
		}
	}

	pr := xgal.Bound(e.Hover.X, e.Hover.Y, 100, 20)

	style.Ink(screen, pr, fmt.Sprintf("%s: (%d,%d): %d",
		e.Name, e.Hover.X, e.Hover.Y, e.Cell))
	di := xgal.Pt(0, 12)
	pr = pr.Add(di)
	if e.Error != nil {
		style.Ink(screen, pr, fmt.Sprintf("Error %s", e.Error))
		pr = pr.Add(di)
	}
	if e.Message != "" {
		style.Ink(screen, pr, e.Message)
		pr = pr.Add(di)
	}
	e.Layer.Render(screen)
}

func (e Editor) Place(bounds xgal.Rectangle) (rw, th int) {
	e.Layer.Place(bounds)
	return e.Camera.Dx() / e.Scale, e.Camera.Dy() / e.Scale
}

func (e *Editor) UpdateChoosers() {
	m := e.ActiveLayer()
	if m == nil || m.Texture == nil {
		return
	}
	for _, sub := range e.Layer.Kids {
		if choose, ok := sub.(*xui.ChooserLayer); ok {
			choose.Image = m.Texture
			choose.TileSize = xgal.Pt(m.TileWidth, m.TileHeight)
		}
	}
}

func (e *Editor) LoadSurface(name string) bool {
	m := e.ActiveLayer()
	if m == nil {
		return false
	}
	if e.TileWatcher != nil {
		e.TileWatcher.Done <- struct{}{}
		e.TileWatcher = nil
	}
	e.TileWatcher = Watch(name)
	err := m.SetSource(e.Fsys, name)
	if err != nil {
		e.UpdateChoosers()
	}
	e.Error = err
	// e.Layer.Error(70, 70, 270, 120, err)
	return e.Error == nil
}

func (e *Editor) LoadSpriteSurface(name string) bool {
	/*
		 TODO
			if e.SpriteWatcher != nil {
				e.SpriteWatcher.Done <- struct{}{}
				e.SpriteWatcher = nil
			}
			e.SpriteWatcher = Watch(name)
			err := e.Zone.Sprites.LoadSurface(name)
			if err != nil {
				e.UpdateTilers()
			}
			e.Error = err
			e.Midget.Error(70, 70, 270, 120, err)
			return e.Error == nil
	*/
	return true
}

func (e *Editor) ShowMessage(msg string, args ...any) {
	e.Message = fmt.Sprintf(msg, args...)
	e.MessageTicks = 60 * 15
}

func (e *Editor) UpdateWatcher() bool {
	if e.MessageTicks > 0 {
		e.MessageTicks--
	} else {
		e.Message = ""
		e.Error = nil
	}
	if e.TileWatcher == nil {
		return false
	}
	m := e.ActiveLayer()
	select {
	case name := <-e.TileWatcher.C:
		err := m.SetSource(e.Fsys, name)
		e.Error = err
		if e.Error == nil {
			e.ShowMessage("Auto update tiles: %s", name)
			e.UpdateChoosers()
		}
		return e.Error == nil
	default:
		break
	}
	/*
		if e.SpriteWatcher != nil {
			select {
			case name := <-e.SpriteWatcher.C:
				err := e.Zone.Sprites.LoadSurface(name)
				e.Error = err
				e.Midget.Error(70, 70, 270, 120, err)
				if e.Error == nil {
					e.ShowMessage("Auto update sprites: %s", name)
					e.UpdateTilers()
				}
				return e.Error == nil
			default:
				break
			}
		}
	*/
	return false
}

func (e *Editor) TileSelected(x, y int) {
	m := e.ActiveLayer()
	if m == nil {
		return
	}

	idx := x + y*255

	e.Cell = xdat.Tile(max(0, idx))
}

func (e *Editor) SpriteSelected(x, y int) {
	/*
		_, h := e.Zone.Surface.Size()
		idx := x + y*(h/e.Zone.Th)
		e.Presence.Offset = max(0, idx)
	*/
}

func (e *Editor) SaveZone(name string) bool {
	err := e.Zone.SaveFile(name)
	e.Error = err
	if e.Error == nil {
		e.Name = name
		e.ShowMessage("Zone saved to %s", name)
		return true
	}
	return false
}

func (e *Editor) LoadZone(name string) bool {
	m, err := xdat.LoadZone(e.Fsys, name)
	e.Error = err
	if e.Error == nil {
		e.Zone = m
		e.UpdateChoosers()
		e.ShowMessage("Zone loaded from %s", name)
		e.Name = name
		return true
	}
	return false
}

func (e *Editor) SetDone(done bool) {
	e.Layer.Done = done
}

func (e Editor) FloodFill(at xgal.Point, cell xdat.Tile) {
	m := e.ActiveLayer()
	if m == nil {
		return
	}
	m.FloodFill(at, cell)
}

const HELP = `HELP
Mouse: Draw, select, drag pop up panes.
Mouse Wheel: Select tile index.
Left Shift+Click: Draw image.
Left Control+Click: Draw flag.
Left Control+Alt: Flood fill.
Pause: Exit without save.
F1: This help.          | F2: Save map.
F3: Show tile selector. | F4: Load map.
F5: Export as basic.    | P: Edit Prefix.
F:  Load tile image.    | M: Toggle flag mode.
H: Horizontal flip      | V: Vertical flip
Y: Yank hovered tile.   | G: Edit flags.
Enter: Confirm dialogs. | Esc: Cancel dialogs.
`

func (e *Editor) Update() error {
	var err error
	e.Hover = xgal.Mouse()
	e.Tile = e.Zone.ToTile(e.Hover, e.Camera)

	_, wheel := ebiten.Wheel()
	if wheel > 0 {
		e.Cell.Index++
	} else if wheel < 0 {
		e.Cell.Index = max(0, e.Cell.Index-1)
	}

	e.UpdateWatcher()

	err = e.Midget.Update()
	if err != nil {
		if err == MidgetOK { // input handled by some active Midget.
			return nil
		}
		return err
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyPause):
		e.Midget.YesNo(50, 50, 250, 100, "Quit", "Y", e.SetDone)
	case inpututil.IsKeyJustPressed(ebiten.KeyY):
		e.Cell = e.Zone.Get(e.Tile)
		e.ShowMessage("Yanked %d %d", e.Cell.Index, e.Cell.Flag)
	case inpututil.IsKeyJustPressed(ebiten.KeyL):
		if e.Zone != nil {
			e.Zone.Flags = !e.Zone.Flags
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyH):
		e.Cell.Flag ^= FlagHorizontalFlip
	case inpututil.IsKeyJustPressed(ebiten.KeyV):
		e.Cell.Flag ^= FlagVerticalFlip
	case inpututil.IsKeyJustPressed(ebiten.KeyN):
		e.Cell.Flag ^= FlagOnTop
	case inpututil.IsKeyJustPressed(ebiten.KeyB):
		e.Cell.Flag ^= FlagSolid
	case inpututil.IsKeyJustPressed(ebiten.KeyG):
		e.Midget.AskText(50, 50, 250, 100, "Flag", &e.Cell.Flag)
	case inpututil.IsKeyJustPressed(ebiten.KeyF1):
		e.Midget.Ask(50, 0, 300, 250, HELP, "", Accept)
	case inpututil.IsKeyJustPressed(ebiten.KeyF2):
		e.Midget.Ask(50, 50, 250, 100, "Save As", e.Name, e.SaveZone)
	case inpututil.IsKeyJustPressed(ebiten.KeyF4):
		e.Midget.Ask(50, 50, 250, 100, "Load From", e.Name, e.LoadZone)
	case inpututil.IsKeyJustPressed(ebiten.KeyU):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Backup.Commit(e.SaveZoneToFile)
		} else {
			e.Midget.YesNo(50, 50, 250, 100, "Restore backup", "Y", e.Restore)
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyF):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Midget.Ask(50, 50, 250, 100, "Sprites", e.Zone.Sprites.From, e.LoadSpriteSurface)
		} else {
			e.Midget.Ask(50, 50, 250, 100, "From", e.Zone.From, e.LoadSurface)
		}

	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		e.Midget.AskString(50, 50, 250, 100, "Prefix", &e.Zone.Prefix)
	case inpututil.IsKeyJustPressed(ebiten.KeyO):
		e.Midget.AskInt(50, 50, 250, 100, "Offset", &e.Zone.Offset)
	case inpututil.IsKeyJustPressed(ebiten.KeyS):
		e.Midget.AskInt(50, 50, 250, 100, "UI Scale", &e.Scale)
	case inpututil.IsKeyJustPressed(ebiten.KeyF3):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			tiler := e.Midget.Tile(200, 100, e.Zone.Sprites.Surface, e.SpriteSelected)
			tiler.SetCaption("Sprite")
		} else {
			tiler := e.Midget.Tile(200, 100, e.Zone.Surface, e.TileSelected)
			tiler.SetCaption("Tile")
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyF5):
		e.ExportBasic()
	case inpututil.IsKeyJustPressed(ebiten.KeyF6):
		e.Midget.AskCommand(10, 10, 300, 250, "Command", e.Commander)
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Zone.PutIndex(e.Tile, e.Cell.Index)
		} else if inpututil.KeyPressDuration(ebiten.KeyControlLeft) > 0 || e.Zone.Flags {
			e.Zone.PutFlag(e.Tile, e.Cell.Flag)
		} else if inpututil.KeyPressDuration(ebiten.KeyAltLeft) > 0 {
			e.Zone.FloodFill(e.Tile, e.Cell)
		} else {
			e.Zone.Put(e.Tile, e.Cell)
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Zone.PutIndex(e.Tile, 0)
		} else if inpututil.KeyPressDuration(ebiten.KeyControlLeft) > 0 || e.Zone.Flags {
			e.Zone.PutFlag(e.Tile, 0)
		} else {
			zero := Cell{}
			e.Zone.Put(e.Tile, zero)
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle):
		e.Zone.PutPresence(e.Tile, e.Presence)
	default:
	}

	if e.Midget.Done {
		return Termination
	}
	return nil
}

func (e *Editor) Wrap(t *Tila, args ...any) any {
	if dx, err := TilaArg[int](args); err != nil {
		return err
	} else {
		if e.Zone != nil {
			e.Zone.Wrap(dx)
			return dx
		}
		return false
	}
}

func (e *Editor) Roll(t *Tila, args ...any) any {
	if dx, err := TilaArg[int](args); err != nil {
		return err
	} else {
		if e.Zone != nil {
			e.Zone.Roll(dx)
			return dx
		}
		return false
	}
}

func (e *Editor) CommandHelp(t *Tila, args ...any) any {
	return "available commands: get, set, wrap, roll, help"
}

func NewEditor(tm *Zone, name string, w, h, scale int) *Editor {

	e := &Editor{Zone: tm, Name: name, Camera: image.Rect(0, 0, w, h),
		Scale:  scale,
		Midget: MakeMidget(image.Rect(0, 0, 0, 0)),
	}
	e.Midget.Lock = true
	if tm.From != "" {
		e.TileWatcher = Watch(tm.From)
	}
	if tm.Sprites.From != "" {
		e.SpriteWatcher = Watch(tm.Sprites.From)
	}
	e.Backup.Pattern = "masite*.xml"
	e.Commander = NewTila()
	e.Commander.Commands["get"] = (*Tila).Get
	e.Commander.Operators["$"] = (*Tila).Get
	e.Commander.Commands["set"] = (*Tila).Set
	e.Commander.Commands["wrap"] = e.Wrap
	e.Commander.Commands["roll"] = e.Roll
	e.Commander.Commands["help"] = e.CommandHelp

	return e
}
