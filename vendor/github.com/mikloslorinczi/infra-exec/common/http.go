package common

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/pkg/errors"
)

// HTTPHeader is a key-value pair of strings
// passed to the HTTPRequest as header.
type HTTPHeader struct {
	key, value string
}

// HTTPHeaderJSON is Content-Type : application/json
var HTTPHeaderJSON = HTTPHeader{
	key:   "Content-Type",
	value: "application/json",
}

// SendRequest creates a HTTP request with the given Method, to the given URL,
// with (optional) body and HTTP headers.
// The function returns the response body's byte-representation (JSON), and on optional error.
func SendRequest(method string, url string, body []byte, headers ...HTTPHeader) ([]byte, error) {
	var response []byte
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		err = errors.Wrapf(err, "Cannot create %v request to %v\n", method, url)
		return response, err
	}
	for _, header := range headers {
		req.Header.Set(header.key, header.value)
	}
	req.Header.Set("APITOKEN", viper.GetString("apiToken"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "Cannot send request\n")
		return response, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			closeErr = errors.Wrap(closeErr, "Error closing response body\n")
			log.Fatal(closeErr)
		}
	}()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading response body\n")
		return response, err
	}
	if resp.StatusCode != 200 {
		return response, errors.Errorf("Server answered with a non-200 status: %v\n", resp.StatusCode)
	}
	return response, nil
}

// GetTasks sends a request to APIURL/tasks and returns
// the fetched tasks as a slice, and on optional http error.
func GetTasks() ([]Task, error) {
	var tasks Tasks
	tasksJSON, err := SendRequest("GET", viper.GetString("apiURL")+"/tasks", nil, HTTPHeaderJSON)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get task list\n")
	}
	if err := FromJSON(&tasks, tasksJSON); err != nil {
		return nil, errors.Wrap(err, "Cannot get task list\n")
	}
	return tasks, nil
}

// UpdateTaskStatus sends a request to APIURL/task/{id}/{status}
// and returns the server Msg and an optional http/encode error.
func UpdateTaskStatus(id, status string) (ResponseMsg, error) {
	var responseMsg ResponseMsg
	responseJSON, err := SendRequest("PUT", viper.GetString("apiURL")+"/task/"+id+"/"+status, nil)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update Task status\n")
	}
	if err := FromJSON(&responseMsg, responseJSON); err != nil {
		return responseMsg, errors.Wrap(err, "Cannot read response\n")
	}
	return responseMsg, nil
}
