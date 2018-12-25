package modell

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StringToDirection(t *testing.T) {
	require.Equal(t, Up, StringToDirection("up"))
	require.Equal(t, Down, StringToDirection("down"))
	require.Equal(t, Left, StringToDirection("left"))
	require.Equal(t, Right, StringToDirection("right"))
}

func Test_ChangeDirection(t *testing.T) {
	mySnake := Snake{
		Direction: Up,
	}
	mySnake.ChangeDirection("left")
	require.Equal(t, Left, mySnake.NextDirection)
	mySnake.ChangeDirection("down")
	require.Equal(t, Left, mySnake.NextDirection)
}
