package database

import (
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)
func (db *appdbimpl) BanUser(bannerId, bannedId structs.Identifier) error {
	_, err := db.Exec("INSERT INTO banTable (banner_id, banned_id) VALUES (?, ?)", bannerId, bannedId)
	if err != nil {
		return err
	}

	return nil

} 