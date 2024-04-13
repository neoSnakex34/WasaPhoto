package serviceutilities

import "regexp"

// [ ] IMPORTANT this should be implemented in login also !!!!
// FIXME
func CheckRegexNewUsername(username string) bool {

	println("regex call")

	usernameRegex := "^[a-z0-9]*?$"
	matched, err := regexp.MatchString(usernameRegex, username)
	if err != nil {
		return false
	}

	println("matched: ", matched)

	return matched
}
