package api

type user struct {
	ID      int64   `json:"userid"`
	Name    string  `json:"username"`
	Date    string  `json:"date"`
	Profile []int64 `json:"posted"`
}
type post struct {
	Image   []byte `json:"photo"`
	ID      int64  `json:"id"`
	Creator int64  `json:"creator"`
	Date    string `json:"date"`
}
type like struct {
	UserID int64 `json:"creator"`
	PostID int64 `json:"post"`
}
type comment struct {
	PostID int64  `json:"post"`
	UserID int64  `json:"creator"`
	Date   string `json:"date"`
	Text   string `json:"content"`
}
type follow struct {
	FollowingID int64 `json:"followingId"`
	FollowedID  int64 `json:"followedId"`
}
type ban struct {
	BannedID int64 `json:"bannedId"`
	BannerID int64 `json:"bannerId"`
}
