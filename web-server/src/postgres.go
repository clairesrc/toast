// web-server/src/postgres.go
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type pgRowUser struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	password string
}

// UserGetter is an interface for getting user data.  This allows us to mock
// the database interaction during testing.
type UserGetter interface {
	GetUser(username string) (pgRowUser, error)
}

type postgresClient struct {
	*sql.DB
}

// Implement the UserGetter interface
func (p *postgresClient) GetUser(username string) (pgRowUser, error) {
	var user pgRowUser
	err := p.QueryRow("SELECT id, name, email, password FROM user WHERE name = $1", username).Scan(&user.Id, &user.Name, &user.Email, &user.password)
	if err != nil {
		return user, fmt.Errorf("cannot get user from database: %w", err)
	}
	return user, nil
}

// newPostgresClient creates a new postgresClient
func newPostgresClient() (UserGetter, error) { // Return the interface, not the concrete type
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