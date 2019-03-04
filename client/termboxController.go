package client

import (
	"time"

	"github.com/pkg/errors"

	termbox "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

// termController is a wrapper for basic Termbox functionality
type termController struct {
	eventQueue chan termbox.Event
	keyCh      chan string
	errorCh    chan error
	stopCh     chan struct{}
	height     int
	width      int
	open       bool
}

// newTermController returns a pointer to a new un-initialized termController struct
func newTermController(keyCh chan string, errorCh chan error, stopCh chan struct{}) *termController {
	return &termController{
		eventQueue: make(chan termbox.Event, 16),
		keyCh:      keyCh,
		errorCh:    errorCh,
		stopCh:     stopCh,
		open:       false,
	}
}

// init sets up the termController and starts Termbox
func (term *termController) init() {
	log.Debug("Initializing Termbox")
	if err := termbox.Init(); err != nil {
		term.errorCh <- errors.Wrap(err, "Cannot initialize Termbox")
		return
	}
	termbox.SetInputMode(termbox.InputEsc)
	term.resize()
	go term.eventFeeder()
	go term.eventReader()
	term.open = true
	log.Debug("Termbox initialized successfully")
}

// close stops Termbox
func (term *termController) close() {
	termbox.Close()
	term.open = false
	log.Debug("Termbox closed successfully")
}

// isOpen reports if Termbox is initialized
func (term *termController) isOpen() bool {
	return term.open
}

// clear clears the termbox with given foreground and background colors
func (term *termController) clear(colors ...termbox.Attribute) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	if len(colors) == 1 {
		fg = colors[0]
		bg = colors[0]
	}
	if len(colors) > 1 {
		fg = colors[0]
		bg = colors[1]
	}
	if err := termbox.Clear(fg, bg); err != nil {
		term.errorCh <- errors.Wrap(err, "Cannot clear Termbox")
	}
}

// render draws the Termbox buffer onto the screen
func (term *termController) render() {
	if err := termbox.Flush(); err != nil {
		term.errorCh <- errors.Wrap(err, "Cannot flush Termbox buffer")
	}
}

// renderAndWait will call render then wait until the given time
func (term *termController) renderAndWait(t time.Time) {
	term.render()
	time.Sleep(time.Until(t))
}

// resize will set termControllers width and height property
// according to Termbox
func (term *termController) resize() {
	term.width, term.height = termbox.Size()
}

// getSize will return the termController's width and height
func (term *termController) getSize() (width, height int) {
	return term.width, term.height
}

// eventFeeder will constantly read Termbox events and push
// them to the termControllers event queue
func (term *termController) eventFeeder() {
	for {
		select {
		case <-term.stopCh:
			log.Debug("Termbox eventFeeder stopped")
			return
		default:
			term.eventQueue <- termbox.PollEvent()
		}
	}
}

// eventReader reads the termbox event queue and feeds the key and error chans
func (term *termController) eventReader() {
	for {
		select {

		case <-term.stopCh:
			log.Debug("Termbox eventReader stopped")
			return

		case ev := <-term.eventQueue:

			switch ev.Type {

			case termbox.EventKey:

				switch ev.Key {
				case termbox.KeyArrowUp:
					term.keyCh <- "up"
				case termbox.KeyArrowDown:
					term.keyCh <- "down"
				case termbox.KeyArrowLeft:
					term.keyCh <- "left"
				case termbox.KeyArrowRight:
					term.keyCh <- "right"
				case termbox.KeyEnter:
					term.keyCh <- "enter"
				case termbox.KeySpace:
					term.keyCh <- "space"
				case termbox.KeyEsc:
					term.keyCh <- "esc"
				}

			case termbox.EventResize:
				term.resize()

			case termbox.EventError:
				term.errorCh <- ev.Err

			case termbox.EventInterrupt:
				term.errorCh <- errors.New("Interrupted")

			}
		}
	}
}
