package serviceutilities

import (
	"errors"
	"log"
	"regexp"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
)

const Folder string = "/tmp/wasaphoto/photofiles/"

// TODO move folder in service and fix dependencies
// [ ] IMPORTANT this should be implemented in login also !!!!
// FIXME
func CheckRegexNewUsername(username string) bool {

	log.Println("regex check for username entered")

	usernameRegex := "^[a-z0-9]{3,12}?$"
	matched, err := regexp.MatchString(usernameRegex, username)
	if err != nil {
		return false
	}

	log.Println("regex check for username match: ", matched)

	return matched
}

func CheckFileType(file []byte) (string, error) {
	println("len of file: ", len(file))
	if len(file) < 8 {
		// TODO handle this error
		return "", errors.New("file is too small to be a photo")
	}
	switch {
	case file[0] == 0xFF &&
		file[1] == 0xD8 &&
		file[2] == 0xFF:
		return "jpg", nil

	case file[0] == 0x89 &&
		file[1] == 'P' &&
		file[2] == 'N' &&
		file[3] == 'G' &&
		file[4] == '\r' &&
		file[5] == '\n' &&
		file[6] == 0x1a &&
		file[7] == '\n':
		return "png", nil

	}
	// HANDLE
	return "", errors.New("file is not a photo")
}

func GetPhotoPath(partialPath string) (string, error) {
	if partialPath == "" {
		return "", customErrors.ErrInvalidPhotoPath
	}
	return Folder + partialPath, nil
}
