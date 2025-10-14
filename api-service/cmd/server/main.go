package main

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/controller"
	repoInterface "api-service/internal/interface/repository"
	serviceInterface "api-service/internal/interface/service"
	"api-service/internal/model"
	repoImpl "api-service/internal/repository"
	"api-service/internal/router"
	serviceImpl "api-service/internal/service"
	customValidator "api-service/internal/validator"
	"api-service/pkg/auth"
	"api-service/pkg/crypto"
	"api-service/pkg/database"
	"api-service/pkg/email"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/gorm"
	// _ "api-service/docs" // This line is necessary for go-swagger to find your docs!
)

//	@title			Websoft9 API Service
//	@version		1.0
//	@description	Enterprise cloud application management platform API service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Websoft9 Support
//	@contact.url	https://www.websoft9.com
//	@contact.email	support@websoft9.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/

//	@securitydefinitions.apikey BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter the token with the 'Bearer ' prefix, e.g. 'Bearer abc123'

func main() {
	// 1. Load configuration first
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2. Initialize logging system with configuration
	zapLogger := logger.NewZapLoggerWithServerConfig(
		cfg.Server.Log.LogPath,
		cfg.Server.Log.LogLevel,
		cfg.Server.Log.LogMaxSize,
		cfg.Server.Log.LogMaxBackups,
		cfg.Server.Log.LogMaxAge,
		cfg.Server.Log.LogCompress,
	)
	logger.SetDefault(zapLogger)

	zapLogger.Info("Application starting")
	zapLogger.Info("Configuration loaded successfully")

	// Set Gin mode based on configuration or environment variable
	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	} else if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Initialize authentication configuration manager
	authConfigManager, err := initAuthConfig()
	if err != nil {
		zapLogger.Fatal("Failed to initialize auth config manager", logger.String("error", err.Error()))
	}
	authConfigJson, _ := json.Marshal(authConfigManager.GetConfig())
	zapLogger.Info("Authentication configuration manager initialized successfully", logger.String("auth_config", string(authConfigJson)))

	// 4. Initialize internationalization system
	i18nInstance, err := initI18n(cfg)
	if err != nil {
		zapLogger.Fatal("Failed to initialize i18n", logger.String("error", err.Error()))
	}
	zapLogger.Info("Internationalization initialized successfully")

	_, err = initCrypto(cfg)
	if err != nil {
		zapLogger.Fatal("Failed to initialize default crypto", logger.String("error", err.Error()))
	}

	// 5. Initialize database connection and perform migrations
	dbWrapper, err := initDatabaseWrapper(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to initialize database", logger.String("error", err.Error()))
	}

	// 6. Initialize Redis, InfluxDB and JWT authentication services
	serviceConns, err := initServices(cfg, authConfigManager.GetConfig(), zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to initialize services", logger.String("error", err.Error()))
	}

	// 7. Initialize repositories, services, controllers and start HTTP server
	if err := startServer(cfg, authConfigManager, zapLogger, i18nInstance, dbWrapper.GetDB(), serviceConns); err != nil {
		zapLogger.Fatal("Failed to start server", logger.String("error", err.Error()))
	}
}

// initAuthConfig creates and returns an authentication configuration manager
// that handles OAuth2 and other authentication provider configurations
func initAuthConfig() (*config.AuthConfigManager, error) {
	authConfigPath := filepath.Join("configs", "auth.yaml")
	return config.NewAuthConfigManager(authConfigPath)
}

// initI18n initializes the internationalization system with default and supported languages
// Returns the i18n instance for use throughout the application
func initI18n(cfg *config.Config) (*i18n.I18n, error) {
	if i18nErr := i18n.InitWithConfig(cfg.I18n.DefaultLanguage, cfg.I18n.SupportedLanguages); i18nErr != nil {
		return nil, i18nErr
	}
	return i18n.GetInstance(), nil
}

func initCrypto(cfg *config.Config) (*crypto.AESCrypto, error) {
	if _, err := crypto.InitDefaultCrypto(cfg.Security.AesKey); err != nil {
		return nil, err
	}
	return crypto.GetDefaultCrypto(), nil
}

