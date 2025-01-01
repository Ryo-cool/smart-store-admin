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
	if delivery.RobotID == "" {
		return errors.New("robot ID is required")
	}
	if delivery.StartLocation == (models.Location{}) {
		return errors.New("start location is required")
	}
	if delivery.EndLocation == (models.Location{}) {
		return errors.New("end location is required")
	}

	// 初期状態の設定
	delivery.Status = models.StatusPending
	delivery.StartedAt = time.Now()
	delivery.CurrentLocation = delivery.StartLocation

	return s.repo.Create(ctx, delivery)
}

// GetDelivery は指定されたIDの配送を取得します
func (s *DeliveryService) GetDelivery(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error) {
	return s.repo.GetByID(ctx, id)
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
func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error {
	delivery, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// ステータスの遷移チェック
	if !isValidStatusTransition(delivery.Status, status) {
		return errors.New("invalid status transition")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// UpdateDeliveryLocation は配送の現在位置を更新します
func (s *DeliveryService) UpdateDeliveryLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
	return s.repo.UpdateLocation(ctx, id, location)
}

// GetActiveDeliveries はアクティブな配送一覧を取得します
func (s *DeliveryService) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	return s.repo.GetActiveDeliveries(ctx)
}

// GetDeliveriesByRobot は指定されたロボットの配送一覧を取得します
func (s *DeliveryService) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
	return s.repo.GetDeliveriesByRobot(ctx, robotID)
}

// isValidStatusTransition はステータスの遷移が有効かどうかをチェックします
func isValidStatusTransition(current, new models.DeliveryStatus) bool {
	switch current {
	case models.StatusPending:
		return new == models.StatusInProgress
	case models.StatusInProgress:
		return new == models.StatusCompleted || new == models.StatusCancelled
	case models.StatusCompleted, models.StatusCancelled:
		return false
	default:
		return false
	}
}
