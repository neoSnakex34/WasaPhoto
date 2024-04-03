package database

import (
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)
func (db *appdbimpl) BanUser(bannerId, bannedId structs.Identifier) error {
	_, err := db.Exec("INSERT INTO bans (bannerId, bannedId) VALUES (?, ?)", bannerId, bannedId)
	if err != nil {
		return err
	}

	return nil

}

func (db *appdbimpl) UnbanUser(bannerId, bannedId structs.Identifier) error {
	_, err := db.Exec("DELETE FROM bans WHERE bannerId = ? AND bannedId = ?", bannerId, bannedId)
	if err != nil {
		return err
	}

	return nil
}