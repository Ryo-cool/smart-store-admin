package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"smart-store-admin/backend/handler"
)

func NewRouter(
	productHandler *handler.ProductHandler,
	saleHandler *handler.SaleHandler,
	deliveryHandler *handler.DeliveryHandler,
) *echo.Echo {
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API グループ
	api := e.Group("/api")

	// 商品関連のエンドポイント
	products := api.Group("/products")
	products.POST("", productHandler.CreateProduct)
	products.GET("", productHandler.ListProducts)
	products.GET("/:id", productHandler.GetProduct)
	products.PUT("/:id", productHandler.UpdateProduct)
	products.DELETE("/:id", productHandler.DeleteProduct)

	// 売上関連のエンドポイント
	sales := api.Group("/sales")
	sales.POST("", saleHandler.CreateSale)
	sales.GET("/daily", saleHandler.GetDailySales)
	sales.GET("/range", saleHandler.GetSalesByDateRange)
	sales.GET("/environmental-impact", saleHandler.GetEnvironmentalImpact)

	// 配送関連のエンドポイント
	deliveries := api.Group("/deliveries")
	deliveries.GET("", deliveryHandler.GetDeliveries)
	deliveries.GET("/:id", deliveryHandler.GetDelivery)
	deliveries.PATCH("/:id", deliveryHandler.UpdateDelivery)
	deliveries.PATCH("/:id/status", deliveryHandler.UpdateDeliveryStatus)
	deliveries.GET("/:id/history", deliveryHandler.GetDeliveryHistory)

	return e
}
