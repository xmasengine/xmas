// Package xlog is the xmas engine GUI debug log on top of slog.
// It is not part of xui to avoid cyclical dependencies and to allow
// standalone use.
// To keep everything relatively simple, there can only be a single active
// xlog, and we use package level variables for the settings.
package xlog

import (
	"image"
	"image/color"
	"log/slog"
	"strings"
	"sync"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var LineSkip = 8
var ShowLines = 24
var Fill = color.RGBA{64, 64, 64, 128}
var Stroke = 5
var Border = color.RGBA{128, 128, 32, 128}

type Log struct {
	Lines []string
	sync.Mutex
	From    int
	Hide    bool
	Size    image.Point
	Pressed []ebiten.Key
	Level   slog.LevelVar
}

func (l *Log) Write(buf []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	s := string(buf)
	lines := strings.Split(s, "\n")
	l.Lines = append(l.Lines, lines...)
	if len(l.Lines) > ShowLines {
		l.From += len(lines)
	}

	return len(buf), nil
}

func (l *Log) NewHandler(opts *slog.HandlerOptions) *slog.TextHandler {
	return slog.NewTextHandler(l, opts)
}

func (l *Log) Logger() *slog.Logger {
	l.Level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{
		Level: &l.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			if a.Key == "level" {
				a.Key = "l"
				if a.Value.String() == "INFO" {
					a.Value = slog.StringValue("I")
				}
				if a.Value.String() == "DEBUG" {
					a.Value = slog.StringValue("D")
				}
				if a.Value.String() == "ERROR" {
					a.Value = slog.StringValue("E")
				}
				if a.Value.String() == "WARN" {
					a.Value = slog.StringValue("W")
				}
			} else if a.Key == "msg" {
				a.Key = "m"
			}
			return a
		},
	}
	return slog.New(l.NewHandler(opts))
}

func (l *Log) Update() error {
	l.Pressed = l.Pressed[:0]
	l.Pressed = inpututil.AppendPressedKeys(l.Pressed)
	for _, k := range l.Pressed {
		if k == ebiten.KeyF10 {
			l.Hide = false
		}
		if k == ebiten.KeyF11 {
			l.Hide = true
		}
		if l.Hide {
			continue
		}
		switch k {
		case ebiten.KeyUp:
			l.From--
			if l.From < 0 {
				l.From = 0
			}
		case ebiten.KeyDown:
			l.From++
			if l.From > len(l.Lines)-1 {
				l.From = len(l.Lines) - 2
			}
		case ebiten.KeyPageUp:
		case ebiten.KeyPageDown:
		case ebiten.KeyE:
			l.Level.Set(slog.LevelError)
		case ebiten.KeyD:
			l.Level.Set(slog.LevelDebug)
		case ebiten.KeyI:
			l.Level.Set(slog.LevelInfo)
		case ebiten.KeyW:
			l.Level.Set(slog.LevelWarn)
		case ebiten.KeyC:
			l.Lines = []string{"Cleared"}
			l.From = 0
		default:
		}
	}

	return nil
}

func (l *Log) DrawBox(surface *ebiten.Image, r image.Rectangle) {
	vector.DrawFilledRect(
		surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()), Fill, false,
	)

	b := r.Inset(Stroke)

	if Stroke > 0 {
		vector.StrokeRect(
			surface, float32(b.Min.X), float32(b.Min.Y),
			float32(b.Dx()), float32(b.Dy()),
			float32(Stroke), Border, false,
		)
	}
}

func (l *Log) Draw(screen *ebiten.Image) {
	if l.Hide {
		return
	}
	x := 0
	y := 0
	l.DrawBox(screen, image.Rect(x, y, l.Size.X, l.Size.Y))
	for i := l.From; i <= l.From+ShowLines && i >= 0 && i < len(l.Lines); i++ {
		ebitenutil.DebugPrintAt(screen, l.Lines[i], x, y)
		y += LineSkip
	}
}

func (l *Log) Layout(width, height int) (screenWidth, screenHeight int) {
	l.Size = image.Pt(width, height)
	return width, height

}
