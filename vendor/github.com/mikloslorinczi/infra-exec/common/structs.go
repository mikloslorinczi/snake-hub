// Structs contain all the structures used by the Infra CLI, Server and Client
// Eache struct has a barshal() and an unmarshal() method to convert them to JSON and back
// They are all the implementation of the JSONable interface, which requires these two methods.

package common

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

// ResponseMsg is a general response message from the Infra Server
// usually related to an error.
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
	Command string `json:"command"`
	Tags    string `json:"tags"`
}

func (s *CommandObj) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *CommandObj) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s)
}

// Task is the structure of a single task. ID is provided by github.com/rs/xid. It is unique and the
// time of creation can be parsed from it. Node is the name of the executing Infra Client, initially it is "none".
// Tags hold the desired tags separated by space " ", only infra-clients with these tags can execute the task.
// Status is the current status of the task. Initially it is "Created".
// Command is the actual command and its argumentums separated by space " ".
type Task struct {
	ID      string `json:"id"`
	Node    string `json:"node"`
	Tags    string `json:"tags"`
	Status  string `json:"status"`
	Command string `json:"command"`
}

func (s *Task) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Task) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s)
}

// Tasks represent a list of Tasks
type Tasks []Task

func (s *Tasks) marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Tasks) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s)
}
