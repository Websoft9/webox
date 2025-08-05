package main

import (
	"testing"
)

func TestNewUserService(t *testing.T) {
	service := NewUserService()
	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if len(service.users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(service.users))
	}
}

func TestGetUser(t *testing.T) {
	service := NewUserService()

	// Test existing user
	user, err := service.GetUser(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("Expected Alice, got %s", user.Name)
	}

	// Test non-existing user
	_, err = service.GetUser(999)
	if err == nil {
		t.Error("Expected error for non-existing user")
	}
}

func TestCreateUser(t *testing.T) {
	service := NewUserService()

	user := service.CreateUser("Charlie")
	if user.Name != "Charlie" {
		t.Errorf("Expected Charlie, got %s", user.Name)
	}

	if len(service.users) != 3 {
		t.Errorf("Expected 3 users after creation, got %d", len(service.users))
	}
}
