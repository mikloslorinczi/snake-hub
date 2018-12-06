package modell

// User obj...
type User struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	SnakeID string `json:"snakeid"`
	Points  int    `json:"points"`
}
