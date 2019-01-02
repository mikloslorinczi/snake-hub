package client

import (
	"encoding/json"
	"sync"

	"github.com/mikloslorinczi/snake-hub/modell"
)

// stateController manages the client-side game-state
type stateController struct {
	state  *modell.State
	loaded bool
	mu     sync.RWMutex
}

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

func (st *stateController) getScene() string {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Scene
}

func (st *stateController) getTextbox() []string {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Textbox
}

func (st *stateController) getSnakes() []modell.Snake {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.state.Snakes
}
