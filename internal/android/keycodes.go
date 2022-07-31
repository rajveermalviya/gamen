//go:build android

package android

/*

#include <game-activity/native_app_glue/android_native_app_glue.h>

*/
import "C"
import "github.com/rajveermalviya/gamen/events"

// https://android.googlesource.com/platform/frameworks/native/+/master/include/android/keycodes.h

func mapKeycode(k C.int32_t) events.VirtualKey {
	switch k {
	// case 0:
	// 	return events.VirtualKeyCancel
	case C.AKEYCODE_HELP:
		return events.VirtualKeyHelp
	case C.AKEYCODE_DEL:
		return events.VirtualKeyBackSpace
	case C.AKEYCODE_TAB:
		return events.VirtualKeyTab
	case C.AKEYCODE_CLEAR:
		return events.VirtualKeyClear
	case C.AKEYCODE_ENTER:
		return events.VirtualKeyReturn
	case C.AKEYCODE_SHIFT_LEFT:
		return events.VirtualKeyLShift
	case C.AKEYCODE_SHIFT_RIGHT:
		return events.VirtualKeyRShift
	case C.AKEYCODE_CTRL_LEFT:
		return events.VirtualKeyLControl
	case C.AKEYCODE_CTRL_RIGHT:
		return events.VirtualKeyRControl
	case C.AKEYCODE_ALT_LEFT:
		return events.VirtualKeyLAlt
	case C.AKEYCODE_ALT_RIGHT:
		return events.VirtualKeyRAlt
	case C.AKEYCODE_BREAK:
		return events.VirtualKeyPause
	case C.AKEYCODE_CAPS_LOCK:
		return events.VirtualKeyCapsLock
	case C.AKEYCODE_KANA:
		return events.VirtualKeyKana
	// case 0:
	// 	return events.VirtualKeyJunja
	// case 0:
	// 	return events.VirtualKeyFinal
	// case 0:
	// 	return events.VirtualKeyKanji
	case C.AKEYCODE_ESCAPE:
		return events.VirtualKeyEscape
	case C.AKEYCODE_HENKAN:
		return events.VirtualKeyConvert
	case C.AKEYCODE_MUHENKAN:
		return events.VirtualKeyNonconvert
	// case 0:
	// 	return events.VirtualKeyAccept
	// case 0:
	// 	return events.VirtualKeyModechange
	case C.AKEYCODE_SPACE:
		return events.VirtualKeySpace
	case C.AKEYCODE_PAGE_UP:
		return events.VirtualKeyPageUp
	case C.AKEYCODE_PAGE_DOWN:
		return events.VirtualKeyPageDown
	case C.AKEYCODE_MOVE_END:
		return events.VirtualKeyEnd
	case C.AKEYCODE_MOVE_HOME:
		return events.VirtualKeyHome
	case C.AKEYCODE_DPAD_LEFT:
		return events.VirtualKeyLeft
	case C.AKEYCODE_DPAD_UP:
		return events.VirtualKeyUp
	case C.AKEYCODE_DPAD_RIGHT:
		return events.VirtualKeyRight
	case C.AKEYCODE_DPAD_DOWN:
		return events.VirtualKeyDown
	// case 0:
	// 	return events.VirtualKeySelect
	case C.AKEYCODE_SYSRQ:
		return events.VirtualKeyPrint
	// case 0:
	// 	return events.VirtualKeyExecute
	case C.AKEYCODE_INSERT:
		return events.VirtualKeyInsert
	case C.AKEYCODE_FORWARD_DEL:
		return events.VirtualKeyDelete
	case C.AKEYCODE_0:
		return events.VirtualKey0
	case C.AKEYCODE_1:
		return events.VirtualKey1
	case C.AKEYCODE_2:
		return events.VirtualKey2
	case C.AKEYCODE_3:
		return events.VirtualKey3
	case C.AKEYCODE_4:
		return events.VirtualKey4
	case C.AKEYCODE_5:
		return events.VirtualKey5
	case C.AKEYCODE_6:
		return events.VirtualKey6
	case C.AKEYCODE_7:
		return events.VirtualKey7
	case C.AKEYCODE_8:
		return events.VirtualKey8
	case C.AKEYCODE_9:
		return events.VirtualKey9
	case C.AKEYCODE_SEMICOLON:
		return events.VirtualKeySemicolon
	case C.AKEYCODE_EQUALS:
		return events.VirtualKeyEquals
	case C.AKEYCODE_A:
		return events.VirtualKeyA
	case C.AKEYCODE_B:
		return events.VirtualKeyB
	case C.AKEYCODE_C:
		return events.VirtualKeyC
	case C.AKEYCODE_D:
		return events.VirtualKeyD
	case C.AKEYCODE_E:
		return events.VirtualKeyE
	case C.AKEYCODE_F:
		return events.VirtualKeyF
	case C.AKEYCODE_G:
		return events.VirtualKeyG
	case C.AKEYCODE_H:
		return events.VirtualKeyH
	case C.AKEYCODE_I:
		return events.VirtualKeyI
	case C.AKEYCODE_J:
		return events.VirtualKeyJ
	case C.AKEYCODE_K:
		return events.VirtualKeyK
	case C.AKEYCODE_L:
		return events.VirtualKeyL
	case C.AKEYCODE_M:
		return events.VirtualKeyM
	case C.AKEYCODE_N:
		return events.VirtualKeyN
	case C.AKEYCODE_O:
		return events.VirtualKeyO
	case C.AKEYCODE_P:
		return events.VirtualKeyP
	case C.AKEYCODE_Q:
		return events.VirtualKeyQ
	case C.AKEYCODE_R:
		return events.VirtualKeyR
	case C.AKEYCODE_S:
		return events.VirtualKeyS
	case C.AKEYCODE_T:
		return events.VirtualKeyT
	case C.AKEYCODE_U:
		return events.VirtualKeyU
	case C.AKEYCODE_V:
		return events.VirtualKeyV
	case C.AKEYCODE_W:
		return events.VirtualKeyW
	case C.AKEYCODE_X:
		return events.VirtualKeyX
	case C.AKEYCODE_Y:
		return events.VirtualKeyY
	case C.AKEYCODE_Z:
		return events.VirtualKeyZ
	case C.AKEYCODE_META_LEFT:
		return events.VirtualKeyLWin
	case C.AKEYCODE_META_RIGHT:
		return events.VirtualKeyRWin
	case 0:
		return events.VirtualKeyContextMenu
	case C.AKEYCODE_SLEEP:
		return events.VirtualKeySleep
	case C.AKEYCODE_NUMPAD_0:
		return events.VirtualKeyNumpad0
	case C.AKEYCODE_NUMPAD_1:
		return events.VirtualKeyNumpad1
	case C.AKEYCODE_NUMPAD_2:
		return events.VirtualKeyNumpad2
	case C.AKEYCODE_NUMPAD_3:
		return events.VirtualKeyNumpad3
	case C.AKEYCODE_NUMPAD_4:
		return events.VirtualKeyNumpad4
	case C.AKEYCODE_NUMPAD_5:
		return events.VirtualKeyNumpad5
	case C.AKEYCODE_NUMPAD_6:
		return events.VirtualKeyNumpad6
	case C.AKEYCODE_NUMPAD_7:
		return events.VirtualKeyNumpad7
	case C.AKEYCODE_NUMPAD_8:
		return events.VirtualKeyNumpad8
	case C.AKEYCODE_NUMPAD_9:
		return events.VirtualKeyNumpad9
	case C.AKEYCODE_NUMPAD_MULTIPLY:
		return events.VirtualKeyMultiply
	case C.AKEYCODE_NUMPAD_ADD:
		return events.VirtualKeyAdd
	case C.AKEYCODE_NUMPAD_SUBTRACT:
		return events.VirtualKeySubtract
	case C.AKEYCODE_NUMPAD_DOT, C.AKEYCODE_NUMPAD_COMMA:
		return events.VirtualKeyDecimal
	case C.AKEYCODE_NUMPAD_DIVIDE:
		return events.VirtualKeyDivide
	case C.AKEYCODE_F1:
		return events.VirtualKeyF1
	case C.AKEYCODE_F2:
		return events.VirtualKeyF2
	case C.AKEYCODE_F3:
		return events.VirtualKeyF3
	case C.AKEYCODE_F4:
		return events.VirtualKeyF4
	case C.AKEYCODE_F5:
		return events.VirtualKeyF5
	case C.AKEYCODE_F6:
		return events.VirtualKeyF6
	case C.AKEYCODE_F7:
		return events.VirtualKeyF7
	case C.AKEYCODE_F8:
		return events.VirtualKeyF8
	case C.AKEYCODE_F9:
		return events.VirtualKeyF9
	case C.AKEYCODE_F10:
		return events.VirtualKeyF10
	case C.AKEYCODE_F11:
		return events.VirtualKeyF11
	case C.AKEYCODE_F12:
		return events.VirtualKeyF12
	// case 0:
	// 	return events.VirtualKeyF13
	// case 0:
	// 	return events.VirtualKeyF14
	// case 0:
	// 	return events.VirtualKeyF15
	// case 0:
	// 	return events.VirtualKeyF16
	// case 0:
	// 	return events.VirtualKeyF17
	// case 0:
	// 	return events.VirtualKeyF18
	// case 0:
	// 	return events.VirtualKeyF19
	// case 0:
	// 	return events.VirtualKeyF20
	// case 0:
	// 	return events.VirtualKeyF21
	// case 0:
	// 	return events.VirtualKeyF22
	// case 0:
	// 	return events.VirtualKeyF23
	// case 0:
	// 	return events.VirtualKeyF24
	case C.AKEYCODE_NUM_LOCK:
		return events.VirtualKeyNumLock
	case C.AKEYCODE_SCROLL_LOCK:
		return events.VirtualKeyScrollLock
	case C.AKEYCODE_MINUS:
		return events.VirtualKeyHyphenMinus
	case C.AKEYCODE_VOLUME_MUTE:
		return events.VirtualKeyVolumeMute
	case C.AKEYCODE_VOLUME_DOWN:
		return events.VirtualKeyVolumeDown
	case C.AKEYCODE_VOLUME_UP:
		return events.VirtualKeyVolumeUp
	case C.AKEYCODE_COMMA:
		return events.VirtualKeyComma
	case C.AKEYCODE_PERIOD:
		return events.VirtualKeyPeriod
	case C.AKEYCODE_SLASH:
		return events.VirtualKeySlash
	case C.AKEYCODE_GRAVE:
		return events.VirtualKeyBackQuote
	case C.AKEYCODE_LEFT_BRACKET:
		return events.VirtualKeyOpenBracket
	case C.AKEYCODE_BACKSLASH:
		return events.VirtualKeyBackSlash
	case C.AKEYCODE_RIGHT_BRACKET:
		return events.VirtualKeyCloseBracket
	case C.AKEYCODE_APOSTROPHE:
		return events.VirtualKeyQuote
		// case 0:
		// 	return events.VirtualKeyAltgr
	}

	return events.VirtualKey(k)
}
