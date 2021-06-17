package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const (
	basePath     string = "/photo_sync/users/"
	postBasePath string = "/objects/"
)

var (
	localPath, _ = homedir.Dir()
	filesPath    = localPath + basePath
)

func CreateObjectFile(object []byte, objectName string, objectExtension string, username string) (string, int64, error) {
	path := filesPath + username + postBasePath
	err := CreateDir(strings.ReplaceAll(path, "\\", "/")) //Creates the directory if doesn't exist
	if err != nil {
		fmt.Println(err)
	}
	objectPath := path + objectName + objectExtension
	objectPath = strings.ReplaceAll(objectPath, "\\", "/")
	exists, _ := FileExists(objectPath) //Checks if the file already exists
	if exists {
		return "", 0, errors.New("file_already_exists")
	}
	fmt.Println("Saving the file at: " + objectPath)
	err = os.WriteFile(objectPath, object, 0666)
	if err != nil {
		fmt.Print(err)
		return "", 0, err
	}
	info, _ := os.Stat(objectPath)
	return objectPath, info.Size(), nil
}

func CreateDir(directory string) error {
	err := os.MkdirAll(directory, 0755)
	return err
}

// Checks if the file already exists
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return true, nil

}
