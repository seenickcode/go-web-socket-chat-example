package web

import (
	"io/ioutil"
	"net/http"
)

func httpResponseToString(r *http.Response) string {
	data, _ := ioutil.ReadAll(r.Body)
	return string(data)
}
