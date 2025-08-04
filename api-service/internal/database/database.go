package database

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"fmt"
	"os"
	"path/filepath"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, constants.DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %v", err)
	}

	return db, nil
}

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return rdb, nil
}

func InitInfluxDB(cfg *config.Config) (influxdb2.Client, error) {
	client := influxdb2.NewClient(cfg.InfluxDB.URL, cfg.InfluxDB.Token)
	return client, nil
}
