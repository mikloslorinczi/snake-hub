package utils

import (
	"fmt"
	"strings"
)

// GetWSURL accepts the HTTP(S) SNAKE_URL, the Client ID, and the SNAKE_SECRET and
// returns the connection-string to the proper WS endpoint
func GetWSURL(url, id, secret string) string {
	prefix := "ws://"
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	}
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
		prefix = "wss://"
	}
	return fmt.Sprintf("%s%s/hub?clientid=%s&snakesecret=%s", prefix, url, id, secret)
}
