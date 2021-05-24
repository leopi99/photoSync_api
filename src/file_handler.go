package main

import (
	"io/ioutil"
	"os"
)

const (
	basePath         string = ""
	videoExtension   string = ".mp4"
	pictureExtension string = ".png"
)

//	Returns the path of the file, the size of the file and then the error (if any)
func CreatePicture(picture []byte, imageName string) (string, int64, error) {
	imagePath := basePath + imageName + pictureExtension
	err := ioutil.WriteFile(imagePath, picture, 0644)
	if err != nil {
		return "", 0, err
	}
	info, _ := os.Stat(imagePath)
	return imagePath, info.Size(), nil
}

//	Returns the path of the file, the size of the file and then the error (if any)
func CreateVideo(video []byte, videoName string) (string, int64, error) {
	videopath := basePath + videoName + pictureExtension
	err := ioutil.WriteFile(videopath, video, 0644)
	if err != nil {
		return "", 0, err
	}
	info, _ := os.Stat(videopath)
	return videopath, info.Size(), nil
}

// Checks if the file already exists
func FileExists(fileName string, fileExtension string) (bool, error) {
	_, err := os.Stat(basePath + fileName + fileExtension)
	if err != nil {
		return false, err
	}
	return true, nil

}
