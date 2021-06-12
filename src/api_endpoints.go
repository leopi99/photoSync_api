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
	authNotNeeded []string = []string{"/login", "/register"}
	apiKeys       map[string]string
)

//	Initialize the listeners on the endoponts
func InitializeApiEndPoints() {
	apiKeys = make(map[string]string)
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
	s.HandleFunc("/logout", handlerLogout)
	s.HandleFunc("/updateDownloadedObject", handlerUpdateDownloadedObjetc)
	s.HandleFunc("/updateProfile", handlerUpdateProfile)
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
		//Checks if the request needs the authentication
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
	contained := containsMap(apiKeys, apiKey)
	if !contained {
		w.Write(ErrorStruct{ErrorType: "Auth", Description: "This operation needs authentication"}.toJSON())

	}
	return contained
}

func writeGenericError(w http.ResponseWriter, r *http.Request, errorStruct ErrorStruct) {
	//Sets the errors if none is provided
	if errorStruct.Description == "" {
		errorStruct.Description = "An error occured"
	}
	if errorStruct.ErrorType == "" {
		errorStruct.ErrorType = "Internal Server Error"
	}

	if errorStruct.errorStatusCode == 999 {
		errorStruct.errorStatusCode = 500
	}
	w.WriteHeader(errorStruct.errorStatusCode)
	w.Write(errorStruct.toJSON())
}

func containsMap(thisMap map[string]string, word string) bool {
	contained := false
	for _, value := range thisMap {
		if value == word {
			contained = true
			break
		}
	}
	return contained
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
		writeGenericError(w, r, ErrorStruct{ErrorType: "user_not_selected", errorStatusCode: 400, Description: "User identification not set"})
		return
	}
	objects, err := GetUserObjectsFiltered(userID, "picture")
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
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
		writeGenericError(w, r, ErrorStruct{ErrorType: "user_not_selected", errorStatusCode: 400, Description: "User identification not set"})
		return
	}
	objects, err := GetUserObjectsFiltered(userID, "video")
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
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
		writeGenericError(w, r, ErrorStruct{ErrorType: "user_not_selected", errorStatusCode: 400, Description: "User identification not set"})
		return
	}
	objects, err := GetUserObjects(userID)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
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
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		return
	}
	err = AddPicture(rawObject)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		return
	}
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "missing_parameter", errorStatusCode: 400, Description: "One or more parameters needed are missing"})
		return
	}
	user, err := databaseLogin(User{password: password, Username: username})
	if err != nil {
		if err.Error() == "user_not_found" {
			writeGenericError(w, r, ErrorStruct{ErrorType: "user_not_found", errorStatusCode: 999, Description: "User not found"})
		} else {
			writeGenericError(w, r, ErrorStruct{ErrorType: "wrong_credentials", errorStatusCode: 400, Description: "Wrong credentials"})
		}
		return
	}
	if user.Username == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "wrong_credentials", errorStatusCode: 400, Description: "Wrong credentials"})
		return
	}
	var key string
	if apiKeys[username] == "" {
		key = tokenGenerator()
		apiKeys[username] = key
	} else {
		key = apiKeys[username]
	}

	user.ApiKey = key
	w.Write(user.toJSON())
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "missing_parameter", errorStatusCode: 400, Description: "One or more parameters needed are missing"})
		return
	}

	var user User
	user.Username = username
	user.password = password
	err := databaseRegister(user)

	if err != nil {
		fmt.Println(err)
		writeGenericError(w, r, ErrorStruct{ErrorType: "registration_error", errorStatusCode: 500, Description: "Registration error"})
		return
	}

	w.Write([]byte("{\"result\": \"ok\"}"))
}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	delete(apiKeys, username)
	w.Write([]byte("{\"result\":\"ok\"}"))
}

func handlerUpdateDownloadedObjetc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	objectID := r.Form.Get("objectID")
	value := r.Form.Get("value")
	if value == "true" {
		value = "1"
	} else {
		value = "0"
	}
	err := databaseUpdateDownloadedObject(objectID, value)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		return
	}
	w.Write([]byte("{\"result\": \"ok\"}"))
}

func handlerUpdateProfile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.Form.Get("password")
	username := r.Form.Get("username")
	if password == "" {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 400, ErrorType: "missing_parameter", Description: "One or more parameters needed are missing"})
		return
	}

	var user User = User{Username: username, password: password}
	err := databaseUpdateProfile(user)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		return
	}
	w.Write([]byte("{\"result\": \"ok\"}"))
}
