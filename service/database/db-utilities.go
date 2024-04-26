package database

// [x] generaliziong validId may have led to a series of bugs in users.go, is anything is broken go check dependencies
// between the two files

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// as stated in api.yaml the identifier is a string of lenght 11 @X000000000
// actor will be mode of the id (U P C)
// not sure if this needs to be exported
// FIXME ^^^^
func GenerateIdentifier(actor string) (structs.Identifier, error) {
	const lenght = 9
	const validChars = "0123456789"
	var actorChar string

	switch actor {
	case "U":
		actorChar = "U"
	case "C":
		actorChar = "C"
	case "P":
		actorChar = "P"
	default:
		actorChar = "E"
	}

	if actorChar == "E" {
		return structs.Identifier{}, customErrors.ErrInvalidIdMode
	}

	rand.Seed(time.Now().UnixNano())

	// had a look online for this, check if it can be improved
	randomChunk := make([]byte, lenght)
	for i := range randomChunk {
		randomChunk[i] = validChars[rand.Intn(len(validChars))]
	}

	randomStringChunk := string(randomChunk)

	generatedId := structs.Identifier{Id: "@" + actorChar + randomStringChunk}

	return generatedId, nil

}

// when building core functionality decide it
func (db *appdbimpl) checkBan(bannerId string, bannedId string) error {
	var counter int

	err := db.c.QueryRow(`SELECT COUNT(*) FROM bans WHERE bannerId = ? AND bannedId = ?`, bannerId, bannedId).Scan(&counter)

	if err != nil {
		return err
	} else if counter == 0 {
		return nil
	} else if counter > 0 {
		return customErrors.ErrIsBanned
	}
	return nil
}

// mode can be U P or C any other is invalid (capital letters only)
func (db *appdbimpl) validId(id string, mode string) (bool, error) {
	// FIXME since single char are unicode byte
	// even if i am sure that those are utf8 1byte chars
	// it is probably better to check them using the appropriate comparator

	// FIXME this is a check thhat works only if the id already exists
	// if mode != "N" {
	// 	var idMode string = string(id[1])

	// 	// first check is modecheck is mode is matched we proceed else we abort
	// 	if idMode != mode {
	// 		return false, invalidIdMode
	// 	}

	// }
	var count int
	var err error = nil

	// here we check if the id is present in the table for the respective mode
	switch mode {
	case "U":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE userId = ?`, id).Scan(&count)

	case "P":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM photos WHERE photoId = ?`, id).Scan(&count)

	case "C":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE commentId = ?`, id).Scan(&count)

	default:
		err = customErrors.ErrInvalidIdMode
	}

	if err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}

	return false, customErrors.ErrInvalidId
}

// TODO important
// FIXME pass those as plain ids instead
func (db *appdbimpl) getUploaderByPhotoId(photoId structs.Identifier) (structs.Identifier, error) {
	var plainUploaderId string

	err := db.c.QueryRow(`SELECT userId FROM photos WHERE photoId = ?`, photoId.Id).Scan(&plainUploaderId)
	return structs.Identifier{Id: plainUploaderId}, err
}

func (db *appdbimpl) getCommenterIdByCommentId(commentId structs.Identifier) (structs.Identifier, error) {
	var plainCommenterId string
	err := db.c.QueryRow(`SELECT userId FROM comments WHERE commentId = ?`, commentId.Id).Scan(&plainCommenterId)
	if err != nil {
		return structs.Identifier{}, err // MANAGE
	}
	return structs.Identifier{Id: plainCommenterId}, nil
}

func (db *appdbimpl) getCommenterUsernameByCommentingId(plainCommenterId string) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE userId = ?`, plainCommenterId).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (db *appdbimpl) getUploaderByCommentId(commentId structs.Identifier) (structs.Identifier, error) {
	var plainPhotoId string
	err := db.c.QueryRow(`SELECT photoId FROM comments WHERE commentId = ?`, commentId.Id).Scan(&plainPhotoId)
	if err != nil {
		return structs.Identifier{}, err
	}
	var plainUploaderId string
	err = db.c.QueryRow(`SELECT userId FROM photos WHERE photoId = ?`, plainPhotoId).Scan(&plainUploaderId)

	return structs.Identifier{Id: plainUploaderId}, err
}

