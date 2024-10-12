package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type pgRowUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	password string
}

// postgresClient is a wrapper around the sql.DB type
// that allows us to mock the sql.DB type in tests
// by embedding this type in another struct
// and overriding the methods we need to mock
type postgresClient struct {
	*sql.DB
}

func (p *postgresClient) getUser(username string) (pgRowUser, error) {
	var user pgRowUser
	err := p.QueryRow("SELECT id, name, email, password FROM user WHERE name = $1", username).Scan(&user.id, &user.name, &user.email, &user.password)
	if err != nil {
		return user, fmt.Errorf("cannot get user from database: %w", err)
	}
	return user, nil
}

// newPostgresClient creates a new postgresClient
func newPostgresClient() (*postgresClient, error) {
	if os.Getenv("POSTGRES_HOST") == "" {
		return nil, fmt.Errorf("POSTGRES_HOST environment variable not set")
	}
	if os.Getenv("POSTGRES_USER") == "" {
		return nil, fmt.Errorf("POSTGRES_USER environment variable not set")
	}
	if os.Getenv("POSTGRES_DB") == "" {
		return nil, fmt.Errorf("POSTGRES_DB environment variable not set")
	}
	if os.Getenv("POSTGRES_PASSWORD") == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD environment variable not set")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", postgresHost, postgresUser, postgresDB, postgresPassword))
	if err != nil {
		return nil, fmt.Errorf("cannot open database connection: %w", err)
	}
	return &postgresClient{db}, nil
}
