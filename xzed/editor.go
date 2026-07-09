package xzed

import (
	"fmt"
	"image"
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
	Scale         int
	Error         error
	Message       string
	TileWatcher   *Watcher
	SpriteWatcher *Watcher
	MessageTicks  int
	// Presence      Presence
	// Backup
	// Commander *Tila
}

// var _ xui.Widget = &Editor{}

func (e Editor) Render(screen *xgal.Surface) {
	if e.Map != nil {
		e.Map.Render(screen, e.Camera)
		if e.Tile.In(image.Rect(0, 0, e.Map.Width-1, e.Map.Height-1)) {
			e.Midget.Style.DrawCursor(screen,
				Bounds(e.Tile.X*e.Map.Tw, e.Tile.Y*e.Map.Th, e.Map.Tw, e.Map.Th).Add(e.Camera.Min),
			)
		}
	}

	kl := len(e.Midget.Kids)
	y := 10
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: (%d,%d): %d %s %d",
		e.Name, e.Hover.X, e.Hover.Y, e.Cell.Index, e.Cell.Flag, kl), e.Map.Width*e.Map.Tw, y)
	y += 12
	if e.Error != nil {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Error: %s [%d]", e.Error, kl),
			e.Map.Width*e.Map.Tw, y)
		y += 12
	}
	if e.Message != "" {
		ebitenutil.DebugPrintAt(screen, e.Message, e.Map.Width*e.Map.Tw, y)
		y += 12
	}
	e.Midget.Draw(screen)
}

func (e Editor) Layout(w, h int) (rw, th int) {
	e.Midget.Layout(w, h)
	return e.Camera.Dx() / e.Scale, e.Camera.Dy() / e.Scale
}

func (e *Editor) UpdateTilers() {
	if e.Map == nil || e.Map.Surface == nil {
		return
	}
	for _, sub := range e.Midget.Kids {
		if tiler, ok := sub.(*Tiler); ok {
			tiler.Surface = e.Map.Surface
		}
	}
}

func (e *Editor) LoadSurface(name string) bool {
	if e.TileWatcher != nil {
		e.TileWatcher.Done <- struct{}{}
		e.TileWatcher = nil
	}
	e.TileWatcher = Watch(name)
	err := e.Map.LoadSurface(name)
	if err != nil {
		e.UpdateTilers()
	}
	e.Error = err
	e.Midget.Error(70, 70, 270, 120, err)
	return e.Error == nil
}

