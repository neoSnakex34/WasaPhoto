package database

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	customErrors "github.com/neoSnakex34/WasaPhoto/service/custom-errors"
	"github.com/neoSnakex34/WasaPhoto/service/structs"
)

// TODO change this to be in home of user
const Folder string = "/tmp/wasaphoto/photofiles/"

// generate the identifier for the photo
// save the photofile path in the database
// save the photo in the database and create a new photo struct
// TODO decide when to use photostruct and comment struct in interactions
// FIXME will fronted give backend uploadphoto the file as a byte stream?
func (db *appdbimpl) UploadPhoto(file []byte, upoloaderUserId structs.Identifier, format string) (structs.Photo, error) {
	println(os.Getwd())
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

	photoPath = Folder + uploaderId + "/" + newPhotoId.Id + "." + format

	// FIRST save the photo file in the filesystem
	err = savePhotoFile(file, photoPath)
	if err != nil {
		return structs.Photo{}, err
	}

	date := time.Now().UTC().Format(time.RFC3339)
	// SECONDLY create the photo struct
	newPhoto := structs.Photo{
		PhotoId: newPhotoId,
		UserId:  upoloaderUserId,
		// Like:      0,                   // defaults not saved in the database
		// Comments:  []structs.Comment{}, // defaults not saved in the database
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

	approximatePhotoPath := Folder + removerUserId + "/" + removedPhotoId + ".*" // removes agnostically the image (without format) since ids are unique
	matchingPhoto, err := filepath.Glob(approximatePhotoPath)
	if err != nil {
		return err
	}

	if len(matchingPhoto) == 0 {
		return customErrors.ErrPhotoDoesNotExist
	}

	if len(matchingPhoto) > 1 {
		return customErrors.ErrCriticDuplicatedId
	}

	photoPath := matchingPhoto[0]

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

// path will be the final path
func savePhotoFile(file []byte, path string) error {
	// retrieve the directory
	dir := filepath.Dir(path)
	// build the directory if it does not exist (it doesn't first time cause there will be a directory for every user)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, file, 0644) // permission may fail on linux but since this is distributed via docker it should be fine
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

func (db *appdbimpl) getPhotosByUploaderId(plainUploaderId string) ([]structs.Photo, error) {
	var photos []structs.Photo
	var userId string = plainUploaderId
	var photoId string
	var date string
	var photoPath string

	// query to retrieve info
	rows, err := db.c.Query(`SELECT photoId, userId, date, photoPath FROM photos WHERE userId = ?`, plainUploaderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// user id is the same for all photos at every call, but for now is re assigned
		err = rows.Scan(&photoId, &userId, &date, &photoPath)
		if err != nil {
			return nil, err
		}
		photo := structs.Photo{
			PhotoId:   structs.Identifier{Id: photoId},
			UserId:    structs.Identifier{Id: userId},
			Date:      date,
			PhotoPath: photoPath,
		}
		photos = append(photos, photo)
	}

	return photos, err
}

// TODO complete me
func (db *appdbimpl) getSortedStreamOfPhotos(followerIdsForUser []string) ([]structs.Photo, error) { // TODO Note it returns a stream of photos, that needs to be displayed by obtaining info from the structs
	// for each follower i should retrieve a (photo slice) in order to build the stream
	// since i will need to sort the stream by date, i should return a complex struct instead of []string
	// and the access datas

	var stream []structs.Photo
	var tmpPhotos [][]structs.Photo
	for _, followerId := range followerIdsForUser {
		photos, err := db.getPhotosByUploaderId(followerId)
		if err != nil {
			return nil, err
		}
		tmpPhotos = append(tmpPhotos, photos)
	}

	// now stream will be an unsorted plain list of photos (no list of lists)
	for _, tmpList := range tmpPhotos {
		stream = append(stream, tmpList...)
	}

	// sort stream by date, i need to parse date with Time type
	// TODO
	sort.SliceStable(stream, func(i int, j int) bool {
		// FIXME should probably handle errors here
		// decide in debugging whether or not enabling error
		// handling here
		date1, _ := time.Parse(time.RFC3339, stream[i].Date)
		date2, _ := time.Parse(time.RFC3339, stream[j].Date)

		return date1.Before(date2)

	})
	// CHECK if stream is sorted and err is actually nil then return err and not nil
	return stream, nil

}
