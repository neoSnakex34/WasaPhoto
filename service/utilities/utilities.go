package utilities

import (
	"errors"
	"math/rand"
	"time"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// as stated in api.yaml the identifier is a string of lenght 11 @X000000000
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
		return Identifier{}, errors.New("Provided invalid actor type string") //TODO handle where needed to be handled
	}

	rand.Seed(time.Now().UnixNano())

	// had a look online for this, check if it can be improved
	randomChunk := make([]byte, lenght)
	for i := range randomChunk {
		randomChunk[i] = validChars[rand.Intn(len(validChars))]
	}

	randomStringChunk := string(randomChunk)

	generatedId := Identifier{identifier: "@" + actorChar + randomStringChunk}

	return generatedId, nil

}
