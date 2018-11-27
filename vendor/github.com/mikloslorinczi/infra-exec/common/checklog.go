package common

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CheckLogFolder will check if logs/ is present,
// and creates it if necessery. May return a path error.
func CheckLogFolder() error {
	logsPath := filepath.Join(".", "logs")
	if err := os.MkdirAll(logsPath, os.ModePerm); err != nil {
		return errors.Wrap(err, "Cannot create logs/ folder")
	}
	return nil
}
