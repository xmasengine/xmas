// package xzed is a GUI Zone Editor for the xmas engine
package xzed

import "image"

/*
import "image/color"
import "fmt"
import "path"
*/

import (
// "github.com/hajimehoshi/ebiten/v2"
)

import (
	//	"github.com/xmasengine/xmas/xmap"
	// "github.com/xmasengine/xmas/xres"
	"github.com/xmasengine/xmas/xui"
)

// Type aliases for ease of use.
type (
	Style      = xui.Style
	Surface    = xui.Surface
	Context    = xui.Root
	Event      = xui.Event
	MouseEvent = xui.MouseEvent
	KeyEvent   = xui.KeyEvent
	Point      = image.Point
)

/*

// NewEditorGroup returns a zone editor with additional boxes attached to it.
func NewEditorGroup(x, y, w, h int, tm *xmap.Map, camera *geom.Camera, loader MapLoader) (*Editor) {
	e := New(x, y, w, h, tm, camera, loader)
	e.Picker = NewPicker(TILE_W*4, TILE_H*2+10, e.Atlas, 3, e)
	e.Picker.Move(image.Pt(w-TILE_W*6, TILE_H*2))

	e.ChoosePanel = xui.NewPanel(128, 64)
	e.ChoosePanel.AddNewPanelHeader(128, 10, "Choose Tile", e.ChoosePanel)
	e.ChoosePanel.AddNewScrollBar(4, 60)
	e.ChoosePanel.AddNewScrollBar(82, 4)
	e.ChoosePanel.AddNewPanelCorner(4, 4, e.ChoosePanel)

	if e.Atlas != nil && e.Atlas.Main != nil {
		e.Choose = e.ChoosePanel.AddNewPicker(128, 64-22, e.Atlas.Main, TILE_W, TILE_H)
		e.Choose.Move(image.Pt(4, 22))
		e.Choose.Select = e.SetSelected
	}

	e.ChoosePanel.Move(image.Pt(TILE_W*12, TILE_H*5))

	e.Help = xui.NewPanel(TILE_W*9, TILE_H*8)
	e.Help.AddNewPanelHeader(TILE_W*9, 10, "Help", e.Help)
	e.Help.Text.WriteString(helpText)
	e.Help.Move(image.Pt(TILE_W*12, TILE_H))
	e.Help.Hidden = true

	e.AddWidget(&e.Picker.Widget)
	e.AddWidget(&e.ChoosePanel.Widget)
	e.AddWidget(&e.Help.widget)

	return group, e
}

*/

