package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/service"
)

type SaleHandler struct {
	saleService service.SaleService
}

func NewSaleHandler(ss service.SaleService) *SaleHandler {
	return &SaleHandler{
		saleService: ss,
	}
}

// CreateSale は新しい売上を記録します
func (h *SaleHandler) CreateSale(c echo.Context) error {
	var sale models.Sale
	if err := json.NewDecoder(c.Request().Body).Decode(&sale); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なリクエストボディです",
		})
	}

	if err := h.saleService.Create(c.Request().Context(), &sale); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "売上の記録に失敗しました",
		})
	}

	return c.JSON(http.StatusCreated, sale)
}

// GetDailySales は日次の売上データを取得します
func (h *SaleHandler) GetDailySales(c echo.Context) error {
	dateStr := c.QueryParam("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な日付形式です",
		})
	}

	sales, err := h.saleService.GetDailySales(c.Request().Context(), date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "売上データの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, sales)
}

// GetSalesByDateRange は期間指定の売上データを取得します
func (h *SaleHandler) GetSalesByDateRange(c echo.Context) error {
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な開始日付です",
		})
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な終了日付です",
		})
	}

	sales, err := h.saleService.GetSalesByDateRange(c.Request().Context(), start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "売上データの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, sales)
}

// GetEnvironmentalImpact は環境影響の分析データを取得します
func (h *SaleHandler) GetEnvironmentalImpact(c echo.Context) error {
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な開始日付です",
		})
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な終了日付です",
		})
	}

	impact, err := h.saleService.GetEnvironmentalImpactAnalytics(c.Request().Context(), start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "環境影響データの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, impact)
}
