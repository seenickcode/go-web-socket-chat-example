package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errgo"
)

const (
	templatesPath     = "web/templates"
	applicationLayout = "application.html"
)

type ResponseError struct {
	StatusCode int    `json:"status"`
	ErrorMsg   string `json:"error"`
}

func renderPage(web *Web, w http.ResponseWriter, r *http.Request, path string, file string, data map[string]interface{}) {
	renderPageWithLayout(web, w, r, path, file, applicationLayout, data)
}

func renderPageWithLayout(web *Web, w http.ResponseWriter, r *http.Request, path string, file string, layout string, data map[string]interface{}) {
	log.Debugf("rendering template with layout %+v at path %+v, page %+v", layout, path, file)

	// TODO cache templates ahead of time https://golang.org/doc/articles/wiki/#tmp_10

	funcs := web.defaultTemplHelpers(w, r)

	// gotcha, .New() must take base filename
	t, err := template.New(file).Funcs(funcs).ParseFiles(
		filepath.Join(templatesPath, path, file),
		filepath.Join(templatesPath, "layouts", layout),
	)
	if err != nil {
		log.Errorf("error rendering html template %+v: %+v", file, err.Error())
		http.Error(w, "error", http.StatusInternalServerError) // TODO show err.Error() if dev else redirect to error page
		return
	}

	// establish default template data
	// (would do this using FuncMap but mux Flash must be saved _before_ templ executes,
	// so that wouldn't work)
	templData := map[string]interface{}{
		"FlashError":  web.session.Flash(w, r, FlashError),
		"FlashNotice": web.session.Flash(w, r, FlashNotice),
	}
	// merge passed in data with our template data
	for k, v := range data {
		templData[k] = v
	}
	log.Debugf("executing template with data: %+v", templData)
	executeTemplate(w, t, templData)
}

func render404(w http.ResponseWriter) {
	t, _ := template.ParseFiles(
		filepath.Join(templatesPath, "404.html"),
	)
	executeTemplate(w, t, nil)
}

func render500(w http.ResponseWriter) {
	t, _ := template.ParseFiles(
		filepath.Join(templatesPath, "500.html"),
	)
	executeTemplate(w, t, nil)
}

func executeTemplate(w http.ResponseWriter, t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err := t.Execute(w, data)
	if err != nil {
		log.Error(err)
	}
}

func renderJSON(w http.ResponseWriter, value interface{}, status int) {
	log.Debugf("rendering JSON: %+v", value)
	body, _ := json.Marshal(value)
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, _ = w.Write(body)
	return
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

func renderJSONIfError(w http.ResponseWriter, err error, status int) (ok bool) {
	if err != nil {
		renderJSONError(w, err, http.StatusInternalServerError)
	} else {
		ok = true
	}
	return
}

func renderText(w http.ResponseWriter, text string, status int) {
	log.Debugf("rendering text: %s", text)
	body := []byte(text)
	w.Header().Set("Content-Type", "text/plain; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, _ = w.Write(body)
	return
}
