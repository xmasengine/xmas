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

type Color = color.Color
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

// Surface is an ebiten.Image
// We use surface to still be able to use image.Image together easily.
type Surface = ebiten.Image

// Face is a font face
type Face = text.Face

// DrawOptions are options for drawing a Surface.
type DrawOptions = ebiten.DrawImageOptions

// A KeyCode represents a physical keyboard key code as if it was typed on a QWERTY
// keyboard.
type KeyCode int

const (
	KeyA              KeyCode = KeyCode(ebiten.KeyA)
	KeyB              KeyCode = KeyCode(ebiten.KeyB)
	KeyC              KeyCode = KeyCode(ebiten.KeyC)
	KeyD              KeyCode = KeyCode(ebiten.KeyD)
	KeyE              KeyCode = KeyCode(ebiten.KeyE)
	KeyF              KeyCode = KeyCode(ebiten.KeyF)
	KeyG              KeyCode = KeyCode(ebiten.KeyG)
	KeyH              KeyCode = KeyCode(ebiten.KeyH)
	KeyI              KeyCode = KeyCode(ebiten.KeyI)
	KeyJ              KeyCode = KeyCode(ebiten.KeyJ)
	KeyK              KeyCode = KeyCode(ebiten.KeyK)
	KeyL              KeyCode = KeyCode(ebiten.KeyL)
	KeyM              KeyCode = KeyCode(ebiten.KeyM)
	KeyN              KeyCode = KeyCode(ebiten.KeyN)
	KeyO              KeyCode = KeyCode(ebiten.KeyO)
	KeyP              KeyCode = KeyCode(ebiten.KeyP)
	KeyQ              KeyCode = KeyCode(ebiten.KeyQ)
	KeyR              KeyCode = KeyCode(ebiten.KeyR)
	KeyS              KeyCode = KeyCode(ebiten.KeyS)
	KeyT              KeyCode = KeyCode(ebiten.KeyT)
	KeyU              KeyCode = KeyCode(ebiten.KeyU)
	KeyV              KeyCode = KeyCode(ebiten.KeyV)
	KeyW              KeyCode = KeyCode(ebiten.KeyW)
	KeyX              KeyCode = KeyCode(ebiten.KeyX)
	KeyY              KeyCode = KeyCode(ebiten.KeyY)
	KeyZ              KeyCode = KeyCode(ebiten.KeyZ)
	KeyAltLeft        KeyCode = KeyCode(ebiten.KeyAltLeft)
	KeyAltRight       KeyCode = KeyCode(ebiten.KeyAltRight)
	KeyArrowDown      KeyCode = KeyCode(ebiten.KeyArrowDown)
	KeyArrowLeft      KeyCode = KeyCode(ebiten.KeyArrowLeft)
	KeyArrowRight     KeyCode = KeyCode(ebiten.KeyArrowRight)
	KeyArrowUp        KeyCode = KeyCode(ebiten.KeyArrowUp)
	KeyBackquote      KeyCode = KeyCode(ebiten.KeyBackquote)
	KeyBackslash      KeyCode = KeyCode(ebiten.KeyBackslash)
	KeyBackspace      KeyCode = KeyCode(ebiten.KeyBackspace)
	KeyBracketLeft    KeyCode = KeyCode(ebiten.KeyBracketLeft)
	KeyBracketRight   KeyCode = KeyCode(ebiten.KeyBracketRight)
	KeyCapsLock       KeyCode = KeyCode(ebiten.KeyCapsLock)
	KeyComma          KeyCode = KeyCode(ebiten.KeyComma)
	KeyContextMenu    KeyCode = KeyCode(ebiten.KeyContextMenu)
	KeyControlLeft    KeyCode = KeyCode(ebiten.KeyControlLeft)
	KeyControlRight   KeyCode = KeyCode(ebiten.KeyControlRight)
	KeyDelete         KeyCode = KeyCode(ebiten.KeyDelete)
	KeyDigit0         KeyCode = KeyCode(ebiten.KeyDigit0)
	KeyDigit1         KeyCode = KeyCode(ebiten.KeyDigit1)
	KeyDigit2         KeyCode = KeyCode(ebiten.KeyDigit2)
	KeyDigit3         KeyCode = KeyCode(ebiten.KeyDigit3)
	KeyDigit4         KeyCode = KeyCode(ebiten.KeyDigit4)
	KeyDigit5         KeyCode = KeyCode(ebiten.KeyDigit5)
	KeyDigit6         KeyCode = KeyCode(ebiten.KeyDigit6)
	KeyDigit7         KeyCode = KeyCode(ebiten.KeyDigit7)
	KeyDigit8         KeyCode = KeyCode(ebiten.KeyDigit8)
	KeyDigit9         KeyCode = KeyCode(ebiten.KeyDigit9)
	KeyEnd            KeyCode = KeyCode(ebiten.KeyEnd)
	KeyEnter          KeyCode = KeyCode(ebiten.KeyEnter)
	KeyEqual          KeyCode = KeyCode(ebiten.KeyEqual)
	KeyEscape         KeyCode = KeyCode(ebiten.KeyEscape)
	KeyF1             KeyCode = KeyCode(ebiten.KeyF1)
	KeyF2             KeyCode = KeyCode(ebiten.KeyF2)
	KeyF3             KeyCode = KeyCode(ebiten.KeyF3)
	KeyF4             KeyCode = KeyCode(ebiten.KeyF4)
	KeyF5             KeyCode = KeyCode(ebiten.KeyF5)
	KeyF6             KeyCode = KeyCode(ebiten.KeyF6)
	KeyF7             KeyCode = KeyCode(ebiten.KeyF7)
	KeyF8             KeyCode = KeyCode(ebiten.KeyF8)
	KeyF9             KeyCode = KeyCode(ebiten.KeyF9)
	KeyF10            KeyCode = KeyCode(ebiten.KeyF10)
	KeyF11            KeyCode = KeyCode(ebiten.KeyF11)
	KeyF12            KeyCode = KeyCode(ebiten.KeyF12)
	KeyF13            KeyCode = KeyCode(ebiten.KeyF13)
	KeyF14            KeyCode = KeyCode(ebiten.KeyF14)
	KeyF15            KeyCode = KeyCode(ebiten.KeyF15)
	KeyF16            KeyCode = KeyCode(ebiten.KeyF16)
	KeyF17            KeyCode = KeyCode(ebiten.KeyF17)
	KeyF18            KeyCode = KeyCode(ebiten.KeyF18)
	KeyF19            KeyCode = KeyCode(ebiten.KeyF19)
	KeyF20            KeyCode = KeyCode(ebiten.KeyF20)
	KeyF21            KeyCode = KeyCode(ebiten.KeyF21)
	KeyF22            KeyCode = KeyCode(ebiten.KeyF22)
	KeyF23            KeyCode = KeyCode(ebiten.KeyF23)
	KeyF24            KeyCode = KeyCode(ebiten.KeyF24)
	KeyHome           KeyCode = KeyCode(ebiten.KeyHome)
	KeyInsert         KeyCode = KeyCode(ebiten.KeyInsert)
	KeyIntlBackslash  KeyCode = KeyCode(ebiten.KeyIntlBackslash)
	KeyMetaLeft       KeyCode = KeyCode(ebiten.KeyMetaLeft)
	KeyMetaRight      KeyCode = KeyCode(ebiten.KeyMetaRight)
	KeyMinus          KeyCode = KeyCode(ebiten.KeyMinus)
	KeyNumLock        KeyCode = KeyCode(ebiten.KeyNumLock)
	KeyNumpad0        KeyCode = KeyCode(ebiten.KeyNumpad0)
	KeyNumpad1        KeyCode = KeyCode(ebiten.KeyNumpad1)
	KeyNumpad2        KeyCode = KeyCode(ebiten.KeyNumpad2)
	KeyNumpad3        KeyCode = KeyCode(ebiten.KeyNumpad3)
	KeyNumpad4        KeyCode = KeyCode(ebiten.KeyNumpad4)
	KeyNumpad5        KeyCode = KeyCode(ebiten.KeyNumpad5)
	KeyNumpad6        KeyCode = KeyCode(ebiten.KeyNumpad6)
	KeyNumpad7        KeyCode = KeyCode(ebiten.KeyNumpad7)
	KeyNumpad8        KeyCode = KeyCode(ebiten.KeyNumpad8)
	KeyNumpad9        KeyCode = KeyCode(ebiten.KeyNumpad9)
	KeyNumpadAdd      KeyCode = KeyCode(ebiten.KeyNumpadAdd)
	KeyNumpadDecimal  KeyCode = KeyCode(ebiten.KeyNumpadDecimal)
	KeyNumpadDivide   KeyCode = KeyCode(ebiten.KeyNumpadDivide)
	KeyNumpadEnter    KeyCode = KeyCode(ebiten.KeyNumpadEnter)
	KeyNumpadEqual    KeyCode = KeyCode(ebiten.KeyNumpadEqual)
	KeyNumpadMultiply KeyCode = KeyCode(ebiten.KeyNumpadMultiply)
	KeyNumpadSubtract KeyCode = KeyCode(ebiten.KeyNumpadSubtract)
	KeyPageDown       KeyCode = KeyCode(ebiten.KeyPageDown)
	KeyPageUp         KeyCode = KeyCode(ebiten.KeyPageUp)
	KeyPause          KeyCode = KeyCode(ebiten.KeyPause)
	KeyPeriod         KeyCode = KeyCode(ebiten.KeyPeriod)
	KeyPrintScreen    KeyCode = KeyCode(ebiten.KeyPrintScreen)
	KeyQuote          KeyCode = KeyCode(ebiten.KeyQuote)
	KeyScrollLock     KeyCode = KeyCode(ebiten.KeyScrollLock)
	KeySemicolon      KeyCode = KeyCode(ebiten.KeySemicolon)
	KeyShiftLeft      KeyCode = KeyCode(ebiten.KeyShiftLeft)
	KeyShiftRight     KeyCode = KeyCode(ebiten.KeyShiftRight)
	KeySlash          KeyCode = KeyCode(ebiten.KeySlash)
	KeySpace          KeyCode = KeyCode(ebiten.KeySpace)
	KeyTab            KeyCode = KeyCode(ebiten.KeyTab)
	KeyAlt            KeyCode = KeyCode(ebiten.KeyAlt)
	KeyControl        KeyCode = KeyCode(ebiten.KeyControl)
	KeyShift          KeyCode = KeyCode(ebiten.KeyShift)
	KeyMeta           KeyCode = KeyCode(ebiten.KeyMeta)
	KeyMax            KeyCode = KeyMeta
)

