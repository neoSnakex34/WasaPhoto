package database

// TODO generaliziong validId may have led to a series of bugs in users.go, is anything is broken go check dependencies
// between the two files

import (
	"errors"
	"math/rand"
	"time"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

var invalidIdMode = errors.New("invalid id mode")
var invalidId = errors.New("invalid id")

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
		return structs.Identifier{}, errors.New("Provided invalid actor type string") //TODO handle where needed to be handled
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

// mode can be U P or C any other is invalid (capital letters only)
func (db *appdbimpl) validId(id string, mode string) (bool, error) {
	// FIXME since single char are unicode byte
	// even if i am sure that those are utf8 1byte chars
	// it is probably better to check them using the appropriate comparator

	var idMode string = string(id[1])

	// first check is modecheck is mode is matched we proceed else we abort
	if idMode != mode {
		return false, invalidIdMode
	}

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
		err = invalidIdMode
	}

	if err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}

	return false, invalidId

}
