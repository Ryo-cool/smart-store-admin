package service

import (
	"context"

	"smart-store-admin/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error)
	List(ctx context.Context, skip, limit int64) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
