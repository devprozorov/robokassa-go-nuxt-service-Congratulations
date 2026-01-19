package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"happy-api/internal/config"
	"happy-api/internal/db"
	httpapi "happy-api/internal/http"
	"happy-api/internal/jobs"
	"happy-api/internal/payments"
	"happy-api/internal/repo"
)

func main() {
	mode := flag.String("mode", "server", "server|worker")
	flag.Parse()

	cfg := config.MustLoad()

	ctx := context.Background()
	d, err := db.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = d.Client.Disconnect(ctx) }()

	r := repo.New(d.DB)
	if err := r.InitIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	if err := r.EnsureAdmin(ctx, cfg.AdminEmail, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Fatal(err)
	}

	pay := payments.FromConfig(cfg)

	if *mode == "worker" {
		log.Println("Happy worker started")
		cleaner := jobs.NewCleaner(cfg, r)
		cleaner.Run(ctx)
		return
	}

	router := httpapi.NewRouter(cfg, r, pay)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("Happy API listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
