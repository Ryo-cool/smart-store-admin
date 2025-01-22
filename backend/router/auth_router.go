package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onoderaryou/smart-store-admin/backend/handler"
)

func SetupAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
	// 認証関連のエンドポイント
	auth := e.Group("/api/auth")
	{
		// Google OAuth2認証
		auth.GET("/google", func(c echo.Context) error {
			url := authHandler.GetGoogleAuthURL()
			return c.JSON(http.StatusOK, echo.Map{"url": url})
		})

		// Google OAuth2コールバック
		auth.GET("/google/callback", authHandler.HandleGoogleCallback)

		// ログアウト
		auth.POST("/logout", func(c echo.Context) error {
			return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
		})

		// 現在のユーザー情報取得
		auth.GET("/me", authHandler.GetCurrentUser)
	}
}
