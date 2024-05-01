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

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var bannerId structs.Identifier
	var bannedId structs.Identifier

	bannerId = structs.Identifier{Id: ps.ByName("bannerId")}
	bannedId = structs.Identifier{Id: ps.ByName("userId")}

	if bannerId.Id == "" || bannedId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("bannerId or bannedId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if bannerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to ban")
		w.Write([]byte("User is not allowed to ban"))
		return
	}

	err := rt.db.BanUser(bannerId, bannedId)
	if errors.Is(err, customErrors.ErrAlreadyBanned) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("user is already banned")
		w.Write([]byte("User is already banned"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while banning user")
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("User banned successfully")
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var bannerId structs.Identifier
	var bannedId structs.Identifier

	bannerId = structs.Identifier{Id: ps.ByName("bannerId")}
	bannedId = structs.Identifier{Id: ps.ByName("userId")}

	if bannerId.Id == "" || bannedId.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("bannerId or bannedId has not been provided")
		return
	}

	authorization := r.Header.Get("Authorization")
	if bannerId.Id != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed to unban")
		// this could happen only if an unlogged user tries to unban someone
		// cause frontend will encapsulate possibility of unbanning someone
		// giving the possibility only to logged user
		// so bannerId will always be equal the loggedUserId
		return
	}

	err := rt.db.UnbanUser(bannerId, bannedId)
	if errors.Is(err, customErrors.ErrNotBanned) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("user is not banned")
		w.Write([]byte("User is not banned"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("an error occured while unbanning user: ", err)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("User unbanned successfully")
}
