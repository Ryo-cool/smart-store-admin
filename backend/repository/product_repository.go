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

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

// Create は新しい商品を作成します
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID は指定されたIDの商品を取得します
func (r *ProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List は商品のリストを取得します
func (r *ProductRepository) List(ctx context.Context, skip, limit int64) ([]*models.Product, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// Update は商品情報を更新します
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete は商品を削除します
func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetByCategory は指定されたカテゴリの商品を取得します
func (r *ProductRepository) GetByCategory(ctx context.Context, category string) ([]*models.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// GetLowStock は在庫が最小在庫レベルを下回っている商品を取得します
func (r *ProductRepository) GetLowStock(ctx context.Context) ([]*models.Product, error) {
	filter := bson.M{
		"stock": bson.M{
			"$lte": bson.M{"$ref": "min_stock_level"},
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
