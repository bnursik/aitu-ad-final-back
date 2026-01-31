package mongorepo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepo struct {
	col *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{col: db.Collection("users")}
}

func (r *UsersRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_email"),
	})
	return err
}

type userDoc struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	Role         string             `bson:"role"`
	Address      string             `bson:"address,omitempty"`
	Phone        string             `bson:"phone,omitempty"`
	Bio          string             `bson:"bio,omitempty"`
	CreatedAt    time.Time          `bson:"created_at"`
}

func (r *UsersRepo) Insert(ctx context.Context, u users.User) (users.User, error) {
	doc := userDoc{
		ID:           primitive.NewObjectID(),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt,
	}

	_, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return users.User{}, users.ErrEmailTaken
		}
		return users.User{}, fmt.Errorf("insert user: %w", err)
	}

	u.ID = doc.ID.Hex()
	return u, nil
}

func (r *UsersRepo) FindByEmail(ctx context.Context, email string) (users.User, error) {
	var doc userDoc
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return users.User{}, users.ErrUserNotFound
		}
		return users.User{}, fmt.Errorf("find by email: %w", err)
	}

	return mapUserDoc(doc), nil
}

func (r *UsersRepo) FindByID(ctx context.Context, id string) (users.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return users.User{}, users.ErrUserNotFound
	}

	var doc userDoc
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return users.User{}, users.ErrUserNotFound
		}
		return users.User{}, fmt.Errorf("find by id: %w", err)
	}

	return mapUserDoc(doc), nil
}

func (r *UsersRepo) Update(ctx context.Context, id string, in users.UpdateProfileInput) (users.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return users.User{}, users.ErrUserNotFound
	}

	set := bson.M{}
	if in.Name != nil {
		set["name"] = *in.Name
	}
	if in.Address != nil {
		set["address"] = *in.Address
	}
	if in.Phone != nil {
		set["phone"] = *in.Phone
	}
	if in.Bio != nil {
		set["bio"] = *in.Bio
	}

	if len(set) == 0 {
		return r.FindByID(ctx, id)
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var doc userDoc
	err = r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": set},
		opts,
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return users.User{}, users.ErrUserNotFound
		}
		return users.User{}, fmt.Errorf("update user: %w", err)
	}

	return mapUserDoc(doc), nil
}

func mapUserDoc(doc userDoc) users.User {
	return users.User{
		ID:           doc.ID.Hex(),
		Name:         doc.Name,
		Email:        doc.Email,
		PasswordHash: doc.PasswordHash,
		Role:         users.Role(doc.Role),
		Address:      doc.Address,
		Phone:        doc.Phone,
		Bio:          doc.Bio,
		CreatedAt:    doc.CreatedAt,
	}
}

func (r *UsersRepo) GetAll(ctx context.Context) ([]users.User, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []userDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}

	result := make([]users.User, len(docs))
	for i, doc := range docs {
		result[i] = mapUserDoc(doc)
	}

	return result, nil
}
