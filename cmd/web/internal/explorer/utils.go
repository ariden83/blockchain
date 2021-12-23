package explorer

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// GenericError Default response when we have an error
//
// swagger:response genericError
// nolint
type GenericError struct {
	// in: body
	Body ErrorResponse `json:"body"`
}

// ErrorResponse structure of error response
type ErrorResponse struct {
	// The status code
	Code int `json:"code"`
	// The error message
	Message string `json:"message"`
}

// fail Respond error to json format
func (e *Explorer) fail(statusCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(statusCode)

	error := ErrorResponse{
		Message: err.Error(),
		Code:    statusCode,
	}
	js, err := json.Marshal(error)
	if err != nil {
		e.log.Error("Fail to json.Marshal in Patch method", zap.Error(err))
		return
	}
	if _, err := w.Write(js); err != nil {
		e.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
	}
}

func (e *Explorer) resp(rw http.ResponseWriter, resp interface{}) {
	if js, err := json.Marshal(resp); err != nil {
		e.log.Error("Fail to json.Marshal", zap.Error(err))
		e.fail(http.StatusInternalServerError, err, rw)
		return

	} else if _, err := rw.Write(js); err != nil {
		e.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
		e.fail(http.StatusInternalServerError, err, rw)
		return
	}
}

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