func (e *Editor) LoadSpriteSurface(name string) bool {
	if e.SpriteWatcher != nil {
		e.SpriteWatcher.Done <- struct{}{}
		e.SpriteWatcher = nil
	}
	e.SpriteWatcher = Watch(name)
	err := e.Map.Sprites.LoadSurface(name)
	if err != nil {
		e.UpdateTilers()
	}
	e.Error = err
	e.Midget.Error(70, 70, 270, 120, err)
	return e.Error == nil
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
	select {
	case name := <-e.TileWatcher.C:
		err := e.Map.LoadSurface(name)
		e.Error = err
		e.Midget.Error(70, 70, 270, 120, err)
		if e.Error == nil {
			e.ShowMessage("Auto update tiles: %s", name)
			e.UpdateTilers()
		}
		return e.Error == nil
	default:
		break
	}
	if e.SpriteWatcher != nil {
		select {
		case name := <-e.SpriteWatcher.C:
			err := e.Map.Sprites.LoadSurface(name)
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
	return false
}

func (e *Editor) TileSelected(x, y int) {
	_, h := e.Map.Surface.Size()
	idx := x + y*(h/e.Map.Th)
	if idx > 255 {
		idx -= 255
		e.Cell.Flag |= FlagExtended
	}
	e.Cell.Index = byte(max(0, idx))
}

func (e *Editor) SpriteSelected(x, y int) {
	_, h := e.Map.Surface.Size()
	idx := x + y*(h/e.Map.Th)
	e.Presence.Offset = max(0, idx)
}

func (e *Editor) SaveMap(name string) bool {
	err := e.Map.Save(name)
	e.Error = err
	e.Midget.Error(70, 70, 270, 120, err)
	if e.Error == nil {
		e.Name = name
		e.ShowMessage("Map saved to %s", name)
		return true
	}
	return false
}

func (e *Editor) ExportBasic() bool {
	name := e.Name + ".bas"
	err := e.Map.Export(name, BasicFormat)
	e.Error = err
	e.Midget.Error(70, 70, 270, 120, err)
	if e.Error == nil {
		e.ShowMessage("Exported to %s", name)
		return true
	}
	return false
}

func (e *Editor) SaveMapToFile(f *os.File) error {
	err := e.Map.SaveToFile(f)
	e.Error = err
	if e.Error == nil {
		e.ShowMessage("Map backed up to %s", f.Name())
		return nil
	}
	return err
}

func (e *Editor) LoadMapFromFile(f *os.File) error {
	m, err := LoadMapFromFile(f)
	e.Error = err
	if e.Error == nil {
		e.Map = m
		e.UpdateTilers()
		e.ShowMessage("Map restored from %s", f.Name())
		return nil
	}
	return err
}

func (e *Editor) LoadMap(name string) bool {
	m, err := LoadMap(name)
	e.Error = err
	e.Midget.Error(70, 70, 270, 120, err)
	if e.Error == nil {
		e.Map = m
		e.UpdateTilers()
		e.ShowMessage("Map loaded from %s", name)
		e.Name = name
		return true
	}
	return false
}

func (e *Editor) Restore(doit bool) {
	if doit {
		e.Backup.Restore(e.LoadMapFromFile)
	}
}

func (e *Editor) SetDone(done bool) {
	e.Midget.Done = done
}

func (e Editor) FloodFill(at Point, cell Cell) {
	e.Map.FloodFill(at, cell)
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
	e.Hover = image.Pt(ebiten.CursorPosition())
	e.Tile = e.Map.ToTile(e.Hover, e.Camera)

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
		e.Cell = e.Map.Get(e.Tile)
		e.ShowMessage("Yanked %d %d", e.Cell.Index, e.Cell.Flag)
	case inpututil.IsKeyJustPressed(ebiten.KeyL):
		if e.Map != nil {
			e.Map.Flags = !e.Map.Flags
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
		e.Midget.Ask(50, 50, 250, 100, "Save As", e.Name, e.SaveMap)
	case inpututil.IsKeyJustPressed(ebiten.KeyF4):
		e.Midget.Ask(50, 50, 250, 100, "Load From", e.Name, e.LoadMap)
	case inpututil.IsKeyJustPressed(ebiten.KeyU):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Backup.Commit(e.SaveMapToFile)
		} else {
			e.Midget.YesNo(50, 50, 250, 100, "Restore backup", "Y", e.Restore)
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyF):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Midget.Ask(50, 50, 250, 100, "Sprites", e.Map.Sprites.From, e.LoadSpriteSurface)
		} else {
			e.Midget.Ask(50, 50, 250, 100, "From", e.Map.From, e.LoadSurface)
		}

	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		e.Midget.AskString(50, 50, 250, 100, "Prefix", &e.Map.Prefix)
	case inpututil.IsKeyJustPressed(ebiten.KeyO):
		e.Midget.AskInt(50, 50, 250, 100, "Offset", &e.Map.Offset)
	case inpututil.IsKeyJustPressed(ebiten.KeyS):
		e.Midget.AskInt(50, 50, 250, 100, "UI Scale", &e.Scale)
	case inpututil.IsKeyJustPressed(ebiten.KeyF3):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			tiler := e.Midget.Tile(200, 100, e.Map.Sprites.Surface, e.SpriteSelected)
			tiler.SetCaption("Sprite")
		} else {
			tiler := e.Midget.Tile(200, 100, e.Map.Surface, e.TileSelected)
			tiler.SetCaption("Tile")
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyF5):
		e.ExportBasic()
	case inpututil.IsKeyJustPressed(ebiten.KeyF6):
		e.Midget.AskCommand(10, 10, 300, 250, "Command", e.Commander)
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Map.PutIndex(e.Tile, e.Cell.Index)
		} else if inpututil.KeyPressDuration(ebiten.KeyControlLeft) > 0 || e.Map.Flags {
			e.Map.PutFlag(e.Tile, e.Cell.Flag)
		} else if inpututil.KeyPressDuration(ebiten.KeyAltLeft) > 0 {
			e.Map.FloodFill(e.Tile, e.Cell)
		} else {
			e.Map.Put(e.Tile, e.Cell)
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight):
		if inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 0 {
			e.Map.PutIndex(e.Tile, 0)
		} else if inpututil.KeyPressDuration(ebiten.KeyControlLeft) > 0 || e.Map.Flags {
			e.Map.PutFlag(e.Tile, 0)
		} else {
			zero := Cell{}
			e.Map.Put(e.Tile, zero)
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle):
		e.Map.PutPresence(e.Tile, e.Presence)
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
		if e.Map != nil {
			e.Map.Wrap(dx)
			return dx
		}
		return false
	}
}

func (e *Editor) Roll(t *Tila, args ...any) any {
	if dx, err := TilaArg[int](args); err != nil {
		return err
	} else {
		if e.Map != nil {
			e.Map.Roll(dx)
			return dx
		}
		return false
	}
}

func (e *Editor) CommandHelp(t *Tila, args ...any) any {
	return "available commands: get, set, wrap, roll, help"
}

func NewEditor(tm *Map, name string, w, h, scale int) *Editor {

	e := &Editor{Map: tm, Name: name, Camera: image.Rect(0, 0, w, h),
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
