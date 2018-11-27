package modell

import (
	"github.com/nsf/termbox-go"

	"github.com/mikloslorinczi/snake-hub/utils"
)

// import "github.com/rs/xid"

// Direction represents a movement relative to x y
type Direction struct {
	// VX X Velocity
	VX int
	// VY Y Velocity
	VY int
}

var (
	// Up VX 0 VY -1
	Up = Direction{VX: 0, VY: -1}
	// Down VX 0 VY 1
	Down = Direction{VX: 0, VY: 1}
	// Left VX -1 VY 0
	Left = Direction{VX: -1, VY: 0}
	// Right VX 1 VY 0
	Right = Direction{VX: 1, VY: 0}
)

// Opposite returns true if new direction is the opposite of current direction
func (current Direction) Opposite(new Direction) bool {
	return current.VX+new.VX == 0 && current.VY+new.VY == 0
}

// Snake represents a snake-object
type Snake struct {
	ID           string
	PlayerID     string
	Color        termbox.Attribute
	Body         []Block
	Direction    Direction
	TargetLength int
}

// NewSnake creates a snake with the given pharameters generates an ID for it and returns a pointer to it
func NewSnake(playerID string, x, y int, color termbox.Attribute, direction Direction) *Snake {

	return &Snake{
		ID:       utils.NewID(),
		PlayerID: playerID,
		Color:    color,
		Body: []Block{
			{
				X:          x,
				Y:          y,
				Color:      color,
				Background: color,
				LeftRune:   ' ',
				RightRune:  ' ',
			},
			{
				X:          x - 1,
				Y:          y,
				Color:      color,
				Background: color,
				LeftRune:   ' ',
				RightRune:  ' ',
			},
			{
				X:          x - 2,
				Y:          y,
				Color:      color,
				Background: color,
				LeftRune:   ' ',
				RightRune:  ' ',
			}},
		Direction:    direction,
		TargetLength: 3,
	}
}
