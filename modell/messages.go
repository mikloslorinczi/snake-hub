package modell

import (
	"encoding/json"
)

// JSONable all struct which implements the marshal and unmarshal methods, to convert it to JSON and back.
type JSONable interface {
	marshal() ([]byte, error)
	unmarshal([]byte) error
}

// FromJSON will convert JSON b into struct j
// it may return on optional decoding error.
func FromJSON(j JSONable, b []byte) error {
	return j.unmarshal(b)
}

// ToJSON will convert struct j into JSON b
// it may return on optional encoding error
func ToJSON(j JSONable, b *[]byte) error {
	var err error
	*b, err = j.marshal()
	return err
}

// ResponseMsg is a general response message from the Snake-hub server
type ResponseMsg struct {
	Msg string `json:"msg"`
}

func (s *ResponseMsg) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ResponseMsg) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s)
}

// CommandObj is the representation of a command ans its tags.
// It is ment to be sent to the Infra Server so it can make a task out of it for the infra-client(s).
type CommandObj struct {
	UserID  string `json:"userid"`
	Command string `json:"command"`
}

func (s *CommandObj) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *CommandObj) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s)
}
