package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// NewServer creates a new router
func (u *UserRouter) NewServer() {
	u.Server = http.Server{
		Addr:         ":8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	router := mux.NewRouter()
	router.HandleFunc("/login", u.Registration).Methods("POST")
	u.Handler = router
}
