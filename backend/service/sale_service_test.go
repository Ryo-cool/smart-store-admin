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

func TestCreateSale(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	service := NewSaleService(mockSaleRepo)
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
					assert.Equal(t, float64(2000), tt.sale.TotalAmount) // 1000 * 2
					assert.Equal(t, float64(10), tt.sale.TotalCO2Saved) // 5.0 * 2
				}
			}
		})
	}
}

func TestGetSalesAnalytics(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	service := NewSaleService(mockSaleRepo)
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
	service := NewSaleService(mockSaleRepo)
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
