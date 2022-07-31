//go:build windows

package win32

import (
	"sync/atomic"
	"unicode/utf16"
	"unsafe"

	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/win32/procs"
)

func mapVK(vk uintptr, scancode uintptr, extended bool) events.VirtualKey {
	switch vk {
	case procs.VK_CANCEL:
		return events.VirtualKeyCancel
	case procs.VK_HELP:
		return events.VirtualKeyHelp
	case procs.VK_BACK:
		return events.VirtualKeyBackSpace
	case procs.VK_TAB:
		return events.VirtualKeyTab
	case procs.VK_CLEAR:
		return events.VirtualKeyClear
	case procs.VK_RETURN:
		return events.VirtualKeyReturn
	case procs.VK_PAUSE:
		return events.VirtualKeyPause
	case procs.VK_CAPITAL:
		return events.VirtualKeyCapsLock
	case procs.VK_KANA:
		return events.VirtualKeyKana
	case procs.VK_JUNJA:
		return events.VirtualKeyJunja
	case procs.VK_FINAL:
		return events.VirtualKeyFinal
	case procs.VK_KANJI:
		return events.VirtualKeyKanji
	case procs.VK_ESCAPE:
		return events.VirtualKeyEscape
	case procs.VK_CONVERT:
		return events.VirtualKeyConvert
	case procs.VK_NONCONVERT:
		return events.VirtualKeyNonconvert
	case procs.VK_ACCEPT:
		return events.VirtualKeyAccept
	case procs.VK_MODECHANGE:
		return events.VirtualKeyModechange
	case procs.VK_SPACE:
		return events.VirtualKeySpace
	case procs.VK_PRIOR:
		return events.VirtualKeyPageUp
	case procs.VK_NEXT:
		return events.VirtualKeyPageDown
	case procs.VK_END:
		return events.VirtualKeyEnd
	case procs.VK_HOME:
		return events.VirtualKeyHome
	case procs.VK_LEFT:
		return events.VirtualKeyLeft
	case procs.VK_UP:
		return events.VirtualKeyUp
	case procs.VK_RIGHT:
		return events.VirtualKeyRight
	case procs.VK_DOWN:
		return events.VirtualKeyDown
	case procs.VK_SELECT:
		return events.VirtualKeySelect
	case procs.VK_PRINT:
		return events.VirtualKeyPrint
	case procs.VK_EXECUTE:
		return events.VirtualKeyExecute
	case procs.VK_INSERT:
		return events.VirtualKeyInsert
	case procs.VK_DELETE:
		return events.VirtualKeyDelete
	case procs.VK_0:
		return events.VirtualKey0
	case procs.VK_1:
		return events.VirtualKey1
	case procs.VK_2:
		return events.VirtualKey2
	case procs.VK_3:
		return events.VirtualKey3
	case procs.VK_4:
		return events.VirtualKey4
	case procs.VK_5:
		return events.VirtualKey5
	case procs.VK_6:
		return events.VirtualKey6
	case procs.VK_7:
		return events.VirtualKey7
	case procs.VK_8:
		return events.VirtualKey8
	case procs.VK_9:
		return events.VirtualKey9
	case procs.VK_A:
		return events.VirtualKeyA
	case procs.VK_B:
		return events.VirtualKeyB
	case procs.VK_C:
		return events.VirtualKeyC
	case procs.VK_D:
		return events.VirtualKeyD
	case procs.VK_E:
		return events.VirtualKeyE
	case procs.VK_F:
		return events.VirtualKeyF
	case procs.VK_G:
		return events.VirtualKeyG
	case procs.VK_H:
		return events.VirtualKeyH
	case procs.VK_I:
		return events.VirtualKeyI
	case procs.VK_J:
		return events.VirtualKeyJ
	case procs.VK_K:
		return events.VirtualKeyK
	case procs.VK_L:
		return events.VirtualKeyL
	case procs.VK_M:
		return events.VirtualKeyM
	case procs.VK_N:
		return events.VirtualKeyN
	case procs.VK_O:
		return events.VirtualKeyO
	case procs.VK_P:
		return events.VirtualKeyP
	case procs.VK_Q:
		return events.VirtualKeyQ
	case procs.VK_R:
		return events.VirtualKeyR
	case procs.VK_S:
		return events.VirtualKeyS
	case procs.VK_T:
		return events.VirtualKeyT
	case procs.VK_U:
		return events.VirtualKeyU
	case procs.VK_V:
		return events.VirtualKeyV
	case procs.VK_W:
		return events.VirtualKeyW
	case procs.VK_X:
		return events.VirtualKeyX
	case procs.VK_Y:
		return events.VirtualKeyY
	case procs.VK_Z:
		return events.VirtualKeyZ
	case procs.VK_LWIN:
		return events.VirtualKeyLWin
	case procs.VK_RWIN:
		return events.VirtualKeyRWin
	case procs.VK_APPS:
		return events.VirtualKeyContextMenu
	case procs.VK_SLEEP:
		return events.VirtualKeySleep
	case procs.VK_NUMPAD0:
		return events.VirtualKeyNumpad0
	case procs.VK_NUMPAD1:
		return events.VirtualKeyNumpad1
	case procs.VK_NUMPAD2:
		return events.VirtualKeyNumpad2
	case procs.VK_NUMPAD3:
		return events.VirtualKeyNumpad3
	case procs.VK_NUMPAD4:
		return events.VirtualKeyNumpad4
	case procs.VK_NUMPAD5:
		return events.VirtualKeyNumpad5
	case procs.VK_NUMPAD6:
		return events.VirtualKeyNumpad6
	case procs.VK_NUMPAD7:
		return events.VirtualKeyNumpad7
	case procs.VK_NUMPAD8:
		return events.VirtualKeyNumpad8
	case procs.VK_NUMPAD9:
		return events.VirtualKeyNumpad9
	case procs.VK_MULTIPLY:
		return events.VirtualKeyMultiply
	case procs.VK_ADD:
		return events.VirtualKeyAdd
	case procs.VK_SUBTRACT:
		return events.VirtualKeySubtract
	case procs.VK_DECIMAL:
		return events.VirtualKeyDecimal
	case procs.VK_DIVIDE:
		return events.VirtualKeyDivide
	case procs.VK_F1:
		return events.VirtualKeyF1
	case procs.VK_F2:
		return events.VirtualKeyF2
	case procs.VK_F3:
		return events.VirtualKeyF3
	case procs.VK_F4:
		return events.VirtualKeyF4
	case procs.VK_F5:
		return events.VirtualKeyF5
	case procs.VK_F6:
		return events.VirtualKeyF6
	case procs.VK_F7:
		return events.VirtualKeyF7
	case procs.VK_F8:
		return events.VirtualKeyF8
	case procs.VK_F9:
		return events.VirtualKeyF9
	case procs.VK_F10:
		return events.VirtualKeyF10
	case procs.VK_F11:
		return events.VirtualKeyF11
	case procs.VK_F12:
		return events.VirtualKeyF12
	case procs.VK_F13:
		return events.VirtualKeyF13
	case procs.VK_F14:
		return events.VirtualKeyF14
	case procs.VK_F15:
		return events.VirtualKeyF15
	case procs.VK_F16:
		return events.VirtualKeyF16
	case procs.VK_F17:
		return events.VirtualKeyF17
	case procs.VK_F18:
		return events.VirtualKeyF18
	case procs.VK_F19:
		return events.VirtualKeyF19
	case procs.VK_F20:
		return events.VirtualKeyF20
	case procs.VK_F21:
		return events.VirtualKeyF21
	case procs.VK_F22:
		return events.VirtualKeyF22
	case procs.VK_F23:
		return events.VirtualKeyF23
	case procs.VK_F24:
		return events.VirtualKeyF24
	case procs.VK_NUMLOCK:
		return events.VirtualKeyNumLock
	case procs.VK_SCROLL:
		return events.VirtualKeyScrollLock
	case procs.VK_VOLUME_MUTE:
		return events.VirtualKeyVolumeMute
	case procs.VK_VOLUME_DOWN:
		return events.VirtualKeyVolumeDown
	case procs.VK_VOLUME_UP:
		return events.VirtualKeyVolumeUp

	case procs.VK_SHIFT:
		switch procs.MapVirtualKeyA(scancode, procs.MAPVK_VSC_TO_VK_EX) {
		case procs.VK_LSHIFT:
			return events.VirtualKeyLShift
		case procs.VK_RSHIFT:
			return events.VirtualKeyRShift
		}

	case procs.VK_CONTROL:
		if extended {
			return events.VirtualKeyRControl
		} else {
			return events.VirtualKeyLControl
		}

	case procs.VK_MENU:
		if extended {
			if layoutUsesAltgr() {
				return events.VirtualKeyAltgr
			} else {
				return events.VirtualKeyRAlt
			}
		} else {
			return events.VirtualKeyLAlt
		}

	default:
		switch rune(procs.MapVirtualKeyA(vk, procs.MAPVK_VK_TO_CHAR) & 0x7FFF) {
		case '-':
			return events.VirtualKeyHyphenMinus
		case ';':
			return events.VirtualKeySemicolon
		case '=':
			return events.VirtualKeyEquals
		case ',':
			return events.VirtualKeyComma
		case '.':
			return events.VirtualKeyPeriod
		case '/':
			return events.VirtualKeySlash
		case '`':
			return events.VirtualKeyBackQuote
		case '[':
			return events.VirtualKeyOpenBracket
		case ']':
			return events.VirtualKeyCloseBracket
		case '\\':
			return events.VirtualKeyBackSlash
		case '\'':
			return events.VirtualKeyQuote
		}
	}

	return events.VirtualKey(vk)
}

