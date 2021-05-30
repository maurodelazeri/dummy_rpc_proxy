package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type QuiknodeProxy struct {
	app_port               string
	downstream_client_addr string
}

var (
	healthy        int32
	quiknode_proxy QuiknodeProxy
)

func init() {
	atomic.StoreInt32(&healthy, 1)

	app_port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		fmt.Println("APP_PORT is not present")
		os.Exit(1)
	}
	downstream_client_addr, ok := os.LookupEnv("CLIENT_DOWN_STREAM")
	if !ok {
		fmt.Println("CLIENT_DOWN_STREAM is not present")
		os.Exit(1)
	}
	quiknode_proxy.app_port = app_port
	quiknode_proxy.downstream_client_addr = downstream_client_addr
}

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", proxyHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/make_it_fail", make_it_failHandler)
	router.HandleFunc("/make_it_work", make_it_workHandler)
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	n.UseHandler(router)
	n.Run(":" + quiknode_proxy.app_port)
}
