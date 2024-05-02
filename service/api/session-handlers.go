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

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var username string

	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("could not retrieve username: ", err)
		return
	}
	defer r.Body.Close()

	userId, err := rt.db.DoLogin(username)
	if errors.Is(customErrors.ErrInvalidRegexUsername, err) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("username regular expression not matched")
		w.Write([]byte("INVALID USERNAME: use only lowercase letters and numbers; min 3, max 12 chars."))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong with login", err)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("User logged in successfully")

}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := structs.Identifier{Id: ps.ByName("userId")}

	if userId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("userId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if userId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("userId cannot retrieve stream")
		return
	}

	photos, err := rt.db.GetMyStream(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		return
	}

	err = json.NewEncoder(w).Encode(photos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	// w.WriteHeader(http.StatusOK)
	log.Println("stream retrieved")

}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requestorUserId := structs.Identifier{Id: r.Header.Get("Requestor")}
	profileUserId := structs.Identifier{Id: ps.ByName("userId")}

	if requestorUserId.Id == "" || profileUserId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("one of the ids has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if requestorUserId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("userId cannot retrieve profile")
		return
	}

	// not found would be returned from server request in the
	// unfortunate case in which user is deleted before accessing profile
	// for now i will not handle that here
	profile, err := rt.db.GetUserProfile(profileUserId, requestorUserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	// commented out cause is flagged as superfluos at runtime
	// w.WriteHeader(http.StatusOK)
	log.Println("profile retrieved")
}
