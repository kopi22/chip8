package emulator

const FONTSET_LOCATION = 0x0
const INITIAL_PC = 0x200

type State struct {
	V        [16]byte
	I        uint16
	PC       uint16
	SP       byte
	Delay    byte
	Sound    byte
	Memory   []byte
	FrameBuf []byte
	Stack    [16]uint16
	Keyboard uint16
}

func InitChipState() *State {
	state := &State{
		PC:     INITIAL_PC,
		Memory: make([]byte, 4096), // 4kb
	}
	state.FrameBuf = state.Memory[0xf00:(0xf00 + 64*32/8)]

	copy(state.Memory[FONTSET_LOCATION:], getFontset())

	return state
}

func getFontset() []byte {
	return []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
}
