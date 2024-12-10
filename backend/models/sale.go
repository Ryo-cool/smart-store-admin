package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SaleItem struct {
	ProductID   primitive.ObjectID `bson:"product_id" json:"productId"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	PriceAtSale float64            `bson:"price_at_sale" json:"priceAtSale"`
}

type Sale struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Items       []SaleItem         `bson:"items" json:"items"`
	TotalAmount float64            `bson:"total_amount" json:"totalAmount"`

	// 環境影響
	TotalCO2Saved float64 `bson:"total_co2_saved" json:"totalCO2Saved"`

	// 分析用データ
	PaymentMethod string `bson:"payment_method" json:"paymentMethod"`
	TimeOfDay     string `bson:"time_of_day" json:"timeOfDay"`
	WeekDay       string `bson:"week_day" json:"weekDay"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}
