package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ShelfStatus struct {
	ShelfID     string    `bson:"shelf_id" json:"shelfId"`
	StockLevel  int       `bson:"stock_level" json:"stockLevel"`
	Temperature float64   `bson:"temperature" json:"temperature"`
	LastChecked time.Time `bson:"last_checked" json:"lastChecked"`
}

type CheckoutStatus struct {
	RegisterID    string    `bson:"register_id" json:"registerId"`
	IsOperational bool      `bson:"is_operational" json:"isOperational"`
	QueueLength   int       `bson:"queue_length" json:"queueLength"`
	LastChecked   time.Time `bson:"last_checked" json:"lastChecked"`
}

type StoreOperation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 店内環境データ
	Temperature  float64 `bson:"temperature" json:"temperature"`
	Humidity     float64 `bson:"humidity" json:"humidity"`
	CrowdDensity float64 `bson:"crowd_density" json:"crowdDensity"`

	// 設備状態
	Shelves   []ShelfStatus    `bson:"shelves" json:"shelves"`
	Checkouts []CheckoutStatus `bson:"checkouts" json:"checkouts"`

	// エネルギー使用
	LightingUsage float64 `bson:"lighting_usage" json:"lightingUsage"`
	ACUsage       float64 `bson:"ac_usage" json:"acUsage"`
	RefrigUsage   float64 `bson:"refrig_usage" json:"refrigUsage"`

	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
