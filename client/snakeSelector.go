package client

/*
import (
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Selector is a TermBox controller for the snake select screen
type Selector struct {
	stopEventLoop chan struct{}
	stopRender    chan struct{}
	stopQueue     chan struct{}
	eventQueue    chan termbox.Event
	height        int
	width         int
}

func (s *Selector) startEventloop() {
	s.eventQueue = make(chan termbox.Event)

	go func() {
		for {
			select {
			case <-s.stopQueue:
				return
			default:
				s.eventQueue <- termbox.PollEvent()
			}
		}
	}()

	log.Info("Starting Termbox event reader")
	termbox.SetInputMode(termbox.InputEsc)
	for {
		select {
		case ev := <-s.eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					errorChan <- errors.New("Terminated by user")
				case termbox.KeyArrowUp:
					postMsg("control", "up")
				case termbox.KeyArrowDown:
					log.Warn("Down pressed")
					postMsg("control", "down")
				case termbox.KeyArrowLeft:
					postMsg("control", "left")
				case termbox.KeyArrowRight:
					postMsg("control", "right")
				}
			case termbox.EventResize:
				s.getSize()
			case termbox.EventError:
				errorChan <- ev.Err
			case termbox.EventInterrupt:
				errorChan <- errors.New("Interrupted")
			}
		case <-s.stopEventLoop:
			return
		}
	}
}

func (s *Selector) clearAll() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		errorChan <- errors.Wrap(err, "Cannot clear Termbox")
	}
}

func (s *Selector) renderAll() {
	if err := termbox.Flush(); err != nil {
		errorChan <- errors.Wrap(err, "Cannot sync Termbox")
	}
}

func (s *Selector) renderAndSleep(t time.Time) {
	s.renderAll()
	time.Sleep(time.Until(t))
}

func (s *Selector) getSize() {
	s.width, s.height = termbox.Size()
}

// Render draws the Scene...
func (s *Selector) startRenderer() {
	log.Info("Start rendering on TermBox")
mainLoop:
	for {

		select {

		case <-s.stopRender:
			return

		default:
			// Calculate next tick
			nextTick := time.Now().Add(time.Millisecond * 30) // ~33.3 FPS
			// Clear everything
			s.clearAll()
			// Draw resize message if terminal is too small
			w, h := 20, 10
			if state.loaded {
				w, h = state.getLvlSize()
			}
			if s.width < w*2 || s.height < h {
				s.drawResizeMsg(w*2, h)
				s.renderAndSleep(nextTick)
				continue mainLoop
			}
			// Do nothing if still not got the first state-update
			if !state.loaded {
				s.renderAndSleep(nextTick)
				continue mainLoop
			}
			// Depending on the game scene call the appropriate draw function
			switch state.getScene() {
			case "game":
				s.drawGame()
			default:
				s.drawTextbox()
			}
			// Display everything we have draw so far and wait the next tick
			s.renderAndSleep(nextTick)
		}
	}
}
*/
