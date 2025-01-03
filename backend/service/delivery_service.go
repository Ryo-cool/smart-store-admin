package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

// DeliveryServiceInterface は配送サービスのインターフェースを定義します
type DeliveryServiceInterface interface {
	CreateDelivery(ctx context.Context, delivery *models.Delivery) error
	GetDeliveryByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error)
	UpdateDeliveryStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error
	UpdateDeliveryLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error
	GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error)
	GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error)
	GetDeliveries(query *models.DeliveryQuery) (*models.DeliveryResponse, error)
	GetDeliveryHistory(id string) (*models.DeliveryHistoryResponse, error)
	UpdateDelivery(id string, delivery *models.Delivery) error
}

// DeliveryService は配送サービスを表します
type DeliveryService struct {
	repo repository.DeliveryRepository
}

// NewDeliveryService は新しい配送サービスを作成します
func NewDeliveryService(repo repository.DeliveryRepository) *DeliveryService {
	return &DeliveryService{
		repo: repo,
	}
}

// CreateDelivery は新しい配送を作成します
func (s *DeliveryService) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	if delivery.DeliveryType == "" {
		return errors.New("delivery type is required")
	}
	if delivery.Address == "" {
		return errors.New("address is required")
	}

	// 初期状態の設定
	delivery.Status = models.StatusPreparing
	delivery.CreatedAt = time.Now()
	delivery.UpdatedAt = time.Now()

	return s.repo.Create(ctx, delivery)
}

// GetDelivery は指定されたIDの配送を取得します
func (s *DeliveryService) GetDelivery(id string) (*models.Delivery, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return s.repo.GetByID(ctx, objectID)
}

// ListDeliveries は配送のリストを取得します
func (s *DeliveryService) ListDeliveries(ctx context.Context, page, pageSize int64) ([]*models.Delivery, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	return s.repo.List(ctx, skip, pageSize)
}

// GetDeliveryByID は指定されたIDの配送を取得します
func (s *DeliveryService) GetDeliveryByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateDeliveryStatus は配送のステータスを更新します
func (s *DeliveryService) UpdateDeliveryStatus(id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	deliveryStatus := models.DeliveryStatus(status)
	if !models.ValidateDeliveryStatus(status) {
		return errors.New("invalid status")
	}

	ctx := context.Background()
	delivery, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}
	if delivery == nil {
		return errors.New("delivery not found")
	}

	if !isValidStatusTransition(delivery.Status, deliveryStatus) {
		return errors.New("invalid status transition")
	}

	return s.repo.UpdateStatus(ctx, objectID, deliveryStatus)
}

// UpdateDeliveryLocation は配送の現在位置を更新します
func (s *DeliveryService) UpdateDeliveryLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
	if location == (models.Location{}) {
		return errors.New("invalid location")
	}
	return s.repo.UpdateLocation(ctx, id, location)
}

// GetActiveDeliveries はアクティブな配送一覧を取得します
func (s *DeliveryService) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	return s.repo.GetActiveDeliveries(ctx)
}

// GetDeliveriesByRobot は指定されたロボットの配送一覧を取得します
func (s *DeliveryService) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
	if robotID == "" {
		return nil, errors.New("robot ID is required")
	}
	return s.repo.GetDeliveriesByRobot(ctx, robotID)
}

// GetDeliveries retrieves deliveries based on the query parameters
func (s *DeliveryService) GetDeliveries(query *models.DeliveryQuery) (*models.DeliveryResponse, error) {
	return s.repo.GetDeliveries(query)
}

// GetDeliveryHistory は配送履歴を取得します
func (s *DeliveryService) GetDeliveryHistory(id string) (*models.DeliveryHistoryResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return s.repo.GetDeliveryHistory(ctx, objectID)
}

// UpdateDelivery は配送情報を更新します
func (s *DeliveryService) UpdateDelivery(id string, delivery *models.Delivery) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return s.repo.Update(ctx, objectID, delivery)
}

// isValidStatusTransition はステータスの遷移が有効かどうかをチェックします
func isValidStatusTransition(current, new models.DeliveryStatus) bool {
	validTransitions := map[models.DeliveryStatus][]models.DeliveryStatus{
		models.StatusPreparing: {
			models.StatusInProgress,
			models.StatusFailed,
		},
		models.StatusInProgress: {
			models.StatusCompleted,
			models.StatusFailed,
		},
		models.StatusCompleted: {}, // 完了状態からの遷移は不可
		models.StatusFailed:    {}, // 失敗状態からの遷移は不可
	}

	allowedStatuses, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == new {
			return true
		}
	}
	return false
}
