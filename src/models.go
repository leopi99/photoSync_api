package main

import "encoding/json"

///
/// 	Structs declaration
///

type ErrorStruct struct {
	ErrorType   string `json:"error"`
	Description string `json:"description"`
}

type MediaAttributes struct {
	Url             string `json:"url"`
	SyncDate        string `json:"sync_date"`
	CreationDate    string `json:"creation_date"`
	UserProperty    string `json:"user"`
	PicturePosition string `json:"position"`
	BytesSize       int64  `json:"bytes_size"`
	LocalPath       string `json:"local_path"`
	databaseID      int    `json:"database_id"`
}

type Object struct {
	Attributes MediaAttributes `json:"attributes"`
	Type       string          `json:"type"`
}

type RawObject struct {
	ObjectStruct Object
	FileBytes    []byte
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

func ObjectfromJSON(jsonBytes []byte) (RawObject, error) {
	var obj RawObject
	err := json.Unmarshal(jsonBytes, &obj)
	return obj, err
}
