package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"smart-store-admin/backend/models"
)

type SaleRepository struct {
	collection *mongo.Collection
}

func NewSaleRepository(db *mongo.Database) *SaleRepository {
	return &SaleRepository{
		collection: db.Collection("sales"),
	}
}

// Create は新しい売上を記録します
func (r *SaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	sale.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, sale)
	if err != nil {
		return err
	}

	sale.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID は指定されたIDの売上を取得します
func (r *SaleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Sale, error) {
	var sale models.Sale
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sale)
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

// GetDailySales は日付の売上を取得します
func (r *SaleRepository) GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*models.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}
	return sales, nil
}

// GetSalesByTimeOfDay は時間帯別の売上を取得します
func (r *SaleRepository) GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"time_of_day": timeOfDay})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*models.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}
	return sales, nil
}

// GetSalesByDateRange は指定期間の売上を取得します
func (r *SaleRepository) GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	opts := options.Find().SetSort(bson.M{"created_at": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*models.Sale
	if err = cursor.All(ctx, &sales); err != nil {
		return nil, err
	}
	return sales, nil
}

// GetTotalSalesAmount は指定期間の総売上金額を取得します
func (r *SaleRepository) GetTotalSalesAmount(ctx context.Context, start, end time.Time) (float64, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$match", Value: bson.M{
				"created_at": bson.M{
					"$gte": start,
					"$lt":  end,
				},
			}},
		},
		bson.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id": nil,
				"total": bson.M{
					"$sum": "$total_amount",
				},
			}},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}
