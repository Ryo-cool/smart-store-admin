package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DeliveryStatus string

const (
	StatusPending    DeliveryStatus = "pending"
	StatusInProgress DeliveryStatus = "in_progress"
	StatusCompleted  DeliveryStatus = "completed"
	StatusFailed     DeliveryStatus = "failed"
)

type Location struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

type Delivery struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RobotID string             `bson:"robot_id" json:"robotId"`
	Status  DeliveryStatus     `bson:"status" json:"status"`

	// ロボット/ドローン情報
	CurrentLocation Location `bson:"current_location" json:"currentLocation"`
	BatteryLevel    float64  `bson:"battery_level" json:"batteryLevel"`

	// 配送情報
	StartLocation Location   `bson:"start_location" json:"startLocation"`
	EndLocation   Location   `bson:"end_location" json:"endLocation"`
	OptimalRoute  []Location `bson:"optimal_route" json:"optimalRoute"`

	// エネルギー効率
	EnergyUsage     float64 `bson:"energy_usage" json:"energyUsage"`
	DistanceCovered float64 `bson:"distance_covered" json:"distanceCovered"`

	StartedAt   time.Time  `bson:"started_at" json:"startedAt"`
	CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completedAt"`
	CreatedAt   time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updatedAt"`
}
