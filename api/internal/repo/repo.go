package repo

import (
	"context"
	"time"

	"happy-api/internal/models"
	"happy-api/internal/security"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	Users     *mongo.Collection
	Greetings *mongo.Collection
	Orders    *mongo.Collection
}

func New(db *mongo.Database) *Repo {
	return &Repo{
		Users:     db.Collection("users"),
		Greetings: db.Collection("greetings"),
		Orders:    db.Collection("orders"),
	}
}

func (r *Repo) InitIndexes(ctx context.Context) error {
	_, err := r.Users.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return err
	}
	_, err = r.Greetings.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "code", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
		{Keys: bson.D{{Key: "subdomain", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
		{Keys: bson.D{{Key: "expiresAt", Value: 1}}},
		{Keys: bson.D{{Key: "ownerId", Value: 1}}},
	})
	if err != nil {
		return err
	}
	_, err = r.Orders.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "ownerId", Value: 1}}},
		{Keys: bson.D{{Key: "providerRef", Value: 1}}, Options: options.Index().SetSparse(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	return err
}

func (r *Repo) EnsureAdmin(ctx context.Context, email, username, password string) error {
	count, err := r.Users.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	pwHash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	recovery := security.RandomCode(10)
	recHash, _ := security.HashPassword(recovery)

	now := time.Now()
	u := models.User{
		Email:            email,
		Username:         username,
		PasswordHash:     pwHash,
		Role:             models.RoleAdmin,
		RecoveryCodeHash: recHash,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	_, err = r.Users.InsertOne(ctx, u)
	return err
}
