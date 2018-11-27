package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// GetInput prints the message to stdout and reads user input from stdin.
// The return values are the inputString and an optional io error.
func GetInput(msg string) (string, error) {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		err = errors.Wrapf(err, "Input error")
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}
