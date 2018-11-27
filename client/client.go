package client

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	termbox "github.com/nsf/termbox-go"
	"github.com/spf13/viper"
)

var (
	eventQueue = make(chan termbox.Event)
	interrupt  = make(chan os.Signal, 1)
	running    = false
)

func getConn() *websocket.Conn {
	// Get the URL of Snake-hub server
	u := url.URL{Scheme: "ws", Host: viper.GetString("SNAKE_URL"), Path: "/echo"}
	wsURL := u.String()
	fmt.Printf("Connecting to %v\n", wsURL)
	// Connect to the Snake-hub server, fatal on error
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return conn
}

// Run starts the Client...
func Run() {

	running = true

	// Get WebSocket connection, defer close
	conn := getConn()
	defer conn.Close()

	// watch signals (such as Ctrl+C ^C)
	go watchSignal(conn)

	// Start the message reader in the background
	go wsReader(conn)

	// Start Termbox event reader in the background
	go eventReader(conn)

	// Start the rendering process
	go render()

	// Block while running
	for running {
	}
	time.Sleep(time.Second)
	fmt.Println("Halted")
}
