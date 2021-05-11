package main

import (
	"fmt"
	"io/ioutil"
)

type Instruction uint16

func (instruction Instruction) getNNN() uint16 {
	return uint16(instruction & 0x0FFF)
}

func (instruction Instruction) getY() byte {
	return byte((instruction & 0x00F0) >> 4)
}

func (instruction Instruction) getX() byte {
	return byte((instruction & 0x0F00) >> 8)
}

func (instruction Instruction) getKK() byte {
	return byte(instruction & 0x00FF)
}

func (instruction Instruction) getN() byte {
	return byte(instruction & 0x000F)
}

func DisassembleInstruction(codeBuffer []byte, pc int) {
	// instructions must be at even locations in the TEXT segment
	if pc%2 == 1 {
		fmt.Printf("Not-aligned PC! (%04X)\n", pc)
	}

	fmt.Printf("0x%03X - ", pc)

	instruction := Instruction((uint16(codeBuffer[pc]) << 8) | uint16(codeBuffer[pc+1]))
	// first NIBBLE determines the instruction type
	switch instruction >> 12 {
	case 0x0:
		switch uint16(instruction) {
		case 0x00E0:
			fmt.Printf("CLS")
		case 0x00EE:
			fmt.Printf("RET")
		default:
			// SYS addr - not implemented (skip)
			fmt.Printf("NOP")
		}
	case 0x1:
		addr := instruction.getNNN()
		fmt.Printf("JP 0x%03X", addr)
	case 0x2:
		addr := instruction.getNNN()
		fmt.Printf("CALL 0x%03X", addr)
	case 0x3:
		reg := instruction.getX()
		val := instruction.getKK()
		fmt.Printf("SE V%X, $%02X", reg, val)
	case 0x4:
		reg := instruction.getX()
		val := instruction.getKK()
		fmt.Printf("SNE V%X, $%02X", reg, val)
	case 0x5:
		switch instruction & 0x000F {
		case 0x0:
			regX := instruction.getX()
			regY := instruction.getY()
			fmt.Printf("SE V%X, V%X", regX, regY)
		default:
			fmt.Printf("Instruction %04X not yet implemented", instruction)
		}
	case 0x6:
		dst := instruction.getX()
		val := instruction.getKK()
		fmt.Printf("LD V%X, $%02X", dst, val)
	case 0x7:
		dst := instruction.getX()
		val := instruction.getKK()
		fmt.Printf("ADD V%X, $%02X", dst, val)
	case 0x8:
		switch instruction & 0x000F {
		case 0x0:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("LD V%X, V%X", dst, src)
		case 0x1:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("OR V%X, V%X", dst, src)
		case 0x2:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("AND V%X, V%X", dst, src)
		case 0x3:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("XOR V%X, V%X", dst, src)
		case 0x4:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("ADD V%X, V%X", dst, src)
		case 0x5:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("SUB V%X, V%X", dst, src)
		case 0x6:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("SHR V%X {, V%X}", dst, src)
		case 0x7:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("SUBN V%X, V%X", dst, src)
		case 0xe:
			dst := instruction.getX()
			src := instruction.getY()
			fmt.Printf("SHL V%X {, V%X}", dst, src)
		default:
			fmt.Printf("Instruction %04X not yet implemented", instruction)
		}
	case 0x9:
		switch instruction & 0x000F {
		case 0x0:
			regX := instruction.getX()
			regY := instruction.getY()
			fmt.Printf("SNE V%X, V%X", regX, regY)
		default:
			fmt.Printf("Instruction %04X not yet implemented", instruction)
		}
	case 0xa:
		addr := instruction.getNNN()
		fmt.Printf("LD I, $%03X", addr)
	case 0xb:
		addr := instruction.getNNN()
		fmt.Printf("JP V0, $%03X", addr)
	case 0xc:
		dst := instruction.getX()
		val := instruction.getKK()
		fmt.Printf("RND V%X, $%02X", dst, val)
	case 0xd:
		x := instruction.getX()
		y := instruction.getY()
		n := instruction.getN()
		fmt.Printf("DRW V%X, V%X, $%X", x, y, n)
	case 0xe:
		switch instruction & 0x00FF {
		case 0x9e:
			reg := instruction.getX()
			fmt.Printf("SKP V%X", reg)
		case 0xA1:
			reg := instruction.getX()
			fmt.Printf("SKNP V%X", reg)
		default:
			fmt.Printf("Instruction %04X not yet implemented", instruction)
		}
	case 0xf:
		switch instruction & 0xFF {
		case 0x07:
			dst := instruction.getX()
			fmt.Printf("LD V%X, DT", dst)
		case 0x0A:
			dst := instruction.getX()
			fmt.Printf("LD V%X, K", dst)
		case 0x15:
			src := instruction.getX()
			fmt.Printf("LD DT, V%X", src)
		case 0x18:
			src := instruction.getX()
			fmt.Printf("LD ST, V%X", src)
		case 0x1e:
			reg := instruction.getX()
			fmt.Printf("ADD I, V%X", reg)
		case 0x29:
			reg := instruction.getX()
			fmt.Printf("LD F, V%X", reg)
		case 0x33:
			reg := instruction.getX()
			fmt.Printf("LD B, V%X", reg)
		case 0x55:
			reg := instruction.getX()
			fmt.Printf("LD [I], V%X", reg)
		case 0x65:
			reg := instruction.getX()
			fmt.Printf("LD V%X, [I]", reg)
		default:
			fmt.Printf("Instruction %04X not yet implemented", instruction)
		}
	default:
		fmt.Printf("Instruction %04X not yet implemented", instruction)
	}
	fmt.Println()
}

func main() {
	sourcecodeFilename := "roms/Sierpinski.ch8"

	// read CHIP-8 instructions
	sourcecode, err := ioutil.ReadFile(sourcecodeFilename)
	if err != nil {
		fmt.Printf("Cannot open the file: %v\n", sourcecodeFilename)
		panic(err)
	}

	// instructions start at 0x200
	TEXT := make([]byte, 0x200+len(sourcecode))
	copy(TEXT[0x200:], sourcecode)

	for pc := 0x200; pc < len(TEXT); pc += 2 {
		DisassembleInstruction(TEXT, pc)
	}
}
