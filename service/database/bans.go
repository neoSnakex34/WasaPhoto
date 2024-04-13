package database

import (
	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (db *appdbimpl) BanUser(bannerId structs.Identifier, bannedId structs.Identifier) error {

	var counter int
	var err error

	// check if user is arleady banned
	err = db.c.QueryRow(`SELECT COUNT(*) FROM bans WHERE bannerId = ? AND bannedId = ?`, bannerId.Id, bannedId.Id).Scan(&counter)

	if err != nil {
		return err
	} else if counter > 0 {
		return customErrors.ErrAlreadyBanned
	} else if counter == 0 { // redundand check just to be paranoid
		// here i ban
		err = db.addBan(bannerId.Id, bannedId.Id)
		if err != nil {
			return err
		}

	}

	println("user successfully banned")
	return nil

}

func (db *appdbimpl) UnbanUser(bannerId structs.Identifier, bannedId structs.Identifier) error {

	var counter int
	var err error
	err = db.c.QueryRow(`COUNT(*) FROM bans WHERE bannerId = ? AND bannedId = ?`, bannerId.Id, bannedId.Id).Scan(&counter)

	if err != nil {
		return err
	} else if counter == 0 {
		return customErrors.ErrNotBanned
	} else if counter > 0 {
		err = db.removeBan(bannerId.Id, bannedId.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) addBan(bannerId string, bannedId string) error {

	_, err := db.c.Exec(`INSERT INTO bans (bannerId, bannedId) VALUES (?, ?)`, bannerId, bannedId)
	return err

}

func (db *appdbimpl) removeBan(bannerId string, bannedId string) error {

	_, err := db.c.Exec(`DELETE FROM bans WHERE bannerId = ? AND bannedId = ?`, bannerId, bannedId)
	return err
}

// TODO probably need to add a function to check if a user is banned
// when building core functionality decide it
