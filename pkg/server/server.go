package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func NewServer() *Server {
	s := &Server{
		Server: http.Server{
			Addr:         ":8080",
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
		shutDown: make(chan bool),
	}
	router := mux.NewRouter()
	router.HandleFunc("/login", s.Login)
	s.Handler = router
	return s
}

// WaitShutdown waits until shutdown's ready. Move a piece of code to main
func (s *Server) WaitShutdown() {
	sigint := make(chan os.Signal)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigint:
		log.Printf("Signal : %v", sig)
	case sig := <-s.shutDown:
		log.Printf("Signal : %v", sig)
	}
	log.Printf("Stop listening to server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown error %v", err)
	}
}
