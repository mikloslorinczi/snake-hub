package server

import (
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/nsf/termbox-go"

	"github.com/mikloslorinczi/snake-hub/modell"

	"github.com/spf13/viper"
)

var (
	errorChan = make(chan error, 10)
	exitChan  = make(chan struct{}, 1)
	stateChan = make(chan modell.State, 5)

	gameState = &stateController{
		stateChan: stateChan,
		closeChan: exitChan,
		state:     *modell.NewState(25, 25, termbox.ColorDefault),
	}

	wsHub = &clientHub{
		stateChan: stateChan,
		closeChan: exitChan,
	}
)

// Run starts the Snake-hub server...
func Run() {

	// Log formatter
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	log.SetFormatter(customFormatter)

	// Log level
	if viper.GetBool("SNAKE_DEBUG") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Log file
	if !viper.GetBool("SNAKE_VERBOSE") {
		logFile, err := os.OpenFile("snake-hub-server.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	go wsHub.start()

	go gameState.updateAndBroadcast()

	http.HandleFunc("/hub", hub)

	http.HandleFunc("/", home)

	log.WithField("PORT", viper.GetInt("SNAKE_PORT")).Info("Snake-hub Server listening...")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("SNAKE_PORT")), nil))

}
