package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

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
	defer r.Body.Close()

	log.Println("newUsername: ", newUsername)

	err = rt.db.SetMyUserName(newUsername, userId, "U")
	if errors.Is(err, customErrors.ErrInvalidRegexUsername) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("new username is not valid: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during db calls in setting username: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("username set successfully")
}

func (rt *_router) getUserList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// retrieve by header the requestor
	reqId := r.Header.Get("Requestor")

	if reqId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("requestor id has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if authorization != reqId {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to get list of users")
		_, err := w.Write([]byte("user is not allowed to get list of users"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	users, err := rt.db.GetUserList(structs.Identifier{Id: reqId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during db calls in getting list of users: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occurred during encoding list of users: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	// w.WriteHeader(http.StatusOK)
	log.Println("list of users retrieved successfully")

}
