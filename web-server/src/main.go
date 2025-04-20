// web-server/src/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type getUserInput struct {
	Username string `json:"username"`
}

func main() {
	pgDb, err := newPostgresClient()
	if err != nil {
		log.Fatal(err)
	}

	_, err = newMongoClient() // Assuming newMongoClient is defined elsewhere.
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to databases")

	// set up web server
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		err := getUsers(w, r, pgDb)
		log.Default().Println(fmt.Errorf("cannot get users: %w", err))
	})

	http.ListenAndServe(":8080", nil)
}

func getUsers(w http.ResponseWriter, r *http.Request, pgDb interface{}) error { // Use interface{} type
	// get username from json post body, getUserInput type
	var input getUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot decode json input"))
		return fmt.Errorf("cannot decode json input %s: %w", input, err)
	}

	// Type assertion to access the GetUser method.
	postgresClient, ok := pgDb.(*postgresClient) // Convert pgDb to *postgresClient type
	if !ok {
		return fmt.Errorf("pgDb is not a *postgresClient")
	}

	// get user from postgres
	user, err := postgresClient.GetUser(input.Username)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("cannot get user from database: %w", err) //return error so that error is handled
	}

	// return user as json
	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot marshal user to json"))
		return fmt.Errorf("cannot marshal user to json %s: %w", input.Username, err) //return error so that error is handled
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)

	return nil
}