// initDatabaseWrapper establishes database connection using our enhanced SQLite manager
// and performs automatic schema migration. Supports SQLite with optimized concurrent access
func initDatabaseWrapper(cfg *config.Config, zapLogger logger.Logger) (*database.DBWrapper, error) {
	dbWrapper, err := database.InitDBWrapper(cfg)
	if err != nil {
		return nil, err
	}
	zapLogger.Info("Database connection successful",
		logger.String("type", cfg.Database.Type),
		logger.Bool("sqlite_optimized", cfg.Database.Type == "sqlite"))

	// Check if database was already initialized by the init script
	flagFile := "data/.websoft9_db_initialized"
	if _, err := os.Stat(flagFile); err == nil {
		zapLogger.Info("Database already initialized by init script, skipping auto-migration")
		return dbWrapper, nil
	}

	// Auto-migrate all database models to ensure schema consistency
	db := dbWrapper.GetDB()
	if migrateErr := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.APIToken{},
		&model.UserTwoFactor{},
		&model.AuditLog{},
		&model.UserLoginHistory{},
		&model.UserProfile{},
		&model.SystemConfig{},
		&model.Tag{},
		&model.Tagging{},
		&model.AlertRecord{},
		&model.AlertRule{},
		&model.SecretKey{},
		&model.NotificationRecord{},
		&model.NotificationChannelConfig{},
		&model.NotificationTemplate{},
	); migrateErr != nil {
		return nil, fmt.Errorf("failed to migrate database models: %v", migrateErr)
	}

	zapLogger.Info("Database schema migration completed successfully")

	// Create flag file to indicate database is initialized
	flagDir := filepath.Dir(flagFile)
	if err := os.MkdirAll(flagDir, constants.DefaultDirPerm); err != nil {
		zapLogger.Warn("Failed to create flag directory", logger.String("error", err.Error()))
	} else {
		if flagFileErr := os.WriteFile(flagFile, []byte("initialized"), constants.DefaultFilePerm); flagFileErr != nil {
			zapLogger.Warn("Failed to create initialization flag file", logger.String("error", flagFileErr.Error()))
		}
	}

	return dbWrapper, nil
}

// ServiceConnections holds all service connections that need to be closed during shutdown
type ServiceConnections struct {
	InfluxDBClient influxdb2.Client
}

// initServices initializes external service connections (Redis, InfluxDB) and JWT authentication
// Returns ServiceConnections struct containing clients that need graceful shutdown
func initServices(cfg *config.Config, authConfig *config.AuthConfig, zapLogger logger.Logger) (*ServiceConnections, error) {
	// Initialize Redis connection pool for caching and session storage
	err := redis.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}
	zapLogger.Info("Redis connection successful")

	// Initialize InfluxDB client for time-series monitoring data storage
	influxDBClient, err := database.InitInfluxDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize InfluxDB: %w", err)
	}
	zapLogger.Info("InfluxDB connection successful")

	// Initialize JWT authentication with secret key and expiration time
	auth.InitJWT(authConfig)

	return &ServiceConnections{
		InfluxDBClient: influxDBClient,
	}, nil
}

