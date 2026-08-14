package xlui

import (
	"os"
	"testing"
	"testing/fstest"
)

// TestSetFrameFS ensures the frame can be loaded from a configurable FS and
// that a missing frame reports an error.
func TestSetFrameFS(t *testing.T) {
	if err := SetFrameFS(os.DirFS("../pack/image/ui")); err != nil {
		t.Fatalf("SetFrameFS: %v", err)
	}
	if DefaultVec == nil || DefaultFrame == nil {
		t.Fatal("SetFrameFS did not load the frame")
	}
	if err := SetFrameFS(fstest.MapFS{}); err == nil {
		t.Fatal("SetFrameFS succeeded without frame.xvec")
	}
}
