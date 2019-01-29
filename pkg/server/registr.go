package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Lighty0410/ekadashi-server/pkg/provider"
)

// UserRouter is a providing struct
type UserRouter struct {
	userService provider.Handler
	http.Server
}

// Registration registrate users.
func (users *UserRouter) Registration(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUser(r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	err = users.userService.Create(&user)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
	}
	JSON(w, http.StatusOK, err)
}

func decodeUser(r *http.Request) (provider.User, error) {
	var u provider.User
	if r.Body == nil {
		return u, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	return u, err
}
