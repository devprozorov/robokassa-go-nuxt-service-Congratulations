package repo

import (
	"context"

	"happy-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) AdminListUsers(ctx context.Context, limit int) ([]models.User, error) {
	cur, err := r.Users.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)))
	if err != nil { return nil, err }
	defer cur.Close(ctx)

	var out []models.User
	for cur.Next(ctx) {
		var u models.User
		if err := cur.Decode(&u); err == nil {
			u.PasswordHash = ""
			u.RecoveryCodeHash = ""
			out = append(out, u)
		}
	}
	return out, nil
}

func (r *Repo) AdminDeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Users.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repo) AdminListGreetings(ctx context.Context, limit int) ([]models.Greeting, error) {
	cur, err := r.Greetings.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)))
	if err != nil { return nil, err }
	defer cur.Close(ctx)

	var out []models.Greeting
	for cur.Next(ctx) {
		var g models.Greeting
		if err := cur.Decode(&g); err == nil {
			out = append(out, g)
		}
	}
	return out, nil
}
