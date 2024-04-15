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
	//TODO absolutely manage this in openapi
	UploadedPhotos []string `json:"uploadedPhotos"` // list of paths to photos
	// BannedUsers    []UserName `json:"bannedUsers"`
}

type User struct {
	UserId   Identifier `json:"userId"`
	Username string     `json:"username"`

	//TODO manage following

}

// type PhotoFile struct {
// 	PhotoByteStream []byte `json:"photoByteStream"`
// }

type Photo struct {
	PhotoId        Identifier `json:"photoId"`
	UploaderUserId Identifier `json:"uploaderUserId"`
	// Like     int        `json:"like"`
	// Comments []Comment  `json:"comments"` // How to manage this?

	Date      string `json:"date"`
	PhotoPath string `json:"photoPath"` // in openapi this is represented as photofile
}

// HANDLE THIS TAGS
type Comment struct {
	CommentId        Identifier `json:"commentId"`
	CommentingUserId Identifier `json:"commentingUserId"` // commenter id
	PhotoId          Identifier `json:"photoId"`
	Body             string
	Date             string `json:"date"`
	//TODO manage others
}

type BodyRequest struct {
	Body string `json:"body"`
}

type StreamInfo struct {
	PhotoPaths []string
	Date       string
}
