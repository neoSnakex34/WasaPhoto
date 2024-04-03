/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// AppDatabase is the high level interface for the DB
// methods are exported ones, hence they are written with capital first letter
// TODO modify the methods to match the requirements
type AppDatabase interface {
	DoLogin(username string) (structs.Identifier, error)

	GetUserProfile(userId structs.Identifier) (structs.UserProfile, error)
	SetMyUserName(username string) error
	GetMyStream(userId structs.Identifier) ([]structs.Photo, error)

	FollowUser(userId structs.Identifier) error
	UnfollowUser(userId structs.Identifier) error

	BanUser(userId structs.Identifier) error
	UnbanUser(userId structs.Identifier) error

	UploadPhoto(file structs.PhotoFile) (structs.Photo, error)
	RemovePhoto(photoId structs.Identifier) error
	CommentPhoto(photoId structs.Identifier, commentBody structs.Comment) error
	UncommentPhoto(photoId structs.Identifier, commentId structs.Identifier) error
	LikePhoto(photoId structs.Identifier) error
	UnlikePhoto(photoId structs.Identifier) error

	//TODO consider adding methods to exract info from structs

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// start creating the AppDatabase if needed

	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// change integer with text or something like that to match regex pattern 
		userTable := `CREATE TABLE users (
			UserId VARCHAR(11) NOT NULL PRIMARY KEY,
			Username VARCHAR(18) NOT NULL UNIQUE//TODO remember to check elsewhere for the username to be unique
		
			)`
				//TODO is it all?
				// add photo counte and followercounter probably 
		
		//followerid will be a userid  
		//TODO check if that's correct
		followerTable := `CREATE TABLE followers (
			followerId VARCHAR(11) NOT NULL PRIMARY KEY,
			followedId VARCHAR(11) NOT NULL,
			FOREIGN KEY followerId REFERENCES userTable(userId),
		)`
		
		bansTable := `CREATE TABLE bans (
			bannerId VARCHAR(11) NOT NULL PRIMARY KEY,
			bannedId VARCHAR(11) NOT NULL,
			FOREIGN KEY bannerId REFERENCES userTable(userId),	
		)`

		photoTable := `CREATE TABLE photos (
			photoId VARCHAR(11) NOT NULL PRIMARY KEY, 
			userId VARCHAR(11) NOT NULL, 
			photo BLOB, 
			date TEXT, 
			FOREIGN KEY userId REFERENCES userTable(userId)
		)`
		
		likeTable := `CREATE TABLE likes (
			likeId VARCHAR(11) NOT NULL PRIMARY KEY,
			photoId VARCHAR(11) NOT NULL,
			FORIEGN KEY likeId REFERENCES userTable(userId)
			FOREIGN KEY photoId REFERENCES photoTable(photoId)
		)`

		commentTable := `CREATE TABLE comments (
			commentId VARCHAR(11) NOT NULL PRIMARY KEY,
			userId VARCHAR(11) NOT NULL,
			photoId VARCHAR(11) NOT NULL,
			body TEXT,
			FOREIGN KEY userId REFERENCES userTable(userId),
			FOREIGN KEY photoId REFERENCES photoTable(photoId)
		)`

		//TODO this would be executed one by one with dedicated errors probably
		//TODO check if i need to check for errors even here (function returns error if something goes wrong)
		runCreateQueries(db, userTable, followerTable, bansTable, photoTable, likeTable, commentTable)
		
	}

	// // Check if table exists. If not, the database is empty, and we need to create the structure
	// var tableName string
	// err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	// if errors.Is(err, sql.ErrNoRows) {
	// 	sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
	// 	_, err = db.Exec(sqlStmt)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error creating database structure: %w", err)
	// 	}
	// }

	return &appdbimpl{
		c: db,
	}, nil
}

func runCreateQueries(db *sql.DB, queries ...string) error {
	for _, query := range queries {
		_, err := db.Exec(query)

		if err != nil {
			return fmt.Errorf("error creating database table: %w", err)
		}
	}
	return nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

