package xzed

import "image"

import (
	"github.com/xmasengine/xmas/xmap"
	"github.com/xmasengine/xmas/xui"
)

// Picker is a Picker for the current tile.
type Picker struct {
	xui.Box
	// Tiles are frames that display the current tiles
	Tiles []*xui.Frame
	// Indexes to the current tiles
	Indexes []int
	// Layer the picker is for
	Layer *xmap.Layer
	// OnPick callback.
	OnPick func(*Picker)
	Index  int // Index result set when OnPick Is called
}

type PickerClass struct {
	*xui.BoxClass
	*Picker
}

func NewPickerClass(b *Picker) *PickerClass {
	res := &PickerClass{Picker: b}
	res.BoxClass = xui.NewBoxClass(&b.Box)
	return res
}

func NewPicker(at Point, layer *xmap.Layer, tiles int, cb func(*Picker)) *Picker {
	t := &Picker{}
	t.Init(at, layer, tiles, cb)
	return t
}

func (t *Picker) Init(at Point, layer *xmap.Layer, tiles int, cb func(*Picker)) *Picker {
	// Set up the tile box.
	w := layer.Tw * tiles * 2
	h := layer.Th*2 + 10
	size := image.Pt(w, h)
	bounds := image.Rectangle{Min: at, Max: at.Add(size)}
	t.Box.Init(bounds)
	t.Tiles = make([]*xui.Frame, tiles)
	t.Indexes = make([]int, len(t.Tiles))
	t.OnPick = cb
	t.AddTitleBar(10, "Tile")
	t.ChangeLayer(layer)
	return t
}

func (t *PickerClass) Render(r *xui.Root, Surface *xui.Surface) {
	t.BoxClass.Render(r, Surface)
}

func (t *Picker) SetSelected(at image.Point, which int) {
	if which < 0 || which >= len(t.Indexes) {
		return
	}

	/*
	   	idx := xmap.Index{X: uint8(at.X), Y: uint8(at.Y)}
	   	t.Indexes[which] = idx
	   	if t.Atlas != nil {
	   		t.Tiles[which].Image = t.Atlas.Get(int(t.Indexes[which].X), int(t.Indexes[which].Y))
	   	}

	   *
	*/
}

func (t *Picker) SelectNext(which int) {
	if which < 0 || which > len(t.Indexes) {
		return
	}
	/*
		if t.Atlas != nil {
			t.Indexes[which].NextIn(*t.Atlas)
			idx := t.Indexes[which]
			t.Tiles[which].Image = t.Atlas.Get(int(idx.X), int(idx.Y))
		}
	*/
}

func (t *Picker) SelectPrev(which int) {
	if which < 0 || which > len(t.Indexes) {
		return
	}
	/*
		if t.Atlas != nil {
			t.Indexes[which].PrevIn(*t.Atlas)
			idx := t.Indexes[which]
			t.Tiles[which].Image = t.Atlas.Get(int(idx.X), int(idx.Y))
		}
	*/
}

func (t *PickerClass) OnMouseWheel(ev MouseEvent) bool {
	for idx := 0; idx < len(t.Tiles); idx++ {
		if ev.Wheel.Y > 0 || ev.Wheel.X > 0 {
			t.SelectNext(idx)
			return true
		} else if ev.Wheel.Y < 0 || ev.Wheel.X < 0 {
			t.SelectPrev(idx)
			return true
		}
	}
	return false
}

func (t *PickerClass) OnMouseRelease(ev MouseEvent) bool {
	for idx := 0; idx < len(t.Tiles); idx++ {
		if t.OnPick != nil {
			t.Index = t.Indexes[idx]
			t.OnPick(t.Picker)
			return true
		}
	}
	return false
}

func (t *Picker) ChangeLayer(layer *xmap.Layer) {
	t.Layer = layer
	if t.Layer == nil {
		return
	}

	for i := 0; i < len(t.Tiles); i++ {
		w := layer.Tw
		h := layer.Th
		size := image.Pt(w, h)
		offset := image.Pt(i*w, 12)
		at := t.Bounds.Min
		bounds := image.Rectangle{Min: at.Add(offset), Max: at.Add(offset).Add(size)}
		t.Tiles[i] = t.AddFrame(bounds, nil)
		t.Indexes[i] = int(i)
		t.Tiles[i].Image = t.Layer.Atlas.SubImage(bounds).(*Surface)
	}
}
