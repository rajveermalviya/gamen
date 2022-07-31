package events

import (
	"strconv"
)

type VirtualKey uint32

const (
	VirtualKeyCancel VirtualKey = iota + 4096
	VirtualKeyHelp
	VirtualKeyBackSpace
	VirtualKeyTab
	VirtualKeyClear
	VirtualKeyReturn
	VirtualKeyLShift
	VirtualKeyRShift
	VirtualKeyLControl
	VirtualKeyRControl
	VirtualKeyLAlt
	VirtualKeyRAlt
	VirtualKeyPause
	VirtualKeyCapsLock
	VirtualKeyKana
	VirtualKeyJunja
	VirtualKeyFinal
	VirtualKeyKanji
	VirtualKeyEscape
	VirtualKeyConvert
	VirtualKeyNonconvert
	VirtualKeyAccept
	VirtualKeyModechange
	VirtualKeySpace
	VirtualKeyPageUp
	VirtualKeyPageDown
	VirtualKeyEnd
	VirtualKeyHome
	VirtualKeyLeft
	VirtualKeyUp
	VirtualKeyRight
	VirtualKeyDown
	VirtualKeySelect
	VirtualKeyPrint
	VirtualKeyExecute
	VirtualKeyInsert
	VirtualKeyDelete
	VirtualKey0
	VirtualKey1
	VirtualKey2
	VirtualKey3
	VirtualKey4
	VirtualKey5
	VirtualKey6
	VirtualKey7
	VirtualKey8
	VirtualKey9
	VirtualKeySemicolon
	VirtualKeyEquals
	VirtualKeyA
	VirtualKeyB
	VirtualKeyC
	VirtualKeyD
	VirtualKeyE
	VirtualKeyF
	VirtualKeyG
	VirtualKeyH
	VirtualKeyI
	VirtualKeyJ
	VirtualKeyK
	VirtualKeyL
	VirtualKeyM
	VirtualKeyN
	VirtualKeyO
	VirtualKeyP
	VirtualKeyQ
	VirtualKeyR
	VirtualKeyS
	VirtualKeyT
	VirtualKeyU
	VirtualKeyV
	VirtualKeyW
	VirtualKeyX
	VirtualKeyY
	VirtualKeyZ
	VirtualKeyLWin
	VirtualKeyRWin
	VirtualKeyContextMenu
	VirtualKeySleep
	VirtualKeyNumpad0
	VirtualKeyNumpad1
	VirtualKeyNumpad2
	VirtualKeyNumpad3
	VirtualKeyNumpad4
	VirtualKeyNumpad5
	VirtualKeyNumpad6
	VirtualKeyNumpad7
	VirtualKeyNumpad8
	VirtualKeyNumpad9
	VirtualKeyMultiply
	VirtualKeyAdd
	VirtualKeySubtract
	VirtualKeyDecimal
	VirtualKeyDivide
	VirtualKeyF1
	VirtualKeyF2
	VirtualKeyF3
	VirtualKeyF4
	VirtualKeyF5
	VirtualKeyF6
	VirtualKeyF7
	VirtualKeyF8
	VirtualKeyF9
	VirtualKeyF10
	VirtualKeyF11
	VirtualKeyF12
	VirtualKeyF13
	VirtualKeyF14
	VirtualKeyF15
	VirtualKeyF16
	VirtualKeyF17
	VirtualKeyF18
	VirtualKeyF19
	VirtualKeyF20
	VirtualKeyF21
	VirtualKeyF22
	VirtualKeyF23
	VirtualKeyF24
	VirtualKeyNumLock
	VirtualKeyScrollLock
	VirtualKeyHyphenMinus
	VirtualKeyVolumeMute
	VirtualKeyVolumeDown
	VirtualKeyVolumeUp
	VirtualKeyComma
	VirtualKeyPeriod
	VirtualKeySlash
	VirtualKeyBackQuote
	VirtualKeyOpenBracket
	VirtualKeyBackSlash
	VirtualKeyCloseBracket
	VirtualKeyQuote
	VirtualKeyAltgr
)

