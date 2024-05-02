package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var followerId structs.Identifier
	var followedId structs.Identifier

	followerId = structs.Identifier{Id: ps.ByName("followerId")}
	followedId = structs.Identifier{Id: ps.ByName("userId")}

	if followerId.Id == "" || followedId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("followerId or followedId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")

	if followerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to follow")
		_, err := w.Write([]byte("User is not allowed to follow"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	// if those checks pass i will followuser in db
	err := rt.db.FollowUser(followerId, followedId)
	if errors.Is(err, customErrors.ErrAlreadyFollowing) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("user is already following")
		_, err = w.Write([]byte("User is already following"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while following user: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("User followed successfully")

}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var followerId structs.Identifier
	var followedId structs.Identifier

	followerId = structs.Identifier{Id: ps.ByName("followerId")}
	followedId = structs.Identifier{Id: ps.ByName("userId")}

	if followerId.Id == "" || followedId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("followerId or followedId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if followerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to unfollow")
		_, err := w.Write([]byte("User is not allowed to unfollow"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	err := rt.db.UnfollowUser(followerId, followedId)
	if errors.Is(err, customErrors.ErrNotFollowing) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("user is not following")
		_, err = w.Write([]byte("User is not following"))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while unfollowing user: ", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.Error("an error occurred while writing response: ", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("User unfollowed successfully")

}
