package database

import (
	"api-service/internal/config"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Version   string    `gorm:"uniqueIndex;size:50"`
	Name      string    `gorm:"size:100"`
	Batch     int       `gorm:"default:1"`
	Applied   bool      `gorm:"default:false"`
	AppliedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// MigrationFunc represents a migration function
type MigrationFunc func(*gorm.DB) error

// MigrationEntry represents a migration entry
type MigrationEntry struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []MigrationEntry
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []MigrationEntry{},
	}
}

// AddMigration adds a migration to the migrator
func (m *Migrator) AddMigration(version, name string, up, down MigrationFunc) {
	m.migrations = append(m.migrations, MigrationEntry{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

// InitMigrationTable initializes the migration table
func (m *Migrator) InitMigrationTable() error {
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migration table: %v", err)
	}
	return nil
}

// Migrate runs all pending migrations
func (m *Migrator) Migrate() error {
	if err := m.InitMigrationTable(); err != nil {
		return err
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// Get applied migrations
	var appliedMigrations []Migration
	if err := m.db.Where("applied = ?", true).Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Version] = true
	}

	// Run pending migrations
	currentBatch := m.getCurrentBatch() + 1
	for _, migration := range m.migrations {
		if appliedMap[migration.Version] {
			continue
		}

		log.Printf("Running migration: %s - %s", migration.Version, migration.Name)

		// Begin transaction
		tx := m.db.Begin()

		// Run migration
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %v", migration.Version, err)
		}

		// Record migration
		migrationRecord := Migration{
			Version:   migration.Version,
			Name:      migration.Name,
			Batch:     currentBatch,
			Applied:   true,
			AppliedAt: time.Now(),
		}

		if err := tx.Create(&migrationRecord).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %v", migration.Version, err)
		}

		// Commit transaction
		tx.Commit()
		log.Printf("Migration completed: %s", migration.Version)
	}

	log.Printf("All migrations completed successfully")
	return nil
}

// Rollback rollbacks the last batch of migrations
func (m *Migrator) Rollback() error {
	if err := m.InitMigrationTable(); err != nil {
		return err
	}

	// Get last batch number
	lastBatch := m.getCurrentBatch()
	if lastBatch == 0 {
		log.Printf("No migrations to rollback")
		return nil
	}

	// Get migrations from last batch
	var migrations []Migration
	if err := m.db.Where("batch = ? AND applied = ?", lastBatch, true).Order("version DESC").Find(&migrations).Error; err != nil {
		return fmt.Errorf("failed to get last batch migrations: %v", err)
	}

	// Create migration map
	migrationMap := make(map[string]MigrationEntry)
	for _, migration := range m.migrations {
		migrationMap[migration.Version] = migration
	}

	// Rollback migrations in reverse order
	for _, migration := range migrations {
		migrationEntry, exists := migrationMap[migration.Version]
		if !exists {
			continue
		}

		log.Printf("Rolling back migration: %s - %s", migration.Version, migration.Name)

		// Begin transaction
		tx := m.db.Begin()

		// Run rollback
		if err := migrationEntry.Down(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("rollback of migration %s failed: %v", migration.Version, err)
		}

		// Update migration record
		migrationCopy := migration
		migrationCopy.Applied = false
		if err := tx.Save(&migrationCopy).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update migration record %s: %v", migration.Version, err)
		}

		// Commit transaction
		tx.Commit()
		log.Printf("Migration rolled back: %s", migration.Version)
	}

	log.Printf("Rollback completed successfully")
	return nil
}

// Status shows migration status
func (m *Migrator) Status() ([]Migration, error) {
	if err := m.InitMigrationTable(); err != nil {
		return nil, err
	}

	var migrations []Migration
	if err := m.db.Order("version").Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get migration status: %v", err)
	}

	return migrations, nil
}

// getCurrentBatch returns the current batch number
func (m *Migrator) getCurrentBatch() int {
	var lastMigration Migration
	if err := m.db.Where("applied = ?", true).Order("batch DESC").First(&lastMigration).Error; err != nil {
		return 0
	}
	return lastMigration.Batch
}

// RunMigrations runs database migrations for the current configuration
func RunMigrations(cfg *config.Config) error {
	db, err := InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	migrator := NewMigrator(db)

	// Add migrations based on database type
	switch strings.ToLower(cfg.Database.Type) {
	case DatabaseTypeSQLite:
		addSQLiteMigrations(migrator)
	case DatabaseTypeMySQL:
		addMySQLMigrations(migrator)
	case DatabaseTypePostgres, DatabaseTypePostgreSQL:
		addPostgresMigrations(migrator)
	default:
		return fmt.Errorf("unsupported database type for migrations: %s", cfg.Database.Type)
	}

	return migrator.Migrate()
}

// addSQLiteMigrations adds SQLite-specific migrations
func addSQLiteMigrations(m *Migrator) {
	m.AddMigration("20240101_000001", "Create initial tables", func(db *gorm.DB) error {
		// Enable foreign key constraints for SQLite
		db.Exec("PRAGMA foreign_keys = ON")

		// Add your SQLite-specific migration logic here
		// For example, auto-migrate all models
		// return db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{})
		return nil
	}, func(db *gorm.DB) error {
		// Rollback logic
		return nil
	})
}

// addMySQLMigrations adds MySQL-specific migrations
func addMySQLMigrations(m *Migrator) {
	m.AddMigration("20240101_000001", "Create initial tables", func(db *gorm.DB) error {
		// Add your MySQL-specific migration logic here
		// For example, auto-migrate all models
		// return db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{})
		return nil
	}, func(db *gorm.DB) error {
		// Rollback logic
		return nil
	})
}

// addPostgresMigrations adds PostgreSQL-specific migrations
func addPostgresMigrations(m *Migrator) {
	m.AddMigration("20240101_000001", "Create initial tables", func(db *gorm.DB) error {
		// Add your PostgreSQL-specific migration logic here
		// For example, auto-migrate all models
		// return db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{})
		return nil
	}, func(db *gorm.DB) error {
		// Rollback logic
		return nil
	})
}
