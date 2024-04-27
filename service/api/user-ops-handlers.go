package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	serviceutilities "github.com/neoSnakex34/WasaPhoto/service/api/service-utilities"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// stream username in U mode (only mode available plain via api; dologin calls only for n mode) set and getprofile
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// PLEASE NOTE that this will call setMyUsername with mode U
	// cause mode N is encapsulated in doLogin signin operation
	// hence it is also obfuscated from openapi design
	log.Println("entered SetMyUsername")

	userId := ps.ByName("userId")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if user is allowed
	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to set username")
		return
	}

	var newUsername string

	// retrieve username from body
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		ctx.Logger.Error("an error occurred during decoding username: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// defer r.Body.Close()

	log.Println("newUsername: ", newUsername)

	// [x]check new username is valid (unicity will be checked in db, is it a good idea? )
	if !serviceutilities.CheckRegexNewUsername(newUsername) {
		w.WriteHeader(http.StatusBadRequest)
		err = customErrors.ErrInvalidRegexUsername // this is done here and not in return of checkregexnewusername cause checkregexnewusername has bool as return
		ctx.Logger.Error("new username is not valid", err)
		return
	}

	// TODO ERRORS if name exists
	// now call db to set username
	err = rt.db.SetMyUserName(newUsername, userId, "U")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during db calls in setting username: ", err)
		w.Write([]byte(err.Error()))
		return
	}

}

func (rt *_router) getListOfUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// retrieve by header the requestor
	reqId := r.Header.Get("Requestor")

	// TODO log that in frontend
	if reqId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authorization := r.Header.Get("Authorization")
	if authorization != reqId {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to get list of users")
		return
	}

	users, err := rt.db.GetUserList(structs.Identifier{Id: reqId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during db calls in getting list of users: ", err)
		return
	}

	json.NewEncoder(w).Encode(users)

}
