package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := structs.Identifier{Id: ps.ByName("commentingId")} // FIXME consistency via db
	photoId := structs.Identifier{Id: ps.ByName("photoId")}

	if userId.Id == "" || photoId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("userId or photoId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if userId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to comment photo") // not logged in
		return
	}

	var commentBodyReq structs.BodyRequest
	err := json.NewDecoder(r.Body).Decode(&commentBodyReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("comment body has not been provided")
		return
	}

	commentBody := commentBodyReq.Body

	// TODO fix this, check return for consistency etc etc
	comment, err := rt.db.CommentPhoto(photoId, userId, commentBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while commenting the photo: ", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	ctx.Logger.Info("photo commented successfully")
	json.NewEncoder(w).Encode(comment)

}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commentId := structs.Identifier{Id: ps.ByName("commentId")}
	commentingUserId := structs.Identifier{Id: ps.ByName("commentingId")}

	if commentId.Id == "" || commentingUserId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("commentId or commentingUserId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if commentingUserId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to uncomment photo") // not logged in
		return
	}

	err := rt.db.UncommentPhoto(commentId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while uncommenting the photo: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("photo uncommented successfully")

}
