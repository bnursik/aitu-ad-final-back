package mongorepo

import (
	"context"
	"fmt"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/categories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductsCounterRepo struct {
	col *mongo.Collection
}

func NewProductsCounterRepo(db *mongo.Database) *ProductsCounterRepo {
	return &ProductsCounterRepo{col: db.Collection("products")}
}

func (r *ProductsCounterRepo) CountByCategoryID(ctx context.Context, categoryID string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return 0, categories.ErrInvalidID
	}

	n, err := r.col.CountDocuments(ctx, bson.M{"categoryId": oid})
	if err != nil {
		return 0, fmt.Errorf("count products by category: %w", err)
	}
	return n, nil
}
