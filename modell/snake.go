package modell

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

// Snake represents a snake-object
type Snake struct {
	ID            string            `json:"id"`
	UserID        string            `json:"userid"`
	Color         termbox.Attribute `json:"color"`
	BgColor       termbox.Attribute `json:"bgcolor"`
	HeadRune      rune              `json:"headrune"`
	LeftRune      rune              `json:"leftrune"`
	RightRune     rune              `json:"rightrune"`
	Body          []Coords          `json:"body"`
	Direction     Direction         `json:"direction"`
	NextDirection Direction         `json:"nextdirection"`
	TargetLength  int               `json:"targetlength"`
}

// GetHeadCoords return the Coords of the first (0th) element of the body
func (snake *Snake) GetHeadCoords() Coords {
	if len(snake.Body) > 0 {
		return snake.Body[0]
	}
	return Coords{}
}

// Move moves the snake
func (snake *Snake) Move(lvl LevelMap) {
	head := snake.GetHeadCoords()
	newCoords := lvl.GetCoords(head.X+snake.Direction.VX, head.Y+snake.Direction.VY)
	if len(snake.Body) < snake.TargetLength {
		snake.Body = append([]Coords{newCoords}, snake.Body...)
	} else {
		snake.Body = append([]Coords{newCoords}, snake.Body[:len(snake.Body)-1]...)
	}
}

// StringToDirection converts a string to a Direction
func StringToDirection(s string) Direction {
	switch {
	case s == "up":
		return Up
	case s == "down":
		return Down
	case s == "left":
		return Left
	case s == "right":
		return Right
	}
	return Direction{}
}

// ChangeDirection updates snake's next direction
func (snake *Snake) ChangeDirection(direction string) {
	newDirection := StringToDirection(direction)
	if !snake.Direction.IsOpposite(newDirection) {
		snake.NextDirection = newDirection
	}
}

// Direction represents a movement relative to x y
type Direction struct {
	// VX X Velocity
	VX int `json:"vx"`
	// VY Y Velocity
	VY int `json:"vy"`
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
	// Directions is a slice of all directions
	Directions = []Direction{Up, Down, Left, Right}
)

// IsOpposite returns true if new direction is the opposite of current direction
func (current Direction) IsOpposite(new Direction) bool {
	return current.VX+new.VX == 0 && current.VY+new.VY == 0
}

// Opposite returns the opposite of the current direction
func (current Direction) Opposite() Direction {
	switch current {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return current
}

// RandomDirection returns a direction choosen at random
func RandomDirection() Direction {
	return Directions[rand.Intn(len(Directions))]
}

// ClaculateSnakeBody will return a slice of Blocks based on given coordinates and direction
func ClaculateSnakeBody(x, y, length int, direction Direction, lvl LevelMap) []Coords {
	body := []Coords{}
	for i := 0; i < length; i++ {
		body = append(body, lvl.GetCoords(x+direction.Opposite().VX*i, y+direction.Opposite().VY*i))
	}
	return body
}
