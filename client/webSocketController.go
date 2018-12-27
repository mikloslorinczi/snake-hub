package client

import (
	"github.com/mikloslorinczi/snake-hub/modell"
	log "github.com/sirupsen/logrus"
)

type wsController struct {
	stopReader chan struct{}
	stopWriter chan struct{}
}

func (ws *wsController) startReader() {
	for {
		msg := modell.ServerMsg{}
		if err := conn.ReadJSON(&msg); err != nil {
			errorChan <- err
		}
		switch msg.Type {
		case "broadcastMsg":
			log.WithField("Msg", msg.Data).Info("Server broadcast message")
		case "stateUpdate":
			if err := state.loadState([]byte(msg.Data)); err != nil {
				log.WithField("Msg", err).Error("Cannot load state")
			} else {
				state.loaded = true
			}
		default:
			log.WithFields(log.Fields{
				"Type": msg.Type,
				"Msg":  msg.Data,
			}).Warn("Unknow message type")
		}
	}
}

func (ws *wsController) startWriter() {
	for {
		select {
		case <-ws.stopWriter:
			return
		case msg := <-clientMsgChan:
			if err := conn.WriteJSON(msg); err != nil {
				errorChan <- err
			}
		}
	}
}

func postMsg(msgType, msgData string) {

	clientMsg := modell.ClientMsg{
		ClientID: clientID,
		Type:     msgType,
		Data:     msgData,
	}

	clientMsgChan <- clientMsg

}
