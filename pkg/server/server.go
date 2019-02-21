package server

import (
	"log"
	"net/http"
	"time"

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
	s.Router.Handle("/register", chainMiddleware(http.HandlerFunc(s.handleRegistration), withLogging, withTracing)).
		Methods("POST")
	s.Router.Handle("/login", chainMiddleware(http.HandlerFunc(s.handleLogin), withLogging, withTracing)).
		Methods("POST")
	s.Router.Handle("/users", chainMiddleware(http.HandlerFunc(s.showAllUsers), withTracing, withLogging)).
		Methods("GET")
	return s, nil
}

type middleware func(http.Handler) http.Handler

func chainMiddleware(handler http.Handler, mw ...middleware) http.Handler {
	for _, m := range mw {
		handler = m(handler)
	}
	return handler
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logged connection from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func withTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(time.Since(start))
		log.Println(r.Method)
	})
}
