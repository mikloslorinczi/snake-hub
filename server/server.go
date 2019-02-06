package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	termbox "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"

	"github.com/mikloslorinczi/snake-hub/modell"
	"github.com/mikloslorinczi/snake-hub/utils"

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

	if strings.HasPrefix(strings.ToUpper(viper.GetString("SNAKE_ENV")), "DEV") {
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("www"))))
	} else {
		http.HandleFunc("/", home)
	}

	http.HandleFunc("/hub", hub)

	log.WithFields(log.Fields{
		"ENV":              viper.GetString("SNAKE_ENV"),
		"Development Mode": strings.HasPrefix(strings.ToUpper(viper.GetString("SNAKE_ENV")), "DEV"),
		"PORT":             viper.GetInt("SNAKE_PORT"),
		"URL":              viper.GetString("SNAKE_URP"),
		"WS URL":           utils.GetWSURL(viper.GetString("SNAKE_URL"), "Client ID", viper.GetString("SNAKE_SECRET")),
		"Secret":           viper.GetString("SNAKE_SECRET"),
		"Max Score":        viper.GetInt("SNAKE_MAX_SCORE"),
		"Min Player":       viper.GetInt("SNAKE_MIN_PLAYER"),
		"Max Player":       viper.GetInt("SNAKE_MAX_PLAYER"),
		"Map Width":        viper.GetInt("SNAKE_MAP_WIDTH"),
		"Map Height":       viper.GetInt("SNAKE_MAP_HEIGHT"),
	}).Info("Snake-hub Server has been started")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("SNAKE_PORT")), nil))

}
