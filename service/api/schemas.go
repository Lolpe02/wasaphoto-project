package api

type user struct {
	ID   int64  `json:"userId"`
	Name string `json:"userNname"`
	Date string `json:"date"`
}
type post struct {
	ID      int    `json:"postId"`
	Creator int64  `json:"creatorId"`
	Image   []byte `json:"image"`
	Date    string `json:"date"`
}
type like struct {
	PostID int64 `json:"postId"`
	UserID int64 `json:"userId"`
}
type comment struct {
	PostID int64  `json:"postId"`
	UserID int64  `json:"userId"`
	Text   string `json:"text"`
	Date   string `json:"date"`
}
type follow struct {
	FollowingID int64 `json:"followingId"`
	FollowedID  int64 `json:"followedId"`
}
type ban struct {
	BannedID int64 `json:"bannedId"`
	BannerID int64 `json:"bannerId"`
}
