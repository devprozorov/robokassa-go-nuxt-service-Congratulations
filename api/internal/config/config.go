package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv       string
	MongoURI     string
	MongoDB      string
	JWTSecret    string
	CookieDomain string
	UploadDir    string

	PublishTTL time.Duration

	PriceCodeRUB int
	PriceSubRUB  int

	AdminEmail    string
	AdminUsername string
	AdminPassword string

	// Payments
	PaymentProvider string

	// Robokassa
	RoboMerchantLogin string
	RoboPass1         string
	RoboPass2         string
	RoboIsTest        bool
	RoboPayURL        string
}

func MustLoad() Config {
	cfg := Config{
		AppEnv:       os.Getenv("APP_ENV"),
		MongoURI:     getenv("MONGO_URI", "mongodb://mongo:27017"),
		MongoDB:      getenv("MONGO_DB", "happy"),
		JWTSecret:    getenv("JWT_SECRET", "change-me-super-secret"),
		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
		UploadDir:    getenv("UPLOAD_DIR", "/data/uploads"),
		PublishTTL:   parseDuration(getenv("PUBLISH_TTL", "168h")),

		PriceCodeRUB: parseInt(getenv("PRICE_CODE_RUB", "59")),
		PriceSubRUB:  parseInt(getenv("PRICE_SUB_RUB", "99")),

		AdminEmail:    getenv("ADMIN_EMAIL", "admin@local"),
		AdminUsername: getenv("ADMIN_USERNAME", "admin"),
		AdminPassword: getenv("ADMIN_PASSWORD", "admin123456"),

		PaymentProvider: getenv("PAYMENT_PROVIDER", "robokassa"),

		RoboMerchantLogin: os.Getenv("ROBO_MERCHANT_LOGIN"),
		RoboPass1:         os.Getenv("ROBO_PASS1"),
		RoboPass2:         os.Getenv("ROBO_PASS2"),
		RoboPayURL:        getenv("ROBO_PAY_URL", "https://auth.robokassa.ru/Merchant/Index.aspx"),
		RoboIsTest:        parseBool(getenv("ROBO_IS_TEST", "0")),
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "dev"
	}
	if len(cfg.JWTSecret) < 16 {
		log.Println("[warn] JWT_SECRET is short; set a long random value in production")
	}

	return cfg
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func parseInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

func parseDuration(v string) time.Duration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0
	}
	return d
}

func parseBool(v string) bool {
	s := strings.ToLower(strings.TrimSpace(v))
	return s == "1" || s == "true" || s == "yes" || s == "y" || s == "on"
}
