package explorer

import (
	"fmt"
	"sync"
)

const (
	DefaultExplorerPort = 4000
)

var (
	explorerPort int
	explorerOnce sync.Once
)

func GetExplorerPort() int {
	if explorerPort == 0 {
		return DefaultExplorerPort
	}

	return explorerPort
}

func SetExplorerPort(port int) {
	explorerOnce.Do(func() {
		explorerPort = port
	})
}

// BuildPort buils a port string from a port number.
func BuildPort(portNum int) string {
	port := fmt.Sprintf(":%d", portNum)

	return port
}
