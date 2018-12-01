package state

import (
	"github.com/mikloslorinczi/snake-hub/modell"
)

type state struct {
	snakes []modell.Snake
	users  []modell.User
	level  modell.LevelMap
}
