package structs

// FIXME check json tags
type Identifier struct {
	Id string `json:"identifier"`
}

// type UserName struct {
// 	PlainUsername string `json:"username"`
// }

type UserProfile struct {
	UserId           Identifier `json:"userId"`
	Username         string     `json:"username"`
	FollowerCounter  int        `json:"followersCounter"`
	FollowingCounter int        `json:"followingCounter"`
	PhotoCounter     int        `json:"photoCounter"`

	Photos []Photo `json:"photos"` // list of photos
	// BannedUsers    []UserName `json:"bannedUsers"`
}

type User struct {
	UserId   Identifier `json:"userId"`
	Username string     `json:"username"`
}

type Photo struct {
	PhotoId            Identifier `json:"photoId"`
	UploaderUserId     Identifier `json:"uploaderUserId"`
	LikeCounter        int        `json:"likeCounter"`
	Comments           []Comment  `json:"comments"`           // How to manage this?
	LikedByCurrentUser bool       `json:"likedByCurrentUser"` // FIXME seems like a problem
	Date               string     `json:"date"`
	PhotoPath          string     `json:"photoPath"` // in openapi this is represented as photofile
}

// HANDLE THIS TAGS
type Comment struct {
	CommentId        Identifier `json:"commentId"`
	CommentingUserId Identifier `json:"commentingUserId"` // commenter id
	PhotoId          Identifier `json:"photoId"`
	Body             string
	Date             string `json:"commentDate"` // FIXME this should be changed in commentdate
	//TODO manage others
}

type BodyRequest struct {
	Body string `json:"body"`
}

// TODO did i use this?
// type StreamInfo struct {
// 	PhotoPaths []string
// 	Date       string
// }
