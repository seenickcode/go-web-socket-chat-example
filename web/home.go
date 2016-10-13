package web

import "net/http"

func (web *Web) index(w http.ResponseWriter, r *http.Request) {
	renderPage(web, w, r, "home", "index.html", nil)
}
