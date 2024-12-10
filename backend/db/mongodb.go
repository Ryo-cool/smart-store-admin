package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smart-store-admin/backend/models"
)

const (
	defaultDBName         = "smart_store"
	defaultConnTimeout    = 10 * time.Second
	defaultContextTimeout = 30 * time.Second
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDB は新しいMongoDBクライアントを作成します
func NewMongoDB(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	// クライアントオプションの設定
	clientOptions := options.Client().ApplyURI(uri)

	// MongoDB接続
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// 接続テスト
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	// データベース取得
	db := client.Database(defaultDBName)

	// インデックスの設���
	if err := models.SetupIndexes(db); err != nil {
		log.Printf("Failed to setup indexes: %v", err)
		return nil, err
	}

	return &MongoDB{
		client: client,
		db:     db,
	}, nil
}

// Close はMongoDBの接続を閉じます
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// GetCollection は指定されたコレクションを返します
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

// GetContext はタイムアウト付きのコンテキストを返します
func (m *MongoDB) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultContextTimeout)
}

// Collections
func (m *MongoDB) Products() *mongo.Collection {
	return m.GetCollection("products")
}

func (m *MongoDB) Deliveries() *mongo.Collection {
	return m.GetCollection("deliveries")
}

func (m *MongoDB) Sales() *mongo.Collection {
	return m.GetCollection("sales")
}

func (m *MongoDB) StoreOperations() *mongo.Collection {
	return m.GetCollection("store_operations")
}
