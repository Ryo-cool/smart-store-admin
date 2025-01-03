package handler

import (
	"context"
	"net/http"

	"smart-store-admin/backend/models"

	"github.com/labstack/echo/v4"
)

type DeliveryHandler struct {
	deliveryService DeliveryService
}

type DeliveryService interface {
	GetDeliveries(query *models.DeliveryQuery) (*models.DeliveryResponse, error)
	GetDelivery(id string) (*models.Delivery, error)
	UpdateDelivery(id string, delivery *models.Delivery) error
	UpdateDeliveryStatus(id string, status string) error
	GetDeliveryHistory(id string) (*models.DeliveryHistoryResponse, error)
	GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error)
	GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error)
}

func NewDeliveryHandler(ds DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{
		deliveryService: ds,
	}
}

// GetDeliveries handles GET /api/deliveries
func (h *DeliveryHandler) GetDeliveries(c echo.Context) error {
	var query models.DeliveryQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なクエリパラメータです",
		})
	}

	response, err := h.deliveryService.GetDeliveries(&query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "配送情報の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetDelivery handles GET /api/deliveries/:id
func (h *DeliveryHandler) GetDelivery(c echo.Context) error {
	id := c.Param("id")
	delivery, err := h.deliveryService.GetDelivery(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "配送情報の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, delivery)
}

// UpdateDelivery handles PATCH /api/deliveries/:id
func (h *DeliveryHandler) UpdateDelivery(c echo.Context) error {
	id := c.Param("id")
	var delivery models.Delivery
	if err := c.Bind(&delivery); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なリクエストボディです",
		})
	}

	if err := h.deliveryService.UpdateDelivery(id, &delivery); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "配送情報の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, delivery)
}

// UpdateDeliveryStatus handles PATCH /api/deliveries/:id/status
func (h *DeliveryHandler) UpdateDeliveryStatus(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なリクエストボディです",
		})
	}

	if !models.ValidateDeliveryStatus(req.Status) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な配送ステータスです",
		})
	}

	if err := h.deliveryService.UpdateDeliveryStatus(id, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "配送ステータスの更新に失敗しました",
		})
	}

	delivery, err := h.deliveryService.GetDelivery(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "配送情報の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, delivery)
}

// GetDeliveryHistory handles GET /api/deliveries/:id/history
func (h *DeliveryHandler) GetDeliveryHistory(c echo.Context) error {
	id := c.Param("id")
	history, err := h.deliveryService.GetDeliveryHistory(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "配送履歴の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, history)
}
