package emulator

import (
	"fmt"
	"math"
	"math/rand"
)

func UnsupportedInstruction(instruction Instruction) {
	// Add PC for debugging?
	fmt.Printf("Instruction %04X is not implemented.", instruction)
}

func Op0(chipState *State, instruction Instruction) {
	switch uint16(instruction) {

	case 0x00E0: // CLS
		for i := range chipState.FrameBuf {
			chipState.FrameBuf[i] = 0
		}

	case 0x00EE: // RET
		chipState.SP--
		chipState.PC = chipState.Stack[chipState.SP]

	default:
		UnsupportedInstruction(instruction)
	}
}

func Op1(chipState *State, instruction Instruction) {
	// JP addr
	chipState.PC = instruction.GetNNN()
}

func Op2(chipState *State, instruction Instruction) {
	// CALL addr
	// PC is already incremented
	chipState.Stack[chipState.SP] = chipState.PC
	chipState.SP++
	chipState.PC = instruction.GetNNN()
}

func Op3(chipState *State, instruction Instruction) {
	// SE Vx, byte
	if chipState.V[instruction.GetX()] == instruction.GetKK() {
		chipState.PC += 2
	}
}

func Op4(chipState *State, instruction Instruction) {
	// SNE Vx, byte
	if chipState.V[instruction.GetX()] != instruction.GetKK() {
		chipState.PC += 2
	}
}

func Op5(chipState *State, instruction Instruction) {
	//  SE Vx, Vy
	if chipState.V[instruction.GetX()] == chipState.V[instruction.GetY()] {
		chipState.PC += 2
	}
}

func Op6(chipState *State, instruction Instruction) {
	// LD Vx, byte
	chipState.V[instruction.GetX()] = instruction.GetKK()
}

func Op7(chipState *State, instruction Instruction) {
	// ADD Vx, byte
	chipState.V[instruction.GetX()] += instruction.GetKK()
}

func Op8(chipState *State, instruction Instruction) {
	switch instruction & 0x000F {
	case 0x0:
		// LD Vx, Vy
		chipState.V[instruction.GetX()] = chipState.V[instruction.GetY()]

	case 0x1:
		// OR Vx, Vy
		chipState.V[instruction.GetX()] |= chipState.V[instruction.GetY()]

	case 0x2:
		// AND Vx, Vy
		chipState.V[instruction.GetX()] &= chipState.V[instruction.GetY()]

	case 0x3:
		// XOR Vx, Vy
		chipState.V[instruction.GetX()] ^= chipState.V[instruction.GetY()]

	case 0x4:
		// ADD Vx, Vy
		result := uint16(chipState.V[instruction.GetX()]) + uint16(chipState.V[instruction.GetY()])
		if result > 0xFF {
			chipState.V[0xF] = 1
		} else {
			chipState.V[0xF] = 0
		}
		chipState.V[instruction.GetX()] = byte(result)

	case 0x5:
		// SUB Vx, Vy
		vx := &chipState.V[instruction.GetX()]
		vy := &chipState.V[instruction.GetY()]

		if *vx > *vy {
			chipState.V[0xF] = 1
		} else {
			chipState.V[0xF] = 0
		}

		*vx -= *vy

	case 0x6:
		// SHR Vx {, Vy}
		vx := &chipState.V[instruction.GetX()]

		chipState.V[0xF] = 0x1 & (*vx)

		*vx >>= 1

	case 0x7:
		// SUBN Vx, Vy
		vx := &chipState.V[instruction.GetX()]
		vy := &chipState.V[instruction.GetY()]

		if *vy > *vx {
			chipState.V[0xF] = 1
		} else {
			chipState.V[0xF] = 0
		}

		*vx = *vy - *vx

	case 0xe:
		// SHL Vx {, Vy}
		vx := &chipState.V[instruction.GetX()]

		chipState.V[0xF] = (0x80 & (*vx)) >> 7

		*vx <<= 1

	default:
		UnsupportedInstruction(instruction)
	}
}

func Op9(chipState *State, instruction Instruction) {
	//  SE Vx, Vy
	if chipState.V[instruction.GetX()] != chipState.V[instruction.GetY()] {
		chipState.PC += 2
	}
}

