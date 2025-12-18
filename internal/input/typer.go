package input

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	sendInput     = user32.NewProc("SendInput")
	getAsyncState = user32.NewProc("GetAsyncKeyState")
)

const (
	inputKeyboard    = 1
	keyeventfUnicode = 0x0004
	keyeventfKeyup   = 0x0002
)

type keyInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	inputType uint32
	ki        keyInput
	padding   [8]byte
}

type Typer struct {
	delay time.Duration
}

func NewTyper(delay time.Duration) *Typer {
	return &Typer{delay: delay}
}

func (t *Typer) TypeText(text string) {
	for _, char := range text {
		t.typeChar(char)
		time.Sleep(t.delay)
	}
}

func (t *Typer) TypeSpace() {
	t.typeKey(0x20)
	time.Sleep(t.delay)
}

func (t *Typer) typeChar(char rune) {
	in := input{
		inputType: inputKeyboard,
		ki: keyInput{
			wScan:   uint16(char),
			dwFlags: keyeventfUnicode,
		},
	}
	sendInput.Call(uintptr(1), uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
	time.Sleep(1 * time.Millisecond)

	in.ki.dwFlags = keyeventfUnicode | keyeventfKeyup
	sendInput.Call(uintptr(1), uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
}

func (t *Typer) typeKey(vk uint16) {
	in := input{
		inputType: inputKeyboard,
		ki:        keyInput{wVk: vk},
	}
	sendInput.Call(uintptr(1), uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
	time.Sleep(1 * time.Millisecond)

	in.ki.dwFlags = keyeventfKeyup
	sendInput.Call(uintptr(1), uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
}

func IsKeyPressed(vk int) bool {
	ret, _, _ := getAsyncState.Call(uintptr(vk))
	return ret&0x8000 != 0
}
