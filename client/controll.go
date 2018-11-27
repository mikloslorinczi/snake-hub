package client

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	termbox "github.com/nsf/termbox-go"
)

func watchSignal(conn *websocket.Conn) {
	// pipe os.Signals such as ctrl+c ^C to to the interrupt channel
	signal.Notify(interrupt, os.Interrupt)
	fmt.Printf("Watching OS signal...\n")
	for running {
		select {
		case <-interrupt: // on os.Signal
			fmt.Printf("Interrupted...\n")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Printf("Error writing close message to WebSocket: %s", err)
			}
			running = false
		}
	}
}

func wsReader(conn *websocket.Conn) {
	// Log every message received on the WebsScket in the background
	fmt.Printf("Reading messages on webSocket...\n")
	for running {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WebSocket read error %s\n", err)
			running = false
		} else {
			fmt.Printf("Message received: %s\n", message)
		}
	}
}

func eventReader(conn *websocket.Conn) {
	fmt.Printf("Reading Termbox events...\n")
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	for running {
		select {
		case event := <-eventQueue:
			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyEsc:
					running = false
				case termbox.KeyEnter:
					conn.WriteMessage(websocket.TextMessage, []byte("Enter pressed\n"))
				case termbox.KeyArrowDown:
					conn.WriteMessage(websocket.TextMessage, []byte("Enter pressed\n"))
				}
			}
		}
	}
}
