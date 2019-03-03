package client

// import (
// 	"encoding/json"
// 	"fmt"

// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	termbox "github.com/nsf/termbox-go"

// 	log "github.com/sirupsen/logrus"

// 	"github.com/pkg/errors"
// 	"github.com/spf13/viper"

// 	"github.com/mikloslorinczi/snake-hub/modell"
// 	"github.com/mikloslorinczi/snake-hub/utils"
// 	"github.com/mikloslorinczi/snake-hub/validator"
// )

// // Run starts the Client...
// func Runer() {

// 	state.userName = getUsername()

// 	// Setup Logger
// 	log.SetFormatter(&log.TextFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05",
// 		FullTimestamp:   true,
// 	})

// 	if viper.GetBool("SNAKE_DEBUG") {
// 		log.SetLevel(log.DebugLevel)
// 	} else {
// 		log.SetLevel(log.InfoLevel)
// 	}

// 	logFile, err := os.OpenFile("snake-hub-client.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer logFile.Close()
// 	log.SetOutput(logFile)

// 	log.Info("Initializing Snake Client")

// 	// Connect to Snake-hub server
// 	log.Info("Connecting to Snake-Hub")
// 	if err := getConn(); err != nil {
// 		gracefulStop(fmt.Sprintf("Cannot connect to the Snake-hub server %s", err), 1)
// 	}
// 	connOpen = true

// 	// Set up Termbox
// 	log.Info("Starting Termbox Controller")
// 	tc = NewtermController(keyCh, errorCh, exitCh)
// 	tc.Init()

// 	// Select Snake-style
// 	/* 	log.Info("Starting SnakeSelector")
// 	   if err := snakeSelector(); err != nil {
// 		   gracefulStop(fmt.Sprintf("Error during Snake Selection %s", err), 1)
// 	   } */

// 	// Join the game
// 	log.Info("Joining game")
// 	if err := login(); err != nil {
// 		gracefulStop(fmt.Sprintf("Snake-hub refused the connection %s", err), 1)
// 	}

// 	// go func() {
// 	// 	for {
// 	// 		select {
// 	// 		case key := <-keyCh:
// 	// 			if key == "esc" {
// 	// 				gracefulStop("Byez", 0)
// 	// 			}
// 	// 		case <-exitCh:
// 	// 			return
// 	// 		}
// 	// 	}
// 	// }()

// 	log.Info("Starting WebSocket Controller")
// 	go ws.startReader()
// 	go ws.startWriter()

// 	log.Info("Entering mainloop")
// 	returnMsg := "gg bb"
// 	returnCode := 0

// mainLoop:
// 	for {
// 		select {
// 		case err := <-errorCh:
// 			log.WithField("Msg", err).Debug("Msg Received on errorCh")
// 			errStr := fmt.Sprintf("Error : %s", err)
// 			close(exitCh)
// 			if !strings.Contains(errStr, "Terminated by user") {
// 				returnCode = 1
// 			}
// 			returnMsg = errStr
// 			break mainLoop
// 		}
// 	}

// 	gracefulStop(returnMsg, returnCode)

// }

// func getUsernames() string {
// 	if confName := viper.GetString("SNAKE_USERNAME"); validator.ValidUsername(confName) {
// 		return confName
// 	}
// 	for {
// 		inputName, err := utils.GetInput("Enter your username (max 8 char) :")
// 		if err != nil {
// 			log.Fatalf("Error reading username %v", err)
// 		}
// 		if validator.ValidUsername(inputName) {
// 			return inputName
// 		}
// 		fmt.Printf("Invalid user name %s\n", inputName)
// 	}
// }

// func getConn() error {
// 	// Get the URL of Snake-hub server, send client ID and Snake Secret as query string
// 	wsURL := utils.GetWSURL(viper.GetString("SNAKE_URL"), clientID, viper.GetString("SNAKE_SECRET"))

// 	log.WithField("Connection string", wsURL).Info("Connecting to WebSocket server")

// 	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

// 	if err != nil {
// 		return errors.Wrapf(err, "Cannot dial %s", wsURL)
// 	}

// 	conn = wsConn
// 	return nil

// }

// func login() error {

// 	log.WithField("Client ID", clientID).Info("Joining game with client ID")

// 	state.userSnakestyle = modell.SnakeStyle{
// 		Color:     termbox.ColorYellow,
// 		BgColor:   termbox.ColorBlack,
// 		HeadRune:  modell.GetRandomHead(),
// 		LeftRune:  '<',
// 		RightRune: '>',
// 	}

// 	bytes, err := json.Marshal(modell.LoginData{
// 		UserName:   state.userName,
// 		SnakeStyle: state.userSnakestyle,
// 	})

// 	if err != nil {
// 		return errors.Wrap(err, "Cannot encode login data")
// 	}

// 	msg := modell.ClientMsg{
// 		ClientID: clientID,
// 		Type:     "login",
// 		Data:     string(bytes),
// 	}

// 	if err := conn.WriteJSON(msg); err != nil {
// 		return errors.Wrap(err, "Cannot write to WebSocket")
// 	}

// 	resp := modell.ServerMsg{}
// 	if err := conn.ReadJSON(&resp); err != nil {
// 		return errors.Wrap(err, "Cannot read from WebSocket")
// 	}

// 	log.WithField("Server Msg", resp).Debug("Login response")

// 	if strings.Contains(resp.Data, "Error") {
// 		log.Info("Cannot join the game")
// 		return errors.New(resp.Data)
// 	}

// 	log.Info("Successfully joined the game")
// 	return nil

// }

// // Log close message, stop WebSocket, stop Termbox, return status code
// func gracefulStop(msg string, returnCode int) {

// 	log.Info(msg)

// 	if tc.IsOpen {
// 		log.Info("Closing Termbox")
// 		if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
// 			log.WithField("Error", err).Error("Cannot clear Termbox")
// 		}
// 		termbox.Close()
// 		time.Sleep(time.Second / 3)
// 	}

// 	if connOpen {
// 		log.Info("Closing WebSocket Connection")
// 		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
// 			log.WithField("Error", err).Error("Cannot send close message on WebSocket")
// 		}
// 		time.Sleep(time.Second / 3)

// 		log.Info("Closing Websocket")
// 		if err := conn.Close(); err != nil {
// 			log.WithField("Error", err).Error("Cannot close WebSocket connection properly")
// 		}
// 		time.Sleep(time.Second / 3)
// 	}

// 	if returnCode != 0 {
// 		fmt.Printf("\nAn error occured during execution, check snake-hub-client.log for details\n")
// 	}

// 	os.Exit(returnCode)
// }