func (v VirtualKey) String() string {
	switch v {
	case VirtualKeyCancel:
		return "Cancel"
	case VirtualKeyHelp:
		return "Help"
	case VirtualKeyBackSpace:
		return "BackSpace"
	case VirtualKeyTab:
		return "Tab"
	case VirtualKeyClear:
		return "Clear"
	case VirtualKeyReturn:
		return "Return"
	case VirtualKeyLShift:
		return "LShift"
	case VirtualKeyRShift:
		return "RShift"
	case VirtualKeyLControl:
		return "LControl"
	case VirtualKeyRControl:
		return "RControl"
	case VirtualKeyLAlt:
		return "LAlt"
	case VirtualKeyRAlt:
		return "RAlt"
	case VirtualKeyPause:
		return "Pause"
	case VirtualKeyCapsLock:
		return "CapsLock"
	case VirtualKeyKana:
		return "Kana"
	case VirtualKeyJunja:
		return "Junja"
	case VirtualKeyFinal:
		return "Final"
	case VirtualKeyKanji:
		return "Kanji"
	case VirtualKeyEscape:
		return "Escape"
	case VirtualKeyConvert:
		return "Convert"
	case VirtualKeyNonconvert:
		return "Nonconvert"
	case VirtualKeyAccept:
		return "Accept"
	case VirtualKeyModechange:
		return "Modechange"
	case VirtualKeySpace:
		return "Space"
	case VirtualKeyPageUp:
		return "PageUp"
	case VirtualKeyPageDown:
		return "PageDown"
	case VirtualKeyEnd:
		return "End"
	case VirtualKeyHome:
		return "Home"
	case VirtualKeyLeft:
		return "Left"
	case VirtualKeyUp:
		return "Up"
	case VirtualKeyRight:
		return "Right"
	case VirtualKeyDown:
		return "Down"
	case VirtualKeySelect:
		return "Select"
	case VirtualKeyPrint:
		return "Print"
	case VirtualKeyExecute:
		return "Execute"
	case VirtualKeyInsert:
		return "Insert"
	case VirtualKeyDelete:
		return "Delete"
	case VirtualKey0:
		return "0"
	case VirtualKey1:
		return "1"
	case VirtualKey2:
		return "2"
	case VirtualKey3:
		return "3"
	case VirtualKey4:
		return "4"
	case VirtualKey5:
		return "5"
	case VirtualKey6:
		return "6"
	case VirtualKey7:
		return "7"
	case VirtualKey8:
		return "8"
	case VirtualKey9:
		return "9"
	case VirtualKeySemicolon:
		return "Semicolon"
	case VirtualKeyEquals:
		return "Equals"
	case VirtualKeyA:
		return "A"
	case VirtualKeyB:
		return "B"
	case VirtualKeyC:
		return "C"
	case VirtualKeyD:
		return "D"
	case VirtualKeyE:
		return "E"
	case VirtualKeyF:
		return "F"
	case VirtualKeyG:
		return "G"
	case VirtualKeyH:
		return "H"
	case VirtualKeyI:
		return "I"
	case VirtualKeyJ:
		return "J"
	case VirtualKeyK:
		return "K"
	case VirtualKeyL:
		return "L"
	case VirtualKeyM:
		return "M"
	case VirtualKeyN:
		return "N"
	case VirtualKeyO:
		return "O"
	case VirtualKeyP:
		return "P"
	case VirtualKeyQ:
		return "Q"
	case VirtualKeyR:
		return "R"
	case VirtualKeyS:
		return "S"
	case VirtualKeyT:
		return "T"
	case VirtualKeyU:
		return "U"
	case VirtualKeyV:
		return "V"
	case VirtualKeyW:
		return "W"
	case VirtualKeyX:
		return "X"
	case VirtualKeyY:
		return "Y"
	case VirtualKeyZ:
		return "Z"
	case VirtualKeyLWin:
		return "LWin"
	case VirtualKeyRWin:
		return "RWin"
	case VirtualKeyContextMenu:
		return "ContextMenu"
	case VirtualKeySleep:
		return "Sleep"
	case VirtualKeyNumpad0:
		return "Numpad0"
	case VirtualKeyNumpad1:
		return "Numpad1"
	case VirtualKeyNumpad2:
		return "Numpad2"
	case VirtualKeyNumpad3:
		return "Numpad3"
	case VirtualKeyNumpad4:
		return "Numpad4"
	case VirtualKeyNumpad5:
		return "Numpad5"
	case VirtualKeyNumpad6:
		return "Numpad6"
	case VirtualKeyNumpad7:
		return "Numpad7"
	case VirtualKeyNumpad8:
		return "Numpad8"
	case VirtualKeyNumpad9:
		return "Numpad9"
	case VirtualKeyMultiply:
		return "Multiply"
	case VirtualKeyAdd:
		return "Add"
	case VirtualKeySubtract:
		return "Subtract"
	case VirtualKeyDecimal:
		return "Decimal"
	case VirtualKeyDivide:
		return "Divide"
	case VirtualKeyF1:
		return "F1"
	case VirtualKeyF2:
		return "F2"
	case VirtualKeyF3:
		return "F3"
	case VirtualKeyF4:
		return "F4"
	case VirtualKeyF5:
		return "F5"
	case VirtualKeyF6:
		return "F6"
	case VirtualKeyF7:
		return "F7"
	case VirtualKeyF8:
		return "F8"
	case VirtualKeyF9:
		return "F9"
	case VirtualKeyF10:
		return "F10"
	case VirtualKeyF11:
		return "F11"
	case VirtualKeyF12:
		return "F12"
	case VirtualKeyF13:
		return "F13"
	case VirtualKeyF14:
		return "F14"
	case VirtualKeyF15:
		return "F15"
	case VirtualKeyF16:
		return "F16"
	case VirtualKeyF17:
		return "F17"
	case VirtualKeyF18:
		return "F18"
	case VirtualKeyF19:
		return "F19"
	case VirtualKeyF20:
		return "F20"
	case VirtualKeyF21:
		return "F21"
	case VirtualKeyF22:
		return "F22"
	case VirtualKeyF23:
		return "F23"
	case VirtualKeyF24:
		return "F24"
	case VirtualKeyNumLock:
		return "NumLock"
	case VirtualKeyScrollLock:
		return "ScrollLock"
	case VirtualKeyHyphenMinus:
		return "HyphenMinus"
	case VirtualKeyVolumeMute:
		return "VolumeMute"
	case VirtualKeyVolumeDown:
		return "VolumeDown"
	case VirtualKeyVolumeUp:
		return "VolumeUp"
	case VirtualKeyComma:
		return "Comma"
	case VirtualKeyPeriod:
		return "Period"
	case VirtualKeySlash:
		return "Slash"
	case VirtualKeyBackQuote:
		return "BackQuote"
	case VirtualKeyOpenBracket:
		return "OpenBracket"
	case VirtualKeyBackSlash:
		return "BackSlash"
	case VirtualKeyCloseBracket:
		return "CloseBracket"
	case VirtualKeyQuote:
		return "Quote"
	case VirtualKeyAltgr:
		return "Altgr"
	}

	return strconv.FormatUint(uint64(v), 16)
}
