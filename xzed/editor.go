package xzed

import (
	"fmt"
	"image"
	"io/fs"
	"log/slog"
	//	"os"
)

import (
	"github.com/xmasengine/xmas/xdat"
	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xlui"
)

type Editor struct {
	Layer         *xlui.Layer // Layer is back pointer to the layer this data is kept in.
	Name          string
	Zone          *xdat.Zone
	Camera        *xgal.Rectangle
	Over          image.Point // Over which tile the mouse is ahovering
	Cell          xdat.Tile
	Depth         int
	Scale         int
	Error         error
	Message       string
	TileWatcher   *Watcher
	SpriteWatcher *Watcher
	MessageTicks  int
	Fsys          fs.FS
	Choosers      xlui.Stack
	Done          bool
	// Presence      Presence
	// Backup
	// Commander *Tila
}

func NewEditorLayer(zone *xdat.Zone, name string, camera *xgal.Rectangle, scale int) *xlui.Layer {
	e := newEditor(zone, name, camera, scale)
	l := xlui.NewLayer(*camera)
	e.Layer = l
	e.Layer.Lock = true
	l.Data = e
	l.Class.Render = e.Render
	l.Class.Click = e.Click
	l.Class.Hover = e.Hover
	l.Class.Tap = e.Tap
	l.Class.Tick = e.Tick
	l.Class.Wheel = e.Wheel
	return l
}

func newEditor(zone *xdat.Zone, name string, camera *xgal.Rectangle, scale int) *Editor {
	e := &Editor{Zone: zone, Name: name, Camera: camera,
		Scale: scale,
	}

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

func (e *Editor) Render(screen *xgal.Surface) {
	style := e.Layer.Style

	m := e.ActiveLayer()
	if m != nil {
		cr := xgal.Bound(e.Over.X*m.TileWidth, e.Over.Y*m.TileHeight,
			m.TileWidth, m.TileHeight).Add(e.Camera.Min)

		if e.Over.In(image.Rect(0, 0, m.Width-1, m.Height-1)) {
			style.DrawRect(screen, cr)
		}
		pr := cr.Min.Add(xgal.Pt(m.TileWidth, 0))

		style.Print(screen, pr, fmt.Sprintf("%s: (%d,%d): %d",
			e.Name, e.Over.X, e.Over.Y, e.Cell))
	}

	pr := xgal.Pt(0, 0)
	di := xgal.Pt(0, 12)
	pr = pr.Add(di)
	if e.Error != nil {
		style.Print(screen, pr, fmt.Sprintf("Error %s", e.Error))
		pr = pr.Add(di)
	}
	if e.Message != "" {
		style.Print(screen, pr, e.Message)
		pr = pr.Add(di)
	}
}

func (e *Editor) UpdateChoosers() {
	m := e.ActiveLayer()
	if m == nil || m.Texture == nil {
		return
	}
	for _, chooser := range e.Choosers.Layers {
		err := chooser.Class.Set(m.Texture)
		if err != nil {
			slog.Error("set texture", "err", err)
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
	e.Done = done
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

func (e *Editor) Hover(at xgal.Point) xlui.Reply {
	layer := e.ActiveLayer()
	if layer != nil {
		e.Over = layer.ToTile(at, *e.Camera)
		return xlui.Accept
	}
	return xlui.Ignore
}

func (e *Editor) Click(at xgal.Point, button int) xlui.Reply {
	layer := e.ActiveLayer()
	if layer == nil {
		return xlui.Ignore
	}
	if xgal.MouseButton(button) == xgal.MouseButtonLeft {
		layer.Set(e.Over, e.Cell)
	}

	if xgal.MouseButton(button) == xgal.MouseButtonMiddle {
		// e.Zone.PutPresence(e.Tile, e.Presence)
	}

	return xlui.Accept
}

func (e *Editor) Wheel(at xgal.Point, delta int) xlui.Reply {
	if delta > 0 {
		e.Cell++
	} else if delta < 0 {
		e.Cell = max(0, e.Cell-1)
	}
	return xlui.Accept
}

func (e *Editor) Tap(key int, mods xlui.Mods) xlui.Reply {
	switch xgal.KeyCode(key) {
	case xgal.KeyEqual:
		e.Cell++
	case xgal.KeyMinus:
		e.Cell = max(0, e.Cell-1)
	case xgal.KeyPause:
		e.Done = true
		// e.Layer.Ask(50, 50, 250, 100, "Quit", "Y", e.SetDone)
	case xgal.KeyY:
		e.Cell = e.ActiveLayer().Get(e.Over)
		e.ShowMessage("Yanked %d", e.Cell)
	/*
		case xgal.Key(xgal.KeyL):
			if e.Zone != nil {
				e.Zone.Flags = !e.Zone.Flags
			}
		case xgal.Key(xgal.KeyH):
			e.Cell.Flag ^= FlagHorizontalFlip
		case xgal.Key(xgal.KeyV):
			e.Cell.Flag ^= FlagVerticalFlip
		case xgal.Key(xgal.KeyN):
			e.Cell.Flag ^= FlagOnTop
		case xgal.Key(xgal.KeyB):
			e.Cell.Flag ^= FlagSolid
		case xgal.Key(xgal.KeyG):
			e.Layer.AskText(50, 50, 250, 100, "Flag", &e.Cell.Flag)
	*/
	case xgal.KeyF1:
		// e.Layer.Ask(50, 0, 300, 250, HELP, "", Accept)
	case xgal.KeyF2:
		// e.Layer.Ask(50, 50, 250, 100, "Save As", e.Name, e.SaveZone)
	case xgal.KeyF4:
		// e.Layer.Ask(50, 50, 250, 100, "Load From", e.Name, e.LoadZone)
	case xgal.KeyU:
		if mods.Shift {
			// e.Backup.Commit(e.SaveZoneToFile)
		} else {
			// e.Layer.YesNo(50, 50, 250, 100, "Restore backup", "Y", e.Restore)
		}
	case xgal.KeyF:
		if mods.Shift {
			// e.Layer.Ask(50, 50, 250, 100, "Sprites", e.Zone.Sprites.From, e.LoadSpriteSurface)
		} else {
			// e.Layer.Ask(50, 50, 250, 100, "From", e.Zone.From, e.LoadSurface)
		}

	case xgal.KeyP:
		// e.Layer.AskString(50, 50, 250, 100, "Prefix", &e.Zone.Prefix)
	case xgal.KeyO:
		// e.Layer.AskInt(50, 50, 250, 100, "Offset", &e.Zone.Offset)
	case xgal.KeyS:
		// e.Layer.AskInt(50, 50, 250, 100, "UI Scale", &e.Scale)
	case xgal.KeyF3:
		if xgal.Key(xgal.KeyShiftLeft) {
			// choose := e.Layer.Chooser(200, 100, e.Zone.Sprites.Surface, e.SpriteSelected)
			// choose.SetCaption("Sprite")
		} else {
			// choose := e.Layer.Chooser(200, 100, e.Zone.Surface, e.TileSelected)
			// choose.SetCaption("Tile")
		}
	case xgal.KeyF5:

	case xgal.KeyF6:
		// e.Layer.AskCommand(10, 10, 300, 250, "Command", e.Commander)
	default:
		return xlui.Ignore
	}

	return xlui.Accept
}

func (e *Editor) Tick(tick int64) xlui.Reply {
	if e.MessageTicks > 0 {
		e.MessageTicks--
	} else {
		e.Message = ""
		e.Error = nil
	}

	e.UpdateWatcher()
	if e.Done {
		return xlui.Finish
	}
	return xlui.Accept
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
