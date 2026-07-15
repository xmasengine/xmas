package xzed

import (
	"fmt"
	"image"
	"io/fs"
	//	"os"
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
			e.Layer.Error(70, 70, 270, 120, err)
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
				e.Layer.Error(70, 70, 270, 120, err)
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

func (e *Editor) Poll() xui.Reply {
	// var err error
	e.Hover = xgal.Cursor()
	layer := e.ActiveLayer()
	if layer != nil {
		e.Tile = layer.ToTile(e.Hover, e.Camera)
	}

	_, wheel := xgal.Wheel()
	if wheel > 0 {
		e.Cell++
	} else if wheel < 0 {
		e.Cell = max(0, e.Cell-1)
	}

	e.UpdateWatcher()

	res := e.Layer.Poll()
	if res == xui.Accept {
		return res
	}

	switch {
	case xgal.Tap(xgal.KeyPause):
		// e.Layer.Ask(50, 50, 250, 100, "Quit", "Y", e.SetDone)
	case xgal.Tap(xgal.KeyY):
		e.Cell = layer.Get(e.Tile)
		e.ShowMessage("Yanked %d", e.Cell)
	/*
		case xgal.Tap(xgal.KeyL):
			if e.Zone != nil {
				e.Zone.Flags = !e.Zone.Flags
			}
		case xgal.Tap(xgal.KeyH):
			e.Cell.Flag ^= FlagHorizontalFlip
		case xgal.Tap(xgal.KeyV):
			e.Cell.Flag ^= FlagVerticalFlip
		case xgal.Tap(xgal.KeyN):
			e.Cell.Flag ^= FlagOnTop
		case xgal.Tap(xgal.KeyB):
			e.Cell.Flag ^= FlagSolid
		case xgal.Tap(xgal.KeyG):
			e.Layer.AskText(50, 50, 250, 100, "Flag", &e.Cell.Flag)
	*/
	case xgal.Tap(xgal.KeyF1):
		// e.Layer.Ask(50, 0, 300, 250, HELP, "", Accept)
	case xgal.Tap(xgal.KeyF2):
		// e.Layer.Ask(50, 50, 250, 100, "Save As", e.Name, e.SaveZone)
	case xgal.Tap(xgal.KeyF4):
		// e.Layer.Ask(50, 50, 250, 100, "Load From", e.Name, e.LoadZone)
	case xgal.Tap(xgal.KeyU):
		if xgal.Key(xgal.KeyShiftLeft) {
			// e.Backup.Commit(e.SaveZoneToFile)
		} else {
			// e.Layer.YesNo(50, 50, 250, 100, "Restore backup", "Y", e.Restore)
		}
	case xgal.Tap(xgal.KeyF):
		if xgal.Key(xgal.KeyShiftLeft) {
			// e.Layer.Ask(50, 50, 250, 100, "Sprites", e.Zone.Sprites.From, e.LoadSpriteSurface)
		} else {
			// e.Layer.Ask(50, 50, 250, 100, "From", e.Zone.From, e.LoadSurface)
		}

	case xgal.Tap(xgal.KeyP):
		// e.Layer.AskString(50, 50, 250, 100, "Prefix", &e.Zone.Prefix)
	case xgal.Tap(xgal.KeyO):
		// e.Layer.AskInt(50, 50, 250, 100, "Offset", &e.Zone.Offset)
	case xgal.Tap(xgal.KeyS):
		// e.Layer.AskInt(50, 50, 250, 100, "UI Scale", &e.Scale)
	case xgal.Tap(xgal.KeyF3):
		if xgal.Key(xgal.KeyShiftLeft) {
			// choose := e.Layer.Chooser(200, 100, e.Zone.Sprites.Surface, e.SpriteSelected)
			// choose.SetCaption("Sprite")
		} else {
			// choose := e.Layer.Chooser(200, 100, e.Zone.Surface, e.TileSelected)
			// choose.SetCaption("Tile")
		}
	case xgal.Tap(xgal.KeyF5):

	case xgal.Tap(xgal.KeyF6):
		// e.Layer.AskCommand(10, 10, 300, 250, "Command", e.Commander)
	case xgal.Grip(xgal.MouseButtonLeft):
		if xgal.Key(xgal.KeyShiftLeft) {
			// e.Zone.PutIndex(e.Tile, e.Cell.Index)
		} else if xgal.Key(xgal.KeyControlLeft) {
			// e.Zone.PutFlag(e.Tile, e.Cell.Flag)
		} else if xgal.Key(xgal.KeyAltLeft) {
			// e.Zone.FloodFill(e.Tile, e.Cell)
		} else {
			// e.Zone.Put(e.Tile, e.Cell)
		}
	case xgal.Grip(xgal.MouseButtonRight):
		if xgal.Key(xgal.KeyShiftLeft) {
			// e.Zone.PutIndex(e.Tile, 0)
		} else if xgal.Key(xgal.KeyControlLeft) {
			// e.Zone.PutFlag(e.Tile, 0)
		} else {
			// zero := Cell{}
			// e.Zone.Put(e.Tile, zero)
		}
	case xgal.Grip(xgal.MouseButtonMiddle):
		// e.Zone.PutPresence(e.Tile, e.Presence)
	default:
		return xui.Ignore
	}

	if e.Layer.Done {
		return xui.Finish
	}
	return xui.Accept
}

/*
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
*/

func NewEditor(zone *xdat.Zone, name string, w, h, scale int) *Editor {

	e := &Editor{Zone: zone, Name: name, Camera: image.Rect(0, 0, w, h),
		Scale: scale,
		Layer: xui.MakeLayer(image.Rect(0, 0, w, h)),
	}
	e.Layer.Lock = true
	/*
		if tm.From != "" {
			e.TileWatcher = Watch(tm.From)
		}
		if tm.Sprites.From != "" {
			e.SpriteWatcher = Watch(tm.Sprites.From)
		}
	*/
	/*
		e.Backup.Pattern = "xmas*.xml"
		e.Commander = NewTila()
		e.Commander.Commands["get"] = (*Tila).Get
		e.Commander.Operators["$"] = (*Tila).Get
		e.Commander.Commands["set"] = (*Tila).Set
		e.Commander.Commands["wrap"] = e.Wrap
		e.Commander.Commands["roll"] = e.Roll
		e.Commander.Commands["help"] = e.CommandHelp
	*/

	return e
}
