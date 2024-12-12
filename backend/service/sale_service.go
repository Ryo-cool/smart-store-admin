package service

import (
	"context"
	"time"

	"smart-store-admin/backend/models"
)

type SaleService interface {
	Create(ctx context.Context, sale *models.Sale) error
	GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error)
	GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error)
	GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error)
}
