package explorer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ariden83/blockchain/cmd/web/internal/utils"
)

var (
	port    string = utils.BuildPort(DefaultExplorerPort)
	baseURL string = "http://localhost" + port
)

func Start() {
	setEnvVars()
	loadTemplates()
	loadFileServer()
	loadRoutes()

	listenOrDie()
}

func setEnvVars() {
	portNum := GetExplorerPort()
	port = utils.BuildPort(portNum)
	baseURL = "http://localhost" + port
}

func listenOrDie() {
	fmt.Printf("🧭 HTML Explorer listening on %s\n", baseURL)
	log.Fatal(http.ListenAndServe(port, router))
}
