# PhotoSync
## What is PhotoSync_api?

This is the backend for the [PhotoSync app](https://github.com/leopi99/photoSync_app).

## Not stable for "production" this is still in alpha, everything can change

## Features

- Handles login and registration
- Adds pictures to the database; this does not saves the image in the database, using the structured paths uses the username and the file creation date.
- Saves the media files in the server
- Serve the syncronized files (currently pictures and videos)
- Available outside your house [TODO] (maybe using ngrok)

## How to setup the server

- To create the database you can import the .sql file from this repository (photosync.sql, probably not finished)
- Create a file inside the src folder (named as you want) with 4 constants (databaseUsername, databasePassword, databaseAddress, databaseName), note that the databaseAddress must contain the port.
- You can change the api port and the base endpoint url from the api_endpoints.go
- To run the api you can use (inside the src folder) run `go run .` or you can build the project using `go build .` and then running the executable created.

## Why PhotoSync?
Since the popular Google photo have dropped the unlimited storage, I thougth to create this to have more control on where my photos are stored and I can upgrade the storage size when I want 
