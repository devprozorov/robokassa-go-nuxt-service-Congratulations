package payments

import (
	"strings"

	"happy-api/internal/config"
	"happy-api/internal/payments/robokassa"
)

func NewProvider(cfg config.Config) Provider {
	switch strings.ToLower(strings.TrimSpace(cfg.PaymentProvider)) {
	case "robokassa":
		return robokassa.New(cfg)
	default:
		return Stub{BaseURL: "/dashboard/orders"}
	}
}
