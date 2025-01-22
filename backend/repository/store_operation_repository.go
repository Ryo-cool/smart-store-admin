package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/onoderaryou/smart-store-admin/backend/models"
)

// StoreOperationRepositoryImpl は店舗運営リポジトリの実装です
type StoreOperationRepositoryImpl struct {
	collection *mongo.Collection
}

// インターフェースが実装されていることを確認
var _ StoreOperationRepository = (*StoreOperationRepositoryImpl)(nil)

func NewStoreOperationRepository(db *mongo.Database) StoreOperationRepository {
	return &StoreOperationRepositoryImpl{
		collection: db.Collection("store_operations"),
	}
}

// Create は新しい店舗運営データを記録します
func (r *StoreOperationRepositoryImpl) Create(ctx context.Context, op *models.StoreOperation) error {
	op.CreatedAt = time.Now()
	op.UpdatedAt = time.Now()
	op.Timestamp = time.Now()

	result, err := r.collection.InsertOne(ctx, op)
	if err != nil {
		return err
	}

	op.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetLatest は最新の店舗運営データを取得します
func (r *StoreOperationRepositoryImpl) GetLatest(ctx context.Context) (*models.StoreOperation, error) {
	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})
	var op models.StoreOperation
	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&op)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

// GetByTimeRange は指定期間の店舗運営データを取得します
func (r *StoreOperationRepositoryImpl) GetByTimeRange(ctx context.Context, start, end time.Time) ([]*models.StoreOperation, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	opts := options.Find().SetSort(bson.M{"timestamp": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var operations []*models.StoreOperation
	if err = cursor.All(ctx, &operations); err != nil {
		return nil, err
	}
	return operations, nil
}

// UpdateShelfStatus は特定の棚の状態を更新します
func (r *StoreOperationRepositoryImpl) UpdateShelfStatus(ctx context.Context, opID primitive.ObjectID, shelfStatus models.ShelfStatus) error {
	update := bson.M{
		"$set": bson.M{
			"shelves.$[shelf]": shelfStatus,
			"updated_at":       time.Now(),
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"shelf.shelf_id": shelfStatus.ShelfID},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": opID}, update, opts)
	return err
}

// UpdateCheckoutStatus は特定のレジの状態を更新します
func (r *StoreOperationRepositoryImpl) UpdateCheckoutStatus(ctx context.Context, opID primitive.ObjectID, checkoutStatus models.CheckoutStatus) error {
	update := bson.M{
		"$set": bson.M{
			"checkouts.$[checkout]": checkoutStatus,
			"updated_at":            time.Now(),
		},
	}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"checkout.register_id": checkoutStatus.RegisterID},
		},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": opID}, update, opts)
	return err
}

// GetAverageEnergyUsage は指定期間の平均エネルギー使用量を取得します
func (r *StoreOperationRepositoryImpl) GetAverageEnergyUsage(ctx context.Context, start, end time.Time) (map[string]float64, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$match", Value: bson.M{
				"timestamp": bson.M{
					"$gte": start,
					"$lt":  end,
				},
			}},
		},
		bson.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id":          nil,
				"avg_lighting": bson.M{"$avg": "$lighting_usage"},
				"avg_ac":       bson.M{"$avg": "$ac_usage"},
				"avg_refrig":   bson.M{"$avg": "$refrig_usage"},
			}},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AvgLighting float64 `bson:"avg_lighting"`
		AvgAC       float64 `bson:"avg_ac"`
		AvgRefrig   float64 `bson:"avg_refrig"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return map[string]float64{
			"lighting": 0,
			"ac":       0,
			"refrig":   0,
		}, nil
	}

	return map[string]float64{
		"lighting": result[0].AvgLighting,
		"ac":       result[0].AvgAC,
		"refrig":   result[0].AvgRefrig,
	}, nil
}