// startServer initializes all application components and starts the HTTP server
// Handles graceful shutdown when receiving interrupt signals
func startServer(
	cfg *config.Config,
	authConfigManager *config.AuthConfigManager,
	zapLogger logger.Logger,
	i18nInstance *i18n.I18n,
	db *gorm.DB,
	serviceConns *ServiceConnections,
) error {
	// Initialize request validator for input validation
	validatorInstance := validator.New()

	// Register custom validators
	if err := customValidator.RegisterCustomValidators(validatorInstance); err != nil {
		return fmt.Errorf("failed to register custom validators: %w", err)
	}

	// Initialize data access layer repositories with database connection
	repos := initRepositories(db, zapLogger)

	// Initialize business logic services with dependencies injection
	services := initBusinessServices(repos, authConfigManager, zapLogger, i18nInstance, db, cfg)

	// Initialize HTTP controllers with services and middleware
	controllers := initControllers(services, validatorInstance, zapLogger, i18nInstance, cfg)

	// Setup Gin router with all routes, middleware and security configurations
	r := router.SetupRouter(controllers, cfg, zapLogger, services.permissionService, services.apiTokenService, services.auditLogService)

	// Get server port from environment variable or configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
		if port == "" {
			port = constants.DefaultPort
		}
	}

	// Create HTTP server with timeout configurations to prevent Slowloris attacks
	// Configure server timeouts to prevent various attack vectors:
	// - ReadTimeout: Maximum duration for reading the entire request
	// - ReadHeaderTimeout: Amount of time allowed to read request headers
	// - WriteTimeout: Maximum duration before timing out writes of the response
	// - IdleTimeout: Maximum amount of time to wait for the next request
	// - MaxHeaderBytes: Maximum size of request headers
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       time.Duration(cfg.Server.Config.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.Server.Config.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.Server.Config.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.Server.Config.IdleTimeout) * time.Second,
		MaxHeaderBytes:    cfg.Server.Config.MaxHeaderBytes,
	}

	// Start HTTP server in a separate goroutine for non-blocking execution
	go func() {
		// Log server startup information including port, mode and database type
		zapLogger.Info("Websoft9 API Service starting",
			logger.String("port", port),
			logger.String("mode", gin.Mode()),
			logger.String("database", cfg.Database.Type))

		// Start listening for HTTP requests, handle server startup errors
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Failed to start server", logger.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signals (SIGINT, SIGTERM) for graceful shutdown
	// Create signal channel to capture OS interrupt signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	// Register signal handlers for SIGINT (Ctrl+C) and SIGTERM (kill)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Block until a signal is received
	<-quit

	zapLogger.Info("Shutting down server...")

	// Graceful shutdown - close services in dependency order
	// Create context with timeout for graceful shutdown operations
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultShutdownTimeout)
	// Ensure context is canceled to free resources
	defer cancel()

	// 1. First shutdown HTTP server to stop accepting new requests
	// Attempt graceful HTTP server shutdown within timeout period
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server forced to shutdown", logger.String("error", err.Error()))
	}

	// 2. Then close all database connections and services
	// Close all service connections in proper dependency order
	shutdownErrors := shutdownServices(ctx, zapLogger, db, serviceConns)
	if len(shutdownErrors) > 0 {
		// Collect and log all shutdown errors for debugging purposes
		errorMessages := make([]string, len(shutdownErrors))
		// Iterate through all shutdown errors and collect error messages
		for i, err := range shutdownErrors {
			errorMessages[i] = err.Error()
		}
		// Log consolidated error messages for operational visibility
		zapLogger.Error("Some services failed to shutdown gracefully",
			logger.String("errors", strings.Join(errorMessages, "; ")))
	}

	// Log successful server exit
	zapLogger.Info("Server exited")
	return nil
}

// repositories struct holds all repository instances for dependency injection
// Provides data access layer abstractions for different entities
type repositories struct {
	userRepo                 repoInterface.UserRepository
	roleRepo                 repoInterface.RoleRepository
	permissionRepo           repoInterface.PermissionRepository
	apiTokenRepo             repoInterface.APITokenRepository
	twoFactorRepo            repoInterface.UserTwoFactorRepository
	auditLogRepo             repoInterface.AuditLogRepository
	userProfileRepo          repoInterface.UserProfileRepository
	systemConfigRepo         repoInterface.SystemConfigRepository
	tagRepo                  repoInterface.TagRepository
	alertRepo                repoInterface.AlertRepository
	secretKeyRepo            repoInterface.SecretKeyRepository
	notificationRepo         repoInterface.NotificationRecordRepository
	notificationChannelRepo  repoInterface.NotificationChannelRepository
	notificationTemplateRepo repoInterface.NotificationTemplateRepository
}

// initRepositories creates and initializes all repository instances
// Each repository handles data access operations for its respective entity
func initRepositories(db *gorm.DB, zapLogger logger.Logger) *repositories {
	return &repositories{
		userRepo:                 repoImpl.NewUserRepository(db),
		roleRepo:                 repoImpl.NewRoleRepository(db),
		permissionRepo:           repoImpl.NewPermissionRepository(db),
		apiTokenRepo:             repoImpl.NewAPITokenRepository(db),
		twoFactorRepo:            repoImpl.NewTwoFactorRepository(db),
		auditLogRepo:             repoImpl.NewAuditLogRepository(db),
		userProfileRepo:          repoImpl.NewUserProfileRepository(db),
		systemConfigRepo:         repoImpl.NewSystemConfigRepository(db),
		tagRepo:                  repoImpl.NewTagRepository(db),
		alertRepo:                repoImpl.NewAlertRepository(db),
		secretKeyRepo:            repoImpl.NewSecretKeyRepository(db),
		notificationRepo:         repoImpl.NewNotificationRecordRepository(db),
		notificationChannelRepo:  repoImpl.NewNotificationChannelRepository(db, zapLogger),
		notificationTemplateRepo: repoImpl.NewNotificationTemplateRepository(db),
	}
}

