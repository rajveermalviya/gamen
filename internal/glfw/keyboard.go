package glfw

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/gamen/events"
)

func glfwKeyToVirtualKey(key glfw.Key) events.VirtualKey {
	var scanCode events.VirtualKey
	switch key {
	case glfw.KeySpace:
		scanCode = events.VirtualKeySpace
	case glfw.KeyApostrophe:
		scanCode = events.VirtualKeyBackQuote
	case glfw.KeyComma:
		scanCode = events.VirtualKeyComma
	case glfw.KeyMinus:
		scanCode = events.VirtualKeyHyphenMinus
	case glfw.KeyPeriod:
		scanCode = events.VirtualKeyPeriod
	case glfw.KeySlash:
		scanCode = events.VirtualKeySlash
	case glfw.Key0:
		scanCode = events.VirtualKey0
	case glfw.Key1:
		scanCode = events.VirtualKey1
	case glfw.Key2:
		scanCode = events.VirtualKey2
	case glfw.Key3:
		scanCode = events.VirtualKey3
	case glfw.Key4:
		scanCode = events.VirtualKey4
	case glfw.Key5:
		scanCode = events.VirtualKey5
	case glfw.Key6:
		scanCode = events.VirtualKey6
	case glfw.Key7:
		scanCode = events.VirtualKey7
	case glfw.Key8:
		scanCode = events.VirtualKey8
	case glfw.Key9:
		scanCode = events.VirtualKey9
	case glfw.KeySemicolon:
		scanCode = events.VirtualKeySemicolon
	case glfw.KeyEqual:
		scanCode = events.VirtualKeyEquals
	case glfw.KeyA:
		scanCode = events.VirtualKeyA
	case glfw.KeyB:
		scanCode = events.VirtualKeyB
	case glfw.KeyC:
		scanCode = events.VirtualKeyC
	case glfw.KeyD:
		scanCode = events.VirtualKeyD
	case glfw.KeyE:
		scanCode = events.VirtualKeyE
	case glfw.KeyF:
		scanCode = events.VirtualKeyF
	case glfw.KeyG:
		scanCode = events.VirtualKeyG
	case glfw.KeyH:
		scanCode = events.VirtualKeyH
	case glfw.KeyI:
		scanCode = events.VirtualKeyI
	case glfw.KeyJ:
		scanCode = events.VirtualKeyJ
	case glfw.KeyK:
		scanCode = events.VirtualKeyK
	case glfw.KeyL:
		scanCode = events.VirtualKeyL
	case glfw.KeyM:
		scanCode = events.VirtualKeyM
	case glfw.KeyN:
		scanCode = events.VirtualKeyN
	case glfw.KeyO:
		scanCode = events.VirtualKeyO
	case glfw.KeyP:
		scanCode = events.VirtualKeyP
	case glfw.KeyQ:
		scanCode = events.VirtualKeyQ
	case glfw.KeyR:
		scanCode = events.VirtualKeyR
	case glfw.KeyS:
		scanCode = events.VirtualKeyS
	case glfw.KeyT:
		scanCode = events.VirtualKeyT
	case glfw.KeyU:
		scanCode = events.VirtualKeyU
	case glfw.KeyV:
		scanCode = events.VirtualKeyV
	case glfw.KeyW:
		scanCode = events.VirtualKeyW
	case glfw.KeyX:
		scanCode = events.VirtualKeyX
	case glfw.KeyY:
		scanCode = events.VirtualKeyY
	case glfw.KeyZ:
		scanCode = events.VirtualKeyZ
	case glfw.KeyLeftBracket:
		scanCode = events.VirtualKeyOpenBracket
	case glfw.KeyBackslash:
		scanCode = events.VirtualKeyBackSlash
	case glfw.KeyRightBracket:
		scanCode = events.VirtualKeyCloseBracket
	case glfw.KeyGraveAccent:
		scanCode = events.VirtualKeyBackQuote
	case glfw.KeyEscape:
		scanCode = events.VirtualKeyEscape
	case glfw.KeyEnter:
		scanCode = events.VirtualKeyReturn
	case glfw.KeyTab:
		scanCode = events.VirtualKeyTab
	case glfw.KeyBackspace:
		scanCode = events.VirtualKeyBackSpace
	case glfw.KeyInsert:
		scanCode = events.VirtualKeyInsert
	case glfw.KeyDelete:
		scanCode = events.VirtualKeyDelete
	case glfw.KeyRight:
		scanCode = events.VirtualKeyRight
	case glfw.KeyLeft:
		scanCode = events.VirtualKeyLeft
	case glfw.KeyDown:
		scanCode = events.VirtualKeyDown
	case glfw.KeyUp:
		scanCode = events.VirtualKeyUp
	case glfw.KeyPageUp:
		scanCode = events.VirtualKeyPageUp
	case glfw.KeyPageDown:
		scanCode = events.VirtualKeyPageDown
	case glfw.KeyHome:
		scanCode = events.VirtualKeyHome
	case glfw.KeyEnd:
		scanCode = events.VirtualKeyEnd
	case glfw.KeyCapsLock:
		scanCode = events.VirtualKeyCapsLock
	case glfw.KeyScrollLock:
		scanCode = events.VirtualKeyScrollLock
	case glfw.KeyNumLock:
		scanCode = events.VirtualKeyNumLock
	case glfw.KeyPrintScreen:
		scanCode = events.VirtualKeyPrint
	case glfw.KeyPause:
		scanCode = events.VirtualKeyPause
	case glfw.KeyF1:
		scanCode = events.VirtualKeyF1
	case glfw.KeyF2:
		scanCode = events.VirtualKeyF2
	case glfw.KeyF3:
		scanCode = events.VirtualKeyF3
	case glfw.KeyF4:
		scanCode = events.VirtualKeyF4
	case glfw.KeyF5:
		scanCode = events.VirtualKeyF5
	case glfw.KeyF6:
		scanCode = events.VirtualKeyF6
	case glfw.KeyF7:
		scanCode = events.VirtualKeyF7
	case glfw.KeyF8:
		scanCode = events.VirtualKeyF8
	case glfw.KeyF9:
		scanCode = events.VirtualKeyF9
	case glfw.KeyF10:
		scanCode = events.VirtualKeyF10
	case glfw.KeyF11:
		scanCode = events.VirtualKeyF11
	case glfw.KeyF12:
		scanCode = events.VirtualKeyF12
	case glfw.KeyF13:
		scanCode = events.VirtualKeyF13
	case glfw.KeyF14:
		scanCode = events.VirtualKeyF14
	case glfw.KeyF15:
		scanCode = events.VirtualKeyF15
	case glfw.KeyF16:
		scanCode = events.VirtualKeyF16
	case glfw.KeyF17:
		scanCode = events.VirtualKeyF17
	case glfw.KeyF18:
		scanCode = events.VirtualKeyF18
	case glfw.KeyF19:
		scanCode = events.VirtualKeyF19
	case glfw.KeyF20:
		scanCode = events.VirtualKeyF20
	case glfw.KeyF21:
		scanCode = events.VirtualKeyF21
	case glfw.KeyF22:
		scanCode = events.VirtualKeyF22
	case glfw.KeyF23:
		scanCode = events.VirtualKeyF23
	case glfw.KeyF24:
		scanCode = events.VirtualKeyF24
	case glfw.KeyKPDecimal:
		scanCode = events.VirtualKeyDecimal
	case glfw.KeyKPDivide:
		scanCode = events.VirtualKeyDivide
	case glfw.KeyKPMultiply:
		scanCode = events.VirtualKeyMultiply
	case glfw.KeyKPSubtract:
		scanCode = events.VirtualKeySubtract
	case glfw.KeyKPAdd:
		scanCode = events.VirtualKeyAdd
	case glfw.KeyKPEnter:
		scanCode = events.VirtualKeyReturn
	case glfw.KeyKPEqual:
		scanCode = events.VirtualKeyEquals
	case glfw.KeyLeftShift:
		scanCode = events.VirtualKeyLShift
	case glfw.KeyLeftControl:
		scanCode = events.VirtualKeyLControl
	case glfw.KeyLeftAlt:
		scanCode = events.VirtualKeyLAlt
	case glfw.KeyLeftSuper:
		scanCode = events.VirtualKeyLWin
	case glfw.KeyRightShift:
		scanCode = events.VirtualKeyRShift
	case glfw.KeyRightControl:
		scanCode = events.VirtualKeyRControl
	case glfw.KeyRightAlt:
		scanCode = events.VirtualKeyRAlt
	case glfw.KeyRightSuper:
		scanCode = events.VirtualKeyRWin
	case glfw.KeyMenu:
		scanCode = events.VirtualKeyContextMenu
	}
	return scanCode
}
