package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/neoSnakex34/WasaPhoto/service/api/reqcontext"
)

// stream username in U mode (in db it calls only for n mode) set and getprofile
func (rt *_router) setMyUsername(w http.ResponseWriter, r http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// PLEASE NOTE that this will call setMyUsername with mode U
	// cause mode N is encapsulated in doLogin signin operation
	// hence it is also obfuscated from openapi design

	// json username check for case
	username := ps.ByName("username")
}
