package database

import (
	"errors"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (db *appdbimpl) FollowUser(followerId structs.Identifier, followedId structs.Identifier) error {
	var counter int

	// check if user is arleady followed by userId
	err := db.checkBan(followerId.Id, followedId.Id) // opposite of following
	// err := db.checkBan(followedId.Id, followerId.Id) // opposite of following
	if errors.Is(err, customErrors.ErrIsBanned) {
		println("user is banned cannot follow") // TODO log this in api
		return err
	} else if err != nil {
		return err
	}

	err = db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followerId = ? AND followedId = ?`, followerId.Id, followedId.Id).Scan(&counter)

	println("counter: ", counter)

	if err != nil {
		return err
	} else if counter > 0 {
		// FIXME this is obviously a bug
		// aaaaaand it is
		return customErrors.ErrAlreadyFollowing
	} else if counter == 0 {

		// then is followable add a func addfollow
		err = db.addFollow(followerId.Id, followedId.Id)

		if err != nil {
			// if this is not hit it will return nil at end of function than user is succesfully
			return err
		}
	}
	println("user successfully followed") // TODO log this in api
	return nil                            // check if this can be nil
}

func (db *appdbimpl) UnfollowUser(followerId structs.Identifier, followedId structs.Identifier) error {
	// check if user is actually followed by userId
	var counter int

	// check if user is banned
	err := db.checkBan(followerId.Id, followedId.Id)
	if errors.Is(err, customErrors.ErrIsBanned) {
		println("user is banned cannot unfollow") // TODO log this in api
		return err
	} else if err != nil {
		return err
	}

	err = db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followerId = ? AND followedId = ?`, followerId.Id, followedId.Id).Scan(&counter)
	if err != nil {
		return err
	} else if counter == 0 {
		return customErrors.ErrNotFollowing
	} else {
		err = db.removeFollow(followerId.Id, followedId.Id)
	}

	println("user successfully unfollowed") // TODO log this in api
	return err

}

func (db *appdbimpl) addFollow(followerId string, followedId string) error {

	_, err := db.c.Exec(`INSERT INTO followers (followerId, followedId) VALUES (?, ?)`, followerId, followedId)
	return err
}

func (db *appdbimpl) removeFollow(followerId string, followedId string) error {

	_, err := db.c.Exec(`DELETE FROM followers WHERE followerId = ? AND followedId = ?`, followerId, followedId)
	return err
}

// [ ] improve and test this
func (db *appdbimpl) getFollowerList(followedId string) ([]string, error) {

	var followerList []string
	var followerId string
	rows, err := db.c.Query(`SELECT followedId FROM followers WHERE followerId = ?`, followedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&followerId)
		if err != nil {
			return nil, err
		}
		followerList = append(followerList, followerId)
	}

	return followerList, nil

}
