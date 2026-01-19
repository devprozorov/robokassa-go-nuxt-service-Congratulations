package payments

import (
	"context"
	"fmt"
	"happy-api/internal/config"
	"happy-api/internal/models"
	"happy-api/internal/payments/robokassa"
)

// Provider abstracts payment gateway.
// Use stub for development; implement real providers (YooKassa/CloudPayments/Stripe) by adding another provider.
type Provider interface {
	Name() string
	CreateCheckout(ctx context.Context, order models.Order) (checkoutURL string, providerRef string, err error)
	HandleWebhook(ctx context.Context, body []byte, headers map[string]string) (providerRef string, paid bool, err error)
}

func FromConfig(cfg config.Config) Provider {
	switch cfg.PaymentProvider {
	case "robokassa":
		return robokassa.New(cfg)
	default:
		return Stub{}
	}
}

type Stub struct {
	BaseURL string
}

func (s Stub) Name() string { return "stub" }

func (s Stub) CreateCheckout(ctx context.Context, order models.Order) (string, string, error) {
	// Redirect to internal order status page.
	return fmt.Sprintf("%s/%s", s.BaseURL, order.ID.Hex()), "", nil
}

func (s Stub) HandleWebhook(ctx context.Context, body []byte, headers map[string]string) (string, bool, error) {
	// no-op
	return "", false, nil
}
