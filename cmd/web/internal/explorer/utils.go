package explorer

import (
	"fmt"
	"time"
)

func debug(i interface{}) string {
	log := fmt.Sprintf("DEBUG: %v\n", i)
	fmt.Print(log)

	return log
}

func increment(number int) int {
	return number + 1
}

func add(a, b int) int {
	return a + b
}

func unixToHuman(unix int64) string {
	return time.Unix(unix, 0).Format(time.UnixDate)
}

func (e *Explorer) homeURL() string {
	return e.baseURL
}

func (e *Explorer) blockURL(hash string) string {
	return e.baseURL + "/blocks/" + hash
}

func (e *Explorer) txURL(hash string) string {
	return e.baseURL + "/transactions/" + hash
}

func (e *Explorer) walletURL(hash string) string {
	return e.baseURL + "/wallets/" + hash
}
