package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// FIXME try to achieve better consistency in using identifier over strings or vice versa
	var followerId structs.Identifier
	var followedId structs.Identifier

	//TODO SHOULD I CHECK HTTP METHOD?

	followerId = structs.Identifier{Id: ps.ByName("followerId")}
	followedId = structs.Identifier{Id: ps.ByName("userId")}

	if followerId.Id == "" || followedId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("followerId or followedId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")

	println("authorization in followUser: ", authorization)
	println("followerId in followUser: ", followerId.Id)

	if followerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		// TODO check if this is enough as a ban check or i need to add another
		ctx.Logger.Error("user is not allowed to follow")
		return
	}

	// if those checks pass i will followuser in db
	err := rt.db.FollowUser(followerId, followedId)
	if errors.Is(err, customErrors.ErrAlreadyFollowing) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("user is already following")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while following user: ", err)
		return
	}

	// TODO check other errors
	w.WriteHeader(http.StatusOK)

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
		return
	}

	err := rt.db.UnfollowUser(followerId, followedId)
	if errors.Is(err, customErrors.ErrNotFollowing) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("user is not following")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while unfollowing user: ", err)
		return
	}

	// TODO check other errors
	w.WriteHeader(http.StatusOK)

}
