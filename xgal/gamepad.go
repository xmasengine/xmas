package xgal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PadID identifies a connected gamepad.
type PadID = ebiten.GamepadID

// Button is a gamepad button.
type Button = ebiten.GamepadButton

// AxisID identifies a gamepad axis.
type AxisID = ebiten.GamepadAxisType

// Plugs returns all gamepads that were just connected this frame.
// If buf is provided, results are appended to it.
func Plugs(buf ...[]PadID) []PadID {
	var b []PadID
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustConnectedGamepadIDs(b)
}

// Yank reports whether the gamepad was just disconnected this frame.
func Yank(pad PadID) bool {
	return inpututil.IsGamepadJustDisconnected(pad)
}

// Pads returns all currently connected gamepads.
// If buf is provided, results are appended to it.
func Pads(buf ...[]PadID) []PadID {
	var b []PadID
	if len(buf) > 0 {
		b = buf[0]
	}
	return ebiten.AppendGamepadIDs(b)
}

// Nudge reports whether the gamepad button was just pressed.
func Nudge(pad PadID, btn Button) bool {
	return inpututil.IsGamepadButtonJustPressed(pad, btn)
}

// Squeeze reports whether the gamepad button is currently held.
func Squeeze(pad PadID, btn Button) bool {
	for _, b := range inpututil.AppendPressedGamepadButtons(pad, nil) {
		if b == btn {
			return true
		}
	}
	return false
}

// Slip reports whether the gamepad button was just released.
func Slip(pad PadID, btn Button) bool {
	return inpututil.IsGamepadButtonJustReleased(pad, btn)
}

// Nudges returns all buttons that were just pressed on the given pad this frame.
// If buf is provided, results are appended to it.
func Nudges(pad PadID, buf ...[]Button) []Button {
	var b []Button
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustPressedGamepadButtons(pad, b)
}

// Squeezes returns all buttons that are currently held on the given pad.
// If buf is provided, results are appended to it.
func Squeezes(pad PadID, buf ...[]Button) []Button {
	var b []Button
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendPressedGamepadButtons(pad, b)
}

// Slips returns all buttons that were just released on the given pad this frame.
// If buf is provided, results are appended to it.
func Slips(pad PadID, buf ...[]Button) []Button {
	var b []Button
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustReleasedGamepadButtons(pad, b)
}

// Axis returns the current value of the given gamepad axis.
func Axis(pad PadID, axis AxisID) float64 {
	return ebiten.GamepadAxisValue(pad, axis)
}

// Axes returns the number of axes on the given gamepad.
func Axes(pad PadID) int {
	return ebiten.GamepadAxisCount(pad)
}
