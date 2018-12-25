package client

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mikloslorinczi/snake-hub/modell"
)

type termboxController struct {
	stopEventLoop chan struct{}
	stopRender    chan struct{}
	eventQueue    chan termbox.Event
	height        int
	width         int
}

func (term *termboxController) startEventloop() {
	term.eventQueue = make(chan termbox.Event)
	go func() {
		for {
			term.eventQueue <- termbox.PollEvent()
		}
	}()
	log.Info("Starting Termbox event reader")
	termbox.SetInputMode(termbox.InputEsc)
	for {
		select {
		case ev := <-term.eventQueue:
			switch ev.Type {
			case termbox.EventKey:
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
func (term *termboxController) startRenderer() {
	for {
		select {

		case <-term.stopRender:
			return

		default:

			nextTick := time.Now().Add(time.Millisecond * 30) // ~33.3 FPS

			if state.loaded {
				if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
					errorChan <- err
				}

				term.draw()

				if err := termbox.Flush(); err != nil {
					errorChan <- errors.Wrap(err, fmt.Sprintf("Cannot sync Termbox"))
				}
			}

			time.Sleep(time.Until(nextTick))

		}
	}
}

func (term *termboxController) getSize() {
	term.width, term.height = termbox.Size()
}

func putBlock(coords modell.Coords, color, bgColor termbox.Attribute, leftRune, rightRune rune) {
	termbox.SetCell(coords.X*2, coords.Y, leftRune, color, bgColor)
	termbox.SetCell(coords.X*2+1, coords.Y, rightRune, color, bgColor)
}

func (term *termboxController) draw() {
	for _, snake := range state.getSnakes() {
		for i, block := range snake.Body {
			if i == 0 {
				putBlock(modell.Coords{
					X: block.X,
					Y: block.Y,
				},
					snake.Color,
					snake.BgColor,
					snake.HeadRune,
					' ',
				)
			} else {
				putBlock(modell.Coords{
					X: block.X,
					Y: block.Y,
				},
					snake.Color,
					snake.BgColor,
					snake.LeftRune,
					snake.RightRune,
				)

			}
		}
	}
}
