package server

import (
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"

	"github.com/mikloslorinczi/snake-hub/modell"
)

type clientController struct {
	conn          *websocket.Conn
	killChan      chan struct{}
	serverMsgChan chan modell.ServerMsg
	clientErrChan chan error
	userID        string
}

type clientHub struct {
	stateChan chan modell.State
	closeChan chan struct{}
	clients   []*clientController
	mu        sync.Mutex
}

func (client *clientController) init() {

	log.WithField("User ID", client.userID).Debug("Inicializing client Controller...")

	defer func() {
		if err := client.conn.Close(); err != nil {
			log.Error("Cannot close WebSocket properly %s\n", err)
			client.clientErrChan <- err
		}
	}()

	go client.msgReader()

	go client.msgWriter()

clientLoop:
	for {
		select {
		case err := <-client.clientErrChan:
			{
				if err.Error() == "websocket: close 1000 (normal)" {
					log.WithFields(log.Fields{
						"User ID": client.userID,
						"Msg":     err,
					}).Debug("Client disconnected")
					break clientLoop
				}
				log.WithFields(log.Fields{
					"User ID": client.userID,
					"Msg":     err,
				}).Error("Client error")
				break clientLoop
			}
		case <-client.killChan:
			{
				log.WithField("User ID", client.userID).Debug("Client killed")
				break clientLoop
			}
		}
	}

	log.WithField("User ID", client.userID).Debug("Client connection closed")

	gameState.RemoveSnake(client.userID)
	gameState.RemoveUser(client.userID)
	wsHub.removeClient(client.userID)
	close(client.killChan)
}

func (client *clientController) newDirection(direction string) {

}

func (client *clientController) handleMsg(msg modell.ClientMsg) {
	log.WithFields(log.Fields{
		"User ID": msg.ClientID,
		"Type":    msg.Type,
		"Body":    msg.Data,
	}).Debug("Incoming WS Messgae")
	switch {
	case msg.Type == "handshake":
		{
			resp := modell.ServerMsg{
				Type: "handshake",
				Data: fmt.Sprintf("Handshake accapted from user %v", client.userID),
			}
			if err := client.conn.WriteJSON(resp); err != nil {
				client.clientErrChan <- err
				return
			}
			gameState.AddUser(modell.User{
				ID: client.userID,
			})
		}
	case msg.Type == "control":
		{
			go gameState.ChangeDirection(client.userID, msg.Data)
		}
	default:
		log.WithFields(log.Fields{
			"Client ID": msg.ClientID,
			"Type":      msg.Type,
			"Body":      msg.Data,
		}).Warn("Unknown message type")
	}
}

// msgReader will read the next message from the client and pass it to the handler
func (client *clientController) msgReader() {
	log.WithField("User ID", client.userID).Debug("Start reading client messages")
	for {
		msg := modell.ClientMsg{}
		if err := client.conn.ReadJSON(&msg); err != nil {
			client.clientErrChan <- err
			return
		}
		client.handleMsg(msg)
	}
}

// msgWriter will write every ServerMessage received on the serverMsgChan
// to the client WebSocket connection
func (client *clientController) msgWriter() {
	log.WithField("User ID", client.userID).Debug("Start write server messages")
	for {
		select {
		case serverMsg := <-client.serverMsgChan:
			{
				// log.WithField("Msg", string(serverMsg.Data)).Debug("Server Msg")
				client.conn.WriteJSON(serverMsg)
			}
		case <-client.killChan:
			{
				log.WithField("User ID", client.userID).Debug("Stop writing server messages")
				return
			}
		}
	}
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
	log.WithField("User ID", id).Info("New User joined")
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
	hub.mu.Lock()
	defer hub.mu.Unlock()
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
			hub.mu.Lock()
			defer hub.mu.Unlock()
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
