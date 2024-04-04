package database

import (
	"database/sql"
	"errors"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
	"github.com/neoSnakex34/WasaPhoto/service/utilities"
)

//TODO i dont think that using errors to mask others is a good idea in debugging
//implement those only if you did enough testing
//var LoginError = errors.New("an error occured during login")

// TODO generalize for any id MAYBE
func (db *appdbimpl) checkUserExists(username structs.UserName) (bool, structs.Identifier, error) {
	var userInTable bool
	var userId structs.Identifier = structs.Identifier{}
	var id string
	// first we check if user is in the database querying his row (given that username is unique)
	err := db.c.QueryRow("SELECT userId FROM users WHERE username = ?", username).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		userInTable = false
		err = nil //else it will fail control in next function, very important to be checked !

	} else if err != nil {
		return false, userId, err
	} else {
		//so the user exist
		//TODO this could be prone to bugs, if something goes south check it out
		userInTable = true
		userId = structs.Identifier{Id: id}
	}

	return userInTable, userId, nil
}

// TODO generalize for any id
func (db *appdbimpl) validId(id structs.Identifier) (bool, error) {

	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE userId = ?`, id).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil

}

func (db *appdbimpl) DoLogin(username structs.UserName) (structs.Identifier, error) {

	var userId structs.Identifier
	idIsValid := false

	exist, userId, err := db.checkUserExists(username)

	//if any error is found i return it (TODO handle)
	if err != nil {
		//check if you need to throw the login error or not
		return structs.Identifier{}, err

	} else {
		//else if the user exist i have to login
		if exist == true {
			//login
			return userId, nil

		} else if exist == false {
			//loop until a valid user or error is found
			for (idIsValid == false) && (err == nil) {
				idIsValid, err = db.validId(userId)
				//TODO warning with this assegnation, it could break all
				tmpId, _ := utilities.GenerateIdentifier("U") //here error can be ignored since we are automatically using a valid actor
				userId = tmpId
			}

			if err != nil {
				return structs.Identifier{}, err
			}
		}

	}
	return userId, nil

}
