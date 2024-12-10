package main

import (
	"log"
	"smart-store-admin/backend/config"
	"smart-store-admin/backend/db"
)

func main() {
	// 設定の読み込み
	cfg := config.NewConfig()

	// データベース接続
	mongodb, err := db.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer mongodb.Close()

	log.Printf("Connected to MongoDB at %s", cfg.MongoURI)
	log.Printf("Server starting on port %s", cfg.Port)
}
