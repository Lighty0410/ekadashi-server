package server

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
	response, date, err := s.controller.ShowEkadashi(cookie.Value)
	switch response {
	case controller.StatusUnauthorized:
		jsonError(w, http.StatusUnauthorized, err)
		return
	case controller.StatusInternalServerError:
		jsonError(w, http.StatusInternalServerError, err)
		return
	case controller.StatusConflict:
		jsonError(w, http.StatusConflict, err)
		return
	case controller.StatusBadRequest:
		jsonError(w, http.StatusBadRequest, err)
		return
	}
	jsonResponse(w, http.StatusOK, date.Format("January 2 2006"))
}
