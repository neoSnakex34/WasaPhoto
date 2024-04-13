package customerrors

import "errors"

// FIXME check errors new in other files and use them here

var ErrInvalidRegexUsername = errors.New("invalid username")
