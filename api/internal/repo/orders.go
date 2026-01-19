package repo

import (
	"context"
	"time"

	"happy-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) CreateOrder(ctx context.Context, owner, greetingID primitive.ObjectID, product models.GreetingType, amountRUB int, provider string) (models.Order, error) {
	now := time.Now()
	o := models.Order{
		OwnerID:     owner,
		GreetingID:  greetingID,
		ProductType: product,
		AmountRUB:   amountRUB,
		Currency:    "RUB",
		Status:      models.OrderPending,
		Provider:    provider,
		CreatedAt:   now,
	}
	res, err := r.Orders.InsertOne(ctx, o)
	if err != nil {
		return models.Order{}, err
	}
	o.ID = res.InsertedID.(primitive.ObjectID)
	return o, nil
}

func (r *Repo) SetOrderProviderRef(ctx context.Context, id primitive.ObjectID, ref string) error {
	_, err := r.Orders.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"providerRef": ref}})
	return err
}

func (r *Repo) ListOrdersByOwner(ctx context.Context, owner primitive.ObjectID) ([]models.Order, error) {
	cur, err := r.Orders.Find(ctx, bson.M{"ownerId": owner}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(200))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Order
	for cur.Next(ctx) {
		var o models.Order
		if err := cur.Decode(&o); err == nil {
			out = append(out, o)
		}
	}
	return out, nil
}

func (r *Repo) GetOrderOwned(ctx context.Context, id, owner primitive.ObjectID) (models.Order, error) {
	var o models.Order
	err := r.Orders.FindOne(ctx, bson.M{"_id": id, "ownerId": owner}).Decode(&o)
	return o, err
}

func (r *Repo) FindOrderByProviderRef(ctx context.Context, ref string) (models.Order, error) {
	var o models.Order
	err := r.Orders.FindOne(ctx, bson.M{"providerRef": ref}).Decode(&o)
	return o, err
}

func (r *Repo) MarkOrderPaid(ctx context.Context, id primitive.ObjectID) (models.Order, error) {
	now := time.Now()
	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out models.Order
	err := r.Orders.FindOneAndUpdate(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": models.OrderPaid, "paidAt": now}},
		after,
	).Decode(&out)
	return out, err
}

func (r *Repo) AdminListOrders(ctx context.Context, limit int) ([]models.Order, error) {
	cur, err := r.Orders.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Order
	for cur.Next(ctx) {
		var o models.Order
		if err := cur.Decode(&o); err == nil {
			out = append(out, o)
		}
	}
	return out, nil
}

func (r *Repo) AdminSalesStats(ctx context.Context) (paidCount int64, paidSum int64, err error) {
	paidCount, err = r.Orders.CountDocuments(ctx, bson.M{"status": models.OrderPaid})
	if err != nil {
		return 0, 0, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"status": models.OrderPaid}},
		{"$group": bson.M{"_id": nil, "sum": bson.M{"$sum": "$amountRUB"}}},
	}
	cur, err := r.Orders.Aggregate(ctx, pipeline)
	if err != nil {
		return paidCount, 0, err
	}
	defer cur.Close(ctx)

	var row struct {
		Sum int64 `bson:"sum"`
	}
	if cur.Next(ctx) {
		_ = cur.Decode(&row)
		paidSum = row.Sum
	}
	return paidCount, paidSum, nil
}

func (r *Repo) GetOrderByID(ctx context.Context, id primitive.ObjectID) (models.Order, error) {
	var o models.Order
	err := r.Orders.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	return o, err
}
