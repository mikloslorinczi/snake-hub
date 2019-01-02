package server

import (
	"fmt"
	"github.com/spf13/viper"
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
	nextRound time.Time
}

func (sc *stateController) getScene() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.state.Scene
}

func (sc *stateController) changeTextbox(newTextbox []string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.state.Textbox = newTextbox
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
			// If collide with food
			if head.X == food.Pos.X && head.Y == food.Pos.Y {
				sc.state.AddScore(id, food.Type.Score)
				sc.state.Snakes[i].TargetLength += food.Type.LengthModifier
				if newSpeed := sc.state.Snakes[i].Speed + food.Type.SpeedModifier; newSpeed >= 1 && newSpeed <= 6 {
					sc.state.Snakes[i].Speed = newSpeed
				}
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

func (sc *stateController) updateScene() {
	switch sc.state.Scene {
	case "wait":
		{
			if len(sc.state.Users) >= viper.GetInt("SNAKE_MIN_PLAYER") {
				log.WithField("Number of Players", len(sc.state.Users)).Info("Enugh player joined to start a game!")
				sc.state.NewRound()
				sc.state.Scene = "game"
				sc.state.Textbox = nil
				return
			}
			sc.state.Textbox = []string{
				"Snake Hub",
				"",
				"Waiting for more player to join...",
				"",
				fmt.Sprintf("Connected Players %d", len(sc.state.Users)),
				"",
				fmt.Sprintf("Minimum number of Players %d", viper.GetInt("SNAKE_MIN_PLAYER")),
				fmt.Sprintf("Maximum number of Players %d", viper.GetInt("SNAKE_MAX_PLAYER")),
			}
		}
	case "scores":
		{
			if sc.nextRound.Before(time.Now()) {
				log.Info("Score viewing time has passed, back to waiting...")
				sc.state.Scene = "wait"
				return
			}
		}
	case "game":
		{
			if len(sc.state.Users) == 0 {
				log.Info("All player left the game. Back to waiting...")
				sc.state.Scene = "wait"
				return
			}
			if id := sc.state.GetWinner(); id != "" {
				log.WithField("User ID", id).Info("A Player has won the game!")
				sc.state.ScoresToTextbox()
				sc.state.Scene = "scores"
				sc.nextRound = time.Now().Add(time.Second * 10)
				return
			}
		}
	}
}

// Update updates the game-state
func (sc *stateController) update() {
	sc.updateScene()
	if sc.state.Scene != "game" {
		return
	}
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
