package mongorepo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/orders"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrdersRepo struct {
	col *mongo.Collection
}

func NewOrdersRepo(db *mongo.Database) *OrdersRepo {
	return &OrdersRepo{col: db.Collection("orders")}
}

type orderItemDoc struct {
	ProductID primitive.ObjectID `bson:"productId"`
	Quantity  int64              `bson:"quantity"`
}

type orderDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	Items     []orderItemDoc     `bson:"items"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func (r *OrdersRepo) List(ctx context.Context, userID *string, f orders.ListFilter) ([]orders.Order, error) {
	filter := bson.M{}
	if userID != nil {
		filter["userId"] = *userID
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(f.Offset).
		SetLimit(f.Limit)

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find orders: %w", err)
	}
	defer cur.Close(ctx)

	var docs []orderDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode orders: %w", err)
	}

	out := make([]orders.Order, 0, len(docs))
	for _, d := range docs {
		out = append(out, mapOrderDoc(d))
	}
	return out, nil
}

func (r *OrdersRepo) GetByID(ctx context.Context, id string, userID *string) (orders.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return orders.Order{}, orders.ErrInvalidID
	}

	filter := bson.M{"_id": oid}
	if userID != nil {
		filter["userId"] = *userID
	}

	var d orderDoc
	if err := r.col.FindOne(ctx, filter).Decode(&d); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return orders.Order{}, orders.ErrNotFound
		}
		return orders.Order{}, fmt.Errorf("find order: %w", err)
	}

	return mapOrderDoc(d), nil
}

func (r *OrdersRepo) Create(ctx context.Context, o orders.Order) (orders.Order, error) {
	items := make([]orderItemDoc, 0, len(o.Items))
	for _, it := range o.Items {
		pid, err := primitive.ObjectIDFromHex(it.ProductID)
		if err != nil {
			return orders.Order{}, orders.ErrInvalidProduct
		}
		items = append(items, orderItemDoc{ProductID: pid, Quantity: it.Quantity})
	}

	doc := orderDoc{
		ID:        primitive.NewObjectID(),
		UserID:    o.UserID,
		Items:     items,
		Status:    string(o.Status),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}

	if _, err := r.col.InsertOne(ctx, doc); err != nil {
		return orders.Order{}, fmt.Errorf("insert order: %w", err)
	}

	o.ID = doc.ID.Hex()
	return o, nil
}

func (r *OrdersRepo) UpdateStatus(ctx context.Context, id string, status orders.Status) (orders.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return orders.Order{}, orders.ErrInvalidID
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var d orderDoc
	err = r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": string(status), "updatedAt": time.Now().UTC()}},
		opts,
	).Decode(&d)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return orders.Order{}, orders.ErrNotFound
		}
		return orders.Order{}, fmt.Errorf("update order status: %w", err)
	}

	return mapOrderDoc(d), nil
}

func mapOrderDoc(d orderDoc) orders.Order {
	items := make([]orders.Item, 0, len(d.Items))
	for _, it := range d.Items {
		items = append(items, orders.Item{
			ProductID: it.ProductID.Hex(),
			Quantity:  it.Quantity,
		})
	}

	return orders.Order{
		ID:        d.ID.Hex(),
		UserID:    d.UserID,
		Items:     items,
		Status:    orders.Status(d.Status),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *OrdersRepo) Count(ctx context.Context, userID *string) (int64, error) {
	filter := bson.M{}
	if userID != nil && strings.TrimSpace(*userID) != "" {
		filter["userId"] = *userID
	}

	n, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count orders: %w", err)
	}
	return n, nil
}
