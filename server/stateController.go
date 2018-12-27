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
					sc.state.Snakes[i].Alive = false
				}
			}
		}
	}
}

func (sc *stateController) removeDeadSnakes() {
	for _, snake := range sc.state.Snakes {
		if !snake.Alive {
			id := snake.UserID
			sc.state.RemoveSnake(id)
			sc.state.AddSnake(*sc.state.GetNewSnake(id))
		}
	}
}

func (sc *stateController) moveSnakes() {
	for i := range sc.state.Snakes {
		sc.state.Snakes[i].Update(sc.state.Level)
	}
}

// Update updates the game-state
func (sc *stateController) update() {
	sc.moveSnakes()
	sc.checkCollosions()
	sc.removeDeadSnakes()
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
