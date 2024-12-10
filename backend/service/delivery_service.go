package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

type DeliveryService struct {
	repo *repository.DeliveryRepository
}

func NewDeliveryService(repo *repository.DeliveryRepository) *DeliveryService {
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

// UpdateDeliveryStatus は配送のステータスを更新します
func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error {
	delivery, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// ステータス遷移の検証
	if !isValidStatusTransition(delivery.Status, status) {
		return errors.New("invalid status transition")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// UpdateLocation はロボット/ドローンの現在位置を更新します
func (s *DeliveryService) UpdateLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
	delivery, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if delivery.Status == models.StatusCompleted || delivery.Status == models.StatusFailed {
		return errors.New("cannot update location of completed or failed delivery")
	}

	// 進行中でない場合は、ステータスを進行中に更新
	if delivery.Status == models.StatusPending {
		if err := s.repo.UpdateStatus(ctx, id, models.StatusInProgress); err != nil {
			return err
		}
	}

	return s.repo.UpdateLocation(ctx, id, location)
}

// GetActiveDeliveries はアクティブな配送を取得します
func (s *DeliveryService) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	return s.repo.GetActiveDeliveries(ctx)
}

// GetDeliveriesByRobot は特定のロボット/ドローンの配送履歴を取得します
func (s *DeliveryService) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
	if robotID == "" {
		return nil, errors.New("robot ID is required")
	}
	return s.repo.GetDeliveriesByRobot(ctx, robotID)
}

// isValidStatusTransition はステータス遷移が有効かどうかを検証します
func isValidStatusTransition(current, next models.DeliveryStatus) bool {
	switch current {
	case models.StatusPending:
		return next == models.StatusInProgress || next == models.StatusFailed
	case models.StatusInProgress:
		return next == models.StatusCompleted || next == models.StatusFailed
	case models.StatusCompleted, models.StatusFailed:
		return false
	default:
		return false
	}
}
