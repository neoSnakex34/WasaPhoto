package database

import (
	"time"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// IMPORTANT likeid = userid of linkingUser;

// ======= verification operations

// ======= comments operations
func (db *appdbimpl) CommentPhoto(commentedPhotoId structs.Identifier, requestorUserId structs.Identifier, body string) (structs.Comment, error) {

	var isValidId bool = false
	// photoId and userId are already verified when firstly created, note that unmasking the use of a function like this
	// may led to some serious bugs if someone manages to use CommentPhoto with an invalid id
	var newCommentId structs.Identifier
	var err error

	for isValidId == false && err == nil {

		newCommentId, err = GenerateIdentifier("C")
		isValidId, err = db.validId(newCommentId.Id, "C")

	}

	if err != nil {

		return structs.Comment{}, err

	}

	// TODO keep the date inference after the validId loop
	commentDate := time.Now().Format(time.RFC3339)

	newComment := structs.Comment{
		CommentId: newCommentId,
		UserId:    requestorUserId,
		PhotoId:   commentedPhotoId,
		Body:      body,
		Date:      commentDate,
	}
	return newComment, nil
}
