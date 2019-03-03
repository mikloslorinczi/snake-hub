package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// GetInput displays prints msg and reads user input
// from stdin as a string. May return I/O error
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
