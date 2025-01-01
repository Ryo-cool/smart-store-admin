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

// MockStoreOperationRepository はrepository.StoreOperationRepositoryインターフェースのモック実装です
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
		Temperature: 20.0,
		StockLevel:  80,
	}

	validCheckoutStatus := models.CheckoutStatus{
		RegisterID:    "REG-001",
		QueueLength:   3,
		IsOperational: true,
	}

	tests := []struct {
		name    string
		op      *models.StoreOperation
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常な店舗運営データ記録",
			op: &models.StoreOperation{
				Temperature: 25.0,
				Humidity:    50.0,
				Shelves:     []models.ShelfStatus{validShelfStatus},
				Checkouts:   []models.CheckoutStatus{validCheckoutStatus},
			},
			mockFn: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.StoreOperation")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "棚データなしでエラー",
			op: &models.StoreOperation{
				Temperature: 25.0,
				Humidity:    50.0,
				Checkouts:   []models.CheckoutStatus{validCheckoutStatus},
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "レジデータなしでエラー",
			op: &models.StoreOperation{
				Temperature: 25.0,
				Humidity:    50.0,
				Shelves:     []models.ShelfStatus{validShelfStatus},
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な温度でエラー",
			op: &models.StoreOperation{
				Temperature: -40.0,
				Humidity:    50.0,
				Shelves:     []models.ShelfStatus{validShelfStatus},
				Checkouts:   []models.CheckoutStatus{validCheckoutStatus},
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な湿度でエラー",
			op: &models.StoreOperation{
				Temperature: 25.0,
				Humidity:    150.0,
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

func TestGetLatestOperation(t *testing.T) {
	mockRepo := new(MockStoreOperationRepository)
	service := NewStoreOperationService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name    string
		mockFn  func()
		want    *models.StoreOperation
		wantErr bool
	}{
		{
			name: "最新の店舗運営データ取得",
			mockFn: func() {
				expectedOp := &models.StoreOperation{
					ID:          primitive.NewObjectID(),
					Temperature: 25.0,
					Humidity:    50.0,
				}
				mockRepo.On("GetLatest", ctx).Return(expectedOp, nil)
			},
			want: &models.StoreOperation{
				Temperature: 25.0,
				Humidity:    50.0,
			},
			wantErr: false,
		},
		{
			name: "データが存在しない場合",
			mockFn: func() {
				mockRepo.On("GetLatest", ctx).Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockFn()
			got, err := service.GetLatestOperation(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Temperature, got.Temperature)
				assert.Equal(t, tt.want.Humidity, got.Humidity)
			}
		})
	}
}

func TestGetOperationsByTimeRange(t *testing.T) {
	mockRepo := new(MockStoreOperationRepository)
	service := NewStoreOperationService(mockRepo)
	ctx := context.Background()

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		mockFn  func()
		want    []*models.StoreOperation
		wantErr bool
	}{
		{
			name:  "正常な期間での取得",
			start: start,
			end:   end,
			mockFn: func() {
				expectedOps := []*models.StoreOperation{
					{
						ID:          primitive.NewObjectID(),
						Temperature: 25.0,
						Humidity:    50.0,
					},
					{
						ID:          primitive.NewObjectID(),
						Temperature: 26.0,
						Humidity:    51.0,
					},
				}
				mockRepo.On("GetByTimeRange", ctx, start, end).Return(expectedOps, nil)
			},
			want: []*models.StoreOperation{
				{
					Temperature: 25.0,
					Humidity:    50.0,
				},
				{
					Temperature: 26.0,
					Humidity:    51.0,
				},
			},
			wantErr: false,
		},
		{
			name:    "無効な期間でエラー",
			start:   end,
			end:     start,
			mockFn:  func() {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetOperationsByTimeRange(ctx, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, len(tt.want))
				for i, op := range got {
					assert.Equal(t, tt.want[i].Temperature, op.Temperature)
					assert.Equal(t, tt.want[i].Humidity, op.Humidity)
				}
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
			name: "正常な棚状態更新",
			status: models.ShelfStatus{
				ShelfID:     "SHELF-001",
				Temperature: 20.0,
				StockLevel:  80,
			},
			mockFn: func() {
				mockRepo.On("UpdateShelfStatus", ctx, opID, mock.AnythingOfType("models.ShelfStatus")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "棚IDなしでエラー",
			status: models.ShelfStatus{
				Temperature: 20.0,
				StockLevel:  80,
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な在庫レベルでエラー",
			status: models.ShelfStatus{
				ShelfID:     "SHELF-001",
				Temperature: 20.0,
				StockLevel:  -1,
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な温度でエラー",
			status: models.ShelfStatus{
				ShelfID:     "SHELF-001",
				Temperature: -40.0,
				StockLevel:  80,
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

func TestUpdateCheckoutStatus(t *testing.T) {
	mockRepo := new(MockStoreOperationRepository)
	service := NewStoreOperationService(mockRepo)
	ctx := context.Background()
	opID := primitive.NewObjectID()

	tests := []struct {
		name    string
		status  models.CheckoutStatus
		mockFn  func()
		wantErr bool
	}{
		{
			name: "正常なレジ状態更新",
			status: models.CheckoutStatus{
				RegisterID:    "REG-001",
				QueueLength:   3,
				IsOperational: true,
			},
			mockFn: func() {
				mockRepo.On("UpdateCheckoutStatus", ctx, opID, mock.AnythingOfType("models.CheckoutStatus")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "レジIDなしでエラー",
			status: models.CheckoutStatus{
				QueueLength:   3,
				IsOperational: true,
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "無効な待ち行列長でエラー",
			status: models.CheckoutStatus{
				RegisterID:    "REG-001",
				QueueLength:   -1,
				IsOperational: true,
			},
			mockFn:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.UpdateCheckoutStatus(ctx, opID, tt.status)
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

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		mockFn  func()
		want    *EnergyUsageAnalytics
		wantErr bool
	}{
		{
			name:  "正常なエネルギー使用量分析",
			start: start,
			end:   end,
			mockFn: func() {
				usage := map[string]float64{
					"lighting": 100.0,
					"ac":       200.0,
					"refrig":   150.0,
				}
				mockRepo.On("GetAverageEnergyUsage", ctx, start, end).Return(usage, nil)
			},
			want: &EnergyUsageAnalytics{
				AverageLightingUsage: 100.0,
				AverageACUsage:       200.0,
				AverageRefrigUsage:   150.0,
				TotalUsage:           450.0,
				UsagePerHour:         18.75, // 450.0 / 24
			},
			wantErr: false,
		},
		{
			name:    "無効な期間でエラー",
			start:   end,
			end:     start,
			mockFn:  func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "0時間の期間でエラー",
			start: now,
			end:   now,
			mockFn: func() {
				usage := map[string]float64{
					"lighting": 100.0,
					"ac":       200.0,
					"refrig":   150.0,
				}
				mockRepo.On("GetAverageEnergyUsage", ctx, now, now).Return(usage, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetEnergyUsageAnalytics(ctx, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.AverageLightingUsage, got.AverageLightingUsage)
				assert.Equal(t, tt.want.AverageACUsage, got.AverageACUsage)
				assert.Equal(t, tt.want.AverageRefrigUsage, got.AverageRefrigUsage)
				assert.Equal(t, tt.want.TotalUsage, got.TotalUsage)
				assert.Equal(t, tt.want.UsagePerHour, got.UsagePerHour)
			}
		})
	}
}
