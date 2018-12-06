package modell

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetCoords(t *testing.T) {
	lvlMap := LevelMap{
		Width:  10,
		Height: 5,
	}
	type testStruct struct {
		coords Coords
		x, y   int
	}
	testStructs := []testStruct{
		{
			coords: Coords{
				X: 0,
				Y: 0,
			},
			x: 0,
			y: 0,
		},
		{
			coords: Coords{
				X: lvlMap.Width - 1,
				Y: lvlMap.Height - 1,
			},
			x: -1,
			y: -1,
		},
		{
			coords: Coords{
				X: lvlMap.Width - 2,
				Y: lvlMap.Height - 2,
			},
			x: -2,
			y: -2,
		},
		{
			coords: Coords{
				X: 3,
				Y: 4,
			},
			x: 3,
			y: 4,
		},
		{
			coords: Coords{
				X: lvlMap.Width - 3,
				Y: 1,
			},
			x: -3,
			y: 11,
		},
	}
	for _, testCase := range testStructs {
		require.Equal(t, testCase.coords.X, lvlMap.GetCoords(testCase.x, testCase.y).X, "GetCoords should return X: %v on input x: %v", testCase.coords.X, testCase.x)
		require.Equal(t, testCase.coords.Y, lvlMap.GetCoords(testCase.y, testCase.y).Y, "GetCoords should return Y: %v on input x: %v", testCase.coords.Y, testCase.y)
	}
}
