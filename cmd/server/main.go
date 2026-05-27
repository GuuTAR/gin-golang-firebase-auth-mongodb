package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/config"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/handlers"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/pkg/db"
)

func main() {
	// Load .env for local development; ignore error when file is absent.
	_ = godotenv.Load()

	cfg := config.Load()

	// ── Process-lifetime context (cancelled on SIGINT / SIGTERM) ─────────────
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ── MongoDB ──────────────────────────────────────────────────────────────
	// ctx is signal-aware: if SIGTERM arrives during connect/ping it cancels
	// immediately rather than waiting out the internal 10 s timeout.
	mongoClient, err := db.NewMongoClient(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}
	log.Printf("mongodb connected (db=%s)", cfg.MongoDB)

	// ── Dependency wiring ────────────────────────────────────────────────────
	healthHandler := handlers.NewHealthHandler()

	// ── Router ───────────────────────────────────────────────────────────────
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/health", healthHandler.Check)

	v1 := r.Group("/api/v1")
	_ = v1 // placeholder until routes are added

	// ── HTTP server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	<-ctx.Done() // block until signal

	// Release signal resources so a second signal kills the process immediately.
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		log.Printf("mongodb disconnect: %v", err)
	}
	log.Println("server exited")
}
