package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	"smart-store-admin/backend/service"
)

type ProductHandler struct {
	productService service.ProductServiceInterface
}

func NewProductHandler(ps service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{
		productService: ps,
	}
}

// CreateProduct は新しい商品を作成します
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product models.Product
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なリクエストボディです",
		})
	}

	if err := h.productService.CreateProduct(c.Request().Context(), &product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品の作成に失敗しました",
		})
	}

	return c.JSON(http.StatusCreated, product)
}

// GetProduct は指定されたIDの商品を取得します
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なIDです",
		})
	}

	product, err := h.productService.GetProductByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品の取得に失敗しました",
		})
	}

	if product == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "商品が見つかりません",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// ListProducts は商品のリストを取得します
func (h *ProductHandler) ListProducts(c echo.Context) error {
	// クエリパラメータからページネーション情報を取得
	page := 1
	limit := 10
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	skip := int64((page - 1) * limit)
	products, err := h.productService.List(c.Request().Context(), skip, int64(limit))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品リストの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, products)
}

// UpdateProduct は商品情報を更新します
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な商品IDです",
		})
	}

	var product models.Product
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なリクエストボディです",
		})
	}

	product.ID = id
	if err := h.productService.Update(c.Request().Context(), &product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct は商品を削除します
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な商品IDです",
		})
	}

	if err := h.productService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品の削除に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "商品を削除しました",
	})
}
