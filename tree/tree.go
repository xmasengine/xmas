// tree is a tree based retained mode UI
package tree

import (
	"image"
	"image/color"
)

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Rectangle is used for sizes and positions.
type Rectangle = image.Rectangle

// Point is used for position and offsets.
type Point = image.Point

// Color is a color.
type Color = color.Color

// Image is an image.Image
type Image = image.Image

// Surface is an ebiten.Image
type Surface = ebiten.Image

// State is the state of an Element.
type State struct {
	Focus bool
	Hover bool
	Pause bool
	Hide  bool
	Clip  bool
	Value any
}

// Bounds are the bounds of an element
type Bounds struct {
	Rectangle
}

// NewBounds return a new bounds for an element
func NewBounds(x, y, w, h int) Bounds {
	b := Bounds{}
	b.Rectangle = image.Rect(x, y, x+w, y+h)
	return b
}

// Element is a basic UI element, component or widget.
type Element interface {
	// Draw is called when the element needs to be drawn
	Draw(screen *ebiten.Image)

	// Place places the widget at the given bounds.
	// It should react by also updating any contained elements if appropriate.
	// It is allowed for the widget to shrink and become smaller than the
	// requested bounds.
	// It should return the actual size of the widget
	Place(bounds Bounds) (size Point)

	// Bounds are the actual absolute visual bounds of the element,
	// in screen coordinates, as should be used for layout,
	// ignoring any popups or overflows.
	Bounds() Bounds

	// State returns the state of the element.
	State() State

	// Modify sets the state of the element.
	Modify(set State)
}

// Applier sets an aspect of an element. Is used when initializing
// or updating Elements.
type Applier interface {
	Apply(element Element)
}

// Layouter applies a layout on an element or container.
type Layouter interface {
	Layout(availableWidth, availableHeight int, element Element) (totalWidth, totalHeight int)
}

// List is a list of Elements
type List []Element

// NewList returns a lew List of elements
func NewList(l ...Element) List {
	return List(l)
}

// Container is an element that also contains other elements.
type Container interface {
	Element
	// Contain returns list of the elements in draw order.
	// Events are processed in the opposite order of draw order.
	Contain() List
}

type Result bool

func FindTop(at Point, c Container) Element {
	elts := c.Contain()
	for i := len(elts) - 1; i >= 0; i-- {
		elt := elts[i]
		if sub, ok := elt.(Container); ok {
			found := FindTop(at, sub)
			if found != nil {
				return found
			}
		}

		bounds := elt.Bounds().Rectangle
		println("bounds", bounds.Min.X, bounds.Min.Y)
		if at.In(bounds) {
			return elt
		}
	}
	return nil
}

// Control is the abstract implementation part of the UI Node.
type Control interface {
	// Draw is called when the tree needs to be drawn.
	Draw(screen *ebiten.Image)
}

// Node is a UI tree node.
// It is a concrete struct with an abstract implementation.
type Node struct {
	Control
	Parent *Node
	Bounds Rectangle
	State  State
	Nodes  []*Node
}