func (db *appdbimpl) alreadyLiked(plainRequestorUserId string, likedPhotoId string) (bool, error) {
	var counter int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE likerId = ? AND photoId = ?`, plainRequestorUserId, likedPhotoId).Scan(&counter)
	if err != nil {
		return false, err
	} else if counter > 0 {
		return true, customErrors.ErrPhotoAlreadyLikedByUser
	}
	return false, nil

}

func (db *appdbimpl) getPhotosByUploaderId(plainUploaderId string) ([]structs.Photo, error) {
	var photos []structs.Photo
	var userId string = plainUploaderId
	var photoId string
	var date string
	var photoPath string

	// query to retrieve info
	rows, err := db.c.Query(`SELECT photoId, userId, date, photoPath FROM photos WHERE userId = ?`, plainUploaderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// user id is the same for all photos at every call, but for now is re assigned
		err = rows.Scan(&photoId, &userId, &date, &photoPath)
		if err != nil {
			return nil, err
		}
		photo := structs.Photo{
			PhotoId:        structs.Identifier{Id: photoId},
			UploaderUserId: structs.Identifier{Id: userId},
			Date:           date,
			PhotoPath:      photoPath,
		}
		photos = append(photos, photo)
	}

	return photos, err
}

func (db *appdbimpl) getPhotoDateByPhotoId(plainPhotoId string) (string, error) {
	var date string

	err := db.c.QueryRow(`SELECT date FROM photos WHERE photoId = ?`, plainPhotoId).Scan(&date)
	// FIXME manage errors
	// TODO change
	// error no rows
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return date, nil
}

func (db *appdbimpl) getNumberOfLikedByPhotoId(plainPhotoId string) (int, error) {
	var likeCounter int
	err := db.c.QueryRow(`SELECT COUNT(likerId) FROM likes WHERE photoId = ?`, plainPhotoId).Scan(&likeCounter)
	// FIXME manage custom error here
	if err != nil {
		return 0, err
	}
	return likeCounter, nil
}

func (db *appdbimpl) getLikedByUserId(plainUserId string, plainphotoId string) (bool, error) {
	var counter int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE likerId = ? AND photoId = ?`, plainUserId, plainphotoId).Scan(&counter)
	if err != nil {
		println("an error occured while checking photolikedbycurrentid")
		return false, err
	}
	return counter > 0, nil
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

		comments, err := db.getCommentsByPhotoId(plainPhotoId)
		// TODO manage errors
		if err != nil {
			log.Println("error in getting comments from db")
			return 0, nil, err
		}

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

	// logrus.Info("comments: ", photos[0].Comments)
	return photoCount, photos, nil
}

// id will be plain cause it is passed as plain
func (db *appdbimpl) getCommentsByPhotoId(plainPhotoId string) ([]structs.Comment, error) {
	var Comments []structs.Comment

	rows, err := db.c.Query(`SELECT commentId, userId, body, date FROM comments WHERE photoId = ?`, plainPhotoId)
	// TODO manage errors
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var plainCommentId string
		var plainUserId string
		var body string
		var date string

		err = rows.Scan(&plainCommentId, &plainUserId, &body, &date)
		if err != nil {
			return nil, err
		}

		commentingUsername, err := db.getCommenterUsernameByCommentingId(plainUserId)
		if err != nil {
			return nil, err
		}

		comment := structs.Comment{
			CommentId:          structs.Identifier{Id: plainCommentId},
			CommentingUserId:   structs.Identifier{Id: plainUserId},
			CommentingUsername: commentingUsername,
			PhotoId:            structs.Identifier{Id: plainPhotoId},
			Body:               body,
			Date:               date,
		}

		Comments = append(Comments, comment)

	}

	return Comments, nil

}
