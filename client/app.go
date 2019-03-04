package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/utils"
	termbox "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type app struct {
	clientID string
	username string

	keyCh       chan string
	stateCh     chan string
	errorCh     chan error
	clientMsgCh chan modell.ClientMsg
	exitCh      chan struct{}

	ws    *wsController
	term  *termController
	state *stateController
}

// NewApp returns a pointer to an initialized app object with the given clientID
func newApp(ID string) *app {
	return &app{
		clientID:    ID,
		keyCh:       make(chan string, 8),
		stateCh:     make(chan string, 8),
		errorCh:     make(chan error, 1),
		clientMsgCh: make(chan modell.ClientMsg, 8),
		exitCh:      make(chan struct{}, 1),
	}
}

// setup ...
func (a *app) setup(username string) {

	go a.errorReader()

	a.ws = newWsController(a.clientID, a.clientMsgCh, a.stateCh, a.errorCh, a.exitCh)
	a.term = newTermController(a.keyCh, a.errorCh, a.exitCh)
	a.state = newStateController(a.stateCh, a.errorCh, a.exitCh)
	a.ws.initConn(utils.GetWSURL(viper.GetString("SNAKE_URL"), a.clientID, viper.GetString("SNAKE_SECRET")))

	a.term.init()

	snakeStyle := a.selectSnake()

	a.login(username, snakeStyle)

	a.ws.initStream()

	a.state.init()

}

// login sends the login data to the Snake-hub server and waits for its response
// any erroor will be writen to the errorCh chanel
func (a *app) login(username string, snakeStyle modell.SnakeStyle) {

	log.WithField("Username", username).Debug("Logging into Snake-hub server")

	loginData := modell.LoginData{
		UserName:   username,
		SnakeStyle: snakeStyle,
	}

	data, err := json.Marshal(loginData)
	if err != nil {
		a.errorCh <- errors.Wrap(err, "Cannot encode login data")
		return
	}

	a.ws.post("login", string(data))

	resp := a.ws.read()

	log.WithField("Server Msg", resp).Debug("Login response")

	if strings.Contains(resp.Data, "Error") {
		log.Info("Cannot join the game")
		a.errorCh <- errors.Errorf("Cannot join the game %s", resp.Data)
		return
	}

	log.Debug("Successfully logged into the Snake-hub server")
}

func (a *app) selectSnake() modell.SnakeStyle {
	return modell.SnakeStyle{
		Color:     termbox.ColorCyan,
		BgColor:   termbox.ColorWhite,
		HeadRune:  '🤩',
		LeftRune:  '(',
		RightRune: ')',
	}
}

func (a *app) run() {

}

// errorReader will constantly read the error channel and call Exit accordingly
func (a *app) errorReader() {
	for {
		select {
		case <-a.exitCh:
			return
		case err := <-a.errorCh:
			if strings.Contains(err.Error(), "Terminated by user") {
				a.exit(err.Error(), 0)
			} else {
				a.exit(err.Error(), 1)
			}
		}
	}
}

// exit will log the exit message, close the exitCh channel,
// close the Termbox controller, wait 200ms
// and call os.exit with the given status code
func (a *app) exit(msg string, statusCode int) {
	close(a.exitCh)
	time.Sleep(time.Millisecond * 100)
	if statusCode == 0 {
		log.Info(msg)
	} else {
		log.Error(msg)
	}
	if a.term.isOpen() {
		log.Info("Closing Termbox")
		a.term.clear()
		a.term.close()
	}
	if a.ws.isConnected() {
		a.ws.close()
	}
	time.Sleep(time.Millisecond * 100)
	if statusCode > 0 {
		fmt.Println("An error occured during execution, check snake-hub-client.log for more detail")
	}
	os.Exit(statusCode)
}
