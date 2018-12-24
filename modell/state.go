package modell

import (
	"github.com/nsf/termbox-go"
)

// State represents the game-state
// It can marshalled to JSON and sent over WebSockets
type State struct {
	Snakes []Snake  `json:"snakes"`
	Users  []User   `json:"users"`
	Level  LevelMap `json:"level"`
}

// NewState creates a new game with an empty map of width * height
func NewState(width, height int, bgColor termbox.Attribute) State {
	return State{
		Level: LevelMap{
			Width:   width,
			Height:  height,
			BgColor: bgColor,
		},
	}
}
