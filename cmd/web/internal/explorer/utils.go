package explorer

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	RequestIDHeaderKey = "X-Request-ID"
	RequestIDKey       = "RequestID"
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

type BodyReceived interface{}

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

func (e *Explorer) JSON(rw http.ResponseWriter, resp interface{}) {
	data := json.NewEncoder(rw)
	data.SetIndent("", "  ")
	data.Encode(resp)

	/* if js, err := json.Marshal(resp); err != nil {
		e.log.Error("Fail to json.Marshal", zap.Error(err))
		e.fail(http.StatusInternalServerError, err, rw)
		return

	} else if _, err := rw.Write(js); err != nil {
		e.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
		e.fail(http.StatusInternalServerError, err, rw)
		return
	}
	*/
}

/*
func (e *Explorer) setTokenHeaders(rw http.ResponseWriter, ts *auth.TokenDetails) {
	rw.Header().Set("Access_Token", ts.AccessToken)
	rw.Header().Set("Token_Type", "Bearer")
	rw.Header().Set("Expires_In", fmt.Sprintf("%d", ts.RtExpires))
	rw.Header().Set("Refresh_Token", ts.RefreshToken)
	if ts.Scope != "" {
		rw.Header().Set("Scope", ts.Scope)
	}
}
*/

func (e *Explorer) authorized(_ http.ResponseWriter, r *http.Request) (string, bool) {
	data, err := e.authServer.ValidationBearerToken(r)
	if err != nil {
		return "", false
	}
	return data.GetUserID(), true
}

func (e *Explorer) decodeBody(rw http.ResponseWriter, logCTX *zap.Logger, body io.ReadCloser, req BodyReceived) error {
	var err error
	body = http.MaxBytesReader(rw, body, 1048)

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err = dec.Decode(&req); err != nil {
		logCTX.Error("fail to decode input", zap.Error(err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return err
	}
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		logCTX.Error(msg, zap.Error(err))
		http.Error(rw, msg, http.StatusBadRequest)
		return err
	}
	return nil
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

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}
