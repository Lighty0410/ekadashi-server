package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

func (s *EkadashiServer) showEkadashiEnpoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, nil)
		return
	}
	err = s.checkAuth(cookie.Value)
	if err == mongo.ErrNoSession {
		jsonResponse(w, http.StatusUnauthorized, nil)
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("cannot check authentification: %v", err))
		return
	}
	ekadashiDate, err := s.db.LastEkadashi(time.Now())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("there is no information in database: %v", err))
	}
	jsonResponse(w, http.StatusOK, ekadashiDate.Date)
}