func getChar(keyboardState *[256]byte, vKey uintptr, hkl uintptr) rune {
	var unicodeBytes [5]uint16
	l := procs.ToUnicodeEx(
		vKey,
		0,
		uintptr(unsafe.Pointer(&keyboardState)),
		uintptr(unsafe.Pointer(&unicodeBytes)),
		uintptr(len(unicodeBytes)),
		0,
		hkl,
	)
	if l >= 1 {
		r := utf16.Decode(unicodeBytes[:l])
		if len(r) > 0 {
			return r[0]
		}
	}
	return -1
}

var currentLayout uintptr
var usesAltGr uint32

func layoutUsesAltgr() bool {
	hkl := procs.GetKeyboardLayout(0)
	oldHkl := atomic.SwapUintptr(&currentLayout, hkl)

	if hkl == oldHkl {
		return atomic.LoadUint32(&usesAltGr) != 0
	}

	var keyboardStateAltgr [256]byte
	var keyboardStateEmpty [256]byte

	keyboardStateAltgr[procs.VK_MENU] = 0x80
	keyboardStateAltgr[procs.VK_CONTROL] = 0x80

	for vKey := uintptr(0); vKey < 255; vKey++ {
		keyNoAltgr := getChar(&keyboardStateEmpty, vKey, hkl)
		keyAltgr := getChar(&keyboardStateAltgr, vKey, hkl)
		if keyAltgr != -1 && keyNoAltgr != -1 {
			if keyAltgr != keyNoAltgr {
				atomic.StoreUint32(&usesAltGr, 1)
				return true
			}
		}
	}

	atomic.StoreUint32(&usesAltGr, 0)
	return false
}
