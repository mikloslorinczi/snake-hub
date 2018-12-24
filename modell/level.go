package modell

import (
	"math"

	termbox "github.com/nsf/termbox-go"
)

// Coords represents a coordinate (X and Y)
type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// LevelMap stores the width and height of the map
type LevelMap struct {
	Width   int               `json:"width"`
	Height  int               `json:"height"`
	BgColor termbox.Attribute `json:"bgcolor"`
}

// GetCoords accepts x y as int and returns respective thorus coordinates
func (lv *LevelMap) GetCoords(x, y int) Coords {
	coords := Coords{}
	if x < 0 {
		coords.X = lv.Width - int(math.Abs(float64(x)))%lv.Width
	} else {
		coords.X = x % lv.Width
	}
	if y < 0 {
		coords.Y = lv.Height - int(math.Abs(float64(y)))%lv.Height
	} else {
		coords.Y = y % lv.Height
	}
	return coords
}
