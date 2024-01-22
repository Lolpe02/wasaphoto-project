package api

type user struct {
	UserId   int64   `json:"userId"`
	Name     string  `json:"userName"`
	Date     string  `json:"date"`
	Profile  []int64 `json:"posted"`
	Follows  []int64 `json:"following"`
	Followed []int64 `json:"followed"`
}
type post struct {
	Creator int64  `json:"userId"`
	Desc    string `json:"description"`
	Date    string `json:"date"`
}
type like struct {
	UserId int64 `json:"creator"`
	PostId int64 `json:"post"`
}
type comment struct {
	CommentId int64  `json:"commentId"`
	UserName  string `json:"creator"`
	Text      string `json:"content"`
	Date      string `json:"date"`
}
type follow struct {
	FollowingId int64 `json:"followingId"`
	FollowedId  int64 `json:"followedId"`
}
type ban struct {
	BannerId int64 `json:"bannerId"`
	BannedId int64 `json:"bannedId"`
}
