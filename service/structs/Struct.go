package structs

type Identifier struct {
	Id string `json:"identifier"`
}

type UserName struct {
	plainUsername string `json:"username"`
}

type UserProfile struct {
	userId           Identifier `json:"userId"`
	username         UserName   `json:"username"`
	followerCounter  uint32     `json:"followersCounter"`
	followingCounter uint32     `json:"followingCounter"`
	photoCounter     uint32     `json:"photoCounter"`
	//TODO absolutely manage this
	bannedUsers []UserName `json:"bannedUsers"`
}

type User struct {
	userId   Identifier `json:"userId"`
	username UserName   `json:"username"`

	//TODO manage following

}

type PhotoFile struct {
	photoByteStream []byte `json:"photoByteStream"`
}

type Photo struct {
	photoId  Identifier `json:"photoId"`
	userId   Identifier `json:"userId"`
	like     int        `json:"like"`
	comments []Comment  `json:"comments"`

	//TODO consider manage type of date time here or convert/parse it later from string and viceversa
	//kinda prone to errors so would be preferred to use a specific type here probably, maybe using a struct
	//with a string and a time.Time field
	date       string    `json:"date"`
	photoBytes PhotoFile `json:"photoBytes"`
}

type Comment struct {
	commentId Identifier `json:"commentId"`
	userId    Identifier `json:"userId"`
	body      string     `json:"body"`
	date      string     `json:"date"`
	//TODO manage others
}
