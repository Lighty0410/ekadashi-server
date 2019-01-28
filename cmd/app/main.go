package main

import (
	"log"

	"github.com/Lighty0410/ekadashi-server/pkg/server"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func main() {
	server := server.NewServer()

	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and server err:", err)
		}
		done <- true
	}()
	server.WaitShutdown()
	<-done
	log.Printf("DONE!")
}
