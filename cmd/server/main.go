package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/config"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/handlers"
)

func main() {
	// Load .env for local development; ignore error when file is absent.
	_ = godotenv.Load()

	cfg := config.Load()

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

	// ── Graceful shutdown ────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the HTTP server in a goroutine so it doesn't block the main thread,
	// allowing the signal listener below to run concurrently.
	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.Port, cfg.Env)
		// ListenAndServe always returns a non-nil error; ignore ErrServerClosed
		// because that is the expected result of a clean Shutdown call.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Block until an OS interrupt or termination signal is received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Give in-flight requests up to 5 seconds to complete before forcing exit.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("server exited")
}
