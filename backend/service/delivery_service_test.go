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

// MockDeliveryRepository はrepository.DeliveryRepositoryインターフェースのモック実装です
type MockDeliveryRepository struct {
	mock.Mock
}

// インターフェースが実装されていることを確認
var _ repository.DeliveryRepository = (*MockDeliveryRepository)(nil)

func (m *MockDeliveryRepository) Create(ctx context.Context, delivery *models.Delivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockDeliveryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) List(ctx context.Context, skip, limit int64) ([]*models.Delivery, error) {
	args := m.Called(ctx, skip, limit)
	return args.Get(0).([]*models.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockDeliveryRepository) UpdateLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
	args := m.Called(ctx, id, location)
	return args.Error(0)
}

func (m *MockDeliveryRepository) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
	args := m.Called(ctx, robotID)
	return args.Get(0).([]*models.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) GetDeliveries(query *models.DeliveryQuery) (*models.DeliveryResponse, error) {
	args := m.Called(query)
	return args.Get(0).(*models.DeliveryResponse), args.Error(1)
}

func (m *MockDeliveryRepository) GetDeliveryHistory(ctx context.Context, id primitive.ObjectID) (*models.DeliveryHistoryResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.DeliveryHistoryResponse), args.Error(1)
}

func (m *MockDeliveryRepository) Update(ctx context.Context, id primitive.ObjectID, delivery *models.Delivery) error {
	args := m.Called(ctx, id, delivery)
	return args.Error(0)
}

func TestCreateDelivery(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name     string
		delivery *models.Delivery
		mockFn   func()
		wantErr  bool
	}{
		{
			name: "正常な配送作成",
			delivery: &models.Delivery{
				DeliveryType:          "ロボット",
				Address:               "東京都渋谷区",
				EstimatedDeliveryTime: time.Now().Add(2 * time.Hour),
				Status:                models.StatusPreparing,
			},
			mockFn: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ロボットIDなしでエラー",
			delivery: &models.Delivery{
				Address: "東京都渋谷区",
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "開始位置なしでエラー",
			delivery: &models.Delivery{
				DeliveryType: "ロボット",
				Address:      "東京都渋谷区",
			},
			mockFn:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.CreateDelivery(ctx, tt.delivery)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, models.StatusPreparing, tt.delivery.Status)
			}
		})
	}
}

