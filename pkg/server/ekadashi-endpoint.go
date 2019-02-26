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
	ekadashiDate, err := s.db.SendEkadashi()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Errorf("there is no information in database: %v", err))
	}

	for _, currentEkadashi := range ekadashiDate {
		nextEkadashi := time.Date(currentEkadashi.Year, currentEkadashi.Month, currentEkadashi.Day,
			10, 0, 0, 0, time.UTC)
		fmt.Println(nextEkadashi)
		if nextEkadashi.After(time.Now()) {
			jsonResponse(w, http.StatusOK, currentEkadashi)
		} else {
			continue
		}
		return
	}
}
