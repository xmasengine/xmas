package xui

import (
	"image"
	"image/color"
)

import (
	"github.com/xmasengine/xmas/xres"
)

type Icon = xres.Icon

const IconSize = 8
const CheckboxSize = 8

var CheckColor = color.RGBA{255, 0, 0, 255}

type Checkbox struct {
	Button
	Checked bool
	Check   func(*Checkbox)
	Value   int
	Icon    *Icon
}

func (c *Checkbox) Init(bounds Rectangle, text string, ch func(*Checkbox)) *Checkbox {
	adapt := func(b *Button) {
		if ch != nil {
			if c != nil {
				c.Checked = !c.Checked
			}
			ch(c)
		}
	}
	c.Button.Init(bounds, text, adapt)
	c.Style.Fill = color.RGBA{0, 0, 0, 0}
	c.Class = NewCheckboxClass(c)
	return c
}

func NewCheckbox(Rectangle Rectangle, text string, ch func(*Checkbox)) *Checkbox {
	cb := &Checkbox{}
	return cb.Init(Rectangle, text, ch)
}

type CheckboxClass struct {
	*Checkbox
	*ButtonClass
}

func NewCheckboxClass(c *Checkbox) *CheckboxClass {
	res := &CheckboxClass{Checkbox: c}
	res.ButtonClass = NewButtonClass(&c.Button)
	return res
}

func (bc CheckboxClass) Render(r *Root, screen *Surface) {
	b := bc.Checkbox
	box := b.Bounds

	tbox := box.Add(image.Pt(CheckboxSize+b.Style.Margin.X*2, CheckboxSize+b.Style.Margin.Y*2))
	at := tbox.Min
	style := b.Style
	if b.State.Hover {
		style = HoverStyle()
	}

	// Draw the checkbox.
	// The style will have been changed depending on how it should be drawn.
	// cb := image.Rect(box.Min.X.+b.Style.Margin, box.Min.Y+b.Style.Margin, box.Min.X+CheckboxSize, box.Min.Y+CheckboxSize)
	style.DrawBox(screen, tbox)

	ibox := image.Rect(tbox.Min.X, tbox.Min.Y, tbox.Min.X+IconSize, tbox.Min.Y+IconSize)
	at = at.Add(image.Pt(IconSize, 0))
	if b.Icon != nil {
		b.Icon.Draw(screen, ibox)
	} else {
		cstyle := CheckStyle()
		if b.Checked {
			cstyle.DrawBox(screen, ibox)
		} else {
			cstyle.DrawRect(screen, ibox)
		}
	}

	style.DrawText(screen, at, b.Text)

}

func (w *Widget) AddCheckbox(bounds Rectangle, text string, ch func(*Checkbox)) *Checkbox {
	b := NewCheckbox(bounds, text, ch)
	w.Widgets = append(w.Widgets, &b.Widget)
	return b
}
