// web-server/src/postgres_test.go
package main

import (
	"errors"
	"os"
	"testing"
)

// MockUserGetter is a mock implementation of the UserGetter interface for testing.
type MockUserGetter struct {
	GetUserFunc func(username string) (pgRowUser, error)
}

// Implement the UserGetter interface for the mock.
func (m *MockUserGetter) GetUser(username string) (pgRowUser, error) {
	return m.GetUserFunc(username)
}

func TestGetUserSuccess(t *testing.T) {
	// Arrange
	expectedUser := pgRowUser{Id: 1, Name: "testuser", Email: "test@example.com", password: "password"}
	mockGetter := &MockUserGetter{
		GetUserFunc: func(username string) (pgRowUser, error) {
			if username == "testuser" {
				return expectedUser, nil
			}
			return pgRowUser{}, errors.New("user not found")
		},
	}

	// Act
	user, err := mockGetter.GetUser("testuser")

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user.Id != expectedUser.Id || user.Name != expectedUser.Name || user.Email != expectedUser.Email {
		t.Errorf("Expected user, got: %+v", user)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// Arrange
	mockGetter := &MockUserGetter{
		GetUserFunc: func(username string) (pgRowUser, error) {
			return pgRowUser{}, errors.New("user not found")
		},
	}

	// Act
	_, err := mockGetter.GetUser("nonexistentuser")

	// Assert
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got: %v", err)
	}
}

func TestNewPostgresClientMissingEnvVars(t *testing.T) {
	// Arrange
	// Temporarily unset environment variables
	originalHost := os.Getenv("POSTGRES_HOST")
	os.Unsetenv("POSTGRES_HOST")
	defer os.Setenv("POSTGRES_HOST", originalHost) // Restore the environment variable

	// Act
	_, err := newPostgresClient()

	// Assert
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Error() != "POSTGRES_HOST environment variable not set" {
		t.Errorf("Expected 'POSTGRES_HOST environment variable not set' error, got: %v", err)
	}
}

func TestNewPostgresClientSuccess(t *testing.T) {
	// Arrange
	// Set environment variables
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_DB", "testdb")
	os.Setenv("POSTGRES_PASSWORD", "testpassword")

	// Act
	client, err := newPostgresClient()

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Errorf("Expected client, got nil")
	}
}