package xgal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyCode is a keyboard key code.
type KeyCode = ebiten.Key

const (
	KeyA              KeyCode = ebiten.KeyA
	KeyB              KeyCode = ebiten.KeyB
	KeyC              KeyCode = ebiten.KeyC
	KeyD              KeyCode = ebiten.KeyD
	KeyE              KeyCode = ebiten.KeyE
	KeyF              KeyCode = ebiten.KeyF
	KeyG              KeyCode = ebiten.KeyG
	KeyH              KeyCode = ebiten.KeyH
	KeyI              KeyCode = ebiten.KeyI
	KeyJ              KeyCode = ebiten.KeyJ
	KeyK              KeyCode = ebiten.KeyK
	KeyL              KeyCode = ebiten.KeyL
	KeyM              KeyCode = ebiten.KeyM
	KeyN              KeyCode = ebiten.KeyN
	KeyO              KeyCode = ebiten.KeyO
	KeyP              KeyCode = ebiten.KeyP
	KeyQ              KeyCode = ebiten.KeyQ
	KeyR              KeyCode = ebiten.KeyR
	KeyS              KeyCode = ebiten.KeyS
	KeyT              KeyCode = ebiten.KeyT
	KeyU              KeyCode = ebiten.KeyU
	KeyV              KeyCode = ebiten.KeyV
	KeyW              KeyCode = ebiten.KeyW
	KeyX              KeyCode = ebiten.KeyX
	KeyY              KeyCode = ebiten.KeyY
	KeyZ              KeyCode = ebiten.KeyZ
	KeyAltLeft        KeyCode = ebiten.KeyAltLeft
	KeyAltRight       KeyCode = ebiten.KeyAltRight
	KeyArrowDown      KeyCode = ebiten.KeyArrowDown
	KeyArrowLeft      KeyCode = ebiten.KeyArrowLeft
	KeyArrowRight     KeyCode = ebiten.KeyArrowRight
	KeyArrowUp        KeyCode = ebiten.KeyArrowUp
	KeyBackquote      KeyCode = ebiten.KeyBackquote
	KeyBackslash      KeyCode = ebiten.KeyBackslash
	KeyBackspace      KeyCode = ebiten.KeyBackspace
	KeyBracketLeft    KeyCode = ebiten.KeyBracketLeft
	KeyBracketRight   KeyCode = ebiten.KeyBracketRight
	KeyCapsLock       KeyCode = ebiten.KeyCapsLock
	KeyComma          KeyCode = ebiten.KeyComma
	KeyContextMenu    KeyCode = ebiten.KeyContextMenu
	KeyControlLeft    KeyCode = ebiten.KeyControlLeft
	KeyControlRight   KeyCode = ebiten.KeyControlRight
	KeyDelete         KeyCode = ebiten.KeyDelete
	KeyDigit0         KeyCode = ebiten.KeyDigit0
	KeyDigit1         KeyCode = ebiten.KeyDigit1
	KeyDigit2         KeyCode = ebiten.KeyDigit2
	KeyDigit3         KeyCode = ebiten.KeyDigit3
	KeyDigit4         KeyCode = ebiten.KeyDigit4
	KeyDigit5         KeyCode = ebiten.KeyDigit5
	KeyDigit6         KeyCode = ebiten.KeyDigit6
	KeyDigit7         KeyCode = ebiten.KeyDigit7
	KeyDigit8         KeyCode = ebiten.KeyDigit8
	KeyDigit9         KeyCode = ebiten.KeyDigit9
	KeyEnd            KeyCode = ebiten.KeyEnd
	KeyEnter          KeyCode = ebiten.KeyEnter
	KeyEqual          KeyCode = ebiten.KeyEqual
	KeyEscape         KeyCode = ebiten.KeyEscape
	KeyF1             KeyCode = ebiten.KeyF1
	KeyF2             KeyCode = ebiten.KeyF2
	KeyF3             KeyCode = ebiten.KeyF3
	KeyF4             KeyCode = ebiten.KeyF4
	KeyF5             KeyCode = ebiten.KeyF5
	KeyF6             KeyCode = ebiten.KeyF6
	KeyF7             KeyCode = ebiten.KeyF7
	KeyF8             KeyCode = ebiten.KeyF8
	KeyF9             KeyCode = ebiten.KeyF9
	KeyF10            KeyCode = ebiten.KeyF10
	KeyF11            KeyCode = ebiten.KeyF11
	KeyF12            KeyCode = ebiten.KeyF12
	KeyF13            KeyCode = ebiten.KeyF13
	KeyF14            KeyCode = ebiten.KeyF14
	KeyF15            KeyCode = ebiten.KeyF15
	KeyF16            KeyCode = ebiten.KeyF16
	KeyF17            KeyCode = ebiten.KeyF17
	KeyF18            KeyCode = ebiten.KeyF18
	KeyF19            KeyCode = ebiten.KeyF19
	KeyF20            KeyCode = ebiten.KeyF20
	KeyF21            KeyCode = ebiten.KeyF21
	KeyF22            KeyCode = ebiten.KeyF22
	KeyF23            KeyCode = ebiten.KeyF23
	KeyF24            KeyCode = ebiten.KeyF24
	KeyHome           KeyCode = ebiten.KeyHome
	KeyInsert         KeyCode = ebiten.KeyInsert
	KeyIntlBackslash  KeyCode = ebiten.KeyIntlBackslash
	KeyMetaLeft       KeyCode = ebiten.KeyMetaLeft
	KeyMetaRight      KeyCode = ebiten.KeyMetaRight
	KeyMinus          KeyCode = ebiten.KeyMinus
	KeyNumLock        KeyCode = ebiten.KeyNumLock
	KeyNumpad0        KeyCode = ebiten.KeyNumpad0
	KeyNumpad1        KeyCode = ebiten.KeyNumpad1
	KeyNumpad2        KeyCode = ebiten.KeyNumpad2
	KeyNumpad3        KeyCode = ebiten.KeyNumpad3
	KeyNumpad4        KeyCode = ebiten.KeyNumpad4
	KeyNumpad5        KeyCode = ebiten.KeyNumpad5
	KeyNumpad6        KeyCode = ebiten.KeyNumpad6
	KeyNumpad7        KeyCode = ebiten.KeyNumpad7
	KeyNumpad8        KeyCode = ebiten.KeyNumpad8
	KeyNumpad9        KeyCode = ebiten.KeyNumpad9
	KeyNumpadAdd      KeyCode = ebiten.KeyNumpadAdd
	KeyNumpadDecimal  KeyCode = ebiten.KeyNumpadDecimal
	KeyNumpadDivide   KeyCode = ebiten.KeyNumpadDivide
	KeyNumpadEnter    KeyCode = ebiten.KeyNumpadEnter
	KeyNumpadEqual    KeyCode = ebiten.KeyNumpadEqual
	KeyNumpadMultiply KeyCode = ebiten.KeyNumpadMultiply
	KeyNumpadSubtract KeyCode = ebiten.KeyNumpadSubtract
	KeyPageDown       KeyCode = ebiten.KeyPageDown
	KeyPageUp         KeyCode = ebiten.KeyPageUp
	KeyPause          KeyCode = ebiten.KeyPause
	KeyPeriod         KeyCode = ebiten.KeyPeriod
	KeyPrintScreen    KeyCode = ebiten.KeyPrintScreen
	KeyQuote          KeyCode = ebiten.KeyQuote
	KeyScrollLock     KeyCode = ebiten.KeyScrollLock
	KeySemicolon      KeyCode = ebiten.KeySemicolon
	KeyShiftLeft      KeyCode = ebiten.KeyShiftLeft
	KeyShiftRight     KeyCode = ebiten.KeyShiftRight
	KeySlash          KeyCode = ebiten.KeySlash
	KeySpace          KeyCode = ebiten.KeySpace
	KeyTab            KeyCode = ebiten.KeyTab
	KeyAlt            KeyCode = ebiten.KeyAlt
	KeyControl        KeyCode = ebiten.KeyControl
	KeyShift          KeyCode = ebiten.KeyShift
	KeyMeta           KeyCode = ebiten.KeyMeta
	KeyMax            KeyCode = ebiten.KeyMax
)

