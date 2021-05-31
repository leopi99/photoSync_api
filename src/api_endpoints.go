package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverBaseEndpoint string = "photoSync/api/v1"
	deployPort         string = ":8080"
)

var (
	authNotNeeded []string = []string{"/login"}
	authApi       string   = "1234"
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
	s.HandleFunc("/addPicture", handlerAddPicture).Methods("POST")
	s.HandleFunc("/login", handlerLogin).Methods("POST")

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
			if checkApiKey(w, r) {
				next.ServeHTTP(w, r)
			} else {
				r.Body.Close()
			}
		}
	})
}

func checkApiKey(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.URL.Query().Get("apiKey")
	if apiKey != authApi {
		w.Write(ErrorStruct{ErrorType: "Auth", Description: "The auth key provided is not correct"}.toJSON())
	}
	return apiKey == authApi
}

func writeGenericError(w http.ResponseWriter, r *http.Request) {
	w.Write(ErrorStruct{ErrorType: "Internal Server Error", Description: "An error occured"}.toJSON())
}

//Generator for the API_KEY
func tokenGenerator() string {
	b := make([]byte, 32)
	rand.Read(b)
	fmt.Println("New api auth generated")
	return fmt.Sprintf("%x", b)
}

///
///	Endpoints handlers
///

func handlerGetPictures(w http.ResponseWriter, r *http.Request) {
	objects, err := GetUserObjectsFiltered("", "picture")
	if err != nil {
		writeGenericError(w, r)
		return
	} else {
		w.Write(objects.toJSON())
	}
}

func handlerGetVideos(w http.ResponseWriter, r *http.Request) {
	objects, err := GetUserObjectsFiltered("", "video")
	if err != nil {
		writeGenericError(w, r)
		return
	} else {
		w.Write(objects.toJSON())
	}
}

func handlerGetObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := GetUserObjects("")
	if err != nil {
		writeGenericError(w, r)
		return
	} else {
		w.Write(objects.toJSON())
	}
}

func handlerAddPicture(w http.ResponseWriter, r *http.Request) {
	// TODO: Add the Object generation from the POST request
	var rawObject RawObject
	var err error
	data := []byte(r.PostForm.Get("data"))
	rawObject, err = ObjectfromJSON(data)
	if err != nil {
		writeGenericError(w, r)
		return
	}
	err = AddPicture(rawObject)
	if err != nil {
		writeGenericError(w, r)
		return
	}
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := databaseLogin(User{Password: password, Username: username})
	if err != nil {
		writeGenericError(w, r)
		return
	}
	authApi = tokenGenerator()
	user.ApiKey = authApi
	w.Write(user.toJSON())
}
