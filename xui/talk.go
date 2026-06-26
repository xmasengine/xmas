package xui

import (
	"time"

	"github.com/xmasengine/xmas/xgal"
)

// TalkLayer is a dialog box for NPC conversations. It displays a
// speaker portrait (optional) and text revealed with a typewriter
// effect.  Advance with confirm key; Finish is returned when the
// last message is dismissed.
type TalkLayer struct {
	Bounds   xgal.Rectangle
	Style    Style
	Portrait *xgal.Surface // speaker face, drawn on the left

	messages  []string
	msgIdx    int
	reveal    int // characters revealed so far in current message
	speed     time.Duration
	lastChar  time.Time
	done      bool
	advanceFn func() bool // returns true when the player wants to advance
}

const talkSpeed = 30 * time.Millisecond

// Talk creates a talk dialog.  advance is called each frame
// and should return true when the player presses confirm (A/Enter/etc).
func Talk(bounds xgal.Rectangle, portrait *xgal.Surface, messages []string, advance func() bool) *TalkLayer {
	return &TalkLayer{
		Bounds:    bounds,
		Style:     DefaultStyle(),
		Portrait:  portrait,
		messages:  messages,
		speed:     talkSpeed,
		advanceFn: advance,
	}
}

var _ Widget = &TalkLayer{}

func (t *TalkLayer) Poll() Reply {
	if t.done {
		return Finish
	}

	msg := t.messages[t.msgIdx]

	// typewriter reveal
	if t.reveal < len(msg) {
		if time.Since(t.lastChar) >= t.speed {
			t.reveal++
			t.lastChar = time.Now()
		}
		return Accept
	}

	// message fully revealed — wait for advance
	if t.advanceFn != nil && t.advanceFn() {
		t.msgIdx++
		t.reveal = 0
		t.lastChar = time.Now()
		if t.msgIdx >= len(t.messages) {
			t.done = true
			return Finish
		}
	}
	return Accept
}

func (t *TalkLayer) Render(s *xgal.Surface) {
	t.Style.DrawBox(s, t.Bounds)

	lx := t.Bounds.Min.X + t.Style.Margin.X
	ty := t.Bounds.Min.Y + t.Style.Margin.Y

	// portrait
	if t.Portrait != nil {
		pb := t.Portrait.Bounds()
		xgal.Blit(s, t.Portrait,
			xgal.Rect(lx, ty, lx+pb.Dx(), ty+pb.Dy()),
			pb)
		lx += pb.Dx() + t.Style.Margin.X
	}

	// text
	if t.msgIdx < len(t.messages) {
		msg := t.messages[t.msgIdx]
		if t.reveal < len(msg) {
			msg = msg[:t.reveal]
		}
		xgal.Ink(s, t.Style.Face, t.Style.Fore, lx, ty, msg)
	}

	// advance indicator when fully revealed
	if t.msgIdx < len(t.messages) && t.reveal >= len(t.messages[t.msgIdx]) {
		// draw a small indicator at the bottom-right as "press to continue"
		indicator := xgal.Rect(t.Bounds.Max.X-16, t.Bounds.Max.Y-12,
			t.Bounds.Max.X-8, t.Bounds.Max.Y-4)
		xgal.Disk(s, xgal.Pt(indicator.Min.X, indicator.Min.Y), 2, t.Style.Fore)
	}
}

func (t *TalkLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	return t.Bounds
}

func (t *TalkLayer) MoveBy(delta xgal.Point) {
	t.Bounds = t.Bounds.Add(delta)
}

// Done returns true after the last message was dismissed.
func (t *TalkLayer) Done() bool { return t.done }