// Key reports whether the key is currently pressed.
func Key(code KeyCode) bool {
	return ebiten.IsKeyPressed(code)
}

// Tap reports whether the key was just pressed this frame.
func Tap(code KeyCode) bool {
	return inpututil.IsKeyJustPressed(code)
}

// Lift reports whether the key was just released this frame.
func Lift(code KeyCode) bool {
	return inpututil.IsKeyJustReleased(code)
}

// Keys returns all keys that are currently pressed.
// If buf is provided, results are appended to it.
func Keys(buf ...[]KeyCode) []KeyCode {
	var b []KeyCode
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendPressedKeys(b)
}

// Taps returns all keys that were just pressed this frame.
// If buf is provided, results are appended to it.
func Taps(buf ...[]KeyCode) []KeyCode {
	var b []KeyCode
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustPressedKeys(b)
}

// Lifts returns all keys that were just released this frame.
// If buf is provided, results are appended to it.
func Lifts(buf ...[]KeyCode) []KeyCode {
	var b []KeyCode
	if len(buf) > 0 {
		b = buf[0]
	}
	return inpututil.AppendJustReleasedKeys(b)
}

// Chars returns the text input (IME) characters entered this frame.
// If buf is provided, results are appended to it.
func Chars(buf ...[]rune) []rune {
	var b []rune
	if len(buf) > 0 {
		b = buf[0]
	}
	return ebiten.AppendInputChars(b)
}

// Age returns how long the key has been pressed, in ticks.
func Age(code KeyCode) int {
	return inpututil.KeyPressDuration(code)
}
