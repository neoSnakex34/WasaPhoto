package database

import (
	"errors"
	"time"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// IMPORTANT likeid = userid of linkingUser;

// ======= comments operations
// FIXME add comments to comments list in photo struct
func (db *appdbimpl) CommentPhoto(commentedPhotoId structs.Identifier, requestorUserId structs.Identifier, body string) (structs.Comment, error) {

	var isValidId bool = false
	// photoId and userId are already verified when firstly created, note that unmasking the use of a function like this
	// may lead to some serious bugs if someone manages to use CommentPhoto with an invalid id
	var newCommentId structs.Identifier
	var err error

	// check ban
	userUploaderId, err := db.getUploaderByPhotoId(commentedPhotoId)
	if err != nil {
		return structs.Comment{}, err
	}

	// check ban
	err = db.checkBan(userUploaderId.Id, requestorUserId.Id)
	if errors.Is(err, customErrors.ErrIsBanned) {
		return structs.Comment{}, err
	} else if err != nil {
		return structs.Comment{}, err
	}

	for !isValidId && err == nil {

		newCommentId, err = GenerateIdentifier("C")
		if err != nil {
			return structs.Comment{}, err
		}
		isValidId, err = db.validId(newCommentId.Id, "C")

	}

	if err != nil {

		return structs.Comment{}, err

	}

	// TODO keep the date inference after the validId loop
	commentDate := time.Now().UTC().Format(time.RFC3339)

	err = db.addComment(newCommentId.Id, requestorUserId.Id, commentedPhotoId.Id, body, commentDate)
	// insert comment in db

	if err != nil {
		return structs.Comment{}, err
	}

	newComment := structs.Comment{
		CommentId:        newCommentId,
		CommentingUserId: requestorUserId,
		PhotoId:          commentedPhotoId,
		Body:             body,
		Date:             commentDate,
	}
	return newComment, nil
}

// TODO make this an external function removeComment to maintain the consistency
func (db *appdbimpl) UncommentPhoto(commentId structs.Identifier) error {

	// get uplader id
	uploaderId, err := db.getUploaderByCommentId(commentId)
	if err != nil {
		return err
	}

	commenterId, err := db.getCommenterByCommentId(commentId)
	if err != nil {
		return err
	}

	// check ban
	err = db.checkBan(uploaderId.Id, commenterId.Id)
	if errors.Is(err, customErrors.ErrIsBanned) {
		return err // TODO log
	} else if err != nil {
		return err
	}

	err = db.removeComment(commentId.Id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) addComment(commentId string, userId string, photoId string, body string, date string) error {
	_, err := db.c.Exec(`INSERT INTO comments (commentId, userId, photoId, body, date) VALUES (?, ?, ?, ?, ?)`, commentId, userId, photoId, body, date)
	return err
}

func (db *appdbimpl) removeComment(commentId string) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE commentId = ?`, commentId)
	return err
}

// ======= likes operations
func (db *appdbimpl) LikePhoto(userId structs.Identifier, photoId structs.Identifier) error {

	var photoIsLiked bool
	var err error

	// retrieve uploader of the photo
	uploaderId, err := db.getUploaderByPhotoId(photoId)
	if err != nil {
		return err
	}

	// check ban
	err = db.checkBan(uploaderId.Id, userId.Id)
	// FIXME try to make this a lil more fancy...
	if errors.Is(err, customErrors.ErrIsBanned) {
		return err
	} else if err != nil {
		return err
	}

	photoIsLiked, err = db.alreadyLiked(userId.Id, photoId.Id)
	println("photo is liked: ", photoIsLiked)
	if err == nil && !photoIsLiked {
		// TODO add like
		err = db.addLike(userId.Id, photoId.Id)
		// if err now is nil will be returned nil at the end of the function
		// TODO check this really happens
	}
	return err
}

func (db *appdbimpl) UnlikePhoto(userId structs.Identifier, photoId structs.Identifier) error {
	var err error

	// retrieve uploader of the photo
	uploaderId, err := db.getUploaderByPhotoId(photoId)
	if err != nil {
		return err
	}

	// check ban
	err = db.checkBan(uploaderId.Id, userId.Id)
	if errors.Is(err, customErrors.ErrIsBanned) {
		return err
	} else if err != nil {
		return err
	}

	liked, err := db.alreadyLiked(userId.Id, photoId.Id)
	if errors.Is(err, customErrors.ErrPhotoAlreadyLikedByUser) && liked {

		err = db.removeLike(userId.Id, photoId.Id)
	} else if err != nil {
		return err
	} else if !liked {
		return customErrors.ErrPhotoNotLikedByUser
	}

	return err
}
func (db *appdbimpl) addLike(requestorUserId string, likedPhotoId string) error {
	_, err := db.c.Exec(`INSERT INTO likes (likerId, photoId) VALUES (?, ?)`, requestorUserId, likedPhotoId)
	return err
}

func (db *appdbimpl) removeLike(requestorUserId string, likedPhotoId string) error {
	_, err := db.c.Exec(`DELETE FROM likes WHERE likerId = ? AND photoId = ?`, requestorUserId, likedPhotoId)
	return err
}

func (db *appdbimpl) alreadyLiked(requestorUserId string, likedPhotoId string) (bool, error) {
	var counter int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE likerId = ? AND photoId = ?`, requestorUserId, likedPhotoId).Scan(&counter)
	if err != nil {
		return false, err
	} else if counter > 0 {
		return true, customErrors.ErrPhotoAlreadyLikedByUser
	}
	return false, nil

}
