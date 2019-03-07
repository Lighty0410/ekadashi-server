package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
)

type ekadashiJSON struct {
	Date string `json:"date"`
}

func (s *EkadashiServer) nextEkadashiHandler(w http.ResponseWriter, r *http.Request) {
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
	ekadashiDate, err := s.db.NextEkadashi(time.Now())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("cannot get next ekadashi day: %v", err))
		return
	}
	jsonResponse(w, http.StatusOK, ekadashiJSON{Date: ekadashiDate.Date.Format("January 2 2006")})
}
