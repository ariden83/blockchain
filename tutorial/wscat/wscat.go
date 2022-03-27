package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var (
	url          = "wss://stream.data.alpaca.markets/v1beta1/crypto"
	extraHeaders = map[string]string{
		"action": "auth",
		"key":    "PKALQ7J68RKIHCHUZT47",
		"secret": "2KevuZB2AjUZ1di6OX9R7IQK5HhS3jV1qODsRKTx",
	}
	oneOnly = false
)

func fail(msg string, o ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, o...)
	os.Exit(1)
}

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	headers := http.Header{}

	for key, value := range extraHeaders {
		headers.Add(key, value)
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fail("failed to connect to %q: %v\n", url, err)
	}
	defer conn.Close()

	doneReading := make(chan bool)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					fail("unexpected read error %v\n", err)
				}
				break
			}
			fmt.Println(string(message))
			if oneOnly {
				break
			}
		}
		doneReading <- true
	}()

	go func() {
		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			conn.WriteMessage(websocket.TextMessage, []byte(stdin.Text()))
		}
	}()

	for {
		select {
		case <-doneReading:
			return
		case <-interrupt:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					fail("unexpected close error %v\n", err)
				}
			}
			return
		}
	}
}
