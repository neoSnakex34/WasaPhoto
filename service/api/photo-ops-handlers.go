package api

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	//TODO
	var userId structs.Identifier // uploader user id

	userId = structs.Identifier{ps.ByName("userId")}

	if userId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("userId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if userId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to upload photo")
		return
	}

	// FIXME add a check for photofile dimension max 10meg
	photoFile, err := io.ReadAll(r.Body)
	if err != nil {
		// [ ] check this and provide a better error
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while reading the photo file: ", err)
		return
	}

	// check file size
	var maxSize int = (2 ^ 20) * 10
	if len(photoFile) > maxSize {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photo file is too big, max size is 10MB")
		return
	}

	// TODO maybe check file type

}
