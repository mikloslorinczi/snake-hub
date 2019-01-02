package server

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

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

	go func() {
		gameState.mu.Lock()
		defer gameState.mu.Unlock()
		gameState.state.RemoveSnake(client.userID)
		gameState.state.RemoveUser(client.userID)
	}()

	wsHub.removeClient(client.userID)
	close(client.killChan)
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

func (client *clientController) handleMsg(msg modell.ClientMsg) {
	log.WithFields(log.Fields{
		"User ID": msg.ClientID,
		"Type":    msg.Type,
		"Body":    msg.Data,
	}).Debug("Incoming WS Messgae")
	switch {
	case msg.Type == "handshake":
		{
			if len(gameState.state.Users) == viper.GetInt("SNAKE_MAX_PLAYER") {
				resp := modell.ServerMsg{
					Type: "handshake",
					Data: fmt.Sprintf("Server is full (max player %d)", viper.GetInt("SNAKE_MAX_PLAYER")),
				}
				if err := client.conn.WriteJSON(resp); err != nil {
					client.clientErrChan <- err
				}
				return
			}
			resp := modell.ServerMsg{
				Type: "handshake",
				Data: fmt.Sprintf("User successfully loged in with User ID %v", client.userID),
			}
			if err := client.conn.WriteJSON(resp); err != nil {
				client.clientErrChan <- err
				return
			}
			go func() {
				gameState.mu.Lock()
				defer gameState.mu.Unlock()
				gameState.state.AddUser(modell.User{
					ID: client.userID,
				})
			}()
		}
	case msg.Type == "control":
		{
			if gameState.getScene() == "game" {
				go gameState.changeDirection(client.userID, msg.Data)
			}
		}
	default:
		log.WithFields(log.Fields{
			"Client ID": msg.ClientID,
			"Type":      msg.Type,
			"Body":      msg.Data,
		}).Warn("Unknown message type")
	}
}
