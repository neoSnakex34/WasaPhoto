package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var username string
	ctx.Logger.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("could not fetch username: ", err)
		return
	}

	defer r.Body.Close()

	userId, err := rt.db.DoLogin(username)
	if errors.Is(customErrors.ErrInvalidRegexUsername, err) {
		w.WriteHeader(http.StatusBadRequest)
		// TODO log all the errors in frontend like this
		// FIXME very important
		w.Write([]byte("INVALID USERNAME: use only lowercase letters and numbers; min 3, max 12 chars."))
		ctx.Logger.Error("username regular expression not matched")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong with login", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		return
	}

}

// FIXME LIMIT DISPLAY TO 20 in db side or here, will ease frontend and respect api design
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

	// FIXME specialize errors
	photos, err := rt.db.GetMyStream(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		return
	}

	println("photos: ", photos)

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("stream retrieved")
	json.NewEncoder(w).Encode(photos)

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

	// TODO not found handle?

	profile, err := rt.db.GetUserProfile(profileUserId, requestorUserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("profile retrieved")
	json.NewEncoder(w).Encode(profile)

}
