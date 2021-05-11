package tcellIO

import (
	"github.com/gdamore/tcell/v2"
	"github.com/kopi22/chip8/emulator"
	"github.com/kopi22/chip8/emulator/io"
	"log"
	"time"
)

const KeyPressDuration = 100 * time.Millisecond

func getDefaultDisplayStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
}

func getBorderStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorBlack)
}

func getPixelOnStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
}

func getPixelOffStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorBlack)
}

type IO struct {
	Screen tcell.Screen
}

func (tcellIO *IO) Init() {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	tcellIO.Screen = s

	tcellIO.Screen.SetStyle(getDefaultDisplayStyle())
	tcellIO.Screen.Clear()

	// Draw Chip Display
	drawBox(tcellIO.Screen, 0, 0, 65, 33, getBorderStyle())
}

func (tcellIO *IO) Fini() {
	tcellIO.Screen.Fini()
}

func (tcellIO *IO) Draw(frameBuffer []byte) {
	pixelOffStyle := getPixelOffStyle()
	pixelOnStyle := getPixelOnStyle()

	for r := 0; r < emulator.DisplayHeight; r++ {
		for c := 0; c < emulator.DisplayWidth; c++ {
			totalOffset := r*emulator.DisplayWidth + c
			byteOffset, bitOffset := totalOffset/8, totalOffset%8
			pixelMask := byte(0x80 >> bitOffset)

			if frameBuffer[byteOffset]&pixelMask == 0 {
				tcellIO.Screen.SetContent(c+1, r+1, ' ', nil, pixelOffStyle)
			} else {
				tcellIO.Screen.SetContent(c+1, r+1, ' ', nil, pixelOnStyle)
			}

		}
	}

	tcellIO.Screen.Show()
}

func (tcellIO *IO) Clear() {
	tcellIO.Screen.Clear()
}

type keyPressTimer struct {
	*time.Timer
	key io.Key
}

func (tcellIO *IO) FetchInputEvents(inputChan chan<- io.InputEvent) {
	pressTimer := keyPressTimer{
		time.NewTimer(0), io.Key(0),
	}

	eventChan := make(chan tcell.Event)
	go func() {
		for {
			ev := tcellIO.Screen.PollEvent()
			eventChan <- ev
		}
	}()

	// Event loop
	for {

		select {

		case ev := <-eventChan:
			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				tcellIO.Screen.Sync()
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlC:
					inputChan <- io.InputEvent{
						EventType: io.Quit,
					}
				case tcell.KeyRune:
					key, ok := io.DefaultKeyboardMap[ev.Rune()]
					if ok {
						// reset key press timer
						pressTimer.Stop()
						pressTimer.Reset(KeyPressDuration)
						pressTimer.key = key

						inputChan <- io.InputEvent{
							EventType: io.KeyDown,
							EventKey:  key,
						}
					}
				}
			}

		// simulated key release
		case <-pressTimer.C:
			// clear the timer
			key := pressTimer.key
			pressTimer.key = io.Key(0)

			// send the KeyUp event
			inputChan <- io.InputEvent{
				EventType: io.KeyUp,
				EventKey:  key,
			}

		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

//func (tcellIO IO) drawPixel(col, row int) {
//	// TODO: IMPLEMENT
//}
