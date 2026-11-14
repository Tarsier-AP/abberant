package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"abberant/backend/internal/api"
	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
	"abberant/backend/internal/bookmarks"
	"abberant/backend/internal/comments"
	"abberant/backend/internal/config"
	"abberant/backend/internal/db"
	"abberant/backend/internal/mailer"
	"abberant/backend/internal/media"
	"abberant/backend/internal/notifications"
	"abberant/backend/internal/ratings"
	"abberant/backend/internal/sections"
	"abberant/backend/internal/session"
)

func main() {
	cfg := config.Load()

	pool, close, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	defer close()

	sessions := session.NewStore(pool, cfg.SessionTTL, cfg.PendingTTL)
	service := auth.NewService(pool, sessions, mailer.New(cfg), cfg)

	mediaStore := media.NewStore(cfg.MediaDir)

	// Data-URL and seed-image migrations plus a daily garbage-collection
	// sweep run alongside the server.
	go media.RunMaintenance(context.Background(), pool, mediaStore, cfg.SeedImagesDir)

	app := api.New(
		cfg,
		service,
		sessions,
		articles.NewRepo(pool),
		comments.NewRepo(pool),
		bookmarks.NewRepo(pool),
		notifications.NewRepo(pool),
		ratings.NewRepo(pool),
		sections.NewRepo(pool),
		mediaStore,
	)

	go func() {
		log.Printf("listening on %s", cfg.Addr)

		if err := app.Listen(cfg.Addr); err != nil {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("shutting down")

	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
