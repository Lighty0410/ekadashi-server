package http

import (
	"net/http"

	"github.com/Lighty0410/ekadashi-server/pkg/server/controller"
)

func (s *EkadashiServer) nextEkadashiHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, nil)
		return
	}
	date, err := s.controller.ShowEkadashi(cookie.Value)
	if err != nil {
		switch err {
		case controller.ErrNotFound:
			jsonError(w, http.StatusUnauthorized, err)
			return
		default:
			jsonError(w, http.StatusInternalServerError, err)
			return
		}
	}
	jsonResponse(w, http.StatusOK, date.Format("January 2 2006"))
}
