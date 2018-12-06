package server

import (
	"fmt"
	"log"
	"sync"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/state"
	"github.com/tidwall/gjson"

	"github.com/gorilla/websocket"
)

var (
	clients      []*clientProcessor
	clientsMutex sync.Mutex
	controlChan  = make(chan modell.ClientMsg, 100)
	errorChan    = make(chan error, 10)
	exitChan     = make(chan struct{}, 1)
	gameControl  = state.NewGame(40, 20)
)

type clientProcessor struct {
	conn      *websocket.Conn
	closeChan chan struct{}
	user      modell.User
}

func (cp *clientProcessor) init() {

	defer func() {
		if err := cp.conn.Close(); err != nil {
			fmt.Printf("Cannot close WebSocket properly %s\n", err)
		}
	}()

wsLoop:
	for {

		mt, message, err := cp.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %s\n", err)
			break wsLoop
		}
		log.Printf("Message type: %d, Message: %s", mt, message)

		result := gjson.GetBytes(message, "type").String()
		switch {
		case result == "handshake":
			{
				fmt.Println("Handshake")
				resp := modell.ServerMsg{
					Type: "handshake",
					Data: "Accepted",
				}
				cp.conn.WriteJSON(resp)
			}
		// case result == "leave":
		// 	{
		// 		fmt.Println("User left...")
		// 		break wsLoop
		// 	}
		default:
		}

	}

	fmt.Println("User disconnected...")

}

func newClient(id string, conn *websocket.Conn) {
	cp := &clientProcessor{
		closeChan: exitChan,
		conn:      conn,
		user: modell.User{
			ID: id,
		},
	}
	cp.init()
	gameControl.AddUser(cp.user)
	ok, foundUser := gameControl.GetUser(id)
	fmt.Printf("Found %v, %+v\n", ok, foundUser)
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients = append(clients, cp)
}

func removeClient(id string) bool {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for i, client := range clients {
		if client.user.ID == id {
			clients = append(clients[:i], clients[i+1:]...)
			return true
		}
	}
	return false
}
