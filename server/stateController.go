package server

import (
	"sync"
	"time"

	"github.com/mikloslorinczi/snake-hub/modell"
	log "github.com/sirupsen/logrus"
)

type stateController struct {
	stateChan chan modell.State
	closeChan chan struct{}
	state     modell.State
	mu        sync.RWMutex
}

func (sc *stateController) getUser(id string) (bool, *modell.User) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.state.GetUser(id)
}

func (sc *stateController) getSnake(userID string) (bool, *modell.Snake) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.state.GetSnake(userID)
}

func (sc *stateController) getFood(id string) (bool, *modell.Food) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.state.GetFood(id)
}

func (sc *stateController) addUser(user modell.User) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.state.AddUser(user)
}

func (sc *stateController) addSnake(snake modell.Snake) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.state.AddSnake(snake)
}

func (sc *stateController) addFood(food modell.Food) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.state.AddFood(food)
}

func (sc *stateController) removeUser(id string) bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.state.RemoveUser(id)
}

func (sc *stateController) removeSnake(userID string) bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.state.RemoveSnake(userID)
}

func (sc *stateController) removeFood(id string) bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.state.RemoveFood(id)
}

func (sc *stateController) newFood() {
	sc.mu.Lock()
	sc.mu.Unlock()
	sc.state.NewFood()
}

func (sc *stateController) changeDirection(userID, direction string) {
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

func (sc *stateController) updateFoods() {
	if len(sc.state.Foods) == 0 {
		sc.state.NewFood()
	}
}

func (sc *stateController) checkCollosions() {
	for i := range sc.state.Snakes {
		id := sc.state.Snakes[i].UserID
		head := sc.state.Snakes[i].GetHeadCoords()
		for _, food := range sc.state.Foods {
			if head.X == food.Pos.X && head.Y == food.Pos.Y {
				sc.state.Snakes[i].TargetLength += 3
				sc.state.RemoveFood(food.ID)
				sc.state.NewFood()
			}
		}
		for _, snake := range sc.state.Snakes {
			for j, block := range snake.Body {
				// Do not collide with own head
				if head.X == block.X && head.Y == block.Y && id == snake.UserID && j == 0 {
					continue
				}
				if head.X == block.X && head.Y == block.Y {
					sc.state.RemoveSnake(id)
					sc.state.AddSnake(*sc.state.GetNewSnake(id))
				}
			}
		}
	}
}

// Update updates the game-state
func (sc *stateController) update() {
	for i := range sc.state.Snakes {
		sc.state.Snakes[i].Update(sc.state.Level)
	}
	sc.checkCollosions()
	sc.updateFoods()
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
				sc.update()
				sc.stateChan <- sc.state
			}()

			time.Sleep(time.Until(nextTick))

		}
	}
	log.Info("State updater stoped")
}
