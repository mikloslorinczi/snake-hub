package client

// import (
// 	"fmt"

// 	termbox "github.com/nsf/termbox-go"

// 	"github.com/mikloslorinczi/snake-hub/modell"
// )

// func (term *termboxController) print(x, y int, fg, bg termbox.Attribute, msg string) {
// 	for i, c := range msg {
// 		termbox.SetCell(x+i, y, c, fg, bg)
// 	}
// }

// func (term *termboxController) putBlock(coords modell.Coords, color, bgColor termbox.Attribute, leftRune, rightRune rune) {
// 	termbox.SetCell(coords.X*2, coords.Y, leftRune, color, bgColor)
// 	termbox.SetCell(coords.X*2+1, coords.Y, rightRune, color, bgColor)
// }

// func (term *termboxController) drawResizeMsg(width, height int) {
// 	msg1 := fmt.Sprintf("Incorrect terminal size %v * %v", term.width, term.height)
// 	msg2 := fmt.Sprintf("Resize your terminal to etleast %v * %v", width, height)
// 	term.print(term.width/2-len(msg1)/2, term.height/2, termbox.ColorBlack, termbox.ColorRed, msg1)
// 	term.print(term.width/2-len(msg2)/2, term.height/2+1, termbox.ColorBlack, termbox.ColorRed, msg2)
// }

// func (term *termboxController) drawTextbox() {
// 	for i, line := range state.getTextbox() {
// 		term.print(2, 2+i, termbox.ColorBlack, termbox.ColorWhite, line)
// 	}
// }

// func (term *termboxController) drawSnake(snake modell.Snake) {
// 	for i, block := range snake.Body {
// 		// Head
// 		if i == 0 {
// 			term.putBlock(modell.Coords{
// 				X: block.X,
// 				Y: block.Y,
// 			},
// 				snake.Style.Color,
// 				snake.Style.BgColor,
// 				snake.Style.HeadRune,
// 				' ',
// 			)
// 			// Rest of the body
// 		} else {
// 			term.putBlock(modell.Coords{
// 				X: block.X,
// 				Y: block.Y,
// 			},
// 				snake.Style.Color,
// 				snake.Style.BgColor,
// 				snake.Style.LeftRune,
// 				snake.Style.RightRune,
// 			)
// 		}
// 	}
// }

// func (term *termboxController) drawSetup() {
// 	// term.drawSnake(state.userSnake)
// 	term.print(2, 2, termbox.ColorWhite, termbox.ColorBlack, "SELECT SNAKE")
// }

// func (term *termboxController) drawSnakes() {
// 	for _, snake := range state.getSnakes() {
// 		term.drawSnake(snake)
// 	}
// }

// func (term *termboxController) drawFoods() {
// 	for _, food := range state.getFoods() {
// 		term.putBlock(food.Pos, termbox.ColorDefault, termbox.ColorDefault, food.Type.LeftRune, food.Type.RightRune)
// 	}
// }

// func (term *termboxController) drawGame() {
// 	term.drawSnakes()
// 	term.drawFoods()
// }
