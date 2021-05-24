package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverBaseEndpoint string = ""
	deployPort         string = ":8080"
)

var (
	authNotNeeded []string = []string{"/login"}
)

//	Initialize the listeners on the endoponts
func InitializeApiEndPoints() {
	fmt.Println("Initializating api endpoints")
	r := mux.NewRouter()
	s := r.PathPrefix(serverBaseEndpoint).Subrouter()
	s.Use(apiMiddleware)
	s.HandleFunc("/getPictures", handlerGetPictures).Methods("GET")
	s.HandleFunc("/getVideos", handlerGetVideos).Methods("GET")
	s.HandleFunc("/getAll", handlerGetObjects).Methods("GET")

	log.Fatal(http.ListenAndServe(deployPort, r))
}

///
///	Not categorized functions
///

//	Middleware for the apis
func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len(serverBaseEndpoint):]
		fmt.Printf("Handling %s for %s", path, r.RemoteAddr)
		w.Header().Add("Content-Type", "application/json")
		var found bool
		for _, checkedPath := range authNotNeeded {
			if checkedPath == path {
				found = true
				break
			}
		}
		if !found {
			checkApiKey(w, r)
		}
	})
}

func checkApiKey(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.URL.Query().Get("apiKey")
	return apiKey == "1234"
}

///
///	Endpoints handlers
///

func handlerGetPictures(w http.ResponseWriter, r *http.Request) {
	objects, err := GetUserObjectsFiltered("", "picture")
	if err != nil {
		w.Write(ErrorStruct{ErrorType: "Internal Server Error", Description: "An error occured"}.toJSON())
	} else {
		w.Write(objects.toJSON())
	}
}

func handlerGetVideos(w http.ResponseWriter, r *http.Request) {
	objects, err := GetUserObjectsFiltered("", "video")
	if err != nil {
		w.Write(ErrorStruct{ErrorType: "Internal Server Error", Description: "An error occured"}.toJSON())
	} else {
		w.Write(objects.toJSON())
	}
}

func handlerGetObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := GetUserObjects("")
	if err != nil {
		w.Write(ErrorStruct{ErrorType: "Internal Server Error", Description: "An error occured"}.toJSON())
	} else {
		w.Write(objects.toJSON())
	}
}
