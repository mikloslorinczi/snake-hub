package client

import (
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/mikloslorinczi/snake-hub/modell"
)

func putBlock(b modell.Block) {
	termbox.SetCell(b.X*2, b.Y, b.LeftRune, b.Color, b.Background)
	termbox.SetCell(b.X*2+1, b.Y, b.RightRune, b.Color, b.Background)
}

func draw() {
	w, h := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y := 0; y < h; y++ {
		for x := 0; x < w/2; x++ {
			putBlock(modell.Block{
				X:          x,
				Y:          y,
				Color:      termbox.ColorDefault,
				Background: termbox.Attribute(rand.Int()%8) + 1,
				LeftRune:   ' ',
				RightRune:  ' ',
			})
		}
	}
	termbox.Flush()
}

// Render draws the scene...
func render() {
	// fmt.Printf("Start renderer...\n")
	// Init Termbox, defer close
	err := termbox.Init()
	if err != nil {
		log.Fatalf("Termbox init error %v\n", err)
	}
	defer termbox.Close()
	for running {
		nextTick := time.Now().Add(time.Millisecond * 30)
		draw()
		time.Sleep(time.Until(nextTick))
	}
}
