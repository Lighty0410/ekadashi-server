package ekadashihttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonError(w http.ResponseWriter, status int, err error) { // TODO basically our
	// Json responses is a view for HTTP. So what ? Is it necessary to move it to another .go file ?
	jsonResponse(w, status, map[string]string{"reason": err.Error()})
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (req *loginRequest) validateRequest() error { // TODO hmmm. Where should it place ?
	const minSymbols = 6
	if !validUsername(req.Username) {
		return fmt.Errorf("field username contain latin characters and numbers without space only")
	}
	if !validPassword(req.Password) {
		return fmt.Errorf("field password contain latin characters and numbers without space only")
	}
	if len(req.Username) < minSymbols {
		return fmt.Errorf("field username could not be less than 6 characters")
	}
	if len(req.Password) < minSymbols {
		return fmt.Errorf("field password could not be less than 6 characters")
	}
	return nil
}
