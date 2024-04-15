package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	serviceutilities "github.com/neoSnakex34/WasaPhoto/service/api/service-utilities"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// TODO test more and check this implementations
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	//TODO

	userId := structs.Identifier{Id: ps.ByName("userId")}

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

	photoFile, err := io.ReadAll(r.Body)
	if err != nil {
		// [ ] check this and provide a better error
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while reading the photo file: ", err)
		return
	}

	format, err := serviceutilities.CheckFileType(photoFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photo file is not a valid image format")
		return
	}

	// check file size
	var maxSize int = 1048576 // 2 power of 20 times 10
	if len(photoFile) > maxSize {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photo file is too big, max size is 10MB")
		return
	}

	// call db
	photo, err := rt.db.UploadPhoto(photoFile, userId, format)
	if err != nil {
		// TODO check this error and log better
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while uploading photo: ", err)
		return
	}
	// FIXME a whole photo with all details is returned, at the moment we only need the photoId
	// probably needs to be changed

	// return photoId
	w.WriteHeader(http.StatusCreated)
	ctx.Logger.Info("photo uploaded successfully")
	json.NewEncoder(w).Encode(photo)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoId := structs.Identifier{Id: ps.ByName("photoId")}
	userId := structs.Identifier{Id: ps.ByName("userId")}

	if photoId.Id == "" || userId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photoId or userId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if userId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to delete photo")
		return
	}

	err := rt.db.RemovePhoto(photoId, userId)

	if errors.Is(err, customErrors.ErrPhotoDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Error("photo does not exist")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while deleting photo: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("photo deleted successfully")

}
