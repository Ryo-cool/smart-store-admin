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

// DeliveryRepositoryImpl は配送リポジトリの実装です
type DeliveryRepositoryImpl struct {
	collection *mongo.Collection
}

// インターフェースが実装されていることを確認
var _ DeliveryRepository = (*DeliveryRepositoryImpl)(nil)

func NewDeliveryRepository(db *mongo.Database) DeliveryRepository {
	return &DeliveryRepositoryImpl{
		collection: db.Collection("deliveries"),
	}
}

// Create は新しい配送を作成します
func (r *DeliveryRepositoryImpl) Create(ctx context.Context, delivery *models.Delivery) error {
	delivery.CreatedAt = time.Now()
	delivery.UpdatedAt = time.Now()
	delivery.Status = models.StatusPreparing

	result, err := r.collection.InsertOne(ctx, delivery)
	if err != nil {
		return err
	}

	delivery.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// GetByID は指定されたIDの配送を取得します
func (r *DeliveryRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error) {
	var delivery models.Delivery
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&delivery)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

// List は配送のリストを取得します
func (r *DeliveryRepositoryImpl) List(ctx context.Context, skip, limit int64) ([]*models.Delivery, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deliveries []*models.Delivery
	if err = cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// UpdateStatus は配送のステータスを更新します
func (r *DeliveryRepositoryImpl) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	if status == models.StatusCompleted {
		update["$set"].(bson.M)["completed_at"] = time.Now()
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateLocation はロボット/ドローンの現在位置を更新します
func (r *DeliveryRepositoryImpl) UpdateLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
	update := bson.M{
		"$set": bson.M{
			"current_location": location,
			"updated_at":       time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// GetActiveDeliveries はアクティブな配送（進行中のもの）を取得します
func (r *DeliveryRepositoryImpl) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	filter := bson.M{
		"status": bson.M{
			"$in": []models.DeliveryStatus{models.StatusPreparing, models.StatusInProgress},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deliveries []*models.Delivery
	if err = cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// GetDeliveriesByRobot は特定のロボット/ドローンの配送履歴を取得します
func (r *DeliveryRepositoryImpl) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"robot_id": robotID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deliveries []*models.Delivery
	if err = cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// GetDeliveries retrieves deliveries based on query parameters
func (r *DeliveryRepositoryImpl) GetDeliveries(query *models.DeliveryQuery) (*models.DeliveryResponse, error) {
	collection := r.collection
	ctx := context.Background()

	filter := bson.M{}
	if query.Status != nil {
		filter["status"] = *query.Status
	}
	if query.Search != nil {
		filter["$or"] = []bson.M{
			{"id": bson.M{"$regex": *query.Search, "$options": "i"}},
			{"address": bson.M{"$regex": *query.Search, "$options": "i"}},
		}
	}

	opts := options.Find()
	if query.Page != nil && query.Limit != nil {
		opts.SetSkip(int64((*query.Page - 1) * *query.Limit))
		opts.SetLimit(int64(*query.Limit))
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deliveries []models.Delivery
	if err := cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.DeliveryResponse{
		Deliveries: deliveries,
		Total:      total,
	}, nil
}

// GetDeliveryHistory は配送履歴を取得します
func (r *DeliveryRepositoryImpl) GetDeliveryHistory(ctx context.Context, id primitive.ObjectID) (*models.DeliveryHistoryResponse, error) {
	collection := r.collection.Database().Collection("delivery_history")

	cursor, err := collection.Find(ctx, bson.M{"delivery_id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.DeliveryHistory
	if err := cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	return &models.DeliveryHistoryResponse{
		History: history,
	}, nil
}

// Update は配送情報を更新します
func (r *DeliveryRepositoryImpl) Update(ctx context.Context, id primitive.ObjectID, delivery *models.Delivery) error {
	delivery.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": delivery},
	)
	return err
}
