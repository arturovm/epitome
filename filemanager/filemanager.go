package filemanager

import (
	"os"

	"github.com/pkg/errors"
)

func TouchDir(path string) error {
	expandedPath := os.ExpandEnv(path)
	err := os.MkdirAll(expandedPath, os.ModeDir|0755)
	if err != nil {
		return errors.Wrap(err, "error creating directory")
	}
	return nil
}
