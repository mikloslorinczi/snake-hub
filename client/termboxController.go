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
					postMsg("controll", "up")
				case termbox.KeyArrowDown:
					postMsg("controll", "down")
				case termbox.KeyArrowLeft:
					postMsg("controll", "left")
				case termbox.KeyArrowRight:
					postMsg("controll", "right")
				}
			case termbox.EventResize:
				termChan <- "resize"
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

func putBlock(b modell.Block) {
	termbox.SetCell(b.X*2, b.Y, b.LeftRune, b.Color, b.Background)
	termbox.SetCell(b.X*2+1, b.Y, b.RightRune, b.Color, b.Background)
}

func draw() {
	w, h := termbox.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w/2; x++ {
			putBlock(modell.Block{
				X:          x,
				Y:          y,
				Color:      termbox.ColorDefault,
				Background: termbox.Attribute(rand.Int()%8) + 1,
				LeftRune:   ' ',
				RightRune:  ' ',
			})
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
			draw()

			if err := termbox.Flush(); err != nil {
				errorChan <- errors.Wrap(err, fmt.Sprintf("Cannot sync Termbox"))
			}

			time.Sleep(time.Until(nextTick))

		}
	}
}
