package server

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/websocket"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/validator"
)

type clientController struct {
	conn          *websocket.Conn
	killChan      chan struct{}
	serverMsgChan chan modell.ServerMsg
	clientErrChan chan error
	userID        string
}

func (client *clientController) init() {

	log.WithField("User ID", client.userID).Debug("Inicializing client Controller")

	defer func() {
		if err := client.conn.Close(); err != nil {
			log.WithField("Error", err).Error("Cannot close WebSocket properly")
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
				client.sendServerMsg(serverMsg.Type, serverMsg.Data)
			}
		case <-client.killChan:
			{
				log.WithField("User ID", client.userID).Debug("Stop writing server messages onto the WebScoket")
				return
			}
		}
	}
}

func (client *clientController) sendServerMsg(msgType, msgData string) {
	msg := modell.ServerMsg{
		Type: msgType,
		Data: msgData,
	}
	if err := client.conn.WriteJSON(msg); err != nil {
		client.clientErrChan <- err
	}
}

func (client *clientController) handleMsg(msg modell.ClientMsg) {
	log.WithFields(log.Fields{
		"User ID": msg.ClientID,
		"Type":    msg.Type,
		"Body":    msg.Data,
	}).Debug("Incoming WS Messgae")
	switch {
	case msg.Type == "login":
		{
			client.handleLogin(msg)
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

func (client *clientController) handleLogin(msg modell.ClientMsg) {
	if len(gameState.state.Users) == viper.GetInt("SNAKE_MAX_PLAYER") {
		client.sendServerMsg("login", fmt.Sprintf("Error: Server is full (max player %d)", viper.GetInt("SNAKE_MAX_PLAYER")))
		return
	}
	if !validator.ValidLogin([]byte(msg.Data), client.userID) {
		client.sendServerMsg("login", fmt.Sprintf("Error: Invalid login data %s", msg.Data))
		return
	}
	userData := modell.LoginData{}
	json.Unmarshal([]byte(msg.Data), &userData)
	client.sendServerMsg("login", fmt.Sprintf("User successfully loged in with UserName: %s and UserID %s", userData.UserName, client.userID))

	var loginData modell.LoginData
	json.Unmarshal([]byte(msg.Data), &loginData)

	go func(data modell.LoginData) {
		gameState.mu.Lock()
		defer gameState.mu.Unlock()
		gameState.state.AddUser(modell.User{
			Name:       data.UserName,
			ID:         client.userID,
			Score:      0,
			SnakeStyle: data.SnakeStyle,
		})
	}(loginData)

}
