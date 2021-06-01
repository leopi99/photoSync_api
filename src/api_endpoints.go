package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverBaseEndpoint string = "/photoSync/api/v1"
	deployPort         string = ":8010"
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
	s.HandleFunc("/getPictures", handlerGetPictures)
	s.HandleFunc("/getVideos", handlerGetVideos)
	s.HandleFunc("/getAll", handlerGetObjects)
	s.HandleFunc("/addPicture", handlerAddPicture)
	s.HandleFunc("/login", handlerLogin)
	s.HandleFunc("/register", handlerRegistration)
	fmt.Println("Running from localhost" + deployPort + serverBaseEndpoint)
	log.Fatal(http.ListenAndServe(deployPort, r))
}

///
///	Functions without category
///

//	Middleware for the apis
func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len(serverBaseEndpoint):]
		fmt.Printf("Handling %s for %s\n", path, r.RemoteAddr)
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
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func checkApiKey(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.URL.Query().Get("apiKey")
	if apiKey == "" {
		r.ParseForm()
		apiKey = r.Form.Get("apiKey")
	}
	if apiKey != authApi {
		w.Write(ErrorStruct{ErrorType: "Auth", Description: "The auth key provided is not correct"}.toJSON())
	}
	return apiKey == authApi
}

func writeGenericError(w http.ResponseWriter, r *http.Request, description string, errorType string) {
	if description == "" {
		description = "An error occured"
	}
	if errorType == "" {
		errorType = "Internal Server Error"
	}
	w.WriteHeader(500)
	w.Write(ErrorStruct{ErrorType: errorType, Description: description}.toJSON())
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
	r.ParseForm()
	userID := r.Form.Get("userID")
	if userID == "" {
		writeGenericError(w, r, "user_not_selected", "User identification not set")
		return
	}
	objects, err := GetUserObjectsFiltered(userID, "picture")
	if err != nil {
		writeGenericError(w, r, "", "")
		fmt.Print(err)
		return
	} else {
		if len(objects) == 0 {
			w.Write([]byte("{}"))
		} else {
			w.Write(objects.toJSON())
		}
	}
}

func handlerGetVideos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := r.Form.Get("userID")
	if userID == "" {
		writeGenericError(w, r, "user_not_selected", "User identification not set")
		return
	}
	objects, err := GetUserObjectsFiltered(userID, "video")
	if err != nil {
		writeGenericError(w, r, "", "")
		fmt.Print(err)
		return
	} else {
		if len(objects) == 0 {
			w.Write([]byte("{}"))
		} else {
			w.Write(objects.toJSON())
		}
	}
}

func handlerGetObjects(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := r.Form.Get("userID")
	if userID == "" {
		writeGenericError(w, r, "user_not_selected", "User identification not set")
		return
	}
	objects, err := GetUserObjects(userID)
	if err != nil {
		writeGenericError(w, r, "", "")
		fmt.Print(err)
		return
	} else {
		if len(objects) == 0 {
			w.Write([]byte("{}"))
		} else {
			w.Write(objects.toJSON())
		}
	}
}

func handlerAddPicture(w http.ResponseWriter, r *http.Request) {
	// TODO: Add the Object generation from the POST request
	var rawObject RawObject
	var err error
	data := []byte(r.PostForm.Get("data"))
	rawObject, err = ObjectfromJSON(data)
	if err != nil {
		writeGenericError(w, r, "", "")
		return
	}
	err = AddPicture(rawObject)
	if err != nil {
		writeGenericError(w, r, "", "")
		return
	}
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		writeGenericError(w, r, "parameter_missing", "One or more parameters needed are missing")
		return
	}
	user, err := databaseLogin(User{Password: password, Username: username})
	if err != nil {
		writeGenericError(w, r, "db_error", "Database error")
		return
	}
	authApi = tokenGenerator()
	if user.Username == "" {
		writeGenericError(w, r, "User not found", "user_not_found_error")
		return
	}
	user.ApiKey = authApi
	w.Write(user.toJSON())
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		writeGenericError(w, r, "parameter_missing", "One or more parameters needed are missing")
		return
	}
}
