package main

import (
	"api-service/internal/config"
	"api-service/internal/dto/request"
	repoInterface "api-service/internal/interface/repository"
	"api-service/internal/model"
	repoImpl "api-service/internal/repository"
	serviceImpl "api-service/internal/service"
	"api-service/pkg/database"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
)

// ExitCodes for consistent error handling
const (
	ExitSuccess = 0
	ExitError   = 1
)

func main() {
	// Define command line flags
	var (
		username = flag.String("username", "", "Username or email for the new user (required)")
		password = flag.String("password", "", "Password for the new user (required)")
		roleCode = flag.String("role-code", "", "Role code for the new user (required)")
		help     = flag.Bool("help", false, "Show help information")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Create a new user in the Websoft9 system.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -username=admin@example.com -password=password123 -role-code=admin\n", os.Args[0])
	}

	flag.Parse()

	// Show help if requested
	if *help {
		flag.Usage()
		os.Exit(ExitSuccess)
	}

	// Validate required parameters
	if *username == "" {
		fmt.Fprintln(os.Stderr, "Error: username is required")
		flag.Usage()
		os.Exit(ExitError)
	}

	if *password == "" {
		fmt.Fprintln(os.Stderr, "Error: password is required")
		flag.Usage()
		os.Exit(ExitError)
	}

	if *roleCode == "" {
		fmt.Fprintln(os.Stderr, "Error: role-code is required")
		flag.Usage()
		os.Exit(ExitError)
	}

	// Create user
	if err := createUser(*username, *password, *roleCode); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
		os.Exit(ExitError)
	}

	fmt.Printf("User '%s' created successfully with role code '%s'\n", *username, *roleCode)
}

// createUser creates a new user using the UserService
func createUser(username, password, roleCode string) error {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize logging system
	zapLogger := logger.NewZapLoggerWithServerConfig(
		cfg.Server.Log.LogPath,
		cfg.Server.Log.LogLevel,
		cfg.Server.Log.LogMaxSize,
		cfg.Server.Log.LogMaxBackups,
		cfg.Server.Log.LogMaxAge,
		cfg.Server.Log.LogCompress,
	)
	logger.SetDefault(zapLogger)

	// 3. Initialize database connection
	dbWrapper, err := database.InitDBWrapper(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	db := dbWrapper.GetDB()

	// 5. Initialize repositories and services
	userRepo := repoImpl.NewUserRepository(db)
	roleRepo := repoImpl.NewRoleRepository(db)
	userService := serviceImpl.NewUserService(userRepo, zapLogger)

	// 6. Validate that the role exists and is active
	role, err := validateRoleByCode(roleRepo, roleCode)
	if err != nil {
		return err
	}

	// 7. Determine email and username
	email, actualUsername := parseUsernameEmail(username)

	// 8. Create user request
	req := &request.UserCreateRequest{
		Username: actualUsername,
		Email:    email,
		Password: password,
		RoleIDs:  []uint{role.ID},
	}

	// 9. Create user
	ctx := context.Background()
	currentUserID := uint(1) // Use system user ID for CLI operations

	_, err = userService.CreateUser(ctx, currentUserID, req)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// validateRoleByCode checks if the specified role code exists and is active in the database
func validateRoleByCode(roleRepo repoInterface.RoleRepository, roleCode string) (*model.Role, error) {
	ctx := context.Background()

	// Get role by code
	role, err := roleRepo.GetByCode(ctx, roleCode)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.CodeRecordNotFound {
			return nil, fmt.Errorf("role with code '%s' not found", roleCode)
		}
		return nil, fmt.Errorf("failed to check role existence: %w", err)
	}

	// Check if role is active (status = 1)
	if !role.IsActive() {
		return nil, fmt.Errorf("role with code '%s' is not active (status: %d)", roleCode, role.Status)
	}

	return role, nil
}

// parseUsernameEmail determines email and username from input
// If input contains @, it's treated as email, otherwise as username
func parseUsernameEmail(input string) (email, username string) {
	if strings.Contains(input, "@") {
		// Input is an email
		email = input
		// Extract username from email (part before @)
		parts := strings.Split(input, "@")
		username = parts[0]
	} else {
		// Input is a username, create a default email
		username = input
		email = fmt.Sprintf("%s@websoft9.local", username)
	}
	return email, username
}
