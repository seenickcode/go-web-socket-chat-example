package web

import (
	"html/template"
	"net/http"
	"time"
)

func (wb *Web) defaultTemplHelpers(w http.ResponseWriter, r *http.Request) template.FuncMap {
	m := template.FuncMap{
		"timeFormattedFull": func(date *time.Time) string {
			if date != nil {
				return date.Format("02/01/2006 15:04:05")
			}
			return ""
		},
	}
	return m
}
