package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

///
///	The constants:
///	databaseUsername	string = ""
/// databasePassword	string = ""
/// databaseAddress		string = ""
/// databaseName		string = ""
///
///	Can be created inside another file (like sensible_data.go) as consts
///

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
func GetUserObjects(userID string) (Objects, error) {
	rows, err := database.Query("SELECT * FROM object WHERE userID = " + userID + ";")
	if err != nil {
		return nil, err
	}
	var objects Objects
	for rows.Next() {
		var currentPicture Object
		rows.Scan(&currentPicture.Attributes.DatabaseID, &currentPicture.Attributes.CreationDate, &currentPicture.Attributes.PicturePosition, &currentPicture.Attributes.SyncDate, &currentPicture.Attributes.Url, &currentPicture.Attributes.UserProperty)
		objects = append(objects, currentPicture)
	}
	return objects, nil
}

//	Returns all the pictures or videos of a user
func GetUserObjectsFiltered(userID string, objType string) (Objects, error) {
	rows, err := database.Query("SELECT * FROM object WHERE userID = " + userID + " AND type = \"" + objType + "\";")
	if err != nil {
		return nil, err
	}
	var objects Objects
	for rows.Next() {
		var currentPicture Object
		rows.Scan(&currentPicture.Attributes.DatabaseID, &currentPicture.Attributes.CreationDate, &currentPicture.Attributes.PicturePosition, &currentPicture.Attributes.SyncDate, &currentPicture.Attributes.Url, &currentPicture.Attributes.UserProperty)
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

func databaseLogin(user User) (User, error) {
	var userFound User
	//Checks if the user exists
	rows, err := database.Query("SELECT username FROM user WHERE username = \"" + user.Username + "\";")
	if err != nil {
		return userFound, err
	}
	for rows.Next() {
		rows.Scan(&userFound.Username)
	}
	if userFound.Username == "" {
		return userFound, errors.New("no_user_found")
	}
	userFound = User{}
	//Actually gets the user
	rows, err = database.Query("SELECT username, userID FROM user WHERE username = \"" + user.Username + "\" AND password = \"" + user.Password + "\"" + ";")

	if err != nil {
		return userFound, err
	}
	for rows.Next() {
		rows.Scan(&userFound.Username, &userFound.UserID)
	}
	return userFound, nil
}
