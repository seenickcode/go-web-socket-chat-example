package main

import (
	"fmt"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/seenickcode/go-web-socket-chat-example/api"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init router
	router := mux.NewRouter()

	// init API
	api.WireupRoutes(core)

	// start server
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	n.UseHandler(router)
	log.Infof("starting chat server")
	n.Run(fmt.Sprintf(":%v", *port))
}
