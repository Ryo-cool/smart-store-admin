package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

// MockSaleRepository はrepository.SaleRepositoryインターフェースのモック実装です
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

func (m *MockSaleRepository) GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentalImpact), args.Error(1)
}

func (m *MockSaleRepository) GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error) {
	args := m.Called(ctx, timeOfDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetSalesByCategory(ctx context.Context, start, end time.Time) (map[string]int, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockSaleRepository) GetTotalSalesAmount(ctx context.Context, start, end time.Time) (float64, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(float64), args.Error(1)
}

func TestCreate(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo)
	ctx := context.Background()

	productID := primitive.NewObjectID()
	validSale := &models.Sale{
		Items: []models.SaleItem{
			{
				ProductID:   productID,
				Quantity:    2,
				PriceAtSale: 1000,
			},
		},
		TotalAmount:   2000,
		TotalCO2Saved: 10.0,
		TimeOfDay:     "morning",
		PaymentMethod: "credit_card",
	}

	tests := []struct {
		name    string
		sale    *models.Sale
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常な売上記録",
			sale: validSale,
			mockFn: func() {
				mockProductRepo.On("GetByID", ctx, productID).Return(&models.Product{
					ID:    productID,
					Name:  "Test Product",
					Price: 1000,
				}, nil)
				mockSaleRepo.On("Create", ctx, mock.AnythingOfType("*models.Sale")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "商品なしでエラー",
			sale: &models.Sale{
				Items:         []models.SaleItem{},
				TotalAmount:   0,
				TotalCO2Saved: 0,
				TimeOfDay:     "morning",
				PaymentMethod: "cash",
			},
			mockFn:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaleRepo.ExpectedCalls = nil
			mockProductRepo.ExpectedCalls = nil
			tt.mockFn()
			err := service.Create(ctx, tt.sale)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDailySales(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo)
	ctx := context.Background()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedSales := []*models.Sale{
		{
			ID:            primitive.NewObjectID(),
			TotalAmount:   2000,
			TimeOfDay:     "morning",
			PaymentMethod: "credit_card",
		},
		{
			ID:            primitive.NewObjectID(),
			TotalAmount:   3000,
			TimeOfDay:     "afternoon",
			PaymentMethod: "cash",
		},
	}

	tests := []struct {
		name    string
		date    time.Time
		mockFn  func()
		want    []*models.Sale
		wantErr bool
	}{
		{
			name: "正常な日次売上取得",
			date: date,
			mockFn: func() {
				mockSaleRepo.On("GetDailySales", ctx, date).Return(expectedSales, nil)
			},
			want:    expectedSales,
			wantErr: false,
		},
		{
			name: "データなしの場合",
			date: date.AddDate(0, 0, 1),
			mockFn: func() {
				mockSaleRepo.On("GetDailySales", ctx, date.AddDate(0, 0, 1)).Return([]*models.Sale{}, nil)
			},
			want:    []*models.Sale{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetDailySales(ctx, tt.date)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetSalesByDateRange(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo)
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	expectedSales := []*models.Sale{
		{
			ID:            primitive.NewObjectID(),
			TotalAmount:   2000,
			TimeOfDay:     "morning",
			PaymentMethod: "credit_card",
		},
		{
			ID:            primitive.NewObjectID(),
			TotalAmount:   3000,
			TimeOfDay:     "afternoon",
			PaymentMethod: "cash",
		},
	}

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		mockFn  func()
		want    []*models.Sale
		wantErr bool
	}{
		{
			name:  "正常な期間売上取得",
			start: start,
			end:   end,
			mockFn: func() {
				mockSaleRepo.On("GetSalesByDateRange", ctx, start, end).Return(expectedSales, nil)
			},
			want:    expectedSales,
			wantErr: false,
		},
		{
			name:  "データなしの期間",
			start: start.AddDate(0, 1, 0),
			end:   end.AddDate(0, 1, 0),
			mockFn: func() {
				mockSaleRepo.On("GetSalesByDateRange", ctx, start.AddDate(0, 1, 0), end.AddDate(0, 1, 0)).Return([]*models.Sale{}, nil)
			},
			want:    []*models.Sale{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetSalesByDateRange(ctx, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetEnvironmentalImpactAnalytics(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo)
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	expectedImpact := &models.EnvironmentalImpact{
		TotalCO2Saved:          100.0,
		AverageCO2SavedPerItem: 10.0,
		TotalEcoFriendlyItems:  10,
	}

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		mockFn  func()
		want    *models.EnvironmentalImpact
		wantErr bool
	}{
		{
			name:  "正常な環境影響分析取得",
			start: start,
			end:   end,
			mockFn: func() {
				mockSaleRepo.On("GetEnvironmentalImpactAnalytics", ctx, start, end).Return(expectedImpact, nil)
			},
			want:    expectedImpact,
			wantErr: false,
		},
		{
			name:  "データなしの期間",
			start: start.AddDate(0, 1, 0),
			end:   end.AddDate(0, 1, 0),
			mockFn: func() {
				mockSaleRepo.On("GetEnvironmentalImpactAnalytics", ctx, start.AddDate(0, 1, 0), end.AddDate(0, 1, 0)).Return(nil, errors.New("no data found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetEnvironmentalImpactAnalytics(ctx, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetSalesByTimeOfDay(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo)
	ctx := context.Background()

	expectedSales := []*models.Sale{
		{
			ID:            primitive.NewObjectID(),
			TotalAmount:   2000,
			TimeOfDay:     "morning",
			PaymentMethod: "credit_card",
		},
		{
			ID:            primitive.NewObjectID(),
			TotalAmount:   3000,
			TimeOfDay:     "morning",
			PaymentMethod: "cash",
		},
	}

	tests := []struct {
		name      string
		timeOfDay string
		mockFn    func()
		want      []*models.Sale
		wantErr   bool
	}{
		{
			name:      "正常な時間帯売上取得",
			timeOfDay: "morning",
			mockFn: func() {
				mockSaleRepo.On("GetSalesByTimeOfDay", ctx, "morning").Return(expectedSales, nil)
			},
			want:    expectedSales,
			wantErr: false,
		},
		{
			name:      "無効な時間帯",
			timeOfDay: "invalid_time",
			mockFn: func() {
				mockSaleRepo.On("GetSalesByTimeOfDay", ctx, "invalid_time").Return(nil, errors.New("invalid time of day"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetSalesByTimeOfDay(ctx, tt.timeOfDay)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetSalesByCategory(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo)
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	expectedCategories := map[string]int{
		"食品": 10,
		"飲料": 5,
	}

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		mockFn  func()
		want    map[string]int
		wantErr bool
	}{
		{
			name:  "正常なカテゴリー別売上取得",
			start: start,
			end:   end,
			mockFn: func() {
				mockSaleRepo.On("GetSalesByCategory", ctx, start, end).Return(expectedCategories, nil)
			},
			want:    expectedCategories,
			wantErr: false,
		},
		{
			name:  "データなしの期間",
			start: start.AddDate(0, 1, 0),
			end:   end.AddDate(0, 1, 0),
			mockFn: func() {
				mockSaleRepo.On("GetSalesByCategory", ctx, start.AddDate(0, 1, 0), end.AddDate(0, 1, 0)).Return(map[string]int{}, nil)
			},
			want:    map[string]int{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetSalesByCategory(ctx, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
