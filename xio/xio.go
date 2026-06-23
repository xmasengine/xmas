// Package xio is a wrapper around the low level ebitengine game library
// and its supporting libraries, as well as the image and image/color library.
// While ebitengine works well, it is a bit of a hassle to import everything
// separately.
// Also this wrapper could be useful in case we have to use a different
// low level game library.
package xio

import (
	"image"
	"image/color"
)

import (
	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RGBA = color.RGBA

type Rectangle = image.Rectangle
type Point = image.Point

var BitmapFontFace = bitmapfont.Face

// TextInputField is an input field for IME text entry.
// It embeds the original textinput.Field
type TextInputField struct {
	textinput.Field
	X int
	Y int
}

// KeyMods are the current key modifers.
type KeyMods struct {
	Alt   bool
	Class bool
	Shift bool
	Meta  bool
}

// Surface is an ebiten.Image
// We use surface to still be able to use image.Image together easily.
type Surface = ebiten.Image

// Face is a font face
type Face = text.Face

// DrawOptions are options for drawing a Surface.
type DrawOptions = ebiten.DrawImageOptions

// A Key represents a physical keyboard key code as if it was typed on a QWERTY
// keyboard.
type Key int

const (
	KeyA              Key = Key(ebiten.KeyA)
	KeyB              Key = Key(ebiten.KeyB)
	KeyC              Key = Key(ebiten.KeyC)
	KeyD              Key = Key(ebiten.KeyD)
	KeyE              Key = Key(ebiten.KeyE)
	KeyF              Key = Key(ebiten.KeyF)
	KeyG              Key = Key(ebiten.KeyG)
	KeyH              Key = Key(ebiten.KeyH)
	KeyI              Key = Key(ebiten.KeyI)
	KeyJ              Key = Key(ebiten.KeyJ)
	KeyK              Key = Key(ebiten.KeyK)
	KeyL              Key = Key(ebiten.KeyL)
	KeyM              Key = Key(ebiten.KeyM)
	KeyN              Key = Key(ebiten.KeyN)
	KeyO              Key = Key(ebiten.KeyO)
	KeyP              Key = Key(ebiten.KeyP)
	KeyQ              Key = Key(ebiten.KeyQ)
	KeyR              Key = Key(ebiten.KeyR)
	KeyS              Key = Key(ebiten.KeyS)
	KeyT              Key = Key(ebiten.KeyT)
	KeyU              Key = Key(ebiten.KeyU)
	KeyV              Key = Key(ebiten.KeyV)
	KeyW              Key = Key(ebiten.KeyW)
	KeyX              Key = Key(ebiten.KeyX)
	KeyY              Key = Key(ebiten.KeyY)
	KeyZ              Key = Key(ebiten.KeyZ)
	KeyAltLeft        Key = Key(ebiten.KeyAltLeft)
	KeyAltRight       Key = Key(ebiten.KeyAltRight)
	KeyArrowDown      Key = Key(ebiten.KeyArrowDown)
	KeyArrowLeft      Key = Key(ebiten.KeyArrowLeft)
	KeyArrowRight     Key = Key(ebiten.KeyArrowRight)
	KeyArrowUp        Key = Key(ebiten.KeyArrowUp)
	KeyBackquote      Key = Key(ebiten.KeyBackquote)
	KeyBackslash      Key = Key(ebiten.KeyBackslash)
	KeyBackspace      Key = Key(ebiten.KeyBackspace)
	KeyBracketLeft    Key = Key(ebiten.KeyBracketLeft)
	KeyBracketRight   Key = Key(ebiten.KeyBracketRight)
	KeyCapsLock       Key = Key(ebiten.KeyCapsLock)
	KeyComma          Key = Key(ebiten.KeyComma)
	KeyContextMenu    Key = Key(ebiten.KeyContextMenu)
	KeyControlLeft    Key = Key(ebiten.KeyControlLeft)
	KeyControlRight   Key = Key(ebiten.KeyControlRight)
	KeyDelete         Key = Key(ebiten.KeyDelete)
	KeyDigit0         Key = Key(ebiten.KeyDigit0)
	KeyDigit1         Key = Key(ebiten.KeyDigit1)
	KeyDigit2         Key = Key(ebiten.KeyDigit2)
	KeyDigit3         Key = Key(ebiten.KeyDigit3)
	KeyDigit4         Key = Key(ebiten.KeyDigit4)
	KeyDigit5         Key = Key(ebiten.KeyDigit5)
	KeyDigit6         Key = Key(ebiten.KeyDigit6)
	KeyDigit7         Key = Key(ebiten.KeyDigit7)
	KeyDigit8         Key = Key(ebiten.KeyDigit8)
	KeyDigit9         Key = Key(ebiten.KeyDigit9)
	KeyEnd            Key = Key(ebiten.KeyEnd)
	KeyEnter          Key = Key(ebiten.KeyEnter)
	KeyEqual          Key = Key(ebiten.KeyEqual)
	KeyEscape         Key = Key(ebiten.KeyEscape)
	KeyF1             Key = Key(ebiten.KeyF1)
	KeyF2             Key = Key(ebiten.KeyF2)
	KeyF3             Key = Key(ebiten.KeyF3)
	KeyF4             Key = Key(ebiten.KeyF4)
	KeyF5             Key = Key(ebiten.KeyF5)
	KeyF6             Key = Key(ebiten.KeyF6)
	KeyF7             Key = Key(ebiten.KeyF7)
	KeyF8             Key = Key(ebiten.KeyF8)
	KeyF9             Key = Key(ebiten.KeyF9)
	KeyF10            Key = Key(ebiten.KeyF10)
	KeyF11            Key = Key(ebiten.KeyF11)
	KeyF12            Key = Key(ebiten.KeyF12)
	KeyF13            Key = Key(ebiten.KeyF13)
	KeyF14            Key = Key(ebiten.KeyF14)
	KeyF15            Key = Key(ebiten.KeyF15)
	KeyF16            Key = Key(ebiten.KeyF16)
	KeyF17            Key = Key(ebiten.KeyF17)
	KeyF18            Key = Key(ebiten.KeyF18)
	KeyF19            Key = Key(ebiten.KeyF19)
	KeyF20            Key = Key(ebiten.KeyF20)
	KeyF21            Key = Key(ebiten.KeyF21)
	KeyF22            Key = Key(ebiten.KeyF22)
	KeyF23            Key = Key(ebiten.KeyF23)
	KeyF24            Key = Key(ebiten.KeyF24)
	KeyHome           Key = Key(ebiten.KeyHome)
	KeyInsert         Key = Key(ebiten.KeyInsert)
	KeyIntlBackslash  Key = Key(ebiten.KeyIntlBackslash)
	KeyMetaLeft       Key = Key(ebiten.KeyMetaLeft)
	KeyMetaRight      Key = Key(ebiten.KeyMetaRight)
	KeyMinus          Key = Key(ebiten.KeyMinus)
	KeyNumLock        Key = Key(ebiten.KeyNumLock)
	KeyNumpad0        Key = Key(ebiten.KeyNumpad0)
	KeyNumpad1        Key = Key(ebiten.KeyNumpad1)
	KeyNumpad2        Key = Key(ebiten.KeyNumpad2)
	KeyNumpad3        Key = Key(ebiten.KeyNumpad3)
	KeyNumpad4        Key = Key(ebiten.KeyNumpad4)
	KeyNumpad5        Key = Key(ebiten.KeyNumpad5)
	KeyNumpad6        Key = Key(ebiten.KeyNumpad6)
	KeyNumpad7        Key = Key(ebiten.KeyNumpad7)
	KeyNumpad8        Key = Key(ebiten.KeyNumpad8)
	KeyNumpad9        Key = Key(ebiten.KeyNumpad9)
	KeyNumpadAdd      Key = Key(ebiten.KeyNumpadAdd)
	KeyNumpadDecimal  Key = Key(ebiten.KeyNumpadDecimal)
	KeyNumpadDivide   Key = Key(ebiten.KeyNumpadDivide)
	KeyNumpadEnter    Key = Key(ebiten.KeyNumpadEnter)
	KeyNumpadEqual    Key = Key(ebiten.KeyNumpadEqual)
	KeyNumpadMultiply Key = Key(ebiten.KeyNumpadMultiply)
	KeyNumpadSubtract Key = Key(ebiten.KeyNumpadSubtract)
	KeyPageDown       Key = Key(ebiten.KeyPageDown)
	KeyPageUp         Key = Key(ebiten.KeyPageUp)
	KeyPause          Key = Key(ebiten.KeyPause)
	KeyPeriod         Key = Key(ebiten.KeyPeriod)
	KeyPrintScreen    Key = Key(ebiten.KeyPrintScreen)
	KeyQuote          Key = Key(ebiten.KeyQuote)
	KeyScrollLock     Key = Key(ebiten.KeyScrollLock)
	KeySemicolon      Key = Key(ebiten.KeySemicolon)
	KeyShiftLeft      Key = Key(ebiten.KeyShiftLeft)
	KeyShiftRight     Key = Key(ebiten.KeyShiftRight)
	KeySlash          Key = Key(ebiten.KeySlash)
	KeySpace          Key = Key(ebiten.KeySpace)
	KeyTab            Key = Key(ebiten.KeyTab)
	KeyAlt            Key = Key(ebiten.KeyAlt)
	KeyControl        Key = Key(ebiten.KeyControl)
	KeyShift          Key = Key(ebiten.KeyShift)
	KeyMeta           Key = Key(ebiten.KeyMeta)
	KeyMax            Key = KeyMeta
)

// MarshalText implements encoding.TextMarshaler.
func (k Key) MarshalText() ([]byte, error) {
	return ebiten.Key(k).MarshalText()
}

// String returns a string representing the key.
// If k is an undefined key, String returns an empty string.
func (k Key) String() string {
	return ebiten.Key(k).String()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (k *Key) UnmarshalText(text []byte) error {
	v := ebiten.Key(0)
	err := v.UnmarshalText(text)
	if err != nil {
		return err
	}
	*k = Key(v)
	return nil
}

func DrawText(dst *Surface, face Face, color color.RGBA, x, y int, str string) {
	opts := text.DrawOptions{}
	opts.LineSpacing = float64(LineHeight(face))
	opts.GeoM.Translate(float64(x), float64(y))
	opts.ColorScale.Scale(
		float32(color.R)/255.0,
		float32(color.G)/255.0,
		float32(color.B)/255.0,
		float32(color.A)/255.0,
	)
	text.Draw(dst, str, face, &opts)
}

func MeasureText(txt string, face Face, lineSpacingInPixels float64) (width, height float64) {
	return text.Measure(txt, face, lineSpacingInPixels)
}

func LineHeight(face Face) int {
	return int(face.Metrics().HAscent + face.Metrics().HDescent + face.Metrics().HLineGap)
}

func FillRect(Surface *Surface, r Rectangle, col RGBA) {
	vector.DrawFilledRect(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
		col, false,
	)
}

func DrawRect(Surface *Surface, r Rectangle, stroke int, col RGBA) {
	vector.StrokeRect(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
		float32(stroke), col, false,
	)
}

// DrawsLine draws a line on the diagonal of the Rectangle r.
func DrawLine(Surface *Surface, r Rectangle, thick int, col RGBA) {
	vector.StrokeLine(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Max.X), float32(r.Max.Y),
		float32(thick), col, false,
	)
}

func DrawBox(Surface *Surface, r Rectangle, stroke int, fill, border, shadow RGBA) {
	if shadow.A != 0 {
		shadow.A = (shadow.A / 2) + 1 // make half transparent
		right := image.Rect(r.Max.X+1, r.Min.Y+1, r.Max.X+1, r.Max.Y+1)
		DrawLine(Surface, right, 1, shadow)
		bottom := image.Rect(r.Min.X+1, r.Max.Y+1, r.Max.X+1, r.Max.Y+1)
		DrawLine(Surface, bottom, 1, shadow)
	}

	vector.DrawFilledRect(
		Surface, float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()), fill, false,
	)

	if stroke > 0 {
		vector.StrokeRect(
			Surface, float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()),
			float32(stroke), border, false,
		)
	}
}

