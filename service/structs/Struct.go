package structs

type Identifier struct {
	Id string `json:"identifier"`
}

type UserName struct {
	PlainUsername string `json:"username"`
}

type UserProfile struct {
	UserId           Identifier `json:"userId"`
	Username         UserName   `json:"username"`
	FollowerCounter  uint32     `json:"followersCounter"`
	FollowingCounter uint32     `json:"followingCounter"`
	PhotoCounter     uint32     `json:"photoCounter"`
	//TODO absolutely manage this in openapi
	UploadedPhotos []Photo    `json:"uploadedPhotos"`
	BannedUsers    []UserName `json:"bannedUsers"`
}

type User struct {
	UserId   Identifier `json:"userId"`
	Username UserName   `json:"username"`

	//TODO manage following

}

// type PhotoFile struct {
// 	PhotoByteStream []byte `json:"photoByteStream"`
// }

type Photo struct {
	PhotoId  Identifier `json:"photoId"`
	UserId   Identifier `json:"userId"`
	Like     int        `json:"like"`
	Comments []Comment  `json:"comments"`

	//TODO consider manage type of date time here or convert/parse it later from string and viceversa
	//kinda prone to errors so would be preferred to use a specific type here probably, maybe using a struct
	//with a string and a time.Time field
	Date       string `json:"date"`
	PhotoBytes []byte `json:"photoBytes"`
}

type Comment struct {
	CommentId Identifier `json:"commentId"`
	UserId    Identifier `json:"userId"`
	PhotoId   Identifier `json:"photoId"`
	Body      string     `json:"body"`
	Date      string     `json:"date"`
	//TODO manage others
}
