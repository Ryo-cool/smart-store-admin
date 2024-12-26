package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	mock_repository "smart-store-admin/backend/repository/mocks"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaleRepo := mock_repository.NewMockSaleRepository(ctrl)
	mockProductRepo := mock_repository.NewMockProductRepository(ctrl)
	service := &SaleService{repo: mockSaleRepo, productRepo: mockProductRepo}
	ctx := context.Background()

	productID := primitive.NewObjectID()
	product := &models.Product{
		ID:          productID,
		Name:        "テスト商品",
		Price:       1000,
		Stock:       10,
		CO2Emission: 5.0,
	}

	tests := []struct {
		name    string
		sale    *models.Sale
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常な売上作成",
			sale: &models.Sale{
				Items: []models.SaleItem{
					{
						ProductID: productID,
						Quantity:  2,
					},
				},
			},
			mockFn: func() {
				mockProductRepo.EXPECT().
					GetByID(gomock.Any(), productID).
					Return(product, nil)
				mockProductRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
				mockSaleRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "商品なしでエラー",
			sale: &models.Sale{
				Items: []models.SaleItem{},
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "在庫不足でエラー",
			sale: &models.Sale{
				Items: []models.SaleItem{
					{
						ProductID: productID,
						Quantity:  20,
					},
				},
			},
			mockFn: func() {
				mockProductRepo.EXPECT().
					GetByID(gomock.Any(), productID).
					Return(product, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.Create(ctx, tt.sale)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(tt.sale.Items) > 0 {
					assert.Equal(t, float64(2000), tt.sale.TotalAmount)
					assert.Equal(t, float64(10), tt.sale.TotalCO2Saved)
					assert.Equal(t, float64(1000), tt.sale.Items[0].PriceAtSale)
				}
			}
		})
	}
}

func TestGetSalesByTimeOfDay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaleRepo := mock_repository.NewMockSaleRepository(ctrl)
	mockProductRepo := mock_repository.NewMockProductRepository(ctrl)
	service := &SaleService{repo: mockSaleRepo, productRepo: mockProductRepo}
	ctx := context.Background()

	tests := []struct {
		name      string
		timeOfDay string
		mockFn    func()
		wantErr   bool
	}{
		{
			name:      "正常な時間帯での取得",
			timeOfDay: "morning",
			mockFn: func() {
				mockSaleRepo.EXPECT().
					GetSalesByTimeOfDay(gomock.Any(), "morning").
					Return([]*models.Sale{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "無効な時間帯でエラー",
			timeOfDay: "invalid_time",
			mockFn:    func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			_, err := service.GetSalesByTimeOfDay(ctx, tt.timeOfDay)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnvironmentalImpactAnalytics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaleRepo := mock_repository.NewMockSaleRepository(ctrl)
	mockProductRepo := mock_repository.NewMockProductRepository(ctrl)
	service := &SaleService{repo: mockSaleRepo, productRepo: mockProductRepo}
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	sales := []*models.Sale{
		{
			TimeOfDay:     "morning",
			TotalCO2Saved: 5.0,
			Items: []models.SaleItem{
				{
					ProductID:   primitive.NewObjectID(),
					Quantity:    2,
					PriceAtSale: 1000,
				},
			},
		},
		{
			TimeOfDay:     "afternoon",
			TotalCO2Saved: 10.0,
			Items: []models.SaleItem{
				{
					ProductID:   primitive.NewObjectID(),
					Quantity:    4,
					PriceAtSale: 1000,
				},
			},
		},
	}

	expectedImpact := &models.EnvironmentalImpact{
		TotalCO2Saved:          15.0,
		AverageCO2SavedPerItem: 2.5,
		TotalEcoFriendlyItems:  6,
	}

	mockSaleRepo.EXPECT().
		GetSalesByDateRange(gomock.Any(), start, end).
		Return(sales, nil)

	mockSaleRepo.EXPECT().
		GetEnvironmentalImpactAnalytics(gomock.Any(), start, end).
		Return(expectedImpact, nil)

	impact, err := service.GetEnvironmentalImpactAnalytics(ctx, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, impact)
	assert.Equal(t, float64(15.0), impact.TotalCO2Saved)
	assert.Equal(t, float64(2.5), impact.AverageCO2SavedPerItem)
	assert.Equal(t, 6, impact.TotalEcoFriendlyItems)
}

func TestGetSalesByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaleRepo := mock_repository.NewMockSaleRepository(ctrl)
	mockProductRepo := mock_repository.NewMockProductRepository(ctrl)
	service := &SaleService{repo: mockSaleRepo, productRepo: mockProductRepo}
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	expectedCategories := map[string]int{
		"食品": 5,
		"飲料": 3,
	}

	mockSaleRepo.EXPECT().
		GetSalesByCategory(gomock.Any(), start, end).
		Return(expectedCategories, nil)

	categories, err := service.GetSalesByCategory(ctx, start, end)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
}
