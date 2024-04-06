package database

import (
	"time"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

const path string = "photofiles/"

// generate the identifier for the photo
// save the photofile path in the database
// save the photo in the database and create a new photo struct
// TODO decide when to use photostruct and comment struct in interactions
func (db *appdbimpl) UploadPhoto(file []byte, upoloaderUserId structs.Identifier) (structs.Photo, error) {

	var isValidId bool = false
	var newPhotoId structs.Identifier
	var err error

	// generate a new photo valid id
	for isValidId == false && err == nil {

		newPhotoId, err = GenerateIdentifier("P")
		isValidId, err = db.validId(newPhotoId.Id, "P")

	}

	if err != nil {
		return structs.Photo{}, err
	}

	// FIRST save the photo file in the filesystem
	err = db.savePhotoFile(file)
	if err != nil {
		return structs.Photo{}, err
	}

	date := time.Now().Format(time.RFC3339)
	// SECONDLY create the photo struct
	newPhoto := structs.Photo{
		PhotoId:   newPhotoId,
		UserId:    upoloaderUserId,
		Like:      0,                   // defaults not saved in the database
		Comments:  []structs.Comment{}, // defaults not saved in the database
		Date:      date,
		PhotoPath: path + newPhotoId.Id + ".jpg",
		// PhotoBytes: file,
	}

	// AFTER FIRST TWO STEPS insert photo in the database
	err = db.insertPhotoInTable(newPhotoId.Id, upoloaderUserId.Id, date, newPhoto.PhotoPath)
	if err != nil {
		return structs.Photo{}, err
	}

	return newPhoto, nil
}

func (db *appdbimpl) savePhotoFile(file []byte) error {

}

func (db *appdbimpl) insertPhotoInTable(photoId string, userId string, date string, path string) error {

	_, err := db.c.Exec(`INSERT INTO photos (photoId, userId, photoPath, date) VALUES (?, ?, ?, ?)`, photoId, userId, path, date)
	return err

}
