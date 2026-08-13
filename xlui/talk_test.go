package xlui

import (
	"strings"
	"testing"

	"github.com/xmasengine/xmas/xgal"
)

func TestRevealedText(t *testing.T) {
	cases := []struct {
		text string
		n    int
		want string
	}{
		{"Hello\nworld", 0, ""},
		{"Hello\nworld", 5, "Hello"},
		{"Hello\nworld", 6, "Hello\nw"},
		{"Hello\nworld", 11, "Hello\nworld"},
		{"Hello\nworld", 99, "Hello\nworld"},
		{"Héllo wörld\nLast", 6, "Héllo "},
		{"Héllo wörld\nLast", 12, "Héllo wörld\nL"},
		{"\n\nThird", 1, "\n\nT"},
		{"\n\nThird", 2, "\n\nTh"},
		{"\n\nThird", 3, "\n\nThi"},
		{"Single", 7, "Single"},
	}
	for _, c := range cases {
		got := revealedText(textToRuneLines(c.text), c.n)
		if got != c.want {
			t.Errorf("revealedText(%q, %d) = %q, want %q", c.text, c.n, got, c.want)
		}
	}
}

func TestRevealCursor(t *testing.T) {
	cases := []struct {
		text string
		n    int
		want xgal.Point
	}{
		{"Hello\nworld", 0, xgal.Pt(0, 0)},
		{"Hello\nworld", 5, xgal.Pt(5, 0)},
		{"Hello\nworld", 6, xgal.Pt(1, 1)},
		{"Hello\nworld", 11, xgal.Pt(5, 1)},
		{"Hello\nworld", 99, xgal.Pt(5, 1)},
		{"Héllo wörld\nLast", 12, xgal.Pt(1, 1)},
		{"\n\nThird", 1, xgal.Pt(1, 2)},
		{"\n\nThird", 2, xgal.Pt(2, 2)},
		{"\n\nThird", 3, xgal.Pt(3, 2)},
		{"Single", 7, xgal.Pt(6, 0)},
	}
	for _, c := range cases {
		got := revealCursor(textToRuneLines(c.text), c.n)
		if got != c.want {
			t.Errorf("revealCursor(%q, %d) = %v, want %v", c.text, c.n, got, c.want)
		}
	}
}

// TestRevealTrace ensures the cursor always sits exactly one character past
// the revealed text, for every reveal step.
func TestRevealTrace(t *testing.T) {
	text := "Héllo wörld\nThis is me\nLast line"
	lines := textToRuneLines(text)
	total := 0
	for _, line := range lines {
		total += len(line)
	}
	for n := 0; n <= total; n++ {
		rendered := revealedText(lines, n)
		lastLine := rendered[strings.LastIndex(rendered, "\n")+1:]
		want := xgal.Pt(len([]rune(lastLine)), strings.Count(rendered, "\n"))
		got := revealCursor(lines, n)
		if got != want {
			t.Fatalf("reveal step %d: cursor %v, want %v (text %q)", n, got, want, rendered)
		}
	}
}
