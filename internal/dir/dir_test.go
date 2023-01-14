package dir

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DirExists(t *testing.T) {
	assert.True(t, DirExists(Options{
		Dir:      "../dir",
		FileMode: os.FileMode(0705),
	}))

	assert.False(t, DirExists(Options{
		Dir:      "./no-exist",
		FileMode: os.FileMode(0705),
	}))
}

func Test_CreateDir(t *testing.T) {
	folderPath := "./test"
	err := CreateDir(Options{
		Dir:      folderPath,
		FileMode: os.FileMode(0705),
	})
	assert.NoError(t, err)

	// clean test
	err = os.RemoveAll(folderPath)
	assert.NoError(t, err)
}
