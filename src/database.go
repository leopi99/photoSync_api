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
	rows, err := database.Query("SELECT objectID, local_path, creation_date, sync_date, picture_position, type, byte_size, downloaded, extension, local_id FROM object WHERE userID = " + userID + ";")
	if err != nil {
		return nil, err
	}
	var objects Objects
	for rows.Next() {
		var currentPicture Object
		rows.Scan(&currentPicture.Attributes.DatabaseID, &currentPicture.Attributes.LocalPath, &currentPicture.Attributes.CreationDate, &currentPicture.Attributes.SyncDate, &currentPicture.Attributes.PicturePosition, &currentPicture.Type, &currentPicture.Attributes.BytesSize, &currentPicture.Attributes.Downloaded, &currentPicture.Attributes.Extension, &currentPicture.Attributes.LocalID)
		objects = append(objects, currentPicture)
	}
	return objects, nil
}

//	Returns all the pictures or videos of a user
func GetUserObjectsFiltered(userID string, objType string) (Objects, error) {
	rows, err := database.Query("SELECT objectID, local_path, creation_date, sync_date, picture_position, type, byte_size, downloaded, extension, local_id FROM object WHERE userID = " + userID + " AND type = \"" + objType + "\";")
	if err != nil {
		return nil, err
	}
	var objects Objects
	for rows.Next() {
		var currentPicture Object
		rows.Scan(&currentPicture.Attributes.DatabaseID, &currentPicture.Attributes.LocalPath, &currentPicture.Attributes.CreationDate, &currentPicture.Attributes.SyncDate, &currentPicture.Attributes.PicturePosition, &currentPicture.Type, &currentPicture.Attributes.BytesSize, &currentPicture.Attributes.Downloaded, &currentPicture.Attributes.Extension, &currentPicture.Attributes.LocalID)
		objects = append(objects, currentPicture)
	}
	return objects, nil
}

//	Adds an object into the db
func AddObject(picture RawObject, username string, userID string) error {
	_, size, err := CreateObjectFile(picture.FileBytes, picture.ObjectStruct.Attributes.CreationDate, picture.ObjectStruct.Attributes.Extension, username)
	if err != nil {
		fmt.Println(err)
		return err
	}
	picture.ObjectStruct.Attributes.BytesSize = size
	//	Generates the query
	query := "INSERT INTO object (userID, local_path, creation_date, sync_date, picture_position, type, byte_size, downloaded, extension, local_id) VALUES ("
	query += userID + ",\"" + picture.ObjectStruct.Attributes.LocalPath + "\", \"" + picture.ObjectStruct.Attributes.CreationDate + "\","
	query += "\"" + picture.ObjectStruct.Attributes.SyncDate + "\",\"" + picture.ObjectStruct.Attributes.PicturePosition + "\", \"" + picture.ObjectStruct.Type + "\","
	query += fmt.Sprintf("%d", size) + ", true, \"" + picture.ObjectStruct.Attributes.Extension + "\","
	query += fmt.Sprintf("%d", picture.ObjectStruct.Attributes.LocalID) + ");"
	//	Executes the query
	_, err = database.Query(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func databaseLogin(user User) (User, error) {
	var userFound User
	rows, err := database.Query("SELECT username, password, userID FROM user WHERE username = \"" + user.Username + "\";")
	if err != nil {
		return userFound, err
	}
	for rows.Next() {
		rows.Scan(&userFound.Username, &userFound.password, &userFound.UserID)
	}
	if userFound.Username == "" {
		return userFound, errors.New("user_not_found")
	}
	userFound = User{}
	rows, err = database.Query("SELECT username, password, userID FROM user WHERE username = \"" + user.Username + "\" AND password = PASSWORD(\"" + user.password + "\");")
	if err != nil {
		return userFound, err
	}

	for rows.Next() {
		rows.Scan(&userFound.Username, &userFound.password, &userFound.UserID)
	}

	if userFound.password == "" {
		return User{}, errors.New("wrong_credentials")
	}
	return userFound, nil
}

func databaseRegister(user User) error {
	_, err := database.Query("INSERT INTO user(username, password) VALUES(\"" + user.Username + "\", PASSWORD(\"" + user.password + "\"));")
	return err
}

func databaseUpdateDownloadedObject(objectID string, value string) error {
	_, err := database.Query("UPDATE object SET downloaded = \"" + value + "\" WHERE objectID = " + objectID + ";")
	if err != nil {
		return err
	}

	return nil
}

func databaseUpdateProfile(user User) error {
	_, err := database.Query("UPDATE user SET password = PASSWORD(\"" + user.password + "\") WHERE username = \"" + user.Username + "\";")
	return err
}
