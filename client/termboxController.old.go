package client

// import (
// 	"time"

// 	"github.com/nsf/termbox-go"
// 	"github.com/pkg/errors"
// 	log "github.com/sirupsen/logrus"
// )

// type termboxController struct {
// 	stopEventLoop chan struct{}
// 	stopRender    chan struct{}
// 	stopQueue     chan struct{}
// 	eventQueue    chan termbox.Event
// 	height        int
// 	width         int
// }

// func (term *termboxController) startEventloop() {
// 	term.eventQueue = make(chan termbox.Event)

// 	go func() {
// 		for {
// 			select {
// 			case <-term.stopQueue:
// 				return
// 			default:
// 				term.eventQueue <- termbox.PollEvent()
// 			}
// 		}
// 	}()

// 	log.Info("Starting Termbox event reader")
// 	termbox.SetInputMode(termbox.InputEsc)
// 	for {
// 		select {
// 		case ev := <-term.eventQueue:
// 			log.Debug(ev)
// 			log.Debug(string(ev.Ch))
// 			switch ev.Type {
// 			case termbox.EventKey:
// 				switch ev.Key {
// 				case termbox.KeyEsc:
// 					errorChan <- errors.New("Terminated by user")
// 				case termbox.KeyArrowUp:
// 					postMsg("control", "up")
// 				case termbox.KeyArrowDown:
// 					log.Warn("Down pressed")
// 					postMsg("control", "down")
// 				case termbox.KeyArrowLeft:
// 					postMsg("control", "left")
// 				case termbox.KeyArrowRight:
// 					postMsg("control", "right")
// 				}
// 			case termbox.EventResize:
// 				term.getSize()
// 			case termbox.EventError:
// 				errorChan <- ev.Err
// 			case termbox.EventInterrupt:
// 				errorChan <- errors.New("Interrupted")
// 			}
// 		case <-term.stopEventLoop:
// 			return
// 		}
// 	}
// }

// func (term *termboxController) clearAll() {
// 	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
// 		errorChan <- errors.Wrap(err, "Cannot clear Termbox")
// 	}
// }

// func (term *termboxController) renderAll() {
// 	if err := termbox.Flush(); err != nil {
// 		errorChan <- errors.Wrap(err, "Cannot sync Termbox")
// 	}
// }

// func (term *termboxController) renderAndSleep(t time.Time) {
// 	term.renderAll()
// 	time.Sleep(time.Until(t))
// }

// func (term *termboxController) getSize() {
// 	term.width, term.height = termbox.Size()
// }

// // Render draws the Scene...
// func (term *termboxController) startRenderer() {
// 	log.Info("Start rendering on TermBox")
// mainLoop:
// 	for {

// 		select {

// 		case <-term.stopRender:
// 			return

// 		default:
// 			// Calculate next tick
// 			nextTick := time.Now().Add(time.Millisecond * 30) // ~33.3 FPS
// 			// Clear everything
// 			term.clearAll()
// 			// Draw resize message if terminal is too small
// 			w, h := 20, 10
// 			if state.loaded {
// 				w, h = state.getLvlSize()
// 			}
// 			if term.width < w*2 || term.height < h {
// 				term.drawResizeMsg(w*2, h)
// 				term.renderAndSleep(nextTick)
// 				continue mainLoop
// 			}
// 			// Do nothing if still not got the first state-update
// 			if !state.loaded {
// 				term.renderAndSleep(nextTick)
// 				continue mainLoop
// 			}
// 			// Depending on the game scene call the appropriate draw function
// 			switch state.getScene() {
// 			case "game":
// 				term.drawGame()
// 			default:
// 				term.drawTextbox()
// 			}
// 			// Display everything we have drawn so far and wait the next tick
// 			term.renderAndSleep(nextTick)
// 		}
// 	}
// }
