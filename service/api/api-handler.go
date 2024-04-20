package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// session login routes
	rt.router.POST("/login", rt.wrap(rt.doLogin))

	// user related routes
	rt.router.PUT("/users/:userId/username", rt.wrap(rt.setMyUsername))
	rt.router.GET("/users", rt.wrap(rt.getListOfUsers))

	// userProfile and stream related routes
	rt.router.GET("/users/:userId/profile", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:userId/stream", rt.wrap(rt.getMyStream))

	// // ban related routes
	rt.router.PUT("/users/:userId/bans/:bannerId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userId/bans/:bannerId", rt.wrap(rt.unbanUser))

	// // follow related routes
	rt.router.PUT("/users/:userId/followers/:followerId", rt.wrap(rt.followUser)) // FIXME follows should be better as followers
	rt.router.DELETE("/users/:userId/followers/:followerId", rt.wrap(rt.unfollowUser))

	// // photo related routes
	rt.router.POST("/users/:userId/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:userId/photos/:photoId", rt.wrap(rt.deletePhoto))

	// // comment related routes
	// userid is the uploader user id
	// TODO custom header will ease this mess dunno if it is a good idea tho
	// FIXME
	rt.router.POST("/users/:userId/photos/:photoId/comments/:commentingId", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:userId/photos/:photoId/comments/:commentingId/:commentId", rt.wrap(rt.uncommentPhoto))

	// // like related routes
	rt.router.PUT("/users/:userId/photos/:photoId/likes/:likerId", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:userId/photos/:photoId/likes/:likerId", rt.wrap(rt.unlikePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
