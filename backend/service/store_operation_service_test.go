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

type MockStoreOperationRepository struct {
	mock.Mock
}

var _ repository.StoreOperationRepository = (*MockStoreOperationRepository)(nil)

func (m *MockStoreOperationRepository) Create(ctx context.Context, op *models.StoreOperation) error {
	args := m.Called(ctx, op)
	return args.Error(0)
}

func (m *MockStoreOperationRepository) GetLatest(ctx context.Context) (*models.StoreOperation, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StoreOperation), args.Error(1)
}

func (m *MockStoreOperationRepository) GetByTimeRange(ctx context.Context, start, end time.Time) ([]*models.StoreOperation, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]*models.StoreOperation), args.Error(1)
}

func (m *MockStoreOperationRepository) UpdateShelfStatus(ctx context.Context, opID primitive.ObjectID, status models.ShelfStatus) error {
	args := m.Called(ctx, opID, status)
	return args.Error(0)
}

func (m *MockStoreOperationRepository) UpdateCheckoutStatus(ctx context.Context, opID primitive.ObjectID, status models.CheckoutStatus) error {
	args := m.Called(ctx, opID, status)
	return args.Error(0)
}

func (m *MockStoreOperationRepository) GetAverageEnergyUsage(ctx context.Context, start, end time.Time) (map[string]float64, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func TestRecordStoreOperation(t *testing.T) {
	mockRepo := new(MockStoreOperationRepository)
	service := NewStoreOperationService(mockRepo)
	ctx := context.Background()

	validShelfStatus := models.ShelfStatus{
		ShelfID:     "SHELF-001",
		StockLevel:  80,
		Temperature: 20.0,
		LastChecked: time.Now(),
	}

	validCheckoutStatus := models.CheckoutStatus{
		RegisterID:    "REG-001",
		IsOperational: true,
		QueueLength:   0,
		LastChecked:   time.Now(),
	}

	tests := []struct {
		name    string
		op      *models.StoreOperation
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常な店舗運営データの記録",
			op: &models.StoreOperation{
				Temperature:  25.0,
				Humidity:     50.0,
				CrowdDensity: 0.3,
				Shelves:      []models.ShelfStatus{validShelfStatus},
				Checkouts:    []models.CheckoutStatus{validCheckoutStatus},
			},
			mockFn: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.StoreOperation")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "棚データなしでエラー",
			op: &models.StoreOperation{
				Temperature:  25.0,
				Humidity:     50.0,
				CrowdDensity: 0.3,
				Checkouts:    []models.CheckoutStatus{validCheckoutStatus},
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な温度でエラー",
			op: &models.StoreOperation{
				Temperature: 60.0, // 範囲外
				Humidity:    50.0,
				Shelves:     []models.ShelfStatus{validShelfStatus},
				Checkouts:   []models.CheckoutStatus{validCheckoutStatus},
			},
			mockFn:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.RecordStoreOperation(ctx, tt.op)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateShelfStatus(t *testing.T) {
	mockRepo := new(MockStoreOperationRepository)
	service := NewStoreOperationService(mockRepo)
	ctx := context.Background()
	opID := primitive.NewObjectID()

	tests := []struct {
		name    string
		status  models.ShelfStatus
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常な棚状態の更新",
			status: models.ShelfStatus{
				ShelfID:     "SHELF-001",
				StockLevel:  80,
				Temperature: 20.0,
				LastChecked: time.Now(),
			},
			mockFn: func() {
				mockRepo.On("UpdateShelfStatus", ctx, opID, mock.AnythingOfType("models.ShelfStatus")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "棚ID無しでエラー",
			status: models.ShelfStatus{
				StockLevel:  80,
				Temperature: 20.0,
				LastChecked: time.Now(),
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な温度でエラー",
			status: models.ShelfStatus{
				ShelfID:     "SHELF-001",
				StockLevel:  80,
				Temperature: -40.0, // 範囲外
				LastChecked: time.Now(),
			},
			mockFn:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.UpdateShelfStatus(ctx, opID, tt.status)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnergyUsageAnalytics(t *testing.T) {
	mockRepo := new(MockStoreOperationRepository)
	service := NewStoreOperationService(mockRepo)
	ctx := context.Background()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	usage := map[string]float64{
		"lighting": 100.0,
		"ac":       200.0,
		"refrig":   300.0,
	}

	mockRepo.On("GetAverageEnergyUsage", ctx, start, end).Return(usage, nil)

	analytics, err := service.GetEnergyUsageAnalytics(ctx, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, 100.0, analytics.AverageLightingUsage)
	assert.Equal(t, 200.0, analytics.AverageACUsage)
	assert.Equal(t, 300.0, analytics.AverageRefrigUsage)
	assert.Equal(t, 600.0, analytics.TotalUsage)
	assert.Equal(t, 25.0, analytics.UsagePerHour) // 600 / 24 hours
}