// MarshalText implements encoding.TextMarshaler.
func (k KeyCode) MarshalText() ([]byte, error) {
	return ebiten.Key(k).MarshalText()
}

// String returns a string representing the key.
// If k is an undefined key, String returns an empty string.
func (k KeyCode) String() string {
	return ebiten.Key(k).String()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (k *KeyCode) UnmarshalText(text []byte) error {
	v := ebiten.Key(0)
	err := v.UnmarshalText(text)
	if err != nil {
		return err
	}
	*k = KeyCode(v)
	return nil
}

func DrawText(dst *Surface, face Face, color RGBA, x, y int, str string) {
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

// Keymods are the current key modifers.
type KeyMods struct {
	Alt     bool
	Control bool
	Shift   bool
	Meta    bool
}

type PadID = ebiten.GamepadID
type TriggerID = ebiten.GamepadButton
type AxisID = ebiten.GamepadAxisType

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

type TriggerInput struct {
	PadInput
	Trigger  TriggerID
	Duration int
}

type TriggerPressInput TriggerInput
type TriggerHoldInput TriggerInput
type TriggerReleaseInput TriggerInput

type AxisInput struct {
	PadInput
	Axis  AxisID
	Value float64
}

type KeyInput struct {
	Pulse
	Key      KeyCode
	Mods     KeyMods
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
	Mods   KeyMods
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
			to <- TriggerPressInput{PadInput: PadInput{Pulse: core, Pad: gid}, Trigger: button}
		}

		buttons = inpututil.AppendPressedGamepadButtons(gid, nil)
		for _, button := range buttons {
			dur := inpututil.GamepadButtonPressDuration(gid, button)
			to <- TriggerHoldInput{PadInput: PadInput{Pulse: core, Pad: gid}, Trigger: button, Duration: dur}
		}

		buttons = inpututil.AppendJustReleasedGamepadButtons(gid, nil)
		for _, button := range buttons {
			to <- TriggerReleaseInput{PadInput: PadInput{Pulse: core, Pad: gid}, Trigger: button}

		}

		count := ebiten.GamepadAxisCount(gid)
		moved := false
		for axis := 0; axis < count; axis++ {
			value := ebiten.GamepadAxisValue(gid, axis)
			moved = ((value > 0.1) || (value < -0.1))
			if moved {
				to <- AxisInput{PadInput: PadInput{Pulse: core, Pad: gid}, Axis: axis, Value: value}
			}
		}
	}

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		switch KeyCode(key) {
		case KeyAlt:
			i.Mods.Alt = true
		case KeyShift:
			i.Mods.Shift = true
		case KeyControl:
			i.Mods.Control = true
		case KeyMeta:
			i.Mods.Meta = true
		}

		to <- KeyPressInput{Pulse: core, Key: KeyCode(key)}
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		dur := inpututil.KeyPressDuration(key)
		to <- KeyHoldInput{Pulse: core, Key: KeyCode(key), Press: true, Duration: dur, Mods: i.Mods}
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		switch KeyCode(key) {
		case KeyAlt:
			i.Mods.Alt = false
		case KeyShift:
			i.Mods.Shift = false
		case KeyControl:
			i.Mods.Control = false
		case KeyMeta:
			i.Mods.Meta = false
		}
		to <- KeyReleaseInput{Pulse: core, Key: KeyCode(key)}
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
