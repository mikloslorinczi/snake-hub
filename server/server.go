package server

import (
	"net/http"
	"os"
	"strconv"

	termbox "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"

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

	gameState.state = *modell.NewState(viper.GetInt("SNAKE_MAP_WIDTH"), viper.GetInt("SNAKE_MAP_HEIGHT"), viper.GetInt("SNAKE_MAX_SCORE"), termbox.ColorDefault)

	go wsHub.start()

	go gameState.updateAndBroadcast()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("www"))))

	// http.HandleFunc("/", http.StripPrefix("/", http.FileServer(http.Dir(viper.GetString("www")))))

	http.HandleFunc("/hub", hub)

	log.WithFields(log.Fields{
		"PORT":       viper.GetInt("SNAKE_PORT"),
		"Secret":     viper.GetString("SNAKE_SECRET"),
		"Max Score":  viper.GetInt("SNAKE_MAX_SCORE"),
		"Min Player": viper.GetInt("SNAKE_MIN_PLAYER"),
		"Max Player": viper.GetInt("SNAKE_MAX_PLAYER"),
		"Map Width":  viper.GetInt("SNAKE_MAP_WIDTH"),
		"Map Height": viper.GetInt("SNAKE_MAP_HEIGHT"),
	}).Info("Snake-hub Server listening")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("SNAKE_PORT")), nil))

}
