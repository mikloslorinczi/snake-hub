package client

import (
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mikloslorinczi/snake-hub/modell"
)

// stateController manages the client-side game-state and userinfo,
// behind a mutex making it thread-safe
type stateController struct {
	userName       string
	userSnakestyle modell.SnakeStyle
	stateCh        chan string
	errorCh        chan error
	stopCh         chan struct{}

	state  *modell.State
	mu     sync.RWMutex
	loaded bool
	setup  bool
}

// newStateController returns a pointer to a stateController loaded with username and snake-style
func newStateController(UserName string, SnakeStyle modell.SnakeStyle, StateCh chan string, ErrorCh chan error, StopCh chan struct{}) *stateController {
	return &stateController{
		userName:       UserName,
		userSnakestyle: SnakeStyle,
		stateCh:        StateCh,
		errorCh:        ErrorCh,
		stopCh:         StopCh,
	}
}

func (st *stateController) init() {
	go st.stateReader()
	log.Debug("StateController intialized successfully")
}

func (st *stateController) stateReader() {
	for {
		select {
		case <-st.stopCh:
			log.Debug("State stateReader stopped")
			return
		case state := <-st.stateCh:
			if err := st.loadState([]byte(state)); err != nil {
				st.errorCh <- errors.Wrap(err, "Cannot load state")
			}
		}
	}
}

// loadState loads the state from a JSON object
func (st *stateController) loadState(bytes []byte) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	return json.Unmarshal(bytes, &st.state)
}

// getLvlSize returns the Level's width height
func (st *stateController) getLvlSize() (int, int) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Level.Width, st.state.Level.Height
}

// getScene returns the actual scene
func (st *stateController) getScene() string {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Scene
}

// getTextbox returns the actual content of the textbox
func (st *stateController) getTextbox() []string {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Textbox
}

// getSnakes returns all the snakes in the game
func (st *stateController) getSnakes() []modell.Snake {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Snakes
}

// getFoods returns all foods in the game
func (st *stateController) getFoods() []modell.Food {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Foods
}

// getUsers returns all the users
func (st *stateController) getUsers() []modell.User {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Users
}
