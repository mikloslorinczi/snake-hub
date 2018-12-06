package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pkg/errors"

	"github.com/mikloslorinczi/snake-hub/modell"
)

type termboxWorker struct {
	stopEventLoop chan struct{}
	stopRender    chan struct{}
	eventQueue    chan termbox.Event
	height        int
	width         int
}

func (term *termboxWorker) startEventloop() {
	term.eventQueue = make(chan termbox.Event)
	go func() {
		for {
			term.eventQueue <- termbox.PollEvent()
		}
	}()
	fmt.Println("Starting Termbox event reader")
	termbox.SetInputMode(termbox.InputEsc)
	for {
		select {
		case ev := <-term.eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				fmt.Printf("Event key %U\n", ev.Key)
				switch ev.Key {
				case termbox.KeyEsc:
					errorChan <- errors.New("Terminated by user")
				case termbox.KeyArrowUp:
					postMsg("control", "up")
				case termbox.KeyArrowDown:
					postMsg("control", "down")
				case termbox.KeyArrowLeft:
					postMsg("control", "left")
				case termbox.KeyArrowRight:
					postMsg("control", "right")
				}
			case termbox.EventResize:
				term.getSize()
			case termbox.EventError:
				errorChan <- ev.Err
			case termbox.EventInterrupt:
				errorChan <- errors.New("Interrupted")
			}
		case <-term.stopEventLoop:
			return
		}
	}
}

// Render draws the scene...
func (term *termboxWorker) startRenderer() {
	for {
		select {

		case <-term.stopRender:
			return

		default:

			nextTick := time.Now().Add(time.Millisecond * 30) // ~33.3 FPS

			if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
				errorChan <- err
			}

			term.draw()

			if err := termbox.Flush(); err != nil {
				errorChan <- errors.Wrap(err, fmt.Sprintf("Cannot sync Termbox"))
			}

			time.Sleep(time.Until(nextTick))

		}
	}
}

func (term *termboxWorker) getSize() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	term.width, term.height = termbox.Size()
}

func putBlock(b modell.Block) {
	termbox.SetCell(b.Coord.X*2, b.Coord.Y, b.LeftRune, b.Color, b.Background)
	termbox.SetCell(b.Coord.X*2+1, b.Coord.Y, b.RightRune, b.Color, b.Background)
}

func (term *termboxWorker) draw() {

	for y := 0; y < term.height; y++ {
		for x := 0; x < term.width/2; x++ {
			putBlock(modell.Block{
				Coord: modell.Coords{
					X: x,
					Y: y,
				},
				Color:      termbox.ColorDefault,
				Background: termbox.Attribute(rand.Int()%8) + 1,
				LeftRune:   ' ',
				RightRune:  ' ',
			})
		}
	}
}
