package client

import (
	"fmt"
	"os"

	"github.com/mikloslorinczi/snake-hub/utils"
	"github.com/mikloslorinczi/snake-hub/validator"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// Run sets up and starts the client
func Run() {

	// Setup Logger
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	if viper.GetBool("SNAKE_DEBUG") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	logFile, err := os.OpenFile("snake-hub-client.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Initialize the Game
	log.Info("Initializing Snake Client")
	id := utils.NewID()
	username := getUsername()
	game := newApp(id)
	game.setup(username)

	// Start the game and block until it exits
	game.run()
}

func getUsername() string {
	if confName := viper.GetString("SNAKE_USERNAME"); validator.ValidUsername(confName) {
		return confName
	}
	for {
		inputName, err := utils.GetInput("Enter your username (max 8 char) :")
		if err != nil {
			log.Fatalf("Error reading username %v", err)
		}
		if validator.ValidUsername(inputName) {
			return inputName
		}
		fmt.Printf("Invalid user name %s\n", inputName)
	}
}
