package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var upgrader = websocket.Upgrader{}

func hub(res http.ResponseWriter, req *http.Request) {

	id, err := auth(req)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Unauthorized\n"))
		return
	}

	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		res.WriteHeader(http.StatusUpgradeRequired)
		res.Write([]byte(fmt.Sprintf("Cannot upgrade HTTP connection to WebSocket %s\n", err)))
		return
	}

	newClient(id, conn)

}

func auth(req *http.Request) (string, error) {

	query := req.URL.Query()
	foundSecret := false
	foundID := false
	id := ""

	for key, value := range query {
		switch key {
		case "snakesecret":
			{
				if len(value) == 1 && value[0] == viper.GetString("SNAKE_SECRET") {
					foundSecret = true
				}
			}
		case "clientid":
			{
				if len(value) == 1 {
					foundID = true
					id = value[0]
				}
			}
		}
	}

	if foundID && foundSecret {
		return id, nil
	}

	return id, errors.New("Secret or ID missing")

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "")
}
