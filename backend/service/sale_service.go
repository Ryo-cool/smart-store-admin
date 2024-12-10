package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/repository"
)

// SaleService は売上サービスを表します
type SaleService struct {
	saleRepo    repository.SaleRepository
	productRepo repository.ProductRepository
}

// NewSaleService は新しい売上サービスを作成します
func NewSaleService(saleRepo repository.SaleRepository, productRepo repository.ProductRepository) *SaleService {
	return &SaleService{
		saleRepo:    saleRepo,
		productRepo: productRepo,
	}
}

// CreateSale は新しい売上を記録します
func (s *SaleService) CreateSale(ctx context.Context, sale *models.Sale) error {
	if len(sale.Items) == 0 {
		return errors.New("sale must have at least one item")
	}

	// 各商品の在庫チェックと更新
	for _, item := range sale.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return err
		}

		if product.Stock < item.Quantity {
			return errors.New("insufficient stock for product: " + product.Name)
		}

		// 在庫を減らす
		product.Stock -= item.Quantity
		if err := s.productRepo.Update(ctx, product); err != nil {
			return err
		}

		// 売上時の価格を記録
		item.PriceAtSale = product.Price
	}

	// 合計金額の計算
	var totalAmount float64
	for _, item := range sale.Items {
		totalAmount += item.PriceAtSale * float64(item.Quantity)
	}
	sale.TotalAmount = totalAmount

	// 環境影響の計算（CO2削減量）
	var totalCO2Saved float64
	for _, item := range sale.Items {
		product, _ := s.productRepo.GetByID(ctx, item.ProductID)
		if product != nil {
			totalCO2Saved += product.CO2Emission * float64(item.Quantity)
		}
	}
	sale.TotalCO2Saved = totalCO2Saved

	// 時間帯情報の設定
	now := time.Now()
	sale.TimeOfDay = getTimeOfDay(now)
	sale.WeekDay = now.Weekday().String()

	return s.saleRepo.Create(ctx, sale)
}

// GetSale は指定されたIDの売上を取得します
func (s *SaleService) GetSale(ctx context.Context, id primitive.ObjectID) (*models.Sale, error) {
	return s.saleRepo.GetByID(ctx, id)
}

// GetDailySales は指定された日の売上を取得します
func (s *SaleService) GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error) {
	return s.saleRepo.GetDailySales(ctx, date)
}

// GetSalesByTimeOfDay は時間帯別の売上を取得します
func (s *SaleService) GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error) {
	if !isValidTimeOfDay(timeOfDay) {
		return nil, errors.New("invalid time of day")
	}
	return s.saleRepo.GetSalesByTimeOfDay(ctx, timeOfDay)
}

// GetSalesAnalytics は期間を指定して売上分析データを取得します
func (s *SaleService) GetSalesAnalytics(ctx context.Context, start, end time.Time) (*SalesAnalytics, error) {
	sales, err := s.saleRepo.GetSalesByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	totalAmount, err := s.saleRepo.GetTotalSalesAmount(ctx, start, end)
	if err != nil {
		return nil, err
	}

	analytics := &SalesAnalytics{
		TotalSales:     len(sales),
		TotalAmount:    totalAmount,
		TimeOfDaySales: make(map[string]int),
		WeekDaySales:   make(map[string]int),
		TotalCO2Saved:  0,
	}

	for _, sale := range sales {
		analytics.TimeOfDaySales[sale.TimeOfDay]++
		analytics.WeekDaySales[sale.WeekDay]++
		analytics.TotalCO2Saved += sale.TotalCO2Saved
	}

	return analytics, nil
}

// SalesAnalytics は売上分析データを表す構造体です
type SalesAnalytics struct {
	TotalSales     int
	TotalAmount    float64
	TimeOfDaySales map[string]int
	WeekDaySales   map[string]int
	TotalCO2Saved  float64
}

// getTimeOfDay は現在時刻から時間帯を定します
func getTimeOfDay(t time.Time) string {
	hour := t.Hour()
	switch {
	case hour < 6:
		return "early_morning"
	case hour < 11:
		return "morning"
	case hour < 14:
		return "lunch"
	case hour < 17:
		return "afternoon"
	case hour < 21:
		return "evening"
	default:
		return "night"
	}
}

// isValidTimeOfDay は時間帯が有効かどうかを検証します
func isValidTimeOfDay(timeOfDay string) bool {
	validTimes := map[string]bool{
		"early_morning": true,
		"morning":       true,
		"lunch":         true,
		"afternoon":     true,
		"evening":       true,
		"night":         true,
	}
	return validTimes[timeOfDay]
}
