package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	//rt.router.GET("/", rt.getName) //??
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/Images/:postId/likes/", rt.getLikes)
	rt.router.POST("/Images/:postId/likes/", rt.like)
	rt.router.DELETE("/Images/:postId/likes/:yourId", rt.unlike)
	rt.router.GET("/Images/:postId/comments/", rt.getComments)
	rt.router.POST("/Images/:postId/comments/", rt.comment)
	rt.router.DELETE("/Images/:postId/comments/:yourId", rt.uncomment)
	rt.router.PATCH("/Users/:userId", rt.setMyUserName)
	rt.router.GET("/Users/", rt.getProfile)
	rt.router.GET("/Users/me/myStream", rt.getMyStream)
	rt.router.POST("/Images", rt.upload)
	rt.router.DELETE("/Images/:postId", rt.deletePhoto)
	rt.router.GET("/Images/:postId", rt.getPhoto)
	rt.router.PUT("/Users/me/following/:followId", rt.follow)
	rt.router.DELETE("/Users/me/following/:followId", rt.unfollow)
	rt.router.PUT("/Users/me/muted/:bannedId", rt.ban)
	rt.router.DELETE("/Users/me/muted/:bannedId", rt.unban)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
