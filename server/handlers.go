package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/mikloslorinczi/snake-hub/modell"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func game(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Cannot close WebSocket properly %s\n", err)
		}
	}()

wsLoop:
	for {
		mt, message, err := conn.ReadMessage()
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
					Data: "handshaker",
				}
				conn.WriteJSON(resp)
			}
		case result == "leave":
			{
				fmt.Println("User left...")
				break wsLoop
			}
		default:
		}
		// err = conn.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
		// websocket.
		// msg := modell.CommandObj{}
		// if err := conn.ReadJSON(msg); err != nil {
		// 	log.Printf("Cannot read message %v", err)
		// 	break
	}
	fmt.Println("game loop broken..")
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/game")
}