func TestUpdateDeliveryStatus(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	deliveryID := primitive.NewObjectID()

	tests := []struct {
		name          string
		currentStatus models.DeliveryStatus
		newStatus     string
		mockFn        func()
		wantErr       bool
	}{
		{
			name:          "Preparing から InProgress への遷移",
			currentStatus: models.StatusPreparing,
			newStatus:     string(models.StatusInProgress),
			mockFn: func() {
				mockRepo.On("GetByID", mock.Anything, deliveryID).Return(&models.Delivery{
					ID:     deliveryID.Hex(),
					Status: models.StatusPreparing,
				}, nil)
				mockRepo.On("UpdateStatus", mock.Anything, deliveryID, models.StatusInProgress).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "InProgress から Completed への遷移",
			currentStatus: models.StatusInProgress,
			newStatus:     string(models.StatusCompleted),
			mockFn: func() {
				mockRepo.On("GetByID", mock.Anything, deliveryID).Return(&models.Delivery{
					ID:     deliveryID.Hex(),
					Status: models.StatusInProgress,
				}, nil)
				mockRepo.On("UpdateStatus", mock.Anything, deliveryID, models.StatusCompleted).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "Completed から InProgress への無効な遷移",
			currentStatus: models.StatusCompleted,
			newStatus:     string(models.StatusInProgress),
			mockFn: func() {
				mockRepo.On("GetByID", mock.Anything, deliveryID).Return(&models.Delivery{
					ID:     deliveryID.Hex(),
					Status: models.StatusCompleted,
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockFn()
			err := service.UpdateDeliveryStatus(deliveryID.Hex(), tt.newStatus)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetActiveDeliveries(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()

	expectedDeliveries := []*models.Delivery{
		{
			ID:           primitive.NewObjectID().Hex(),
			DeliveryType: "ロボット",
			Status:       models.StatusPreparing,
		},
		{
			ID:           primitive.NewObjectID().Hex(),
			DeliveryType: "ロボット",
			Status:       models.StatusInProgress,
		},
	}

	mockRepo.On("GetActiveDeliveries", ctx).Return(expectedDeliveries, nil)

	deliveries, err := service.GetActiveDeliveries(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedDeliveries, deliveries)
	assert.Len(t, deliveries, 2)
	for _, d := range deliveries {
		assert.True(t, d.Status == models.StatusPreparing || d.Status == models.StatusInProgress)
	}
}

func TestGetDeliveryByID(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()
	deliveryID := primitive.NewObjectID()

	tests := []struct {
		name    string
		id      primitive.ObjectID
		mockFn  func()
		want    *models.Delivery
		wantErr bool
	}{
		{
			name: "正常な配送取得",
			id:   deliveryID,
			mockFn: func() {
				expectedDelivery := &models.Delivery{
					ID:           deliveryID.Hex(),
					DeliveryType: "ロボット",
					Status:       models.StatusInProgress,
				}
				mockRepo.On("GetByID", ctx, deliveryID).Return(expectedDelivery, nil)
			},
			want: &models.Delivery{
				ID:           deliveryID.Hex(),
				DeliveryType: "ロボット",
				Status:       models.StatusInProgress,
			},
			wantErr: false,
		},
		{
			name: "存在しない配送ID",
			id:   primitive.NewObjectID(),
			mockFn: func() {
				mockRepo.On("GetByID", ctx, mock.AnythingOfType("primitive.ObjectID")).Return(nil, errors.New("delivery not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetDeliveryByID(ctx, tt.id)
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

func TestUpdateDeliveryLocation(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()
	deliveryID := primitive.NewObjectID()

	validLocation := models.Location{
		Latitude:  35.6895,
		Longitude: 139.6917,
	}

	tests := []struct {
		name     string
		id       primitive.ObjectID
		location models.Location
		mockFn   func()
		wantErr  bool
	}{
		{
			name:     "正常な位置更新",
			id:       deliveryID,
			location: validLocation,
			mockFn: func() {
				mockRepo.On("UpdateLocation", ctx, deliveryID, validLocation).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "無効な位置",
			id:   deliveryID,
			location: models.Location{
				Latitude:  0,
				Longitude: 0,
			},
			mockFn: func() {
				mockRepo.On("UpdateLocation", ctx, deliveryID, models.Location{}).Return(errors.New("invalid location"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.UpdateDeliveryLocation(ctx, tt.id, tt.location)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDeliveriesByRobot(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()

	robotID := "ROBOT-001"
	expectedDeliveries := []*models.Delivery{
		{
			ID:           primitive.NewObjectID().Hex(),
			DeliveryType: "ロボット",
			Status:       models.StatusCompleted,
		},
		{
			ID:           primitive.NewObjectID().Hex(),
			DeliveryType: "ロボット",
			Status:       models.StatusInProgress,
		},
	}

	tests := []struct {
		name    string
		robotID string
		mockFn  func()
		want    []*models.Delivery
		wantErr bool
	}{
		{
			name:    "正常な配送履歴取得",
			robotID: robotID,
			mockFn: func() {
				mockRepo.On("GetDeliveriesByRobot", ctx, robotID).Return(expectedDeliveries, nil)
			},
			want:    expectedDeliveries,
			wantErr: false,
		},
		{
			name:    "存在しないロボットID",
			robotID: "NONEXISTENT",
			mockFn: func() {
				mockRepo.On("GetDeliveriesByRobot", ctx, "NONEXISTENT").Return([]*models.Delivery{}, nil)
			},
			want:    []*models.Delivery{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			got, err := service.GetDeliveriesByRobot(ctx, tt.robotID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
