package jobs

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"happy-api/internal/config"
	"happy-api/internal/models"
	"happy-api/internal/repo"
)

type Cleaner struct {
	cfg config.Config
	r   *repo.Repo
}

func NewCleaner(cfg config.Config, r *repo.Repo) *Cleaner {
	return &Cleaner{cfg: cfg, r: r}
}

func (c *Cleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// first run immediately
	c.sweep(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.sweep(ctx)
		}
	}
}

func (c *Cleaner) sweep(ctx context.Context) {
	items, err := c.r.FindExpiredGreetings(ctx, 100)
	if err != nil {
		log.Println("[cleaner] find expired:", err)
		return
	}
	for _, g := range items {
		if err := c.cleanupOne(ctx, g); err != nil {
			log.Println("[cleaner] cleanup one:", err)
		}
	}
}

func (c *Cleaner) cleanupOne(ctx context.Context, g models.Greeting) error {

	// delete uploads dir
	if c.cfg.UploadDir != "" {
		_ = os.RemoveAll(filepath.Join(c.cfg.UploadDir, g.ID.Hex()))
	}

	// delete related orders
	_ = c.r.DeleteOrdersByGreeting(ctx, g.ID)

	// delete greeting itself
	return c.r.DeleteGreeting(ctx, g.ID)
}
