package mongorepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/wishlist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WishlistRepo struct {
	col *mongo.Collection
}

func NewWishlistRepo(db *mongo.Database) *WishlistRepo {
	return &WishlistRepo{col: db.Collection("wishlist")}
}

type wishlistItemDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	ProductID primitive.ObjectID `bson:"productId"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (r *WishlistRepo) EnsureIndexes(ctx context.Context) error {
	// Create compound unique index on userId and productId
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "productId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (r *WishlistRepo) Add(ctx context.Context, item wishlist.WishlistItem) (wishlist.WishlistItem, error) {
	productOID, err := primitive.ObjectIDFromHex(item.ProductID)
	if err != nil {
		return wishlist.WishlistItem{}, wishlist.ErrInvalidProduct
	}

	doc := wishlistItemDoc{
		ID:        primitive.NewObjectID(),
		UserID:    item.UserID,
		ProductID: productOID,
		CreatedAt: item.CreatedAt,
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return wishlist.WishlistItem{}, wishlist.ErrAlreadyExists
		}
		return wishlist.WishlistItem{}, fmt.Errorf("insert wishlist item: %w", err)
	}

	item.ID = doc.ID.Hex()
	return item, nil
}

func (r *WishlistRepo) List(ctx context.Context, userID string, f wishlist.ListFilter) ([]wishlist.WishlistItem, error) {
	filter := bson.M{"userId": userID}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(f.Offset).
		SetLimit(f.Limit)

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find wishlist items: %w", err)
	}
	defer cur.Close(ctx)

	var docs []wishlistItemDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode wishlist items: %w", err)
	}

	out := make([]wishlist.WishlistItem, 0, len(docs))
	for _, d := range docs {
		out = append(out, mapWishlistItemDoc(d))
	}
	return out, nil
}

func (r *WishlistRepo) Delete(ctx context.Context, userID string, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return wishlist.ErrInvalidID
	}

	filter := bson.M{
		"_id":    oid,
		"userId": userID,
	}

	res, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("delete wishlist item: %w", err)
	}
	if res.DeletedCount == 0 {
		return wishlist.ErrNotFound
	}
	return nil
}

func (r *WishlistRepo) Exists(ctx context.Context, userID string, productID string) (bool, error) {
	productOID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return false, wishlist.ErrInvalidProduct
	}

	filter := bson.M{
		"userId":    userID,
		"productId": productOID,
	}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("count wishlist items: %w", err)
	}
	return count > 0, nil
}

func mapWishlistItemDoc(d wishlistItemDoc) wishlist.WishlistItem {
	return wishlist.WishlistItem{
		ID:        d.ID.Hex(),
		UserID:    d.UserID,
		ProductID: d.ProductID.Hex(),
		CreatedAt: d.CreatedAt,
	}
}

func (r *WishlistRepo) Count(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"userId": userID}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count wishlist items: %w", err)
	}
	return count, nil
}
