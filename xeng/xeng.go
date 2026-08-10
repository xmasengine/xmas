package xeng

import (
	"fmt"
	"image"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"
)

import (
	"github.com/xmasengine/xmas/xdat"
	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xlog"
	// "github.com/xmasengine/xmas/xres"
	"github.com/xmasengine/xmas/xlui"
	"github.com/xmasengine/xmas/xzed"
)

// We want a resolution similar to the SMS: 256 × 192
// However most screens are wide view now so we will do 320 x 192 for now.

const ViewWidth = 320
const ViewHeight = 192

// const ViewHeight = 240 * 2

type Engine struct {
	Log         xlog.Log
	Msg         string
	Pressed     []xgal.KeyCode
	Script      strings.Reader
	DebugRow    int
	ScreenSize  image.Point
	Zone        *xdat.Zone
	FS          fs.FS
	Debug       bool
	Camera      xgal.Rectangle
	EditorLayer *xlui.Layer
	Editor      *xzed.Editor
	Windowed    bool
}

func New(sw, sh int) *Engine {
	engine := &Engine{ScreenSize: image.Point{X: sw, Y: sh}, Msg: "!"}
	engine.Camera = image.Rect(0, 0, ViewWidth, ViewHeight)
	engine.Pressed = make([]xgal.KeyCode, 16)
	engine.Log.Hide = true
	wd, _ := os.Getwd()
	engine.FS = os.DirFS(wd)
	engine.loadFirstZone()
	return engine
}

func dprintln(msg string, vars ...any) {
	slog.Info(msg, "vars", vars)
}

func (engine *Engine) loadFirstZone() {
	_, err := engine.LoadZone("map_0001.xml")
	if err != nil {
		slog.Error("loading zone", "err", err)
	}
}

func (g *Engine) Update() error {
	g.Log.Update()

	res := xlui.Poll()
	if res == xlui.Finish {
		return nil
	}

	g.Pressed = g.Pressed[:0]
	g.Pressed = xgal.Keys(g.Pressed)
	var delta image.Point
	var mdelta image.Point
	// act := xdat.Stand
	// var dir xdat.Direction
	if g.Zone != nil {
		// dir = g.Zone.Player.Direction
	}
	for _, k := range g.Pressed {
		switch k {
		case xgal.KeyArrowUp:
			delta.Y = -1
			// dir = xdat.North
			// act = xdat.Walk
		case xgal.KeyArrowDown:
			delta.Y = 1
			// dir = xdat.South
			// act = xdat.Walk
		case xgal.KeyArrowLeft:
			delta.X = -1
			// dir = xdat.West
			// act = xdat.Walk
		case xgal.KeyArrowRight:
			delta.X = 1
			// dir = xdat.East
			// act = xdat.Walk
		case xgal.KeyPageUp:
			mdelta.Y = -1
		case xgal.KeyPageDown:
			mdelta.Y = 1
		case xgal.KeyHome:
			mdelta.X = -1
		case xgal.KeyEnd:
			mdelta.X = 1

		default:
		}
	}

	if g.Zone != nil {
		g.Camera = g.Camera.Add(mdelta)
		// g.Zone.Player.At = g.Zone.Player.At.Add(delta)
		// pose := g.Zone.Player.BestPose(dir, act)
		// g.Zone.Player.Pose = pose
		// g.Zone.Player.Update()
	}

	switch {
	case g.Editor != nil && g.Editor.Done:
		return xgal.Quit

	case xgal.Tap(xgal.KeyEscape):
		return xgal.Quit

	case xgal.Tap(xgal.KeyF10):
		if g.Zone != nil {
			if g.Editor == nil && g.EditorLayer == nil {
				g.EditorLayer = xzed.NewEditorLayer(g, g.Zone, "map_0001.xml", &g.Camera, 1)
				g.Editor = g.EditorLayer.Data.(*xzed.Editor)
				xlui.Append(g.EditorLayer)
			} else {
				xlui.CloseLayer(g.EditorLayer)
				g.EditorLayer = nil
				g.Editor = nil
			}
		}
	case xgal.Tap(xgal.KeyF):
		g.Debug = !g.Debug
		xlui.Ask(50, 50, 250, 100, "Debug", "debug", func(string) bool { return true })

	case xgal.Tap(xgal.KeyPrintScreen):
		g.Windowed = !g.Windowed
		xgal.Expand(!g.Windowed)
	default:
	}

	return nil
}

const tileDebug = false

func (g *Engine) Draw(screen *xgal.Surface) {
	if g.Zone != nil {
		g.RenderZone(screen, g.Camera)
		if g.Debug {
			// pose := g.Zone.Player.Pose
			// xgal.Debug(screen, fmt.Sprintf("pose: %d %d %d %d %d",
			//	pose.Direction, pose.Action, pose.Phase, pose.Frames, pose.Tick), 0, 0)
		}
	}

	xlui.Render(screen)

	if g.Debug {
		xgal.Debug(screen, fmt.Sprintf("\n%f\n", xgal.FPS()), 0, 0)
	}
	g.Log.Draw(screen)
}

func (g *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Log.Layout(ViewWidth, ViewHeight)
	return ViewWidth, ViewHeight
}

const ZoneDir = "pack/map"

func (g *Engine) LoadZone(name string) (*xdat.Zone, error) {
	z, err := xdat.LoadZone(g.FS, path.Join(ZoneDir, name))
	if err != nil {
		return nil, err
	}

	g.Zone = z
	return z, nil
}

func (g *Engine) SetLayerSource(layer *xdat.Layer, name string) error {
	return layer.SetSource(g.FS, name)
}

func (g *Engine) GetLayer(depth int) *xdat.Layer {
	if g.Zone == nil {
		return nil
	}
	if depth < 0 || depth >= len(g.Zone.Layers) {
		return nil
	}
	return g.Zone.Layers[depth]
}
