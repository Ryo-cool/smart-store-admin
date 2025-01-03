package models

import "time"

// DeliveryStatus represents the status of a delivery
type DeliveryStatus string

const (
	StatusPreparing  DeliveryStatus = "配送準備中"
	StatusInProgress DeliveryStatus = "配送中"
	StatusCompleted  DeliveryStatus = "配送完了"
	StatusFailed     DeliveryStatus = "配送失敗"
)

// ValidateDeliveryStatus checks if the given status is valid
func ValidateDeliveryStatus(status string) bool {
	switch DeliveryStatus(status) {
	case StatusPreparing, StatusInProgress, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

// Delivery represents a delivery record
type Delivery struct {
	ID                    string         `json:"id" db:"id"`
	DeliveryType          string         `json:"deliveryType" db:"delivery_type"`
	Address               string         `json:"address" db:"address"`
	EstimatedDeliveryTime time.Time      `json:"estimatedDeliveryTime" db:"estimated_delivery_time"`
	ActualDeliveryTime    *time.Time     `json:"actualDeliveryTime,omitempty" db:"actual_delivery_time"`
	Status                DeliveryStatus `json:"status" db:"status"`
	Notes                 *string        `json:"notes,omitempty" db:"notes"`
	TrackingInfo          *TrackingInfo  `json:"trackingInfo,omitempty" db:"-"`
	CreatedAt             time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time      `json:"updatedAt" db:"updated_at"`
}

// TrackingInfo represents the current tracking information of a delivery
type TrackingInfo struct {
	CurrentLocation *Location `json:"currentLocation,omitempty"`
	BatteryLevel    *float64  `json:"batteryLevel,omitempty"`
	Speed           *float64  `json:"speed,omitempty"`
}

// Location represents a geographical location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DeliveryHistory represents a historical record of delivery status changes
type DeliveryHistory struct {
	ID         string    `json:"id" db:"id"`
	DeliveryID string    `json:"deliveryId" db:"delivery_id"`
	Status     string    `json:"status" db:"status"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Location   *Location `json:"location,omitempty" db:"-"`
	Note       *string   `json:"note,omitempty" db:"note"`
}

// DeliveryQuery represents query parameters for filtering deliveries
type DeliveryQuery struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Status *string `json:"status,omitempty"`
	Search *string `json:"search,omitempty"`
}

// DeliveryResponse represents the response structure for delivery queries
type DeliveryResponse struct {
	Deliveries []Delivery `json:"deliveries"`
	Total      int64      `json:"total"`
}

// DeliveryHistoryResponse represents the response structure for delivery history queries
type DeliveryHistoryResponse struct {
	History []DeliveryHistory `json:"history"`
}
