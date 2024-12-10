package config

import (
	"os"
)

// Config はアプリケーションの設定を保持します
type Config struct {
	MongoURI string
	Port     string
}

// NewConfig は新しい設定を作成します
func NewConfig() *Config {
	return &Config{
		// 環境変数から設定を読み込み、デフォルト値を設定
		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		Port:     getEnv("PORT", "8080"),
	}
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
