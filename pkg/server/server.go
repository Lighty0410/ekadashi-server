package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
	"github.com/gorilla/mux"
)

// EkadashiServer serves ekadashi HTTP requests.
type EkadashiServer struct {
	*mux.Router
	db *mongo.Service
}

// NewEkadashiServer sets up http routs and returns server ready to use in http.ListenAndServe.
func NewEkadashiServer(db *mongo.Service) (*EkadashiServer, error) {
	s := &EkadashiServer{
		Router: mux.NewRouter(),
		db:     db,
	}
	s.Use(withLogging)
	s.Methods("POST").Path("/register").HandlerFunc(s.handleRegistration)
	s.Methods("POST").Path("/login").HandlerFunc(s.handleLogin)
	s.Methods("GET").Path("/ekadashi/next").HandlerFunc(s.nextEkadashiHandler)
	s.Methods("GET").Path("/ekadashi/next").HandlerFunc(s.nextEkadashiHandler)
	err := s.fillEkadashi()
	if err != nil {
		return nil, fmt.Errorf("cannot fill ekadashiAPI: %v", err)
	}
	return s, nil
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
