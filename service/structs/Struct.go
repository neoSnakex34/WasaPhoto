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
	PhotoId Identifier `json:"photoId"`
	UserId  Identifier `json:"userId"`
	// Like     int        `json:"like"`
	// Comments []Comment  `json:"comments"`

	//TODO consider manage type of date time here or convert/parse it later from string and viceversa
	//kinda prone to errors so would be preferred to use a specific type here probably, maybe using a struct
	//with a string and a time.Time field
	Date      string `json:"date"`
	PhotoPath string `json:"photoPath"` // FIXME May be incosistent with openapi
}

// HANDLE THIS TAGS
type Comment struct {
	CommentId Identifier `json:"commentId"`
	UserId    Identifier `json:"userId"` // commenter id
	PhotoId   Identifier `json:"photoId"`
	Body      string
	Date      string `json:"date"`
	//TODO manage others
}

type BodyRequest struct {
	Body string `json:"body"`
}

type StreamInfo struct {
	PhotoPaths []string
	Date       string
}
