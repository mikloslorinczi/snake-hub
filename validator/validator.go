package validator

import (
	"encoding/json"

	"github.com/mikloslorinczi/snake-hub/modell"
)

// ValidUsername will check if username is etleast 1 char long
// and not longer than 8
func ValidUsername(username string) bool {
	return len(username) > 0 && len(username) < 9
}

// ValidLogin ...
func ValidLogin(data []byte, id string) bool {
	var loginData modell.LoginData
	err := json.Unmarshal(data, &loginData)
	if err != nil {
		return false
	}
	if !ValidUsername(loginData.UserName) {
		return false
	}
	return true
}