func DrawCircleInBox(Surface *Surface, box Rectangle, stroke int, fill, border RGBA) {
	r := box.Dx()
	if box.Dy() < r {
		r = box.Dy()
	}
	r = r / 2
	c := image.Pt((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2)
	DrawCircle(Surface, c, r, stroke, fill, border)
}

func DrawCircle(Surface *Surface, c Point, r int, stroke int, fill, border RGBA) {
	if r < 0 {
		r = 1
	}
	vector.DrawFilledCircle(Surface, float32(c.X), float32(c.Y),
		float32(r), fill, false)

	if stroke > 0 {
		vector.StrokeCircle(
			Surface, float32(c.X), float32(c.Y),
			float32(r), float32(stroke), border, false,
		)
	}
}

type PadID = ebiten.GamepadID
type PadButton = ebiten.GamepadButton
type PadAxis = ebiten.GamepadAxisType

type Inputter interface {
	Input() Inputter
}

type Pulse struct {
	Tick int64
}

func (c Pulse) Input() Inputter {
	return c
}

type PadInput struct {
	Pulse
	Pad PadID
}

type PadPlug PadInput
type PadYank PadInput

type PadButtonInput struct {
	PadInput
	Button   PadButton
	Duration int
}

type PadButtonPressInput PadButtonInput
type PadButtonHoldInput PadButtonInput
type PadButtonReleaseInput PadButtonInput

type PadAxisInput struct {
	PadInput
	Axis  PadAxis
	Value float64
}

type KeyInput struct {
	Pulse
	Key      Key
	Press    bool
	Release  bool
	Duration int
}

type KeyPressInput KeyInput
type KeyHoldInput KeyInput
type KeyReleaseInput KeyInput

type TextInput struct {
	Pulse
	ID   int
	Text []rune
}

type TouchID = ebiten.TouchID

type TouchInput struct {
	Pulse
	Touch    TouchID
	X        int
	Y        int
	DX       int
	DY       int
	Duration int
}

type TouchPressInput TouchInput
type TouchHoldInput TouchInput
type TouchReleaseInput TouchInput

type MouseInput struct {
	Pulse
	X        int
	Y        int
	DX       int
	DY       int
	Press    bool
	Duration int
}

type MousePressInput MouseInput
type MouseHoldInput MouseInput
type MouseMoveInput MouseInput
type MouseReleaseInput MouseInput

type WheelInput struct {
	Pulse
	DX float64
	DY float64
}

type InputState struct {
	PadIDs []PadID
	Fields []TextInputField
	MouseX int
	MouseY int
}

// Update updates the input state.
func (i *InputState) Poll(to chan<- Inputter) error {
	core := Pulse{Tick: ebiten.Tick()}

	for _, gid := range i.PadIDs {
		if inpututil.IsGamepadJustDisconnected(gid) {
			to <- PadYank{Pulse: core, Pad: gid}
		}
	}

	gamepadConnected := inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, gid := range gamepadConnected {
		to <- PadPlug{Pulse: core, Pad: gid}
	}

	i.PadIDs = i.PadIDs[0:0]
	i.PadIDs = ebiten.AppendGamepadIDs(i.PadIDs)
	for _, gid := range i.PadIDs {
		buttons := inpututil.AppendJustPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			to <- PadButtonPressInput{PadInput: PadInput{Pulse: core, Pad: gid}, Button: button}
		}

		buttons = inpututil.AppendPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			dur := inpututil.GamepadButtonPressDuration(gid, button)
			to <- PadButtonHoldInput{PadInput: PadInput{Pulse: core, Pad: gid}, Button: button, Duration: dur}
		}

		buttons = inpututil.AppendJustReleasedGamepadButtons(gid, nil)
		for _, button := range buttons {
			to <- PadButtonReleaseInput{PadInput: PadInput{Pulse: core, Pad: gid}, Button: button}

		}

		count := ebiten.GamepadAxisCount(gid)
		moved := false
		for axis := 0; axis < count; axis++ {
			value := ebiten.GamepadAxisValue(gid, axis)
			moved = ((value > 0.1) || (value < -0.1))
			if moved {
				to <- PadAxisInput{PadInput: PadInput{Pulse: core, Pad: gid}, Axis: axis, Value: value}
			}
		}
	}

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		to <- KeyPressInput{Pulse: core, Key: Key(key)}
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		dur := inpututil.KeyPressDuration(key)
		to <- KeyHoldInput{Pulse: core, Key: Key(key), Press: true, Duration: dur}
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		to <- KeyReleaseInput{Pulse: core, Key: Key(key)}
	}

	chars := ebiten.AppendInputChars(nil)
	if len(chars) > 0 {
		to <- TextInput{Pulse: core, Text: chars, ID: -1}
	}

	for id, field := range i.Fields {
		if field.IsFocused() {
			handled, _ := field.HandleInput(field.X, field.Y)
			if handled {
				to <- TextInput{Pulse: core, ID: id, Text: []rune(field.Text())}
			}
		}
	}

	touches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		to <- TouchPressInput{Pulse: core, Touch: touch, X: x, Y: y}
	}

	touches = ebiten.AppendTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		px, py := inpututil.TouchPositionInPreviousTick(touch)
		dx, dy := x-px, y-py
		dur := inpututil.TouchPressDuration(touch)
		to <- TouchHoldInput{Pulse: core, Touch: touch, X: x, Y: y,
			DX: dx, DY: dy, Duration: dur}
	}

	touches = inpututil.AppendJustReleasedTouchIDs(nil)
	for _, touch := range touches {
		x, y := ebiten.TouchPosition(touch)
		to <- TouchReleaseInput{Pulse: core, Touch: touch, X: x, Y: y}
	}

	x, y := ebiten.CursorPosition()
	dx, dy := x-i.MouseX, y-i.MouseY
	i.MouseX = x
	i.MouseY = y
	me := MouseInput{Pulse: core, X: x, Y: y, DX: dx, DY: dy}

	for mb := ebiten.MouseButton(0); mb < ebiten.MouseButtonMax; mb++ {
		if inpututil.IsMouseButtonJustPressed(mb) {
			to <- MousePressInput(me)
		}
		if ebiten.IsMouseButtonPressed(mb) {
			dur := inpututil.MouseButtonPressDuration(mb)
			mh := MouseHoldInput(me)
			mh.Duration = dur
			to <- mh
		}
		if inpututil.IsMouseButtonJustReleased(mb) {
			to <- MouseReleaseInput(me)
		}
	}
	if dx != 0 || dy != 0 {
		to <- MouseMoveInput(me)
	}

	wx, wy := ebiten.Wheel()
	if wx != 0 || wy != 0 {
		to <- WheelInput{Pulse: core, DX: wx, DY: wy}
	}

	return nil
}
