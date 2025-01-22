package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/onoderaryou/smart-store-admin/backend/models"
	"github.com/onoderaryou/smart-store-admin/backend/repository"
)

type MockProductRepository struct {
	mock.Mock
}

var _ repository.ProductRepository = (*MockProductRepository)(nil)

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context, skip, limit int64) ([]*models.Product, error) {
	args := m.Called(ctx, skip, limit)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) GetByCategory(ctx context.Context, category string) ([]*models.Product, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetLowStock(ctx context.Context) ([]*models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := ProductService{repo: mockRepo}
	ctx := context.Background()

	tests := []struct {
		name    string
		product *models.Product
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常な商品作成",
			product: &models.Product{
				Name:        "テスト商品",
				Price:       1000,
				Stock:       10,
				Category:    "テストカテゴリ",
				CO2Emission: 5.0,
				RecycleRate: 80.0,
			},
			mockFn: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "商品名なしでエラー",
			product: &models.Product{
				Price: 1000,
				Stock: 10,
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "負の価格でエラー",
			product: &models.Product{
				Name:  "テスト商品",
				Price: -1000,
				Stock: 10,
			},
			mockFn:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.CreateProduct(ctx, tt.product)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := ProductService{repo: mockRepo}
	ctx := context.Background()
	productID := primitive.NewObjectID()

	tests := []struct {
		name     string
		product  *models.Product
		quantity int
		mockFn   func()
		wantErr  bool
	}{
		{
			name: "在庫の正常な更新",
			product: &models.Product{
				ID:    productID,
				Name:  "テスト商品",
				Stock: 10,
			},
			quantity: -5,
			mockFn: func() {
				mockRepo.On("GetByID", ctx, productID).Return(&models.Product{
					ID:    productID,
					Name:  "テスト商品",
					Stock: 10,
				}, nil)
				mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "在庫不足でエラー",
			product: &models.Product{
				ID:    productID,
				Name:  "テスト商品",
				Stock: 5,
			},
			quantity: -10,
			mockFn: func() {
				mockRepo.On("GetByID", ctx, productID).Return(&models.Product{
					ID:    productID,
					Name:  "テスト商品",
					Stock: 5,
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.UpdateStock(ctx, tt.product.ID, tt.quantity)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetProductsByCategory(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := ProductService{repo: mockRepo}
	ctx := context.Background()

	expectedProducts := []*models.Product{
		{
			ID:       primitive.NewObjectID(),
			Name:     "テスト商品1",
			Category: "テストカテゴリ",
		},
		{
			ID:       primitive.NewObjectID(),
			Name:     "テスト商品2",
			Category: "テストカテゴリ",
		},
	}

	tests := []struct {
		name     string
		category string
		mockFn   func()
		want     []*models.Product
		wantErr  bool
	}{
		{
			name:     "カテゴリによる商品取得",
			category: "テストカテゴリ",
			mockFn: func() {
				mockRepo.On("GetByCategory", ctx, "テストカテゴリ").Return(expectedProducts, nil)
			},
			want:    expectedProducts,
			wantErr: false,
		},
		{
			name:     "空のカテゴリでエラー",
			category: "",
			mockFn:   func() {},
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetProductsByCategory(ctx, tt.category)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
