package database

import (
	"database/sql"
	"errors"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (db *appdbimpl) FollowUser(userId structs.Identifier, followerdId structs.Identifier) error {
	var counter int
	var err error
	// check if user is arleady followed by userId
	err = db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followerId = ? AND followedId = ?`, userId.Id, followerdId.Id).Scan(&counter)

	if errors.Is(err, sql.ErrNoRows) {
		// then is followable add a func addfollow
		err = db.addFollow(userId.Id, followerdId.Id)

		if err != nil {
			// if this is not hit it will return nil at end of function than user is succesfully
			return err
		}
	} else if err != nil {
		return err
	} else {
		return errors.New("user is already followed")
	}

	return nil
}

func (db *appdbimpl) UnfollowUser(userId structs.Identifier, follwedId structs.Identifier) error {
	// check if user is actually followed by userId
	var counter int
	var err error

	err = db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followerId = ? AND followedId = ?`, userId.Id, follwedId.Id).Scan(&counter)
	if err != nil {
		return err
	} else {
		err = db.removeFollow(userId.Id, follwedId.Id)
	}

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
