package modell

import termbox "github.com/nsf/termbox-go"

type Block struct {
	X          int
	Y          int
	Color      termbox.Attribute
	Background termbox.Attribute
	LeftRune   rune
	RightRune  rune
}
