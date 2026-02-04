package mongorepo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/products"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductsRepo struct {
	col *mongo.Collection
}

func NewProductsRepo(db *mongo.Database) *ProductsRepo {
	return &ProductsRepo{col: db.Collection("products")}
}

type reviewDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	Rating    int64              `bson:"rating"`
	Comment   string             `bson:"comment,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"`
}

type productDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID `bson:"categoryId"`
	Name        string             `bson:"name"`
	Description string             `bson:"description,omitempty"`
	Price       float64            `bson:"price"`
	Stock       int64              `bson:"stock"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	Reviews     []reviewDoc        `bson:"reviews,omitempty"`
}

func (r *ProductsRepo) List(ctx context.Context, f products.ListFilter) ([]products.Product, error) {
	filter := bson.M{}
	if f.CategoryID != nil && strings.TrimSpace(*f.CategoryID) != "" {
		oid, err := primitive.ObjectIDFromHex(*f.CategoryID)
		if err != nil {
			return nil, products.ErrInvalidCategory
		}
		filter["categoryId"] = oid
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(f.Offset).
		SetLimit(f.Limit)

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find products: %w", err)
	}
	defer cur.Close(ctx)

	var docs []productDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode products: %w", err)
	}

	out := make([]products.Product, 0, len(docs))
	for _, d := range docs {
		out = append(out, mapProductDoc(d))
	}
	return out, nil
}

func (r *ProductsRepo) GetByID(ctx context.Context, id string) (products.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return products.Product{}, products.ErrInvalidID
	}

	var d productDoc
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return products.Product{}, products.ErrNotFound
		}
		return products.Product{}, fmt.Errorf("find product: %w", err)
	}
	return mapProductDoc(d), nil
}

func (r *ProductsRepo) Create(ctx context.Context, p products.Product) (products.Product, error) {
	catOID, err := primitive.ObjectIDFromHex(p.CategoryID)
	if err != nil {
		return products.Product{}, products.ErrInvalidCategory
	}

	doc := productDoc{
		ID:          primitive.NewObjectID(),
		CategoryID:  catOID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Reviews:     []reviewDoc{},
	}

	if _, err := r.col.InsertOne(ctx, doc); err != nil {
		return products.Product{}, fmt.Errorf("insert product: %w", err)
	}

	p.ID = doc.ID.Hex()
	return p, nil
}

func (r *ProductsRepo) Update(ctx context.Context, id string, in products.UpdateInput) (products.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return products.Product{}, products.ErrInvalidID
	}

	set := bson.M{
		"updatedAt": time.Now().UTC(),
	}

	if in.CategoryID != nil {
		catOID, err := primitive.ObjectIDFromHex(*in.CategoryID)
		if err != nil {
			return products.Product{}, products.ErrInvalidCategory
		}
		set["categoryId"] = catOID
	}
	if in.Name != nil {
		set["name"] = *in.Name
	}
	if in.Description != nil {
		set["description"] = *in.Description
	}
	if in.Price != nil {
		set["price"] = *in.Price
	}
	if in.Stock != nil {
		set["stock"] = *in.Stock
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var d productDoc
	err = r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": set},
		opts,
	).Decode(&d)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return products.Product{}, products.ErrNotFound
		}
		return products.Product{}, fmt.Errorf("update product: %w", err)
	}

	return mapProductDoc(d), nil
}

func (r *ProductsRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return products.ErrInvalidID
	}

	filter := bson.M{"_id": oid, "stock": bson.M{"$lt": 1}}
	res, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if res.DeletedCount == 0 {
		n, err := r.col.CountDocuments(ctx, bson.M{"_id": oid})
		if err != nil {
			return fmt.Errorf("check product: %w", err)
		}
		if n == 0 {
			return products.ErrNotFound
		}
		return products.ErrCannotDeleteProduct
	}
	return nil
}

func (r *ProductsRepo) DecrementStock(ctx context.Context, productID string, qty int64) error {
	oid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return products.ErrInvalidID
	}
	if qty <= 0 {
		return nil
	}

	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": oid, "stock": bson.M{"$gte": qty}},
		bson.M{
			"$inc":  bson.M{"stock": -qty},
			"$set":  bson.M{"updatedAt": time.Now().UTC()},
		},
	)
	if err != nil {
		return fmt.Errorf("decrement stock: %w", err)
	}
	if res.MatchedCount == 0 {
		n, err := r.col.CountDocuments(ctx, bson.M{"_id": oid})
		if err != nil {
			return fmt.Errorf("check product: %w", err)
		}
		if n == 0 {
			return products.ErrNotFound
		}
		return products.ErrInsufficientStock
	}
	return nil
}

func (r *ProductsRepo) AddReview(ctx context.Context, productID string, rev products.Review) (products.Review, error) {
	pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return products.Review{}, products.ErrInvalidID
	}

	doc := reviewDoc{
		ID:        primitive.NewObjectID(),
		UserID:    rev.UserID,
		Rating:    rev.Rating,
		Comment:   rev.Comment,
		CreatedAt: rev.CreatedAt,
	}

	// $push embedded review
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": pid},
		bson.M{"$push": bson.M{"reviews": doc}},
	)
	if err != nil {
		return products.Review{}, fmt.Errorf("push review: %w", err)
	}
	if res.MatchedCount == 0 {
		return products.Review{}, products.ErrNotFound
	}

	return products.Review{
		ID:        doc.ID.Hex(),
		UserID:    doc.UserID,
		Rating:    doc.Rating,
		Comment:   doc.Comment,
		CreatedAt: doc.CreatedAt,
	}, nil
}

func (r *ProductsRepo) DeleteReview(ctx context.Context, productID string, reviewID string) error {
	pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return products.ErrInvalidID
	}
	rid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return products.ErrInvalidReviewID
	}

	// $pull embedded review by _id
	res, err := r.col.UpdateOne(ctx,
		bson.M{"_id": pid},
		bson.M{"$pull": bson.M{"reviews": bson.M{"_id": rid}}},
	)
	if err != nil {
		return fmt.Errorf("pull review: %w", err)
	}
	if res.MatchedCount == 0 {
		return products.ErrNotFound
	}

	// ModifiedCount could be 0 if review not found; treat as not found for review
	if res.ModifiedCount == 0 {
		return products.ErrNotFound
	}
	return nil
}

func mapProductDoc(d productDoc) products.Product {
	out := products.Product{
		ID:          d.ID.Hex(),
		CategoryID:  d.CategoryID.Hex(),
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Stock:       d.Stock,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		Reviews:     make([]products.Review, 0, len(d.Reviews)),
	}
	for _, r := range d.Reviews {
		out.Reviews = append(out.Reviews, products.Review{
			ID:        r.ID.Hex(),
			UserID:    r.UserID,
			Rating:    r.Rating,
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt,
		})
	}
	return out
}

func (r *ProductsRepo) Count(ctx context.Context, f products.ListFilter) (int64, error) {
	filter := bson.M{}
	if f.CategoryID != nil && strings.TrimSpace(*f.CategoryID) != "" {
		oid, err := primitive.ObjectIDFromHex(*f.CategoryID)
		if err != nil {
			return 0, products.ErrInvalidCategory
		}
		filter["categoryId"] = oid
	}

	n, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}
	return n, nil
}
