package database

// FIXME most likely everything inside a struct must be unpacked to be used successfully in queries, i have to handle it
// TODO the fixme just above was checked and now everything i use except return statement in dologin is a plain string (instead of a wrapper struct)
// doing such thing is useful in queries but since those structs are useful
import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

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
	println("user exist: ", exist)
	// if any error is found i return it (TODO handle)
	if err != nil {
		// check if you need to throw the login error or not
		return structs.Identifier{}, err
	}

	// else if the user exist i have to login
	if exist {
		// login
		println("user exist!")
		return structs.Identifier{Id: userId}, nil

	} else if !exist {

		// loop until a valid user or error is found
		for (!idIsValid) && (err == nil) {
			println("im here now // debugging")
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
		println("im out")
		// here i actually create the user by setting is username in N mode
		// setting username for the first time is part of the action of generating the userId
		// that it has been verified in the for (while) loop on line 38
		println(userId, " ", username)
		// [ ] GIVEN that here will be called setMyUsername regex check will be done after generating id, it is not slow but neither is
		// efficient or clean
		// i should modify that
		db.SetMyUserName(username, userId, "N")
	}

	return structs.Identifier{Id: userId}, nil

}

func (db *appdbimpl) SetMyUserName(newUsername string, userId string, mode string) error {

	// TODO all this cheks must be lowercase, also i need an efficient way to loop over errors till
	// a valid name pops up

	println("entered setmyusername")

	var count int
	var valid bool = false

	//  if user is new one MODE = N i need to do inser
	// if user is already signed MODE = U i need to update by id

	//  i check if newUsername is taken
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, newUsername).Scan(&count)

	println("err: ", err)
	matched := serviceutilities.CheckRegexNewUsername(newUsername)

	if count == 0 && matched {

		println("username is valid")
		valid = true

	} else {
		if !matched {
			err = customErrors.ErrInvalidRegexUsername
			println("username is not valid", err)
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

			println("updating username")

			_, err := db.c.Exec(`UPDATE users SET username = ? WHERE userId = ?`, newUsername, userId)
			return err

		default:
			return errors.New("error in parsing mode or invalid mode for userame operation")

		}

	}

	return err
}

// TODO getmystream and getmyuserprofile
// FIXME when removing, adding photos and folloerts counters should be updated
// probably i can add a refreshUserProfile function that updates all the counters
func (db *appdbimpl) GetUserProfile(userId structs.Identifier) (structs.UserProfile, error) {

	id := userId.Id
	var username string
	var followerCounter int
	var followingCounter int
	var photoCounter int

	// [ ] add photo list
	var photoList []string

	usernameQuery := `SELECT username FROM users WHERE userId = ?`
	followerCounterQuery := `SELECT COUNT(*) FROM followers WHERE followedId = ?`
	followingCounterQuery := `SELECT COUNT(*) FROM followers WHERE followerId = ?`

	err := db.retriveProfileQueries(usernameQuery, followerCounterQuery, followingCounterQuery)
	if err != nil {
		return structs.UserProfile{}, err
	}

	// photo count via os directory counter
	photoCounter, photoList, err = getPhotoCountAndList(id)
	if err != nil {
		return structs.UserProfile{}, err
	}

	profileRetrieved := structs.UserProfile{
		UserId:           userId,
		Username:         username,
		FollowerCounter:  followerCounter,
		FollowingCounter: followingCounter,
		PhotoCounter:     photoCounter,
		UploadedPhotos:   photoList,
	}

	return profileRetrieved, nil
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
	println("user created")
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

func (db *appdbimpl) retriveProfileQueries(queries ...string) error {
	for _, query := range queries {
		_, err := db.c.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPhotoCountAndList(userId string) (int, []string, error) {
	path := Folder + userId + "/"
	photoFsDirs, err := os.ReadDir(path)
	if err != nil {
		return 0, nil, err
	}

	photoCount := len(photoFsDirs)

	var photoList []string
	for _, photo := range photoFsDirs {
		photoList = append(photoList, filepath.Join(path, photo.Name()))
	}

	return photoCount, photoList, nil
}
