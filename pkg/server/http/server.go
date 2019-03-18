package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
	"github.com/gorilla/mux"
)

// EkadashiServer serves ekadashi HTTP requests.
type EkadashiServer struct {
	*mux.Router
	controller *controller.Controller
}

// NewEkadashiServer sets up http routs and returns server.
func NewServer(c *controller.Controller) (*http.Server, error) {
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
		return nil, fmt.Errorf("cannot fill ekadashi dates: %v", err)
	}
	server := &http.Server{
		Addr:    ":9000",
		Handler: s.Router,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Could not listen: %v", err)
		}
	}()
	return server, nil
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
