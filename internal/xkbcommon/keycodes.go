//go:build linux && !android

package xkbcommon

import "github.com/rajveermalviya/gamen/events"

/*

#include <xkbcommon/xkbcommon.h>

*/
import "C"

func KeySymToVirtualKey(sym KeySym) events.VirtualKey {
	switch sym {
	case C.XKB_KEY_Cancel:
		return events.VirtualKeyCancel
	case C.XKB_KEY_BackSpace:
		return events.VirtualKeyBackSpace
	case C.XKB_KEY_Tab,
		C.XKB_KEY_ISO_Left_Tab:
		return events.VirtualKeyTab
	case C.XKB_KEY_Clear, C.XKB_KEY_KP_Begin:
		return events.VirtualKeyClear
	case C.XKB_KEY_Return, C.XKB_KEY_KP_Enter:
		return events.VirtualKeyReturn
	case C.XKB_KEY_Shift_L:
		return events.VirtualKeyLShift
	case C.XKB_KEY_Shift_R:
		return events.VirtualKeyRShift
	case C.XKB_KEY_Control_L:
		return events.VirtualKeyLControl
	case C.XKB_KEY_Control_R:
		return events.VirtualKeyRControl
	case C.XKB_KEY_Alt_L:
		return events.VirtualKeyLAlt
	case C.XKB_KEY_Alt_R:
		return events.VirtualKeyRAlt
	case C.XKB_KEY_Super_L,
		C.XKB_KEY_Hyper_L:
		return events.VirtualKeyLWin
	case C.XKB_KEY_Super_R,
		C.XKB_KEY_Hyper_R:
		return events.VirtualKeyRWin
	case C.XKB_KEY_ISO_Level3_Shift,
		C.XKB_KEY_ISO_Level5_Shift,
		C.XKB_KEY_ISO_Group_Shift:
		return events.VirtualKeyAltgr
	case C.XKB_KEY_Pause:
		return events.VirtualKeyPause
	case C.XKB_KEY_Caps_Lock:
		return events.VirtualKeyCapsLock
	case C.XKB_KEY_Kana_Lock,
		C.XKB_KEY_Kana_Shift:
		return events.VirtualKeyKana
	// case 0:
	// 	return events.VirtualKeyJunja
	// case 0:
	// 	return events.VirtualKeyFinal
	case C.XKB_KEY_Kanji:
		return events.VirtualKeyKanji
	case C.XKB_KEY_Escape:
		return events.VirtualKeyEscape
	case C.XKB_KEY_Henkan:
		return events.VirtualKeyConvert
	case C.XKB_KEY_Muhenkan:
		return events.VirtualKeyNonconvert
	// case 0:
	// 	return events.VirtualKeyAccept
	// case 0:
	// 	return events.VirtualKeyModechange
	case C.XKB_KEY_Page_Up, C.XKB_KEY_KP_Page_Up:
		return events.VirtualKeyPageUp
	case C.XKB_KEY_Page_Down, C.XKB_KEY_KP_Page_Down:
		return events.VirtualKeyPageDown
	case C.XKB_KEY_End, C.XKB_KEY_KP_End:
		return events.VirtualKeyEnd
	case C.XKB_KEY_Home, C.XKB_KEY_KP_Home:
		return events.VirtualKeyHome
	case C.XKB_KEY_Left, C.XKB_KEY_KP_Left:
		return events.VirtualKeyLeft
	case C.XKB_KEY_Up, C.XKB_KEY_KP_Up:
		return events.VirtualKeyUp
	case C.XKB_KEY_Right, C.XKB_KEY_KP_Right:
		return events.VirtualKeyRight
	case C.XKB_KEY_Down, C.XKB_KEY_KP_Down:
		return events.VirtualKeyDown

	case C.XKB_KEY_Select:
		return events.VirtualKeySelect
	case C.XKB_KEY_Print:
		return events.VirtualKeyPrint
	case C.XKB_KEY_Execute:
		return events.VirtualKeyExecute
	case C.XKB_KEY_Insert, C.XKB_KEY_KP_Insert:
		return events.VirtualKeyInsert
	case C.XKB_KEY_Delete, C.XKB_KEY_KP_Delete:
		return events.VirtualKeyDelete
	case C.XKB_KEY_Help:
		return events.VirtualKeyHelp
	case C.XKB_KEY_Num_Lock:
		return events.VirtualKeyNumLock
	case C.XKB_KEY_Scroll_Lock:
		return events.VirtualKeyScrollLock
	case C.XKB_KEY_F1:
		return events.VirtualKeyF1
	case C.XKB_KEY_F2:
		return events.VirtualKeyF2
	case C.XKB_KEY_F3:
		return events.VirtualKeyF3
	case C.XKB_KEY_F4:
		return events.VirtualKeyF4
	case C.XKB_KEY_F5:
		return events.VirtualKeyF5
	case C.XKB_KEY_F6:
		return events.VirtualKeyF6
	case C.XKB_KEY_F7:
		return events.VirtualKeyF7
	case C.XKB_KEY_F8:
		return events.VirtualKeyF8
	case C.XKB_KEY_F9:
		return events.VirtualKeyF9
	case C.XKB_KEY_F10:
		return events.VirtualKeyF10
	case C.XKB_KEY_F11:
		return events.VirtualKeyF11
	case C.XKB_KEY_F12:
		return events.VirtualKeyF12
	case C.XKB_KEY_F13:
		return events.VirtualKeyF13
	case C.XKB_KEY_F14:
		return events.VirtualKeyF14
	case C.XKB_KEY_F15:
		return events.VirtualKeyF15
	case C.XKB_KEY_F16:
		return events.VirtualKeyF16
	case C.XKB_KEY_F17:
		return events.VirtualKeyF17
	case C.XKB_KEY_F18:
		return events.VirtualKeyF18
	case C.XKB_KEY_F19:
		return events.VirtualKeyF19
	case C.XKB_KEY_F20:
		return events.VirtualKeyF20
	case C.XKB_KEY_F21:
		return events.VirtualKeyF21
	case C.XKB_KEY_F22:
		return events.VirtualKeyF22
	case C.XKB_KEY_F23:
		return events.VirtualKeyF23
	case C.XKB_KEY_F24:
		return events.VirtualKeyF24

	case C.XKB_KEY_Menu:
		return events.VirtualKeyContextMenu
	// case 0:
	// 	return events.VirtualKeySleep

	case C.XKB_KEY_space, C.XKB_KEY_KP_Space:
		return events.VirtualKeySpace

	case C.XKB_KEY_0:
		return events.VirtualKey0
	case C.XKB_KEY_1:
		return events.VirtualKey1
	case C.XKB_KEY_2:
		return events.VirtualKey2
	case C.XKB_KEY_3:
		return events.VirtualKey3
	case C.XKB_KEY_4:
		return events.VirtualKey4
	case C.XKB_KEY_5:
		return events.VirtualKey5
	case C.XKB_KEY_6:
		return events.VirtualKey6
	case C.XKB_KEY_7:
		return events.VirtualKey7
	case C.XKB_KEY_8:
		return events.VirtualKey8
	case C.XKB_KEY_9:
		return events.VirtualKey9

	case C.XKB_KEY_semicolon:
		return events.VirtualKeySemicolon
	case C.XKB_KEY_equal:
		return events.VirtualKeyEquals
	case C.XKB_KEY_A, C.XKB_KEY_a:
		return events.VirtualKeyA
	case C.XKB_KEY_B, C.XKB_KEY_b:
		return events.VirtualKeyB
	case C.XKB_KEY_C, C.XKB_KEY_c:
		return events.VirtualKeyC
	case C.XKB_KEY_D, C.XKB_KEY_d:
		return events.VirtualKeyD
	case C.XKB_KEY_E, C.XKB_KEY_e:
		return events.VirtualKeyE
	case C.XKB_KEY_F, C.XKB_KEY_f:
		return events.VirtualKeyF
	case C.XKB_KEY_G, C.XKB_KEY_g:
		return events.VirtualKeyG
	case C.XKB_KEY_H, C.XKB_KEY_h:
		return events.VirtualKeyH
	case C.XKB_KEY_I, C.XKB_KEY_i:
		return events.VirtualKeyI
	case C.XKB_KEY_J, C.XKB_KEY_j:
		return events.VirtualKeyJ
	case C.XKB_KEY_K, C.XKB_KEY_k:
		return events.VirtualKeyK
	case C.XKB_KEY_L, C.XKB_KEY_l:
		return events.VirtualKeyL
	case C.XKB_KEY_M, C.XKB_KEY_m:
		return events.VirtualKeyM
	case C.XKB_KEY_N, C.XKB_KEY_n:
		return events.VirtualKeyN
	case C.XKB_KEY_O, C.XKB_KEY_o:
		return events.VirtualKeyO
	case C.XKB_KEY_P, C.XKB_KEY_p:
		return events.VirtualKeyP
	case C.XKB_KEY_Q, C.XKB_KEY_q:
		return events.VirtualKeyQ
	case C.XKB_KEY_R, C.XKB_KEY_r:
		return events.VirtualKeyR
	case C.XKB_KEY_S, C.XKB_KEY_s:
		return events.VirtualKeyS
	case C.XKB_KEY_T, C.XKB_KEY_t:
		return events.VirtualKeyT
	case C.XKB_KEY_U, C.XKB_KEY_u:
		return events.VirtualKeyU
	case C.XKB_KEY_V, C.XKB_KEY_v:
		return events.VirtualKeyV
	case C.XKB_KEY_W, C.XKB_KEY_w:
		return events.VirtualKeyW
	case C.XKB_KEY_X, C.XKB_KEY_x:
		return events.VirtualKeyX
	case C.XKB_KEY_Y, C.XKB_KEY_y:
		return events.VirtualKeyY
	case C.XKB_KEY_Z, C.XKB_KEY_z:
		return events.VirtualKeyZ

	case C.XKB_KEY_KP_0:
		return events.VirtualKeyNumpad0
	case C.XKB_KEY_KP_1:
		return events.VirtualKeyNumpad1
	case C.XKB_KEY_KP_2:
		return events.VirtualKeyNumpad2
	case C.XKB_KEY_KP_3:
		return events.VirtualKeyNumpad3
	case C.XKB_KEY_KP_4:
		return events.VirtualKeyNumpad4
	case C.XKB_KEY_KP_5:
		return events.VirtualKeyNumpad5
	case C.XKB_KEY_KP_6:
		return events.VirtualKeyNumpad6
	case C.XKB_KEY_KP_7:
		return events.VirtualKeyNumpad7
	case C.XKB_KEY_KP_8:
		return events.VirtualKeyNumpad8
	case C.XKB_KEY_KP_9:
		return events.VirtualKeyNumpad9

	case C.XKB_KEY_KP_Multiply:
		return events.VirtualKeyMultiply
	case C.XKB_KEY_KP_Add:
		return events.VirtualKeyAdd
	case C.XKB_KEY_KP_Subtract:
		return events.VirtualKeySubtract
	case C.XKB_KEY_KP_Decimal:
		return events.VirtualKeyDecimal
	case C.XKB_KEY_KP_Divide:
		return events.VirtualKeyDivide

	case C.XKB_KEY_minus:
		return events.VirtualKeyHyphenMinus
	case C.XKB_KEY_XF86AudioMute:
		return events.VirtualKeyVolumeMute
	case C.XKB_KEY_XF86AudioLowerVolume:
		return events.VirtualKeyVolumeDown
	case C.XKB_KEY_XF86AudioRaiseVolume:
		return events.VirtualKeyVolumeUp
	case C.XKB_KEY_comma:
		return events.VirtualKeyComma
	case C.XKB_KEY_period:
		return events.VirtualKeyPeriod
	case C.XKB_KEY_slash:
		return events.VirtualKeySlash
	case C.XKB_KEY_grave:
		return events.VirtualKeyBackQuote
	case C.XKB_KEY_bracketleft:
		return events.VirtualKeyOpenBracket
	case C.XKB_KEY_backslash:
		return events.VirtualKeyBackSlash
	case C.XKB_KEY_bracketright:
		return events.VirtualKeyCloseBracket
	case C.XKB_KEY_apostrophe:
		return events.VirtualKeyQuote
	}

	return events.VirtualKey(sym)
}
