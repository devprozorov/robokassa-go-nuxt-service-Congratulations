package httpapi

import (
	"happy-api/internal/config"
	"happy-api/internal/http/handlers"
	"happy-api/internal/http/mw"
	"happy-api/internal/payments"
	"happy-api/internal/payments/robokassa"
	"happy-api/internal/repo"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, r *repo.Repo, pay payments.Provider) *gin.Engine {
	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.Logger())
	e.MaxMultipartMemory = 100 << 20 // 100MB

	h := handlers.New(cfg, r, pay)

	// robokassa client
	rb := &robokassa.Client{
		MerchantLogin: cfg.RoboMerchantLogin,
		Pass1:         cfg.RoboPass1,
		Pass2:         cfg.RoboPass2,
		IsTest:        cfg.RoboIsTest,
		PayURL:        cfg.RoboPayURL,
	}
	h.SetRobokassa(rb)

	e.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	api := e.Group("/api")
	{

		api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
		// Auth
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
		api.POST("/auth/logout", mw.Auth(cfg, r), h.Logout)
		api.GET("/auth/me", mw.Auth(cfg, r), h.Me)
		api.POST("/auth/recover/reset", h.RecoverReset)
		api.POST("/auth/password/change", mw.Auth(cfg, r), h.ChangePassword)

		// Public greetings
		api.GET("/public/code/:code", h.PublicByCode)
		api.GET("/public/subdomain/:sub", h.PublicBySubdomain)

		// Media
		api.GET("/media/:gid/:filename", h.Media)

		// Webhooks (optional)
		api.POST("/webhooks/payment", h.PaymentWebhook) // ResultURL робокассы

		// webhook публичный (ResultURL)
		api.POST("/webhooks/robokassa/result", h.RoboResultWebhook)
		api.GET("/webhooks/robokassa/result", h.RoboResultWebhook)

		// Protected: greetings, orders
		protected := api.Group("")
		protected.Use(mw.Auth(cfg, r))
		{

			protected.POST("/greetings", h.CreateGreeting)
			protected.GET("/greetings", h.ListGreetings)
			protected.GET("/greetings/:id", h.GetGreeting)
			protected.PUT("/greetings/:id", h.UpdateGreeting)
			protected.DELETE("/greetings/:id", h.DeleteGreeting)

			protected.POST("/greetings/:id/photos", h.UploadPhotos)
			protected.POST("/greetings/:id/publish", h.PublishGreeting)

			protected.GET("/orders/:id", h.GetOrder)

			protected.POST("/payments/robokassa/init", h.RoboInitPayment)

		}

		// Admin
		admin := api.Group("/admin")
		admin.Use(mw.Auth(cfg, r), mw.AdminOnly())
		{
			admin.GET("/users", h.AdminUsers)
			admin.DELETE("/users/:id", h.AdminDeleteUser)

			admin.GET("/orders", h.AdminOrders)
			admin.POST("/orders/:id/mark-paid", h.AdminMarkPaid)

			admin.GET("/greetings", h.AdminGreetings)
			admin.DELETE("/greetings/:id", h.AdminDeleteGreeting)

			admin.GET("/stats", h.AdminStats)
		}
	}

	// Fallback
	e.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"ok": false, "error": "NOT_FOUND"})
	})
	return e
}
