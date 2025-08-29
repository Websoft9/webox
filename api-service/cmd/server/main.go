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
	"api-service/pkg/auth"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	_ "api-service/docs" // This line is necessary for go-swagger to find your docs!
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

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. Load configuration first
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
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

	// Set Gin mode
	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	} else if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

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
	controllers := initControllers(services, validatorInstance, zapLogger, i18nInstance, cfg)

	// Initialize router with complete functionality
	r := router.SetupRouter(controllers, cfg, zapLogger, services.permissionService)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
		if port == "" {
			port = constants.DefaultPort
		}
	}

	// 创建 HTTP 服务器 - 根据配置参数设置服务器参数，防止 Slowloris 攻击
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       time.Duration(cfg.Server.Config.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.Server.Config.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.Server.Config.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.Server.Config.IdleTimeout) * time.Second,
		MaxHeaderBytes:    cfg.Server.Config.MaxHeaderBytes,
	}

	// 启动服务器
	go func() {
		zapLogger.Info("Websoft9 API Service starting",
			logger.String("port", port),
			logger.String("mode", gin.Mode()),
			logger.String("database", cfg.Database.Type))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server forced to shutdown", logger.String("error", err.Error()))
		return err
	}

	zapLogger.Info("Server exited")
	return nil
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

func initControllers(services *businessServices, validatorInstance *validator.Validate, zapLogger logger.Logger, i18nInstance *i18n.I18n, cfg *config.Config) *router.Controllers {
	return &router.Controllers{
		UserController:       controller.NewUserController(services.userService, zapLogger),
		I18nController:       controller.NewI18nController(),
		RoleController:       controller.NewRoleController(services.roleService, validatorInstance, zapLogger, i18nInstance),
		PermissionController: controller.NewPermissionController(services.permissionService, validatorInstance, zapLogger, i18nInstance),
		APITokenController:   controller.NewAPITokenController(services.apiTokenService, validatorInstance, zapLogger, i18nInstance),
		AuthConfigController: controller.NewAuthConfigController(services.authConfigService, validatorInstance, zapLogger, i18nInstance),
		TwoFactorController:  controller.NewTwoFactorController(services.twoFactorService, validatorInstance, zapLogger, i18nInstance),
		HealthController:     controller.NewHealthController(cfg),
	}
}
