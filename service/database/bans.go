package database

import (
	"database/sql"
	"errors"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

func (db *appdbimpl) BanUser(bannerId structs.Identifier, bannedId structs.Identifier) error {

	var counter int
	var err error

	// check if user is arleady banned
	err = db.c.QueryRow(`SELECT COUNT(*) FROM bans WHERE bannerId = ? AND bannedId = ?`, bannerId.Id, bannedId.Id).Scan(&counter)
	if errors.Is(err, sql.ErrNoRows) {

		err = db.addBan(bannerId.Id, bannedId.Id)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		return errors.New("user is already banned")
	}

	println("user successfully banned")
	return nil

}

func (db *appdbimpl) UnbanUser(bannerId, bannedId structs.Identifier) error {

	_, err := db.c.Exec("DELETE FROM bans WHERE bannerId = ? AND bannedId = ?", bannerId, bannedId)
	if err != nil {
		return err
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
