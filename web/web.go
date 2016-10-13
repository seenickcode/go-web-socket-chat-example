package web

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
)

// Web .
type Web struct {
	session Session
}

type webConfig struct {
}

const (
	sessionKeyUserID       = "userID"
	sessionKeyReturnURL    = "return_url"
	sessionKeyReturnAction = "return_action"
)

// New .
func New(router *mux.Router) (w *Web) {

	w = &Web{}

	session := NewSession(sessions.NewCookieStore([]byte("addyoursecretkeyhere")))
	w.session = session

	wireupRoutes(router, w)

	return
}

func (wb *Web) redirectToURLIfErr(r *http.Request, w http.ResponseWriter, err error, url string) (ok bool) {
	if err == nil {
		ok = true
		return
	}
	wb.session.SetFlash(w, r, err.Error(), FlashError)
	log.Debugf("redirecting to %v with flash error: %v", url, err.Error())
	http.Redirect(w, r, url, http.StatusSeeOther)
	return
}

func (wb *Web) renderPageIfErr(r *http.Request, w http.ResponseWriter, err error, path string, filename string, data map[string]interface{}) (ok bool) {
	if err == nil {
		ok = true
		return
	}
	wb.session.SetFlash(w, r, err.Error(), FlashError)
	log.Debugf("rendering page %v with flash error: %v", filename, err.Error())
	renderPage(wb, w, r, path, filename, data)
	return
}

func (wb *Web) createUserSession(r *http.Request, w http.ResponseWriter, userID string) {
	log.Debugf("creating user session for: %+v", userID)
	if len(userID) == 0 {
		panic("userID is empty")
	}
	wb.session.SetValueWithMaxAge(r, w, sessionKeyUserID, userID, -1)
}

func (wb *Web) destroyUserSession(r *http.Request, w http.ResponseWriter) {
	log.Debugf("destroying user session")
	wb.session.RemoveValue(r, w, sessionKeyUserID)
}

func (wb *Web) hasUserSession(r *http.Request) bool {
	log.Debugf("checking for user session")
	return wb.session.HasKey(r, sessionKeyUserID)
}
