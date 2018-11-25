package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
)

// Run starts the server...
func Run() {

	http.HandleFunc("/echo", echo)

	http.HandleFunc("/", home)

	fmt.Printf("\nSnake-hub listening on PORT %v\n", viper.GetInt("SNAKE_PORT"))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("SNAKE_PORT")), nil))

}
