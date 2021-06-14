package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	basePath     string = "/photo_sync/users/"
	postBasePath string = "/objects/"
)

func CreateObjectFile(object []byte, objectName string, objectExtension string, username string) (string, int64, error) {
	objectPath := basePath + username + postBasePath + objectName + objectExtension
	exists, err := FileExists(objectPath)
	if err != nil {
		fmt.Print(err)
		return "", 0, err
	}
	if exists {
		return "", 0, errors.New("already_exists")
	}
	err = ioutil.WriteFile(objectPath, object, 0644)
	if err != nil {
		fmt.Print(err)
		return "", 0, err
	}
	info, _ := os.Stat(objectPath)
	return objectPath, info.Size(), nil
}

// Checks if the file already exists
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return true, nil

}
