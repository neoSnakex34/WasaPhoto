package database

// FIXME most likely everything inside a struct must be unpacked to be used successfully in queries, i have to handle it
// TODO the fixme just above was checked and now everything i use except return statement in dologin is a plain string (instead of a wrapper struct)
// doing such thing is useful in queries but since those structs are useful
import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"

	serviceutilities "github.com/neoSnakex34/WasaPhoto/service/api/service-utilities"
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// TODO i dont think that using errors to mask others is a good idea in debugging
// implement those only if you did enough testing
// var LoginError = errors.New("an error occured during login")
// FIXME log error in console in order to unmask badrequest ecc in apis
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

			// println("id: ", userId)

			// TODO warning with this assignation, it could break everything
			tmpId, _ := GenerateIdentifier("U") // here error can be ignored since we are automatically using a valid actor

			// println("tmpId: ", tmpId.Id)

			userId = tmpId.Id

			// println("userId: ", userId)
		}

		if err != nil {
			return structs.Identifier{}, err
		}
		// here i actually create the user by setting is username in N mode
		// setting username for the first time is part of the action of generating the userId
		// that it has been verified in the for (while) loop on line 38
		log.Println(userId, " ", username)

		// [ ] GIVEN that here will be called setMyUsername regex check will be done after generating id, it is not slow but neither is
		// efficient or clean
		// i should modify that
		err = db.SetMyUserName(username, userId, "N")
		// added this check cause it would login with invalid id in frontend
		// resulting in a brocken  backend
		if err != nil {
			return structs.Identifier{}, err
		}
	}

	return structs.Identifier{Id: userId}, nil

}

func (db *appdbimpl) SetMyUserName(newUsername string, userId string, mode string) error {

	// TODO all this cheks must be lowercase, also i need an efficient way to loop over errors till
	// a valid name pops up

	var count int
	var valid bool = false

	//  if user is new one MODE = N i need to do inser
	// if user is already signed MODE = U i need to update by id

	//  i check if newUsername is taken
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, newUsername).Scan(&count)

	matched := serviceutilities.CheckRegexNewUsername(newUsername)

	if count == 0 && matched {

		log.Printf("username [%s] is valid\n", newUsername)
		valid = true

	} else {
		if !matched {
			err = customErrors.ErrInvalidRegexUsername
			log.Printf("username [%s] is valid\n", newUsername)
			return err
		}

		// if any other error occurred i return it
		if err != nil {
			return err
		}
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
			// FIXME add custom error
			return errors.New("error in parsing mode or invalid mode for userame operation")

		}

	}

	return err
}

// probably i can add a refreshUserProfile function that updates all the counters
func (db *appdbimpl) GetUserProfile(profileUserId structs.Identifier, requestorUserId structs.Identifier) (structs.UserProfile, error) {

	println("profileUserId, requestorUserId: ", profileUserId.Id, requestorUserId.Id)
	plainUserId := profileUserId.Id
	plainRequestorUserId := requestorUserId.Id
	// check requestor banned by profile user
	err := db.checkBan(plainUserId, plainRequestorUserId)

	if errors.Is(err, customErrors.ErrIsBanned) {
		log.Println("requestor is banned by user") // TODO log this in api
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

func (db *appdbimpl) GetUserList() ([]structs.User, error) {

	var userList []structs.User

	rows, err := db.c.Query(`SELECT userId, username FROM users`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userId string
		var username string
		err = rows.Scan(&userId, &username)
		if err != nil {
			return nil, err
		}

		userList = append(userList, structs.User{UserId: structs.Identifier{Id: userId}, Username: username})
	}

	return userList, nil
}

// TODO sort photosbydate
func (db *appdbimpl) GetMyStream(userId structs.Identifier) ([]structs.Photo, error) {

	// first i obtain a followerlist
	var followerIdList []string
	followerIdList, err := db.getFollowerList(userId.Id)

	println("followerIdList: ", followerIdList)

	if err != nil {
		return nil, err
	}

	streamOfPhotoStructs, err := db.getSortedStreamOfPhotos(followerIdList)
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

// TODO generalize for any id MAYBE
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

func (db *appdbimpl) getUsernameByUserId(plainUserId string) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE userId = ?`, plainUserId).Scan(&username)
	return username, err
}

func (db *appdbimpl) getFollowersCounterByUserId(plainUserId string) (int, error) {
	var followers int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followedId = ?`, plainUserId).Scan(&followers)
	return followers, err
}

func (db *appdbimpl) getFollowingCounterByUserId(plainUserId string) (int, error) {
	var following int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followerId = ?`, plainUserId).Scan(&following)
	return following, err
}

// FIXME let it return a slice of Photo with metadatas when retrieving profile in frontend
func (db *appdbimpl) getPhotosAndInfoByUserId(plainUserId string, plainRequestorUserId string) (int, []structs.Photo, error) {

	path := Folder + plainUserId + "/"
	photoFsDirs, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		log.Println("folder not found or does not exist counters set to 0")
		return 0, nil, nil
	} else if err != nil {
		return 0, nil, err
	}

	photoCount := len(photoFsDirs)
	var photoName string
	var plainPhotoId string
	var photoDate string
	var likeCounter int
	var liked bool
	var photoPath string

	var photos []structs.Photo

	var tmpPhoto structs.Photo

	// for each photo in the folder i get the metadata
	// and extract the path to the photo
	// absolute path could be retrieved via db but i need the partial one
	// to be used in frontend
	for _, photo := range photoFsDirs {

		photoName = photo.Name()
		plainPhotoId = strings.Split(photo.Name(), ".")[0]

		// partial photo path
		photoPath = plainUserId + "/" + photoName
		// photoPathList = append(photoPathList, photoPath)

		photoDate, err = db.getPhotoDateByPhotoId(plainPhotoId)
		if err != nil {
			log.Println("error in getting photo date from db")
			return 0, nil, err
		}

		likeCounter, err = db.getNumberOfLikedByPhotoId(plainPhotoId)
		if err != nil {
			log.Println("error in getting like counter from db")
			return 0, nil, err
		}

		liked, err = db.getLikedByUserId(plainRequestorUserId, plainPhotoId)
		if err != nil {
			log.Println("error in getting info of like by user from db")
			return 0, nil, err
		}

		// FIXME i have to manage comments

		comments := []structs.Comment{}

		tmpPhoto = structs.Photo{
			PhotoId:            structs.Identifier{Id: plainPhotoId},
			UploaderUserId:     structs.Identifier{Id: plainUserId},
			LikeCounter:        likeCounter,
			Comments:           comments,
			LikedByCurrentUser: liked,
			Date:               photoDate,
			PhotoPath:          photoPath,
		}

		photos = append(photos, tmpPhoto)

	}

	return photoCount, photos, nil
}
