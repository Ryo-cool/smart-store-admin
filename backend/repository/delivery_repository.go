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

type DeliveryRepository struct {
	collection *mongo.Collection
}

func NewDeliveryRepository(db *mongo.Database) *DeliveryRepository {
	return &DeliveryRepository{
		collection: db.Collection("deliveries"),
	}
}

// Create は新しい配送を作成します
func (r *DeliveryRepository) Create(ctx context.Context, delivery *models.Delivery) error {
	delivery.CreatedAt = time.Now()
	delivery.UpdatedAt = time.Now()
	delivery.Status = models.StatusPending

	result, err := r.collection.InsertOne(ctx, delivery)
	if err != nil {
		return err
	}

	delivery.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID は指定されたIDの配送を取得します
func (r *DeliveryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error) {
	var delivery models.Delivery
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&delivery)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

// List は配送のリストを取得します
func (r *DeliveryRepository) List(ctx context.Context, skip, limit int64) ([]*models.Delivery, error) {
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
func (r *DeliveryRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error {
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
func (r *DeliveryRepository) UpdateLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
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
func (r *DeliveryRepository) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	filter := bson.M{
		"status": bson.M{
			"$in": []models.DeliveryStatus{models.StatusPending, models.StatusInProgress},
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
func (r *DeliveryRepository) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
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
