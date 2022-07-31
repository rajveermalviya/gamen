//go:build js

package web

import (
	"github.com/rajveermalviya/gamen/events"
)

func mapKeyCode(code string) (events.VirtualKey, bool) {
	switch code {
	// case "":
	// 	return events.VirtualKeyCancel, true
	case "Help":
		return events.VirtualKeyHelp, true
	case "Backspace":
		return events.VirtualKeyBackSpace, true
	case "Tab":
		return events.VirtualKeyTab, true
	// case "":
	// 	return events.VirtualKeyClear, true
	case "Enter":
		return events.VirtualKeyReturn, true
	case "ShiftLeft":
		return events.VirtualKeyLShift, true
	case "ShiftRight":
		return events.VirtualKeyRShift, true
	case "ControlLeft":
		return events.VirtualKeyLControl, true
	case "ControlRight":
		return events.VirtualKeyRControl, true
	case "AltLeft":
		return events.VirtualKeyLAlt, true
	case "AltRight":
		return events.VirtualKeyRAlt, true
	case "Pause":
		return events.VirtualKeyPause, true
	case "CapsLock":
		return events.VirtualKeyCapsLock, true
	case "KanaMode":
		return events.VirtualKeyKana, true
	// case "":
	// 	return events.VirtualKeyJunja, true
	// case "":
	// 	return events.VirtualKeyFinal, true
	// case "":
	// 	return events.VirtualKeyKanji, true
	case "Escape":
		return events.VirtualKeyEscape, true
	case "Convert":
		return events.VirtualKeyConvert, true
	case "NonConvert":
		return events.VirtualKeyNonconvert, true
	// case "":
	// 	return events.VirtualKeyAccept, true
	// case "":
	// 	return events.VirtualKeyModechange, true
	case "Space":
		return events.VirtualKeySpace, true
	case "PageUp":
		return events.VirtualKeyPageUp, true
	case "PageDown":
		return events.VirtualKeyPageDown, true
	case "End":
		return events.VirtualKeyEnd, true
	case "Home":
		return events.VirtualKeyHome, true
	case "ArrowLeft":
		return events.VirtualKeyLeft, true
	case "ArrowUp":
		return events.VirtualKeyUp, true
	case "ArrowRight":
		return events.VirtualKeyRight, true
	case "ArrowDown":
		return events.VirtualKeyDown, true
	case "Select":
		return events.VirtualKeySelect, true
	// case "":
	// 	return events.VirtualKeyPrint, true
	// case "":
	// 	return events.VirtualKeyExecute, true
	case "Insert":
		return events.VirtualKeyInsert, true
	case "Delete":
		return events.VirtualKeyDelete, true
	case "Digit0":
		return events.VirtualKey0, true
	case "Digit1":
		return events.VirtualKey1, true
	case "Digit2":
		return events.VirtualKey2, true
	case "Digit3":
		return events.VirtualKey3, true
	case "Digit4":
		return events.VirtualKey4, true
	case "Digit5":
		return events.VirtualKey5, true
	case "Digit6":
		return events.VirtualKey6, true
	case "Digit7":
		return events.VirtualKey7, true
	case "Digit8":
		return events.VirtualKey8, true
	case "Digit9":
		return events.VirtualKey9, true
	case "Semicolon":
		return events.VirtualKeySemicolon, true
	case "Equal":
		return events.VirtualKeyEquals, true
	case "KeyA":
		return events.VirtualKeyA, true
	case "KeyB":
		return events.VirtualKeyB, true
	case "KeyC":
		return events.VirtualKeyC, true
	case "KeyD":
		return events.VirtualKeyD, true
	case "KeyE":
		return events.VirtualKeyE, true
	case "KeyF":
		return events.VirtualKeyF, true
	case "KeyG":
		return events.VirtualKeyG, true
	case "KeyH":
		return events.VirtualKeyH, true
	case "KeyI":
		return events.VirtualKeyI, true
	case "KeyJ":
		return events.VirtualKeyJ, true
	case "KeyK":
		return events.VirtualKeyK, true
	case "KeyL":
		return events.VirtualKeyL, true
	case "KeyM":
		return events.VirtualKeyM, true
	case "KeyN":
		return events.VirtualKeyN, true
	case "KeyO":
		return events.VirtualKeyO, true
	case "KeyP":
		return events.VirtualKeyP, true
	case "KeyQ":
		return events.VirtualKeyQ, true
	case "KeyR":
		return events.VirtualKeyR, true
	case "KeyS":
		return events.VirtualKeyS, true
	case "KeyT":
		return events.VirtualKeyT, true
	case "KeyU":
		return events.VirtualKeyU, true
	case "KeyV":
		return events.VirtualKeyV, true
	case "KeyW":
		return events.VirtualKeyW, true
	case "KeyX":
		return events.VirtualKeyX, true
	case "KeyY":
		return events.VirtualKeyY, true
	case "KeyZ":
		return events.VirtualKeyZ, true
	case "OSLeft", "MetaLeft":
		return events.VirtualKeyLWin, true
	case "OSRight", "MetaRight":
		return events.VirtualKeyRWin, true
	case "ContextMenu":
		return events.VirtualKeyContextMenu, true
	case "Sleep":
		return events.VirtualKeySleep, true
	case "Numpad0":
		return events.VirtualKeyNumpad0, true
	case "Numpad1":
		return events.VirtualKeyNumpad1, true
	case "Numpad2":
		return events.VirtualKeyNumpad2, true
	case "Numpad3":
		return events.VirtualKeyNumpad3, true
	case "Numpad4":
		return events.VirtualKeyNumpad4, true
	case "Numpad5":
		return events.VirtualKeyNumpad5, true
	case "Numpad6":
		return events.VirtualKeyNumpad6, true
	case "Numpad7":
		return events.VirtualKeyNumpad7, true
	case "Numpad8":
		return events.VirtualKeyNumpad8, true
	case "Numpad9":
		return events.VirtualKeyNumpad9, true
	case "NumpadMultiply":
		return events.VirtualKeyMultiply, true
	case "NumpadAdd":
		return events.VirtualKeyAdd, true
	case "NumpadSubtract":
		return events.VirtualKeySubtract, true
	case "NumpadDecimal":
		return events.VirtualKeyDecimal, true
	case "NumpadDivide":
		return events.VirtualKeyDivide, true
	case "F1":
		return events.VirtualKeyF1, true
	case "F2":
		return events.VirtualKeyF2, true
	case "F3":
		return events.VirtualKeyF3, true
	case "F4":
		return events.VirtualKeyF4, true
	case "F5":
		return events.VirtualKeyF5, true
	case "F6":
		return events.VirtualKeyF6, true
	case "F7":
		return events.VirtualKeyF7, true
	case "F8":
		return events.VirtualKeyF8, true
	case "F9":
		return events.VirtualKeyF9, true
	case "F10":
		return events.VirtualKeyF10, true
	case "F11":
		return events.VirtualKeyF11, true
	case "F12":
		return events.VirtualKeyF12, true
	case "F13":
		return events.VirtualKeyF13, true
	case "F14":
		return events.VirtualKeyF14, true
	case "F15":
		return events.VirtualKeyF15, true
	case "F16":
		return events.VirtualKeyF16, true
	case "F17":
		return events.VirtualKeyF17, true
	case "F18":
		return events.VirtualKeyF18, true
	case "F19":
		return events.VirtualKeyF19, true
	case "F20":
		return events.VirtualKeyF20, true
	case "F21":
		return events.VirtualKeyF21, true
	case "F22":
		return events.VirtualKeyF22, true
	case "F23":
		return events.VirtualKeyF23, true
	case "F24":
		return events.VirtualKeyF24, true
	case "NumLock":
		return events.VirtualKeyNumLock, true
	case "ScrollLock":
		return events.VirtualKeyScrollLock, true
	case "Minus":
		return events.VirtualKeyHyphenMinus, true
	case "AudioVolumeMute", "VolumeMute":
		return events.VirtualKeyVolumeMute, true
	case "AudioVolumeDown", "VolumeDown":
		return events.VirtualKeyVolumeDown, true
	case "AudioVolumeUp", "VolumeUp":
		return events.VirtualKeyVolumeUp, true
	case "Comma":
		return events.VirtualKeyComma, true
	case "Period":
		return events.VirtualKeyPeriod, true
	case "Slash":
		return events.VirtualKeySlash, true
	case "Backquote":
		return events.VirtualKeyBackQuote, true
	case "BracketLeft":
		return events.VirtualKeyOpenBracket, true
	case "Backslash":
		return events.VirtualKeyBackSlash, true
	case "BracketRight":
		return events.VirtualKeyCloseBracket, true
	case "Quote":
		return events.VirtualKeyQuote, true
		// case "":
		// 	return events.VirtualKeyAltgr, true
	}

	return 0, false
}
