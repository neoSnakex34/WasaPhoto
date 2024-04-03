package Utilities

import (
	"math/rand"
	"github.com/NeoSnakex34/WasaPhoto/service/structs"
)

// as stated in api.yaml the identifier is a string of lenght 11 @X000000000
func GenerateIdentifier (actor char) (structs.Identifier, nil) {
	const len = 9 
	const validChars = "0123456789"
	var actorChar string;

	switch actor {
		case 'U':
			actorChar = "U"
		case 'C':
			actorChar = "C"
		case 'P':
			actorChar = "P"
		default:
			actorChar = nil 
	}

	if actorChar == nil {
		return nil //TODO handle where needed to be handled 
	}

	rand.Seed(time.Now().UnixNano())

	// had a look online for this, check if it can be improved 
	randomChunk := make([]byte, len)
	for i := range randomChunk {
		randomChunk[i] = validChars[rand.Intn(len(validChars))]
	}

	randomStringChunk := string(randomChunk)
	
	generatedId := Identifier[ identifier: '@'+ actorChar + randomStringChunk]
	
	return generatedId
	
	
}