package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"happy-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repo) CreateUser(ctx context.Context, email, username, pwHash, recoveryHash string) (models.User, error) {
	now := time.Now()
	u := models.User{
		Email:            strings.ToLower(email),
		Username:         username,
		PasswordHash:     pwHash,
		Role:             models.RoleUser,
		RecoveryCodeHash: recoveryHash,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	res, err := r.Users.InsertOne(ctx, u)
	if err != nil {
		return models.User{}, err
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
	return u, nil
}

func (r *Repo) FindUserByLogin(ctx context.Context, loginOrEmail string) (models.User, error) {
	var u models.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": strings.ToLower(loginOrEmail)},
			{"username": loginOrEmail},
		},
	}
	err := r.Users.FindOne(ctx, filter).Decode(&u)
	return u, err
}

func (r *Repo) FindUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var u models.User
	err := r.Users.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	return u, err
}

func (r *Repo) UpdateUserPasswordAndRecovery(ctx context.Context, id primitive.ObjectID, pwHash, recoveryHash string, usedAt *time.Time) error {
	set := bson.M{
		"passwordHash":     pwHash,
		"recoveryCodeHash": recoveryHash,
		"updatedAt":        time.Now(),
	}

	if usedAt == nil {
		set["recoveryCodeUsedAt"] = nil
	} else {
		set["recoveryCodeUsedAt"] = *usedAt
	}

	_, err := r.Users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": set})
	return err
}

// helper: duplicate key checker (optional)
func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	var we mongo.WriteException
	if ok := errors.As(err, &we); ok {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}
