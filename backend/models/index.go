package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupIndexes はMongoDBのインデックスを設定します
func SetupIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Products collection indexes
	productIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"name": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: map[string]interface{}{
				"category": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"shelf_location": 1,
			},
		},
	}

	if _, err := db.Collection("products").Indexes().CreateMany(ctx, productIndexes); err != nil {
		log.Printf("Failed to create product indexes: %v", err)
		return err
	}

	// Deliveries collection indexes
	deliveryIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"robot_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"status": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"created_at": 1,
			},
		},
	}

	if _, err := db.Collection("deliveries").Indexes().CreateMany(ctx, deliveryIndexes); err != nil {
		log.Printf("Failed to create delivery indexes: %v", err)
		return err
	}

	// Sales collection indexes
	saleIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"created_at": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"payment_method": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"time_of_day": 1,
			},
		},
	}

	if _, err := db.Collection("sales").Indexes().CreateMany(ctx, saleIndexes); err != nil {
		log.Printf("Failed to create sale indexes: %v", err)
		return err
	}

	// Store operations collection indexes
	storeOpIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"timestamp": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"shelves.shelf_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"checkouts.register_id": 1,
			},
		},
	}

	if _, err := db.Collection("store_operations").Indexes().CreateMany(ctx, storeOpIndexes); err != nil {
		log.Printf("Failed to create store operation indexes: %v", err)
		return err
	}

	return nil
}
