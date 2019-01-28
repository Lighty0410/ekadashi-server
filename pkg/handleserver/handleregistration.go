package handleserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type AuthName struct {
	username string
	password string
}

type Server struct {
	http.Server
	shutDown chan bool
	reqCount uint32
}

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
		log.Printf("Shutdown error", err)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	info := &AuthName{}
	if r.Method == "GET" {
		t, err := template.ParseFiles("handleserver/login.gtpl")
		if err != nil {
			fmt.Println("something went wrong")
		} else {
			t.Execute(w, nil)
		}
	} else {
		r.ParseForm()
		for _, word := range r.Form["username"] {
			fmt.Fprintf(w, "Username %v \n", word)
			info.username += word
		}
		for _, word := range r.Form["password"] {
			fmt.Fprintf(w, "Password %v \n", word)
			info.password += word
		}
	}
}

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
