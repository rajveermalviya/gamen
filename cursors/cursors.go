package cursors

type Icon uint8

const (
	Default Icon = iota + 1
	ContextMenu
	Help
	Pointer
	Progress
	Wait
	Cell
	Crosshair
	Text
	VerticalText
	Alias
	Copy
	Move
	NoDrop
	NotAllowed
	Grab
	Grabbing
	AllScroll
	ColResize
	RowResize
	NResize
	EResize
	SResize
	WResize
	NEResize
	NWResize
	SEResize
	SWResize
	EWResize
	NSResize
	NESWResize
	NWSEResize
	ZoomIn
	ZoomOut
)

func (i Icon) String() string {
	switch i {
	case Default:
		return "default"
	case ContextMenu:
		return "context-menu"
	case Help:
		return "help"
	case Pointer:
		return "pointer"
	case Progress:
		return "progress"
	case Wait:
		return "wait"
	case Cell:
		return "cell"
	case Crosshair:
		return "crosshair"
	case Text:
		return "text"
	case VerticalText:
		return "vertical-text"
	case Alias:
		return "alias"
	case Copy:
		return "copy"
	case Move:
		return "move"
	case NoDrop:
		return "no-drop"
	case NotAllowed:
		return "not-allowed"
	case Grab:
		return "grab"
	case Grabbing:
		return "grabbing"
	case AllScroll:
		return "all-scroll"
	case ColResize:
		return "col-resize"
	case RowResize:
		return "row-resize"
	case NResize:
		return "n-resize"
	case EResize:
		return "e-resize"
	case SResize:
		return "s-resize"
	case WResize:
		return "w-resize"
	case NEResize:
		return "ne-resize"
	case NWResize:
		return "nw-resize"
	case SEResize:
		return "se-resize"
	case SWResize:
		return "sw-resize"
	case EWResize:
		return "ew-resize"
	case NSResize:
		return "ns-resize"
	case NESWResize:
		return "nesw-resize"
	case NWSEResize:
		return "nwse-resize"
	case ZoomIn:
		return "zoom-in"
	case ZoomOut:
		return "zoom-out"
	}

	panic("unreachable")
}
