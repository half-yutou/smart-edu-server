package main

import (
	"log"
	"smarteduhub/internal/pkg/oss"

	"smarteduhub/internal/config"
	"smarteduhub/internal/pkg/auth"
	"smarteduhub/internal/pkg/database"
	"smarteduhub/internal/pkg/validator"
	"smarteduhub/internal/router"
)

func main() {
	// 1. 初始化配置
	configPath := "./config.json"
	if err := config.Load(configPath); err != nil {
		log.Fatalf("Config load failed: %v", err)
	}

	auth.InitSaToken()

	oss.Init()

	if err := validator.Init(); err != nil {
		log.Fatalf("Validator init failed: %v", err)
	}

	// 2. 初始化数据库
	database.InitDB()

	// 3. 初始化路由
	r := router.InitRouter()

	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
