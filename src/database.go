package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	databaseUsername string = ""
	databasePassword string = ""
	databaseAddress  string = ""
	databaseName     string = ""
)

var (
	database *sql.DB
)

//	Initialize the database connection
func InitializeDatabaseConnection() error {
	fmt.Println("Initializating database connection")
	var err error
	database, err = sql.Open("mysql", databaseUsername+":"+databasePassword+"@tcp("+databaseAddress+")/"+databaseName)
	return err
}

//	Returns all the pictures and the videos saved for a user
func GetUserObjects(username string) (Objects, error) {
	rows, err := database.Query("SELECT * FROM picture WHERE username = " + username + ";")
	if err != nil {
		return nil, err
	}
	var objects Objects
	for rows.Next() {
		var currentPicture Object
		rows.Scan(&currentPicture.Attributes.databaseID, &currentPicture.Attributes.CreationDate, &currentPicture.Attributes.PicturePosition, &currentPicture.Attributes.SyncDate, &currentPicture.Attributes.Url, &currentPicture.Attributes.UserProperty)
		objects = append(objects, currentPicture)
	}
	return objects, nil
}

//	Returns all the pictures or videos of a user
func GetUserObjectsFiltered(username string, objType string) (Objects, error) {
	rows, err := database.Query("SELECT * FROM picture WHERE username = " + username + "AND type = \"" + objType + "\"" + ";")
	if err != nil {
		return nil, err
	}
	var objects Objects
	for rows.Next() {
		var currentPicture Object
		rows.Scan(&currentPicture.Attributes.databaseID, &currentPicture.Attributes.CreationDate, &currentPicture.Attributes.PicturePosition, &currentPicture.Attributes.SyncDate, &currentPicture.Attributes.Url, &currentPicture.Attributes.UserProperty)
		objects = append(objects, currentPicture)
	}
	return objects, nil
}

//	Adds a picture into the db
func AddPicture(picture RawObject) error {
	imagePath, size, err := CreatePicture(picture.FileBytes, picture.ObjectStruct.Attributes.CreationDate)
	if err != nil {
		return err
	}
	picture.ObjectStruct.Attributes.BytesSize = size
	picture.ObjectStruct.Attributes.Url = imagePath
	fmt.Println("imagePath: " + imagePath)
	fmt.Printf("imageSize: %d", size)
	//TODO: implement the image save into the db
	return nil
}
