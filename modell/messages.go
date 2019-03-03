package modell

// ServerMsg is a JSON formatted message sent by the Snake-hub server to a client
type ServerMsg struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// ClientMsg is a JSON formatted message sent by a client to the Snake-hub server
type ClientMsg struct {
	ClientID string `json:"clientid"`
	Type     string `json:"type"`
	Data     string `json:"data"`
}

// LoginData sent by the client as ClientMsg.Data when logging in
type LoginData struct {
	UserName   string     `json:"username"`
	SnakeStyle SnakeStyle `json:"snakestyle"`
}
