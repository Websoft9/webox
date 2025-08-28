package main

import (
	"api-service/internal/config"
	"api-service/internal/controller"
	repoInterface "api-service/internal/interface/repository"
	serviceInterface "api-service/internal/interface/service"
	"api-service/internal/model"
	repoImpl "api-service/internal/repository"
	"api-service/internal/router"
	serviceImpl "api-service/internal/service"
	"api-service/pkg/auth"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	_ "api-service/docs" // This line is necessary for go-swagger to find your docs!
)

func main() {
	// 1. Initialize logging system
	zapLogger := logger.NewDefaultZapLogger()
	logger.SetDefault(zapLogger)

	zapLogger.Info("Application starting")

	// 2. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	zapLogger.Info("Configuration loaded successfully")

	// 3. Initialize authentication configuration manager
	authConfigManager, err := initAuthConfig()
	if err != nil {
		log.Fatal("Failed to initialize auth config manager:", err)
	}
	zapLogger.Info("Authentication configuration manager initialized successfully")

	// 4. Initialize i18n
	i18nInstance, err := initI18n(cfg)
	if err != nil {
		log.Fatal("Failed to initialize i18n:", err)
	}
	zapLogger.Info("Internationalization initialized successfully")

	// 5. Initialize database and perform migrations
	db, err := initDatabase(cfg, zapLogger)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 6. Initialize other services
	if err := initServices(cfg, zapLogger); err != nil {
		log.Fatal("Failed to initialize services:", err)
	}

	// 7. Initialize and start server
	if err := startServer(cfg, authConfigManager, zapLogger, i18nInstance, db); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initAuthConfig() (*config.AuthConfigManager, error) {
	authConfigPath := filepath.Join("configs", "auth.yaml")
	return config.NewAuthConfigManager(authConfigPath)
}

func initI18n(cfg *config.Config) (*i18n.I18n, error) {
	if i18nErr := i18n.InitWithConfig(cfg.I18n.DefaultLanguage, cfg.I18n.SupportedLanguages); i18nErr != nil {
		return nil, i18nErr
	}
	return i18n.GetInstance(), nil
}

func initDatabase(cfg *config.Config, zapLogger logger.Logger) (*gorm.DB, error) {
	db, err := utils.InitDB(cfg)
	if err != nil {
		return nil, err
	}
	zapLogger.Info("Database connection successful")

	// Database migration - add all required models
	if migrateErr := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.APIToken{},
		&model.UserTwoFactor{},
	); migrateErr != nil {
		return nil, migrateErr
	}
	zapLogger.Info("Database migration completed")
	return db, nil
}

func initServices(cfg *config.Config, zapLogger logger.Logger) error {
	// Initialize Redis
	_, err := utils.InitRedis(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}
	zapLogger.Info("Redis connection successful")

	// Initialize InfluxDB
	_, err = utils.InitInfluxDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize InfluxDB: %w", err)
	}
	zapLogger.Info("InfluxDB connection successful")

	// Initialize JWT authentication
	auth.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpireTime)
	return nil
}

func startServer(cfg *config.Config, authConfigManager *config.AuthConfigManager, zapLogger logger.Logger, i18nInstance *i18n.I18n, db *gorm.DB) error {
	// Initialize validator
	validatorInstance := validator.New()

	// Initialize repositories
	repos := initRepositories(db)

	// Initialize services
	services := initBusinessServices(repos, authConfigManager, zapLogger, i18nInstance, db)

	// Initialize controllers
	controllers := initControllers(services, validatorInstance, zapLogger, i18nInstance)

	// Initialize router
	r := router.SetupRouter(controllers, cfg, zapLogger, services.permissionService)

	// Start server
	zapLogger.Info("Server starting", logger.String("port", cfg.Server.Port))
	return r.Run(":" + cfg.Server.Port)
}

type repositories struct {
	userRepo       repoInterface.UserRepository
	roleRepo       repoInterface.RoleRepository
	permissionRepo repoInterface.PermissionRepository
	apiTokenRepo   repoInterface.APITokenRepository
	twoFactorRepo  repoInterface.UserTwoFactorRepository
}

func initRepositories(db *gorm.DB) *repositories {
	return &repositories{
		userRepo:       repoImpl.NewUserRepository(db),
		roleRepo:       repoImpl.NewRoleRepository(db),
		permissionRepo: repoImpl.NewPermissionRepository(db),
		apiTokenRepo:   repoImpl.NewAPITokenRepository(db),
		twoFactorRepo:  repoImpl.NewTwoFactorRepository(db),
	}
}

type businessServices struct {
	userService       serviceInterface.UserService
	roleService       serviceInterface.RoleService
	permissionService serviceInterface.PermissionService
	apiTokenService   serviceInterface.APITokenService
	authConfigService serviceInterface.AuthConfigService
	twoFactorService  serviceInterface.TwoFactorService
}

func initBusinessServices(repos *repositories, authConfigManager *config.AuthConfigManager, zapLogger logger.Logger, i18nInstance *i18n.I18n, db *gorm.DB) *businessServices {
	return &businessServices{
		userService:       serviceImpl.NewUserService(repos.userRepo, zapLogger),
		roleService:       serviceImpl.NewRoleService(repos.roleRepo, repos.permissionRepo, db, zapLogger, i18nInstance),
		permissionService: serviceImpl.NewPermissionService(repos.permissionRepo, db, zapLogger, i18nInstance),
		apiTokenService:   serviceImpl.NewAPITokenService(repos.apiTokenRepo, db, zapLogger, i18nInstance),
		authConfigService: serviceImpl.NewAuthConfigService(authConfigManager, zapLogger),
		twoFactorService:  serviceImpl.NewTwoFactorService(repos.twoFactorRepo, db, zapLogger, i18nInstance),
	}
}

func initControllers(services *businessServices, validatorInstance *validator.Validate, zapLogger logger.Logger, i18nInstance *i18n.I18n) *router.Controllers {
	return &router.Controllers{
		UserController:       controller.NewUserController(services.userService, zapLogger),
		I18nController:       controller.NewI18nController(),
		RoleController:       controller.NewRoleController(services.roleService, validatorInstance, zapLogger, i18nInstance),
		PermissionController: controller.NewPermissionController(services.permissionService, validatorInstance, zapLogger, i18nInstance),
		APITokenController:   controller.NewAPITokenController(services.apiTokenService, validatorInstance, zapLogger, i18nInstance),
		AuthConfigController: controller.NewAuthConfigController(services.authConfigService, validatorInstance, zapLogger, i18nInstance),
		TwoFactorController:  controller.NewTwoFactorController(services.twoFactorService, validatorInstance, zapLogger, i18nInstance),
	}
}
