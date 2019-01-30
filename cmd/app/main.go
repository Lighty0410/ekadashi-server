package main

import (
	"net/http"
	"testsmt/pkg/server"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", server.Registration).Methods("POST")
	http.ListenAndServe(":9000", r)
}
