package database

// [x] generaliziong validId may have led to a series of bugs in users.go, is anything is broken go check dependencies
// between the two files

import (
	"math/rand"
	"time"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// as stated in api.yaml the identifier is a string of lenght 11 @X000000000
// actor will be mode of the id (U P C)
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
	println("id: ", id)
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

func (db *appdbimpl) getUploaderByPhotoId(photoId structs.Identifier) (structs.Identifier, error) {
	var plainUploaderId string
	println("photoId ", photoId.Id)

	err := db.c.QueryRow(`SELECT userId FROM photos WHERE photoId = ?`, photoId.Id).Scan(&plainUploaderId)
	return structs.Identifier{Id: plainUploaderId}, err
}

func (db *appdbimpl) getCommenterByCommentId(commentId structs.Identifier) (structs.Identifier, error) {
	var plainCommenterId string
	err := db.c.QueryRow(`SELECT userId FROM comments WHERE commentId = ?`, commentId.Id).Scan(&plainCommenterId)
	return structs.Identifier{Id: plainCommenterId}, err
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
