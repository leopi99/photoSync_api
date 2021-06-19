package main

import (
	"encoding/json"
)

///
/// 	Structs declaration
///

type ErrorStruct struct {
	ErrorType       string `json:"error"`
	Description     string `json:"description"`
	errorStatusCode int
}

type MediaAttributes struct {
	SyncDate        string `json:"sync_date"`
	CreationDate    string `json:"creation_date"`
	PicturePosition string `json:"position"`
	BytesSize       int64  `json:"bytes_size"`
	LocalPath       string `json:"local_path"`
	DatabaseID      int    `json:"database_id"`
	Downloaded      bool   `json:"downloaded"`
	Extension       string `json:"extension"`
	LocalID         int    `json:"local_id"`
}

type Object struct {
	Attributes MediaAttributes `json:"attributes"`
	Type       string          `json:"type"`
}

type RawObject struct {
	ObjectStruct Object `json:"object"`
	FileBytes    []byte `json:"fileBytes"`
}

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	password string
	ApiKey   string `json:"apiKey"`
}

type Objects []Object

///
///	Marshal functions => toJSON
///

// Returns the bytes of the json
// TODO: uncomment when needed
// func (object Object) toJSON() []byte {
// 	json, error := json.Marshal(object)
// 	if error != nil {
// 		panic(error)
// 	}
// 	return json
// }

// Returns the bytes of the json
func (err ErrorStruct) toJSON() []byte {
	json, error := json.Marshal(err)
	if error != nil {
		panic(error)
	}
	return json
}

func (objects Objects) toJSON() []byte {
	json, error := json.Marshal(objects)
	if error != nil {
		panic(error)
	}
	return json
}

func (user User) toJSON() []byte {
	json, error := json.Marshal(user)
	if error != nil {
		panic(error)
	}
	return json
}

//
// 	UnMarshal function (object) => fromJSON
//

func RawObjectfromJSON(jsonBytes []byte) (RawObject, error) {
	var obj RawObject
	err := json.Unmarshal(jsonBytes, &obj)
	return obj, err
}