// businessServices struct holds all service instances for dependency injection
// Provides business logic layer abstractions for different domains
type businessServices struct {
	userService                 serviceInterface.UserService
	userAuthService             serviceInterface.UserAuthService
	roleService                 serviceInterface.RoleService
	permissionService           serviceInterface.PermissionService
	apiTokenService             serviceInterface.APITokenService
	authConfigService           serviceInterface.AuthConfigService
	twoFactorService            serviceInterface.TwoFactorService
	auditLogService             serviceInterface.AuditLogService
	userProfileService          serviceInterface.UserProfileService
	systemConfigService         serviceInterface.SystemConfigService
	tagService                  serviceInterface.TagService
	alertServices               serviceInterface.AlertService
	secretKeyService            serviceInterface.SecretKeyService
	i18nService                 serviceInterface.I18nService
	notificationRecordService   serviceInterface.NotificationRecordService
	notificationChannelService  serviceInterface.NotificationChannelService
	notificationTemplateService serviceInterface.NotificationTemplateService
}

// initBusinessServices creates and initializes all service instances with their dependencies
// Wires up the service layer with repositories and other required components
func initBusinessServices(
	repos *repositories,
	authConfigManager *config.AuthConfigManager,
	zapLogger logger.Logger,
	i18nInstance *i18n.I18n,
	db *gorm.DB,
	cfg *config.Config,
) *businessServices {
	// Create OAuth2 service for external authentication providers
	oauth2Service := serviceImpl.NewOAuth2Service(authConfigManager, zapLogger)
	userService := serviceImpl.NewUserService(repos.userRepo, zapLogger)

	// Create email service first
	emailService := email.NewEmailService(cfg, zapLogger)

	services := &businessServices{
		userService: userService,
		userAuthService: serviceImpl.NewUserAuthService(
			repos.userRepo,
			repos.apiTokenRepo,
			repos.userProfileRepo,
			repos.systemConfigRepo,
			oauth2Service,
			zapLogger,
			cfg,
			authConfigManager,
		),
		roleService:                serviceImpl.NewRoleService(repos.roleRepo, repos.permissionRepo, db, zapLogger),
		permissionService:          serviceImpl.NewPermissionService(repos.permissionRepo, db, zapLogger),
		apiTokenService:            serviceImpl.NewAPITokenService(repos.apiTokenRepo, authConfigManager, db, zapLogger),
		authConfigService:          serviceImpl.NewAuthConfigService(authConfigManager, zapLogger),
		twoFactorService:           serviceImpl.NewTwoFactorService(repos.twoFactorRepo, db, zapLogger),
		auditLogService:            serviceImpl.NewAuditLogService(repos.auditLogRepo, userService, db, zapLogger, cfg),
		systemConfigService:        serviceImpl.NewSystemConfigService(repos.systemConfigRepo, cfg, db, zapLogger),
		userProfileService:         serviceImpl.NewUserProfileService(repos.userProfileRepo, zapLogger, i18nInstance),
		tagService:                 serviceImpl.NewTagService(repos.tagRepo, db, zapLogger, i18nInstance),
		alertServices:              serviceImpl.NewAlertService(repos.alertRepo, zapLogger, i18nInstance),
		secretKeyService:           serviceImpl.NewSecretKeyService(repos.secretKeyRepo, zapLogger, i18nInstance, cfg),
		i18nService:                serviceImpl.NewI18nService(repos.userProfileRepo, zapLogger),
		notificationRecordService:  serviceImpl.NewNotificationRecordService(repos.notificationRepo, zapLogger),
		notificationChannelService: serviceImpl.NewNotificationChannelService(repos.notificationChannelRepo, zapLogger, emailService),
	}

	// Create notification template service after channel service is created
	services.notificationTemplateService = serviceImpl.NewNotificationTemplateService(repos.notificationTemplateRepo, services.notificationChannelService, zapLogger)

	return services
}

