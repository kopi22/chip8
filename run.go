package main

import (
	"github.com/kopi22/chip8/emulator"
	"github.com/kopi22/chip8/emulator/io/tcellIO"
)

// TODO:
// - add sound support
// - add scaling support
// - fix first key press issue

func main() {
	// set up emulator
	emu := emulator.NewEmulator().ConnectIO(new(tcellIO.IO))

	emu.Launch("Pong1.ch8")
}
