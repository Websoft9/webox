package main

import (
	"fmt"
	"log"
	"net/http"
)

// User represents a user in the system
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserService handles user operations
type UserService struct {
	users []User
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		users: []User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		},
	}
}

// GetUser returns a user by ID
func (s *UserService) GetUser(id int) (*User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name string) *User {
	user := User{
		ID:   len(s.users) + 1,
		Name: name,
	}
	s.users = append(s.users, user)
	return &user
}

func main() {
	service := NewUserService()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Users: %+v", service.users)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
