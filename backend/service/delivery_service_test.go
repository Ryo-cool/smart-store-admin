package service

import (
	"context"
	"testing"

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

func TestCreateDelivery(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()

	validLocation := models.Location{
		Latitude:  35.6895,
		Longitude: 139.6917,
	}

	tests := []struct {
		name     string
		delivery *models.Delivery
		mockFn   func()
		wantErr  bool
	}{
		{
			name: "正常な配送作成",
			delivery: &models.Delivery{
				RobotID:       "ROBOT-001",
				StartLocation: validLocation,
				EndLocation:   validLocation,
				BatteryLevel:  100.0,
				EnergyUsage:   0.0,
			},
			mockFn: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ロボットIDなしでエラー",
			delivery: &models.Delivery{
				StartLocation: validLocation,
				EndLocation:   validLocation,
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "開始位置なしでエラー",
			delivery: &models.Delivery{
				RobotID:     "ROBOT-001",
				EndLocation: validLocation,
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
				assert.Equal(t, models.StatusPending, tt.delivery.Status)
				assert.Equal(t, tt.delivery.StartLocation, tt.delivery.CurrentLocation)
			}
		})
	}
}

func TestUpdateDeliveryStatus(t *testing.T) {
	mockRepo := new(MockDeliveryRepository)
	service := NewDeliveryService(mockRepo)
	ctx := context.Background()
	deliveryID := primitive.NewObjectID()

	tests := []struct {
		name          string
		currentStatus models.DeliveryStatus
		newStatus     models.DeliveryStatus
		mockFn        func()
		wantErr       bool
	}{
		{
			name:          "Pending から InProgress への遷移",
			currentStatus: models.StatusPending,
			newStatus:     models.StatusInProgress,
			mockFn: func() {
				mockRepo.On("GetByID", ctx, deliveryID).Return(&models.Delivery{
					ID:     deliveryID,
					Status: models.StatusPending,
				}, nil)
				mockRepo.On("UpdateStatus", ctx, deliveryID, models.StatusInProgress).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "InProgress から Completed への遷移",
			currentStatus: models.StatusInProgress,
			newStatus:     models.StatusCompleted,
			mockFn: func() {
				mockRepo.On("GetByID", ctx, deliveryID).Return(&models.Delivery{
					ID:     deliveryID,
					Status: models.StatusInProgress,
				}, nil)
				mockRepo.On("UpdateStatus", ctx, deliveryID, models.StatusCompleted).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "Completed から InProgress への無効な遷移",
			currentStatus: models.StatusCompleted,
			newStatus:     models.StatusInProgress,
			mockFn: func() {
				mockRepo.On("GetByID", ctx, deliveryID).Return(&models.Delivery{
					ID:     deliveryID,
					Status: models.StatusCompleted,
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := service.UpdateDeliveryStatus(ctx, deliveryID, tt.newStatus)
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
			ID:      primitive.NewObjectID(),
			RobotID: "ROBOT-001",
			Status:  models.StatusPending,
		},
		{
			ID:      primitive.NewObjectID(),
			RobotID: "ROBOT-002",
			Status:  models.StatusInProgress,
		},
	}

	mockRepo.On("GetActiveDeliveries", ctx).Return(expectedDeliveries, nil)

	deliveries, err := service.GetActiveDeliveries(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedDeliveries, deliveries)
	assert.Len(t, deliveries, 2)
	for _, d := range deliveries {
		assert.True(t, d.Status == models.StatusPending || d.Status == models.StatusInProgress)
	}
}
