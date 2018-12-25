package server

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/mikloslorinczi/snake-hub/modell"
)

type clientHub struct {
	stateChan chan modell.State
	closeChan chan struct{}
	clients   []*clientController
	mu        sync.Mutex
}

func (hub *clientHub) newClient(id string, conn *websocket.Conn) {
	client := &clientController{
		killChan:      make(chan struct{}, 1),
		serverMsgChan: make(chan modell.ServerMsg, 5),
		clientErrChan: make(chan error, 1),
		conn:          conn,
		userID:        id,
	}
	go client.init()
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.clients = append(hub.clients, client)
	log.WithField("User ID", id).Info("New User joined the game")
}

func (hub *clientHub) removeClient(id string) bool {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for i, client := range hub.clients {
		if client.userID == id {
			hub.clients = append(hub.clients[:i], hub.clients[i+1:]...)
			log.WithField("User ID", id).Info("User disconnected")
			return true
		}
	}
	return false
}

// Broadcast a Server Message to all clients
func (hub *clientHub) broadcast(msg modell.ServerMsg) {
	for _, client := range hub.clients {
		go func(c *clientController) {
			c.serverMsgChan <- msg
		}(client)
	}
}

func (hub *clientHub) start() {
	log.Info("Hub starting")
hubLoop:
	for {
		select {

		case state := <-hub.stateChan:
			bytes, err := json.Marshal(state)
			if err != nil {
				log.WithField("Error", err).Error("Cannot marshal JSON")
			}
			msg := &modell.ServerMsg{
				Type: "stateUpdate",
				Data: string(bytes),
			}
			hub.broadcast(*msg)

		case <-hub.closeChan:
			for _, client := range hub.clients {
				go func(c *clientController) {
					close(c.killChan)
				}(client)
			}
			break hubLoop
		}

	}
	log.Info("Hub stopped")
}