/*
// MapLoader is an interface for a map loader.
type MapLoader interface {
	MapLoad(*xmap.Zone)
}

// Editor is a map editor view.
// It doesn't display the map itself, but keeps a reference to it with a MapLoader.
type Editor struct {
	xui.Box
	Map *xmap.Zone
	// Layer is the layer index we are currently editing.
	Layer int
	// Tile Coordinates of the mouse pointer including camera offset.
	MouseTile image.Point
	// Tile Coordinates of the camera alone.
	CameraTile image.Point

	// Loader is the map loader
	Loader MapLoader

	// Below are references to related gadgets and panels
	Picker    *Picker      // Picker is a reference to the box to display the current tiles.
	Cot       *xui.Cot     // Graphical tile cursor cot.
	Help      *xui.Box     // Help box.
	Choose    *xui.Chooser // Tile Chooser.
	ChooseBox *xui.Box     // Box around the the tile Chooser.
	Hidden    bool
}

// New creates a new map editor. The size should be that of the Surface.
// The editor istelf is invisible but it has activatable child panels.
func New(x, y, w, h int, mz *xmap.Zone, loader MapLoader) *Editor {
	e := &Editor{}
	e.Loader = loader
	e.Panel.Init(w, h)
	e.Panel.Style.Fill = color.RGBA{0, 0, 0, 0}
	e.Panel.Style.Stroke = 0
	e.Map = mz
	e.Layer = 0

	// Set up a cursor which is a box that marks the tile to edit in the map.
	e.Cursor = xui.NewCursor(TILE_W, TILE_H)
	e.Gadgets = append(e.Gadgets, e.Cursor)

	// Set up the atlas.
	e.setupAtlas()

	return e
}

func (e *Editor) setupAtlas() {
	if e.Map != nil && len(e.Map.Layers) > e.Layer && e.Map.Layers[e.Layer].Atlas != nil {
		e.Atlas = e.Map.Layers[e.Layer].Atlas
	} else {
		dprintln("setupAtlas no atlas")
	}

	if e.Picker != nil {
		e.Picker.ChangeAtlas(e.Atlas)
	} else {
		dprintln("setupAtlas no Picker")
	}

	if e.Choose != nil && e.Atlas != nil {
		e.Choose.SetImage(e.Atlas.Main)
	} else {
		dprintln("setupAtlas no Picker or atlas")
	}

}

func (e *Editor) EventHandle(ctx *Context, ev Event) bool {
	if e.Hidden {
		return false
	}

	ok := xui.Dispatch(ctx, e)(ev)
	if ok {
		return true
	}

	return e.Panel.EventHandle(ctx, ev)
}

func (e *Editor) Draw(Surface *Surface) {
	if e.Hidden {
		return
	}

	mct := e.MouseCameraTile()
	msg := fmt.Sprintf("%0d:%0d:%0d", e.Layer, mct.X, mct.Y)
	at := image.Pt(TILE_W, TILE_H)
	if e.Cursor != nil {
		at = at.Add(e.Cursor.Rectangle.Min)
	}
	e.Style.DrawText(Surface, at, msg)
	e.Panel.Draw(Surface)
}

func (b *Editor) drawTileWithMouse(ctx *Context, ev MouseEvent) bool {
	if b.Map != nil && b.Layer >= b.Layer && b.Layer < len(b.Map.Layers) {
		layer := b.Map.Layers[b.Layer]
		mcb := b.MouseCameraTile()

		button := int(ev.Button)
		if (button < 0) || (button >= len(b.Picker.Indexes)) {
			return true
		} else {
			index := b.Picker.Indexes[button]
			layer.SetTile(int(mcb.X), int(mcb.Y), index)
			return true
		}
	}
	return false
}

func (b *Editor) MousePress(ctx *Context, ev MouseEvent) bool {
	if !b.AcceptMouse(ev) {
		return false
	}

	return b.drawTileWithMouse(ctx, ev)
}

func (b *Editor) MouseWheel(ctx *Context, ev MouseEvent) bool {
	if !b.AcceptMouse(ev) {
		return false
	}

	if ev.Wheel.Y > 0 {
		b.NextLayer()
		return true
	} else if ev.Wheel.Y < 0 {
		b.PrevLayer()
		return true
	}
	return false
}

func (b *Editor) SetSelected(at image.Point, which int) {
	b.Picker.SetSelected(at, which)
}

func (b *Editor) AddLayer() bool {
	if b.Map == nil {
		return false
	}
	b.Map.AddNewLayerFromLayer(b.Layer)
	b.setupAtlas()
	return true
}

func (b *Editor) DeleteLayer(ctx *Context) bool {
	if b.Map == nil {
		return false
	}

	onDelete := func(dia *xui.Dialog) bool {
		if dia.Result != 0 {
			return true
		}
		err := b.Map.DeleteLayer(b.Layer)
		if err == nil {
			b.Layer = 0
			b.setupAtlas()
		} else {
			ctx.ErrorDialog(err, "Delete Layer")
		}
		return true
	}
	ctx.Dialog(onDelete, "Delete Layer", "OK", "Cancel")
	return true
}

func (b *Editor) NextLayer() bool {
	if b.Layer < len(b.Map.Layers)-1 {
		b.Layer++
		b.setupAtlas()
		return true
	}
	return false
}

func (b *Editor) PrevLayer() bool {
	if b.Layer > 0 {
		b.Layer--
		b.setupAtlas()
		return true
	}
	return false
}

func (b *Editor) LoadMap(ctx *Context) bool {
	var name string
	var mapNames []string

	onLoad := func(dia *xui.Dialog) bool {
		if dia.Result != 0 {
			return true
		}
		m, err := xmap.Load(amos.Manager().Overlay, path.Join("map", name))
		if err != nil {
			ctx.ErrorDialog(err, "Map Load")
		} else {
			b.Map = m
			b.setupAtlas()
			if b.Loader != nil {
				b.Loader.MapLoad(b.Map)
			}
		}
		return true
	}
	mapFiles, err := amos.Manager().Overlay.ReadDir("map")
	if err != nil {
		ctx.ErrorDialog(err, "Map Dir")
		return false
	}

	mapNames = amos.DirNames(mapFiles...)
	dia := ctx.Dialog(onLoad, "Map Load", "OK", "Cancel")
	dia.StringListWithItems(len(mapNames)*10, &name, mapNames...)

	return true
}

const DefaultMapFileName = "map_000001.cma"

func (b *Editor) SaveMap(ctx *Context) bool {
	name := b.Map.FileName
	if name == "" {
		name = DefaultMapFileName
	}
	var entry *xui.Entry

	onSave := func(dia *xui.Dialog) bool {
		if dia.Result != 0 {
			return true
		}
		if len(entry.Input) > 0 {
			name = string(entry.Input)
			b.Map.FileName = name
		}
		err := b.Map.SaveAs(amos.Manager().Overlay, path.Join("map", name))
		ctx.ErrorDialog(err, "Map Save")
		return true
	}
	dia := ctx.Dialog(onSave, "Map Save", "OK", "Cancel")
	entry = dia.AddEntry(100, 10, name).SetLabel("File")
	// entry.Move(image.Pt(10, 40))
	return true
}

func (b *Editor) ChangeLayer(ctx *Context) bool {
	layer := b.currentLayer()
	if layer == nil {
		return false
	}

	alt := &xmap.Layer{
		Image:    path.Base(layer.Image),
		TileSize: layer.TileSize,
	}
	var tileNames []string

	tileFiles, err := amos.Manager().Overlay.ReadDir("tile")
	if err != nil {
		ctx.ErrorDialog(err, "Tile Dir")
		return false
	}
	tileNames = amos.DirNames(tileFiles...)

	onChange := func(dia *xui.Dialog) bool {
		if dia.Result != 0 {
			return true
		}
		if alt.Image != "" {
			alt.Image = path.Join("tile", alt.Image)
		}
		alt.Size.X = alt.TileSize.X * int(TILE_W)
		alt.Size.Y = alt.TileSize.Y * TILE_H

		err := layer.ChangeLayer(alt)
		if err == nil {
			b.setupAtlas()
		} else {
			ctx.ErrorDialog(err, "Layer Change")
		}
		return true
	}
	dia := ctx.Dialog(onChange, "Change Layer", "OK", "Cancel")
	dia.IntEntry(100, 10, &alt.TileSize.X).SetLabel("Width")
	dia.IntEntry(100, 10, &alt.TileSize.Y).SetLabel("Height")
	dia.StringListWithItems(len(tileNames)*6, &alt.Image, tileNames...)
	return true
}

func (b *Editor) ChangeMap(ctx *Context) bool {
	alt := &xmap.Map{
		FileName:   b.Map.FileName,
		Title:      b.Map.Title,
		Background: b.Map.Background,
		TileSize:   b.Map.TileSize,
	}
	var bgNames []string

	bgFiles, err := amos.Manager().Overlay.ReadDir("background")
	if err != nil {
		ctx.ErrorDialog(err, "Background Dir")
		return false
	}
	bgNames = amos.DirNames(bgFiles...)

	if alt.FileName == "" {
		alt.FileName = DefaultMapFileName
	}

	onChange := func(dia *xui.Dialog) bool {
		if dia.Result != 0 {
			return true
		}
		alt.Size.X = alt.TileSize.X * int(TILE_W)
		alt.Size.Y = alt.TileSize.Y * TILE_H

		err := b.Map.ChangeMap(alt)
		ctx.ErrorDialog(err, "Map Change")
		return true
	}
	dia := ctx.Dialog(onChange, "Change Map", "OK", "Cancel")
	dia.StringEntry(100, 10, &alt.FileName).SetLabel("File")
	dia.StringEntry(100, 10, &alt.Title).SetLabel("Title")
	dia.IntEntry(100, 10, &alt.TileSize.X).SetLabel("Width")
	dia.IntEntry(100, 10, &alt.TileSize.Y).SetLabel("Height")
	dia.StringListWithItems(len(bgNames)*4, &alt.Background, bgNames...)

	// entry.Move(image.Pt(10, 40))
	return true
}

func (b Editor) currentLayer() *xmap.Layer {
	if b.Map == nil || b.Layer < 0 && b.Layer >= len(b.Map.Layers) {
		return nil
	}
	return b.Map.Layers[b.Layer]
}

func (b *Editor) TileChange(ctx *Context, index int) bool {
	layer := b.currentLayer()
	if layer == nil {
		return false
	}

	animation := int(layer.Animation(int(index.X), int(index.Y)))
	kind := int(layer.Kind(int(index.X), int(index.Y)))

	onChange := func(dia *xui.Dialog) bool {
		if dia.Result != 0 {
			return true
		}
		layer.SetAnimation(int(index.X), int(index.Y), xmap.Animation(animation))
		layer.SetKind(int(index.X), int(index.Y), xmap.Kind(kind))
		return true
	}
	dia := ctx.Dialog(onChange, "Change Tile", "OK", "Cancel")
	dia.Resize(image.Pt(10, 20))
	dia.IntEntry(50, 10, &animation).SetLabel("Animation")
	// dia.IntCheckBoxes(100, 10, &animation, "flip", xmap.AnimationPingPong)
	dia.IntCheckBoxes(100, 10, &kind,
		"wall", xmap.KindWall,
		"stair", xmap.KindStair,
		"jump", xmap.KindJump,
		"drop", xmap.KindDrop,
		"north", xmap.KindNorth,
		"east", xmap.KindEast,
		"west", xmap.KindWest,
		"south", xmap.KindSouth,
	)

	return true
}

func (b *Editor) ChangeTile(ctx *Context) bool {
	index := b.Picker.Indexes[0]
	return b.TileChange(ctx, index)
}

func (b *Editor) KeyPress(ctx *Context, ev KeyEvent) bool {
	switch ev.Key {
	case ebiten.KeyF1:
		b.Help.Hidden = !b.Help.Hidden
		if !b.Help.Hidden {
			ctx.Top(b.Help)
		}
		return true
	case ebiten.KeyNumpadAdd:
		b.Picker.SelectNext(0)
		return true
	case ebiten.KeyNumpadSubtract:
		b.Picker.SelectPrev(0)
		return true
	case ebiten.KeyC:
		if ev.Modifiers.Control {
			b.ChoosePanel.Hidden = !b.ChoosePanel.Hidden
			return true
		}
	case ebiten.KeyS:
		if ev.Modifiers.Control {
			return b.SaveMap(ctx)
		}
	case ebiten.KeyM:
		if ev.Modifiers.Control {
			return b.LoadMap(ctx)
		}

	case ebiten.KeyDelete:
		if ev.Modifiers.Control {
			return b.DeleteLayer(ctx)
		}

	case ebiten.KeyL:
		if ev.Modifiers.Control {
			return b.ChangeLayer(ctx)
		}
		if ev.Modifiers.Alt {
			return b.AddLayer()
		}
		if ev.Modifiers.Shift {
			return b.NextLayer()
		}
	case ebiten.KeyK:
		if ev.Modifiers.Shift {
			return b.PrevLayer()
		}
	}
	return false
}

func (b *Editor) KeyHold(ctx *Context, ev KeyEvent) bool {
	if ev.Duration > 60 && (ev.Duration%30) == 0 {
		return b.KeyPress(ctx, ev)
	}
	return false
}

func (b Editor) MouseCameraTile() image.Point {
	raw := b.MouseTile.Add(b.CameraTile)

	// Clamp to layer size
	raw.X = max(0, raw.X)
	raw.Y = max(0, raw.Y)
	if b.Map != nil {
		layer := b.Map.Layers[0]
		if b.Layer < len(b.Map.Layers) {
			layer = b.Map.Layers[b.Layer]
		}
		raw.X = min(raw.X, int(layer.TileSize.X))
		raw.Y = min(raw.Y, int(layer.TileSize.Y))
	}
	return raw
}

func (b *Editor) updateMouse(ctx *Context, ev MouseEvent) bool {
	// New mouse tile location.
	mx := (ev.X / TILE_W)
	my := (ev.Y / TILE_H)

	// don't clamp this, clamp MouseCameraTile
	b.MouseTile.X = mx
	b.MouseTile.Y = my

	return false
}

func (b *Editor) updateCursor(ctx *Context) bool {
	cx := 0
	cy := 0
	if b.Camera != nil {
		cx = int(b.Camera.X())
		cy = int(b.Camera.Y())
	}
	mct := b.MouseCameraTile()

	rmin := image.Pt(mct.X*TILE_W, mct.Y*TILE_H)
	rmin = rmin.Sub(image.Pt(cx, cy))
	rmax := rmin.Add(image.Pt(TILE_W, TILE_H))
	b.Cursor.Rectangle.Min = rmin
	b.Cursor.Rectangle.Max = rmax

	return true
}

func (b *Editor) UpdateCamera(ctx *Context, camera *image.Rectangle) bool {
	if camera == nil {
		return false
	}

	// Update camera but calculate delta.
	cx := int(camera.X()) / TILE_W
	cy := int(camera.Y()) / TILE_H

	b.CameraTile.X = cx
	b.CameraTile.Y = cy

	// don't clamp this, clamp MouseCameraTile

	b.updateCursor(ctx)
	return true
}

func (b *Editor) MouseMove(ctx *Context, ev MouseEvent) bool {
	if !b.AcceptMouse(ev) {
		return false
	}
	b.updateMouse(ctx, ev)
	b.updateCursor(ctx)
	return true
}

func (b *Editor) MouseRelease(ctx *Context, ev MouseEvent) bool {
	if !b.AcceptMouse(ev) {
		return false
	}
	return true
}

func (b *Editor) MouseHold(ctx *Context, ev MouseEvent) bool {
	if !b.AcceptMouse(ev) {
		return false
	}
	return b.drawTileWithMouse(ctx, ev)
}

func (b *Editor) Hide() {
	b.Hidden = true
}

func (b *Editor) Show() {
	b.Hidden = false
}

func AddEditor(c *Context, x, y, w, h int, mz *xmap.Zone, loader MapLoader) *Editor {
	editorGroup, editor := NewEditorGroup(x, y, w, h, tm, camera, loader)
	c.AddGroupWithNodes(editorGroup)
	return editorGroup, editor
}

*/
