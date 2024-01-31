package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getContextReply) // ??
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/Images/:postId/likes/", rt.getLikes)
	rt.router.PUT("/Images/:postId/likes/:yourId", rt.like)
	rt.router.DELETE("/Images/:postId/likes/:yourId", rt.unlike)
	rt.router.GET("/Images/:postId/comments/", rt.getComments)
	rt.router.POST("/Images/:postId/comments/", rt.comment)
	rt.router.DELETE("/Images/:postId/comments/:commentId", rt.uncomment)
	rt.router.GET("/Images/:postId/metadata/", rt.GetPostMetadata)
	rt.router.PATCH("/Users/", rt.setMyUserName)

	rt.router.GET("/Users/me/myStream", rt.getMyStream)
	rt.router.POST("/Images/", rt.upload)
	rt.router.DELETE("/Images/:postId", rt.deletePhoto)
	rt.router.GET("/Images/:postId", rt.getPhoto)
	rt.router.POST("/Users/me/following/", rt.follow)
	rt.router.DELETE("/Users/me/following/:followId", rt.unfollow)
	rt.router.POST("/Users/me/muted/", rt.ban)
	rt.router.DELETE("/Users/me/muted/:bannedId", rt.unban)

	// New routes
	rt.router.GET("/Users/profile", rt.getProfile)
	rt.router.GET("/Users/me/followers/", rt.getFollowersOf)
	rt.router.GET("/Users/me/following/", rt.getFollowingNames)
	rt.router.GET("/Users/", rt.getInfo)
	// Special routes

	rt.router.GET("/liveness", rt.liveness)
	rt.router.PUT("/Database/", rt.omniPotence1)
	rt.router.POST("/Database/", rt.omniPotence2)
	return rt.router
}

// useful curl commands:
// curl -X POST -H 'Content-Type: application/json' -d '"Ermenegildo24"' http://lomalhost:3000/session
// curl -X PATCH -H 'Content-Type: application/json' -H "Authorization: Bearer 2" -d '"Ermenegildo24"' http://lomalhost:3000/Users/2
