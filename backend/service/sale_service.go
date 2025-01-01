package service

import (
	"context"
	"errors"
	"time"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

type SaleServiceInterface interface {
	Create(ctx context.Context, sale *models.Sale) error
	GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error)
	GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error)
	GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error)
	GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error)
	GetSalesByCategory(ctx context.Context, start, end time.Time) (map[string]int, error)
}

type SaleService struct {
	repo        repository.SaleRepository
	productRepo repository.ProductRepository
}

// オプション: コンストラクタ
func NewSaleService(repo repository.SaleRepository, productRepo repository.ProductRepository) *SaleService {
	return &SaleService{
		repo:        repo,
		productRepo: productRepo,
	}
}

// Create は新しい売上を記録します
func (ss *SaleService) Create(ctx context.Context, sale *models.Sale) error {
	if sale == nil || len(sale.Items) == 0 {
		return errors.New("商品が指定されていません")
	}

	// 商品の存在チェック
	for _, item := range sale.Items {
		if item.ProductID.IsZero() {
			return errors.New("無効な商品IDです")
		}
		p, err := ss.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return err
		}
		if p == nil {
			return errors.New("指定された商品が存在しません")
		}
	}

	return ss.repo.Create(ctx, sale)
}

func (ss *SaleService) GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error) {
	return ss.repo.GetDailySales(ctx, date)
}

func (ss *SaleService) GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error) {
	return ss.repo.GetSalesByDateRange(ctx, start, end)
}

func (ss *SaleService) GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error) {
	return ss.repo.GetEnvironmentalImpactAnalytics(ctx, start, end)
}

func (ss *SaleService) GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error) {
	return ss.repo.GetSalesByTimeOfDay(ctx, timeOfDay)
}

func (ss *SaleService) GetSalesByCategory(ctx context.Context, start, end time.Time) (map[string]int, error) {
	return ss.repo.GetSalesByCategory(ctx, start, end)
}
