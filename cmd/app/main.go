package main

import (
	"log"

	"github.com/Lighty0410/ekadashi-server/pkg/server"
)

func main() {
	http := &server.UserRouter{}
	http.NewServer()
	err := http.ListenAndServe()
	if err != nil {
		log.Printf("Listen and serve err: %v", err)
	}
}
