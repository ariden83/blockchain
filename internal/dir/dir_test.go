package dir

import (
	"os"
	"testing"
)

func Test_DirExists(t *testing.T) {
	// Create a temporary directory for testing
	testDir := "testdir"
	os.Mkdir(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Test when directory exists
	opt := Options{Dir: testDir}
	if !DirExists(opt) {
		t.Errorf("Expected directory %s to exist, but it does not", testDir)
	}

	// Test when directory does not exist
	opt.Dir = "nonexistentdir"
	if DirExists(opt) {
		t.Errorf("Expected directory %s to not exist, but it does", opt.Dir)
	}
}

func Test_CreateDir(t *testing.T) {
	// Create a temporary directory for testing
	testDir := "testdir"
	defer os.RemoveAll(testDir)

	// Test creating a directory that doesn't exist
	opt := Options{Dir: testDir}
	err := CreateDir(opt)
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}

	// Test creating a directory that already exists
	err = CreateDir(opt)
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}

	// Test creating a directory with read-only mode
	opt.ReadOnly = true
	err = CreateDir(opt)
	if err == nil {
		t.Error("Expected error when creating directory in read-only mode, but none occurred")
	}
}
