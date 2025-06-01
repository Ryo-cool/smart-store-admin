package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	SKU         string             `bson:"sku" json:"sku"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	Status      string             `bson:"status" json:"status"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Weight      string             `bson:"weight" json:"weight"`
	Dimensions  string             `bson:"dimensions" json:"dimensions"`
	Images      []string           `bson:"images" json:"images"`

	// 環境負荷関連
	CO2Emission float64 `bson:"co2_emission" json:"co2Emission"`
	RecycleRate float64 `bson:"recycle_rate" json:"recycleRate"`

	// 商品配置
	ShelfLocation string `bson:"shelf_location" json:"shelfLocation"`

	// 在庫管理
	MinStockLevel int `bson:"min_stock_level" json:"minStockLevel"`
	ReorderPoint  int `bson:"reorder_point" json:"reorderPoint"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
