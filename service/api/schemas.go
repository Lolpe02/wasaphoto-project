package api

type user struct {
	Id       int64   `json:"userId"`
	Name     string  `json:"userName"`
	Date     string  `json:"date"`
	Profile  []int64 `json:"posted"`
	Follows  []int64 `json:"following"`
	Followed []int64 `json:"followed"`
}
type post struct {
	Image   []byte `json:"photo"`
	Id      int64  `json:"id"`
	Creator int64  `json:"creator"`
	Date    string `json:"date"`
}
type like struct {
	UserId int64 `json:"creator"`
	PostId int64 `json:"post"`
}
type comment struct {
	PostId int64  `json:"post"`
	UserId int64  `json:"creator"`
	Date   string `json:"date"`
	Text   string `json:"content"`
}
type follow struct {
	FollowingId int64 `json:"followingId"`
	FollowedId  int64 `json:"followedId"`
}
type ban struct {
	BannerId int64 `json:"bannerId"`
	BannedId int64 `json:"bannedId"`
}
