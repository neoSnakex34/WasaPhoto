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

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requestorUserId := structs.Identifier{Id: r.Header.Get("Requestor")}
	photoId := structs.Identifier{Id: ps.ByName("photoId")}

	if requestorUserId.Id == "" || photoId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("userId or photoId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if requestorUserId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to comment photo") // not logged in
		w.Write([]byte("User is not allowed to comment photo"))
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

	comment, err := rt.db.CommentPhoto(photoId, requestorUserId, commentBody)
	if errors.Is(err, customErrors.ErrIsBanned) {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is banned and cannot comment photo")
		w.Write([]byte("User is banned and cannot comment photo"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while commenting the photo: ", err)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("Photo commented successfully")
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	// w.WriteHeader(http.StatusCreated)

}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commentId := structs.Identifier{Id: ps.ByName("commentId")}
	requestorUserId := structs.Identifier{Id: r.Header.Get("Requestor")}

	if commentId.Id == "" || requestorUserId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("commentId or commentingUserId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if requestorUserId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to uncomment photo") // not logged in
		w.Write([]byte("User is not allowed to uncomment photo"))
		return
	}

	err := rt.db.UncommentPhoto(commentId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while uncommenting the photo: ", err)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Photo uncommented successfully")

}
