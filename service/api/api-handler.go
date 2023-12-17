package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getName)
	rt.router.POST("/session", rt.doLogin) //.wrap(rt.getContextReply))
	rt.router.GET("/Images/:postId/likes/", rt.getLikesPhoto)
	rt.router.POST("/Images/:postId/likes/", rt.likePhoto)
	rt.router.DELETE("/Images/:postId/likes/:yourId", rt.unlikePhoto)
	rt.router.GET("/Images/:postId/comments/", rt.getCommentsPhoto)
	rt.router.POST("/Images/:postId/comments/", rt.commentPhoto)
	rt.router.DELETE("/Images/:postId/comments/:yourId", rt.uncommentPhoto)
	rt.router.PATCH("/Users/:userId", rt.setMyUserName)
	rt.router.GET("/Users/", rt.getUserProfile)
	rt.router.GET("/Users/me/myStream", rt.getMyStream)
	rt.router.POST("/Images", rt.uploadPhoto)
	rt.router.DELETE("/Images/:postId", rt.deletePhoto)
	rt.router.GET("/Images/:postId", rt.getPhoto)
	rt.router.DELETE("/Users/me/following/:followId", rt.followUser)
	rt.router.PUT("/Users/me/following/:followId", rt.unfollowUser)
	rt.router.DELETE("/Users/me/muted/:bannedId", rt.banUser)
	rt.router.PUT("/Users/me/muted/:bannedId", rt.unbanUser)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
