package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

//AuthName is a struct of users that's gonna connect to the server
type AuthName struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Server is a struct wich we can to handle to
type Server struct {
	http.Server
	shutDown chan bool
	reqCount uint32
}

var users []AuthName

//NewServer creates\handle basic server router\logic
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
	router.HandleFunc("/shutdown", s.ShutdownHandler)
	s.Handler = router
	return s
}

//WaitShutdown waits until shutdown's ready
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

//Login authorize users
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user AuthName
	json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

//ShutdownHandler shutdown the server
func (s *Server) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server's downed"))
	if !atomic.CompareAndSwapUint32(&s.reqCount, 0, 1) {
		log.Printf("Shutdown through API call in progress...")
		return
	}
	go func() {
		s.shutDown <- true
	}()
}
