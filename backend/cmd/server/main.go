package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"abberant/backend/internal/api"
	"abberant/backend/internal/articles"
	"abberant/backend/internal/auth"
	"abberant/backend/internal/config"
	"abberant/backend/internal/db"
	"abberant/backend/internal/mailer"
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

	app := api.New(cfg, service, sessions, articles.NewRepo(pool))

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