func OpA(chipState *State, instruction Instruction) {
	// LD I, addr
	chipState.I = instruction.GetNNN()
}

func OpB(chipState *State, instruction Instruction) {
	// JP V0, addr
	chipState.PC = uint16(chipState.V[0]) + instruction.GetNNN()
}

func OpC(chipState *State, instruction Instruction) {
	// RND Vx, byte
	chipState.V[instruction.GetX()] = instruction.GetKK() & byte(rand.Int())
}

func OpD(chipState *State, instruction Instruction) {
	// DRW Vx, Vy, nibble

	// clear collision indicator
	chipState.V[0xF] = 0

	bytesToRead := int(instruction.GetN())
	initX, initY := int(chipState.V[instruction.GetX()]), int(chipState.V[instruction.GetY()])

	sprite := chipState.Memory[chipState.I : chipState.I+uint16(bytesToRead)]

	for r := 0; r < bytesToRead; r++ {
		for c := 0; c < SpriteWidth; c++ {
			// draw only if sprite bit is 1
			if sprite[r]&(0x80>>c) != 0 {
				// position of the pixel on the screen (with wrap-around)
				x, y := (c+initX)%DisplayWidth, (r+initY)%DisplayHeight

				// computer address of corresponding place in memory
				// (with respect to the beginning of the frame buffer segment)
				totalOffset := y*DisplayWidth + x
				byteOffset, bitOffset := totalOffset/8, totalOffset%8
				pixelMask := byte(0x80 >> bitOffset)

				chipState.FrameBuf[byteOffset] ^= pixelMask

				// mark collision if display pixel is OFF
				if chipState.FrameBuf[byteOffset]&pixelMask == 0 {
					chipState.V[0xF] = 1
				}
			}
		}
	}
}

func OpE(chipState *State, instruction Instruction) {
	switch instruction & 0x00FF {
	case 0x9e:
		// SKP Vx
		key := chipState.V[instruction.GetX()]
		var keyMask uint16 = 1 << key

		if chipState.Keyboard&keyMask != 0 {
			chipState.PC += 2
		}

	case 0xA1:
		// SKNP Vx
		key := chipState.V[instruction.GetX()]
		var keyMask uint16 = 1 << key

		if chipState.Keyboard&keyMask == 0 {
			chipState.PC += 2
		}

	default:
		UnsupportedInstruction(instruction)
	}
}

func OpF(chipState *State, instruction Instruction) {
	switch instruction & 0xFF {
	case 0x07:
		// LD Vx, DT
		chipState.V[instruction.GetX()] = chipState.Delay

	case 0x0A:
		//  LD Vx, K
		if chipState.Keyboard != 0 {
			// determine the key based on which bit is set
			chipState.V[instruction.GetX()] = byte(math.Log2(float64(chipState.Keyboard)))
		} else {
			// repeat instruction
			chipState.PC -= 2
		}

	case 0x15:
		// LD DT, Vx
		chipState.Delay = chipState.V[instruction.GetX()]

	case 0x18:
		// LD ST, Vx
		chipState.Sound = chipState.V[instruction.GetX()]

	case 0x1e:
		// ADD I, Vx
		chipState.I += uint16(chipState.V[instruction.GetX()])

	case 0x29:
		//  LD F, Vx
		chipState.I = FONTSET_LOCATION + uint16(chipState.V[instruction.GetX()]*5) // each font sprite takes 5 bytes

	case 0x33:
		// LD B, Vx
		vx := chipState.V[instruction.GetX()]
		digits := fmt.Sprintf("%03d", vx)

		// convert from string (bytes) to int values
		chipState.Memory[chipState.I] = digits[0] - 0x30
		chipState.Memory[chipState.I+1] = digits[1] - 0x30
		chipState.Memory[chipState.I+2] = digits[2] - 0x30

	case 0x55:
		lastRegToStore := uint16(instruction.GetX())

		for i := uint16(0); i <= lastRegToStore; i++ {
			chipState.Memory[chipState.I+i] = chipState.V[i]
		}

	case 0x65:
		lastRegToLoad := uint16(instruction.GetX())

		for i := uint16(0); i <= lastRegToLoad; i++ {
			chipState.V[i] = chipState.Memory[chipState.I+i]
		}

	default:
		UnsupportedInstruction(instruction)
	}
}
