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

// newWsController returns a pointer to an un-initialized wsController object
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

// initConn tries to connect to the given websocket url
func (ws *wsController) initConn(url string) {
	log.WithField("Connection string", url).Info("Connecting to WebSocket server")
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		ws.errorCh <- errors.Wrapf(err, "Cannot dial Snake-hub at %s", url)
		return
	}
	ws.conn = wsConn
	ws.connected = true
	log.WithField("Connection string", url).Debug("Successfully connected to WebSocket server")
}

// isConnected reports if wsController is connected to a WebSocket server
func (ws *wsController) isConnected() bool {
	return ws.connected
}

// close sends a close message on the websocket and tries to close it properly
func (ws *wsController) close() {
	log.Info("Closing WebSocket Connection")
	if err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.WithField("Error", err).Error("Cannot send close message on WebSocket")
	}
	log.Info("Closing Websocket")
	if err := ws.conn.Close(); err != nil {
		log.WithField("Error", err).Error("Cannot close WebSocket connection properly")
	}
}

func (ws *wsController) initStream() {
	log.Debug("WebSocket controller starting reader")
	go ws.reader()
	log.Debug("WebSocket controller starting writer")
	go ws.writer()
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

// enquePost writes creates a client message with the given type and data
// and sends it to the clientMsgCh channel
func (ws *wsController) enquePost(msgType, msgData string) {
	clientMsg := modell.ClientMsg{
		ClientID: ws.clientID,
		Type:     msgType,
		Data:     msgData,
	}
	ws.clientMsgCh <- clientMsg
}

// post writes one client message to the websocket
// use this before stream is initialized
func (ws *wsController) post(msgType, msgData string) {
	clientMsg := modell.ClientMsg{
		ClientID: ws.clientID,
		Type:     msgType,
		Data:     msgData,
	}
	if err := ws.conn.WriteJSON(clientMsg); err != nil {
		ws.errorCh <- err
	}
}

// read reads the next server message from the websocket
// use this before stream is initialized
func (ws *wsController) read() modell.ServerMsg {
	msg := modell.ServerMsg{}
	if err := ws.conn.ReadJSON(&msg); err != nil {
		ws.errorCh <- err
	}
	return msg
}
