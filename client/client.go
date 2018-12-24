package client

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	termbox "github.com/nsf/termbox-go"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/utils"
)

var (
	clientID    = utils.NewID()
	conn        *websocket.Conn
	connMutex   sync.RWMutex
	connOpen    = false
	termboxOpen = false

	exitChan      = make(chan struct{}, 1)
	errorChan     = make(chan error, 10)
	termChan      = make(chan string, 10)
	serverMsgChan = make(chan modell.ServerMsg, 10)
	clientMsgChan = make(chan modell.ClientMsg, 10)

	ws = &wsController{
		stopReader: exitChan,
		stopWriter: exitChan,
	}
	tc = &termboxController{
		stopEventLoop: exitChan,
		stopRender:    exitChan,
	}
	state = &stateController{
		loaded: false,
	}
)

// Run starts the Client...
func Run() {

	// Log formatter
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	log.SetFormatter(customFormatter)

	// Log level
	if viper.GetBool("SNAKE_DEBUG") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Log file
	logFile, err := os.OpenFile("snake-hub-client.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Info("Initializing Snake Client...")

	log.Info("Connecting to Snake-Hub")
	if err := getConn(); err != nil {
		gracefulStop(fmt.Sprintf("Cannot connect to the Snake-hub server %s", err), 1)
	}
	connOpen = true

	log.Info("Joining game...")
	if err := login(); err != nil {
		gracefulStop(fmt.Sprintf("Snake-hub refused the connection %s", err), 1)
	}

	log.Info("Initializeing TermBox...")
	if err := termbox.Init(); err != nil {
		gracefulStop(fmt.Sprintf("Termbox init error %v", err), 1)
	}
	termboxOpen = true

	log.Info("Starting WebSocket Controller...")
	go ws.startReader()
	go ws.startWriter()

	log.Info("Starting Termbox Controller...")
	tc.getSize()
	go tc.startEventloop()
	go tc.startRenderer()

	log.Info("Entering mainloop")
	returnMsg := "gg bb"
	returnCode := 0

mainLoop:
	for {
		select {
		case err := <-errorChan:
			log.WithField("Msg", err).Debug("Msg Received on errorChan")
			errStr := fmt.Sprintf("Error : %s", err)
			close(exitChan)
			if !strings.Contains(errStr, "Terminated by user") {
				returnCode = 1
			}
			returnMsg = errStr
			break mainLoop
		}
	}

	gracefulStop(returnMsg, returnCode)

}

func getConn() error {
	// Get the URL of Snake-hub server, send client ID and Snake Secret as query string
	u := url.URL{
		Scheme:   "ws",
		Host:     viper.GetString("SNAKE_URL"),
		Path:     "/hub",
		RawQuery: fmt.Sprintf("clientid=%s&snakesecret=%s", clientID, viper.Get("SNAKE_SECRET")),
	}
	wsURL := u.String()

	log.WithField("Connection string", wsURL).Info("Connecting to WebSocket server")

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return errors.Wrapf(err, "Cannot dial %s", wsURL)
	}

	conn = wsConn
	return nil

}

func login() error {

	log.WithField("Client ID", clientID).Info("Joining game with client ID")
	msg := modell.ClientMsg{
		ClientID: clientID,
		Type:     "handshake",
		Data:     viper.GetString("SNAKE_SECRET"),
	}
	if err := conn.WriteJSON(msg); err != nil {
		return errors.Wrap(err, "Cannot write to WebSocket")
	}

	resp := modell.ServerMsg{}
	if err := conn.ReadJSON(&resp); err != nil {
		return errors.Wrap(err, "Cannot read from WebSocket")
	}

	log.WithField("Server Msg", resp).Info("Hanshake response")
	return nil

}

// Log close message, stop WebSocket, stop Termbox, return status code
func gracefulStop(msg string, returnCode int) {

	log.Info(msg)

	if termboxOpen {

		log.Info("Closing Termbox...")
		if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
			log.WithField("Error", err).Error("Cannot clear Termbox")
		}
		termbox.Close()
		time.Sleep(time.Second / 3)

	}

	if connOpen {

		log.Info("Closing WebSocket Connection...")
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			log.WithField("Error", err).Error("Cannot send close message on WebSocket")
		}
		time.Sleep(time.Second / 3)

		log.Info("Closing Websocket...")
		if err := conn.Close(); err != nil {
			log.WithField("Error", err).Error("Cannot close WebSocket connection properly")
		}
		time.Sleep(time.Second / 3)

	}

	os.Exit(returnCode)

}
