package state

import (
	"math/rand"
	"sync"

	"github.com/nsf/termbox-go"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/utils"
)

// State represents the game-state
type State struct {
	Snakes     []modell.Snake  `json:"snakes"`
	Users      []modell.User   `json:"users"`
	Level      modell.LevelMap `json:"level"`
	StateMutex sync.RWMutex
}

// Controller interface...
type Controller interface {
	GetUser(id string) (bool, *modell.User)
	AddUser(user modell.User)
	RemoveUser(id string) bool
	AddSnake(userID string)
	RemoveSnake(id string) bool
	GetNewSnake(userID string) *modell.Snake
	Update()
}

// NewGame creates a new game with an empty map of width * height
func NewGame(width, height int) Controller {
	return &State{
		Level: modell.LevelMap{
			Width:  width,
			Height: height,
		},
	}
}

// GetUser returns the user associated with the given ID
func (st *State) GetUser(id string) (bool, *modell.User) {
	for _, user := range st.Users {
		if user.ID == id {
			return true, &user
		}
	}
	return false, nil
}

// AddUser adds a new user to the game
func (st *State) AddUser(user modell.User) {
	st.StateMutex.Lock()
	defer st.StateMutex.Unlock()
	snake := st.GetNewSnake(user.ID)
	user.SnakeID = snake.ID
	st.Users = append(st.Users, user)
}

// RemoveUser removes an user from the game
func (st *State) RemoveUser(id string) bool {
	st.StateMutex.Lock()
	defer st.StateMutex.Unlock()
	for i, user := range st.Users {
		if user.ID == id {
			st.Users = append(st.Users[:i], st.Users[i+1:]...)
			return true
		}
	}
	return false
}

// AddSnake adds a new snake to the game
func (st *State) AddSnake(userID string) {
	// st.StateMutex.Lock()
	// defer st.StateMutex.Unlock()

	// snake := st.GetNewSnake(user.ID)
	// user.SnakeID = snake.ID
	// st.Users = append(st.Users, user)
}

// RemoveSnake removes a snake from the game
func (st *State) RemoveSnake(id string) bool {
	st.StateMutex.Lock()
	defer st.StateMutex.Unlock()
	for i, snake := range st.Snakes {
		if snake.ID == id {
			for _, user := range st.Users {
				if user.SnakeID == id {
					user.SnakeID = ""
				}
			}
			st.Snakes = append(st.Snakes[:i], st.Snakes[i+1:]...)
			return true
		}
	}
	return false
}

// GetNewSnake generates a new snake and returns its snake-pointer
func (st *State) GetNewSnake(userID string) *modell.Snake {

	x, y := 0, 0
	direction := modell.Up
	for {
		direction = modell.RandomDirection()
		x, y = rand.Intn(st.Level.Width), rand.Intn(st.Level.Height)
		if st.validSnakePos(x, y, direction) {
			break
		}
	}

	snake := &modell.Snake{
		ID:           utils.NewID(),
		PlayerID:     userID,
		Color:        termbox.ColorWhite,
		Body:         modell.ClaculateSnakeBody(x, y, termbox.ColorWhite, direction),
		Direction:    direction,
		TargetLength: 3,
	}

	return snake

}

func (st *State) validSnakePos(x, y int, direction modell.Direction) bool {

	boxX, boxY, boxW, boxH := 0, 0, 0, 0

	switch direction {

	case modell.Up:
		boxX = x - 2
		boxY = y - 3
		boxW = 5
		boxH = 7

	case modell.Down:
		boxX = x - 2
		boxY = y - 4
		boxW = 5
		boxH = 7

	case modell.Left:
		boxX = x - 3
		boxY = y - 2
		boxW = 7
		boxH = 5

	case modell.Right:
		boxX = x + 3
		boxY = y - 2
		boxW = 7
		boxH = 5
	}

	for by := 0; by < boxH; by++ {
		for bx := 0; bx < boxW; bx++ {
			pos := st.Level.GetCoords(boxX+bx, boxY+by)
			for _, snake := range st.Snakes {
				for _, block := range snake.Body {
					if block.Coord.X == pos.X && block.Coord.Y == pos.Y {
						return false
					}
				}
			}
		}
	}

	return true

}

// Update updates the game-state
func (st *State) Update() {

}
