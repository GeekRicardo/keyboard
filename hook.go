package keyboardhook

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")

	hook HHOOK
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
)

type KBDLLHOOKSTRUCT struct {
	VKCode     uint32
	ScanCode   uint32
	Flags      uint32
	Time       uint32
	DwExtraInfo uintptr
}

type HHOOK uintptr

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X, Y int32
}

// KeyEvent represents a keyboard event
type KeyEvent struct {
	VKCode uint32 // Virtual key code
	Name   string // Readable name of the key
}

// Key codes and names map
var keyNames = map[uint32]string{
	0x41: "A",
	0x42: "B",
	0x43: "C",
	0x44: "D",
	0x45: "E",
	0x46: "F",
	0x47: "G",
	0x48: "H",
	0x49: "I",
	0x4A: "J",
	0x4B: "K",
	0x4C: "L",
	0x4D: "M",
	0x4E: "N",
	0x4F: "O",
	0x50: "P",
	0x51: "Q",
	0x52: "R",
	0x53: "S",
	0x54: "T",
	0x55: "U",
	0x56: "V",
	0x57: "W",
	0x58: "X",
	0x59: "Y",
	0x5A: "Z",
	0x30: "0",
	0x31: "1",
	0x32: "2",
	0x33: "3",
	0x34: "4",
	0x35: "5",
	0x36: "6",
	0x37: "7",
	0x38: "8",
	0x39: "9",
	0x0D: "Enter",
	0x1B: "Escape",
	0x08: "Backspace",
	0x09: "Tab",
	0x20: "Space",
	0xA0: "Left Shift",
	0xA1: "Right Shift",
	0xA2: "Left Control",
	0xA3: "Right Control",
	0xA4: "Left Alt",
	0xA5: "Right Alt",
	0x70: "F1",
	0x71: "F2",
	0x72: "F3",
	0x73: "F4",
	0x74: "F5",
	0x75: "F6",
	0x76: "F7",
	0x77: "F8",
	0x78: "F9",
	0x79: "F10",
	0x7A: "F11",
	0x7B: "F12",
	// Add more keys as needed
}

// GetKeyName returns the readable name of the key
func GetKeyName(vkCode uint32) string {
	if name, exists := keyNames[vkCode]; exists {
		return name
	}
	return "Unknown"
}

// Callback functions
var (
	onKeyDown func(KeyEvent)
	onKeyUp   func(KeyEvent)
)

// SetKeyDownCallback sets the callback function for key down events
func SetKeyDownCallback(callback func(KeyEvent)) {
	onKeyDown = callback
}

// SetKeyUpCallback sets the callback function for key up events
func SetKeyUpCallback(callback func(KeyEvent)) {
	onKeyUp = callback
}

// Start starts the keyboard hook
func Start() {
	hook = setHook()
	defer unhook()

	var msg MSG
	for {
		ret, _, err := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(ret) == -1 {
			fmt.Println("GetMessage 错误:", err)
			break
		}
	}
}

func setHook() HHOOK {
	h, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		syscall.NewCallback(keyboardProc),
		0,
		0,
	)
	if h == 0 {
		fmt.Println("设置钩子失败:", err)
	}
	return HHOOK(h)
}

func unhook() {
	procUnhookWindowsHookEx.Call(uintptr(hook))
}

func keyboardProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 {
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		event := KeyEvent{
			VKCode: kbd.VKCode,
			Name:   GetKeyName(kbd.VKCode),
		}

		switch wParam {
		case WM_KEYDOWN:
			if onKeyDown != nil {
				onKeyDown(event)
			}
		case WM_KEYUP:
			if onKeyUp != nil {
				onKeyUp(event)
			}
		}
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}
