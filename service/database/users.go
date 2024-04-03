package database

import (
	"database/sql"
	"errors"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
	"github.com/neoSnakex34/WasaPhoto/service/utilities"
)

func (db *appdbimpl) DoLogin(username structs.UserName) (structs.Identifier, error) {

	var userId structs.Identifier
	// first we check if user is in the database querying his row (given that username is unique)
	err := db.c.QueryRow("SELECT userId FROM users WHERE username = ?", username).Scan(&userId)

	if errors.Is(err, sql.ErrNoRows) {
		//then user does not exist on the system
		//we need to create a new user
		userId, err = utilities.GenerateIdentifier("U")
		if err != nil {
			return structs.Identifier{}, err
		}

		var count int
		//check if id is unique (valid)
		err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE userId = ?`, userId).Scan(&count)
		//TODO print this count anywhere to check if it is working
		if err != nil{
		 	if errors.Is(err, sql.ErrNoRows) {
			//the id was not found
			//we can insert the new user
			//TODO INSERT USER
		 	} else {
				return structs.Identifier{}, err
			}
		
		}else{
			//the id was found
			//we need to generate a new one
			//log this on console for statistic purposes on collisions 
			//TODO make this thing an extern function to loop until a valid userId is generated
		
		}


	} 



	}

}
