package main

import (
	"api-service/internal/config"
	"api-service/pkg/utils"
	"fmt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	err = utils.RunMigrations(cfg)
	if err != nil {
		panic(fmt.Sprintf("Migration failed: %v", err))
	}

	fmt.Println("Migrations completed successfully!")
}
