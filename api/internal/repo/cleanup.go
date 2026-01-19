package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Repo) DeleteOrdersByGreeting(ctx context.Context, gid primitive.ObjectID) error {
	_, err := r.Orders.DeleteMany(ctx, bson.M{"greetingId": gid})
	return err
}
