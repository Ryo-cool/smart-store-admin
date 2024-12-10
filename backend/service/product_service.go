package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

// ProductService は商品サービスを表します
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService は新しい商品サービスを作成します
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// CreateProduct は新しい商品を作成します
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price < 0 {
		return errors.New("price must be non-negative")
	}
	if product.Stock < 0 {
		return errors.New("stock must be non-negative")
	}

	return s.repo.Create(ctx, product)
}

// GetProduct は指定されたIDの商品を取得します
func (s *ProductService) GetProduct(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// ListProducts は商品のリストを取得します
func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int64) ([]*models.Product, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	return s.repo.List(ctx, skip, pageSize)
}

// UpdateProduct は商品情報を更新します
func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if product.ID.IsZero() {
		return errors.New("product ID is required")
	}
	if product.Price < 0 {
		return errors.New("price must be non-negative")
	}
	if product.Stock < 0 {
		return errors.New("stock must be non-negative")
	}

	return s.repo.Update(ctx, product)
}

// DeleteProduct は商品を削除します
func (s *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

// GetProductsByCategory は指定されたカテゴリの商品を取得します
func (s *ProductService) GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}
	return s.repo.GetByCategory(ctx, category)
}

// CheckLowStock は在庫が少ない商品を取得し、必要に応じてアラートを生成します
func (s *ProductService) CheckLowStock(ctx context.Context) ([]*models.Product, error) {
	products, err := s.repo.GetLowStock(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: 在庫アラートの生成ロジックを実装
	// 例: 通知システムとの連携、自動発注システムとの連携など

	return products, nil
}

// UpdateStock は商品の在庫数を更新します
func (s *ProductService) UpdateStock(ctx context.Context, id primitive.ObjectID, quantity int) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if product.Stock+quantity < 0 {
		return errors.New("insufficient stock")
	}

	product.Stock += quantity
	return s.repo.Update(ctx, product)
}
