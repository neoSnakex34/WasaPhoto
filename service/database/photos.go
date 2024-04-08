package database

import (
	"os"
	"time"

	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

const Folder string = "photofiles/"

// generate the identifier for the photo
// save the photofile path in the database
// save the photo in the database and create a new photo struct
// TODO decide when to use photostruct and comment struct in interactions
// FIXME will fronted give backend uploadphoto the file as a byte stream?
func (db *appdbimpl) UploadPhoto(file []byte, upoloaderUserId structs.Identifier) (structs.Photo, error) {

	var isValidId bool = false
	var newPhotoId structs.Identifier
	var err error
	var photoPath string
	var uploaderId string = upoloaderUserId.Id
	// generate a new photo valid id
	for isValidId == false && err == nil {

		newPhotoId, err = GenerateIdentifier("P")
		isValidId, err = db.validId(newPhotoId.Id, "P")

	}

	if err != nil {
		return structs.Photo{}, err
	}

	photoPath = Folder + uploaderId + "/" + newPhotoId.Id + ".jpg"

	// FIRST save the photo file in the filesystem
	err = savePhotoFile(file, photoPath)
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
		PhotoPath: photoPath,
		// PhotoBytes: file,
	}

	// AFTER FIRST TWO STEPS insert photo in the database
	err = db.insertPhotoInTable(newPhotoId.Id, upoloaderUserId.Id, date, newPhoto.PhotoPath)
	if err != nil {
		return structs.Photo{}, err
	}

	return newPhoto, nil
}

// [ ] check you built the path correctly
// maybe add a little func to build it
func (db *appdbimpl) RemovePhoto(photoId structs.Identifier, userId structs.Identifier) error {
	removedPhotoId := photoId.Id
	removerUserId := userId.Id
	// TODO for now the control on user right in removing photo is done elsewhere
	// FIXME remember to handle it
	photoPath := Folder + removerUserId + "/" + removedPhotoId + ".jpg"
	var err error
	err = db.removePhotoFromTable(removedPhotoId)
	if err != nil {
		return err
	}
	err = deletePhotoFile(photoPath)
	if err != nil {
		return err
	}
	return nil
}

func savePhotoFile(file []byte, path string) error {
	err := os.WriteFile(path, file, 0644) // permission may fail on linux but since this is distributed via docker it should be fine
	return err
}

func deletePhotoFile(path string) error {
	err := os.Remove(path)
	return err
}

func (db *appdbimpl) removePhotoFromTable(photoId string) error {
	_, err := db.c.Exec(`DELETE FROM photos WHERE photoId = ?`, photoId)
	return err
}

func (db *appdbimpl) insertPhotoInTable(photoId string, userId string, date string, path string) error {
	_, err := db.c.Exec(`INSERT INTO photos (photoId, userId, photoPath, date) VALUES (?, ?, ?, ?)`, photoId, userId, path, date)
	return err
}

// TODO complete me
func (db *appdbimpl) getStreamPhotoListForUser(followerIdsForUser []string) ([]string, error) {
	// for each follower i should retrieve a (photopath, date) in order to build the stream
	// since i will need to sort the stream by date, i should return a complex struct instead of []string
	// and the access datas

	return nil, nil
}
