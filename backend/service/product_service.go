package service

import (
	"context"
	"errors"

	"github.com/onoderaryou/smart-store-admin/backend/models"
	"github.com/onoderaryou/smart-store-admin/backend/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	repo repository.ProductRepository
}

type ProductServiceInterface interface {
	GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateStock(ctx context.Context, id primitive.ObjectID, quantity int) error
	GetProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error)
	List(ctx context.Context, skip, limit int64) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (ps *ProductService) GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}
	products, err := ps.repo.GetByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price < 0 {
		return errors.New("price must be non-negative")
	}
	return ps.repo.Create(ctx, product)
}

func (ps *ProductService) UpdateStock(ctx context.Context, id primitive.ObjectID, quantity int) error {
	product, err := ps.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	product.Stock = newStock
	return ps.repo.Update(ctx, product)
}

func (ps *ProductService) GetProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	return ps.repo.GetByID(ctx, id)
}

func (ps *ProductService) List(ctx context.Context, skip, limit int64) ([]*models.Product, error) {
	if skip < 0 {
		return nil, errors.New("skip must be non-negative")
	}
	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	return ps.repo.List(ctx, skip, limit)
}

func (ps *ProductService) Update(ctx context.Context, product *models.Product) error {
	if product.ID.IsZero() {
		return errors.New("product ID is required")
	}
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price < 0 {
		return errors.New("price must be non-negative")
	}
	return ps.repo.Update(ctx, product)
}

func (ps *ProductService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if id.IsZero() {
		return errors.New("product ID is required")
	}
	return ps.repo.Delete(ctx, id)
}
