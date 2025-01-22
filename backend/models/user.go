package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleStaff  Role = "staff"
	RoleViewer Role = "viewer"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`
	Name        string             `bson:"name" json:"name"`
	Picture     string             `bson:"picture" json:"picture"`
	GoogleID    string             `bson:"google_id" json:"google_id"`
	Role        Role               `bson:"role" json:"role"`
	LastLoginAt time.Time          `bson:"last_login_at" json:"last_login_at"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID      primitive.ObjectID `json:"id"`
	Email   string             `json:"email"`
	Name    string             `json:"name"`
	Picture string             `json:"picture"`
	Role    Role               `json:"role"`
}
