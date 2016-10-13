package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errgo"
)

type ResponseError struct {
	StatusCode int    `json:"status"`
	ErrorMsg   string `json:"error"`
}

func renderJSON(w http.ResponseWriter, value interface{}, status int) {
	body, _ := json.Marshal(value)
	//log.Debugf("responding with json: %v", string(body))
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.WriteHeader(status)
	_, _ = w.Write(body)
	return
}

func renderError(w http.ResponseWriter, err error, status int) {
	log.Debugf("responding with error: %s", err.Error())
	renderErrorWithMessage(w, status, "Sorry, we couldn't process your request. Please try again later.")
}

func renderErrorWithMessage(w http.ResponseWriter, status int, message string) {
	log.Debugf("responding with error status %v message: %s", status, message)
	respErr := &ResponseError{}
	respErr.StatusCode = status
	respErr.ErrorMsg = message
	renderJSON(w, respErr, status)
}

func renderJSONError(w http.ResponseWriter, err error, status int) {
	log.Printf(errgo.Mask(err).Error())
	r := &ResponseError{}
	r.StatusCode = status
	if status == http.StatusInternalServerError {
		r.ErrorMsg = "Sorry, something went wrong. Try again later."
	} else {
		r.ErrorMsg = err.Error()
	}
	renderJSON(w, r, status)
}

func renderJSONErrorWithMessage(w http.ResponseWriter, err error, status int, message string) {
	//log.Debugf("responding with json error: %s", err.Error())
	r := &ResponseError{}
	r.StatusCode = status
	r.ErrorMsg = message
	renderJSON(w, r, status)
}

func renderString(w http.ResponseWriter, msg string, status int) {
	log.Debugf("responding with string: %v", msg)
	w.Header().Set("Content-Type", "text/plain; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg + "\n"))
	return
}
