package service

import (
	"context"
	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"

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
	products, err := ps.repo.GetByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	return ps.repo.Create(ctx, product)
}

func (ps *ProductService) UpdateStock(ctx context.Context, id primitive.ObjectID, quantity int) error {
	product, err := ps.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	product.Stock = quantity
	return ps.repo.Update(ctx, product)
}

func (ps *ProductService) GetProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	return ps.repo.GetByID(ctx, id)
}

func (ps *ProductService) List(ctx context.Context, skip, limit int64) ([]*models.Product, error) {
	return ps.repo.List(ctx, skip, limit)
}

func (ps *ProductService) Update(ctx context.Context, product *models.Product) error {
	return ps.repo.Update(ctx, product)
}

func (ps *ProductService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return ps.repo.Delete(ctx, id)
}
