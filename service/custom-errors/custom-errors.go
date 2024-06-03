package customerrors

import "errors"

var ErrAlreadyTakenUsername = errors.New("username already taken")
var ErrInvalidRegexUsername = errors.New("invalid username")
var ErrCriticDuplicatedId = errors.New("duplicated id multiple actors with the same id")
var ErrAlreadyFollowing = errors.New("already following")
var ErrNotFollowing = errors.New("not following")
var ErrIsBanned = errors.New("cannot interact, user is banned")
var ErrAlreadyBanned = errors.New("already banned")
var ErrNotBanned = errors.New("not banned")
var ErrPhotoDoesNotExist = errors.New("photo does not exist")
var ErrPhotoAlreadyLikedByUser = errors.New("photo already liked by user")
var ErrPhotoNotLikedByUser = errors.New("photo not liked by user")
var ErrInvalidIdMode = errors.New("invalid id mode")
var ErrInvalidId = errors.New("invalid id")
var ErrInvalidPhotoPath = errors.New("invalid photo path")
var ErrInvalidPhotoFile = errors.New("invalid file or unsupported photo Format")
var ErrInvalidIdentifierMode = errors.New("invalid identifier mode can only create U P C ids")
var ErrUserNotFound = errors.New("user not found")
