package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

///Current todo work:
/// - Create the user existance check during signup
/// - Create the file server

const (
	serverBaseEndpoint string = "/photoSync/api/v1"
	deployPort         string = ":8010"
)

var (
	//Contains the api endpoints that doesn't need an apiKey to have the access
	authNotNeeded []string = []string{"/login", "/register"}
	//Contains the apiKeys assigned to the users (username - apiKey)
	apiKeys map[string]User
)

//	Initialize the listeners on the endoponts
func InitializeApiEndPoints() {
	apiKeys = make(map[string]User)
	fmt.Println("Initializating api endpoints")
	//Creates the router
	r := mux.NewRouter()
	//Creates the subRouter from the baseEndpoint
	s := r.PathPrefix(serverBaseEndpoint).Subrouter()
	//Sets the middleware that checks the apiKey
	s.Use(apiMiddleware)
	//Sets the handlers for the endpoints
	s.HandleFunc("/getPictures", handlerGetPictures)
	s.HandleFunc("/getVideos", handlerGetVideos)
	s.HandleFunc("/getAll", handlerGetObjects)
	s.HandleFunc("/addObject", handlerAddObject)
	s.HandleFunc("/login", handlerLogin)
	s.HandleFunc("/register", handlerRegistration)
	s.HandleFunc("/logout", handlerLogout)
	s.HandleFunc("/updateProfile", handlerUpdateProfile)
	s.HandleFunc("/object/{fileName}", handlerServeObject)
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
		//Handles the request; checks the apikey if needed, then proceeds to handle the request
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

//Checks that the apiKey is correct
func checkApiKey(w http.ResponseWriter, r *http.Request) bool {
	apiKey := getApiKey(r)
	contained := containsMap(apiKeys, apiKey)
	if !contained {
		fmt.Println("Api key not correct")
		w.Write(ErrorStruct{ErrorType: "Auth", Description: "This operation needs authentication"}.toJSON())

	}
	return contained
}

//Returns the username from the apiKey
func getUsernameFromApiKey(apiKey string) string {
	var username string
	for currentKey, value := range apiKeys {
		if value.ApiKey == apiKey {
			username = currentKey
			break
		}
	}
	return username
}

//Returns the userID from the apiKey
func getUserIDFromApiKey(apiKey string) int {
	var userID int
	for _, value := range apiKeys {
		if value.ApiKey == apiKey {
			userID = value.UserID
			break
		}
	}
	return userID
}

//Writes and error into the response
//This does not close the request, needs to add a return after the call to the function
func writeGenericError(w http.ResponseWriter, r *http.Request, errorStruct ErrorStruct) {
	//Sets the errors if none is provided
	if errorStruct.Description == "" {
		errorStruct.Description = "Internal Server Error"
	}
	if errorStruct.ErrorType == "" {
		errorStruct.ErrorType = "internal_server_error"
	}

	if errorStruct.errorStatusCode == 999 {
		errorStruct.errorStatusCode = 500
	}
	w.WriteHeader(errorStruct.errorStatusCode)
	w.Write(errorStruct.toJSON())
}

//Returns true if the apiKey is in the map of users
func containsMap(thisMap map[string]User, apiKey string) bool {
	contained := false
	for _, value := range thisMap {
		if value.ApiKey == apiKey {
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

//Returns the apiKey from the request
func getApiKey(r *http.Request) string {
	apiKey := r.URL.Query().Get("apiKey")
	if apiKey == "" {
		r.ParseForm()
		apiKey = r.Form.Get("apiKey")
	}
	return apiKey
}

///
///	Endpoints handlers
///

func handlerGetPictures(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := strconv.Itoa(getUserIDFromApiKey(getApiKey(r)))
	//Checks if the userID is in the request
	if userID == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "user_error", errorStatusCode: 400, Description: "User identification error"})
		return
	}
	//Gets the objects for a user
	objects, err := GetUserObjectsFiltered(userID, "picture")
	//Handles the response
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
	userID := strconv.Itoa(getUserIDFromApiKey(getApiKey(r)))
	//Gets the objects for a user
	if userID == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "user_error", errorStatusCode: 400, Description: "User identification error"})
		return
	}
	objects, err := GetUserObjectsFiltered(userID, "video")
	//Handles the response
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
	userID := strconv.Itoa(getUserIDFromApiKey(getApiKey(r)))
	//Gets the objects for a user
	if userID == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "user_error", errorStatusCode: 400, Description: "User identification error"})
		return
	}
	objects, err := GetUserObjects(userID)
	//Handles the response
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

func handlerAddObject(w http.ResponseWriter, r *http.Request) {
	var rawObject RawObject
	var err error
	r.ParseForm()
	userID := strconv.Itoa(getUserIDFromApiKey(getApiKey(r)))
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		fmt.Print(err)
		return
	}
	rawObject, err = RawObjectfromJSON(bytes)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		fmt.Print(err)
		return
	}
	err = AddObject(rawObject, getUsernameFromApiKey(r.Form.Get("apiKey")), userID)
	if err != nil {
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		fmt.Print(err)
		return
	}

	w.Write([]byte("{\"result\":\"ok\"}"))

}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	//Checks that all the info needed are in the request
	if username == "" || password == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "missing_parameter", errorStatusCode: 400, Description: "One or more parameters needed are missing"})
		return
	}
	//Checks the user existence and that the password is correct
	user, err := databaseLogin(User{password: password, Username: username})
	if err != nil {
		if err.Error() == "user_not_found" {
			writeGenericError(w, r, ErrorStruct{ErrorType: "user_not_found", errorStatusCode: 999, Description: "User not found"})
		} else {
			writeGenericError(w, r, ErrorStruct{ErrorType: "wrong_credentials", errorStatusCode: 400, Description: "Wrong credentials"})
		}
		return
	}
	//If no user is returned => the query didn't produce any rows => no user with the password
	if user.Username == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "wrong_credentials", errorStatusCode: 400, Description: "Wrong credentials"})
		return
	}
	var key string
	//Generates the apiKey if is the first login since the boot of the api
	if apiKeys[username].ApiKey == "" {
		key = tokenGenerator()
		apiKeys[username] = User{Username: username, ApiKey: key, UserID: user.UserID}
	} else {
		key = apiKeys[username].ApiKey
	}

	user.ApiKey = key
	w.Write(user.toJSON())
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	//Checks that all the info needed are in the request
	if username == "" || password == "" {
		writeGenericError(w, r, ErrorStruct{ErrorType: "missing_parameter", errorStatusCode: 400, Description: "One or more parameters needed are missing"})
		return
	}

	var user User
	user.Username = username
	user.password = password
	//Creates the user into the db
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
	username := getUsernameFromApiKey(getApiKey(r))
	//Removes the apiKey (and username) from the apiKeys map
	delete(apiKeys, username)
	w.Write([]byte("{\"result\":\"ok\"}"))
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

func handlerServeObject(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromApiKey(getApiKey(r))
	fileName := mux.Vars(r)["fileName"]
	path := filesPath + username + postBasePath + fileName
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		writeGenericError(w, r, ErrorStruct{errorStatusCode: 999})
		return
	}
	str := base64.StdEncoding.EncodeToString(file)
	w.Write([]byte(str))
}
