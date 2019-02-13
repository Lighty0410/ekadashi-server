package server

import (
	"fmt"
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
	err := db.CreateIndex()
	if err != nil {
		return s, fmt.Errorf("cannot create an index: %v", err)
	}
	s.Methods("POST").Path("/register").HandlerFunc(s.handleRegistration)
	s.Methods("POST").Path("/login").HandlerFunc(s.handleLogin)
	return s, nil
}
