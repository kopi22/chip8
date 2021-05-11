package io

var DefaultKeyboardMap = map[rune]Key{
	'1': Key1, '2': Key2, '3': Key3, '4': KeyC,
	'q': Key4, 'w': Key5, 'e': Key6, 'r': KeyD,
	'a': Key7, 's': Key8, 'd': Key9, 'f': KeyE,
	'z': KeyA, 'x': Key0, 'c': KeyB, 'v': KeyF,
}

type Display interface {
	Draw([]byte)
	Clear()
}

type Keyboard interface {
	FetchInputEvents(chan<- InputEvent)
}

type IO interface {
	Init()
	Fini()
	Display
	Keyboard
}

type EventType string

const (
	KeyDown EventType = "KeyDown"
	KeyUp   EventType = "KeyUp"
	Quit    EventType = "Quit"
)

type Key uint16

const (
	Key0 Key = 0x0001
	Key1 Key = 0x0002
	Key2 Key = 0x0004
	Key3 Key = 0x0008
	Key4 Key = 0x0010
	Key5 Key = 0x0020
	Key6 Key = 0x0040
	Key7 Key = 0x0080
	Key8 Key = 0x0100
	Key9 Key = 0x0200
	KeyA Key = 0x0400
	KeyB Key = 0x0800
	KeyC Key = 0x1000
	KeyD Key = 0x2000
	KeyE Key = 0x4000
	KeyF Key = 0x8000
)

type InputEvent struct {
	EventType EventType
	EventKey  Key
}
