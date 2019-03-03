package client

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pkg/errors"

	"github.com/mikloslorinczi/snake-hub/utils"
	"github.com/mikloslorinczi/snake-hub/validator"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/mikloslorinczi/snake-hub/modell"
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

// init ...
func (a *app) init(username string) {

	go a.errorReader()

	a.ws = newWsController(a.clientID, a.clientMsgCh, a.stateCh, a.errorCh, a.exitCh)
	a.ws.init(utils.GetWSURL(viper.GetString("SNAKE_URL"), a.clientID, viper.GetString("SNAKE_SECRET")))

	a.term = newTermController(a.keyCh, a.errorCh, a.exitCh)
	a.term.init()

	a.state = newStateController(username, modell.SnakeStyle{
		Color:     termbox.ColorCyan,
		BgColor:   termbox.ColorWhite,
		HeadRune:  '🤩',
		LeftRune:  '(',
		RightRune: ')',
	}, a.stateCh, a.errorCh, a.exitCh)
	a.state.init()

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
	time.Sleep(time.Millisecond * 100)
	if statusCode > 0 {
		fmt.Println("An error occured during execution, check snake-hub-client.log for more detail")
	}
	os.Exit(statusCode)
}

func getUsername() string {
	if confName := viper.GetString("SNAKE_USERNAME"); validator.ValidUsername(confName) {
		return confName
	}
	for {
		inputName, err := utils.GetInput("Enter your username (max 8 char) :")
		if err != nil {
			log.Fatalf("Error reading username %v", err)
		}
		if validator.ValidUsername(inputName) {
			return inputName
		}
		fmt.Printf("Invalid user name %s\n", inputName)
	}
}

// Run sets up and starts the client
func Run() {

	// Setup Logger
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	if viper.GetBool("SNAKE_DEBUG") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	logFile, err := os.OpenFile("snake-hub-client.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Initialize the Game
	log.Info("Initializing Snake Client")
	id := utils.NewID()
	username := getUsername()
	game := newApp(id)
	game.init(username)
	for {
		select {
		case key := <-game.keyCh:
			if key == "esc" {
				game.errorCh <- errors.New("Terminated by user")
			}
		}
	}
}
