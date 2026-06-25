// Package xio is a wrapper around the low level ebitengine game library
// and its supporting libraries, as well as the image and image/color library.
// While ebitengine works well, it is a bit of a hassle to import everything
// separately.
// Also this wrapper could be useful in case we have to use a different
// low level game library.
package xio

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputUtil struct{}

type GamepadUtil struct{}

func (GamepadUtil) Axis(id ebiten.GamepadID, axis ebiten.GamepadAxisType) float64 {
	return ebiten.GamepadAxis(id, axis)
}

func (GamepadUtil) AxisCount(id ebiten.GamepadID) int {
	return ebiten.GamepadAxisCount(id)
}

func (GamepadUtil) AxisNum(id ebiten.GamepadID) int {
	return ebiten.GamepadAxisNum(id)
}

func (GamepadUtil) AxisValue(id ebiten.GamepadID, axis ebiten.GamepadAxisType) float64 {
	return ebiten.GamepadAxisValue(id, axis)
}

func (GamepadUtil) ButtonCount(id ebiten.GamepadID) int {
	return ebiten.GamepadButtonCount(id)
}

func (GamepadUtil) ButtonNum(id ebiten.GamepadID) int {
	return ebiten.GamepadButtonNum(id)
}

func (GamepadUtil) Name(id ebiten.GamepadID) string {
	return ebiten.GamepadName(id)
}

func (GamepadUtil) SDLID(id ebiten.GamepadID) string {
	return ebiten.GamepadSDLID(id)
}

func (GamepadUtil) AppendJustConnectedIDs(gamepadIDs []ebiten.GamepadID) []ebiten.GamepadID {
	return inpututil.AppendJustConnectedGamepadIDs(gamepadIDs)
}

func (GamepadUtil) AppendJustPressedButtons(id ebiten.GamepadID, buttons []ebiten.GamepadButton) []ebiten.GamepadButton {
	return inpututil.AppendJustPressedGamepadButtons(id, buttons)
}

func (GamepadUtil) AppendJustReleasedButtons(id ebiten.GamepadID, buttons []ebiten.GamepadButton) []ebiten.GamepadButton {
	return inpututil.AppendJustReleasedGamepadButtons(id, buttons)
}

func (GamepadUtil) AppendPressedButtons(id ebiten.GamepadID, buttons []ebiten.GamepadButton) []ebiten.GamepadButton {
	return inpututil.AppendPressedGamepadButtons(id, buttons)
}

func (GamepadUtil) ButtonPressDuration(id ebiten.GamepadID, button ebiten.GamepadButton) int {
	return inpututil.GamepadButtonPressDuration(id, button)
}

func (GamepadUtil) IsButtonJustPressed(id ebiten.GamepadID, button ebiten.GamepadButton) bool {
	return inpututil.IsGamepadButtonJustPressed(id, button)
}

func (GamepadUtil) IsButtonJustReleased(id ebiten.GamepadID, button ebiten.GamepadButton) bool {
	return inpututil.IsGamepadButtonJustReleased(id, button)
}

func (GamepadUtil) IsJustDisconnected(id ebiten.GamepadID) bool {
	return inpututil.IsGamepadJustDisconnected(id)
}

func (GamepadUtil) JustConnectedIDs() []ebiten.GamepadID {
	return inpututil.JustConnectedGamepadIDs()
}

type StandardGamepadUtil struct{}

func (StandardGamepadUtil) AppendPressedButtons(id ebiten.GamepadID, buttons []ebiten.StandardGamepadButton) []ebiten.StandardGamepadButton {
	return inpututil.AppendPressedStandardGamepadButtons(id, buttons)
}

func (StandardGamepadUtil) AppendJustReleasedButtons(id ebiten.GamepadID, buttons []ebiten.StandardGamepadButton) []ebiten.StandardGamepadButton {
	return inpututil.AppendJustReleasedStandardGamepadButtons(id, buttons)
}

func (StandardGamepadUtil) AppendJustPressedButtons(id ebiten.GamepadID, buttons []ebiten.StandardGamepadButton) []ebiten.StandardGamepadButton {
	return inpututil.AppendJustPressedStandardGamepadButtons(id, buttons)
}

