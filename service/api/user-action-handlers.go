package api

import (
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
)

// stream username in U mode (in db it calls only for n mode) set and getprofile
func (rt *_router) setMyUsername(w http.ResponseWriter, r http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// PLEASE NOTE that this will call setMyUsername with mode U
	// cause mode N is encapsulated in doLogin signin operation
	// hence it is also obfuscated from openapi design

	userId := ps.ByName("userId")
	// [x] handle empty user id
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if user is allowed
	authorization := r.Header.Get("Authorization")
	if userId != authorization {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// retrieve username from body

	// check new username is valid (unicity will be checked in db)

}

func checkRegexNewUsername(username string) (bool, error) {
	usernameRegex := "^[a-z0-9]*?$"
	matched, err := regexp.MatchString(usernameRegex, username)
	if err != nil {
		return false, err
	}
	return matched, nil
}
