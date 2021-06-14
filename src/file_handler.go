package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	basePath     string = "/users/"
	postBasePath string = "/objects/"
)

var (
	_, b, _, _ = runtime.Caller(0)
	localPath  = filepath.Dir(b)
)

func CreateObjectFile(object []byte, objectName string, objectExtension string, username string) (string, int64, error) {
	localPath = strings.ReplaceAll(localPath, "\\", "/")
	CreateDir(localPath + basePath + username + postBasePath) //Creates the directory if doesn't exist
	objectPath := localPath + basePath + username + postBasePath + objectName + objectExtension
	exists, err := FileExists(objectPath) //Checks if the file already exists
	if err != nil {
		fmt.Print(err)
		return "", 0, err
	}
	if exists {
		return "", 0, errors.New("file_already_exists")
	}
	err = os.WriteFile(objectPath, object, 0644)
	if err != nil {
		fmt.Println("Write file error")
		fmt.Print(err)
		return "", 0, err
	}
	info, _ := os.Stat(objectPath)
	return objectPath, info.Size(), nil
}

func CreateDir(directory string) error {
	err := os.Mkdir(directory, 0755)
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
