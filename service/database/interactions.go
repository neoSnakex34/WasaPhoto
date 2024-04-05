package database

import (
	"time"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// IMPORTANT likeid = userid of linkingUser;

// ======= verification operations

// ======= comments operations
func (db *appdbimpl) CommentPhoto(photoId structs.Identifier, userId structs.Identifier) (structs.Comment, error) {

	// TODO generate valid commentId

	// TODO keep the date inference after the validId loop

	date := time.Now().Format(time.RFC3339)
	// FIXME just delete this placeholder after finishing the function
	println(date)
	return structs.Comment{}, nil
}
