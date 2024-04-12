package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
)

// stream username in U mode (in db it calls only for n mode) set and getprofile
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// PLEASE NOTE that this will call setMyUsername with mode U
	// cause mode N is encapsulated in doLogin signin operation
	// hence it is also obfuscated from openapi design
	println("setMyUsername called")

	userId := ps.ByName("userId")

	println("userId: ", userId)

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

	println("has authorization")

	var newUsername string

	// retrieve username from body
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	println("newUsername: ", newUsername)

	// [x]check new username is valid (unicity will be checked in db, is it a good idea? )
	if !checkRegexNewUsername(newUsername) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("new username is not valid")
		return
	}
	// TODO ERRORS if name exists
	// now call db to set username
	err = rt.db.SetMyUserName(newUsername, userId, "U")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during db calls in setting username: ", err)
		return
	}

}

// [ ] IMPORTANT this should be implemented in login also !!!!
// FIXME
func checkRegexNewUsername(username string) bool {

	println("regex call")

	usernameRegex := "^[a-z0-9]*?$"
	matched, err := regexp.MatchString(usernameRegex, username)
	if err != nil {
		return false
	}

	println("matched: ", matched)

	return matched
}
