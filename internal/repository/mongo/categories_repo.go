package mongorepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/categories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoriesRepo struct {
	col *mongo.Collection
}

func NewCategoriesRepo(db *mongo.Database) *CategoriesRepo {
	return &CategoriesRepo{col: db.Collection("categories")}
}

type categoryDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func (r *CategoriesRepo) List(ctx context.Context, f categories.ListFilter) ([]categories.Category, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(f.Offset).
		SetLimit(f.Limit)

	cur, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find categories: %w", err)
	}
	defer cur.Close(ctx)

	var docs []categoryDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode categories: %w", err)
	}

	out := make([]categories.Category, 0, len(docs))
	for _, d := range docs {
		out = append(out, categories.Category{
			ID:          d.ID.Hex(),
			Name:        d.Name,
			Description: d.Description,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}
	return out, nil
}

func (r *CategoriesRepo) GetByID(ctx context.Context, id string) (categories.Category, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return categories.Category{}, categories.ErrInvalidID
	}

	var d categoryDoc
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return categories.Category{}, categories.ErrNotFound
		}
		return categories.Category{}, fmt.Errorf("find category: %w", err)
	}

	return categories.Category{
		ID:          d.ID.Hex(),
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

func (r *CategoriesRepo) Create(ctx context.Context, c categories.Category) (categories.Category, error) {
	doc := categoryDoc{
		ID:          primitive.NewObjectID(),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	if _, err := r.col.InsertOne(ctx, doc); err != nil {
		return categories.Category{}, fmt.Errorf("insert category: %w", err)
	}

	c.ID = doc.ID.Hex()
	return c, nil
}

func (r *CategoriesRepo) Update(ctx context.Context, id string, in categories.UpdateInput) (categories.Category, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return categories.Category{}, categories.ErrInvalidID
	}

	set := bson.M{
		"updatedAt": time.Now().UTC(),
	}
	if in.Name != nil {
		set["name"] = *in.Name
	}
	if in.Description != nil {
		set["description"] = *in.Description
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var d categoryDoc
	err = r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": set},
		opts,
	).Decode(&d)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return categories.Category{}, categories.ErrNotFound
		}
		return categories.Category{}, fmt.Errorf("update category: %w", err)
	}

	return categories.Category{
		ID:          d.ID.Hex(),
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

func (r *CategoriesRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return categories.ErrInvalidID
	}

	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if res.DeletedCount == 0 {
		return categories.ErrNotFound
	}
	return nil
}
