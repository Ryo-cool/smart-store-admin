package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

// StoreOperationServiceInterface は店舗運営サービスのインターフェースを定義します
type StoreOperationServiceInterface interface {
	RecordStoreOperation(ctx context.Context, op *models.StoreOperation) error
	GetLatestOperation(ctx context.Context) (*models.StoreOperation, error)
	GetOperationsByTimeRange(ctx context.Context, start, end time.Time) ([]*models.StoreOperation, error)
	UpdateShelfStatus(ctx context.Context, opID primitive.ObjectID, status models.ShelfStatus) error
	UpdateCheckoutStatus(ctx context.Context, opID primitive.ObjectID, status models.CheckoutStatus) error
	GetEnergyUsageAnalytics(ctx context.Context, start, end time.Time) (*EnergyUsageAnalytics, error)
}

// EnergyUsageAnalytics はエネルギー使用量分析データを表す構造体です
type EnergyUsageAnalytics struct {
	AverageLightingUsage float64
	AverageACUsage       float64
	AverageRefrigUsage   float64
	TotalUsage           float64
	UsagePerHour         float64
}

// StoreOperationService は店舗運営サービスを表します
type StoreOperationService struct {
	repo repository.StoreOperationRepository
}

// NewStoreOperationService は新しい店舗運営サービスを作成します
func NewStoreOperationService(repo repository.StoreOperationRepository) *StoreOperationService {
	return &StoreOperationService{
		repo: repo,
	}
}

// RecordStoreOperation は店舗運営データを記録します
func (s *StoreOperationService) RecordStoreOperation(ctx context.Context, op *models.StoreOperation) error {
	if len(op.Shelves) == 0 {
		return errors.New("at least one shelf status is required")
	}
	if len(op.Checkouts) == 0 {
		return errors.New("at least one checkout status is required")
	}

	// 温度の妥当性チェック
	if !isValidTemperature(op.Temperature) {
		return errors.New("invalid temperature range")
	}

	// 湿度の妥当性チェック
	if !isValidHumidity(op.Humidity) {
		return errors.New("invalid humidity range")
	}

	return s.repo.Create(ctx, op)
}

// GetLatestOperation は最新の店舗運営データを取得します
func (s *StoreOperationService) GetLatestOperation(ctx context.Context) (*models.StoreOperation, error) {
	op, err := s.repo.GetLatest(ctx)
	if err != nil {
		return nil, err
	}
	if op == nil {
		return nil, errors.New("no data found")
	}
	return op, nil
}

// GetOperationsByTimeRange は指定期間の店舗運営データを取得します
func (s *StoreOperationService) GetOperationsByTimeRange(ctx context.Context, start, end time.Time) ([]*models.StoreOperation, error) {
	if end.Before(start) {
		return nil, errors.New("end time must be after start time")
	}
	return s.repo.GetByTimeRange(ctx, start, end)
}

// UpdateShelfStatus は特定の棚の状態を更新します
func (s *StoreOperationService) UpdateShelfStatus(ctx context.Context, opID primitive.ObjectID, status models.ShelfStatus) error {
	if status.ShelfID == "" {
		return errors.New("shelf ID is required")
	}
	if status.StockLevel < 0 {
		return errors.New("stock level must be non-negative")
	}
	if !isValidTemperature(status.Temperature) {
		return errors.New("invalid temperature range")
	}

	return s.repo.UpdateShelfStatus(ctx, opID, status)
}

// UpdateCheckoutStatus は特定のレジの状態を更新します
func (s *StoreOperationService) UpdateCheckoutStatus(ctx context.Context, opID primitive.ObjectID, status models.CheckoutStatus) error {
	if status.RegisterID == "" {
		return errors.New("register ID is required")
	}
	if status.QueueLength < 0 {
		return errors.New("queue length must be non-negative")
	}

	return s.repo.UpdateCheckoutStatus(ctx, opID, status)
}

// GetEnergyUsageAnalytics は指定期間のエネルギー使用量分析を取得します
func (s *StoreOperationService) GetEnergyUsageAnalytics(ctx context.Context, start, end time.Time) (*EnergyUsageAnalytics, error) {
	if end.Before(start) {
		return nil, errors.New("end time must be after start time")
	}

	usage, err := s.repo.GetAverageEnergyUsage(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// 時間あたりの使用量を計算
	duration := end.Sub(start).Hours()
	if duration == 0 {
		return nil, errors.New("time range must be greater than zero")
	}

	return &EnergyUsageAnalytics{
		AverageLightingUsage: usage["lighting"],
		AverageACUsage:       usage["ac"],
		AverageRefrigUsage:   usage["refrig"],
		TotalUsage:           usage["lighting"] + usage["ac"] + usage["refrig"],
		UsagePerHour:         (usage["lighting"] + usage["ac"] + usage["refrig"]) / duration,
	}, nil
}

// isValidTemperature は温度が妥当な範囲内かどうかを検証します
func isValidTemperature(temp float64) bool {
	return temp >= -30 && temp <= 50 // 冷凍庫から店内まで想定
}

// isValidHumidity は湿度が妥当な範囲内かどうかを検証します
func isValidHumidity(humidity float64) bool {
	return humidity >= 0 && humidity <= 100
}
