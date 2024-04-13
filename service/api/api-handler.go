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

	// userProfile and stream related routes
	// rt.router.GET("/users/:userId/profile", rt.wrap(rt.getUserProfile))
	// rt.router.GET("/users/:userId/stream", rt.wrap(rt.getMyStream))

	// // ban related routes
	rt.router.PUT("/users/:userId/bans/:bannedId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userId/bans/:bannedId", rt.wrap(rt.unbanUser))

	// // follow related routes
	rt.router.PUT("/users/:userId/follows/:followerId", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userId/follows/:followerId", rt.wrap(rt.unfollowUser))

	// // photo related routes
	// rt.router.POST("/users/:userId/photos/:photoId", rt.wrap(rt.uploadPhoto))
	// rt.router.DELETE("/users/:userId/photos/:photoId", rt.wrap(rt.deletePhoto))

	// // comment related routes
	// rt.router.POST("/users/:userId/photos/:photoId/comments/:commentId", rt.wrap(rt.commentPhoto))
	// rt.router.DELETE("/users/:userId/photos/:photoId/comments/:commentId", rt.wrap(rt.uncommentPhoto))

	// // like related routes
	// rt.router.PUT("/users/photos/:photoId/likes/:userId", rt.wrap(rt.likePhoto))
	// rt.router.DELETE("/users/photos/:photoId/likes/:userId", rt.wrap(rt.unlikePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
