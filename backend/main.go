package main

import (
	"log"

	"github.com/onoderaryou/smart-store-admin/backend/config"
	"github.com/onoderaryou/smart-store-admin/backend/db"
	"github.com/onoderaryou/smart-store-admin/backend/handler"
	"github.com/onoderaryou/smart-store-admin/backend/repository"
	"github.com/onoderaryou/smart-store-admin/backend/router"
	"github.com/onoderaryou/smart-store-admin/backend/service"
)

func main() {
	// 設定の読み込み
	cfg := config.NewConfig()

	// データベース接続
	mongodb, err := db.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer mongodb.Close()

	// リポジトリの作成
	productRepo := repository.NewProductRepository(mongodb.GetDB())
	saleRepo := repository.NewSaleRepository(mongodb.GetDB())
	deliveryRepo := repository.NewDeliveryRepository(mongodb.GetDB())
	// サービスの作成
	productService := service.NewProductService(productRepo)
	saleService := service.NewSaleService(saleRepo, productRepo)
	deliveryService := service.NewDeliveryService(deliveryRepo)
	// ハンドラーの作成
	productHandler := handler.NewProductHandler(productService)
	saleHandler := handler.NewSaleHandler(saleService)
	deliveryHandler := handler.NewDeliveryHandler(deliveryService)
	// ルーターの設定
	r := router.NewRouter(productHandler, saleHandler, deliveryHandler)

	// サーバーの起動
	if err := r.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
