package database

// FIXME most likely everything inside a struct must be unpacked to be used successfully in queries, i have to handle it
// TODO the fixme just above was checked and now everything i use except return statement in dologin is a plain string (instead of a wrapper struct)
// doing such thing is useful in queries but since those structs are useful
import (
	"database/sql"
	"errors"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// TODO i dont think that using errors to mask others is a good idea in debugging
// implement those only if you did enough testing
// var LoginError = errors.New("an error occured during login")

func (db *appdbimpl) DoLogin(username string) (structs.Identifier, error) {

	var userId string
	idIsValid := false

	exist, userId, err := db.checkUserExists(username)

	// if any error is found i return it (TODO handle)
	if err != nil {
		// check if you need to throw the login error or not
		return structs.Identifier{}, err

	} else {
		// else if the user exist i have to login
		if exist == true {
			// login
			return structs.Identifier{Id: userId}, nil

		} else if exist == false {

			// loop until a valid user or error is found
			for (idIsValid == false) && (err == nil) {
				idIsValid, err = db.validId(userId, "U")
				// TODO warning with this assignation, it could break everything
				tmpId, _ := GenerateIdentifier("U") // here error can be ignored since we are automatically using a valid actor
				userId = tmpId.Id
			}

			if err != nil {
				return structs.Identifier{}, err
			}

			// here i actually create the user by setting is username in N mode
			// setting username for the first time is part of the action of generating the userId
			// that it has been verified in the for (while) loop on line 38
			db.SetMyUserName(username, userId, "N")
		}

	}
	return structs.Identifier{Id: userId}, nil

}

func (db *appdbimpl) SetMyUserName(newUsername string, userId string, mode string) error {

	// TODO all this cheks must be lowercase, also i need an efficient way to loop over errors till
	// a valid name pops up

	var count int
	var valid bool = false

	//  if user is new one MODE = N i need to do inser
	// if user is already signed MODE = U i need to update by id

	//  i check if newUsername is taken
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, newUsername).Scan(&count)

	if errors.Is(err, sql.ErrNoRows) && count == 0 {
		valid = true
	} else {
		if err != nil {
			return err
		}
	}

	if valid == true {

		switch mode {

		case "N":
			// this could probably be the worst design i could ever imagine
			// TODO evaluate if fixing it is worth the effort
			err := db.createUser(newUsername, userId)
			return err

		case "U":
			_, err := db.c.Exec(`UPDATE users SET username = ? WHERE userId = ?`, newUsername, userId)
			return err
		default:
			return errors.New("error in parsing mode or invalid mode for userame operation")

		}

	}

	return err
}

// TODO getmystream and getmyuserprofile

// ========== private functions from here
func (db *appdbimpl) createUser(username string, userId string) error {

	_, err := db.c.Exec("INSERT INTO users (username, userId) VALUES (?, ?)", username, userId)
	return err
}

// TODO generalize for any id MAYBE
func (db *appdbimpl) checkUserExists(username string) (bool, string, error) {
	var userInTable bool
	// var userId structs.Identifier = structs.Identifier{}
	var id string
	// first we check if user is in the database querying his row (given that username is unique)
	err := db.c.QueryRow(`SELECT userId FROM users WHERE username = ?`, username).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		userInTable = false
		err = nil // else it will fail control in next function, very important to be checked !

	} else if err != nil {
		return false, id, err
	} else {
		// so the user exist
		// todo this could be prone to bugs, if something goes south check it out
		userInTable = true
		// userId = structs.Identifier{Id: id}
	}

	return userInTable, id, nil
}
