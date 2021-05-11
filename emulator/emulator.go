package emulator

import (
	"github.com/kopi22/chip8/emulator/io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const DisplayWidth = 64
const DisplayHeight = 32

const SpriteWidth = 8

const DefaultEmuSpeed = 2 * time.Millisecond
const TimerPeriod = 17 * time.Millisecond

type Emulator struct {
	chipState *State
	io        io.IO
}

func NewEmulator() *Emulator {
	return &Emulator{
		chipState: InitChipState(),
	}
}

func (emu *Emulator) ConnectIO(io io.IO) *Emulator {
	// detach current IO
	if emu.io != nil {
		emu.io.Clear()
		emu.io.Fini()
	}

	// initialize new IO
	emu.io = io
	emu.io.Init()

	return emu
}

func (emu *Emulator) Launch(romFilename string) {
	emu.LoadRom(romFilename)

	inputChan := make(chan io.InputEvent)
	go emu.io.FetchInputEvents(inputChan)

	timerTicker := time.NewTicker(TimerPeriod)
	cpuTicker := time.NewTicker(DefaultEmuSpeed)

	for {
		select {
		case <-cpuTicker.C:
			emu.Step()
			// Update screen
			emu.io.Draw(emu.chipState.FrameBuf)
		case <-timerTicker.C:
			if emu.chipState.Delay > 0 {
				emu.chipState.Delay -= 1
			}
			if emu.chipState.Sound > 0 {
				emu.chipState.Sound -= 1
			}
		case ev := <-inputChan:
			emu.handleInputEvent(ev)
		}
	}
}

func (emu *Emulator) exit(exitCode int) {
	emu.io.Fini()
	os.Exit(exitCode)
}

func (emu *Emulator) exitWithError(exitCode int, err error) {
	log.Printf("%+v", err)
	emu.exit(exitCode)

}

func (emu *Emulator) handleInputEvent(event io.InputEvent) {
	switch event.EventType {
	case io.KeyDown:
		emu.chipState.Keyboard = uint16(event.EventKey)
	case io.KeyUp:
		if event.EventKey == io.Key(emu.chipState.Keyboard) {
			emu.chipState.Keyboard = 0
		}
	case io.Quit:
		emu.exit(0)
	}
}

func (emu *Emulator) LoadRom(filename string) {
	// read CHIP-8 instructions
	filepath := "roms/" + filename
	sourcecode, err := ioutil.ReadFile(filepath)
	if err != nil {
		emu.exitWithError(1, err)
	}

	copy(emu.chipState.Memory[INITIAL_PC:], sourcecode)
}

func (emu *Emulator) Step() {
	// fetch instruction from Memory
	instruction := FetchInstruction(emu.chipState.Memory, emu.chipState.PC)

	// increase program counter
	emu.chipState.PC += 2

	emu.executeInstruction(instruction)
}

func (emu *Emulator) executeInstruction(instruction Instruction) {
	switch instruction >> 12 {
	case 0x0:
		Op0(emu.chipState, instruction)
	case 0x1:
		Op1(emu.chipState, instruction)
	case 0x2:
		Op2(emu.chipState, instruction)
	case 0x3:
		Op3(emu.chipState, instruction)
	case 0x4:
		Op4(emu.chipState, instruction)
	case 0x5:
		Op5(emu.chipState, instruction)
	case 0x6:
		Op6(emu.chipState, instruction)
	case 0x7:
		Op7(emu.chipState, instruction)
	case 0x8:
		Op8(emu.chipState, instruction)
	case 0x9:
		Op9(emu.chipState, instruction)
	case 0xA:
		OpA(emu.chipState, instruction)
	case 0xB:
		OpB(emu.chipState, instruction)
	case 0xC:
		OpC(emu.chipState, instruction)
	case 0xD:
		OpD(emu.chipState, instruction)
	case 0xE:
		OpE(emu.chipState, instruction)
	case 0xF:
		OpF(emu.chipState, instruction)
	default:
		UnsupportedInstruction(instruction)
	}
}
