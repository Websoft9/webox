package main

import (
	"api-service/internal/config"
	"api-service/pkg/database"
	"fmt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	_, err = database.InitDB(cfg)
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	fmt.Println("Database connection successful!")
}
