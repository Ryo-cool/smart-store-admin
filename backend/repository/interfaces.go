package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
)

// ProductRepository は商品リポジトリのインターフェースを定義します
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error)
	List(ctx context.Context, skip, limit int64) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByCategory(ctx context.Context, category string) ([]*models.Product, error)
	GetLowStock(ctx context.Context) ([]*models.Product, error)
}

// DeliveryRepository は配送リポジトリのインターフェースを定義します
type DeliveryRepository interface {
	Create(ctx context.Context, delivery *models.Delivery) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error)
	List(ctx context.Context, skip, limit int64) ([]*models.Delivery, error)
	Update(ctx context.Context, id primitive.ObjectID, delivery *models.Delivery) error
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error
	UpdateLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error
	GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error)
	GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error)
	GetDeliveries(query *models.DeliveryQuery) (*models.DeliveryResponse, error)
	GetDeliveryHistory(ctx context.Context, id primitive.ObjectID) (*models.DeliveryHistoryResponse, error)
}

// SaleRepository は売上リポジトリのインターフェースを定義します
type SaleRepository interface {
	Create(ctx context.Context, sale *models.Sale) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Sale, error)
	GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error)
	GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error)
	GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error)
	GetTotalSalesAmount(ctx context.Context, start, end time.Time) (float64, error)
	GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error)
	GetSalesByCategory(ctx context.Context, start, end time.Time) (map[string]int, error)
}

// StoreOperationRepository は店舗運営リポジトリのインターフェースを定義します
type StoreOperationRepository interface {
	Create(ctx context.Context, op *models.StoreOperation) error
	GetLatest(ctx context.Context) (*models.StoreOperation, error)
	GetByTimeRange(ctx context.Context, start, end time.Time) ([]*models.StoreOperation, error)
	UpdateShelfStatus(ctx context.Context, opID primitive.ObjectID, status models.ShelfStatus) error
	UpdateCheckoutStatus(ctx context.Context, opID primitive.ObjectID, status models.CheckoutStatus) error
	GetAverageEnergyUsage(ctx context.Context, start, end time.Time) (map[string]float64, error)
}
