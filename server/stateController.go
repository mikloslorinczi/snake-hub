package server

import (
	"math/rand"
	"sync"
	"time"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/utils"
	termbox "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

type stateController struct {
	stateChan chan modell.State
	closeChan chan struct{}
	state     modell.State
	mu        sync.RWMutex
}

// GetUser returns the user associated with the given ID
func (sc *stateController) GetUser(id string) (bool, *modell.User) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	for _, user := range sc.state.Users {
		if user.ID == id {
			return true, &user
		}
	}
	return false, nil
}

// GetUser returns the user associated with the given ID
func (sc *stateController) GetSnake(userID string) (bool, *modell.Snake) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	for _, snake := range sc.state.Snakes {
		if snake.UserID == userID {
			return true, &snake
		}
	}
	return false, nil
}

// AddUser adds a new user to the game
func (sc *stateController) AddUser(user modell.User) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	snake := sc.GetNewSnake(user.ID)
	user.SnakeID = snake.ID
	sc.state.Users = append(sc.state.Users, user)
	go sc.AddSnake(*snake)
}

// AddSnake adds a new snake to the game
func (sc *stateController) AddSnake(snake modell.Snake) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.state.Snakes = append(sc.state.Snakes, snake)
}

// RemoveUser removes an user from the game
func (sc *stateController) RemoveUser(id string) bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	for i, user := range sc.state.Users {
		if user.ID == id {
			sc.state.Users = append(sc.state.Users[:i], sc.state.Users[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveSnake removes the snake associated with the given User ID
func (sc *stateController) RemoveSnake(userID string) bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	for i, snake := range sc.state.Snakes {
		if snake.UserID == userID {
			for _, user := range sc.state.Users {
				if user.SnakeID == snake.ID {
					user.SnakeID = ""
				}
			}
			sc.state.Snakes = append(sc.state.Snakes[:i], sc.state.Snakes[i+1:]...)
			return true
		}
	}
	return false
}

// GetNewSnake generates a new snake, on a valid position, facing to a random direction
func (sc *stateController) GetNewSnake(userID string) *modell.Snake {

	x, y := 0, 0
	direction := modell.Up
	for {
		direction = modell.RandomDirection()
		x, y = rand.Intn(sc.state.Level.Width), rand.Intn(sc.state.Level.Height)
		if sc.validSnakePos(x, y, direction) {
			break
		}
	}

	leftRune, rightRune := modell.GetRandomTexture()

	snake := &modell.Snake{
		ID:            utils.NewID(),
		UserID:        userID,
		Color:         termbox.Attribute(rand.Int()%8) + 1,
		BgColor:       termbox.Attribute(rand.Int()%8) + 1,
		HeadRune:      modell.GetRandomHead(),
		LeftRune:      leftRune,
		RightRune:     rightRune,
		Body:          modell.ClaculateSnakeBody(x, y, 3, direction, sc.state.Level),
		Direction:     direction,
		NextDirection: direction,
		TargetLength:  3,
		Speed:         rand.Intn(6) + 1,
		StepSize:      0,
	}

	return snake

}

func (sc *stateController) validSnakePos(x, y int, direction modell.Direction) bool {

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
			pos := sc.state.Level.GetCoords(boxX+bx, boxY+by)
			for _, snake := range sc.state.Snakes {
				for _, block := range snake.Body {
					if block.X == pos.X && block.Y == pos.Y {
						return false
					}
				}
			}
		}
	}

	return true

}

func (sc *stateController) ChangeDirection(userID, direction string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	for i, snake := range sc.state.Snakes {
		if snake.UserID == userID {
			newDirection := modell.StringToDirection(direction)
			if !snake.Direction.IsOpposite(newDirection) {
				sc.state.Snakes[i].NextDirection = newDirection
			}
			return
		}
	}
}

// Update updates the game-state
func (sc *stateController) Update() {
	for i := range sc.state.Snakes {
		sc.state.Snakes[i].Update(sc.state.Level)
	}
}

// updateAndBroadcast the game-state to all clients
func (sc *stateController) updateAndBroadcast() {
	log.Info("State updater starting")
updateLoop:
	for {
		select {

		case <-sc.closeChan:
			break updateLoop

		default:

			nextTick := time.Now().Add(time.Millisecond * 30) // ~33.3 FPS

			go func() {
				sc.mu.Lock()
				defer sc.mu.Unlock()
				sc.Update()
				sc.stateChan <- sc.state
			}()

			time.Sleep(time.Until(nextTick))

		}
	}
	log.Info("State updater stoped")
}
