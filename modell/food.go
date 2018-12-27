package modell

import (
	"math/rand"

	termbox "github.com/nsf/termbox-go"
)

// FoodType represents a type of food
type FoodType struct {
	Color     termbox.Attribute `json:"color"`
	BgColor   termbox.Attribute `json:"bgcolor"`
	LeftRune  rune              `json:"leftrune"`
	RightRune rune              `json:"rightrune"`
	Score     int               `json:"score"`
}

// Food represents a bonus object that can be eaten by snakes
type Food struct {
	ID   string   `json:"id"`
	Pos  Coords   `json:"pos"`
	Type FoodType `json:"type"`
}

var (
	// Apple ...
	Apple = FoodType{
		LeftRune:  '🍎',
		RightRune: ' ',
		Score:     1,
	}
	// Apple ...
	// Apple = &FoodType{
	// 	LeftRune:  ' ',
	// 	RightRune: ' ',
	// 	Score:     1,
	// }

	// FoodTypes is a slice of all available food type
	FoodTypes = []FoodType{
		Apple,
	}
)

// GetRandomFood returns a random FoodType
func GetRandomFood() FoodType {
	return FoodTypes[rand.Intn(len(FoodTypes))]
}