func (StandardGamepadUtil) IsStandardGamepadButtonJustPressed(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool {
	return inpututil.IsStandardGamepadButtonJustPressed(id, button)
}

func (StandardGamepadUtil) IsStandardGamepadButtonJustReleased(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool {
	return inpututil.IsStandardGamepadButtonJustReleased(id, button)
}

func (StandardGamepadUtil) StandardGamepadButtonPressDuration(id ebiten.GamepadID, button ebiten.StandardGamepadButton) int {
	return inpututil.StandardGamepadButtonPressDuration(id, button)
}

type KeyUtil struct{}

func (KeyUtil) AppendJustPressed(keys []ebiten.Key) []ebiten.Key {
	return inpututil.AppendJustPressedKeys(keys)
}

func (KeyUtil) AppendJustReleased(keys []ebiten.Key) []ebiten.Key {
	return inpututil.AppendJustReleasedKeys(keys)
}

func (KeyUtil) AppendPressed(keys []ebiten.Key) []ebiten.Key {
	return inpututil.AppendPressedKeys(keys)
}

func (KeyUtil) IsJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func (KeyUtil) IsJustReleased(key ebiten.Key) bool {
	return inpututil.IsKeyJustReleased(key)
}

func (KeyUtil) PressDuration(key ebiten.Key) int {
	return inpututil.KeyPressDuration(key)
}

func (KeyUtil) PressedKeys() []ebiten.Key {
	return inpututil.PressedKeys()
}

type MouseUtil struct{}

func (MouseUtil) IsButtonJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

func (MouseUtil) IsButtonJustReleased(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustReleased(button)
}

func (MouseUtil) ButtonPressDuration(button ebiten.MouseButton) int {
	return inpututil.MouseButtonPressDuration(button)
}

func (MouseUtil) Position() (x, y int) {
	return ebiten.CursorPosition()
}

type TouchUtil struct{}

func (TouchUtil) PositionInPreviousTick(id ebiten.TouchID) (int, int) {
	return inpututil.TouchPositionInPreviousTick(id)
}

func (TouchUtil) PressDuration(id ebiten.TouchID) int {
	return inpututil.TouchPressDuration(id)
}

func (TouchUtil) AppendJustPressedIDs(touchIDs []ebiten.TouchID) []ebiten.TouchID {
	return inpututil.AppendJustPressedTouchIDs(touchIDs)
}

func (TouchUtil) AppendJustReleasedIDs(touchIDs []ebiten.TouchID) []ebiten.TouchID {
	return inpututil.AppendJustReleasedTouchIDs(touchIDs)
}

func (TouchUtil) IsJustReleased(id ebiten.TouchID) bool {
	return inpututil.IsTouchJustReleased(id)
}

func (TouchUtil) JustPressedIDs() []ebiten.TouchID {
	return inpututil.JustPressedTouchIDs()
}

/*

func AppendInputChars(runes []rune) []rune


func InputChars() []rune
func IsFocused() bool
func IsGamepadButtonPressed(id GamepadID, button GamepadButton) bool
func IsKeyPressed(key Key) bool
func IsMouseButtonPressed(mouseButton MouseButton) bool
func IsStandardGamepadAxisAvailable(id GamepadID, axis StandardGamepadAxis) bool
func IsStandardGamepadButtonAvailable(id GamepadID, button StandardGamepadButton) bool
func IsStandardGamepadButtonPressed(id GamepadID, button StandardGamepadButton) bool
func IsStandardGamepadLayoutAvailable(id GamepadID) bool
func KeyName(key Key) string
func StandardGamepadAxisValue(id GamepadID, axis StandardGamepadAxis) float64
func StandardGamepadButtonValue(id GamepadID, button StandardGamepadButton) float64
func Tick() int64
func TouchPosition(id TouchID) (int, int)
func UpdateStandardGamepadLayoutMappings(mappings string) (bool, error)
func Vibrate(options *VibrateOptions)
func VibrateGamepad(gamepadID GamepadID, options *VibrateGamepadOptions)
func Wheel() (xoff, yoff float64)

type GamepadButton = gamepad.Button
    const GamepadButton0 GamepadButton = gamepad.Button0 ...
type GamepadID = gamepad.ID
    func AppendGamepadIDs(gamepadIDs []GamepadID) []GamepadID
    func GamepadIDs() []GamepadID
type Key int
    const KeyA Key = Key(ui.KeyA) ...
type MouseButton int
    const MouseButtonLeft MouseButton = MouseButton0 ...
type StandardGamepadAxis = gamepaddb.StandardAxis
    const StandardGamepadAxisLeftStickHorizontal StandardGamepadAxis = gamepaddb.StandardAxisLeftStickHorizontal ...
type StandardGamepadButton = gamepaddb.StandardButton
    const StandardGamepadButtonRightBottom StandardGamepadButton = gamepaddb.StandardButtonRightBottom ...
type TouchID int

    func AppendTouchIDs(touches []TouchID) []TouchID
    func TouchIDs() []TouchID

*/
