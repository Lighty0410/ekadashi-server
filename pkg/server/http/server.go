package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	"github.com/gorilla/mux"
)

// EkadashiServer serves ekadashi HTTP requests.
type EkadashiServer struct {
	*mux.Router
	controller *controller.Controller
}

// NewEkadashiServer sets up http routs and returns server ready to use in http.ListenAndServe.
func NewServer(c *controller.Controller) error {
	s := &EkadashiServer{
		Router:     mux.NewRouter(),
		controller: c,
	}
	s.Use(withLogging)
	s.Methods("POST").Path("/register").HandlerFunc(s.handleRegistration)
	s.Methods("POST").Path("/login").HandlerFunc(s.handleLogin)
	s.Methods("GET").Path("/ekadashi/next").HandlerFunc(s.nextEkadashiHandler)
	err := s.controller.FillEkadashi(context.Background())
	if err != nil {
		return fmt.Errorf("cannot fill ekadashi dates: %v", err)
	}
	server := &http.Server{
		Addr:    ":9000",
		Handler: s.Router,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Could not listen: %v", err)
		}
	}()
	sig := <-stop
	log.Printf("Shutting down due to signal: %v", sig)
	err = server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("cannot shutdown the server: %v", err)
	}
	return nil
}

func withLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			log.Print("bodyErr <--", bodyErr.Error())
			http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
			return
		}
		err := r.Body.Close()
		if err != nil {
			log.Println("error occurred while closing file: ", err)
		}
		log.Printf("Request --> \n%s", buf)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		log.Printf("Method, URL --> %s %s", r.Method, r.URL.Path)
		log.Println("User agent -->", r.UserAgent())
		log.Printf("Logged connection from --> %s\n", r.RemoteAddr)
		log.Printf("Header --> %s", r.Header)
		wrappedHandler.ServeHTTP(w, r)
	})
}
