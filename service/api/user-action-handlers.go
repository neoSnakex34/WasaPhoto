package api

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var newUsername string

	// retrieve username from body
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// [x]check new username is valid (unicity will be checked in db, is it a good idea? )
	if !checkRegexNewUsername(newUsername) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// now call db to set username
	err = rt.db.SetMyUserName(newUsername, userId, "U")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during db calls in setting username: ", err)
		return
	}

}

func checkRegexNewUsername(username string) bool {
	usernameRegex := "^[a-z0-9]*?$"
	matched, err := regexp.MatchString(usernameRegex, username)
	if err != nil {
		return false
	}
	return matched
}
