package modell

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

// Direction represents a movement relative to x y
type Direction struct {
	// VX X Velocity
	VX int `json:"vx"`
	// VY Y Velocity
	VY int `json:"vy"`
}

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
	Speed         int               `json:"speed"`
	StepSize      int               `json:"stepsize"`
	Alive         bool              `josn:"alive"`
}

type snakeTexture struct {
	leftRune  rune
	rightRune rune
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
	directions = []Direction{Up, Down, Left, Right}
	// Heads is a slice of all available head rune
	heads = []rune{'🐸', '😗', '😡', '🤢', '😈', '💀', '🤖', '😸', '👽', '🐷'}
	// Textures is a slice of all available snake texture
	textures = []snakeTexture{
		{
			leftRune:  '[',
			rightRune: ']',
		},
		{
			leftRune:  '(',
			rightRune: ')',
		},
		{
			leftRune:  '<',
			rightRune: '>',
		},
		{
			leftRune:  'o',
			rightRune: 'O',
		},
		{
			leftRune:  '~',
			rightRune: '~',
		},
	}
)

// GetRandomHead returns a random head rune from all available head runes
func GetRandomHead() rune {
	return heads[rand.Intn(len(heads))]
}

// GetRandomTexture returns the left and right rune of a random texture
func GetRandomTexture() (rune, rune) {
	randomTexture := textures[rand.Intn(len(textures))]
	return randomTexture.leftRune, randomTexture.rightRune
}

// RandomDirection returns a direction choosen at random
func RandomDirection() Direction {
	return directions[rand.Intn(len(directions))]
}

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

// GetHeadCoords return the Coords of the first (0th) element of the body
func (snake *Snake) GetHeadCoords() Coords {
	if len(snake.Body) > 0 {
		return snake.Body[0]
	}
	return Coords{}
}

// Update moves the snake, and changes direction to nextDirection
func (snake *Snake) Update(lvl LevelMap) {

	snake.StepSize += snake.Speed

	if snake.StepSize >= 12 {

		snake.StepSize = 0

		head := snake.GetHeadCoords()

		snake.Direction = snake.NextDirection

		newCoords := lvl.GetCoords(head.X+snake.Direction.VX, head.Y+snake.Direction.VY)

		if len(snake.Body) < snake.TargetLength { // If the snake is still growing, just append the body
			snake.Body = append([]Coords{newCoords}, snake.Body...)
		} else {
			snake.Body = append([]Coords{newCoords}, snake.Body[:len(snake.Body)-1]...) // If it reached the target length, cut the last block
		}

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

// ClaculateSnakeBody will return a slice of Blocks based on given coordinates and direction
func ClaculateSnakeBody(x, y, length int, direction Direction, lvl LevelMap) []Coords {
	body := []Coords{}
	for i := 0; i < length; i++ {
		body = append(body, lvl.GetCoords(x+direction.Opposite().VX*i, y+direction.Opposite().VY*i))
	}
	return body
}
