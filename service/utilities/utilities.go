package Utilities

import (
	"math/rand"
	"github.com/NeoSnakex34/WasaPhoto/service/structs"
)

// as stated in api.yaml the identifier is a string of lenght 11 @X000000000
func GenerateIdentifier (actor char) structs.Identifier {
	const len = 9 
	const validChars = "0123456789"
	var actorChar;

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
		//TODO handle
	}
	

	rand.Seed(time.Now().UnixNano())

	
}