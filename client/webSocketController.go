package client

import (
	"fmt"

	"github.com/mikloslorinczi/snake-hub/modell"
)

type wsHandler struct {
	stopReader chan struct{}
	stopWriter chan struct{}
}

func (ws *wsHandler) startReader() {
	for {
		select {
		case <-ws.stopReader:
			return
		default:
			resp := modell.ServerMsg{}
			if err := conn.ReadJSON(&resp); err != nil {
				errorChan <- err
			}
			switch resp.Type {
			case "broadcastMsg":
				fmt.Printf("Server broadcast message: %s", resp.Data)
			}
		}
	}
}

func (ws *wsHandler) startWriter() {
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
