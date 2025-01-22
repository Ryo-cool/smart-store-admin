package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/onoderaryou/smart-store-admin/backend/config"
	"github.com/onoderaryou/smart-store-admin/backend/handler"
	authmw "github.com/onoderaryou/smart-store-admin/backend/middleware"
)

func SetupRouter(e *echo.Echo, authHandler *handler.AuthHandler, authConfig *config.AuthConfig) {
	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 認証不要のルート
	SetupAuthRoutes(e, authHandler)

	// 認証が必要なルート
	api := e.Group("/api")
	api.Use(authmw.AuthMiddleware(authConfig))
	{
		// 管理者のみアクセス可能
		admin := api.Group("/admin")
		admin.Use(authmw.RequireRole("admin"))
		{
			// TODO: 管理者用エンドポイントの追加
		}

		// スタッフ以上がアクセス可能
		staff := api.Group("/staff")
		staff.Use(authmw.RequireRole("admin", "staff"))
		{
			// TODO: スタッフ用エンドポイントの追加
		}

		// 一般ユーザーがアクセス可能
		userGroup := api.Group("/user")
		{
			// TODO: 一般ユーザー用エンドポイントの追加
			userGroup.GET("/profile", authHandler.GetCurrentUser)
		}
	}
}

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
