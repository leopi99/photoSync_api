package main

import "fmt"

const (
	apiVersion string = "0.8.0"
)

func main() {
	err := InitializeDatabaseConnection()
	if err != nil {
		fmt.Println("Error initializating the database connection")
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("Photo sync initialized v%s\n", apiVersion)
	InitializeApiEndPoints()
}
