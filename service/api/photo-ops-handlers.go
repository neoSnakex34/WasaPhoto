package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	serviceutilities "github.com/neoSnakex34/WasaPhoto/service/service-utilities"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while reading the photo file: ", err)
		w.Write([]byte("an error when attempting to read photofile"))
		return
	}

	format, err := serviceutilities.CheckFileType(photoFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photo file is not a valid image format")
		w.Write([]byte(err.Error()))
		return
	}

	// check file size
	var maxSize int = 10485760 // 2 power of 20 times 10
	if len(photoFile) > maxSize {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photo file is too big, max size is 10MB")
		w.Write([]byte("file size exceeds limit, max size is 10MB"))
		return
	}

	// call db
	newPhoto, err := rt.db.UploadPhoto(photoFile, userId, format)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while uploading photo: ", err)
		w.Write([]byte("an error occured while uploading photo"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newPhoto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	// w.WriteHeader(http.StatusOK)
	log.Println("Photo uploaded successfully")
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
	log.Println("photo deleted successfully")

}

func (rt *_router) servePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId") // user is the photo owner, requestor can be different if guest to a profile or in stream
	photoId := ps.ByName("photoId")
	requestorId := r.Header.Get("Requestor")

	if userId == "" || photoId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("userId or photoId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")

	if requestorId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to view photo")
		w.Write([]byte("could not serve photo, user is not allowed to view photo"))
		return
	}

	photoPath, err := serviceutilities.GetPhotoPath(userId + "/" + photoId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while getting photo path: ", err)
		w.Write([]byte("an error occured while getting photo path"))
		return
	}

	http.ServeFile(w, r, photoPath)
}
