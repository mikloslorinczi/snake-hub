package modell

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/nsf/termbox-go"

	"github.com/mikloslorinczi/snake-hub/utils"
)

// State represents the game-state
// It can marshalled to JSON and sent over WebSockets
type State struct {
	Users    []User   `json:"users"`
	Snakes   []Snake  `json:"snakes"`
	Foods    []Food   `json:"foods"`
	Level    LevelMap `json:"level"`
	Scene    string   `json:"scene"`
	Textbox  []string `json:"textbox"`
	MaxScore int      `json:"maxscore"`
}

// NewState creates a new game with an empty map of width * height, in initial Scene
func NewState(width, height, maxScore int, bgColor termbox.Attribute) *State {
	return &State{
		Level: LevelMap{
			Width:   width,
			Height:  height,
			BgColor: bgColor,
		},
		MaxScore: maxScore,
		Scene:    "wait",
	}
}

// GetUser returns the user associated with the given ID
func (s *State) GetUser(id string) (bool, *User) {
	for _, user := range s.Users {
		if user.ID == id {
			return true, &user
		}
	}
	return false, nil
}

// GetSnake returns the user associated with the given ID
func (s *State) GetSnake(userID string) (bool, *Snake) {
	for _, snake := range s.Snakes {
		if snake.UserID == userID {
			return true, &snake
		}
	}
	return false, nil
}

// GetFood returns the food associated with the given ID
func (s *State) GetFood(id string) (bool, *Food) {
	for _, food := range s.Foods {
		if food.ID == id {
			return true, &food
		}
	}
	return false, nil
}

// AddUser adds a new user to the game
func (s *State) AddUser(user User) {
	snake := s.GetNewSnake(user.ID)
	user.SnakeID = snake.ID
	s.Users = append(s.Users, user)
	s.AddSnake(*snake)
}

// AddSnake adds a new snake to the game
func (s *State) AddSnake(snake Snake) {
	s.Snakes = append(s.Snakes, snake)
}

// AddFood adds a new food to the game
func (s *State) AddFood(food Food) {
	s.Foods = append(s.Foods, food)
}

// RemoveUser removes an user from the game
func (s *State) RemoveUser(id string) bool {
	for i, user := range s.Users {
		if user.ID == id {
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveSnake removes the snake associated with the given User ID
func (s *State) RemoveSnake(userID string) bool {
	for i, snake := range s.Snakes {
		if snake.UserID == userID {
			for _, user := range s.Users {
				if user.SnakeID == snake.ID {
					user.SnakeID = ""
				}
			}
			s.Snakes = append(s.Snakes[:i], s.Snakes[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveFood removes a food from the game
func (s *State) RemoveFood(id string) bool {
	for i, food := range s.Foods {
		if food.ID == id {
			s.Foods = append(s.Foods[:i], s.Foods[i+1:]...)
			return true
		}
	}
	return false
}

// GetNewSnake generates a new snake, on a valid position, facing to a random direction
func (s *State) GetNewSnake(userID string) *Snake {

	x, y := 0, 0
	direction := Up
	for {
		direction = RandomDirection()
		x, y = rand.Intn(s.Level.Width), rand.Intn(s.Level.Height)
		if s.validSnakePos(x, y, direction) {
			break
		}
	}

	leftRune, rightRune := GetRandomTexture()

	bgColor := termbox.Attribute(rand.Int()%8) + 1
	for bgColor == termbox.ColorDefault || bgColor == termbox.ColorBlack {
		bgColor = termbox.Attribute(rand.Int()%8) + 1
	}

	snake := &Snake{
		ID:            utils.NewID(),
		UserID:        userID,
		Color:         termbox.Attribute(rand.Int()%8) + 1,
		BgColor:       bgColor,
		HeadRune:      GetRandomHead(),
		LeftRune:      leftRune,
		RightRune:     rightRune,
		Body:          ClaculateSnakeBody(x, y, 3, direction, s.Level),
		Direction:     direction,
		NextDirection: direction,
		TargetLength:  3,
		Speed:         rand.Intn(6) + 1,
		StepSize:      0,
		Alive:         true,
	}

	return snake

}

// validSnakePos will return true if there is enugh space for a new snake in the given position
func (s *State) validSnakePos(x, y int, direction Direction) bool {

	boxX, boxY, boxW, boxH := 0, 0, 0, 0

	switch direction {

	case Up:
		boxX = x - 2
		boxY = y - 3
		boxW = 5
		boxH = 7

	case Down:
		boxX = x - 2
		boxY = y - 4
		boxW = 5
		boxH = 7

	case Left:
		boxX = x - 3
		boxY = y - 2
		boxW = 7
		boxH = 5

	case Right:
		boxX = x + 3
		boxY = y - 2
		boxW = 7
		boxH = 5
	}

	return s.isEmpty(boxX, boxY, boxW, boxH)

}

// isEmpty will return true if there is no snake or food in the given area
func (s *State) isEmpty(x, y, width, height int) bool {

	for by := 0; by < height; by++ {
		for bx := 0; bx < width; bx++ {
			pos := s.Level.GetCoords(x+bx, y+by)
			for _, snake := range s.Snakes {
				for _, block := range snake.Body {
					if block.X == pos.X && block.Y == pos.Y {
						return false
					}
				}
			}
			for _, food := range s.Foods {
				if food.Pos.X == pos.X && food.Pos.Y == pos.Y {
					return false
				}
			}
		}
	}

	return true

}

// NewFood places a new food on a random position
func (s *State) NewFood() {
	x, y := 0, 0
	for {
		x, y = rand.Intn(s.Level.Width), rand.Intn(s.Level.Height)
		if s.isEmpty(x-1, y-1, 3, 3) {
			break
		}
	}
	food := Food{
		ID: utils.NewID(),
		Pos: Coords{
			X: x,
			Y: y,
		},
		Type: GetRandomFood(),
	}
	s.AddFood(food)
}

// AddScore search for User with the given ID and adds score to his/her scores
func (s *State) AddScore(id string, score int) {
	for i, user := range s.Users {
		if user.ID == id {
			s.Users[i].Score += score
			return
		}
	}
}

// GetWinner check if any player has reached the State's max score, and return that player's ID
func (s *State) GetWinner() string {
	for _, user := range s.Users {
		if user.Score >= s.MaxScore {
			return user.ID
		}
	}
	return ""
}

// ScoresToTextbox will format the State's Textbox to represent Player scores
func (s *State) ScoresToTextbox() {
	sort.Slice(s.Users, func(i, j int) bool {
		return s.Users[i].Score > s.Users[j].Score
	})
	s.Textbox = []string{
		"Round Over",
		"",
	}
	for _, user := range s.Users {
		s.Textbox = append(s.Textbox, fmt.Sprintf("Name: %s\tScore: %d", user.Name, user.Score))
		s.Textbox = append(s.Textbox, "")
	}
}

// NewRound removes all food and snake from the State, sets all player's score to 0,
// generates a new sanke for every player, and a new food
func (s *State) NewRound() {
	for i := range s.Foods {
		s.RemoveFood(s.Foods[i].ID)
	}
	for i := range s.Users {
		s.Users[i].Score = 0
		s.RemoveSnake(s.Users[i].ID)
	}
	for i := range s.Users {
		snake := s.GetNewSnake(s.Users[i].ID)
		s.Users[i].SnakeID = snake.ID
		s.AddSnake(*snake)
	}
	s.NewFood()
}
