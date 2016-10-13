package web

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
)

type Session struct {
	store *sessions.CookieStore
}

const (
	FlashError  string = "error"
	FlashNotice string = "notice"
)

func NewSession(store *sessions.CookieStore) Session {
	s := Session{}
	s.store = store
	return s
}

func (s *Session) GetValue(r *http.Request, k string) string {
	log.Debugf("getting session value for key '%v'", k)
	store := s.getStore(r)
	if val, ok := store.Values[k]; ok {
		return val.(string)
	}
	return ""
}

func (s *Session) PluckValue(r *http.Request, w http.ResponseWriter, k string) string {
	log.Debugf("plucking session value for key '%v'", k)
	store := s.getStore(r)
	if val, ok := store.Values[k]; ok {
		return val.(string)
	}
	// get then remove the value immediately
	v := s.GetValue(r, k)
	if len(v) > 0 {
		s.SetValue(r, w, k, "")
	}
	return v
}

func (s *Session) SetValue(r *http.Request, w http.ResponseWriter, k string, v string) {
	log.Debugf("setting session key '%v' with value: %v", k, v)
	store := s.getStore(r)
	store.Values[k] = v
	store.Save(r, w)
}

func (s *Session) SetValueWithMaxAge(r *http.Request, w http.ResponseWriter, k string, v string, maxAge int) {
	log.Debugf("setting session key '%v' with value: %v", k, v)
	store := s.getStore(r)
	store.Values[k] = v
	store.Options.MaxAge = (10 * 365 * 24 * 60 * 60) // in secs, 10 yrs from now
	store.Save(r, w)
}

func (s *Session) SetMaxAge(r *http.Request, w http.ResponseWriter, maxAge int) {
	store := s.getStore(r)
	store.Options.MaxAge = -1
	store.Save(r, w)
}

func (s *Session) RemoveValue(r *http.Request, w http.ResponseWriter, k string) {
	store := s.getStore(r)
	delete(store.Values, k)
	store.Save(r, w)
}

func (s *Session) HasKey(r *http.Request, k string) bool {
	store := s.getStore(r)
	_, ok := store.Values[k]
	return ok
}

func (s *Session) Flash(w http.ResponseWriter, r *http.Request, flashType string) (msg string) {
	store := s.getStore(r)
	if flashes := store.Flashes(flashType); len(flashes) > 0 {
		msg = flashes[0].(string)
	}
	store.Save(r, w)
	return msg
}

func (s *Session) SetFlash(w http.ResponseWriter, r *http.Request, msg string, flashType string) {
	store := s.getStore(r)
	store.AddFlash(msg, flashType)
	store.Save(r, w)
}

func (s *Session) getStore(r *http.Request) (store *sessions.Session) {
	store, _ = s.store.Get(r, "session")
	return store
}
