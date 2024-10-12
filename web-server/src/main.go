package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type getUserInput struct {
	username string `json:"username"`
}

func main() {
	pgDb, err := newPostgresClient()
	if err != nil {
		log.Fatal(err)
	}

	_, err = newMongoClient()
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

func getUsers(w http.ResponseWriter, r *http.Request, pgDb *postgresClient) error {
	// get username from json post body, getUserInput type
	var input getUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot decode json input"))
		return fmt.Errorf("cannot decode json input %s: %w", input, err)
	}

	// get user from postgres
	user, err := pgDb.getUser(input.username)
	if err != nil {
		log.Fatal(err)
	}

	// return user as json
	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot marshal user to json"))
		return fmt.Errorf("cannot marshal user to json %s: %w", input.username, err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)

	return nil
}
