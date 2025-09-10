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

	err = database.RunMigrations(cfg)
	if err != nil {
		panic(fmt.Sprintf("Migration failed: %v", err))
	}

	fmt.Println("Migrations completed successfully!")
}
