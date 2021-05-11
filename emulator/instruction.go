package emulator

type Instruction uint16

func FetchInstruction(memory []byte, pc uint16) Instruction {
	return Instruction((uint16(memory[pc]) << 8) | uint16(memory[pc+1]))
}

func (instruction Instruction) GetNNN() uint16 {
	return uint16(instruction & 0x0FFF)
}

func (instruction Instruction) GetY() byte {
	return byte((instruction & 0x00F0) >> 4)
}

func (instruction Instruction) GetX() byte {
	return byte((instruction & 0x0F00) >> 8)
}

func (instruction Instruction) GetKK() byte {
	return byte(instruction & 0x00FF)
}

func (instruction Instruction) GetN() byte {
	return byte(instruction & 0x000F)
}
