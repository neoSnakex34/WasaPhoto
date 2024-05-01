package database

import (
	"database/sql"
	"errors"
	"log"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	serviceutilities "github.com/neoSnakex34/WasaPhoto/service/service-utilities"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (db *appdbimpl) DoLogin(username string) (structs.Identifier, error) {

	var userId string
	idIsValid := false

	exist, userId, err := db.checkUserExists(username)
	// if any error is found i return it (TODO handle)
	if err != nil {
		// check if you need to throw the login error or not
		return structs.Identifier{}, err
	}

	// else if the user exist i have to login
	if exist {
		// login
		log.Println("user exist")
		return structs.Identifier{Id: userId}, nil

	} else if !exist {

		// loop until a valid user or error is found
		for (!idIsValid) && (err == nil) {
			idIsValid, err = db.validId(userId, "U")

			tmpId, _ := generateIdentifier("U") // here error can be ignored since we are automatically using a valid actor

			userId = tmpId.Id

		}

		if err != nil {
			return structs.Identifier{}, err
		}
		// here i actually create the user by setting is username in N mode
		// setting username for the first time is part of the action of generating the userId
		// that it has been verified in the for (while) loop on line 38
		log.Println(userId, " ", username)

		err = db.SetMyUserName(username, userId, "N")

		if err != nil {
			return structs.Identifier{}, err
		}
	}

	return structs.Identifier{Id: userId}, nil

}

func (db *appdbimpl) SetMyUserName(newUsername string, userId string, mode string) error {

	var count int
	var valid bool = false

	//  if user is new one MODE = N i need to do inser
	// if user is already signed MODE = U i need to update by id

	//  i check if newUsername is taken
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, newUsername).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 { // cannot check with err (count will always return something) i just checked err for the query
		return customErrors.ErrAlreadyTakenUsername
	}

	matched := serviceutilities.CheckRegexNewUsername(newUsername)

	if count == 0 && matched {

		log.Printf("username [%s] is valid\n", newUsername)
		valid = true

	} else if !matched {
		return customErrors.ErrInvalidRegexUsername
	}

	if valid {

		switch mode {

		case "N":

			err := db.createUser(newUsername, userId)
			return err

		case "U":

			log.Println("updating username")
			_, err := db.c.Exec(`UPDATE users SET username = ? WHERE userId = ?`, newUsername, userId)
			return err

		default:
			return customErrors.ErrInvalidIdentifierMode
		}

	}

	return nil
}

// probably i can add a refreshUserProfile function that updates all the counters
func (db *appdbimpl) GetUserProfile(profileUserId structs.Identifier, requestorUserId structs.Identifier) (structs.UserProfile, error) {

	println("profileUserId, requestorUserId: ", profileUserId.Id, requestorUserId.Id)
	plainUserId := profileUserId.Id
	plainRequestorUserId := requestorUserId.Id
	// check requestor banned by profile user
	err := db.checkBan(plainUserId, plainRequestorUserId)

	if errors.Is(err, customErrors.ErrIsBanned) {
		log.Println("requestor is banned by user")
		return structs.UserProfile{}, err
	} else if err != nil {
		return structs.UserProfile{}, err
	}

	var username string
	var followerCounter int
	var followingCounter int
	var photoCounter int

	// queries func
	username, err = db.getUsernameByUserId(plainUserId)
	if err != nil {
		return structs.UserProfile{}, err
	}

	// follower count
	followerCounter, err = db.getFollowersCounterByUserId(plainUserId)
	if err != nil {
		return structs.UserProfile{}, err
	}

	// following count
	followingCounter, err = db.getFollowingCounterByUserId(plainUserId)
	if err != nil {
		return structs.UserProfile{}, err
	}

	// photo count via os directory counter
	photoCounter, photos, err := db.getPhotosAndInfoByUserId(plainUserId, plainRequestorUserId)
	if err != nil {
		return structs.UserProfile{}, err
	}

	profileRetrieved := structs.UserProfile{
		UserId:           profileUserId,
		Username:         username,
		FollowerCounter:  followerCounter,
		FollowingCounter: followingCounter,
		PhotoCounter:     photoCounter,

		Photos: photos,
	}
	// log.Println("THE PROFILE: ", profileRetrieved)
	return profileRetrieved, nil
}

// is it a security vulnerability? even though no real auth is performed nor needed here
func (db *appdbimpl) GetUserList(requestorUserId structs.Identifier) ([]structs.UserFromQuery, error) {

	var userFromQueryList []structs.UserFromQuery
	var userId string
	var username string
	var isRequestorBanned bool
	var requestorHasBanned bool
	var requestorHasFollowed bool

	rows, err := db.c.Query(`SELECT userId, username FROM users WHERE userId != ?`, requestorUserId.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&userId, &username)
		if err != nil {
			return nil, err
		}

		err = db.checkBan(userId, requestorUserId.Id)
		println("err: ", err)
		if errors.Is(err, customErrors.ErrIsBanned) {
			isRequestorBanned = true
		} else if err != nil {
			return nil, err
		} else {
			isRequestorBanned = false
		}

		err = db.checkBan(requestorUserId.Id, userId)
		if errors.Is(err, customErrors.ErrIsBanned) {
			requestorHasBanned = true
		} else if err != nil {
			return nil, err
		} else {
			requestorHasBanned = false
		}

		requestorHasFollowed, err = db.follows(requestorUserId.Id, userId)
		if err != nil {
			return nil, err
		}

		userFromQueryList = append(userFromQueryList, structs.UserFromQuery{
			User: structs.User{
				UserId:   structs.Identifier{Id: userId},
				Username: username,
			},
			IsRequestorBanned:    isRequestorBanned,
			RequestorHasBanned:   requestorHasBanned,
			RequestorHasFollowed: requestorHasFollowed,
		})
	}

	return userFromQueryList, nil
}

func (db *appdbimpl) GetMyStream(userId structs.Identifier) ([]structs.Photo, error) {

	// first i obtain a followerlist
	var followerIdList []string
	followerIdList, err := db.getFollowerList(userId.Id)
	plainRequestorId := userId.Id

	println("followerIdList: ", followerIdList)

	if err != nil {
		return nil, err
	}

	streamOfPhotoStructs, err := db.getStreamOfPhotos(followerIdList, plainRequestorId)
	if err != nil {
		return nil, err
	}

	// photo struct will be returned, in order to use it in frontend i need to parse the path
	// and retrieve the actual photo
	return streamOfPhotoStructs, nil
}

// ========== private functions from here
func (db *appdbimpl) createUser(username string, userId string) error {

	_, err := db.c.Exec("INSERT INTO users (username, userId) VALUES (?, ?)", username, userId)

	log.Println("user created with error returned: ", err)
	return err
}

func (db *appdbimpl) checkUserExists(username string) (bool, string, error) {
	var userInTable bool
	// var userId structs.Identifier = structs.Identifier{}
	var id string
	// first we check if user is in the database querying his row (given that username is unique)
	err := db.c.QueryRow(`SELECT userId FROM users WHERE username = ?`, username).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		userInTable = false
		err = nil // else it will fail control in next function, very important to be checked !

	} else if err != nil {
		return false, id, err
	} else {
		// so the user exist
		// todo this could be prone to bugs, if something goes south check it out
		userInTable = true
		// userId = structs.Identifier{Id: id}
	}

	return userInTable, id, nil
}
