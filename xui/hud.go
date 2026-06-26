package xui

import (
	"fmt"

	"github.com/xmasengine/xmas/xgal"
)

// HUDLayer is a heads-up display overlay showing player stats. It renders
// at the top/bottom of the screen and does not consume input.
// It consists of "bars" on the left, for exampe for HP display,
// and "stats" on the right.
type HUDLayer struct {
	Bounds xgal.Rectangle
	Style  Style
	Bars   []HUDBar
	Stats  []HUDStat
}

// HUDBAR is a bar on a HUD to diplay a ratio.
type HUDBar struct {
	Name  string    // Name of the bar.
	Fill  xgal.RGBA // Fill color of the bar. Outline is white.
	Value int       // Curent value
	Max   int       // Maximum value
}

// HUDStat is a statistic to display on the HUD
type HUDStat struct {
	Form   string // Format string
	Values []any  // Values to display using the format string.
}

func (h *HUDLayer) AddBar(name string, value int, max int, fill xgal.RGBA) *HUDBar {
	bar := HUDBar{Name: name, Value: value, Max: max, Fill: fill}
	h.Bars = append(h.Bars, bar)
	return &bar
}

func (h *HUDLayer) AddStat(form string, values ...any) *HUDStat {
	stat := HUDStat{Form: form, Values: values}
	h.Stats = append(h.Stats, stat)
	return &stat
}

const hudMinH = 32

// HUD creates a full-width HUDLayer bar at the top.
func HUD(screenW int) *HUDLayer {
	return &HUDLayer{
		Bounds: xgal.Rect(0, 0, screenW, hudMinH),
		Style:  DefaultStyle(),
	}
}

var _ Widget = &HUDLayer{}

func (h *HUDLayer) Poll() Reply { return Ignore }

const barMinW = 40

func (h *HUDLayer) Render(s *xgal.Surface) {
	// background bar
	xgal.Box(s, h.Bounds, xgal.Wash(0, 0, 0, 180))
	lineSize := xgal.Stride(h.Style.Face) + h.Style.Margin.Y

	// Draw bars on the left hand size
	for i, bar := range h.Bars {
		// HP bar

		barSpace := (h.Bounds.Dx() - barMinW) / 4

		fillWidth := (barSpace * bar.Value) / max(bar.Max, 1)
		outlWidth := barSpace

		fillBox := xgal.Rect(h.Bounds.Min.X+8, h.Bounds.Min.Y+4+i*lineSize,
			h.Bounds.Min.X+8+fillWidth, h.Bounds.Min.Y+h.Bounds.Min.Y+4+(i+1)*lineSize)
		outlBox := xgal.Rect(h.Bounds.Min.X+8, h.Bounds.Min.Y+4+i*lineSize,
			h.Bounds.Min.X+8+outlWidth, h.Bounds.Min.Y+h.Bounds.Min.Y+4+(i+1)*lineSize)

		xgal.Box(s, fillBox, bar.Fill)
		xgal.Outline(s, outlBox, 1, h.Style.Fore)

		// labels
		label := fmt.Sprintf("%s %d/%d", bar.Name, bar.Value, bar.Max)
		xgal.Ink(s, h.Style.Face, h.Style.Fore,
			h.Bounds.Min.X+8+outlWidth+4, h.Bounds.Min.Y+4+i*lineSize,
			label)

	}

	// Draw stats on the right side
	right := h.Bounds.Max.X - 8
	for i, stat := range h.Stats {
		xgal.Ink(s, h.Style.Face, h.Style.Fore,
			right-h.Style.Margin.X-100, h.Bounds.Min.Y+4+i*lineSize,
			fmt.Sprintf(stat.Form, stat.Values...))
	}

}

func (h *HUDLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	h.Bounds.Min.X = bounds.Min.X
	h.Bounds.Max.X = bounds.Max.X
	lineSize := xgal.Stride(h.Style.Face) + h.Style.Margin.Y
	h.Bounds.Max.Y = h.Bounds.Min.Y + max(hudMinH, max(len(h.Bars), len(h.Stats))*lineSize) + h.Style.Margin.Y
	return h.Bounds
}

func (h *HUDLayer) MoveBy(delta xgal.Point) {
	h.Bounds = h.Bounds.Add(delta)
}
