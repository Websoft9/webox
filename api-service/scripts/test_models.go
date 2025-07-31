package main

import (
	"api-service/internal/config"
	"api-service/internal/database"
	"api-service/internal/model"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Testing Websoft9 Models...")

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 执行自动迁移
	fmt.Println("Running AutoMigrate...")
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ AutoMigrate completed successfully!")

	// 测试创建一些基础数据
	fmt.Println("Creating test data...")

	// 创建用户组
	userGroup := &model.UserGroup{
		Name:        "测试用户组",
		Code:        "test_group",
		Description: "这是一个测试用户组",
	}
	if err := db.Create(userGroup).Error; err != nil {
		log.Printf("Failed to create user group: %v", err)
	} else {
		fmt.Printf("✅ Created user group: %s (ID: %d)\n", userGroup.Name, userGroup.ID)
	}

	// 创建角色
	role := &model.Role{
		Name:        "测试角色",
		Code:        "test_role",
		Description: "这是一个测试角色",
	}
	if err := db.Create(role).Error; err != nil {
		log.Printf("Failed to create role: %v", err)
	} else {
		fmt.Printf("✅ Created role: %s (ID: %d)\n", role.Name, role.ID)
	}

	// 创建权限
	permission := &model.Permission{
		Name:        "测试权限",
		Code:        "test.permission",
		Module:      "test",
		Action:      "read",
		Description: "这是一个测试权限",
	}
	if err := db.Create(permission).Error; err != nil {
		log.Printf("Failed to create permission: %v", err)
	} else {
		fmt.Printf("✅ Created permission: %s (ID: %d)\n", permission.Name, permission.ID)
	}

	// 创建应用分类
	category := &model.AppStoreCategory{
		Name:        "测试分类",
		Code:        "test_category",
		Description: "这是一个测试应用分类",
	}
	if err := db.Create(category).Error; err != nil {
		log.Printf("Failed to create app category: %v", err)
	} else {
		fmt.Printf("✅ Created app category: %s (ID: %d)\n", category.Name, category.ID)
	}

	// 创建系统配置
	config := &model.SystemConfig{
		ConfigKey:   "test.config",
		ConfigValue: "test_value",
		ConfigType:  "STRING",
		Category:    "test",
		Description: "这是一个测试配置",
	}
	if err := db.Create(config).Error; err != nil {
		log.Printf("Failed to create system config: %v", err)
	} else {
		fmt.Printf("✅ Created system config: %s = %s\n", config.ConfigKey, config.ConfigValue)
	}

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("Database models are working correctly.")
}
