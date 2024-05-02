package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoId := structs.Identifier{Id: ps.ByName("photoId")}
	likerId := structs.Identifier{Id: ps.ByName("likerId")}

	if photoId.Id == "" || likerId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photoId or likerId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")

	if likerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to like photo") // not loggeed in
		_, err := w.Write([]byte("User is not allowed to like photo"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	err := rt.db.LikePhoto(likerId, photoId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while liking the photo: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Photo liked successfully")
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoId := structs.Identifier{Id: ps.ByName("photoId")}
	likerId := structs.Identifier{Id: ps.ByName("likerId")}

	if photoId.Id == "" || likerId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photoId or likerId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if likerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to unlike photo") // not loggeed in
		_, err := w.Write([]byte("User is not allowed to unlike photo"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	err := rt.db.UnlikePhoto(likerId, photoId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while unliking the photo: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Photo unliked successfully")
}
