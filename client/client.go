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
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/utils"
)

var (
	clientID      = utils.NewID()
	conn          *websocket.Conn
	connMutex     sync.RWMutex
	connOpen      = false
	termboxOpen   = false
	exitChan      = make(chan struct{}, 1)
	errorChan     = make(chan error, 10)
	termChan      = make(chan string, 10)
	serverMsgChan = make(chan modell.ServerMsg, 10)
	clientMsgChan = make(chan modell.ClientMsg, 10)
)

// Run starts the Client...
func Run() {

	fmt.Println("Initializing Snake Client...")

	fmt.Println("Connecting to Snake-Hub")
	if err := getConn(); err != nil {
		gracefulStop(fmt.Sprintf("Cannot connect to the Snake-hub server %s\n", err), 1)
	}
	connOpen = true

	fmt.Println("Joining game...")
	if err := login(); err != nil {
		gracefulStop(fmt.Sprintf("Snake-hub refused the connection %s\n", err), 1)
	}

	fmt.Println("Initializeing TermBox...")
	err := termbox.Init()
	if err != nil {
		gracefulStop(fmt.Sprintf("Termbox init error %v\n", err), 1)
	}
	termboxOpen = true

	fmt.Println("Starting WebSocket Handler...")
	ws := &wsHandler{
		stopReader: exitChan,
		stopWriter: exitChan,
	}
	go ws.startReader()
	go ws.startWriter()

	fmt.Println("Starting Termbox loop...")
	tw := &termboxWorker{
		stopEventLoop: exitChan,
		stopRender:    exitChan,
	}
	go tw.startEventloop()
	go tw.startRenderer()

	fmt.Println("Entering mainloop")
	returnMsg := ""
	returnCode := 0

mainLoop:
	for {
		select {
		case err := <-errorChan:
			errStr := fmt.Sprintf("Error : %s\n", err)
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
	// Get the URL of Snake-hub server
	u := url.URL{Scheme: "ws", Host: viper.GetString("SNAKE_URL"), Path: "/game"}
	wsURL := u.String()
	fmt.Printf("Connecting to %v...\n", wsURL)

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return errors.Wrapf(err, "Cannot dial %s", wsURL)
	}

	conn = wsConn
	return nil

}

func login() error {

	fmt.Printf("Joining game with client ID: %s\n", clientID)
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

	fmt.Printf("Hanshake response %v\n", resp)
	return nil

}

// Print close message, stop WebSocket, stop Termbox, return status code
func gracefulStop(msg string, returnCode int) {

	fmt.Println(msg)

	if connOpen {

		fmt.Println("Closing WebSocket Connection...")
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			fmt.Printf("Cannot send close message on WebSocket %s\n", err)
		}
		time.Sleep(time.Second / 3)

		fmt.Println("Closing Websocket...")
		if err := conn.Close(); err != nil {
			fmt.Printf("Cannot close WebSocket connection properly %s\n", err)
		}
		time.Sleep(time.Second / 3)

	}

	if termboxOpen {

		fmt.Println("Closing Termbox...")
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Close()
		time.Sleep(time.Second / 3)

	}

	os.Exit(returnCode)

}
