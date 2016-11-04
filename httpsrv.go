package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func smHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "smhandler")
}
func pingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pingHandler")
}
func db2Handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "db2Handler")
}

func httpSrv(port int) error {
	http.Handle("/", regRouter())

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func dbRead(db *bolt.DB) error {
	return nil
}

func regRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sm monitor %s", time.Now().Format(time.RFC3339))
	})

	r.HandleFunc("/sm", smHandler)
	r.HandleFunc("/ping/{tar}", pingHandler)
	r.HandleFunc("/db2", db2Handler)

	return r
}
