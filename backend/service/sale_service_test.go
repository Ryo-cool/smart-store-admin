package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

type MockSaleRepository struct {
	mock.Mock
}

var _ repository.SaleRepository = (*MockSaleRepository)(nil)

func (m *MockSaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Sale, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error) {
	args := m.Called(ctx, timeOfDay)
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetTotalSalesAmount(ctx context.Context, start, end time.Time) (float64, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockSaleRepository) GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentalImpact), args.Error(1)
}

func (m *MockSaleRepository) GetSalesByCategory(ctx context.Context, start, end time.Time) (map[string]int, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(map[string]int), args.Error(1)
}

func TestCreateSale(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
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
				mockProductRepo.On("GetByID", ctx, productID).Return(product, nil)
				mockProductRepo.On("Update", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
				mockSaleRepo.On("Create", ctx, mock.AnythingOfType("*models.Sale")).Return(nil)
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
				mockProductRepo.On("GetByID", ctx, productID).Return(product, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.CreateSale(ctx, tt.sale)
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

func TestGetSalesAnalytics(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := &SaleService{repo: mockSaleRepo, productRepo: mockProductRepo}
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	sales := []*models.Sale{
		{
			TimeOfDay:     "morning",
			WeekDay:       "Monday",
			TotalAmount:   1000,
			TotalCO2Saved: 5.0,
		},
		{
			TimeOfDay:     "afternoon",
			WeekDay:       "Monday",
			TotalAmount:   2000,
			TotalCO2Saved: 10.0,
		},
	}

	mockSaleRepo.On("GetSalesByDateRange", ctx, start, end).Return(sales, nil)
	mockSaleRepo.On("GetTotalSalesAmount", ctx, start, end).Return(float64(3000), nil)

	analytics, err := service.GetSalesAnalytics(ctx, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, 2, analytics.TotalSales)
	assert.Equal(t, float64(3000), analytics.TotalAmount)
	assert.Equal(t, float64(15.0), analytics.TotalCO2Saved)
	assert.Equal(t, 1, analytics.TimeOfDaySales["morning"])
	assert.Equal(t, 1, analytics.TimeOfDaySales["afternoon"])
	assert.Equal(t, 2, analytics.WeekDaySales["Monday"])
}

func TestGetSalesByTimeOfDay(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
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
				mockSaleRepo.On("GetSalesByTimeOfDay", ctx, "morning").Return([]*models.Sale{}, nil)
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
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
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

	mockSaleRepo.On("GetSalesByDateRange", ctx, start, end).Return(sales, nil)

	expectedImpact := &models.EnvironmentalImpact{
		TotalCO2Saved:          15.0,
		AverageCO2SavedPerItem: 2.5,
		TotalEcoFriendlyItems:  6,
	}

	mockSaleRepo.On("GetEnvironmentalImpactAnalytics", ctx, start, end).Return(expectedImpact, nil)

	impact, err := service.GetEnvironmentalImpactAnalytics(ctx, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, impact)
	assert.Equal(t, float64(15.0), impact.TotalCO2Saved)
	assert.Equal(t, float64(2.5), impact.AverageCO2SavedPerItem)
	assert.Equal(t, 6, impact.TotalEcoFriendlyItems)
}

func TestGetSalesByCategory(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := &SaleService{repo: mockSaleRepo, productRepo: mockProductRepo}
	ctx := context.Background()

	productID := primitive.NewObjectID()
	product := &models.Product{
		ID:          productID,
		Name:        "エコ商品",
		Category:    "食品",
		CO2Emission: 2.0,
	}

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	sales := []*models.Sale{
		{
			Items: []models.SaleItem{
				{
					ProductID:   productID,
					Quantity:    2,
					PriceAtSale: 1000,
				},
			},
		},
	}

	expectedCategorySales := map[string]int{
		"食品": 2,
	}

	mockSaleRepo.On("GetSalesByDateRange", ctx, start, end).Return(sales, nil)
	mockProductRepo.On("GetByID", ctx, productID).Return(product, nil)
	mockSaleRepo.On("GetSalesByCategory", ctx, start, end).Return(expectedCategorySales, nil)

	categorySales, err := service.GetSalesByCategory(ctx, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, categorySales)
	assert.Equal(t, 2, categorySales["食品"])
}
