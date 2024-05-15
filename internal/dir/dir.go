package dir

import (
	"fmt"
	"os"
)

// Options represents the configuration options for directory operations.
type Options struct {
	// Required options.
	Dir      string
	File     string
	ReadOnly bool
	FileMode os.FileMode
}

// DirExists checks if the directory specified in the Options struct exists.
func DirExists(opt Options) bool {
	if _, err := os.Stat(opt.Dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDir creates a directory based on the provided Options.
func CreateDir(opt Options) error {
	dirExists := DirExists(opt)
	if !dirExists {
		if opt.ReadOnly {
			return fmt.Errorf("Cannot find directory %q for read-only open", opt.Dir)
		}
		// Try to create the directory
		if err := os.MkdirAll(opt.Dir, opt.FileMode); err != nil {
			return fmt.Errorf("Error Creating Dir: %s error: %+v", opt.Dir, err)
		}
	}

	return nil
}
