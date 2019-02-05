package server

import (
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
	s.Methods("POST").Path("/register").HandlerFunc(s.handleRegistration)
	//s.Methods("POST").Path("/login").HandlerFunc(s.login)

	return s, nil
}
