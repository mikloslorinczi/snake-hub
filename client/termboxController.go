package client

import (
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

// Render draws the Scene...
func (term *termboxController) startRenderer() {
	for {
		select {

		case <-term.stopRender:
			return

		default:

			nextTick := time.Now().Add(time.Millisecond * 30) // ~33.3 FPS

			if !state.loaded {
				time.Sleep(time.Until(nextTick))
				continue
			}

			if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
				errorChan <- errors.Wrap(err, "Cannot clear Termbox")
			}

			w, h := state.getLvlSize()
			if term.width < w*2 || term.height < h {
				term.drawResizeMsg(w*2, h)
			} else {
				switch state.getScene() {
				case "game":
					term.drawGame()
				default:
					term.drawTextbox()
				}
			}

			if err := termbox.Flush(); err != nil {
				errorChan <- errors.Wrap(err, "Cannot sync Termbox")
			}

			time.Sleep(time.Until(nextTick))

		}
	}
}

func (term *termboxController) getSize() {
	term.width, term.height = termbox.Size()
}
