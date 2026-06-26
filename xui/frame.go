package xui

import "github.com/xmasengine/xmas/xgal"

// FrameLayer is a widget that displays an image. When Image is nil,
// it renders as a plain box via [Layer].
type FrameLayer struct {
	Layer
	Image *xgal.Surface
}

// Frame returns a new [FrameLayer]. If img is non-nil the bounds are
// auto-sized to the image dimensions.
func Frame(bounds xgal.Rectangle, img *xgal.Surface) *FrameLayer {
	f := &FrameLayer{Image: img}
	if img != nil {
		r := img.Bounds()
		bounds.Max.X = bounds.Min.X + r.Dx()
		bounds.Max.Y = bounds.Min.Y + r.Dy()
	}
	f.Layer = MakeLayer(bounds)
	return f
}

func (f *FrameLayer) Render(s *xgal.Surface) {
	f.Layer.Render(s)
	if f.Image != nil {
		xgal.Blit(s, f.Image, f.Bounds, f.Image.Bounds())
	}
}

func (f *FrameLayer) Place(bounds xgal.Rectangle) xgal.Rectangle {
	if f.Image != nil {
		iw := f.Image.Bounds().Dx()
		ih := f.Image.Bounds().Dy()
		f.Bounds = xgal.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+iw, bounds.Min.Y+ih)
		return f.Bounds
	}
	return f.Layer.Place(bounds)
}

func (m *Layer) AddFrame(bounds xgal.Rectangle, img *xgal.Surface) *FrameLayer {
	f := Frame(bounds, img)
	m.Add(f)
	return f
}
