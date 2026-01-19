package repo

import (
	"context"
	"strings"
	"time"

	"happy-api/internal/models"
	"happy-api/internal/security"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) CreateGreeting(ctx context.Context, owner primitive.ObjectID, in models.Greeting) (models.Greeting, error) {
	now := time.Now()

	g := models.Greeting{
		OwnerID:   owner,
		Type:      in.Type,
		Code:      "",
		Subdomain: strings.ToLower(strings.TrimSpace(in.Subdomain)),
		Title:     strings.TrimSpace(in.Title),
		Body:      strings.TrimSpace(in.Body),
		Theme:     in.Theme,
		Gift:      in.Gift,
		Photos:    []string{},
		Paid:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if g.Title == "" {
		g.Title = "С поздравлением!"
	}
	if g.Body == "" {
		g.Body = "Пусть этот день будет особенным ✨"
	}

	if g.Type == models.GreetingCode {
		g.Code = generateCode()
	}
	if g.Theme.Kind == "" {
		g.Theme = models.Theme{Kind: "birthday"}
	}
	if g.Gift.Kind == "" {
		g.Gift = models.Gift{Kind: models.GiftText, Value: "Открой подарок и улыбнись 🙂"}
	}

	res, err := r.Greetings.InsertOne(ctx, g)
	if err != nil {
		return models.Greeting{}, err
	}
	g.ID = res.InsertedID.(primitive.ObjectID)
	return g, nil
}

func generateCode() string {
	// create 6-char uppercase safe code from random base64url
	raw := strings.ToUpper(security.RandomCode(10))
	raw = strings.NewReplacer("-", "", "_", "", "O", "0").Replace(raw)
	var out strings.Builder
	for _, ch := range raw {
		if (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			out.WriteRune(ch)
		}
		if out.Len() >= 6 {
			break
		}
	}
	s := out.String()
	if len(s) < 6 {
		return generateCode()
	}
	return s
}

func (r *Repo) ListGreetingsByOwner(ctx context.Context, owner primitive.ObjectID) ([]models.Greeting, error) {
	cur, err := r.Greetings.Find(ctx, bson.M{"ownerId": owner}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(200))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []models.Greeting
	for cur.Next(ctx) {
		var g models.Greeting
		if err := cur.Decode(&g); err == nil {
			items = append(items, g)
		}
	}
	return items, nil
}

func (r *Repo) GetGreetingOwned(ctx context.Context, id, owner primitive.ObjectID) (models.Greeting, error) {
	var g models.Greeting
	err := r.Greetings.FindOne(ctx, bson.M{"_id": id, "ownerId": owner}).Decode(&g)
	return g, err
}

func (r *Repo) GetGreetingByID(ctx context.Context, id primitive.ObjectID) (models.Greeting, error) {
	var g models.Greeting
	err := r.Greetings.FindOne(ctx, bson.M{"_id": id}).Decode(&g)
	return g, err
}

func (r *Repo) UpdateGreetingOwned(ctx context.Context, id, owner primitive.ObjectID, patch bson.M) error {
	patch["updatedAt"] = time.Now()
	_, err := r.Greetings.UpdateOne(ctx, bson.M{"_id": id, "ownerId": owner}, bson.M{"$set": patch})
	return err
}

func (r *Repo) DeleteGreeting(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Greetings.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repo) MarkGreetingPaidAndPublish(ctx context.Context, id primitive.ObjectID, ttl time.Duration, dns *models.DNSRecord) error {
	now := time.Now()
	exp := now.Add(ttl)
	set := bson.M{
		"paid":        true,
		"publishedAt": now,
		"expiresAt":   exp,
		"updatedAt":   now,
	}
	if dns != nil {
		set["dns"] = dns
	}
	_, err := r.Greetings.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": set})
	return err
}

func (r *Repo) FindPublicByCode(ctx context.Context, code string) (models.Greeting, error) {
	now := time.Now()
	var g models.Greeting
	filter := bson.M{
		"type": "code",
		"code": strings.ToUpper(strings.TrimSpace(code)),
		"paid": true,
		"$or": []bson.M{
			{"expiresAt": bson.M{"$gt": now}},
			{"expiresAt": bson.M{"$exists": false}},
		},
	}
	err := r.Greetings.FindOne(ctx, filter).Decode(&g)
	return g, err
}

func (r *Repo) FindPublicBySubdomain(ctx context.Context, sub string) (models.Greeting, error) {
	now := time.Now()
	var g models.Greeting
	filter := bson.M{
		"type":      "subdomain",
		"subdomain": strings.ToLower(strings.TrimSpace(sub)),
		"paid":      true,
		"$or": []bson.M{
			{"expiresAt": bson.M{"$gt": now}},
			{"expiresAt": bson.M{"$exists": false}},
		},
	}
	err := r.Greetings.FindOne(ctx, filter).Decode(&g)
	return g, err
}

func (r *Repo) FindExpiredGreetings(ctx context.Context, limit int) ([]models.Greeting, error) {
	now := time.Now()
	cur, err := r.Greetings.Find(ctx, bson.M{
		"paid":      true,
		"expiresAt": bson.M{"$lte": now},
	}, options.Find().SetSort(bson.D{{Key: "expiresAt", Value: 1}}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
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

func (r *Repo) MarkGreetingPaidPublished(ctx context.Context, gid primitive.ObjectID, days int) error {
	now := time.Now()
	exp := now.Add(time.Duration(days) * 24 * time.Hour)

	_, err := r.Greetings.UpdateOne(ctx,
		bson.M{"_id": gid},
		bson.M{"$set": bson.M{
			"paid":        true,
			"publishedAt": now,
			"expiresAt":   exp,
			"updatedAt":   now,
		}},
	)
	return err
}
