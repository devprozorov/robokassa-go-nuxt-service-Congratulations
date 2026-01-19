package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email              string             `bson:"email" json:"email"`
	Username           string             `bson:"username" json:"username"`
	PasswordHash       string             `bson:"passwordHash" json:"-"`
	Role               Role               `bson:"role" json:"role"`
	RecoveryCodeHash   string             `bson:"recoveryCodeHash" json:"-"`
	RecoveryCodeUsedAt *time.Time         `bson:"recoveryCodeUsedAt,omitempty" json:"-"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type GreetingType string

const (
	GreetingCode      GreetingType = "code"
	GreetingSubdomain GreetingType = "subdomain"
)

type Theme struct {
	Kind       string `bson:"kind" json:"kind"`
	Primary    string `bson:"primary,omitempty" json:"primary,omitempty"`
	Secondary  string `bson:"secondary,omitempty" json:"secondary,omitempty"`
	Background string `bson:"background,omitempty" json:"background,omitempty"`
}

type GiftKind string

const (
	GiftPromo    GiftKind = "promo"
	GiftImage    GiftKind = "image"
	GiftText     GiftKind = "text"
	GiftRedirect GiftKind = "redirect"
)

type Gift struct {
	Kind  GiftKind `bson:"kind" json:"kind"`
	Value string   `bson:"value" json:"value"`
}

type DNSRecord struct {
	Domain    string    `bson:"domain" json:"domain"`
	Type      string    `bson:"type" json:"type"`
	Name      string    `bson:"name" json:"name"`
	Value     string    `bson:"value" json:"value"`
	TTL       int       `bson:"ttl" json:"ttl"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Greeting struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID    primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	Type       GreetingType       `bson:"type" json:"type"`
	Code       string             `bson:"code,omitempty" json:"code,omitempty"`
	Subdomain  string             `bson:"subdomain,omitempty" json:"subdomain,omitempty"`
	Title      string             `bson:"title" json:"title"`
	Body       string             `bson:"body" json:"body"`
	Theme      Theme              `bson:"theme" json:"theme"`
	Gift       Gift               `bson:"gift" json:"gift"`
	Photos     []string           `bson:"photos" json:"photos"`
	Paid       bool               `bson:"paid" json:"paid"`
	PublishedAt *time.Time        `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	ExpiresAt  *time.Time         `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	DNS        *DNSRecord         `bson:"dns,omitempty" json:"dns,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type OrderStatus string

const (
	OrderPending OrderStatus = "pending"
	OrderPaid    OrderStatus = "paid"
	OrderFailed  OrderStatus = "failed"
)

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID     primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	GreetingID  primitive.ObjectID `bson:"greetingId" json:"greetingId"`
	ProductType GreetingType       `bson:"productType" json:"productType"`
	AmountRUB   int                `bson:"amountRUB" json:"amountRUB"`
	Currency    string             `bson:"currency" json:"currency"`
	Status      OrderStatus        `bson:"status" json:"status"`
	Provider    string             `bson:"provider" json:"provider"`
	ProviderRef string             `bson:"providerRef,omitempty" json:"providerRef,omitempty"`
	PaidAt      *time.Time         `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}
