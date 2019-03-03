package client

import (
	"github.com/gorilla/websocket"
	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// wsController is a wrapper object around WebSocket functionality
type wsController struct {
	clientMsgCh chan modell.ClientMsg
	stateCh     chan string
	errorCh     chan error
	stopCh      chan struct{}

	clientID  string
	connected bool
	conn      *websocket.Conn
}

// newWsController returns a pointer to an initialized wsController object
func newWsController(ID string, ClientMsgCh chan modell.ClientMsg, StateCh chan string, ErrorCh chan error, StopCh chan struct{}) *wsController {
	return &wsController{
		clientMsgCh: ClientMsgCh,
		stateCh:     StateCh,
		errorCh:     ErrorCh,
		stopCh:      StopCh,
		clientID:    ID,
		connected:   false,
	}
}

// init tries to connect to the given websocket url
func (ws *wsController) init(url string) {
	log.WithField("Connection string", url).Info("Connecting to WebSocket server")
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		ws.errorCh <- errors.Wrapf(err, "Cannot dial Snake-hub at %s", url)
		return
	}
	ws.conn = wsConn
	go ws.reader()
	go ws.writer()
	ws.connected = true
	log.WithField("Connection string", url).Debug("Successfully connected to WebSocket server")
}

// isConnected reports id wsController is connected to a WebSocket server
func (ws *wsController) isConnected() bool {
	return ws.connected
}

// reader will constantly try to read server messages from the websocket
// and depending on type feeding the data to stateCh or ...
func (ws *wsController) reader() {
	for {
		select {
		case <-ws.stopCh:
			log.Debug("WebSocket Controller reader stopped")
			return
		default:
			msg := modell.ServerMsg{}
			if err := ws.conn.ReadJSON(&msg); err != nil {
				ws.errorCh <- errors.Wrap(err, "Cannot read server message from websocket")
			}
			switch msg.Type {
			case "broadcastMsg":
				log.WithField("Msg", msg.Data).Info("Server broadcast message")
			case "stateUpdate":
				ws.stateCh <- msg.Data
			default:
				log.WithFields(log.Fields{
					"Type": msg.Type,
					"Msg":  msg.Data,
				}).Warn("Unknow message type")
			}
		}
	}
}

// writer will write every message from the clientMsgCh to the websocket
func (ws *wsController) writer() {
	for {
		select {
		case <-ws.stopCh:
			log.Debug("WebSocket Controller writer stopped")
			return
		case msg := <-ws.clientMsgCh:
			if err := ws.conn.WriteJSON(msg); err != nil {
				ws.errorCh <- errors.Wrap(err, "Cannot write to websocket")
			}
		}
	}
}

// post writes creates a client message with the given type and data
// and sends it to the clientMsgCh channel
func (ws *wsController) post(msgType, msgData string) {

	clientMsg := modell.ClientMsg{
		ClientID: ws.clientID,
		Type:     msgType,
		Data:     msgData,
	}

	ws.clientMsgCh <- clientMsg

}
