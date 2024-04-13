package customerrors

import "errors"

// FIXME check errors new in other files and use them here

var ErrInvalidRegexUsername = errors.New("invalid username")

var ErrAlreadyFollowing = errors.New("already following")
var ErrNotFollowing = errors.New("not following")
var ErrIsBanned = errors.New("cannot follow, user is banned")

var ErrAlreadyBanned = errors.New("already banned")
var ErrNotBanned = errors.New("not banned")