// initControllers creates and initializes all HTTP controllers with their dependencies
// Controllers handle HTTP requests and responses, delegating business logic to services
func initControllers(
	services *businessServices,
	validatorInstance *validator.Validate,
	zapLogger logger.Logger,
	i18nInstance *i18n.I18n,
	cfg *config.Config,
) *router.Controllers {
	// Create OAuth2 service for SecurityController (placeholder for future implementation)
	var oauth2ServiceForSecurity *serviceImpl.OAuth2Service

	return &router.Controllers{
		UserController:     controller.NewUserController(services.userService, zapLogger, i18nInstance, validatorInstance),
		UserAuthController: controller.NewUserAuthController(services.userAuthService, validatorInstance, zapLogger),
		I18nController:     controller.NewI18nController(services.i18nService, validatorInstance, zapLogger),
		RolePermissionController: controller.NewRolePermissionController(
			services.roleService,
			services.permissionService,
			validatorInstance,
			zapLogger,
		),
		SecurityController: controller.NewSecurityController(
			services.apiTokenService,
			services.authConfigService,
			oauth2ServiceForSecurity,
			services.twoFactorService,
			validatorInstance,
			zapLogger,
		),
		HealthController:      controller.NewHealthController(cfg),
		AuditLogController:    controller.NewAuditLogController(services.auditLogService, validatorInstance, zapLogger),
		UserProfileController: controller.NewUserProfileController(services.userProfileService, zapLogger, i18nInstance, validatorInstance),
		SystemConfigController: controller.NewSystemConfigController(
			services.systemConfigService,
			validatorInstance,
			zapLogger,
		),
		AlertController: controller.NewAlertController(services.alertServices, zapLogger, i18nInstance, validatorInstance),
		TagController: controller.NewTagController(
			services.tagService,
			validatorInstance,
			zapLogger,
			i18nInstance,
		),
		SecretKeyController:            controller.NewSecretKeyController(services.secretKeyService, zapLogger, i18nInstance, validatorInstance),
		NotificationRecordController:   controller.NewNotificationRecordController(services.notificationRecordService, validatorInstance, zapLogger),
		NotificationChannelController:  controller.NewNotificationChannelController(services.notificationChannelService, zapLogger, validatorInstance),
		NotificationTemplateController: controller.NewNotificationTemplateController(services.notificationTemplateService, validatorInstance, zapLogger),
	}
}

// shutdownServices gracefully shuts down all services and database connections
// Returns a slice of errors encountered during shutdown for logging purposes
func shutdownServices(ctx context.Context, zapLogger logger.Logger, db *gorm.DB, serviceConns *ServiceConnections) []error {
	var shutdownErrors []error

	zapLogger.Info("Starting graceful shutdown of services...")

	// 1. Close Redis connection pool
	zapLogger.Info("Closing Redis connection...")
	if err := redis.Close(); err != nil {
		shutdownErr := errors.NewAppErrorWrapError(err, errors.CodeInternalError)
		shutdownErrors = append(shutdownErrors, shutdownErr)
		zapLogger.Error("Redis shutdown error", logger.String("error", shutdownErr.Error()))
	} else {
		zapLogger.Info("Redis connection closed successfully")
	}

	// 2. Close InfluxDB client connection
	if serviceConns.InfluxDBClient != nil {
		zapLogger.Info("Closing InfluxDB connection...")
		serviceConns.InfluxDBClient.Close()
		zapLogger.Info("InfluxDB connection closed successfully")
	}

	// 3. Close database connection (GORM with underlying sql.DB)
	if db != nil {
		zapLogger.Info("Closing database connection...")
		sqlDB, err := db.DB()
		if err != nil {
			shutdownErr := errors.NewAppErrorWrapError(err, errors.CodeInternalError)
			shutdownErrors = append(shutdownErrors, shutdownErr)
			zapLogger.Error("Database connection retrieval error", logger.String("error", shutdownErr.Error()))
		} else {
			// Set timeout for database connection closure to prevent hanging
			dbCloseCtx, dbCancel := context.WithTimeout(ctx, constants.DefaultReadHeaderTimeout)
			defer dbCancel()

			// Use goroutine with timeout context to close database connection safely
			done := make(chan error, 1)
			go func() {
				done <- sqlDB.Close()
			}()

			select {
			case err := <-done:
				if err != nil {
					shutdownErr := errors.NewAppErrorWrapError(err, errors.CodeInternalError)
					shutdownErrors = append(shutdownErrors, shutdownErr)
					zapLogger.Error("Database shutdown error", logger.String("error", shutdownErr.Error()))
				} else {
					zapLogger.Info("Database connection closed successfully")
				}
			case <-dbCloseCtx.Done():
				shutdownErr := errors.NewAppError(errors.CodeInternalError)
				shutdownErrors = append(shutdownErrors, shutdownErr)
				zapLogger.Error("Database shutdown timeout", logger.String("error", shutdownErr.Error()))
			}
		}
	}

	if len(shutdownErrors) == 0 {
		zapLogger.Info("All services shut down successfully")
	} else {
		zapLogger.Warn("Service shutdown completed with errors", logger.Int("error_count", len(shutdownErrors)))
	}

	return shutdownErrors
}
