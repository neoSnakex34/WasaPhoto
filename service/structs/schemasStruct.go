package structs

type Identifier struct {
	Id string `json:"identifier"`
}

type UserProfile struct {
	userId       string `json:"userId"`
	username     string `json:"username"`
	followers    uint32 `json:"followers"`
	following    uint32 `json:"following"`
	photoCounter uint32 `json:"photoCounter"`
}

type User struct {
	userId      string `json:"userId"`
	username    string `json:"username"`
	followers   []User `json:"followers"`
	bannedUsers []User `json:"bannedUsers"`
	//TODO manage following

}

type PhotoFile struct {
	photoByteStream []byte `json:"photoByteStream"`
}
type Photo struct {
	photoId  string    `json:"photoId"`
	userId   string    `json:"userId"`
	like     int       `json:"like"`
	comments []Comment `json:"comments"`

	//TODO consider manage type of date time here or convert/parse it later from string and viceversa
	//kinda prone to errors so would be preferred to use a specific type here probably, maybe using a struct
	//with a string and a time.Time field
	date       string    `json:"date"`
	photoBytes PhotoFile `json:"photoBytes"`
}

type Comment struct {
	commentId string `json:"commentId"`
	userId    string `json:"userId"`
	body      string `json:"body"`
	date      string `json:"date"`
	//TODO manage others
}
